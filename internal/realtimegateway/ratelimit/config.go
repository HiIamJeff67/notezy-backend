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

var defaultUpgradeConfig = Config{
	RateLimit:                            rate.Limit(5),
	Burst:                                10,
	UserLimit:                            60,
	WindowDuration:                       time.Minute,
	BackendServerName:                    platformredis.BackendServerName_EastAsia,
	CacheServerRange:                     types.Range[int, int]{Start: 9, Size: 4},
	RequestFrequencyExtraCapacity:        2,
	MinIntervalTimeOfLastRequest:         time.Microsecond,
	SynchronizationToWindowDurationRatio: 10,
	MinSynchronizationInterval:           time.Second,
}

func DefaultUpgradeConfig() Config {
	return defaultUpgradeConfig
}
