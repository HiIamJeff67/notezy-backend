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

var defaultAuthorizedConfig = Config{
	RateLimit:         rate.Limit(100),
	Burst:             20,
	UserLimit:         3000,
	WindowDuration:    time.Minute,
	BackendServerName: platformredis.BackendServerName_EastAsia,
}

var defaultUnauthorizedConfig = Config{
	RateLimit:         rate.Limit(100),
	Burst:             10,
	UserLimit:         1000,
	WindowDuration:    time.Minute,
	BackendServerName: platformredis.BackendServerName_EastAsia,
}

func DefaultAuthorizedConfig() Config {
	return defaultAuthorizedConfig
}

func DefaultUnauthorizedConfig() Config {
	return defaultUnauthorizedConfig
}
