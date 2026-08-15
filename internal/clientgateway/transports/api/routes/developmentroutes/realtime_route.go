package developmentroutes

import (
	"time"

	"github.com/gin-gonic/gin"

	cookies "github.com/HiIamJeff67/notezy-backend/shared/cookies"

	enumcontract "github.com/HiIamJeff67/notezy-backend/contracts/types/enums"

	binders "github.com/HiIamJeff67/notezy-backend/internal/clientgateway/transports/api/binders"
	controllers "github.com/HiIamJeff67/notezy-backend/internal/clientgateway/transports/api/controllers"
	interceptors "github.com/HiIamJeff67/notezy-backend/internal/clientgateway/transports/api/interceptors"
	middlewares "github.com/HiIamJeff67/notezy-backend/internal/clientgateway/transports/api/middlewares"
	coreadapters "github.com/HiIamJeff67/notezy-backend/internal/clientgateway/transports/core/adapters"
)

type RealtimeRouteDependencies struct {
	CoreClient                *coreadapters.CoreAdapter
	AccessTokenCookieHandler  *cookies.CookieHandler
	RefreshTokenCookieHandler *cookies.CookieHandler
	RateLimiters              RateLimiters
}

func configureDevelopmentRealtimeRoutes(
	router *gin.RouterGroup,
	deps RealtimeRouteDependencies,
) {
	coreClient, accessTokenCookieHandler, refreshTokenCookieHandler, rateLimiters := deps.CoreClient, deps.AccessTokenCookieHandler, deps.RefreshTokenCookieHandler, deps.RateLimiters
	if router == nil {
		router = DevelopmentAPIRouterGroup
	}

	realtimeBinder := binders.NewRealtimeBinder()
	realtimeController := controllers.NewRealtimeController(coreClient)
	realtimeRoutes := router.Group("/realtime")
	connectionRouterGroup := realtimeRoutes.Group("/connection")
	connectionMiddlewares := []gin.HandlerFunc{
		middlewares.UnauthorizedRateLimitMiddleware(rateLimiters.Unauthorized),
		middlewares.TimeoutMiddleware(3 * time.Second),
		middlewares.GatewayAuthenticationMiddleware(accessTokenCookieHandler, refreshTokenCookieHandler),
		interceptors.ShareableResponseWriterInterceptor(
			interceptors.RefreshTokenInterceptor(accessTokenCookieHandler),
			interceptors.EmbeddedInterceptor,
		),
	}
	{
		connectionRouterGroup.POST(
			"/ticket",
			middlewares.Reposition(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("createMyRealtimeConnectionTicket"),
					middlewares.ApplyMeterMiddleware("server.requests.realtime.createMyRealtimeConnectionTicket"),
				},
				append(
					connectionMiddlewares,
					middlewares.AllowedPermissionsAbove(enumcontract.AccessControlPermission_Read),
				),
				realtimeBinder.BindCreateMyRealtimeConnectionTicket(realtimeController.CreateMyRealtimeConnectionTicket),
			)...,
		)
	}

	channelRouterGroup := realtimeRoutes.Group("/channel")
	channelMiddlewares := []gin.HandlerFunc{
		middlewares.UnauthorizedRateLimitMiddleware(rateLimiters.Unauthorized),
		middlewares.TimeoutMiddleware(3 * time.Second),
		middlewares.GatewayAuthenticationMiddleware(accessTokenCookieHandler, refreshTokenCookieHandler),
		interceptors.ShareableResponseWriterInterceptor(
			interceptors.RefreshTokenInterceptor(accessTokenCookieHandler),
			interceptors.EmbeddedInterceptor,
		),
	}
	{
		channelRouterGroup.POST(
			"/block-pack/ticket",
			middlewares.Reposition(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("createMyBlockPackChannelTicket"),
					middlewares.ApplyMeterMiddleware("server.requests.realtime.createMyBlockPackChannelTicket"),
				},
				append(
					channelMiddlewares,
					middlewares.AllowedPermissionsAbove(enumcontract.AccessControlPermission_Read),
				),
				realtimeBinder.BindCreateMyBlockPackChannelTicket(realtimeController.CreateMyBlockPackChannelTicket),
			)...,
		)
	}
}
