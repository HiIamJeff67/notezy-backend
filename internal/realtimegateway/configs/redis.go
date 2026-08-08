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
	RealtimeLeaseServerStart   int
	RealtimeLeaseServerSize    int
}

func loadRedisConfig() (RedisConfig, error) {
	rateLimitRecordServerStart, err := parseNonNegativeInteger("REALTIME_GATEWAY_RATE_LIMIT_RECORD_CACHE_SERVER_START")
	if err != nil {
		return RedisConfig{}, err
	}
	rateLimitRecordServerSize, err := strconv.Atoi(strings.TrimSpace(os.Getenv("REALTIME_GATEWAY_RATE_LIMIT_RECORD_CACHE_SERVER_SIZE")))
	if err != nil || rateLimitRecordServerSize < 4 {
		return RedisConfig{}, fmt.Errorf("REALTIME_GATEWAY_RATE_LIMIT_RECORD_CACHE_SERVER_SIZE must be at least 4")
	}
	realtimeLeaseServerStart, err := parseNonNegativeInteger("REALTIME_GATEWAY_REALTIME_LEASE_CACHE_SERVER_START")
	if err != nil {
		return RedisConfig{}, err
	}
	realtimeLeaseServerSize, err := strconv.Atoi(strings.TrimSpace(os.Getenv("REALTIME_GATEWAY_REALTIME_LEASE_CACHE_SERVER_SIZE")))
	if err != nil || realtimeLeaseServerSize <= 0 {
		return RedisConfig{}, fmt.Errorf("REALTIME_GATEWAY_REALTIME_LEASE_CACHE_SERVER_SIZE must be a positive integer")
	}

	return RedisConfig{
		RateLimitRecordServerStart: rateLimitRecordServerStart,
		RateLimitRecordServerSize:  rateLimitRecordServerSize,
		RealtimeLeaseServerStart:   realtimeLeaseServerStart,
		RealtimeLeaseServerSize:    realtimeLeaseServerSize,
	}, nil
}

func parseNonNegativeInteger(name string) (int, error) {
	value, err := strconv.Atoi(strings.TrimSpace(os.Getenv(name)))
	if err != nil || value < 0 {
		return 0, fmt.Errorf("%s must be a non-negative integer", name)
	}

	return value, nil
}
