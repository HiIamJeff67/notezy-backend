package developmentroutes

import (
	"time"

	"github.com/gin-gonic/gin"

	interceptors "github.com/HiIamJeff67/notezy-backend/app/interceptors"
	middlewares "github.com/HiIamJeff67/notezy-backend/app/middlewares"
	enums "github.com/HiIamJeff67/notezy-backend/app/models/schemas/enums"
	modules "github.com/HiIamJeff67/notezy-backend/app/modules"
	realtime "github.com/HiIamJeff67/notezy-backend/app/realtime"
	constants "github.com/HiIamJeff67/notezy-backend/shared/constants"
)

func ConfigureRealtimeRoutes() {
	gateway := realtime.NewGateway()

	DevelopmentRouter.GET(
		"/"+constants.RealtimeDevelopmentBaseURL,
		middlewares.DomainWhiteListMiddleware(),
		middlewares.RealtimeUpgradeRateLimitMiddleware(),
		gateway.Handle,
	)
}

func configureDevelopmentRealtimeRoutes(router *gin.RouterGroup) {
	if router == nil {
		router = DevelopmentAPIRouterGroup
	}

	realtimeModule := modules.NewRealtimeModule()
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
			"/blockPacks/:blockPackId/participants",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("getMyBlockPackRealtimeParticipants"),
					middlewares.ApplyMeterMiddleware("server.requests.realtime.getMyBlockPackRealtimeParticipants"),
				},
				append(
					defaultMiddlewares,
					middlewares.AllowedPermissionsAbove(enums.AccessControlPermission_Read),
				),
				realtimeModule.Binder.BindGetMyBlockPackRealtimeParticipants(
					realtimeModule.Controller.GetMyBlockPackRealtimeParticipants,
				),
			)...,
		)
		realtimeRoutes.POST(
			"/createMyRealtimeConnectionTicket",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("createMyRealtimeConnectionTicket"),
					middlewares.ApplyMeterMiddleware("server.requests.realtime.createMyRealtimeConnectionTicket"),
				},
				append(
					defaultMiddlewares,
					middlewares.AllowedPermissionsAbove(enums.AccessControlPermission_Read),
				),
				realtimeModule.Binder.BindCreateMyRealtimeConnectionTicket(
					realtimeModule.Controller.CreateMyRealtimeConnectionTicket,
				),
			)...,
		)
		realtimeRoutes.POST(
			"/createMyBlockPackChannelTicket",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("createMyBlockPackChannelTicket"),
					middlewares.ApplyMeterMiddleware("server.requests.realtime.createMyBlockPackChannelTicket"),
				},
				append(
					defaultMiddlewares,
					middlewares.AllowedPermissionsAbove(enums.AccessControlPermission_Read),
				),
				realtimeModule.Binder.BindCreateMyBlockPackChannelTicket(
					realtimeModule.Controller.CreateMyBlockPackChannelTicket,
				),
			)...,
		)
	}
}
