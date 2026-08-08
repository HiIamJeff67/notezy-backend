package config

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type UserDataCacheConfig struct {
	CacheExpiresIn time.Duration
}

func loadUserDataCacheConfig() (UserDataCacheConfig, error) {
	cacheExpiresIn, err := time.ParseDuration(strings.TrimSpace(os.Getenv("CORE_USER_DATA_CACHE_EXPIRES_IN")))
	if err != nil || cacheExpiresIn <= 0 {
		return UserDataCacheConfig{}, fmt.Errorf("CORE_USER_DATA_CACHE_EXPIRES_IN must be a positive Go duration")
	}

	return UserDataCacheConfig{
		CacheExpiresIn: cacheExpiresIn,
	}, nil
}
