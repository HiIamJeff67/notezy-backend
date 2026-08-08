package developmentroutes

import (
	"time"

	"github.com/gin-gonic/gin"

	cookies "github.com/HiIamJeff67/notezy-backend/shared/cookies"

	enumcontract "github.com/HiIamJeff67/notezy-backend/contracts/types/enums"

	binders "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/api/binders"
	controllers "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/api/controllers"
	interceptors "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/api/interceptors"
	middlewares "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/api/middlewares"
	coreadapters "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/core/adapters"
)

func configureDevelopmentRealtimeRoutes(
	router *gin.RouterGroup,
	coreClient *coreadapters.CoreAdapter,
	accessTokenCookieHandler *cookies.CookieHandler,
	refreshTokenCookieHandler *cookies.CookieHandler,
) {
	if router == nil {
		router = DevelopmentAPIRouterGroup
	}

	realtimeBinder := binders.NewRealtimeBinder()
	realtimeController := controllers.NewRealtimeController(coreClient)
	realtimeRoutes := router.Group("/realtime")
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
