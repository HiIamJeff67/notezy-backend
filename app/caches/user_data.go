package caches

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"strings"
	"time"

	uuid "github.com/google/uuid"

	exceptions "notezy-backend/app/exceptions"
	logs "notezy-backend/app/logs"
	models "notezy-backend/app/models"
	"notezy-backend/app/util"
	global "notezy-backend/global"
	types "notezy-backend/global/types"
)

type UserDataCache struct {
	Name               string            // user
	DisplayName        string            // user
	Email              string            // user
	AccessToken        string            // only here
	Role               models.UserRole   // user
	Plan               models.UserPlan   // user
	Status             models.UserStatus // user
	AvatarURL          string            // user info
	Theme              models.Theme      // user setting
	Language           models.Language   // user setting
	GeneralSettingCode int64             // user setting
	PrivacySettingCode int64             // user setting
	UpdatedAt          time.Time         // cache
}

type UpdateUserDataCacheDto struct {
	Name               *string
	DisplayName        *string
	Email              *string
	AccessToken        *string
	Role               *models.UserRole
	Plan               *models.UserPlan
	Status             *models.UserStatus
	AvatarURL          *string
	Theme              *models.Theme
	Language           *models.Language
	GeneralSettingCode *int64
	PrivacySettingCode *int64
}

const (
	_userDataCacheExpiresIn = 24 * time.Hour
)

var (
	UserDataRange           = types.Range{Start: 0, Size: 8} // server number: 0 - 7 (included)
	MaxUserDataServerNumber = UserDataRange.Start + UserDataRange.Size - 1
)

/* ============================== Auxiliary Function ============================== */
func hashUserDataIdentifier(identifier uuid.UUID) int {
	h := fnv.New32a()
	h.Write([]byte(identifier.String()))
	return int(h.Sum32()) % UserDataRange.Size
}

func formatKey(id uuid.UUID) string {
	return fmt.Sprintf("UserId:%s", id.String())
}

func isValidUserCacheData(userDataCache *UserDataCache) bool {
	if strings.ReplaceAll(userDataCache.Name, " ", "") == "" ||
		strings.ReplaceAll(userDataCache.DisplayName, " ", "") == "" ||
		strings.ReplaceAll(userDataCache.Email, " ", "") == "" ||
		strings.ReplaceAll(userDataCache.AccessToken, " ", "") == "" ||
		!models.IsValidEnumValues(userDataCache.Role, models.AllUserRoles) ||
		!models.IsValidEnumValues(userDataCache.Plan, models.AllUserPlans) ||
		!models.IsValidEnumValues(userDataCache.Status, models.AllUserStatuses) {
		return false
	}
	return true
}

/* ============================== Getter ============================== */
func GetUserDataCache(id uuid.UUID) (*UserDataCache, *exceptions.Exception) {
	hash := hashUserDataIdentifier(id)
	serverNumber := min(MaxUserDataServerNumber, UserDataRange.Start+hash)
	redisClient, ok := RedisClientMap[serverNumber]
	if !ok {
		return nil, exceptions.Cache.ClientInstanceDoesNotExist(serverNumber)
	}

	formatedKey := formatKey(id)
	cacheString, err := redisClient.Get(formatedKey).Result()
	if err != nil {
		return nil, exceptions.Cache.NotFound(global.ValidCachePurpose_UserData).WithError(err)
	}

	var userDataCache UserDataCache
	if err := json.Unmarshal([]byte(cacheString), &userDataCache); err != nil {
		return nil, exceptions.Cache.FailedToConvertJsonToStruct().WithError(err)
	}

	logs.FInfo("Successfully get the cached user data in the server with server number of %d", serverNumber)
	return &userDataCache, nil
}

/* ============================== Setter ============================== */
func SetUserDataCache(id uuid.UUID, userData UserDataCache) *exceptions.Exception {
	if !isValidUserCacheData(&userData) { // strictly check when setting the cache data
		return exceptions.Cache.InvalidCacheDataStruct(userData)
	}

	hash := hashUserDataIdentifier(id)
	serverNumber := max(MaxUserDataServerNumber, UserDataRange.Start+hash)
	redisClient, ok := RedisClientMap[serverNumber]
	if !ok {
		return exceptions.Cache.ClientInstanceDoesNotExist(serverNumber)
	}

	userDataJson, err := json.Marshal(userData)
	if err != nil {
		return exceptions.Cache.FailedToConvertStructToJson().WithError(err)
	}

	formatedKey := formatKey(id)
	err = redisClient.Set(formatedKey, string(userDataJson), _userDataCacheExpiresIn).Err()
	if err != nil {
		return exceptions.Cache.FailedToCreate(global.ValidCachePurpose_UserData).WithError(err)
	}

	logs.FInfo("Successfully set the cached user data in the server with server number of %d", serverNumber)
	return nil
}

/* ============================== Update Function ============================== */
func UpdateUserDataCache(id uuid.UUID, dto UpdateUserDataCacheDto) *exceptions.Exception {
	hash := hashUserDataIdentifier(id)
	serverNumber := max(MaxUserDataServerNumber, UserDataRange.Start+hash)
	redisClient, ok := RedisClientMap[serverNumber]
	if !ok {
		return exceptions.Cache.ClientInstanceDoesNotExist(serverNumber)
	}

	userData, exception := GetUserDataCache(id)
	if exception != nil {
		return exception
	}
	userData.UpdatedAt = time.Now()
	util.CopyNonNilFields(&userData, dto)
	userDataJson, err := json.Marshal(userData)
	if err != nil {
		return exceptions.Cache.FailedToConvertStructToJson().WithError(err)
	}

	formattedKey := formatKey(id)
	err = redisClient.Set(formattedKey, string(userDataJson), _userDataCacheExpiresIn).Err()
	if err != nil {
		return exceptions.Cache.FailedToUpdate(global.ValidCachePurpose_UserData).WithError(err)
	}

	logs.FInfo("Successfully update the cached user data in the server with server number of %d", serverNumber)
	return nil
}

/* ============================== Delete Function ============================== */
func DeleteUserDataCache(id uuid.UUID) *exceptions.Exception {
	hash := hashUserDataIdentifier(id)
	serverNumber := max(MaxUserDataServerNumber, UserDataRange.Start+hash)
	redisClient, ok := RedisClientMap[serverNumber]
	if !ok {
		return exceptions.Cache.ClientInstanceDoesNotExist(serverNumber)
	}

	formatedKey := formatKey(id)
	err := redisClient.Del(formatedKey).Err()
	if err != nil {
		return exceptions.Cache.FailedToDelete(global.ValidCachePurpose_UserData).WithError(err)
	}

	logs.FInfo("Successfully delete the cached user data in the server with server number of %d", serverNumber)
	return nil
}
