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
	rateLimitRecordServerSize, err := parseMinimumInteger("REALTIME_GATEWAY_RATE_LIMIT_RECORD_CACHE_SERVER_SIZE", 4)
	if err != nil {
		return RedisConfig{}, err
	}
	realtimeLeaseServerStart, err := parseNonNegativeInteger("REALTIME_GATEWAY_REALTIME_LEASE_CACHE_SERVER_START")
	if err != nil {
		return RedisConfig{}, err
	}
	realtimeLeaseServerSize, err := parsePositiveInteger("REALTIME_GATEWAY_REALTIME_LEASE_CACHE_SERVER_SIZE")
	if err != nil {
		return RedisConfig{}, err
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

func parsePositiveInteger(name string) (int, error) {
	value, err := strconv.Atoi(strings.TrimSpace(os.Getenv(name)))
	if err != nil || value <= 0 {
		return 0, fmt.Errorf("%s must be a positive integer", name)
	}

	return value, nil
}

func parseMinimumInteger(name string, minimum int) (int, error) {
	value, err := strconv.Atoi(strings.TrimSpace(os.Getenv(name)))
	if err != nil || value < minimum {
		return 0, fmt.Errorf("%s must be at least %d", name, minimum)
	}

	return value, nil
}
