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

func configureDevelopmentRoutineTaskRecordRoutes(router *gin.RouterGroup) {
	if router == nil {
		router = DevelopmentAPIRouterGroup
	}

	routineTaskRecordModule := modules.NewRoutineTaskRecordModule()

	routineTaskRecordRoutes := router.Group("/routineTaskRecord")
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
		routineTaskRecordRoutes.GET(
			"/getAllMyRoutineTaskRecordsByRoutineTaskId",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "getAllMyRoutineTaskRecordsByRoutineTaskId"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.RoutineTaskRecord.GetAllMyRoutineTaskRecordsByRoutineTaskId,
					),
				},
				defaultMiddlewares,
				routineTaskRecordModule.Binder.BindGetAllMyRoutineTaskRecordsByRoutineTaskId(
					routineTaskRecordModule.Controller.GetAllMyRoutineTaskRecordsByRoutineTaskId,
				),
			)...,
		)
		routineTaskRecordRoutes.GET(
			"/visualizeMyRoutineTaskRecordStatusCount",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "visualizeMyRoutineTaskRecordStatusCount"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.RoutineTaskRecord.VisualizeMyRoutineTaskRecordStatusCount,
					),
				},
				defaultMiddlewares,
				routineTaskRecordModule.Binder.BindVisualizeMyRoutineTaskRecordStatusCount(
					routineTaskRecordModule.Controller.VisualizeMyRoutineTaskRecordStatusCount,
				),
			)...,
		)
		routineTaskRecordRoutes.GET(
			"/visualizeMyRoutineTaskRecordPurposeCount",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "visualizeMyRoutineTaskRecordPurposeCount"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.RoutineTaskRecord.VisualizeMyRoutineTaskRecordPurposeCount,
					),
				},
				defaultMiddlewares,
				routineTaskRecordModule.Binder.BindVisualizeMyRoutineTaskRecordPurposeCount(
					routineTaskRecordModule.Controller.VisualizeMyRoutineTaskRecordPurposeCount,
				),
			)...,
		)
		routineTaskRecordRoutes.GET(
			"/visualizeMyRoutineTaskRecordScheduledAtCount",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "visualizeMyRoutineTaskRecordScheduledAtCount"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.RoutineTaskRecord.VisualizeMyRoutineTaskRecordScheduledAtCount,
					),
				},
				defaultMiddlewares,
				routineTaskRecordModule.Binder.BindVisualizeMyRoutineTaskRecordScheduledAtCount(
					routineTaskRecordModule.Controller.VisualizeMyRoutineTaskRecordScheduledAtCount,
				),
			)...,
		)
		routineTaskRecordRoutes.GET(
			"/visualizeMyRoutineTaskRecordActualStartedAtCount",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "visualizeMyRoutineTaskRecordActualStartedAtCount"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.RoutineTaskRecord.VisualizeMyRoutineTaskRecordActualStartedAtCount,
					),
				},
				defaultMiddlewares,
				routineTaskRecordModule.Binder.BindVisualizeMyRoutineTaskRecordActualStartedAtCount(
					routineTaskRecordModule.Controller.VisualizeMyRoutineTaskRecordActualStartedAtCount,
				),
			)...,
		)
		routineTaskRecordRoutes.GET(
			"/visualizeMyRoutineTaskRecordActualEndedAtCount",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "visualizeMyRoutineTaskRecordActualEndedAtCount"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.RoutineTaskRecord.VisualizeMyRoutineTaskRecordActualEndedAtCount,
					),
				},
				defaultMiddlewares,
				routineTaskRecordModule.Binder.BindVisualizeMyRoutineTaskRecordActualEndedAtCount(
					routineTaskRecordModule.Controller.VisualizeMyRoutineTaskRecordActualEndedAtCount,
				),
			)...,
		)
	}
}
