package userdata

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-redis/redis"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"

	exceptions "github.com/HiIamJeff67/notezy-backend/contracts/types/exceptions"
	logs "github.com/HiIamJeff67/notezy-backend/shared/platform/observability/logs"
	platformredis "github.com/HiIamJeff67/notezy-backend/shared/platform/redis"

	coreconfig "github.com/HiIamJeff67/notezy-backend/internal/core/configs"
	cacheinputs "github.com/HiIamJeff67/notezy-backend/internal/core/data/cache/userdata/inputs"
	redislibraries "github.com/HiIamJeff67/notezy-backend/internal/core/data/cache/userdata/libraries"
	enums "github.com/HiIamJeff67/notezy-backend/internal/core/data/database/schemas/enums"
)

type UserDataCache struct {
	Id                 uuid.UUID        `json:"id"`
	PublicId           uuid.UUID        `json:"publicId"`
	Name               string           `json:"name"`
	DisplayName        string           `json:"displayName"`
	Email              string           `json:"email"`
	AccessToken        string           `json:"accessToken"`
	CSRFToken          string           `json:"csrfToken"`
	Role               enums.UserRole   `json:"role"`
	Plan               enums.UserPlan   `json:"plan"`
	Status             enums.UserStatus `json:"status"`
	AvatarURL          string           `json:"avatarURL"`
	Language           enums.Language   `json:"language"`
	GeneralSettingCode int64            `json:"generalSettingCode"`
	PrivacySettingCode int64            `json:"privacySettingCode"`
	CreatedAt          time.Time        `json:"createdAt"`
	UpdatedAt          time.Time        `json:"updatedAt"`
}

type UserDataCacheClient struct {
	cacheStore *UserDataCacheStore

	cacheExpiresIn                                       time.Duration
	batchCheckAndUpdateQuotasByFormattedKeysArgvPerKey   int
	batchCheckAndUpdateQuotasByFormattedKeyBaseNumOfArgv int
}

/* ============================== Constructor ============================== */

func NewUserDataCacheClient(config coreconfig.UserDataCacheConfig, cacheStore *UserDataCacheStore) *UserDataCacheClient {
	return &UserDataCacheClient{
		cacheStore: cacheStore,

		cacheExpiresIn: config.CacheExpiresIn,
		batchCheckAndUpdateQuotasByFormattedKeysArgvPerKey:   4,
		batchCheckAndUpdateQuotasByFormattedKeyBaseNumOfArgv: 4,
	}
}

/* ============================== Auxiliary Methods ============================== */

func (s *UserDataCacheClient) getRedisClient(identifier string) (*redis.Client, int, *exceptions.Exception) {
	if s == nil || s.cacheStore == nil {
		return nil, 0, exceptions.New(
			"CacheClientUnavailable",
			"Cache",
			"GetRedisClient",
			"User data cache client is unavailable",
			http.StatusInternalServerError,
			true,
		)
	}
	redisClient, shardIndex, err := s.cacheStore.ClientSet().ClientForKey(identifier)
	if err != nil {
		return nil, 0, exceptions.New(
			"CacheClientUnavailable",
			"Cache",
			"GetRedisClient",
			"User data cache client is unavailable",
			http.StatusInternalServerError,
			true,
		).WithOrigin(err)
	}
	return redisClient, shardIndex, nil
}

func formatUserDataKey(identifier string) string {
	return fmt.Sprintf("%s:%s", platformredis.CachePurpose_UserData.String(), identifier)
}

/* ============================== Extend Methods ============================== */

func (s *UserDataCacheClient) Extend(identifier string) *exceptions.Exception {
	redisClient, _, exception := s.getRedisClient(identifier)
	if exception != nil {
		return exception
	}

	updated, err := redisClient.Expire(formatUserDataKey(identifier), s.cacheExpiresIn).Result()
	if err != nil {
		return exceptions.New(
			"FailedToExtendTTL",
			"Cache",
			"ExtendUserData",
			"Failed to extend cached user data",
			http.StatusInternalServerError,
			true,
		).WithOrigin(err)
	}
	if !updated {
		return exceptions.New(
			"NotFound",
			"Cache",
			"ExtendUserData",
			"Cached user data was not found",
			http.StatusNotFound,
			true,
		)
	}

	return nil
}

/* ============================== Quota Method ============================== */

func (s *UserDataCacheClient) CheckAndUpdateQuota(
	identifier string,
	input cacheinputs.CheckAndUpdateUserQuotaInput,
) *exceptions.Exception {
	redisClient, _, exception := s.getRedisClient(identifier)
	if exception != nil {
		return exception
	}

	arguments := []interface{}{
		"FCALL",
		redislibraries.CheckAndUpdateUserQuotaByFormattedKeyFunction,
		1,
		formatUserDataKey(identifier),
		input.Field,
		input.ChangeAmount,
		input.MaxLimit,
		int(time.Until(input.ExpiresIn).Seconds()),
	}
	if _, err := redisClient.Do(arguments...).Result(); err != nil {
		return exceptions.New(
			"FailedToUpdate",
			"Cache",
			"CheckAndUpdateUserQuota",
			"Failed to update the user quota cache",
			http.StatusInternalServerError,
			true,
		).WithOrigin(err)
	}

	return nil
}

func (s *UserDataCacheClient) BestEffortBatchCheckAndUpdateQuotas(
	inputs []cacheinputs.BatchCheckAndUpdateUserQuotaInput,
) *exceptions.Exception {
	if len(inputs) == 0 {
		return nil
	}

	inputsByShard := make(map[int][]cacheinputs.BatchCheckAndUpdateUserQuotaInput)
	for _, input := range inputs {
		_, shardIndex, err := s.cacheStore.ClientSet().ClientForKey(input.Identifier)
		if err != nil {
			continue
		}
		inputsByShard[shardIndex] = append(inputsByShard[shardIndex], input)
	}

	for shardIndex, groupedInputs := range inputsByShard {
		redisClient, exists := s.cacheStore.ClientSet().Client(shardIndex)
		if !exists {
			continue
		}

		keys := make([]interface{}, 0, len(groupedInputs))
		arguments := make([]interface{}, 0, len(groupedInputs)*s.batchCheckAndUpdateQuotasByFormattedKeysArgvPerKey)
		for _, input := range groupedInputs {
			keys = append(keys, formatUserDataKey(input.Identifier))
			arguments = append(arguments,
				input.Input.Field,
				input.Input.ChangeAmount,
				input.Input.MaxLimit,
				int(time.Until(input.Input.ExpiresIn).Seconds()),
			)
		}

		command := []interface{}{
			"FCALL",
			redislibraries.BestEffortBatchCheckAndUpdateUserQuotasByFormattedKeysFunction,
			len(keys),
		}
		command = append(command, keys...)
		command = append(command, arguments...)
		if _, err := redisClient.Do(command...).Result(); err != nil {
			return exceptions.New(
				"FailedToUpdate",
				"Cache",
				"BestEffortBatchCheckAndUpdateUserQuotas",
				"Failed to update the user quota cache",
				http.StatusInternalServerError,
				true,
			).WithOrigin(err)
		}
	}

	return nil
}

func (s *UserDataCacheClient) BestEffortBatchCheckAndUpdateQuotasByIdentifier(
	identifier string,
	inputs []cacheinputs.CheckAndUpdateUserQuotaInput,
) *exceptions.Exception {
	if len(inputs) == 0 {
		return nil
	}

	redisClient, _, exception := s.getRedisClient(identifier)
	if exception != nil {
		return exception
	}

	arguments := make([]interface{}, 0, len(inputs)*s.batchCheckAndUpdateQuotasByFormattedKeyBaseNumOfArgv)
	for _, input := range inputs {
		arguments = append(arguments,
			input.Field,
			input.ChangeAmount,
			input.MaxLimit,
			int(time.Until(input.ExpiresIn).Seconds()),
		)
	}

	command := []interface{}{
		"FCALL",
		redislibraries.BestEffortBatchCheckAndUpdateUserQuotasByFormattedKeyFunction,
		1,
		formatUserDataKey(identifier),
	}
	command = append(command, arguments...)
	if _, err := redisClient.Do(command...).Result(); err != nil {
		return exceptions.New(
			"FailedToUpdate",
			"Cache",
			"BestEffortBatchCheckAndUpdateUserQuotasByIdentifier",
			"Failed to update the user quota cache",
			http.StatusInternalServerError,
			true,
		).WithOrigin(err)
	}

	return nil
}

/* ============================== CRUD Method ============================== */

func (s *UserDataCacheClient) Get(identifier string) (*UserDataCache, *exceptions.Exception) {
	redisClient, serverNumber, exception := s.getRedisClient(identifier)
	if exception != nil {
		return nil, exception
	}

	cacheString, err := redisClient.Get(formatUserDataKey(identifier)).Result()
	if err != nil {
		return nil, exceptions.New(
			"NotFound",
			"Cache",
			"GetUserData",
			"Cached user data was not found",
			http.StatusNotFound,
			true,
		).WithOrigin(err)
	}

	var userDataCache UserDataCache
	if err := json.Unmarshal([]byte(cacheString), &userDataCache); err != nil {
		return nil, exceptions.New(
			"DeserializationFailed",
			"Cache",
			"GetUserData",
			"Failed to decode cached user data",
			http.StatusInternalServerError,
			true,
		).WithOrigin(err)
	}

	logs.NotezyLogger.Debug(context.Background(), fmt.Sprintf("Successfully got cached user data from server %d", serverNumber))
	return &userDataCache, nil
}

func (s *UserDataCacheClient) Set(identifier string, userDataCache UserDataCache) *exceptions.Exception {
	if userDataCache.PublicId == uuid.Nil ||
		strings.TrimSpace(userDataCache.Name) == "" ||
		strings.TrimSpace(userDataCache.DisplayName) == "" ||
		strings.TrimSpace(userDataCache.Email) == "" ||
		strings.TrimSpace(userDataCache.AccessToken) == "" ||
		!userDataCache.Role.IsValidEnum() ||
		!userDataCache.Plan.IsValidEnum() ||
		!userDataCache.Status.IsValidEnum() {
		return exceptions.New(
			"InvalidCacheData",
			"Cache",
			"SetUserData",
			"Cached user data is invalid",
			http.StatusInternalServerError,
			true,
		)
	}

	redisClient, serverNumber, exception := s.getRedisClient(identifier)
	if exception != nil {
		return exception
	}

	value, err := json.Marshal(userDataCache)
	if err != nil {
		return exceptions.New(
			"SerializationFailed",
			"Cache",
			"SetUserData",
			"Failed to encode cached user data",
			http.StatusInternalServerError,
			true,
		).WithOrigin(err)
	}

	if err := redisClient.Set(formatUserDataKey(identifier), string(value), s.cacheExpiresIn).Err(); err != nil {
		return exceptions.New(
			"FailedToCreate",
			"Cache",
			"SetUserData",
			"Failed to store cached user data",
			http.StatusInternalServerError,
			true,
		).WithOrigin(err)
	}

	logs.NotezyLogger.Debug(context.Background(), fmt.Sprintf("Successfully set cached user data in server %d", serverNumber))
	return nil
}

func (s *UserDataCacheClient) Update(identifier string, input cacheinputs.UpdateUserDataCacheInput) *exceptions.Exception {
	userDataCache, exception := s.Get(identifier)
	if exception != nil {
		return exception
	}

	userDataCache.UpdatedAt = time.Now()
	if err := copier.Copy(userDataCache, &input); err != nil {
		return exceptions.New(
			"SerializationFailed",
			"Cache",
			"UpdateUserData",
			"Failed to copy cached user data",
			http.StatusInternalServerError,
			true,
		).WithOrigin(err)
	}

	redisClient, serverNumber, exception := s.getRedisClient(identifier)
	if exception != nil {
		return exception
	}

	value, err := json.Marshal(userDataCache)
	if err != nil {
		return exceptions.New(
			"SerializationFailed",
			"Cache",
			"UpdateUserData",
			"Failed to encode cached user data",
			http.StatusInternalServerError,
			true,
		).WithOrigin(err)
	}

	if err := redisClient.Set(formatUserDataKey(identifier), string(value), s.cacheExpiresIn).Err(); err != nil {
		return exceptions.New(
			"FailedToUpdate",
			"Cache",
			"UpdateUserData",
			"Failed to update cached user data",
			http.StatusInternalServerError,
			true,
		).WithOrigin(err)
	}

	logs.NotezyLogger.Debug(context.Background(), fmt.Sprintf("Successfully updated cached user data in server %d", serverNumber))
	return nil
}

func (s *UserDataCacheClient) Delete(identifier string) *exceptions.Exception {
	redisClient, serverNumber, exception := s.getRedisClient(identifier)
	if exception != nil {
		return exception
	}

	if err := redisClient.Del(formatUserDataKey(identifier)).Err(); err != nil {
		return exceptions.New(
			"FailedToDelete",
			"Cache",
			"DeleteUserData",
			"Failed to delete cached user data",
			http.StatusInternalServerError,
			true,
		).WithOrigin(err)
	}

	logs.NotezyLogger.Debug(context.Background(), fmt.Sprintf("Successfully deleted cached user data from server %d", serverNumber))
	return nil
}
