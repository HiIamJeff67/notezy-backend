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

func configureDevelopmentUserInfoRoutes() {
	userInfoModule := modules.NewUserInfoModule()

	userInfoRoutes := DevelopmentRouterGroup.Group("/userInfo")
	defaultsMiddlewares := []gin.HandlerFunc{
		middlewares.TimeoutMiddleware(1 * time.Second),
		middlewares.AuthMiddleware(),
		middlewares.AuthorizedRateLimitMiddleware(),
		interceptors.RefreshTokenInterceptor(),
	}
	{
		userInfoRoutes.GET(
			"/getMyInfo",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "getMyInfo"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.UserInfo.GetMyInfo,
					),
				},
				defaultsMiddlewares,
				userInfoModule.Binder.BindGetMyInfo(
					userInfoModule.Controller.GetMyInfo,
				),
			)...,
		)
		userInfoRoutes.PUT(
			"/updateMyInfo",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "updateMyInfo"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.UserInfo.UpdateMyInfo,
					),
				},
				defaultsMiddlewares,
				userInfoModule.Binder.BindUpdateMyInfo(
					userInfoModule.Controller.UpdateMyInfo,
				),
			)...,
		)
	}
}
