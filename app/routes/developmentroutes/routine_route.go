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

func configureDevelopmentRoutineRoutes() {
	routineModule := modules.NewRoutineModule()

	routineRoutes := DevelopmentRouterGroup.Group("/routine")
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
		routineRoutes.GET(
			"/getMyRoutineById",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "getMyRoutineById"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.Routine.GetMyRoutineById,
					),
				},
				defaultMiddlewares,
				routineModule.Binder.BindGetMyRoutineById(
					routineModule.Controller.GetMyRoutineById,
				),
			)...,
		)
		routineRoutes.GET(
			"getMyRoutinesByStationId",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "getMyRoutinesByStationId"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.Routine.GetMyRoutinesByStationId,
					),
				},
				defaultMiddlewares,
				routineModule.Binder.BindGetMyRoutinesByStationId(
					routineModule.Controller.GetMyRoutinesByStationId,
				),
			)...,
		)
		routineRoutes.GET(
			"/getAllMyRoutinesByTimeRange",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "getAllMyRoutinesByTimeRange"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.Routine.GetAllMyRoutinesByTimeRange,
					),
				},
				defaultMiddlewares,
				routineModule.Binder.BindGetAllMyRoutinesByTimeRange(
					routineModule.Controller.GetAllMyRoutinesByTimeRange,
				),
			)...,
		)
		routineRoutes.POST(
			"/createRoutineByStationId",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "createRoutineByStationId"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.Routine.CreateRoutineByStationId,
					),
				},
				defaultMiddlewares,
				routineModule.Binder.BindCreateRoutineByStationId(
					routineModule.Controller.CreateRoutineByStationId,
				),
			)...,
		)
		routineRoutes.POST(
			"/createRoutinesByStationIds",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "createRoutinesByStationIds"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.Routine.CreateRoutinesByStationIds,
					),
				},
				defaultMiddlewares,
				routineModule.Binder.BindCreateRoutinesByStationIds(
					routineModule.Controller.CreateRoutinesByStationIds,
				),
			)...,
		)
		routineRoutes.PUT(
			"/updateMyRoutineById",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "updateMyRoutineById"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.Routine.UpdateMyRoutineById,
					),
				},
				defaultMiddlewares,
				routineModule.Binder.BindUpdateMyRoutineById(
					routineModule.Controller.UpdateMyRoutineById,
				),
			)...,
		)
		routineRoutes.PUT(
			"/updateMyRoutinesByIds",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "updateMyRoutinesByIds"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.Routine.UpdateMyRoutinesByIds,
					),
				},
				defaultMiddlewares,
				routineModule.Binder.BindUpdateMyRoutinesByIds(
					routineModule.Controller.UpdateMyRoutinesByIds,
				),
			)...,
		)
		routineRoutes.POST(
			"/linkRoutineTagById",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "linkRoutineTagById"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.Routine.LinkRoutineTagById,
					),
				},
				defaultMiddlewares,
				routineModule.Binder.BindLinkRoutineTagById(
					routineModule.Controller.LinkRoutineTagById,
				),
			)...,
		)
		routineRoutes.POST(
			"/bulkLinkRoutineTagsByIds",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "bulkLinkRoutineTagsByIds"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.Routine.BulkLinkRoutineTagsByIds,
					),
				},
				defaultMiddlewares,
				routineModule.Binder.BindBulkLinkRoutineTagsByIds(
					routineModule.Controller.BulkLinkRoutineTagsByIds,
				),
			)...,
		)
		routineRoutes.POST(
			"/linkRoutineTaskById",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "linkRoutineTaskById"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.Routine.LinkRoutineTaskById,
					),
				},
				defaultMiddlewares,
				routineModule.Binder.BindLinkRoutineTaskById(
					routineModule.Controller.LinkRoutineTaskById,
				),
			)...,
		)
		routineRoutes.POST(
			"/bulkLinkRoutineTasksByIds",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "bulkLinkRoutineTasksByIds"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.Routine.BulkLinkRoutineTasksByIds,
					),
				},
				defaultMiddlewares,
				routineModule.Binder.BindBulkLinkRoutineTasksByIds(
					routineModule.Controller.BulkLinkRoutineTasksByIds,
				),
			)...,
		)
		routineRoutes.POST(
			"/linkRoutineItemById",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "linkRoutineItemById"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.Routine.LinkRoutineItemById,
					),
				},
				defaultMiddlewares,
				routineModule.Binder.BindLinkRoutineItemById(
					routineModule.Controller.LinkRoutineItemById,
				),
			)...,
		)
		routineRoutes.POST(
			"/bulkLinkRoutineItemsByIds",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "bulkLinkRoutineItemsByIds"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.Routine.BulkLinkRoutineItemsByIds,
					),
				},
				defaultMiddlewares,
				routineModule.Binder.BindBulkLinkRoutineItemsByIds(
					routineModule.Controller.BulkLinkRoutineItemsByIds,
				),
			)...,
		)
		routineRoutes.PATCH(
			"/restoreMyRoutineById",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "restoreMyRoutineById"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.Routine.RestoreMyRoutineById,
					),
				},
				defaultMiddlewares,
				routineModule.Binder.BindRestoreMyRoutineById(
					routineModule.Controller.RestoreMyRoutineById,
				),
			)...,
		)
		routineRoutes.PATCH(
			"/restoreMyRoutinesByIds",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "restoreMyRoutinesByIds"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.Routine.RestoreMyRoutinesByIds,
					),
				},
				defaultMiddlewares,
				routineModule.Binder.BindRestoreMyRoutinesByIds(
					routineModule.Controller.RestoreMyRoutinesByIds,
				),
			)...,
		)
		routineRoutes.DELETE(
			"/deleteMyRoutineById",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "deleteMyRoutineById"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.Routine.DeleteMyRoutineById,
					),
				},
				defaultMiddlewares,
				routineModule.Binder.BindDeleteMyRoutineById(
					routineModule.Controller.DeleteMyRoutineById,
				),
			)...,
		)
		routineRoutes.DELETE(
			"/deleteMyRoutinesByIds",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "deleteMyRoutinesByIds"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.Routine.DeleteMyRoutinesByIds,
					),
				},
				defaultMiddlewares,
				routineModule.Binder.BindDeleteMyRoutinesByIds(
					routineModule.Controller.DeleteMyRoutinesByIds,
				),
			)...,
		)
		routineRoutes.DELETE(
			"/hardDeleteMyRoutineById",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "hardDeleteMyRoutineById"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.Routine.HardDeleteMyRoutineById,
					),
				},
				defaultMiddlewares,
				routineModule.Binder.BindHardDeleteMyRoutineById(
					routineModule.Controller.HardDeleteMyRoutineById,
				),
			)...,
		)
		routineRoutes.DELETE(
			"/hardDeleteMyRoutinesByIds",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "hardDeleteMyRoutinesByIds"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.Routine.HardDeleteMyRoutinesByIds,
					),
				},
				defaultMiddlewares,
				routineModule.Binder.BindHardDeleteMyRoutinesByIds(
					routineModule.Controller.HardDeleteMyRoutinesByIds,
				),
			)...,
		)
	}
}
