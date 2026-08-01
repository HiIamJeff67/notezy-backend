package caches

import (
	"context"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"net/http"
	"strings"
	"time"

	"github.com/go-redis/redis"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"

	cacheinputs "github.com/HiIamJeff67/notezy-backend/internal/caches/inputs"
	redislibraries "github.com/HiIamJeff67/notezy-backend/internal/caches/libraries"
	exceptions "github.com/HiIamJeff67/notezy-backend/internal/exceptions"
	logs "github.com/HiIamJeff67/notezy-backend/internal/platform/observability/logs"
	enums "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/schemas/enums"
	types "github.com/HiIamJeff67/notezy-backend/internal/shared/types"
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

type UserDataCacheStore struct {
	redisClientMap map[int]*redis.Client

	Range           types.Range[int, int]
	MaxServerNumber int

	cacheExpiresIn                                       time.Duration
	batchCheckAndUpdateQuotasByFormattedKeysArgvPerKey   int
	batchCheckAndUpdateQuotasByFormattedKeyBaseNumOfArgv int
}

/* ============================== Constructor ============================== */

func NewUserDataCacheStore(redisClientMap map[int]*redis.Client) *UserDataCacheStore {
	rangeValue := types.Range[int, int]{Start: 0, Size: 4}

	return &UserDataCacheStore{
		redisClientMap: redisClientMap,

		Range:           rangeValue,
		MaxServerNumber: rangeValue.Start + rangeValue.Size - 1,

		cacheExpiresIn: time.Hour,
		batchCheckAndUpdateQuotasByFormattedKeysArgvPerKey:   4,
		batchCheckAndUpdateQuotasByFormattedKeyBaseNumOfArgv: 4,
	}
}

/* ============================== Auxiliary Methods ============================== */

func (s *UserDataCacheStore) getRedisClient(identifier string) (*redis.Client, int, *exceptions.Exception) {
	hash := s.hashIdentifier(identifier)
	serverNumber := min(s.MaxServerNumber, s.Range.Start+hash)
	redisClient, ok := s.redisClientMap[serverNumber]
	if !ok || redisClient == nil {
		return nil, 0, exceptions.New(
			"CacheClientUnavailable",
			"Cache",
			"GetRedisClient",
			"User data cache client is unavailable",
			http.StatusInternalServerError,
			true,
		)
	}

	return redisClient, serverNumber, nil
}

func (s *UserDataCacheStore) hashIdentifier(identifier string) int {
	hash := fnv.New32a()
	_, _ = hash.Write([]byte(identifier))

	return int(hash.Sum32()) % s.Range.Size
}

func formatUserDataKey(identifier string) string {
	return fmt.Sprintf("%s:%s", types.ValidCachePurpose_UserData.String(), identifier)
}

func isValidUserDataCache(userDataCache *UserDataCache) bool {
	return userDataCache.PublicId != uuid.Nil &&
		strings.TrimSpace(userDataCache.Name) != "" &&
		strings.TrimSpace(userDataCache.DisplayName) != "" &&
		strings.TrimSpace(userDataCache.Email) != "" &&
		strings.TrimSpace(userDataCache.AccessToken) != "" &&
		userDataCache.Role.IsValidEnum() &&
		userDataCache.Plan.IsValidEnum() &&
		userDataCache.Status.IsValidEnum()
}

/* ============================== Extend Methods ============================== */

func (s *UserDataCacheStore) Extend(identifier string) *exceptions.Exception {
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

func (s *UserDataCacheStore) CheckAndUpdateQuota(
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

func (s *UserDataCacheStore) BestEffortBatchCheckAndUpdateQuotas(
	inputs []cacheinputs.BatchCheckAndUpdateUserQuotaInput,
) *exceptions.Exception {
	if len(inputs) == 0 {
		return nil
	}

	inputsByServerNumber := make(map[int][]cacheinputs.BatchCheckAndUpdateUserQuotaInput)
	for _, input := range inputs {
		hash := s.hashIdentifier(input.Identifier)
		serverNumber := min(s.MaxServerNumber, s.Range.Start+hash)
		inputsByServerNumber[serverNumber] = append(inputsByServerNumber[serverNumber], input)
	}

	for serverNumber, groupedInputs := range inputsByServerNumber {
		redisClient, ok := s.redisClientMap[serverNumber]
		if !ok || redisClient == nil {
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

func (s *UserDataCacheStore) BestEffortBatchCheckAndUpdateQuotasByIdentifier(
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

func (s *UserDataCacheStore) Get(identifier string) (*UserDataCache, *exceptions.Exception) {
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

func (s *UserDataCacheStore) Set(identifier string, userDataCache UserDataCache) *exceptions.Exception {
	if !isValidUserDataCache(&userDataCache) {
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

func (s *UserDataCacheStore) Update(identifier string, input cacheinputs.UpdateUserDataCacheInput) *exceptions.Exception {
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

func (s *UserDataCacheStore) Delete(identifier string) *exceptions.Exception {
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
