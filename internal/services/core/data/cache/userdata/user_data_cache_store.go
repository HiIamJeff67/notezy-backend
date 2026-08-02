package userdata

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-redis/redis"

	platformredis "github.com/HiIamJeff67/notezy-backend/internal/platform/redis"
	redislibraries "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/cache/userdata/libraries"
)

type UserDataCacheStore struct {
	databaseNumber int
	redisClient    *redis.Client
}

func NewUserDataCacheStore(
	databaseNumber int,
	redisClient *redis.Client,
) *UserDataCacheStore {
	return &UserDataCacheStore{
		databaseNumber: databaseNumber,
		redisClient:    redisClient,
	}
}

func Register(
	ctx context.Context,
	clientManager *platformredis.ClientManager,
) error {
	cacheClient := NewUserDataCacheClient()
	if err := clientManager.ConnectAll(cacheClient.Range); err != nil {
		return err
	}

	cacheStores := make([]platformredis.RedisCacheStore, 0, cacheClient.Range.Size)
	for databaseNumber := cacheClient.Range.Start; databaseNumber < cacheClient.Range.Start+cacheClient.Range.Size; databaseNumber++ {
		redisClient, exists := clientManager.Client(databaseNumber)
		if !exists {
			return fmt.Errorf("Redis client for database %d is unavailable", databaseNumber)
		}
		cacheStores = append(cacheStores, NewUserDataCacheStore(databaseNumber, redisClient))
	}

	return platformredis.RegisterCacheStores(ctx, cacheStores...)
}

func (s *UserDataCacheStore) DatabaseNumber() int {
	return s.databaseNumber
}

func (s *UserDataCacheStore) Initialize(_ context.Context) error {
	userQuotaLibraryContent := strings.Join([]string{
		"#!lua name=" + redislibraries.UserQuotaLibrary,
		redislibraries.CheckAndUpdateUserQuotaByFormattedKeyContent,
		redislibraries.BestEffortBatchCheckAndUpdateUserQuotasByFormattedKeysContent,
		redislibraries.AllOrNothingBatchCheckAndUpdateUserQuotasByFormattedKeysContent,
		redislibraries.BestEffortBatchCheckAndUpdateUserQuotasByFormattedKeyContent,
		redislibraries.AllOrNothingBatchCheckAndUpdateUserQuotasByFormattedKeyContent,
	}, "\n\n")

	return s.redisClient.Do(
		"FUNCTION",
		"LOAD",
		"REPLACE",
		userQuotaLibraryContent,
	).Err()
}
