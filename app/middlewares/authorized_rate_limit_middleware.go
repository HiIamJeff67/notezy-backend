package middlewares

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"

	configs "github.com/HiIamJeff67/notezy-backend/app/configs"
	contexts "github.com/HiIamJeff67/notezy-backend/app/contexts"
	exceptions "github.com/HiIamJeff67/notezy-backend/app/exceptions"
	logs "github.com/HiIamJeff67/notezy-backend/app/monitor/logs"
	ratelimit "github.com/HiIamJeff67/notezy-backend/shared/lib/ratelimit"
	types "github.com/HiIamJeff67/notezy-backend/shared/types"
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

	logs.NotezyLogger.Info(context.Background(), fmt.Sprintf("Authorized rate limiter initialized with rate: %v, burst: %d, user limit: %d, window: %v", config.RateLimit, config.Burst, config.UserLimit, config.WindowDuration))
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
			logs.NotezyLogger.Debug(ctx.Request.Context(), fmt.Sprintf("Rate limit exceeded for user: %s", userId.String()))
			exceptions.Auth.PermissionDeniedDueToTooManyRequests().Log().SafelyAbortAndResponseWithJSON(ctx, "server.responses.failed.rateLimit")
			return
		}

		setRateLimitHeaders(ctx, remaining, authorizedRateLimiter)

		ctx.Next()
	}
}

func StopAuthorizedRateLimiter() {
	if authorizedRateLimiter != nil {
		authorizedRateLimiter.Stop()
		authorizedRateLimiter = nil
		logs.NotezyLogger.Info(context.Background(), "Authorized rate limiter stopped")
	}
}
