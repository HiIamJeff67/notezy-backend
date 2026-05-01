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

func configureDevelopmentUserAccountRoutes() {
	userAccountModule := modules.NewUserAccountModule()

	userAccountRoutes := DevelopmentRouterGroup.Group("/userAccount")
	defaultMiddlewares := []gin.HandlerFunc{
		middlewares.TimeoutMiddleware(3 * time.Second),
		middlewares.AuthMiddleware(),
		middlewares.AuthorizedRateLimitMiddleware(),
		interceptors.ShareableResponseWriterInterceptor(
			interceptors.RefreshTokenInterceptor,
			interceptors.EmbeddedInterceptor,
		),
	}
	{
		userAccountRoutes.GET(
			"/getMyAccount",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "getMyAccount"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.UserAccount.GetMyAccount,
					),
				},
				defaultMiddlewares,
				userAccountModule.Binder.BindGetMyAccount(
					userAccountModule.Controller.GetMyAccount,
				),
			)...,
		)
		userAccountRoutes.PUT(
			"/updateMyAccount",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "updateMyAccount"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.UserAccount.UpdateMyAccount,
					),
					middlewares.CSRFMiddleware(),
				},
				defaultMiddlewares,
				userAccountModule.Binder.BindUpdateMyAccount(
					userAccountModule.Controller.UpdateMyAccount,
				),
			)...,
		)
		userAccountRoutes.PUT(
			"/bindGoogleAccount",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "bindGoogleAccount"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.UserAccount.BindGoogleAccount,
					),
				},
				defaultMiddlewares,
				userAccountModule.Binder.BindBindGoogleAccount(
					userAccountModule.Controller.BindGoogleAccount,
				),
			)...,
		)
		userAccountRoutes.PUT(
			"/unbindGoogleAccount",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "unbindGoogleAccount"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.UserAccount.UnbindGoogleAccount,
					),
				},
				defaultMiddlewares,
				userAccountModule.Binder.BindUnbindGoogleAccount(
					userAccountModule.Controller.UnbindGoogleAccount,
				),
			)...,
		)
	}
}
