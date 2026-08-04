package realtimelease

import (
	"context"
	"errors"

	"github.com/go-redis/redis"

	platformredis "github.com/HiIamJeff67/notezy-backend/internal/platform/redis"
	constants "github.com/HiIamJeff67/notezy-backend/shared/constants"
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
) error {
	cacheClient := NewRealtimeLeaseCacheClient()
	if err := clientManager.ConnectAll(cacheClient.Range); err != nil {
		return err
	}

	redisClient, exists := clientManager.Client(constants.RealtimeRedisServerNumber)
	if !exists {
		return errors.New("realtime Redis client is unavailable")
	}

	return platformredis.RegisterCacheStores(
		ctx,
		NewRealtimeLeaseCacheStore(constants.RealtimeRedisServerNumber, redisClient),
	)
}

func (s *RealtimeLeaseCacheStore) DatabaseNumber() int {
	return s.databaseNumber
}

func (s *RealtimeLeaseCacheStore) Initialize(_ context.Context) error {
	return nil
}
