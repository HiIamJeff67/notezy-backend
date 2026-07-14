package middlewares

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"

	configs "github.com/HiIamJeff67/notezy-backend/app/configs"
	exceptions "github.com/HiIamJeff67/notezy-backend/app/exceptions"
	logs "github.com/HiIamJeff67/notezy-backend/app/monitor/logs"
	ratelimit "github.com/HiIamJeff67/notezy-backend/shared/lib/ratelimit"
)

var unauthorizedRateLimiter *ratelimit.HybridRateLimiter

func InitUnauthorizedRateLimiter(config configs.RateLimitConfig) {
	if unauthorizedRateLimiter != nil {
		unauthorizedRateLimiter.Stop()
	}

	unauthorizedRateLimiter = ratelimit.NewHybridRateLimiter(
		config.RateLimit,
		config.Burst,
		config.UserLimit,
		config.WindowDuration,
		config.BackendServerName,
		false,
	)

	logs.NotezyLogger.Info(context.Background(), fmt.Sprintf("Unauthorized rate limiter initialized with rate: %v, burst: %d, user limit: %d, window: %v", config.RateLimit, config.Burst, config.UserLimit, config.WindowDuration))
}

func UnauthorizedRateLimitMiddleware(config ...configs.RateLimitConfig) gin.HandlerFunc {
	cfg := configs.DefaultUnauthorizedRateLimitConfig
	if len(config) > 0 {
		cfg = config[0]
	}

	if unauthorizedRateLimiter == nil {
		InitUnauthorizedRateLimiter(cfg)
	}

	return func(ctx *gin.Context) {
		fingerprint := getClientFingerprint(ctx)

		allowed, remaining := unauthorizedRateLimiter.AllowByFingerprint(fingerprint)
		if !allowed {
			setRateLimitHeaders(ctx, remaining, unauthorizedRateLimiter)
			logs.NotezyLogger.Debug(ctx.Request.Context(), fmt.Sprintf("Rate limit exceeded for fingerprint: %s", fingerprint))
			exceptions.Auth.PermissionDeniedDueToTooManyRequests().Log().SafelyAbortAndResponseWithJSON(ctx, "server.responses.failed.rateLimit")
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
		logs.NotezyLogger.Info(context.Background(), "Unauthorized rate limiter stopped")
	}
}
