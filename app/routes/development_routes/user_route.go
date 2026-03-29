package developmentroutes

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"

	interceptors "notezy-backend/app/interceptors"
	middlewares "notezy-backend/app/middlewares"
	modules "notezy-backend/app/modules"
	metrics "notezy-backend/app/monitor/metrics"
	constants "notezy-backend/shared/constants"
)

func configureDevelopmentUserRoutes() {
	userModule := modules.NewUserModule()

	userRoutes := DevelopmentRouterGroup.Group("/user")
	defaultMiddlewares := []gin.HandlerFunc{
		middlewares.TimeoutMiddleware(1 * time.Second),
		middlewares.AuthMiddleware(),
		middlewares.AuthorizedRateLimitMiddleware(),
		interceptors.RefreshTokenInterceptor(),
	}
	{
		userRoutes.GET(
			"/getUserData",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "getUserData"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.User.GetUserData,
					),
				},
				defaultMiddlewares,
				userModule.Binder.BindGetUserData(
					userModule.Controller.GetUserData,
				),
			)...,
		)
		userRoutes.GET(
			"/getMe",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "getMe"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.User.GetMe,
					),
				},
				defaultMiddlewares,
				userModule.Binder.BindGetMe(
					userModule.Controller.GetMe,
				),
			)...,
		)
		userRoutes.PUT(
			"/updateMe",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "updateMe"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.User.UpdateMe,
					),
				},
				defaultMiddlewares,
				userModule.Binder.BindUpdateMe(
					userModule.Controller.UpdateMe,
				),
			)...,
		)
	}
}
