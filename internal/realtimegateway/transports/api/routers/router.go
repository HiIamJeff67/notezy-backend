package routers

import (
	"github.com/gin-gonic/gin"

	realtimelease "github.com/HiIamJeff67/notezy-backend/internal/realtimegateway/data/cache/realtimelease"
	cookies "github.com/HiIamJeff67/notezy-backend/shared/cookies"
)

func ConfigureRoutes(
	router *gin.RouterGroup,
	realtimeLeaseCache *realtimelease.RealtimeLeaseCacheClient,
	accessTokenCookieHandler *cookies.CookieHandler,
	refreshTokenCookieHandler *cookies.CookieHandler,
) {
	ConfigureBlockPackRoutes(router, realtimeLeaseCache, accessTokenCookieHandler, refreshTokenCookieHandler)
}
