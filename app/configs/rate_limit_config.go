package configs

import (
	"time"

	"golang.org/x/time/rate"

	types "notezy-backend/shared/types"
)

type RateLimitConfig struct {
	RateLimit         rate.Limit
	Burst             int
	UserLimit         int32
	WindowDuration    time.Duration
	BackendServerName types.BackendServerName
}

var DefaultAuthorizedRateLimitConfig = RateLimitConfig{
	RateLimit:         rate.Limit(100),                  // 100 requests/second
	Burst:             20,                               // allowed 20 additional requests/second for burst
	UserLimit:         3000,                             // 300 requests/each life time of the bucket (= 300 requests/`WindowDuration`) for each users
	WindowDuration:    time.Minute,                      // 1 minutes to reset the bucket
	BackendServerName: types.BackendServerName_EastAsia, // the current server
}

var DefaultUnauthorizedRateLimitConfig = RateLimitConfig{
	RateLimit:         rate.Limit(10),                   // 10 requests/second
	Burst:             5,                                // allowed 20 additional requests/second for burst
	UserLimit:         300,                              // 300 requests/each life time of the bucket (= 300 requests/`WindowDuration`) for each users
	WindowDuration:    time.Minute,                      // 1 minutes to reset the bucket
	BackendServerName: types.BackendServerName_EastAsia, // the current server
}
