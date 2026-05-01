package middlewares

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	configs "notezy-backend/app/configs"
	contexts "notezy-backend/app/contexts"
	exceptions "notezy-backend/app/exceptions"
	logs "notezy-backend/app/monitor/logs"
	metrics "notezy-backend/app/monitor/metrics"
	traces "notezy-backend/app/monitor/traces"
	ratelimit "notezy-backend/shared/lib/ratelimit"
	types "notezy-backend/shared/types"
)

var authorizedRateLimiter *ratelimit.HybridRateLimiter // use the hybrid one which including token bucket and cross server request management by redis

func InitAuthorizedRateLimiter(config configs.RateLimitConfig) {
	if authorizedRateLimiter != nil {
		authorizedRateLimiter.Stop()
	}

	authorizedRateLimiter = ratelimit.NewHybridRateLimiter(
		config.RateLimit,
		config.Burst,
		config.UserLimit,
		config.WindowDuration,
		config.BackendServerName,
		true,
	)

	logs.FInfo(traces.GetTrace(0).FileLineString(),
		"Authorized rate limiter initialized with rate: %v, burst: %d, user limit: %d, window: %v",
		config.RateLimit, config.Burst, config.UserLimit, config.WindowDuration)
}

func AuthorizedRateLimitMiddleware(config ...configs.RateLimitConfig) gin.HandlerFunc {
	cfg := configs.DefaultAuthorizedRateLimitConfig
	if len(config) > 0 {
		cfg = config[0]
	}

	if authorizedRateLimiter == nil {
		InitAuthorizedRateLimiter(cfg)
	}

	return func(ctx *gin.Context) {
		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil || userId == nil {
			exceptions.Context.MissPlacingOrWrongMiddlewareOrder(
				"Cannot find the userId, " +
					"please make sure the AuthMiddleware() is placing before the AuthorizedRateLimitMiddleware()",
			).Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}

		allowed, remaining := authorizedRateLimiter.AllowByUserId(*userId)
		if !allowed {
			setRateLimitHeaders(ctx, remaining, authorizedRateLimiter)
			logs.FDebug(traces.GetTrace(0).FileLineString(), "Rate limit exceeded for user: %s", userId.String())
			exceptions.Auth.PermissionDeniedDueToTooManyRequests().Log().SafelyAbortAndResponseWithJSON(ctx, metrics.MetricNames.Server.Responses.Failed.RateLimit)
			return
		}

		setRateLimitHeaders(ctx, remaining, authorizedRateLimiter)

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

func StopAuthorizedRateLimiter() {
	if authorizedRateLimiter != nil {
		authorizedRateLimiter.Stop()
		authorizedRateLimiter = nil
		logs.FInfo(traces.GetTrace(0).FileLineString(), "Authorized rate limiter stopped")
	}
}
