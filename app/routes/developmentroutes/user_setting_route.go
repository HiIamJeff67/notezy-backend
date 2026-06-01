package developmentroutes

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"

	interceptors "github.com/HiIamJeff67/notezy-backend/app/interceptors"
	middlewares "github.com/HiIamJeff67/notezy-backend/app/middlewares"
	modules "github.com/HiIamJeff67/notezy-backend/app/modules"
	metrics "github.com/HiIamJeff67/notezy-backend/app/monitor/metrics"
	constants "github.com/HiIamJeff67/notezy-backend/shared/constants"
)

func configureUserSettingRoutes() {
	userSettingModule := modules.NewUserSettingModule()

	userSettingRoutes := DevelopmentRouterGroup.Group("/userSetting")
	defaultMiddlewares := []gin.HandlerFunc{
		middlewares.UnauthorizedRateLimitMiddleware(),
		middlewares.TimeoutMiddleware(1 * time.Second),
		middlewares.AuthMiddleware(),
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
