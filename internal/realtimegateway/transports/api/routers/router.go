package routers

import (
	"github.com/gin-gonic/gin"

	realtimelease "github.com/HiIamJeff67/notezy-backend/internal/realtimegateway/data/cache/realtimelease"
	ratelimit "github.com/HiIamJeff67/notezy-backend/internal/realtimegateway/ratelimit"
	cookies "github.com/HiIamJeff67/notezy-backend/shared/cookies"
)

func ConfigureRoutes(
	router *gin.RouterGroup,
	realtimeLeaseCache *realtimelease.RealtimeLeaseCacheClient,
	accessTokenCookieHandler *cookies.CookieHandler,
	refreshTokenCookieHandler *cookies.CookieHandler,
	authorizedRateLimiter *ratelimit.HybridRateLimiter,
) {
	ConfigureBlockPackRoutes(router, realtimeLeaseCache, accessTokenCookieHandler, refreshTokenCookieHandler, authorizedRateLimiter)
}
