package developmentroutes

import (
	"time"

	"github.com/gin-gonic/gin"

	enumcontract "github.com/HiIamJeff67/notezy-backend/contracts/types/enums"
	binders "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/api/binders"
	controllers "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/api/controllers"
	interceptors "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/api/interceptors"
	middlewares "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/api/middlewares"
	coreadapters "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/core/adapters"
	realtimegatewayadapters "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/realtimegateway/adapters"
	cookies "github.com/HiIamJeff67/notezy-backend/shared/cookies"
)

func configureDevelopmentRealtimeRoutes(
	router *gin.RouterGroup,
	coreClient *coreadapters.CoreClient,
	realtimeGatewayClient *realtimegatewayadapters.RealtimeGatewayClient,
	accessTokenCookieHandler *cookies.CookieHandler,
	refreshTokenCookieHandler *cookies.CookieHandler,
) {
	if router == nil {
		router = DevelopmentAPIRouterGroup
	}

	realtimeBinder := binders.NewRealtimeBinder()
	realtimeController := controllers.NewRealtimeController(coreClient, realtimeGatewayClient)
	realtimeRoutes := router.Group("/realtime")
	defaultMiddlewares := []gin.HandlerFunc{
		middlewares.UnauthorizedRateLimitMiddleware(),
		middlewares.TimeoutMiddleware(3 * time.Second),
		middlewares.JWTMiddleware(accessTokenCookieHandler, refreshTokenCookieHandler),
		interceptors.ShareableResponseWriterInterceptor(
			interceptors.RefreshTokenInterceptor(accessTokenCookieHandler),
			interceptors.EmbeddedInterceptor,
		),
	}
	{
		realtimeRoutes.GET(
			"/block-pack/:blockPackId/participants",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("getMyBlockPackRealtimeParticipants"),
					middlewares.ApplyMeterMiddleware("server.requests.realtime.getMyBlockPackRealtimeParticipants"),
				},
				append(
					defaultMiddlewares,
					middlewares.AllowedPermissionsAbove(enumcontract.AccessControlPermission_Read),
				),
				realtimeBinder.BindGetMyBlockPackRealtimeParticipants(realtimeController.GetMyBlockPackRealtimeParticipants),
			)...,
		)
	}

	connectionRouterGroup := realtimeRoutes.Group("/connection")
	connectionMiddlewares := []gin.HandlerFunc{
		middlewares.UnauthorizedRateLimitMiddleware(),
		middlewares.TimeoutMiddleware(3 * time.Second),
		middlewares.JWTMiddleware(accessTokenCookieHandler, refreshTokenCookieHandler),
		interceptors.ShareableResponseWriterInterceptor(
			interceptors.RefreshTokenInterceptor(accessTokenCookieHandler),
			interceptors.EmbeddedInterceptor,
		),
	}
	{
		connectionRouterGroup.POST(
			"/ticket",
			middlewares.RepositionMiddleware(
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
		middlewares.UnauthorizedRateLimitMiddleware(),
		middlewares.TimeoutMiddleware(3 * time.Second),
		middlewares.JWTMiddleware(accessTokenCookieHandler, refreshTokenCookieHandler),
		interceptors.ShareableResponseWriterInterceptor(
			interceptors.RefreshTokenInterceptor(accessTokenCookieHandler),
			interceptors.EmbeddedInterceptor,
		),
	}
	{
		channelRouterGroup.POST(
			"/block-pack/ticket",
			middlewares.RepositionMiddleware(
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
