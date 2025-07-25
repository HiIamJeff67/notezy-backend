package caches

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"strings"
	"time"

	uuid "github.com/google/uuid"
	"github.com/jinzhu/copier"

	exceptions "notezy-backend/app/exceptions"
	logs "notezy-backend/app/logs"
	enums "notezy-backend/app/models/schemas/enums"
	shared "notezy-backend/shared"
	types "notezy-backend/shared/types"
)

type UserDataCache struct {
	PublicId           string           `json:"publicId"`           // user
	Name               string           `json:"name"`               // user
	DisplayName        string           `json:"displayName"`        // user
	Email              string           `json:"email"`              // user
	AccessToken        string           `json:"accessToken"`        // only here
	Role               enums.UserRole   `json:"role"`               // user
	Plan               enums.UserPlan   `json:"plan"`               // user
	Status             enums.UserStatus `json:"status"`             // user
	AvatarURL          string           `json:"avatarURL"`          // user info
	Language           enums.Language   `json:"language"`           // user setting
	GeneralSettingCode int64            `json:"generalSettingCode"` // user setting
	PrivacySettingCode int64            `json:"privacySettingCode"` // user setting
	UpdatedAt          time.Time        `json:"updatedAt"`          // cache
}

type UpdateUserDataCacheDto struct {
	PublicId           *string
	Name               *string
	DisplayName        *string
	Email              *string
	AccessToken        *string
	Role               *enums.UserRole
	Plan               *enums.UserPlan
	Status             *enums.UserStatus
	AvatarURL          *string
	Language           *enums.Language
	GeneralSettingCode *int64
	PrivacySettingCode *int64
}

const (
	_userDataCacheExpiresIn = 24 * time.Hour
)

var (
	UserDataRange           = types.Range{Start: 0, Size: 4} // server number: 0 - 3 (included)
	MaxUserDataServerNumber = UserDataRange.Start + UserDataRange.Size - 1
)

/* ========================= Auxiliary Function ========================= */

func hashUserDataIdentifier(identifier uuid.UUID) int {
	h := fnv.New32a()
	h.Write([]byte(identifier.String()))
	return int(h.Sum32()) % UserDataRange.Size
}

func formatKey(id uuid.UUID) string {
	return fmt.Sprintf("UserId:%s", id.String())
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

/* ========================= CRUD Operations of the Cache ========================= */

func GetUserDataCache(id uuid.UUID) (*UserDataCache, *exceptions.Exception) {
	hash := hashUserDataIdentifier(id)
	serverNumber := min(MaxUserDataServerNumber, UserDataRange.Start+hash)
	redisClient, ok := RedisClientMap[serverNumber]
	if !ok {
		return nil, exceptions.Cache.ClientInstanceDoesNotExist(serverNumber)
	}

	formattedKey := formatKey(id)
	cacheString, err := redisClient.Get(formattedKey).Result()
	if err != nil {
		return nil, exceptions.Cache.NotFound(shared.ValidCachePurpose_UserData).WithError(err)
	}

	var userDataCache UserDataCache
	if err := json.Unmarshal([]byte(cacheString), &userDataCache); err != nil {
		// note that the json.Unmarshal() automatically return InvalidUnmarshalError if the userDataCache is nil
		return nil, exceptions.Cache.FailedToConvertJsonToStruct().WithError(err)
	}

	logs.FInfo("Successfully get the cached user data in the server with server number of %d", serverNumber)
	return &userDataCache, nil
}

func SetUserDataCache(id uuid.UUID, userData UserDataCache) *exceptions.Exception {
	if !isValidUserCacheData(&userData) { // strictly check when setting the cache data
		return exceptions.Cache.InvalidCacheDataStruct(userData)
	}

	hash := hashUserDataIdentifier(id)
	serverNumber := min(MaxUserDataServerNumber, UserDataRange.Start+hash)
	redisClient, ok := RedisClientMap[serverNumber]
	if !ok {
		return exceptions.Cache.ClientInstanceDoesNotExist(serverNumber)
	}

	userDataJson, err := json.Marshal(userData)
	if err != nil {
		return exceptions.Cache.FailedToConvertStructToJson().WithError(err)
	}

	formattedKey := formatKey(id)
	err = redisClient.Set(formattedKey, string(userDataJson), _userDataCacheExpiresIn).Err()
	if err != nil {
		return exceptions.Cache.FailedToCreate(shared.ValidCachePurpose_UserData).WithError(err)
	}

	logs.FInfo("Successfully set the cached user data in the server with server number of %d", serverNumber)
	return nil
}

func UpdateUserDataCache(id uuid.UUID, dto UpdateUserDataCacheDto) *exceptions.Exception {
	hash := hashUserDataIdentifier(id)
	serverNumber := min(MaxUserDataServerNumber, UserDataRange.Start+hash)
	redisClient, ok := RedisClientMap[serverNumber]
	if !ok {
		return exceptions.Cache.ClientInstanceDoesNotExist(serverNumber)
	}

	userData, exception := GetUserDataCache(id)
	if exception != nil {
		return exception
	}
	userData.UpdatedAt = time.Now()
	if err := copier.Copy(&userData, &dto); err != nil {
		return exceptions.Cache.FailedToConvertStructToJson().WithError(err)
	}
	userDataJson, err := json.Marshal(userData)
	if err != nil {
		return exceptions.Cache.FailedToConvertStructToJson().WithError(err)
	}

	formattedKey := formatKey(id)
	err = redisClient.Set(formattedKey, string(userDataJson), _userDataCacheExpiresIn).Err()
	if err != nil {
		return exceptions.Cache.FailedToUpdate(shared.ValidCachePurpose_UserData).WithError(err)
	}

	logs.FInfo("Successfully update the cached user data in the server with server number of %d", serverNumber)
	return nil
}

func DeleteUserDataCache(id uuid.UUID) *exceptions.Exception {
	hash := hashUserDataIdentifier(id)
	serverNumber := min(MaxUserDataServerNumber, UserDataRange.Start+hash)
	redisClient, ok := RedisClientMap[serverNumber]
	if !ok {
		return exceptions.Cache.ClientInstanceDoesNotExist(serverNumber)
	}

	formattedKey := formatKey(id)
	err := redisClient.Del(formattedKey).Err()
	if err != nil {
		return exceptions.Cache.FailedToDelete(shared.ValidCachePurpose_UserData).WithError(err)
	}

	logs.FInfo("Successfully delete the cached user data in the server with server number of %d", serverNumber)
	return nil
}
