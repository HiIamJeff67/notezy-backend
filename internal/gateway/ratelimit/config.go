package ratelimiter

import (
	"time"

	"golang.org/x/time/rate"

	types "github.com/HiIamJeff67/notezy-backend/shared/types"

	platformredis "github.com/HiIamJeff67/notezy-backend/shared/platform/redis"
)

type Config struct {
	RateLimit         rate.Limit
	Burst             int
	UserLimit         int32
	WindowDuration    time.Duration
	BackendServerName platformredis.BackendServerName
	CacheServerRange  types.Range[int, int]

	RequestFrequencyExtraCapacity        int
	MinIntervalTimeOfLastRequest         time.Duration
	SynchronizationToWindowDurationRatio int64
	MinSynchronizationInterval           time.Duration
}

var defaultAuthorizedConfig = Config{
	RateLimit:                            rate.Limit(100),
	Burst:                                20,
	UserLimit:                            3000,
	WindowDuration:                       time.Minute,
	BackendServerName:                    platformredis.BackendServerName_EastAsia,
	CacheServerRange:                     types.Range[int, int]{Start: 4, Size: 4},
	RequestFrequencyExtraCapacity:        2,
	MinIntervalTimeOfLastRequest:         time.Microsecond,
	SynchronizationToWindowDurationRatio: 10,
	MinSynchronizationInterval:           time.Second,
}

var defaultUnauthorizedConfig = Config{
	RateLimit:                            rate.Limit(100),
	Burst:                                10,
	UserLimit:                            1000,
	WindowDuration:                       time.Minute,
	BackendServerName:                    platformredis.BackendServerName_EastAsia,
	CacheServerRange:                     types.Range[int, int]{Start: 4, Size: 4},
	RequestFrequencyExtraCapacity:        2,
	MinIntervalTimeOfLastRequest:         time.Microsecond,
	SynchronizationToWindowDurationRatio: 10,
	MinSynchronizationInterval:           time.Second,
}

func DefaultAuthorizedConfig() Config {
	return defaultAuthorizedConfig
}

func DefaultUnauthorizedConfig() Config {
	return defaultUnauthorizedConfig
}
