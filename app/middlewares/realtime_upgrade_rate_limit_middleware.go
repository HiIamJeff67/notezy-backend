package middlewares

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"

	configs "github.com/HiIamJeff67/notezy-backend/app/configs"
	logs "github.com/HiIamJeff67/notezy-backend/app/monitor/logs"
	ratelimit "github.com/HiIamJeff67/notezy-backend/shared/lib/ratelimit"
)

var realtimeUpgradeRateLimiter *ratelimit.HybridRateLimiter

func InitRealtimeUpgradeRateLimiter(config configs.RateLimitConfig) {
	if realtimeUpgradeRateLimiter != nil {
		realtimeUpgradeRateLimiter.Stop()
	}

	realtimeUpgradeRateLimiter = ratelimit.NewHybridRateLimiter(
		config.RateLimit,
		config.Burst,
		config.UserLimit,
		config.WindowDuration,
		config.BackendServerName,
		false,
	)

	logs.NotezyLogger.Info(context.Background(), fmt.Sprintf("Realtime upgrade rate limiter initialized with rate: %v, burst: %d, user limit: %d, window: %v", config.RateLimit, config.Burst, config.UserLimit, config.WindowDuration))
}

func RealtimeUpgradeRateLimitMiddleware(config ...configs.RateLimitConfig) gin.HandlerFunc {
	cfg := configs.DefaultRealtimeUpgradeRateLimitConfig
	if len(config) > 0 {
		cfg = config[0]
	}

	if realtimeUpgradeRateLimiter == nil {
		InitRealtimeUpgradeRateLimiter(cfg)
	}

	return func(ctx *gin.Context) {
		fingerprint := ctx.ClientIP()

		allowed, remaining := realtimeUpgradeRateLimiter.AllowByFingerprint(fingerprint)
		if !allowed {
			setRateLimitHeaders(ctx, remaining, realtimeUpgradeRateLimiter)
			logs.NotezyLogger.Debug(ctx.Request.Context(), fmt.Sprintf("Realtime upgrade rate limit exceeded for fingerprint: %s", fingerprint))
			ctx.AbortWithStatus(429)

			return
		}

		setRateLimitHeaders(ctx, remaining, realtimeUpgradeRateLimiter)

		ctx.Next()
	}
}

func StopRealtimeUpgradeRateLimiter() {
	if realtimeUpgradeRateLimiter != nil {
		realtimeUpgradeRateLimiter.Stop()
		realtimeUpgradeRateLimiter = nil
		logs.NotezyLogger.Info(context.Background(), "Realtime upgrade rate limiter stopped")
	}
}
