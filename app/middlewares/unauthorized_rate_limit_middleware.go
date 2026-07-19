package middlewares

import (
	"context"
	"fmt"
	"strconv"
	"time"

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
		fingerprint := ctx.ClientIP()

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

func setRateLimitHeaders(ctx *gin.Context, remaining int32, limiter *ratelimit.HybridRateLimiter) {
	// standard information
	ctx.Header("X-RateLimit-Limit", strconv.Itoa(int(limiter.UserLimit)))
	ctx.Header("X-RateLimit-Remaining", strconv.Itoa(int(remaining)))
	ctx.Header("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(limiter.WindowDuration).Unix(), 10))

	// extra information
	ctx.Header("X-RateLimit-Window", limiter.WindowDuration.String())
	ctx.Header("X-RateLimit-Policy", "hybrid-token-bucket")
}

func StopUnauthorizedRateLimiter() {
	if unauthorizedRateLimiter != nil {
		unauthorizedRateLimiter.Stop()
		unauthorizedRateLimiter = nil
		logs.NotezyLogger.Info(context.Background(), "Unauthorized rate limiter stopped")
	}
}
