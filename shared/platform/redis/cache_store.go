package redis

import (
	"context"
	"fmt"
)

type RedisCacheStore interface {
	DatabaseNumber() int
	Initialize(ctx context.Context) error
}

var RedisCacheStores = map[int]RedisCacheStore{}

func RegisterCacheStores(
	ctx context.Context,
	cacheStores ...RedisCacheStore,
) error {
	for _, cacheStore := range cacheStores {
		if err := cacheStore.Initialize(ctx); err != nil {
			return err
		}
		RedisCacheStores[cacheStore.DatabaseNumber()] = cacheStore
	}

	return nil
}

func GetRedisCacheStore(databaseNumber int) (RedisCacheStore, error) {
	cacheStore, exists := RedisCacheStores[databaseNumber]
	if !exists {
		return nil, fmt.Errorf("Redis cache store for database %d is unavailable", databaseNumber)
	}

	return cacheStore, nil
}
