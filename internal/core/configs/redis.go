package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type RedisConfig struct {
	UserDataCacheServerStart int
	UserDataCacheServerSize  int
	UserDataCacheExpiresIn   time.Duration
}

func loadRedisConfig() (RedisConfig, error) {
	serverStart, err := strconv.Atoi(strings.TrimSpace(os.Getenv("CORE_USER_DATA_CACHE_SERVER_START")))
	if err != nil || serverStart < 0 {
		return RedisConfig{}, fmt.Errorf("CORE_USER_DATA_CACHE_SERVER_START must be a non-negative integer")
	}
	serverSize, err := strconv.Atoi(strings.TrimSpace(os.Getenv("CORE_USER_DATA_CACHE_SERVER_SIZE")))
	if err != nil || serverSize <= 0 {
		return RedisConfig{}, fmt.Errorf("CORE_USER_DATA_CACHE_SERVER_SIZE must be a positive integer")
	}
	cacheExpiresIn, err := time.ParseDuration(strings.TrimSpace(os.Getenv("CORE_USER_DATA_CACHE_EXPIRES_IN")))
	if err != nil || cacheExpiresIn <= 0 {
		return RedisConfig{}, fmt.Errorf("CORE_USER_DATA_CACHE_EXPIRES_IN must be a positive Go duration")
	}

	return RedisConfig{
		UserDataCacheServerStart: serverStart,
		UserDataCacheServerSize:  serverSize,
		UserDataCacheExpiresIn:   cacheExpiresIn,
	}, nil
}
