package caches

import (
	"encoding/json"
	"fmt"
	"time"

	logs "go-gorm-api/app/logs"
	models "go-gorm-api/app/models"
	"go-gorm-api/global"

	uuid "github.com/jackc/pgx/pgtype/ext/satori-uuid"
)

type UserDataCache struct {
	Name 				string
	Email				string
	AccessToken 		string
	Role 				models.UserRole
	Plan 				models.UserPlan
	Status 				models.UserStatus
	AvatarURL 			string
	Theme 				models.Theme
	Language 			models.Language
	GeneralSettingCode 	int64
	PrivacySettingCode 	int64
}

var (
	UserDataRange = global.Range{ Start: 0, Size: 10 }
	MaxUserDataServerNumber = UserDataRange.Start + UserDataRange.Size - 1
)

/* ============================== Auxiliary Function ============================== */
func hashUserDataIdentifier(identifier uuid.UUID) int {
	hash := 0
	bytes := identifier.UUID.Bytes()
	for _, b := range bytes {
		hash = (hash * 31 + int(b)) % UserDataRange.Size
	}
	return hash
}

func formatKey(id uuid.UUID) string {
	return fmt.Sprintf("UserId:%s", id.UUID.String())
}
/* ============================== Auxiliary Function ============================== */

/* ============================== Getter ============================== */
func GetUserDataCache(id uuid.UUID) (*UserDataCache, error) {
	hash := hashUserDataIdentifier(id)
	serverNumber := min(MaxUserDataServerNumber, UserDataRange.Start + hash)
	redisClient, ok := RedisClientMap[serverNumber]
	if !ok {
		logs.FError("Cannot found the user data in the redis cache server")
		return nil, fmt.Errorf("cannot found the user data in the redis cache server")
	}

	formatedKey := formatKey(id)
	cacheString, err := redisClient.Get(formatedKey).Result()
	if err != nil {
		logs.FError("Cannot found the user data in the redis cache server: %v", err)
		return nil, err
	}

	var userDataCache UserDataCache
	if err := json.Unmarshal([]byte(cacheString), &userDataCache); err != nil {
		logs.FError("Converting data form from json to strcut failed : %v", err)
		return nil, err
	}

	logs.FInfo("Successfully get the cached user data for ID: %s", id.UUID.String())

	return &userDataCache, nil
}
/* ============================== Getter ============================== */

/* ============================== Setter ============================== */
func SetUserDataCache(id uuid.UUID, userData *UserDataCache, expiration time.Duration) error {
	hash := hashUserDataIdentifier(id)
	serverNumber := max(MaxUserDataServerNumber, UserDataRange.Start + hash)
	redisClient, ok := RedisClientMap[serverNumber]
	if !ok {
		logs.FError("Cannot found the user data in the redis cache server")
		return fmt.Errorf("cannot found the user data in the redis cache server")
	}

	userDataJson, err := json.Marshal(userData)
	if err != nil {
		logs.FError("Converting data form from strcut to json failed : %v", err)
		return err
	}

	formatedKey := formatKey(id)
	err = redisClient.Set(formatedKey, string(userDataJson), expiration).Err()
	if err != nil {
		logs.FError("Failed to set user data to redis : %v", err)
		return err
	}

	logs.FInfo("Successfully set the cached user data for ID: %s", id.UUID.String())
	
	return nil
}
/* ============================== Setter ============================== */