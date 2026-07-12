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

func configureDevelopmentRoutineTaskRoutes(router *gin.RouterGroup) {
	if router == nil {
		router = DevelopmentAPIRouterGroup
	}

	routineTaskModule := modules.NewRoutineTaskModule()

	routineTaskRoutes := router.Group("/routineTask")
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
		routineTaskRoutes.GET(
			"/getMyRoutineTaskById",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "getMyRoutineTaskById"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.RoutineTask.GetMyRoutineTaskById,
					),
				},
				defaultMiddlewares,
				routineTaskModule.Binder.BindGetMyRoutineTaskById(
					routineTaskModule.Controller.GetMyRoutineTaskById,
				),
			)...,
		)
		routineTaskRoutes.GET(
			"/getAllMyRoutineTasksByRoutineIds",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "getAllMyRoutineTasksByRoutineIds"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.RoutineTask.GetAllMyRoutineTasksByRoutineIds,
					),
				},
				defaultMiddlewares,
				routineTaskModule.Binder.BindGetAllMyRoutineTasksByRoutineIds(
					routineTaskModule.Controller.GetAllMyRoutineTasksByRoutineIds,
				),
			)...,
		)
		routineTaskRoutes.GET(
			"/getAllMyRoutineTasks",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "getAllMyRoutineTasks"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.RoutineTask.GetAllMyRoutineTasks,
					),
				},
				defaultMiddlewares,
				routineTaskModule.Binder.BindGetAllMyRoutineTasks(
					routineTaskModule.Controller.GetAllMyRoutineTasks,
				),
			)...,
		)
		routineTaskRoutes.GET(
			"/visualizeMyRoutineTaskStatusCount",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "visualizeMyRoutineTaskStatusCount"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.RoutineTask.VisualizeMyRoutineTaskStatusCount,
					),
				},
				defaultMiddlewares,
				routineTaskModule.Binder.BindVisualizeMyRoutineTaskStatusCount(
					routineTaskModule.Controller.VisualizeMyRoutineTaskStatusCount,
				),
			)...,
		)
		routineTaskRoutes.GET(
			"/visualizeMyRoutineTaskPurposeCount",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "visualizeMyRoutineTaskPurposeCount"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.RoutineTask.VisualizeMyRoutineTaskPurposeCount,
					),
				},
				defaultMiddlewares,
				routineTaskModule.Binder.BindVisualizeMyRoutineTaskPurposeCount(
					routineTaskModule.Controller.VisualizeMyRoutineTaskPurposeCount,
				),
			)...,
		)
		routineTaskRoutes.GET(
			"/visualizeMyRoutineTaskScheduledAtCount",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "visualizeMyRoutineTaskScheduledAtCount"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.RoutineTask.VisualizeMyRoutineTaskScheduledAtCount,
					),
				},
				defaultMiddlewares,
				routineTaskModule.Binder.BindVisualizeMyRoutineTaskScheduledAtCount(
					routineTaskModule.Controller.VisualizeMyRoutineTaskScheduledAtCount,
				),
			)...,
		)
		routineTaskRoutes.GET(
			"/visualizeMyRoutineTaskActualStartedAtCount",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "visualizeMyRoutineTaskActualStartedAtCount"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.RoutineTask.VisualizeMyRoutineTaskActualStartedAtCount,
					),
				},
				defaultMiddlewares,
				routineTaskModule.Binder.BindVisualizeMyRoutineTaskActualStartedAtCount(
					routineTaskModule.Controller.VisualizeMyRoutineTaskActualStartedAtCount,
				),
			)...,
		)
		routineTaskRoutes.GET(
			"/visualizeMyRoutineTaskActualEndedAtCount",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "visualizeMyRoutineTaskActualEndedAtCount"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.RoutineTask.VisualizeMyRoutineTaskActualEndedAtCount,
					),
				},
				defaultMiddlewares,
				routineTaskModule.Binder.BindVisualizeMyRoutineTaskActualEndedAtCount(
					routineTaskModule.Controller.VisualizeMyRoutineTaskActualEndedAtCount,
				),
			)...,
		)
		routineTaskRoutes.POST(
			"/createRoutineTaskByRoutineId",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "createRoutineTaskByRoutineId"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.RoutineTask.CreateRoutineTaskByRoutineId,
					),
				},
				defaultMiddlewares,
				routineTaskModule.Binder.BindCreateRoutineTaskByRoutineId(
					routineTaskModule.Controller.CreateRoutineTaskByRoutineId,
				),
			)...,
		)
		routineTaskRoutes.PUT(
			"/updateMyRoutineTaskById",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "updateMyRoutineTaskById"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.RoutineTask.UpdateMyRoutineTaskById,
					),
				},
				defaultMiddlewares,
				routineTaskModule.Binder.BindUpdateMyRoutineTaskById(
					routineTaskModule.Controller.UpdateMyRoutineTaskById,
				),
			)...,
		)
		routineTaskRoutes.PUT(
			"/pauseMyRoutineTaskById",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "pauseMyRoutineTaskById"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.RoutineTask.PauseMyRoutineTaskById,
					),
				},
				defaultMiddlewares,
				routineTaskModule.Binder.BindPauseMyRoutineTaskById(
					routineTaskModule.Controller.PauseMyRoutineTaskById,
				),
			)...,
		)
		routineTaskRoutes.PUT(
			"/resumeMyRoutineTaskById",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "resumeMyRoutineTaskById"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.RoutineTask.ResumeMyRoutineTaskById,
					),
				},
				defaultMiddlewares,
				routineTaskModule.Binder.BindResumeMyRoutineTaskById(
					routineTaskModule.Controller.ResumeMyRoutineTaskById,
				),
			)...,
		)
		routineTaskRoutes.DELETE(
			"/hardDeleteMyRoutineTaskById",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "hardDeleteMyRoutineTaskById"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.RoutineTask.HardDeleteMyRoutineTaskById,
					),
				},
				defaultMiddlewares,
				routineTaskModule.Binder.BindHardDeleteMyRoutineTaskById(
					routineTaskModule.Controller.HardDeleteMyRoutineTaskById,
				),
			)...,
		)
		routineTaskRoutes.DELETE(
			"/hardDeleteMyRoutineTasksByIds",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "hardDeleteMyRoutineTasksByIds"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.RoutineTask.HardDeleteMyRoutineTasksByIds,
					),
				},
				defaultMiddlewares,
				routineTaskModule.Binder.BindHardDeleteMyRoutineTasksByIds(
					routineTaskModule.Controller.HardDeleteMyRoutineTasksByIds,
				),
			)...,
		)
	}
}
