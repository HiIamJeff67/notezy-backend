package middlewares

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	exceptions "github.com/HiIamJeff67/notezy-backend/internal/exceptions"
	contexts "github.com/HiIamJeff67/notezy-backend/internal/gateway/contexts"
	ratelimit "github.com/HiIamJeff67/notezy-backend/internal/gateway/ratelimit"
	logs "github.com/HiIamJeff67/notezy-backend/internal/platform/observability/logs"
	exceptionwriter "github.com/HiIamJeff67/notezy-backend/shared/exceptionwriter"
	sharedcontexts "github.com/HiIamJeff67/notezy-backend/shared/lib/contexts"
)

var authorizedRateLimiter *ratelimit.HybridRateLimiter // use the hybrid one which including token bucket and cross server request management by redis

func InitAuthorizedRateLimiter(config ratelimit.Config) {
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

func AuthorizedRateLimitMiddleware(config ...ratelimit.Config) gin.HandlerFunc {
	cfg := ratelimit.DefaultAuthorizedConfig()
	if len(config) > 0 {
		cfg = config[0]
	}

	if authorizedRateLimiter == nil {
		InitAuthorizedRateLimiter(cfg)
	}

	return func(ctx *gin.Context) {
		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, sharedcontexts.ContextFieldName_User_Id)
		if exception != nil || userId == nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.New(
				"WrongMiddlewareOrder",
				"Context",
				"Middleware",
				"Cannot find the userId, "+
					"please make sure JWTMiddleware() is placed before AuthorizedRateLimitMiddleware()",
				http.StatusInternalServerError,
				true,
			), ctx)
			return
		}

		allowed, remaining := authorizedRateLimiter.AllowByUserId(*userId)
		if !allowed {
			setRateLimitHeaders(ctx, remaining, authorizedRateLimiter)
			logs.NotezyLogger.Debug(ctx.Request.Context(), fmt.Sprintf("Rate limit exceeded for user: %s", userId.String()))
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.New(
				"PermissionDeniedDueToTooManyRequests",
				"Auth",
				"Authorize",
				"Too many requests; please wait before trying again",
				http.StatusTooManyRequests,
			), ctx, "server.responses.failed.rateLimit")
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
