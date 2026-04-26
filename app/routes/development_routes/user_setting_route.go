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

func configureUserSettingRoutes() {
	userSettingModule := modules.NewUserSettingModule()

	userSettingRoutes := DevelopmentRouterGroup.Group("/userSetting")
	defaultMiddlewares := []gin.HandlerFunc{
		middlewares.TimeoutMiddleware(1 * time.Second),
		middlewares.AuthMiddleware(),
		middlewares.AuthorizedRateLimitMiddleware(),
		interceptors.ShareableResponseWriterInterceptor(
			interceptors.RefreshTokenInterceptor,
			interceptors.EmbeddedInterceptor,
		),
	}
	{
		userSettingRoutes.GET(
			"/getMySetting",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "getMySetting"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.UserSetting.GetMySetting,
					),
				},
				defaultMiddlewares,
				userSettingModule.Binder.BindGetMySetting(
					userSettingModule.Controller.GetMySetting,
				),
			)...,
		)
	}
}
