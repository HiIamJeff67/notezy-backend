package routers

import (
	"time"

	"github.com/gin-gonic/gin"

	cookies "github.com/HiIamJeff67/notezy-backend/shared/cookies"

	realtimelease "github.com/HiIamJeff67/notezy-backend/internal/realtimegateway/data/cache/realtimelease"
	endpoints "github.com/HiIamJeff67/notezy-backend/internal/realtimegateway/transports/api/endpoints"
	middlewares "github.com/HiIamJeff67/notezy-backend/internal/realtimegateway/transports/api/middlewares"
)

func ConfigureBlockPackRoutes(
	router *gin.RouterGroup,
	realtimeLeaseCache *realtimelease.RealtimeLeaseCacheClient,
	accessTokenCookieHandler *cookies.CookieHandler,
	refreshTokenCookieHandler *cookies.CookieHandler,
) {
	endpoint := endpoints.NewBlockPackEndpoint(realtimeLeaseCache)

	router.GET(
		"/block-pack/:blockPackId/participants",
		middlewares.JWTMiddleware(accessTokenCookieHandler, refreshTokenCookieHandler),
		middlewares.AuthorizedRateLimitMiddleware(),
		middlewares.TimeoutMiddleware(3*time.Second),
		middlewares.ApplyTracerMiddleware("getRealtimeBlockPackParticipants"),
		middlewares.ApplyMeterMiddleware("server.requests.realtime.blockPackParticipants"),
		endpoint.GetParticipants,
	)
}
