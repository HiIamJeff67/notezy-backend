package caches

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"math/rand"
	"time"

	"github.com/go-redis/redis"
	uuid "github.com/google/uuid"

	exceptions "notezy-backend/app/exceptions"
	logs "notezy-backend/app/logs"
	types "notezy-backend/shared/types"
)

type RateLimitRecordCache struct {
	NumOfTokens     int32         `json:"numOfTokens"`
	WindowStartTime time.Time     `json:"windowStartTime"`
	WindowDuration  time.Duration `json:"windowDuration"`
	UpdatedAt       time.Time     `json:"updatedAt"`
}

type SynchronizeRateLimitRecordCacheDto struct {
	NumOfChangingTokens int32 `json:"numOfChangingTokens"`
	IsAccumulated       bool  `json:"isAccumulated"`
}

const (
	_defaultRateLimitWindow = 1 * time.Minute
	_jitterMaxOffset        = 5 * time.Second
)

var (
	// if the rate limit range(the number of redis server is not enough, we may use another docker serivce for the rate limit redis cache)
	RateLimitRange                         = types.Range{Start: 4, Size: 4} // server number: 4 - 7 (included)
	MaxRateLimitServerNumber               = RateLimitRange.Size - 1
	BackendServerNameToRateLimitRedisIndex = map[types.BackendServerName]int{
		types.BackendServerName_EastAsia:    4,
		types.BackendServerName_EastAmerica: 5,
		types.BackendServerName_WestAmerica: 6,
		types.BackendServerName_WestEurope:  7,
	}
)

/* ============================== Auxiliary Function ============================== */

func formatRateLimitKeyByFingerprint(fingerprint string) string {
	return fmt.Sprintf("%s:%s", types.ValidCachePurpose_RateLimite.String(), fingerprint)
}

func formateRateLimitKeyByUserId(id uuid.UUID) string {
	return fmt.Sprintf("%s:%s", types.ValidCachePurpose_RateLimite.String(), id.String())
}

func calculateExpirationTime(fingerprint string, windowStart time.Time, windowDuration time.Duration) time.Duration {
	nextResetTime := windowStart.Add(windowDuration)
	now := time.Now()

	baseExpirationTime := nextResetTime.Sub(now)
	if baseExpirationTime < 0 {
		return baseExpirationTime
	}

	h := fnv.New32a()
	h.Write([]byte(fingerprint))
	seed := int64(h.Sum32())

	rng := rand.New(rand.NewSource(seed))
	jitterOffset := time.Duration(rng.Int63n(int64(_jitterMaxOffset)))
	expirationTime := baseExpirationTime + jitterOffset

	return expirationTime
}

/* ============================== CRUD Operations By Client IP ============================== */

func GetRateLimitRecordCacheByFingerprint(fingerprint string, backendServerName types.BackendServerName) (*RateLimitRecordCache, *exceptions.Exception) {
	serverNumber, exist := BackendServerNameToRateLimitRedisIndex[backendServerName]
	if !exist {
		return nil, exceptions.Cache.BackendServerNameNotReferenced(types.ValidCachePurpose_RateLimite.String())
	}

	redisClient, exist := RedisClientMap[serverNumber]
	if !exist {
		return nil, exceptions.Cache.RedisServerNumberNotFound()
	}

	formattedKey := formatRateLimitKeyByFingerprint(fingerprint)
	cacheString, err := redisClient.Get(formattedKey).Result()
	if err != nil {
		return nil, exceptions.Cache.NotFound(string(types.ValidCachePurpose_RateLimite)).WithError(err)
	}

	var rateLimitRecordCache RateLimitRecordCache
	if err := json.Unmarshal([]byte(cacheString), &rateLimitRecordCache); err != nil {
		return nil, exceptions.Cache.FailedToConvertJsonToStruct().WithError(err)
	}

	logs.FDebug("Successfully get the cached rate limit in the server with server number of %d", serverNumber)
	return &rateLimitRecordCache, nil
}

func SetRateLimitRecordCacheByFingerprint(fingerprint string, backendServerName types.BackendServerName, rateLimitRecordCache RateLimitRecordCache) *exceptions.Exception {
	serverNumber, exist := BackendServerNameToRateLimitRedisIndex[backendServerName]
	if !exist {
		return exceptions.Cache.BackendServerNameNotReferenced(types.ValidCachePurpose_RateLimite.String())
	}

	redisClient, exist := RedisClientMap[serverNumber]
	if !exist {
		return exceptions.Cache.RedisServerNumberNotFound()
	}

	rateLimitJson, err := json.Marshal(rateLimitRecordCache)
	if err != nil {
		return exceptions.Cache.FailedToConvertJsonToStruct().WithError(err)
	}

	expirationTime := calculateExpirationTime(
		fingerprint,
		rateLimitRecordCache.WindowStartTime,
		rateLimitRecordCache.WindowDuration,
	)

	formattedKey := formatRateLimitKeyByFingerprint(fingerprint)
	if err = redisClient.Set(formattedKey, string(rateLimitJson), expirationTime).Err(); err != nil {
		return exceptions.Cache.FailedToCreate(types.ValidCachePurpose_RateLimite.String()).WithError(err)
	}

	logs.FDebug("Successfully set the cached rate limit record in the server with server number of %d", serverNumber)
	return nil
}

func UpdateSyncrhronizeRateLimitRecordCacheByFingerprint(fingerprint string, backendServerName types.BackendServerName, dto SynchronizeRateLimitRecordCacheDto) *exceptions.Exception {
	// TODO: since we use get and set which means more than or equal to two operations in this single operation,
	// 		 so we may need to use transaction to ensure the atomic

	serverNumber, exist := BackendServerNameToRateLimitRedisIndex[backendServerName]
	if !exist {
		return exceptions.Cache.BackendServerNameNotReferenced(types.ValidCachePurpose_RateLimite.String())
	}

	redisClient, exist := RedisClientMap[serverNumber]
	if !exist {
		return exceptions.Cache.RedisServerNumberNotFound()
	}

	rateLimitRecordCache, exception := GetRateLimitRecordCacheByFingerprint(fingerprint, backendServerName)
	if exception != nil {
		return exception
	}

	if (!dto.IsAccumulated && rateLimitRecordCache.NumOfTokens < dto.NumOfChangingTokens) || rateLimitRecordCache.NumOfTokens < 0 {
		return exceptions.Auth.InvalidRateLimitTokenCount()
	}

	if dto.IsAccumulated {
		rateLimitRecordCache.NumOfTokens += dto.NumOfChangingTokens
	} else {
		rateLimitRecordCache.NumOfTokens -= dto.NumOfChangingTokens
	}

	rateLimitJson, err := json.Marshal(rateLimitRecordCache)
	if err != nil {
		return exceptions.Cache.FailedToConvertStructToJson().WithError(err)
	}

	newExpirationTime := calculateExpirationTime(
		fingerprint,
		rateLimitRecordCache.WindowStartTime,
		rateLimitRecordCache.WindowDuration,
	)

	formattedKey := formatRateLimitKeyByFingerprint(fingerprint)
	if err = redisClient.Set(formattedKey, string(rateLimitJson), newExpirationTime).Err(); err != nil {
		return exceptions.Cache.FailedToUpdate(types.ValidCachePurpose_RateLimite.String()).WithError(err)
	}

	logs.FDebug("Successfully update the cached rate limit record in the server with server number of %d", serverNumber)
	return nil
}

func DeleteRateLimitRecordCacheByFingerprint(fingerprint string, backendServerName types.BackendServerName) *exceptions.Exception {
	serverNumber, exist := BackendServerNameToRateLimitRedisIndex[backendServerName]
	if !exist {
		return exceptions.Cache.BackendServerNameNotReferenced(types.ValidCachePurpose_RateLimite.String())
	}

	redisClient, exist := RedisClientMap[serverNumber]
	if !exist {
		return exceptions.Cache.RedisServerNumberNotFound()
	}

	formattedKey := formatRateLimitKeyByFingerprint(fingerprint)
	if err := redisClient.Del(formattedKey).Err(); err != nil {
		return exceptions.Cache.FailedToDelete(types.ValidCachePurpose_RateLimite.String()).WithError(err)
	}

	logs.FDebug("Successfully delete the cached rate limit record in the server with server number of %d", serverNumber)
	return nil
}

func BatchSynchronizeRateLimitRecordCachesByFingerprints(
	dtos []struct {
		Fingerprint    string                             `json:"fingerprint"`
		SynchronizeDto SynchronizeRateLimitRecordCacheDto `json:"synchronizeDto"`
	},
	backendServerName types.BackendServerName,
) *exceptions.Exception {
	if len(dtos) == 0 {
		return nil
	}

	serverNumber, exist := BackendServerNameToRateLimitRedisIndex[backendServerName]
	if !exist {
		return exceptions.Cache.BackendServerNameNotReferenced(types.ValidCachePurpose_RateLimite.String())
	}

	redisClient, exist := RedisClientMap[serverNumber]
	if !exist {
		return exceptions.Cache.RedisServerNumberNotFound()
	}

	updateMap := make(map[string]struct {
		Fingerprint    string
		SynchronizeDto SynchronizeRateLimitRecordCacheDto
	})
	watchKeys := make([]string, 0)
	for _, dto := range dtos {
		formattedKey := formatRateLimitKeyByFingerprint(dto.Fingerprint)
		if _, exists := updateMap[formattedKey]; exists {
			existing := updateMap[formattedKey]
			existing.SynchronizeDto.NumOfChangingTokens += dto.SynchronizeDto.NumOfChangingTokens
			updateMap[formattedKey] = existing
		}
		updateMap[formattedKey] = struct {
			Fingerprint    string
			SynchronizeDto SynchronizeRateLimitRecordCacheDto
		}{
			Fingerprint:    dto.Fingerprint,
			SynchronizeDto: dto.SynchronizeDto,
		}
		watchKeys = append(watchKeys, formattedKey)
	}

	txf := func(tx *redis.Tx) error {
		// first batch fetching the caches
		pipe := tx.TxPipeline()

		getCmds := make(map[string]*redis.StringCmd)
		for formattedKey := range updateMap {
			getCmds[formattedKey] = pipe.Get(formattedKey)
		}

		if _, err := pipe.Exec(); err != nil {
			return err
		}

		// then using the caches fetched on the top and batch updating the caches
		mutatePipe := tx.TxPipeline()

		for formattedKey, val := range updateMap {
			getCmd := getCmds[formattedKey]
			cacheString, err := getCmd.Result()
			if err != nil {
				continue
			}

			var rateLimitRecordCache RateLimitRecordCache
			if err := json.Unmarshal([]byte(cacheString), &rateLimitRecordCache); err != nil {
				continue
			}

			if (!val.SynchronizeDto.IsAccumulated && rateLimitRecordCache.NumOfTokens < val.SynchronizeDto.NumOfChangingTokens) || rateLimitRecordCache.NumOfTokens < 0 {
				continue
			}

			if val.SynchronizeDto.IsAccumulated {
				rateLimitRecordCache.NumOfTokens += val.SynchronizeDto.NumOfChangingTokens
			} else {
				rateLimitRecordCache.NumOfTokens -= val.SynchronizeDto.NumOfChangingTokens
			}

			rateLimitRecordCache.UpdatedAt = time.Now()

			rateLimitJson, err := json.Marshal(rateLimitRecordCache)
			if err != nil {
				continue
			}

			newExpirationTime := calculateExpirationTime(
				val.Fingerprint,
				rateLimitRecordCache.WindowStartTime,
				rateLimitRecordCache.WindowDuration,
			)

			mutatePipe.Set(formattedKey, string(rateLimitJson), newExpirationTime)
		}

		if _, err := mutatePipe.Exec(); err != nil {
			return err
		}

		return nil
	}

	err := redisClient.Watch(txf, watchKeys...)
	if err != nil {
		return exceptions.Cache.FailedToUpdate(types.ValidCachePurpose_RateLimite.String()).WithError(err)
	}

	logs.FDebug("Successfully batch update cached rate limit records in the server with server number of %d", serverNumber)
	return nil
}

func BatchDeleteRateLimiteCachesByFingerprints(Fingerprints []string, backendServerName types.BackendServerName) *exceptions.Exception {
	if len(Fingerprints) == 0 {
		return nil
	}

	serverNumber, exist := BackendServerNameToRateLimitRedisIndex[backendServerName]
	if !exist {
		return exceptions.Cache.BackendServerNameNotReferenced(types.ValidCachePurpose_RateLimite.String())
	}

	redisClient, exist := RedisClientMap[serverNumber]
	if !exist {
		return exceptions.Cache.RedisServerNumberNotFound()
	}

	txf := func(tx *redis.Tx) error {
		formattedKeys := make([]string, len(Fingerprints))
		seen := make(map[string]bool)
		for _, fingerprint := range Fingerprints {
			if seen[fingerprint] {
				continue
			}
			seen[fingerprint] = true
			formattedKey := formatRateLimitKeyByFingerprint(fingerprint)
			formattedKeys = append(formattedKeys, formattedKey)
		}

		pipe := tx.TxPipeline()
		pipe.Del(formattedKeys...).Result()
		if _, err := pipe.Exec(); err != nil {
			return err
		}

		return nil
	}

	err := redisClient.Watch(txf)
	if err != nil {
		return exceptions.Cache.FailedToDelete(types.ValidCachePurpose_RateLimite.String()).WithError(err)
	}

	logs.FDebug("Successfully batch delete cached rate limit records in the server with server number of %d", serverNumber)
	return nil
}

/* ============================== CRUD Operations By UserId ============================== */

func GetRateLimitRecordCacheByUserId(userId uuid.UUID, fingerprint string, backendServerName types.BackendServerName) (*RateLimitRecordCache, *exceptions.Exception) {
	serverNumber, exist := BackendServerNameToRateLimitRedisIndex[backendServerName]
	if !exist {
		return nil, exceptions.Cache.BackendServerNameNotReferenced(types.ValidCachePurpose_RateLimite.String())
	}

	redisClient, exist := RedisClientMap[serverNumber]
	if !exist {
		return nil, exceptions.Cache.RedisServerNumberNotFound()
	}

	formattedKey := formateRateLimitKeyByUserId(userId)
	cacheString, err := redisClient.Get(formattedKey).Result()
	if err != nil {
		return nil, exceptions.Cache.NotFound(string(types.ValidCachePurpose_RateLimite)).WithError(err)
	}

	var rateLimitRecordCache RateLimitRecordCache
	if err := json.Unmarshal([]byte(cacheString), &rateLimitRecordCache); err != nil {
		return nil, exceptions.Cache.FailedToConvertJsonToStruct().WithError(err)
	}

	logs.FDebug("Successfully get the cached rate limit in the server with server number of %d", serverNumber)
	return &rateLimitRecordCache, nil
}

func SetRateLimitRecordCacheByUserId(userId uuid.UUID, fingerprint string, backendServerName types.BackendServerName, rateLimitRecordCache RateLimitRecordCache) *exceptions.Exception {
	serverNumber, exist := BackendServerNameToRateLimitRedisIndex[backendServerName]
	if !exist {
		return exceptions.Cache.BackendServerNameNotReferenced(types.ValidCachePurpose_RateLimite.String())
	}

	redisClient, exist := RedisClientMap[serverNumber]
	if !exist {
		return exceptions.Cache.RedisServerNumberNotFound()
	}

	rateLimitJson, err := json.Marshal(rateLimitRecordCache)
	if err != nil {
		return exceptions.Cache.FailedToConvertJsonToStruct().WithError(err)
	}

	expirationTime := calculateExpirationTime(
		fingerprint,
		rateLimitRecordCache.WindowStartTime,
		rateLimitRecordCache.WindowDuration,
	)

	formattedKey := formateRateLimitKeyByUserId(userId)
	if err = redisClient.Set(formattedKey, string(rateLimitJson), expirationTime).Err(); err != nil {
		return exceptions.Cache.FailedToCreate(types.ValidCachePurpose_RateLimite.String()).WithError(err)
	}

	logs.FDebug("Successfully set the cached rate limit record in the server with server number of %d", serverNumber)
	return nil
}

func UpdateRateLimitRecordCacheByUserId(userId uuid.UUID, fingerprint string, backendServerName types.BackendServerName, dto SynchronizeRateLimitRecordCacheDto) *exceptions.Exception {
	// TODO: since we use get and set which means more than or equal to two operations in this single operation,
	// 		 so we may need to use transaction to ensure the atomic

	serverNumber, exist := BackendServerNameToRateLimitRedisIndex[backendServerName]
	if !exist {
		return exceptions.Cache.BackendServerNameNotReferenced(types.ValidCachePurpose_RateLimite.String())
	}

	redisClient, exist := RedisClientMap[serverNumber]
	if !exist {
		return exceptions.Cache.RedisServerNumberNotFound()
	}

	rateLimitRecordCache, exception := GetRateLimitRecordCacheByUserId(userId, fingerprint, backendServerName)
	if exception != nil {
		return exception
	}

	if (!dto.IsAccumulated && rateLimitRecordCache.NumOfTokens < dto.NumOfChangingTokens) || rateLimitRecordCache.NumOfTokens < 0 {
		return exceptions.Auth.InvalidRateLimitTokenCount()
	}

	if dto.IsAccumulated {
		rateLimitRecordCache.NumOfTokens += dto.NumOfChangingTokens
	} else {
		rateLimitRecordCache.NumOfTokens -= dto.NumOfChangingTokens
	}

	rateLimitJson, err := json.Marshal(rateLimitRecordCache)
	if err != nil {
		return exceptions.Cache.FailedToConvertStructToJson().WithError(err)
	}

	newExpirationTime := calculateExpirationTime(
		fingerprint,
		rateLimitRecordCache.WindowStartTime,
		rateLimitRecordCache.WindowDuration,
	)

	formattedKey := formateRateLimitKeyByUserId(userId)
	if err = redisClient.Set(formattedKey, string(rateLimitJson), newExpirationTime).Err(); err != nil {
		return exceptions.Cache.FailedToUpdate(types.ValidCachePurpose_RateLimite.String()).WithError(err)
	}

	logs.FDebug("Successfully update the cached rate limit record in the server with server number of %d", serverNumber)
	return nil
}

func DeleteRateLimitRecordCacheByUserId(userId uuid.UUID, fingerprint string, backendServerName types.BackendServerName) *exceptions.Exception {
	serverNumber, exist := BackendServerNameToRateLimitRedisIndex[backendServerName]
	if !exist {
		return exceptions.Cache.BackendServerNameNotReferenced(types.ValidCachePurpose_RateLimite.String())
	}

	redisClient, exist := RedisClientMap[serverNumber]
	if !exist {
		return exceptions.Cache.RedisServerNumberNotFound()
	}

	formattedKey := formateRateLimitKeyByUserId(userId)
	if err := redisClient.Del(formattedKey).Err(); err != nil {
		return exceptions.Cache.FailedToDelete(types.ValidCachePurpose_RateLimite.String()).WithError(err)
	}

	logs.FDebug("Successfully delete the cached rate limit record in the server with server number of %d", serverNumber)
	return nil
}

func BatchSynchronizeRateLimitRecordCachesByUserIds(
	dtos []struct {
		UserId         uuid.UUID                          `json:"userId"`
		SynchronizeDto SynchronizeRateLimitRecordCacheDto `json:"synchronizeDto"`
	},
	backendServerName types.BackendServerName,
) *exceptions.Exception {
	if len(dtos) == 0 {
		return nil
	}

	serverNumber, exist := BackendServerNameToRateLimitRedisIndex[backendServerName]
	if !exist {
		return exceptions.Cache.BackendServerNameNotReferenced(types.ValidCachePurpose_RateLimite.String())
	}

	redisClient, exist := RedisClientMap[serverNumber]
	if !exist {
		return exceptions.Cache.RedisServerNumberNotFound()
	}

	updateMap := make(map[string]struct {
		UserId         uuid.UUID
		SynchronizeDto SynchronizeRateLimitRecordCacheDto
	})
	watchKeys := make([]string, 0)
	for _, dto := range dtos {
		formattedKey := formateRateLimitKeyByUserId(dto.UserId)
		updateMap[formattedKey] = struct {
			UserId         uuid.UUID
			SynchronizeDto SynchronizeRateLimitRecordCacheDto
		}{
			UserId:         dto.UserId,
			SynchronizeDto: dto.SynchronizeDto,
		}
		watchKeys = append(watchKeys, formattedKey)
	}

	txf := func(tx *redis.Tx) error {
		// first batch fetching the caches
		pipe := tx.TxPipeline()

		getCmds := make(map[string]*redis.StringCmd)
		for formattedKey := range updateMap {
			getCmds[formattedKey] = pipe.Get(formattedKey)
		}

		pipe.Exec() // ignore if there's any error of getting the caches of the request

		// then using the caches fetched on the top and batch updating the caches
		mutatePipe := tx.TxPipeline()

		for formattedKey, val := range updateMap {
			getCmd := getCmds[formattedKey]
			cacheString, err := getCmd.Result()
			if err != nil {
				if err == redis.Nil { // if the key does not exist
					logs.FDebug("Key not found, creating new record for user: %s", val.UserId.String())
					newCache := RateLimitRecordCache{
						NumOfTokens:     val.SynchronizeDto.NumOfChangingTokens,
						WindowStartTime: time.Now(),
						WindowDuration:  _defaultRateLimitWindow,
						UpdatedAt:       time.Now(),
					}

					rateLimitJson, marshalErr := json.Marshal(newCache)
					if marshalErr != nil {
						logs.FError("Failed to marshal new cache: %v", marshalErr)
						continue
					}

					newExpirationTime := calculateExpirationTime(
						val.UserId.String(),
						newCache.WindowStartTime,
						newCache.WindowDuration,
					)

					mutatePipe.Set(formattedKey, string(rateLimitJson), newExpirationTime)
					continue
				} else {
					logs.FError("Failed to get cache for key %s: %v", formattedKey, err)
					continue
				}
			}

			var rateLimitRecordCache RateLimitRecordCache
			if err := json.Unmarshal([]byte(cacheString), &rateLimitRecordCache); err != nil {
				continue
			}

			if (!val.SynchronizeDto.IsAccumulated && rateLimitRecordCache.NumOfTokens < val.SynchronizeDto.NumOfChangingTokens) || rateLimitRecordCache.NumOfTokens < 0 {
				continue
			}

			if val.SynchronizeDto.IsAccumulated {
				rateLimitRecordCache.NumOfTokens += val.SynchronizeDto.NumOfChangingTokens
			} else {
				rateLimitRecordCache.NumOfTokens -= val.SynchronizeDto.NumOfChangingTokens
			}

			rateLimitRecordCache.UpdatedAt = time.Now()

			rateLimitJson, err := json.Marshal(rateLimitRecordCache)
			if err != nil {
				continue
			}

			newExpirationTime := calculateExpirationTime(
				val.UserId.String(),
				rateLimitRecordCache.WindowStartTime,
				rateLimitRecordCache.WindowDuration,
			)

			mutatePipe.Set(formattedKey, string(rateLimitJson), newExpirationTime)
		}

		if _, err := mutatePipe.Exec(); err != nil {
			return err
		}

		return nil
	}

	err := redisClient.Watch(txf, watchKeys...)
	if err != nil {
		return exceptions.Cache.FailedToUpdate(types.ValidCachePurpose_RateLimite.String()).WithError(err)
	}

	logs.FDebug("Successfully batch update cached rate limit records in the server with server number of %d", serverNumber)
	return nil
}

func BatchDeleteRateLimiteCachesByUserIds(userIds []uuid.UUID, backendServerName types.BackendServerName) *exceptions.Exception {
	// the batch delete operation required redis transaction and pipeline

	serverNumber, exist := BackendServerNameToRateLimitRedisIndex[backendServerName]
	if !exist {
		return exceptions.Cache.BackendServerNameNotReferenced(types.ValidCachePurpose_RateLimite.String())
	}

	redisClient, exist := RedisClientMap[serverNumber]
	if !exist {
		return exceptions.Cache.RedisServerNumberNotFound()
	}

	if len(userIds) == 0 {
		return nil
	}

	txf := func(tx *redis.Tx) error {
		formattedKeys := make([]string, len(userIds))
		seen := make(map[uuid.UUID]bool)
		for _, id := range userIds {
			if seen[id] {
				continue
			}
			seen[id] = true
			formattedKey := formateRateLimitKeyByUserId(id)
			formattedKeys = append(formattedKeys, formattedKey)
		}

		pipe := tx.TxPipeline()
		pipe.Del(formattedKeys...).Result()
		if _, err := pipe.Exec(); err != nil {
			return err
		}

		return nil
	}

	err := redisClient.Watch(txf)
	if err != nil {
		return exceptions.Cache.FailedToDelete(types.ValidCachePurpose_RateLimite.String()).WithError(err)
	}

	logs.FDebug("Successfully delete cached rate limit records in the server with server number of %d", serverNumber)
	return nil
}
