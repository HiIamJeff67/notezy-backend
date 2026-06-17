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

func configureDevelopmentRoutineTaskRoutes() {
	routineTaskModule := modules.NewRoutineTaskModule()

	routineTaskRoutes := DevelopmentRouterGroup.Group("/routineTask")
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
			"/getAllMyRoutineTasksByStationIds",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "getAllMyRoutineTasksByStationIds"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.RoutineTask.GetAllMyRoutineTasksByStationIds,
					),
				},
				defaultMiddlewares,
				routineTaskModule.Binder.BindGetAllMyRoutineTasksByStationIds(
					routineTaskModule.Controller.GetAllMyRoutineTasksByStationIds,
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
		routineTaskRoutes.POST(
			"/createRoutineTaskByStationId",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "createRoutineTaskByStationId"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.RoutineTask.CreateRoutineTaskByStationId,
					),
				},
				defaultMiddlewares,
				routineTaskModule.Binder.BindCreateRoutineTaskByStationId(
					routineTaskModule.Controller.CreateRoutineTaskByStationId,
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
