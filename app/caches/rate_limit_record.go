package caches

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"math/rand"
	"time"

	uuid "github.com/google/uuid"

	redislibraries "notezy-backend/app/caches/libraries"
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
	_jitterMaxOffset = 5 * time.Second

	batchSynchronizeFunctionArgvPerKey = 2
	batchDeleteFunctionArgvPerKey      = 0
)

var (
	// if the rate limit range(the number of redis server is not enough, we may use another docker serivce for the rate limit redis cache)
	RateLimitRange                         = types.Range[int, int]{Start: 4, Size: 4} // server number: 4 - 7 (included)
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
		return 1
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

func GetRateLimitRecordCacheByFingerprint(
	fingerprint string,
	backendServerName types.BackendServerName,
) (*RateLimitRecordCache, *exceptions.Exception) {
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

func SetRateLimitRecordCacheByFingerprint(
	fingerprint string,
	backendServerName types.BackendServerName,
	rateLimitRecordCache RateLimitRecordCache,
) *exceptions.Exception {
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

func UpdateSyncrhronizeRateLimitRecordCacheByFingerprint(
	fingerprint string,
	backendServerName types.BackendServerName,
	dto SynchronizeRateLimitRecordCacheDto,
) *exceptions.Exception {
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

func DeleteRateLimitRecordCacheByFingerprint(
	fingerprint string,
	backendServerName types.BackendServerName,
) *exceptions.Exception {
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

	keys := make([]interface{}, 0, len(dtos))
	argv := make([]interface{}, 0, len(dtos)*batchSynchronizeFunctionArgvPerKey)
	for _, dto := range dtos {
		keys = append(keys, formatRateLimitKeyByFingerprint(dto.Fingerprint))
		argv = append(argv,
			dto.SynchronizeDto.NumOfChangingTokens,
			dto.SynchronizeDto.IsAccumulated,
		)
	}

	arguments := make([]interface{}, 0)
	arguments = append(arguments, "FCALL")
	arguments = append(arguments, redislibraries.BatchSynchronizeRateLimitRecordByFormattedKeysFunction)
	arguments = append(arguments, len(dtos))
	arguments = append(arguments, keys...)
	arguments = append(arguments, argv...)
	if _, err := redisClient.Do(arguments...).Result(); err != nil {
		return exceptions.Cache.FailedToUpdate(types.ValidCachePurpose_RateLimite.String()).WithError(err)
	}

	logs.FDebug("Successfully batch update cached rate limit records in the server with server number of %d", serverNumber)
	return nil
}

func BatchDeleteRateLimiteCachesByFingerprints(fingerprints []string, backendServerName types.BackendServerName) *exceptions.Exception {
	if len(fingerprints) == 0 {
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

	keys := make([]interface{}, 0, len(fingerprints))
	for _, fingerprint := range fingerprints {
		keys = append(keys, formatRateLimitKeyByFingerprint(fingerprint))
	}

	arguments := make([]interface{}, 0)
	arguments = append(arguments, "FCALL")
	arguments = append(arguments, redislibraries.BatchDeleteRateLimitRecordByFormattedKeysFunction)
	arguments = append(arguments, len(fingerprints))
	arguments = append(arguments, keys...)
	if _, err := redisClient.Do(arguments...).Result(); err != nil {
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

	keys := make([]interface{}, 0, len(dtos))
	argv := make([]interface{}, 0, len(dtos)*batchSynchronizeFunctionArgvPerKey)

	for _, dto := range dtos {
		keys = append(keys, formateRateLimitKeyByUserId(dto.UserId))
		argv = append(argv,
			dto.SynchronizeDto.NumOfChangingTokens,
			dto.SynchronizeDto.IsAccumulated,
		)
	}

	arguments := make([]interface{}, 0)
	arguments = append(arguments, "FCALL")
	arguments = append(arguments, redislibraries.BatchSynchronizeRateLimitRecordByFormattedKeysFunction)
	arguments = append(arguments, len(keys))
	arguments = append(arguments, keys...)
	arguments = append(arguments, argv...)
	if _, err := redisClient.Do(arguments...).Result(); err != nil {
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

	keys := make([]interface{}, 0, len(userIds))
	for _, userId := range userIds {
		keys = append(keys, formateRateLimitKeyByUserId(userId))
	}

	arguments := make([]interface{}, 0)
	arguments = append(arguments, "FCALL")
	arguments = append(arguments, redislibraries.BatchDeleteRateLimitRecordByFormattedKeysFunction)
	arguments = append(arguments, len(userIds))
	arguments = append(arguments, keys...)
	if _, err := redisClient.Do(arguments...).Result(); err != nil {
		return exceptions.Cache.FailedToDelete(types.ValidCachePurpose_RateLimite.String()).WithError(err)
	}

	logs.FDebug("Successfully delete cached rate limit records in the server with server number of %d", serverNumber)
	return nil
}
