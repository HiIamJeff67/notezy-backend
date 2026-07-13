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

func configureDevelopmentBlockRoutes(router *gin.RouterGroup) {
	if router == nil {
		router = DevelopmentAPIRouterGroup
	}

	blockModule := modules.NewBlockModule()
	blockRoutes := router.Group("/block")
	defaultMiddlewares := []gin.HandlerFunc{
		middlewares.UnauthorizedRateLimitMiddleware(),
		middlewares.TimeoutMiddleware(3 * time.Second),
		middlewares.AuthMiddleware(),
		interceptors.ShareableResponseWriterInterceptor(
			interceptors.RefreshTokenInterceptor,
			interceptors.EmbeddedInterceptor,
		),
	}

	blockRoutes.GET(
		"/getMyBlockById",
		middlewares.RepositionMiddleware(
			[]gin.HandlerFunc{
				middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "getMyBlockById"),
				middlewares.ApplyMeterMiddleware(
					otel.Meter(constants.ServiceName),
					metrics.MetricNames.Server.Requests.Block.GetMyBlockById,
				),
			},
			defaultMiddlewares,
			blockModule.Binder.BindGetMyBlockById(
				blockModule.Controller.GetMyBlockById,
			),
		)...,
	)
	blockRoutes.GET(
		"/getMyBlocksByIds",
		middlewares.RepositionMiddleware(
			[]gin.HandlerFunc{
				middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "getMyBlocksByIds"),
				middlewares.ApplyMeterMiddleware(
					otel.Meter(constants.ServiceName),
					metrics.MetricNames.Server.Requests.Block.GetMyBlocksByIds,
				),
			},
			defaultMiddlewares,
			blockModule.Binder.BindGetMyBlocksByIds(
				blockModule.Controller.GetMyBlocksByIds,
			),
		)...,
	)
	blockRoutes.GET(
		"/getMyBlocksByBlockPackId",
		middlewares.RepositionMiddleware(
			[]gin.HandlerFunc{
				middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "getMyBlocksByBlockPackId"),
				middlewares.ApplyMeterMiddleware(
					otel.Meter(constants.ServiceName),
					metrics.MetricNames.Server.Requests.Block.GetMyBlocksByBlockPackId,
				),
			},
			defaultMiddlewares,
			blockModule.Binder.BindGetMyBlocksByBlockPackId(
				blockModule.Controller.GetMyBlocksByBlockPackId,
			),
		)...,
	)
}
