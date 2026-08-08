package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	exceptions "github.com/HiIamJeff67/notezy-backend/contracts/types/exceptions"
	sharedcontexts "github.com/HiIamJeff67/notezy-backend/shared/lib/contexts"
	logs "github.com/HiIamJeff67/notezy-backend/shared/platform/observability/logs"
	exceptionwriter "github.com/HiIamJeff67/notezy-backend/shared/util/exceptionwriter"

	ratelimit "github.com/HiIamJeff67/notezy-backend/internal/realtimegateway/ratelimit"
)

var authorizedRateLimiter *ratelimit.HybridRateLimiter

func InitAuthorizedRateLimiter(config ratelimit.Config) {
	if authorizedRateLimiter != nil {
		authorizedRateLimiter.Stop()
	}

	authorizedRateLimiter = ratelimit.NewHybridRateLimiter(config, true)
	logs.NotezyLogger.Info(context.Background(), fmt.Sprintf("Authorized realtime rate limiter initialized with rate: %v, burst: %d, user limit: %d, window: %v", config.RateLimit, config.Burst, config.UserLimit, config.WindowDuration))
}

func AuthorizedRateLimitMiddleware() gin.HandlerFunc {
	if authorizedRateLimiter == nil {
		InitAuthorizedRateLimiter(ratelimit.DefaultUpgradeConfig())
	}

	return func(ctx *gin.Context) {
		userPublicId, exists := ctx.Get(sharedcontexts.ContextFieldName_User_PublicId.String())
		if !exists || userPublicId == nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.New(
				"WrongMiddlewareOrder",
				"Context",
				"Middleware",
				"Cannot find the user public ID; JWTMiddleware() must run first",
				http.StatusInternalServerError,
				true,
			), ctx)
			return
		}
		publicId, err := uuid.Parse(fmt.Sprint(userPublicId))
		if err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("User").WithOrigin(err), ctx)
			return
		}

		allowed, remaining := authorizedRateLimiter.AllowByUserId(publicId)
		ctx.Header("X-RateLimit-Limit", strconv.Itoa(int(authorizedRateLimiter.UserLimit)))
		ctx.Header("X-RateLimit-Remaining", strconv.Itoa(int(remaining)))
		ctx.Header("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(authorizedRateLimiter.WindowDuration).Unix(), 10))
		ctx.Header("X-RateLimit-Window", authorizedRateLimiter.WindowDuration.String())
		ctx.Header("X-RateLimit-Policy", "hybrid-token-bucket")
		if !allowed {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.New(
				"PermissionDeniedDueToTooManyRequests",
				"RealtimeGateway",
				"Authorize",
				"Too many requests; please wait before trying again",
				http.StatusTooManyRequests,
			), ctx, "server.responses.failed.rateLimit")
			return
		}

		ctx.Next()
	}
}

func StopAuthorizedRateLimiter() {
	if authorizedRateLimiter != nil {
		authorizedRateLimiter.Stop()
		authorizedRateLimiter = nil
		logs.NotezyLogger.Info(context.Background(), "Authorized realtime rate limiter stopped")
	}
}
