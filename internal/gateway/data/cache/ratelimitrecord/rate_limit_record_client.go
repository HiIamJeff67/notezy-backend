package ratelimitrecord

import (
	"context"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"math/rand"
	"net/http"
	"time"

	"github.com/go-redis/redis"

	exceptions "github.com/HiIamJeff67/notezy-backend/internal/exceptions"
	cacheinputs "github.com/HiIamJeff67/notezy-backend/internal/gateway/data/cache/ratelimitrecord/inputs"
	redislibraries "github.com/HiIamJeff67/notezy-backend/internal/gateway/data/cache/ratelimitrecord/libraries"
	configs "github.com/HiIamJeff67/notezy-backend/internal/platform/config"
	logs "github.com/HiIamJeff67/notezy-backend/internal/platform/observability/logs"
	platformredis "github.com/HiIamJeff67/notezy-backend/internal/platform/redis"
	types "github.com/HiIamJeff67/notezy-backend/shared/types"
)

type RateLimitRecordCache struct {
	NumOfTokens     int32         `json:"numOfTokens"`
	WindowStartTime time.Time     `json:"windowStartTime"`
	WindowDuration  time.Duration `json:"windowDuration"`
	UpdatedAt       time.Time     `json:"updatedAt"`
}

type RateLimitRecordCacheClient struct {
	Range                         types.Range[int, int]
	MaxServerNumber               int
	backendServerNameToRedisIndex map[configs.BackendServerName]int

	jitterMaxOffset                    time.Duration
	batchSynchronizeFunctionArgvPerKey int
}

/* ============================== Constructor ============================== */

func NewRateLimitRecordCacheClient() *RateLimitRecordCacheClient {
	rangeValue := types.Range[int, int]{Start: 4, Size: 4}

	return &RateLimitRecordCacheClient{
		Range:           rangeValue,
		MaxServerNumber: rangeValue.Start + rangeValue.Size - 1,
		backendServerNameToRedisIndex: map[configs.BackendServerName]int{
			configs.BackendServerName_EastAsia:    4,
			configs.BackendServerName_EastAmerica: 5,
			configs.BackendServerName_WestAmerica: 6,
			configs.BackendServerName_WestEurope:  7,
		},

		jitterMaxOffset:                    5 * time.Second,
		batchSynchronizeFunctionArgvPerKey: 2,
	}
}

/* ============================== Auxiliary Methods ============================== */

func (s *RateLimitRecordCacheClient) getRedisClient(backendServerName configs.BackendServerName) (*redis.Client, int, *exceptions.Exception) {
	serverNumber, ok := s.backendServerNameToRedisIndex[backendServerName]
	if !ok {
		return nil, 0, exceptions.New(
			"BackendServerNameNotMapped",
			"Cache",
			"GetRedisClient",
			"Rate limit cache backend server is not configured",
			http.StatusInternalServerError,
			true,
		)
	}

	cacheStore, err := platformredis.GetRedisCacheStore(serverNumber)
	if err != nil {
		return nil, 0, exceptions.New(
			"CacheClientUnavailable",
			"Cache",
			"GetRedisClient",
			"Rate limit cache client is unavailable",
			http.StatusInternalServerError,
			true,
		).WithOrigin(err)
	}
	rateLimitRecordCacheStore, ok := cacheStore.(*RateLimitRecordCacheStore)
	if !ok || rateLimitRecordCacheStore.redisClient == nil {
		return nil, 0, exceptions.New(
			"CacheClientUnavailable",
			"Cache",
			"GetRedisClient",
			"Rate limit cache client is unavailable",
			http.StatusInternalServerError,
			true,
		)
	}

	return rateLimitRecordCacheStore.redisClient, serverNumber, nil
}

func formatRateLimitRecordKey(identifier string) string {
	return fmt.Sprintf("%s:%s", platformredis.CachePurpose_RateLimit.String(), identifier)
}

func (s *RateLimitRecordCacheClient) calculateExpiration(identifier string, windowStart time.Time, windowDuration time.Duration) time.Duration {
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

func (s *RateLimitRecordCacheClient) Get(
	identifier string,
	backendServerName configs.BackendServerName,
) (*RateLimitRecordCache, *exceptions.Exception) {
	redisClient, serverNumber, exception := s.getRedisClient(backendServerName)
	if exception != nil {
		return nil, exception
	}

	cacheString, err := redisClient.Get(formatRateLimitRecordKey(identifier)).Result()
	if err != nil {
		return nil, exceptions.New(
			"NotFound",
			"Cache",
			"GetRateLimitRecord",
			"Rate limit record was not found",
			http.StatusNotFound,
			true,
		).WithOrigin(err)
	}

	var rateLimitRecordCache RateLimitRecordCache
	if err := json.Unmarshal([]byte(cacheString), &rateLimitRecordCache); err != nil {
		return nil, exceptions.New(
			"DeserializationFailed",
			"Cache",
			"GetRateLimitRecord",
			"Failed to decode the rate limit record",
			http.StatusInternalServerError,
			true,
		).WithOrigin(err)
	}

	logs.NotezyLogger.Debug(context.Background(), fmt.Sprintf("Successfully got cached rate limit record from server %d", serverNumber))
	return &rateLimitRecordCache, nil
}

func (s *RateLimitRecordCacheClient) Set(
	identifier string,
	backendServerName configs.BackendServerName,
	rateLimitRecordCache RateLimitRecordCache,
) *exceptions.Exception {
	redisClient, serverNumber, exception := s.getRedisClient(backendServerName)
	if exception != nil {
		return exception
	}

	value, err := json.Marshal(rateLimitRecordCache)
	if err != nil {
		return exceptions.New(
			"SerializationFailed",
			"Cache",
			"SetRateLimitRecord",
			"Failed to encode the rate limit record",
			http.StatusInternalServerError,
			true,
		).WithOrigin(err)
	}

	expiresIn := s.calculateExpiration(identifier, rateLimitRecordCache.WindowStartTime, rateLimitRecordCache.WindowDuration)
	if err := redisClient.Set(formatRateLimitRecordKey(identifier), string(value), expiresIn).Err(); err != nil {
		return exceptions.New(
			"FailedToCreate",
			"Cache",
			"SetRateLimitRecord",
			"Failed to store the rate limit record",
			http.StatusInternalServerError,
			true,
		).WithOrigin(err)
	}

	logs.NotezyLogger.Debug(context.Background(), fmt.Sprintf("Successfully set cached rate limit record in server %d", serverNumber))
	return nil
}

func (s *RateLimitRecordCacheClient) Update(
	identifier string,
	backendServerName configs.BackendServerName,
	input cacheinputs.SynchronizeRateLimitRecordCacheInput,
) *exceptions.Exception {
	rateLimitRecordCache, exception := s.Get(identifier, backendServerName)
	if exception != nil {
		return exception
	}

	if (!input.IsAccumulated && rateLimitRecordCache.NumOfTokens < input.NumOfChangingTokens) || rateLimitRecordCache.NumOfTokens < 0 {
		return exceptions.New(
			"InvalidRateLimitTokenCount",
			"RateLimit",
			"UpdateRateLimitRecord",
			"Rate limit token count is invalid",
			http.StatusInternalServerError,
			true,
		)
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
		return exceptions.New(
			"SerializationFailed",
			"Cache",
			"UpdateRateLimitRecord",
			"Failed to encode the rate limit record",
			http.StatusInternalServerError,
			true,
		).WithOrigin(err)
	}

	expiresIn := s.calculateExpiration(identifier, rateLimitRecordCache.WindowStartTime, rateLimitRecordCache.WindowDuration)
	if err := redisClient.Set(formatRateLimitRecordKey(identifier), string(value), expiresIn).Err(); err != nil {
		return exceptions.New(
			"FailedToUpdate",
			"Cache",
			"UpdateRateLimitRecord",
			"Failed to update the rate limit record",
			http.StatusInternalServerError,
			true,
		).WithOrigin(err)
	}

	logs.NotezyLogger.Debug(context.Background(), fmt.Sprintf("Successfully updated cached rate limit record in server %d", serverNumber))
	return nil
}

func (s *RateLimitRecordCacheClient) Delete(
	identifier string,
	backendServerName configs.BackendServerName,
) *exceptions.Exception {
	redisClient, serverNumber, exception := s.getRedisClient(backendServerName)
	if exception != nil {
		return exception
	}

	if err := redisClient.Del(formatRateLimitRecordKey(identifier)).Err(); err != nil {
		return exceptions.New(
			"FailedToDelete",
			"Cache",
			"DeleteRateLimitRecord",
			"Failed to delete the rate limit record",
			http.StatusInternalServerError,
			true,
		).WithOrigin(err)
	}

	logs.NotezyLogger.Debug(context.Background(), fmt.Sprintf("Successfully deleted cached rate limit record from server %d", serverNumber))
	return nil
}

/* ============================== Batch Method ============================== */

func (s *RateLimitRecordCacheClient) BatchSynchronize(
	inputs []cacheinputs.BatchSynchronizeRateLimitRecordCacheInput,
	backendServerName configs.BackendServerName,
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
		return exceptions.New(
			"FailedToUpdate",
			"Cache",
			"BatchSynchronizeRateLimitRecords",
			"Failed to synchronize rate limit records",
			http.StatusInternalServerError,
			true,
		).WithOrigin(err)
	}

	logs.NotezyLogger.Debug(context.Background(), fmt.Sprintf("Successfully batch synchronized rate limit records in server %d", serverNumber))
	return nil
}

func (s *RateLimitRecordCacheClient) BatchDelete(
	identifiers []string,
	backendServerName configs.BackendServerName,
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
		return exceptions.New(
			"FailedToDelete",
			"Cache",
			"BatchDeleteRateLimitRecords",
			"Failed to delete rate limit records",
			http.StatusInternalServerError,
			true,
		).WithOrigin(err)
	}

	logs.NotezyLogger.Debug(context.Background(), fmt.Sprintf("Successfully batch deleted rate limit records from server %d", serverNumber))
	return nil
}
