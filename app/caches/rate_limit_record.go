package caches

import (
	"context"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"math/rand"
	"time"

	"github.com/go-redis/redis"

	cacheinputs "github.com/HiIamJeff67/notezy-backend/app/caches/inputs"
	redislibraries "github.com/HiIamJeff67/notezy-backend/app/caches/libraries"
	exceptions "github.com/HiIamJeff67/notezy-backend/app/exceptions"
	logs "github.com/HiIamJeff67/notezy-backend/app/monitor/logs"
	types "github.com/HiIamJeff67/notezy-backend/shared/types"
)

type RateLimitRecordCache struct {
	NumOfTokens     int32         `json:"numOfTokens"`
	WindowStartTime time.Time     `json:"windowStartTime"`
	WindowDuration  time.Duration `json:"windowDuration"`
	UpdatedAt       time.Time     `json:"updatedAt"`
}

type RateLimitRecordCacheStore struct {
	redisClientMap map[int]*redis.Client

	Range                         types.Range[int, int]
	MaxServerNumber               int
	backendServerNameToRedisIndex map[types.BackendServerName]int

	jitterMaxOffset                    time.Duration
	batchSynchronizeFunctionArgvPerKey int
}

/* ============================== Constructor ============================== */

func NewRateLimitRecordCacheStore(redisClientMap map[int]*redis.Client) *RateLimitRecordCacheStore {
	rangeValue := types.Range[int, int]{Start: 4, Size: 4}

	return &RateLimitRecordCacheStore{
		redisClientMap: redisClientMap,

		Range:           rangeValue,
		MaxServerNumber: rangeValue.Start + rangeValue.Size - 1,
		backendServerNameToRedisIndex: map[types.BackendServerName]int{
			types.BackendServerName_EastAsia:    4,
			types.BackendServerName_EastAmerica: 5,
			types.BackendServerName_WestAmerica: 6,
			types.BackendServerName_WestEurope:  7,
		},

		jitterMaxOffset:                    5 * time.Second,
		batchSynchronizeFunctionArgvPerKey: 2,
	}
}

/* ============================== Auxiliary Methods ============================== */

func (s *RateLimitRecordCacheStore) getRedisClient(backendServerName types.BackendServerName) (*redis.Client, int, *exceptions.Exception) {
	serverNumber, ok := s.backendServerNameToRedisIndex[backendServerName]
	if !ok {
		return nil, 0, exceptions.Cache.BackendServerNameNotReferenced(types.ValidCachePurpose_RateLimite.String())
	}

	redisClient, ok := s.redisClientMap[serverNumber]
	if !ok || redisClient == nil {
		return nil, 0, exceptions.Cache.RedisServerNumberNotFound()
	}

	return redisClient, serverNumber, nil
}

func formatRateLimitRecordKey(identifier string) string {
	return fmt.Sprintf("%s:%s", types.ValidCachePurpose_RateLimite.String(), identifier)
}

func (s *RateLimitRecordCacheStore) calculateExpiration(identifier string, windowStart time.Time, windowDuration time.Duration) time.Duration {
	baseExpiration := windowStart.Add(windowDuration).Sub(time.Now())
	if baseExpiration < 0 {
		return 1
	}

	hash := fnv.New32a()
	_, _ = hash.Write([]byte(identifier))
	random := rand.New(rand.NewSource(int64(hash.Sum32())))

	return baseExpiration + time.Duration(random.Int63n(int64(s.jitterMaxOffset)))
}

/* ============================== CRUD Method ============================== */

func (s *RateLimitRecordCacheStore) Get(
	identifier string,
	backendServerName types.BackendServerName,
) (*RateLimitRecordCache, *exceptions.Exception) {
	redisClient, serverNumber, exception := s.getRedisClient(backendServerName)
	if exception != nil {
		return nil, exception
	}

	cacheString, err := redisClient.Get(formatRateLimitRecordKey(identifier)).Result()
	if err != nil {
		return nil, exceptions.Cache.NotFound(types.ValidCachePurpose_RateLimite.String()).WithOrigin(err)
	}

	var rateLimitRecordCache RateLimitRecordCache
	if err := json.Unmarshal([]byte(cacheString), &rateLimitRecordCache); err != nil {
		return nil, exceptions.Cache.FailedToConvertJsonToStruct().WithOrigin(err)
	}

	logs.NotezyLogger.Debug(context.Background(), fmt.Sprintf("Successfully got cached rate limit record from server %d", serverNumber))
	return &rateLimitRecordCache, nil
}

func (s *RateLimitRecordCacheStore) Set(
	identifier string,
	backendServerName types.BackendServerName,
	rateLimitRecordCache RateLimitRecordCache,
) *exceptions.Exception {
	redisClient, serverNumber, exception := s.getRedisClient(backendServerName)
	if exception != nil {
		return exception
	}

	value, err := json.Marshal(rateLimitRecordCache)
	if err != nil {
		return exceptions.Cache.FailedToConvertStructToJson().WithOrigin(err)
	}

	expiresIn := s.calculateExpiration(identifier, rateLimitRecordCache.WindowStartTime, rateLimitRecordCache.WindowDuration)
	if err := redisClient.Set(formatRateLimitRecordKey(identifier), string(value), expiresIn).Err(); err != nil {
		return exceptions.Cache.FailedToCreate(types.ValidCachePurpose_RateLimite.String()).WithOrigin(err)
	}

	logs.NotezyLogger.Debug(context.Background(), fmt.Sprintf("Successfully set cached rate limit record in server %d", serverNumber))
	return nil
}

func (s *RateLimitRecordCacheStore) Update(
	identifier string,
	backendServerName types.BackendServerName,
	input cacheinputs.SynchronizeRateLimitRecordCacheInput,
) *exceptions.Exception {
	rateLimitRecordCache, exception := s.Get(identifier, backendServerName)
	if exception != nil {
		return exception
	}

	if (!input.IsAccumulated && rateLimitRecordCache.NumOfTokens < input.NumOfChangingTokens) || rateLimitRecordCache.NumOfTokens < 0 {
		return exceptions.Auth.InvalidRateLimitTokenCount()
	}

	if input.IsAccumulated {
		rateLimitRecordCache.NumOfTokens += input.NumOfChangingTokens
	} else {
		rateLimitRecordCache.NumOfTokens -= input.NumOfChangingTokens
	}
	rateLimitRecordCache.UpdatedAt = time.Now()

	redisClient, serverNumber, exception := s.getRedisClient(backendServerName)
	if exception != nil {
		return exception
	}

	value, err := json.Marshal(rateLimitRecordCache)
	if err != nil {
		return exceptions.Cache.FailedToConvertStructToJson().WithOrigin(err)
	}

	expiresIn := s.calculateExpiration(identifier, rateLimitRecordCache.WindowStartTime, rateLimitRecordCache.WindowDuration)
	if err := redisClient.Set(formatRateLimitRecordKey(identifier), string(value), expiresIn).Err(); err != nil {
		return exceptions.Cache.FailedToUpdate(types.ValidCachePurpose_RateLimite.String()).WithOrigin(err)
	}

	logs.NotezyLogger.Debug(context.Background(), fmt.Sprintf("Successfully updated cached rate limit record in server %d", serverNumber))
	return nil
}

func (s *RateLimitRecordCacheStore) Delete(
	identifier string,
	backendServerName types.BackendServerName,
) *exceptions.Exception {
	redisClient, serverNumber, exception := s.getRedisClient(backendServerName)
	if exception != nil {
		return exception
	}

	if err := redisClient.Del(formatRateLimitRecordKey(identifier)).Err(); err != nil {
		return exceptions.Cache.FailedToDelete(types.ValidCachePurpose_RateLimite.String()).WithOrigin(err)
	}

	logs.NotezyLogger.Debug(context.Background(), fmt.Sprintf("Successfully deleted cached rate limit record from server %d", serverNumber))
	return nil
}

/* ============================== Batch Method ============================== */

func (s *RateLimitRecordCacheStore) BatchSynchronize(
	inputs []cacheinputs.BatchSynchronizeRateLimitRecordCacheInput,
	backendServerName types.BackendServerName,
) *exceptions.Exception {
	if len(inputs) == 0 {
		return nil
	}

	redisClient, serverNumber, exception := s.getRedisClient(backendServerName)
	if exception != nil {
		return exception
	}

	keys := make([]interface{}, 0, len(inputs))
	arguments := make([]interface{}, 0, len(inputs)*s.batchSynchronizeFunctionArgvPerKey)
	for _, input := range inputs {
		keys = append(keys, formatRateLimitRecordKey(input.Identifier))
		arguments = append(arguments, input.Input.NumOfChangingTokens, input.Input.IsAccumulated)
	}

	command := []interface{}{
		"FCALL",
		redislibraries.BatchSynchronizeRateLimitRecordByFormattedKeysFunction,
		len(keys),
	}
	command = append(command, keys...)
	command = append(command, arguments...)
	if _, err := redisClient.Do(command...).Result(); err != nil {
		return exceptions.Cache.FailedToUpdate(types.ValidCachePurpose_RateLimite.String()).WithOrigin(err)
	}

	logs.NotezyLogger.Debug(context.Background(), fmt.Sprintf("Successfully batch synchronized rate limit records in server %d", serverNumber))
	return nil
}

func (s *RateLimitRecordCacheStore) BatchDelete(
	identifiers []string,
	backendServerName types.BackendServerName,
) *exceptions.Exception {
	if len(identifiers) == 0 {
		return nil
	}

	redisClient, serverNumber, exception := s.getRedisClient(backendServerName)
	if exception != nil {
		return exception
	}

	keys := make([]interface{}, 0, len(identifiers))
	for _, identifier := range identifiers {
		keys = append(keys, formatRateLimitRecordKey(identifier))
	}

	command := []interface{}{
		"FCALL",
		redislibraries.BatchDeleteRateLimitRecordByFormattedKeysFunction,
		len(keys),
	}
	command = append(command, keys...)
	if _, err := redisClient.Do(command...).Result(); err != nil {
		return exceptions.Cache.FailedToDelete(types.ValidCachePurpose_RateLimite.String()).WithOrigin(err)
	}

	logs.NotezyLogger.Debug(context.Background(), fmt.Sprintf("Successfully batch deleted rate limit records from server %d", serverNumber))
	return nil
}
