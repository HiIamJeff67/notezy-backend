package realtimelease

import (
	"context"
	"fmt"

	"github.com/go-redis/redis"

	platformredis "github.com/HiIamJeff67/notezy-backend/shared/platform/redis"
)

type RealtimeLeaseCacheStore struct {
	databaseNumber int
	redisClient    *redis.Client
}

func NewRealtimeLeaseCacheStore(
	databaseNumber int,
	redisClient *redis.Client,
) *RealtimeLeaseCacheStore {
	return &RealtimeLeaseCacheStore{
		databaseNumber: databaseNumber,
		redisClient:    redisClient,
	}
}

func Register(
	ctx context.Context,
	clientManager *platformredis.ClientManager,
	cacheClient *RealtimeLeaseCacheClient,
) error {
	if err := clientManager.ConnectAll(cacheClient.Range); err != nil {
		return err
	}

	cacheStores := make([]platformredis.RedisCacheStore, 0, cacheClient.Range.Size)
	for databaseNumber := cacheClient.Range.Start; databaseNumber < cacheClient.Range.Start+cacheClient.Range.Size; databaseNumber++ {
		redisClient, exists := clientManager.Client(databaseNumber)
		if !exists {
			return fmt.Errorf("realtime Redis client for database %d is unavailable", databaseNumber)
		}
		cacheStores = append(cacheStores, NewRealtimeLeaseCacheStore(databaseNumber, redisClient))
	}

	return platformredis.RegisterCacheStores(ctx, cacheStores...)
}

func (s *RealtimeLeaseCacheStore) DatabaseNumber() int {
	return s.databaseNumber
}

func (s *RealtimeLeaseCacheStore) Initialize(_ context.Context) error {
	return nil
}
