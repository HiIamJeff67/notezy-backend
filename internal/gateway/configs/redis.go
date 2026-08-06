package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type RedisConfig struct {
	RateLimitRecordServerStart int
	RateLimitRecordServerSize  int
}

func loadRedisConfig() (RedisConfig, error) {
	serverStart, err := strconv.Atoi(strings.TrimSpace(os.Getenv("GATEWAY_RATE_LIMIT_RECORD_CACHE_SERVER_START")))
	if err != nil || serverStart < 0 {
		return RedisConfig{}, fmt.Errorf("GATEWAY_RATE_LIMIT_RECORD_CACHE_SERVER_START must be a non-negative integer")
	}
	serverSize, err := strconv.Atoi(strings.TrimSpace(os.Getenv("GATEWAY_RATE_LIMIT_RECORD_CACHE_SERVER_SIZE")))
	if err != nil || serverSize < 4 {
		return RedisConfig{}, fmt.Errorf("GATEWAY_RATE_LIMIT_RECORD_CACHE_SERVER_SIZE must be at least 4")
	}

	return RedisConfig{
		RateLimitRecordServerStart: serverStart,
		RateLimitRecordServerSize:  serverSize,
	}, nil
}
