package ratelimiter

import (
	"time"

	"golang.org/x/time/rate"

	platformredis "github.com/HiIamJeff67/notezy-backend/internal/platform/redis"
)

type Config struct {
	RateLimit         rate.Limit
	Burst             int
	UserLimit         int32
	WindowDuration    time.Duration
	BackendServerName platformredis.BackendServerName
}

var defaultUpgradeConfig = Config{
	RateLimit:         rate.Limit(5),
	Burst:             10,
	UserLimit:         60,
	WindowDuration:    time.Minute,
	BackendServerName: platformredis.BackendServerName_EastAsia,
}

func DefaultUpgradeConfig() Config {
	return defaultUpgradeConfig
}
