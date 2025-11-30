package caches

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"strings"
	"time"

	uuid "github.com/google/uuid"
	"github.com/jinzhu/copier"

	redislibraries "notezy-backend/app/caches/libraries"
	exceptions "notezy-backend/app/exceptions"
	logs "notezy-backend/app/logs"
	enums "notezy-backend/app/models/schemas/enums"
	types "notezy-backend/shared/types"
)

type UserDataCache struct {
	PublicId           string           `json:"publicId"`           // user
	Name               string           `json:"name"`               // user
	DisplayName        string           `json:"displayName"`        // user
	Email              string           `json:"email"`              // user
	AccessToken        string           `json:"accessToken"`        // !only here: note that it may be expired, but it is the newest one
	CSRFToken          string           `json:"csrfToken"`          // !only here: note that it may be expired, but it is the newest one
	Role               enums.UserRole   `json:"role"`               // user
	Plan               enums.UserPlan   `json:"plan"`               // user
	Status             enums.UserStatus `json:"status"`             // user
	AvatarURL          string           `json:"avatarURL"`          // user info
	Language           enums.Language   `json:"language"`           // user setting
	GeneralSettingCode int64            `json:"generalSettingCode"` // user setting
	PrivacySettingCode int64            `json:"privacySettingCode"` // user setting
	CreatedAt          time.Time        `json:"createdAt"`          // user
	UpdatedAt          time.Time        `json:"updatedAt"`          // user
}

type UpdateUserDataCacheDto struct {
	PublicId           *string
	Name               *string
	DisplayName        *string
	Email              *string
	AccessToken        *string
	CSRFToken          *string
	Role               *enums.UserRole
	Plan               *enums.UserPlan
	Status             *enums.UserStatus
	AvatarURL          *string
	Language           *enums.Language
	GeneralSettingCode *int64
	PrivacySettingCode *int64
}

type CheckAndUpdateUserQuotaDto struct {
	Field        types.UserQuotaField
	ChangeAmount int32
	MaxLimit     int32
	ExpiresIn    time.Time
}

const (
	_userDataCacheExpiresIn = 1 * time.Hour

	batchCheckAndUpdateUserQuotasByFormattedKeysArgvPerKey   = 4
	batchCheckAndUpdateUserQuotasByFormattedKeyBaseNumOfArgv = 4
)

var (
	UserDataRange           = types.Range[int, int]{Start: 0, Size: 4} // server number: 0 - 3 (included)
	MaxUserDataServerNumber = UserDataRange.Start + UserDataRange.Size - 1
)

/* ========================= Auxiliary Function ========================= */

func hashUserDataIdentifier(identifier uuid.UUID) int {
	h := fnv.New32a()
	h.Write([]byte(identifier.String()))
	return int(h.Sum32()) % UserDataRange.Size
}

func formatUserDataKey(id uuid.UUID) string {
	return fmt.Sprintf("%s:%s", types.ValidCachePurpose_UserData.String(), id.String())
}

func isValidUserCacheData(userDataCache *UserDataCache) bool {
	if strings.ReplaceAll(userDataCache.PublicId, " ", "") == "" ||
		strings.ReplaceAll(userDataCache.Name, " ", "") == "" ||
		strings.ReplaceAll(userDataCache.DisplayName, " ", "") == "" ||
		strings.ReplaceAll(userDataCache.Email, " ", "") == "" ||
		strings.ReplaceAll(userDataCache.AccessToken, " ", "") == "" ||
		!enums.IsValidEnumValues(userDataCache.Role, enums.AllUserRoles) ||
		!enums.IsValidEnumValues(userDataCache.Plan, enums.AllUserPlans) ||
		!enums.IsValidEnumValues(userDataCache.Status, enums.AllUserStatuses) {
		return false
	}
	return true
}

/* ============================== Extend Cache TTL Operation ============================== */

func ExtendUserDataCacheTTL(id uuid.UUID) *exceptions.Exception {
	hash := hashUserDataIdentifier(id)
	serverNumber := min(MaxUserDataServerNumber, UserDataRange.Start+hash)
	redisClient, ok := RedisClientMap[serverNumber]
	if !ok {
		return exceptions.Cache.ClientInstanceDoesNotExist()
	}

	formattedKey := formatUserDataKey(id)

	updated, err := redisClient.Expire(formattedKey, _userDataCacheExpiresIn).Result()
	if err != nil {
		return exceptions.Cache.FailedToUpdate("UserDataTTL").WithError(err)
	}

	if !updated {
		return exceptions.Cache.NotFound(string(types.ValidCachePurpose_UserData))
	}

	return nil
}

/* ============================== Automatic Operations for Accouting ============================== */

func CheckAndUpdateUserQuotaByFormattedKey(
	id uuid.UUID,
	dto CheckAndUpdateUserQuotaDto,
) *exceptions.Exception {
	hash := hashUserDataIdentifier(id)
	serverNumber := min(MaxUserDataServerNumber, UserDataRange.Start+hash)
	redisClient, ok := RedisClientMap[serverNumber]
	if !ok {
		return exceptions.Cache.ClientInstanceDoesNotExist()
	}

	formattedKey := formatUserDataKey(id)

	keys := []interface{}{formattedKey}
	argv := []interface{}{
		dto.Field,
		dto.ChangeAmount,
		dto.MaxLimit,
		int(time.Until(dto.ExpiresIn).Seconds()),
	}

	arguments := []interface{}{
		"FCALL",
		redislibraries.CheckAndUpdateUserQuotaByFormattedKeyFunction,
		len(keys),
	}
	arguments = append(arguments, keys...)
	arguments = append(arguments, argv...)
	if _, err := redisClient.Do(arguments...).Result(); err != nil {
		return exceptions.Cache.FailedToUpdate(types.ValidCachePurpose_UserData.String()).WithError(err)
	}

	return nil
}

func BestEffortBatchCheckAndUpdateUserQuotasByFormattedKeys(
	dtos []struct {
		Id                uuid.UUID                  `json:"id"`
		CheckAndUpdateDto CheckAndUpdateUserQuotaDto `json:"checkAndUpdateDto"`
	},
) *exceptions.Exception {
	if len(dtos) == 0 {
		return nil
	}
	serverNumberToUserIdMap := make(map[int][]struct {
		Id                uuid.UUID                  `json:"id"`
		CheckAndUpdateDto CheckAndUpdateUserQuotaDto `json:"checkAndUpdateDto"`
	})

	for _, dto := range dtos {
		hash := hashUserDataIdentifier(dto.Id)
		serverNumber := min(MaxUserDataServerNumber, UserDataRange.Start+hash)
		serverNumberToUserIdMap[serverNumber] = append(serverNumberToUserIdMap[serverNumber], dto)
	}

	for serverNumber, dtos := range serverNumberToUserIdMap {
		redisClient, exist := RedisClientMap[serverNumber]
		if !exist {
			continue // for the strategy of "Best Effort"
		}

		keys := make([]interface{}, 0, len(dtos))
		argv := make([]interface{}, 0, len(dtos)*batchCheckAndUpdateUserQuotasByFormattedKeysArgvPerKey)
		for _, dto := range dtos {
			keys = append(keys, dto.Id)
			argv = append(argv, dto.CheckAndUpdateDto.Field)
			argv = append(argv, dto.CheckAndUpdateDto.ChangeAmount)
			argv = append(argv, dto.CheckAndUpdateDto.MaxLimit)
			argv = append(argv, int(time.Until(dto.CheckAndUpdateDto.ExpiresIn).Seconds()))
		}

		arguments := []interface{}{
			"FCALL",
			redislibraries.BestEffortBatchCheckAndUpdateUserQuotasByFormattedKeysFunction,
			len(dtos),
		}
		arguments = append(arguments, keys...)
		arguments = append(arguments, argv...)
		if _, err := redisClient.Do(arguments...).Result(); err != nil {
			return exceptions.Cache.FailedToDelete(types.ValidCachePurpose_UserData.String()).WithError(err)
		}
	}

	return nil
}

func BestEffortBatchCheckAndUpdateUserQuotasByFormattedKey(
	id uuid.UUID,
	dtos []CheckAndUpdateUserQuotaDto,
) *exceptions.Exception {
	if len(dtos) == 0 {
		return nil
	}

	hash := hashUserDataIdentifier(id)
	serverNumber := min(MaxUserDataServerNumber, UserDataRange.Start+hash)
	redisClient, ok := RedisClientMap[serverNumber]
	if !ok {
		return exceptions.Cache.ClientInstanceDoesNotExist()
	}

	formattedKey := formatUserDataKey(id)

	keys := []interface{}{formattedKey}
	argv := make([]interface{}, 0, len(dtos)*batchCheckAndUpdateUserQuotasByFormattedKeyBaseNumOfArgv)
	for _, dto := range dtos {
		argv = append(argv, dto.Field)
		argv = append(argv, dto.ChangeAmount)
		argv = append(argv, dto.MaxLimit)
		argv = append(argv, int(time.Until(dto.ExpiresIn).Seconds()))
	}

	arguments := []interface{}{
		"FCALL",
		redislibraries.BestEffortBatchCheckAndUpdateUserQuotasByFormattedKeyFunction,
		len(keys),
	}
	arguments = append(arguments, keys...)
	arguments = append(arguments, argv...)
	if _, err := redisClient.Do(arguments...).Result(); err != nil {
		return exceptions.Cache.FailedToUpdate(types.ValidCachePurpose_UserData.String()).WithError(err)
	}

	return nil
}

/* ========================= CRUD Operations ========================= */

func GetUserDataCache(id uuid.UUID) (*UserDataCache, *exceptions.Exception) {
	hash := hashUserDataIdentifier(id)
	serverNumber := min(MaxUserDataServerNumber, UserDataRange.Start+hash)
	redisClient, ok := RedisClientMap[serverNumber]
	if !ok {
		return nil, exceptions.Cache.ClientInstanceDoesNotExist()
	}

	formattedKey := formatUserDataKey(id)
	cacheString, err := redisClient.Get(formattedKey).Result()
	if err != nil {
		return nil, exceptions.Cache.NotFound(string(types.ValidCachePurpose_UserData)).WithError(err)
	}

	var userDataCache UserDataCache
	if err := json.Unmarshal([]byte(cacheString), &userDataCache); err != nil {
		// note that the json.Unmarshal() automatically return InvalidUnmarshalError if the userDataCache is nil
		return nil, exceptions.Cache.FailedToConvertJsonToStruct().WithError(err)
	}

	logs.FDebug("Successfully get the cached user data in the server with server number of %d", serverNumber)
	return &userDataCache, nil
}

func SetUserDataCache(id uuid.UUID, userDataCache UserDataCache) *exceptions.Exception {
	if !isValidUserCacheData(&userDataCache) { // strictly check when setting the cache data
		return exceptions.Cache.InvalidCacheDataStruct(userDataCache)
	}

	hash := hashUserDataIdentifier(id)
	serverNumber := min(MaxUserDataServerNumber, UserDataRange.Start+hash)
	redisClient, ok := RedisClientMap[serverNumber]
	if !ok {
		return exceptions.Cache.ClientInstanceDoesNotExist()
	}

	userDataJson, err := json.Marshal(userDataCache)
	if err != nil {
		return exceptions.Cache.FailedToConvertStructToJson().WithError(err)
	}

	formattedKey := formatUserDataKey(id)
	if err = redisClient.Set(formattedKey, string(userDataJson), _userDataCacheExpiresIn).Err(); err != nil {
		return exceptions.Cache.FailedToCreate(types.ValidCachePurpose_UserData.String()).WithError(err)
	}

	logs.FDebug("Successfully set the cached user data in the server with server number of %d", serverNumber)
	return nil
}

func UpdateUserDataCache(id uuid.UUID, dto UpdateUserDataCacheDto) *exceptions.Exception {
	hash := hashUserDataIdentifier(id)
	serverNumber := min(MaxUserDataServerNumber, UserDataRange.Start+hash)
	redisClient, ok := RedisClientMap[serverNumber]
	if !ok {
		return exceptions.Cache.ClientInstanceDoesNotExist()
	}

	userDataCache, exception := GetUserDataCache(id)
	if exception != nil {
		return exception
	}
	userDataCache.UpdatedAt = time.Now()
	if err := copier.Copy(&userDataCache, &dto); err != nil {
		return exceptions.Cache.FailedToConvertStructToJson().WithError(err)
	}
	userDataJson, err := json.Marshal(userDataCache)
	if err != nil {
		return exceptions.Cache.FailedToConvertStructToJson().WithError(err)
	}

	formattedKey := formatUserDataKey(id)
	if err = redisClient.Set(formattedKey, string(userDataJson), _userDataCacheExpiresIn).Err(); err != nil {
		return exceptions.Cache.FailedToUpdate(string(types.ValidCachePurpose_UserData)).WithError(err)
	}

	logs.FDebug("Successfully update the cached user data in the server with server number of %d", serverNumber)
	return nil
}

func DeleteUserDataCache(id uuid.UUID) *exceptions.Exception {
	hash := hashUserDataIdentifier(id)
	serverNumber := min(MaxUserDataServerNumber, UserDataRange.Start+hash)
	redisClient, ok := RedisClientMap[serverNumber]
	if !ok {
		return exceptions.Cache.ClientInstanceDoesNotExist()
	}

	formattedKey := formatUserDataKey(id)
	err := redisClient.Del(formattedKey).Err()
	if err != nil {
		return exceptions.Cache.FailedToDelete(string(types.ValidCachePurpose_UserData)).WithError(err)
	}

	logs.FDebug("Successfully delete the cached user data in the server with server number of %d", serverNumber)
	return nil
}
