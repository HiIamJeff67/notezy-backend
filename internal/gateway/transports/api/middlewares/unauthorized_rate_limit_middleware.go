package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	exceptions "github.com/HiIamJeff67/notezy-backend/contracts/types/exceptions"

	exceptionwriter "github.com/HiIamJeff67/notezy-backend/shared/util/exceptionwriter"

	logs "github.com/HiIamJeff67/notezy-backend/shared/platform/observability/logs"

	ratelimit "github.com/HiIamJeff67/notezy-backend/internal/gateway/ratelimit"
)

var unauthorizedRateLimiter *ratelimit.HybridRateLimiter

func InitUnauthorizedRateLimiter(config ratelimit.Config) {
	if unauthorizedRateLimiter != nil {
		unauthorizedRateLimiter.Stop()
	}

	unauthorizedRateLimiter = ratelimit.NewHybridRateLimiter(config, false)

	logs.NotezyLogger.Info(context.Background(), fmt.Sprintf("Unauthorized rate limiter initialized with rate: %v, burst: %d, user limit: %d, window: %v", config.RateLimit, config.Burst, config.UserLimit, config.WindowDuration))
}

func UnauthorizedRateLimitMiddleware(config ...ratelimit.Config) gin.HandlerFunc {
	cfg := ratelimit.DefaultUnauthorizedConfig()
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
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.New(
				"PermissionDeniedDueToTooManyRequests",
				"Auth",
				"Authorize",
				"Too many requests; please wait before trying again",
				http.StatusTooManyRequests,
			), ctx, "server.responses.failed.rateLimit")
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
