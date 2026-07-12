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

func configureDevelopmentUserAccountRoutes(router *gin.RouterGroup) {
	if router == nil {
		router = DevelopmentAPIRouterGroup
	}

	userAccountModule := modules.NewUserAccountModule()

	userAccountRoutes := router.Group("/userAccount")
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
