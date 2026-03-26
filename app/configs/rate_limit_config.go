package configs

import (
	"time"

	"golang.org/x/time/rate"

	types "notezy-backend/shared/types"
)

type AuthorizedRateLimitConfig struct {
	RateLimit         rate.Limit
	Burst             int
	UserLimit         int32
	WindowDuration    time.Duration
	BackendServerName types.BackendServerName
}
