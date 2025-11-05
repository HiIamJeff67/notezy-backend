package middlewares

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	rate "golang.org/x/time/rate"

	contexts "notezy-backend/app/contexts"
	exceptions "notezy-backend/app/exceptions"
	lib "notezy-backend/app/lib"
	logs "notezy-backend/app/logs"
	constants "notezy-backend/shared/constants"
	types "notezy-backend/shared/types"
)

type AuthorizedRateLimitConfig struct {
	RateLimit         rate.Limit
	Burst             int
	UserLimit         int32
	WindowDuration    time.Duration
	BackendServerName types.BackendServerName
}

var (
	authorizedRateLimiter            *lib.HybridRateLimiter // use the bybrid one which including token bucket and cross server request management by redis
	DefaultAuthorizedRateLimitConfig = AuthorizedRateLimitConfig{
		RateLimit:         rate.Limit(100),                  // 100 requests/second
		Burst:             20,                               // allowed 20 additional requests/second for burst
		UserLimit:         300,                              // 300 requests/each life time of the bucket (= 300 requests/`WindowDuration`) for each users
		WindowDuration:    time.Minute,                      // 1 minutes to reset the bucket
		BackendServerName: types.BackendServerName_EastAsia, // the current server
	}
)

func InitAuthorizedRateLimiter(config AuthorizedRateLimitConfig) {
	if authorizedRateLimiter != nil {
		authorizedRateLimiter.Stop()
	}

	authorizedRateLimiter = lib.NewHybridRateLimiter(
		config.RateLimit,
		config.Burst,
		config.UserLimit,
		config.WindowDuration,
		config.BackendServerName,
		true,
	)

	logs.FInfo("Authorized rate limiter initialized with rate: %v, burst: %d, user limit: %d, window: %v",
		config.RateLimit, config.Burst, config.UserLimit, config.WindowDuration)
}

func AuthorizedRateLimitMiddleware(config ...AuthorizedRateLimitConfig) gin.HandlerFunc {
	cfg := DefaultAuthorizedRateLimitConfig
	if len(config) > 0 {
		cfg = config[0]
	}

	if authorizedRateLimiter == nil {
		InitAuthorizedRateLimiter(cfg)
	}

	return func(ctx *gin.Context) {
		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
		if exception != nil || userId == nil {
			exceptions.Auth.MissPlacingOrWrongMiddlewareOrder(
				"Cannot find the userId, " +
					"please make sure the AuthMiddleware() is placing before the AuthorizedRateLimitMiddleware()",
			).Log().SafelyResponseWithJSON(ctx)
			return
		}

		fingerprint := getClientFingerprint(ctx)

		allowed, remaining := authorizedRateLimiter.AllowByUserId(*userId, fingerprint)
		if !allowed {
			setRateLimitHeaders(ctx, remaining, authorizedRateLimiter)

			logs.FDebug("Rate limit exceeded for user: %s, fingerprint: %s", userId.String(), fingerprint)
			ctx.JSON(http.StatusTooManyRequests, exceptions.Auth.PermissionDeniedDueToTooManyRequests().GetGinH())
			ctx.Abort()
			return
		}

		setRateLimitHeaders(ctx, remaining, authorizedRateLimiter)

		ctx.Next()
	}
}

func getClientFingerprint(c *gin.Context) string {
	// TODO: use other complex stuff or algorithm or even the machine learning model to generate or get the fingerprint of each clients
	return c.ClientIP()
}

func setRateLimitHeaders(ctx *gin.Context, remaining int32, limiter *lib.HybridRateLimiter) {
	// standard informations
	ctx.Header("X-RateLimit-Limit", strconv.Itoa(int(limiter.UserLimit)))
	ctx.Header("X-RateLimit-Remaining", strconv.Itoa(int(remaining)))
	ctx.Header("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(limiter.WindowDuration).Unix(), 10))

	// extra informations
	ctx.Header("X-RateLimit-Window", limiter.WindowDuration.String())
	ctx.Header("X-RateLimit-Policy", "hybrid-token-bucket")
}

func StopAuthorizedRateLimiter() {
	if authorizedRateLimiter != nil {
		authorizedRateLimiter.Stop()
		authorizedRateLimiter = nil
		logs.FInfo("Authorized rate limiter stopped")
	}
}
