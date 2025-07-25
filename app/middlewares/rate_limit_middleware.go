package middlewares

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	contexts "notezy-backend/app/contexts"
	exceptions "notezy-backend/app/exceptions"
	logs "notezy-backend/app/logs"
	constants "notezy-backend/shared/constants"
)

/* ============================== Implementation of Leaky Bucket Algorithm ============================== */

type LeakyBucket struct {
	requestArrivalTimes []time.Time
	capacity            int
	minInterval         time.Duration
	mutex               sync.Mutex
}

func NewLeakyBucket(requestsPerSecond int) *LeakyBucket {
	minInterval := time.Second / time.Duration(requestsPerSecond)
	return &LeakyBucket{
		requestArrivalTimes: make([]time.Time, 0),
		capacity:            requestsPerSecond + constants.ExtraCapacity,
		minInterval:         minInterval,
	}
}

func (lb *LeakyBucket) Allow() bool {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()

	now := time.Now()

	validRequests := make([]time.Time, 0)
	for _, reqArrivalTime := range lb.requestArrivalTimes {
		if now.Sub(reqArrivalTime) < constants.MinIntervalTimeOfLastRequest {
			validRequests = append(validRequests, reqArrivalTime)
		}
	}
	lb.requestArrivalTimes = validRequests

	if len(lb.requestArrivalTimes) >= lb.capacity {
		return false
	}

	if len(lb.requestArrivalTimes) > 0 {
		lastReqArrivalTime := lb.requestArrivalTimes[len(lb.requestArrivalTimes)-1]
		if now.Sub(lastReqArrivalTime) < lb.minInterval {
			return false
		}
	}

	lb.requestArrivalTimes = append(lb.requestArrivalTimes, now)
	return true
}

/* ============================== RateLimitMiddleware ============================== */

var (
	rateLimiters     = make(map[string]*LeakyBucket)
	rateLimiterMutex sync.RWMutex
)

func getRateLimiter(key string, requestsPerSecond int) *LeakyBucket {
	// check if there's already in the map
	rateLimiterMutex.RLock()
	if limiter, exists := rateLimiters[key]; exists {
		rateLimiterMutex.RUnlock()
		return limiter
	}
	rateLimiterMutex.RUnlock()

	rateLimiterMutex.Lock()
	defer rateLimiterMutex.Unlock()

	// check if there's already in the map again
	if limiter, exists := rateLimiters[key]; exists {
		return limiter
	}

	limiter := NewLeakyBucket(requestsPerSecond)
	rateLimiters[key] = limiter
	return limiter
}

/*
 * Use in different types of rate limit control middlewares, once apply this function,
 * the middleware MUST be use after AuthMiddleware() so that the publicId is ensured
 * Note that this generating function is safe enough, since it uses publicId to get the rate limit key
 */
func generateRateLimitKeyByPublicId(ctx *gin.Context) (string, *exceptions.Exception) {
	if publicIdInterface, exists := ctx.Get("publicId"); exists {
		if publicId, ok := publicIdInterface.(string); ok && len(strings.TrimSpace(publicId)) > 0 {
			return fmt.Sprintf("user:%s", publicId), nil
		}
	}
	return "", exceptions.Auth.WrongAccessToken()
}

func generateRateLimitKeyByClientIP(ctx *gin.Context) (string, *exceptions.Exception) {
	clientIP := contexts.GetRealClientIP(ctx)
	if len(strings.TrimSpace(clientIP)) > 0 {
		return fmt.Sprintf("unauthorizedUser:%s", clientIP), nil
	}
	return "", exceptions.Auth.NoClientIPOrReferenceToClient()
}

func logRateLimitExceeded(ctx *gin.Context, key string, limit int) {
	clientIP := contexts.GetRealClientIP(ctx)
	path := ctx.Request.URL.Path
	userAgent := ctx.GetHeader("User-Agent")

	var publicId string
	if publicIdInterface, exists := ctx.Get("publicId"); exists {
		if publicIdStr, ok := publicIdInterface.(string); ok && len(strings.TrimSpace(publicIdStr)) > 0 {
			publicId = fmt.Sprintf("publicId: %s", publicIdStr)
		}
	}
	if publicId == "" {
		publicId = fmt.Sprintf("IP: %s", clientIP)
	}

	logs.FWarn("[RATE_LIMIT] %s, Path: %s, Key: %s, Limit: %d/s, UserAgent: %s\n",
		publicId, path, key, limit, userAgent)
}

func RateLimitMiddleware(requestsPerSecond int) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		key, exception := generateRateLimitKeyByPublicId(ctx)
		if exception != nil {
			exception.Log()
			ctx.AbortWithStatusJSON(exception.HTTPStatusCode, exception.GetGinH())
			return
		}
		limiter := getRateLimiter(key, requestsPerSecond)

		if !limiter.Allow() {
			logRateLimitExceeded(ctx, key, requestsPerSecond)

			ctx.Header("Retry-After", "1")
			ctx.AbortWithStatusJSON(http.StatusTooManyRequests,
				exceptions.Auth.PermissionDeniedDueToTooManyRequests().GetGinH())
			return
		}

		ctx.Next()
	}
}

func UnauthorizedRateLimitMiddleware(requestsPerSecond int) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		key, exception := generateRateLimitKeyByClientIP(ctx)
		if exception != nil {
			exception.Log()
			ctx.AbortWithStatusJSON(exception.HTTPStatusCode, exception.GetGinH())
			return
		}
		limiter := getRateLimiter(key, requestsPerSecond)

		if !limiter.Allow() {
			logRateLimitExceeded(ctx, key, requestsPerSecond)

			ctx.Header("Retry-After", "1")
			exception := exceptions.Auth.PermissionDeniedDueToTooManyRequests()
			ctx.AbortWithStatusJSON(exception.HTTPStatusCode, exception.GetGinH())
			return
		}

		ctx.Next()
	}
}
