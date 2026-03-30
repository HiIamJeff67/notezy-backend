package middlewares

import (
	configs "notezy-backend/app/configs"
	exceptions "notezy-backend/app/exceptions"
	ratelimiter "notezy-backend/app/lib/ratelimiter"
	logs "notezy-backend/app/monitor/logs"
	metrics "notezy-backend/app/monitor/metrics"
	traces "notezy-backend/app/monitor/traces"

	"github.com/gin-gonic/gin"
)

var unauthorizedRateLimiter *ratelimiter.HybridRateLimiter

func InitUnauthorizedRateLimiter(config configs.RateLimitConfig) {
	if unauthorizedRateLimiter != nil {
		unauthorizedRateLimiter.Stop()
	}

	unauthorizedRateLimiter = ratelimiter.NewHybridRateLimiter(
		config.RateLimit,
		config.Burst,
		config.UserLimit,
		config.WindowDuration,
		config.BackendServerName,
		false,
	)

	logs.FInfo(traces.GetTrace(0).FileLineString(),
		"Unauthorized rate limiter initialized with rate: %v, burst: %d, user limit: %d, window: %v",
		config.RateLimit, config.Burst, config.UserLimit, config.WindowDuration)
}

func UnauthorizedRateLimitMiddleware(config ...configs.RateLimitConfig) gin.HandlerFunc {
	cfg := configs.DefaultUnauthorizedRateLimitConfig
	if len(config) > 0 {
		cfg = config[0]
	}

	if unauthorizedRateLimiter == nil {
		InitAuthorizedRateLimiter(cfg)
	}

	return func(ctx *gin.Context) {
		fingerprint := getClientFingerprint(ctx)

		allowed, remaining := unauthorizedRateLimiter.AllowByFingerprint(fingerprint)
		if !allowed {
			setRateLimitHeaders(ctx, remaining, unauthorizedRateLimiter)
			logs.FDebug(traces.GetTrace(0).FileLineString(), "Rate limit exceeded for fingerprint: %s", fingerprint)
			exceptions.Auth.PermissionDeniedDueToTooManyRequests().Log().SafelyAbortAndResponseWithJSON(ctx, metrics.MetricNames.Server.Responses.Failed.RateLimit)
			return
		}

		setRateLimitHeaders(ctx, remaining, unauthorizedRateLimiter)

		ctx.Next()
	}
}

func getClientFingerprint(c *gin.Context) string {
	// TODO: use other complex stuff or algorithm or even the machine learning model to generate or get the fingerprint of each clients
	return c.ClientIP()
}

func StopUnauthorizedRateLimiter() {
	if unauthorizedRateLimiter != nil {
		unauthorizedRateLimiter.Stop()
		unauthorizedRateLimiter = nil
		logs.FInfo(traces.GetTrace(0).FileLineString(), "Unauthoized rate limiter stopped")
	}
}
