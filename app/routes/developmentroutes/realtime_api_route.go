package developmentroutes

import (
	"time"

	"github.com/gin-gonic/gin"

	interceptors "github.com/HiIamJeff67/notezy-backend/app/interceptors"
	middlewares "github.com/HiIamJeff67/notezy-backend/app/middlewares"
	modules "github.com/HiIamJeff67/notezy-backend/app/modules"
)

func configureDevelopmentRealtimeAPIRoutes(router *gin.RouterGroup) {
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
		realtimeRoutes.POST(
			"/createMyRealtimeConnectionTicket",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("createMyRealtimeConnectionTicket"),
					middlewares.ApplyMeterMiddleware("server.requests.realtime.createMyRealtimeConnectionTicket"),
				},
				defaultMiddlewares,
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
				defaultMiddlewares,
				realtimeModule.Binder.BindCreateMyBlockPackChannelTicket(
					realtimeModule.Controller.CreateMyBlockPackChannelTicket,
				),
			)...,
		)
	}
}
