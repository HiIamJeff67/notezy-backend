package developmentroutes

import (
	"time"

	"github.com/gin-gonic/gin"

	binders "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/api/binders"
	controllers "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/api/controllers"
	interceptors "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/api/interceptors"
	middlewares "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/api/middlewares"
	coreadapters "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/core/adapters"
	sharedtypes "github.com/HiIamJeff67/notezy-backend/internal/shared/types"
)

func configureDevelopmentRealtimeRoutes(router *gin.RouterGroup, coreClient *coreadapters.CoreClient) {
	if router == nil {
		router = DevelopmentAPIRouterGroup
	}

	realtimeBinder := binders.NewRealtimeBinder()
	realtimeController := controllers.NewRealtimeController(coreClient)
	realtimeRoutes := router.Group("/realtime")
	defaultMiddlewares := []gin.HandlerFunc{
		middlewares.UnauthorizedRateLimitMiddleware(),
		middlewares.TimeoutMiddleware(3 * time.Second),
		middlewares.AuthMiddleware(),
		interceptors.ShareableResponseWriterInterceptor(
			interceptors.RefreshTokenInterceptor,
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
					middlewares.AllowedPermissionsAbove(sharedtypes.AccessControlPermission_Read),
				),
				realtimeBinder.BindGetMyBlockPackRealtimeParticipants(realtimeController.GetMyBlockPackRealtimeParticipants),
			)...,
		)
	}

	connectionRouterGroup := realtimeRoutes.Group("/connection")
	connectionMiddlewares := []gin.HandlerFunc{
		middlewares.UnauthorizedRateLimitMiddleware(),
		middlewares.TimeoutMiddleware(3 * time.Second),
		middlewares.AuthMiddleware(),
		interceptors.ShareableResponseWriterInterceptor(
			interceptors.RefreshTokenInterceptor,
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
					middlewares.AllowedPermissionsAbove(sharedtypes.AccessControlPermission_Read),
				),
				realtimeBinder.BindCreateMyRealtimeConnectionTicket(realtimeController.CreateMyRealtimeConnectionTicket),
			)...,
		)
	}

	channelRouterGroup := realtimeRoutes.Group("/channel")
	channelMiddlewares := []gin.HandlerFunc{
		middlewares.UnauthorizedRateLimitMiddleware(),
		middlewares.TimeoutMiddleware(3 * time.Second),
		middlewares.AuthMiddleware(),
		interceptors.ShareableResponseWriterInterceptor(
			interceptors.RefreshTokenInterceptor,
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
					middlewares.AllowedPermissionsAbove(sharedtypes.AccessControlPermission_Read),
				),
				realtimeBinder.BindCreateMyBlockPackChannelTicket(realtimeController.CreateMyBlockPackChannelTicket),
			)...,
		)
	}
}
