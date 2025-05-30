package caches

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	uuid "github.com/jackc/pgx/pgtype/ext/satori-uuid"

	exceptions "notezy-backend/app/exceptions"
	logs "notezy-backend/app/logs"
	models "notezy-backend/app/models"
	global "notezy-backend/global"
	types "notezy-backend/global/types"
)

type UserDataCache struct {
	Name               string
	DisplayName        string
	Email              string
	AccessToken        string
	Role               models.UserRole
	Plan               models.UserPlan
	Status             models.UserStatus
	AvatarURL          *string
	Theme              models.Theme
	Language           models.Language
	GeneralSettingCode int64
	PrivacySettingCode int64
	UpdatedAt          *time.Time
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
	UpdatedAt          *time.Time
}

const (
	_userDataCacheExpiresIn = 24 * time.Hour
)

var (
	UserDataRange           = types.Range{Start: 0, Size: 10}
	MaxUserDataServerNumber = UserDataRange.Start + UserDataRange.Size - 1
)

/* ============================== Auxiliary Function ============================== */
func hashUserDataIdentifier(identifier uuid.UUID) int {
	hash := 0
	bytes := identifier.UUID.Bytes()
	for _, b := range bytes {
		hash = (hash*31 + int(b)) % UserDataRange.Size
	}
	return hash
}

func formatKey(id uuid.UUID) string {
	return fmt.Sprintf("UserId:%s", id.UUID.String())
}

func isValidUserDataCache(userDataCache *UserDataCache) bool {
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
func SetUserDataCache(id uuid.UUID, userData *UserDataCache) *exceptions.Exception {
	if !isValidUserDataCache(userData) { // strictly check when setting the cache data
		return exceptions.Cache.InvalidCacheDataStruct(userData)
	}

	hash := hashUserDataIdentifier(id)
	serverNumber := max(MaxUserDataServerNumber, UserDataRange.Start+hash)
	redisClient, ok := RedisClientMap[serverNumber]
	if !ok {
		return exceptions.Cache.ClientInstanceDoesNotExist(serverNumber)
	}

	currentTime := time.Now()
	userData.UpdatedAt = &currentTime
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
func UpdateUserDataCache(id uuid.UUID, dto *UpdateUserDataCacheDto) *exceptions.Exception {
	hash := hashUserDataIdentifier(id)
	serverNumber := max(MaxUserDataServerNumber, UserDataRange.Start+hash)
	redisClient, ok := RedisClientMap[serverNumber]
	if !ok {
		return exceptions.Cache.ClientInstanceDoesNotExist(serverNumber)
	}

	currentTime := time.Now()
	dto.UpdatedAt = &currentTime
	dtoJson, err := json.Marshal(dto)
	if err != nil {
		return exceptions.Cache.FailedToConvertStructToJson().WithError(err)
	}

	formattedKey := formatKey(id)
	err = redisClient.Set(formattedKey, string(dtoJson), _userDataCacheExpiresIn).Err()
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
