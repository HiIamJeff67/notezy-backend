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

func configureDevelopmentSubShelfRoutes() {
	subShelfModule := modules.NewSubShelfModule()

	subShelfRoutes := DevelopmentRouterGroup.Group("/subShelf")
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
		subShelfRoutes.GET(
			"/getMySubShelfById",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "getMySubShelfById"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.SubShelf.GetMySubShelfById,
					),
				},
				defaultMiddlewares,
				subShelfModule.Binder.BindGetMySubShelfById(
					subShelfModule.Controller.GetMySubShelfById,
				),
			)...,
		)
		subShelfRoutes.GET(
			"/getMySubShelvesByPrevSubShelfId",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "getMySubShelvesByPrevSubShelfId"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.SubShelf.GetMySubShelvesByPrevSubShelfId,
					),
				},
				defaultMiddlewares,
				subShelfModule.Binder.BindGetMySubShelvesByPrevSubShelfId(
					subShelfModule.Controller.GetMySubShelvesByPrevSubShelfId,
				),
			)...,
		)
		subShelfRoutes.GET(
			"/getAllMySubShelvesByRootShelfId",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "getAllMySubShelvesByRootShelfId"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.SubShelf.GetAllMySubShelvesByRootShelfId,
					),
				},
				defaultMiddlewares,
				subShelfModule.Binder.BindGetAllMySubShelvesByRootShelfId(
					subShelfModule.Controller.GetAllMySubShelvesByRootShelfId,
				),
			)...,
		)
		subShelfRoutes.GET(
			"/getMySubShelvesAndItemsByPrevSubShelfId",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "getMySubShelvesAndItemsByPrevSubShelfId"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.SubShelf.GetMySubShelvesAndItemsByPrevSubShelfId,
					),
				},
				defaultMiddlewares,
				subShelfModule.Binder.BindGetMySubShelvesAndItemsByPrevSubShelfId(
					subShelfModule.Controller.GetMySubShelvesAndItemsByPrevSubShelfId,
				),
			)...,
		)
		subShelfRoutes.POST(
			"/createSubShelfByRootShelfId",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "createSubShelfByRootShelfId"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.SubShelf.CreateSubShelfByRootShelfId,
					),
				},
				defaultMiddlewares,
				subShelfModule.Binder.BindCreateSubShelfByRootShelfId(
					subShelfModule.Controller.CreateSubShelfByRootShelfId,
				),
			)...,
		)
		subShelfRoutes.POST(
			"/createSubShelvesByRootShelfIds",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "createSubShelvesByRootShelfIds"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.SubShelf.CreateSubShelvesByRootShelfIds,
					),
				},
				defaultMiddlewares,
				subShelfModule.Binder.BindCreateSubShelvesByRootShelfIds(
					subShelfModule.Controller.CreateSubShelvesByRootShelfIds,
				),
			)...,
		)
		subShelfRoutes.PUT(
			"/updateMySubShelfById",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "updateMySubShelfById"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.SubShelf.UpdateMySubShelfById,
					),
				},
				defaultMiddlewares,
				subShelfModule.Binder.BindUpdateMySubShelfById(
					subShelfModule.Controller.UpdateMySubShelfById,
				),
			)...,
		)
		subShelfRoutes.PUT(
			"/updateMySubShelvesByIds",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "updateMySubShelvesByIds"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.SubShelf.UpdateMySubShelvesByIds,
					),
				},
				defaultMiddlewares,
				subShelfModule.Binder.BindUpdateMySubShelvesByIds(
					subShelfModule.Controller.UpdateMySubShelvesByIds,
				),
			)...,
		)
		subShelfRoutes.PUT(
			"/moveMySubShelf",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "moveMySubShelf"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.SubShelf.MoveMySubShelf,
					),
				},
				defaultMiddlewares,
				subShelfModule.Binder.BindMoveMySubShelf(
					subShelfModule.Controller.MoveMySubShelf,
				),
			)...,
		)
		subShelfRoutes.PUT(
			"/moveMySubShelves",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "moveMySubShelves"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.SubShelf.MoveMySubShelves,
					),
				},
				defaultMiddlewares,
				subShelfModule.Binder.BindMoveMySubShelves(
					subShelfModule.Controller.MoveMySubShelves,
				),
			)...,
		)
		subShelfRoutes.PUT(
			"/batchMoveMySubShelves",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "batchMoveMySubShelves"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.SubShelf.BatchMoveMySubShelves,
					),
				},
				defaultMiddlewares,
				subShelfModule.Binder.BindBatchMoveMySubShelves(
					subShelfModule.Controller.BatchMoveMySubShelves,
				),
			)...,
		)
		subShelfRoutes.PATCH(
			"/restoreMySubShelfById",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "restoreMySubShelfById"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.SubShelf.RestoreMySubShelfById,
					),
				},
				defaultMiddlewares,
				subShelfModule.Binder.BindRestoreMySubShelfById(
					subShelfModule.Controller.RestoreMySubShelfById,
				),
			)...,
		)
		subShelfRoutes.PATCH(
			"/restoreMySubShelvesByIds",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "restoreMySubShelvesByIds"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.SubShelf.RestoreMySubShelvesByIds,
					),
				},
				defaultMiddlewares,
				subShelfModule.Binder.BindRestoreMySubShelvesByIds(
					subShelfModule.Controller.RestoreMySubShelvesByIds,
				),
			)...,
		)
		subShelfRoutes.DELETE(
			"/deleteMySubShelfById",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "deleteMySubShelfById"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.SubShelf.DeleteMySubShelfById,
					),
				},
				defaultMiddlewares,
				subShelfModule.Binder.BindDeleteMySubShelfById(
					subShelfModule.Controller.DeleteMySubShelfById,
				),
			)...,
		)
		subShelfRoutes.DELETE(
			"/deleteMySubShelvesByIds",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "deleteMySubShelvesByIds"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.SubShelf.DeleteMySubShelvesByIds,
					),
				},
				defaultMiddlewares,
				subShelfModule.Binder.BindDeleteMySubShelvesByIds(
					subShelfModule.Controller.DeleteMySubShelvesByIds,
				),
			)...,
		)
	}
}
