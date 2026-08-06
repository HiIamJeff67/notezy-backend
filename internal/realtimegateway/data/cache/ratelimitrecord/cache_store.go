package ratelimitrecord

import (
	"context"
	"fmt"

	"github.com/go-redis/redis"

	platformredis "github.com/HiIamJeff67/notezy-backend/shared/platform/redis"

	redislibraries "github.com/HiIamJeff67/notezy-backend/internal/realtimegateway/data/cache/ratelimitrecord/libraries"
)

type RateLimitRecordCacheStore struct {
	databaseNumber int
	redisClient    *redis.Client
}

func NewRateLimitRecordCacheStore(
	databaseNumber int,
	redisClient *redis.Client,
) *RateLimitRecordCacheStore {
	return &RateLimitRecordCacheStore{
		databaseNumber: databaseNumber,
		redisClient:    redisClient,
	}
}

func Register(
	ctx context.Context,
	clientManager *platformredis.ClientManager,
	cacheClient *RateLimitRecordCacheClient,
) error {
	if err := clientManager.ConnectAll(cacheClient.Range); err != nil {
		return err
	}

	cacheStores := make([]platformredis.RedisCacheStore, 0, cacheClient.Range.Size)
	for databaseNumber := cacheClient.Range.Start; databaseNumber < cacheClient.Range.Start+cacheClient.Range.Size; databaseNumber++ {
		redisClient, exists := clientManager.Client(databaseNumber)
		if !exists {
			return fmt.Errorf("Redis client for database %d is unavailable", databaseNumber)
		}
		cacheStores = append(cacheStores, NewRateLimitRecordCacheStore(databaseNumber, redisClient))
	}

	return platformredis.RegisterCacheStores(ctx, cacheStores...)
}

func (s *RateLimitRecordCacheStore) DatabaseNumber() int {
	return s.databaseNumber
}

func (s *RateLimitRecordCacheStore) Initialize(_ context.Context) error {
	return s.redisClient.Do(
		"FUNCTION",
		"LOAD",
		"REPLACE",
		redislibraries.RateLimitRecordLibraryContent,
	).Err()
}
