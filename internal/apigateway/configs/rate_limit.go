package config

import (
	"time"

	rate "golang.org/x/time/rate"

	platformredis "github.com/HiIamJeff67/notezy-backend/shared/platform/redis"

	ratelimitrecord "github.com/HiIamJeff67/notezy-backend/internal/apigateway/data/cache/ratelimitrecord"
)

type RateLimitConfig struct {
	RateLimit         rate.Limit
	Burst             int
	UserLimit         int32
	WindowDuration    time.Duration
	BackendServerName platformredis.BackendServerName
	CacheClient       *ratelimitrecord.RateLimitRecordCacheClient

	RequestFrequencyExtraCapacity        int
	MinIntervalTimeOfLastRequest         time.Duration
	SynchronizationToWindowDurationRatio int64
	MinSynchronizationInterval           time.Duration
}

func DefaultUnauthorizedRateLimitConfig() RateLimitConfig {
	return RateLimitConfig{
		RateLimit:                            100,
		Burst:                                10,
		UserLimit:                            1000,
		WindowDuration:                       time.Minute,
		BackendServerName:                    platformredis.BackendServerName_EastAsia,
		RequestFrequencyExtraCapacity:        2,
		MinIntervalTimeOfLastRequest:         time.Microsecond,
		SynchronizationToWindowDurationRatio: 10,
		MinSynchronizationInterval:           time.Second,
	}
}
