package developmentroutes

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"

	adapters "notezy-backend/app/adapters"
	interceptors "notezy-backend/app/interceptors"
	middlewares "notezy-backend/app/middlewares"
	modules "notezy-backend/app/modules"
	metrics "notezy-backend/app/monitor/metrics"
	constants "notezy-backend/shared/constants"
)

func configureDevelopmentMaterialRoutes() {
	materialModule := modules.NewMaterialModule()

	materialRoutes := DevelopmentRouterGroup.Group("/material")
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
		materialRoutes.GET(
			"/getMyMaterialById",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "getMyMaterialById"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.Material.GetMyMaterialById,
					),
				},
				defaultMiddlewares,
				materialModule.Binder.BindGetMyMaterialById(
					materialModule.Controller.GetMyMaterialById,
				),
			)...,
		)
		materialRoutes.GET(
			"/getMyMaterialAndItsParentById",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "getMyMaterialAndItsParentById"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.Material.GetMyMaterialAndItsParentById,
					),
				},
				defaultMiddlewares,
				materialModule.Binder.BindGetMyMaterialAndItsParentById(
					materialModule.Controller.GetMyMaterialAndItsParentById,
				),
			)...,
		)
		materialRoutes.GET(
			"/getMyMaterialsByParentSubShelfId",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "getMyMaterialsByParentSubShelfId"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.Material.GetMyMaterialsByParentSubShelfId,
					),
				},
				defaultMiddlewares,
				materialModule.Binder.BindGetMyMaterialsByParentSubShelfId(
					materialModule.Controller.GetMyMaterialsByParentSubShelfId,
				),
			)...,
		)
		materialRoutes.GET(
			"/getAllMyMaterialsByRootShelfId",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "getAllMyMaterialsByRootShelfId"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.Material.GetAllMyMaterialsByRootShelfId,
					),
				},
				defaultMiddlewares,
				materialModule.Binder.BindGetAllMyMaterialsByRootShelfId(
					materialModule.Controller.GetAllMyMaterialsByRootShelfId,
				),
			)...,
		)
		materialRoutes.POST(
			"/createTextbookMaterial",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "createTextbookMaterial"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.Material.CreateTextbookMaterial,
					),
				},
				defaultMiddlewares,
				materialModule.Binder.BindCreateTextbookMaterial(
					materialModule.Controller.CreateTextbookMaterial,
				),
			)...,
		)
		materialRoutes.POST(
			"/createNotebookMaterial",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "createNotebookMaterial"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.Material.CreateNotebookMaterial,
					),
				},
				defaultMiddlewares,
				materialModule.Binder.BindCreateNotebookMaterial(
					materialModule.Controller.CreateNotebookMaterial,
				),
			)...,
		)
		materialRoutes.PUT(
			"/updateMyMaterialById",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "updateMyMaterialById"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.Material.UpdateMyMaterialById,
					),
				},
				defaultMiddlewares,
				materialModule.Binder.BindUpdateMyMaterialById(
					materialModule.Controller.UpdateMyMaterialById,
				),
			)...,
		)
		materialRoutes.PUT(
			"/saveMyNotebookMaterialById",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "saveMyNotebookMaterialById"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.Material.SaveMyNotebookMaterialById,
					),
					adapters.MultipartAdapter(),
				},
				defaultMiddlewares,
				materialModule.Binder.BindSaveMyNotebookMaterialById(
					materialModule.Controller.SaveMyNotebookMaterialById,
				),
			)...,
		)
		materialRoutes.PUT(
			"/moveMyMaterialById",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "moveMyMaterialById"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.Material.MoveMyMaterialById,
					),
				},
				defaultMiddlewares,
				materialModule.Binder.BindMoveMyMaterialById(
					materialModule.Controller.MoveMyMaterialById,
				),
			)...,
		)
		materialRoutes.PUT(
			"/moveMyMaterialsByIds",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "moveMyMaterialsByIds"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.Material.MoveMyMaterialsByIds,
					),
				},
				defaultMiddlewares,
				materialModule.Binder.BindMoveMyMaterialsByIds(
					materialModule.Controller.MoveMyMaterialsByIds,
				),
			)...,
		)
		materialRoutes.PATCH(
			"/restoreMyMaterialById",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "restoreMyMaterialById"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.Material.RestoreMyMaterialById,
					),
				},
				defaultMiddlewares,
				materialModule.Binder.BindRestoreMyMaterialById(
					materialModule.Controller.RestoreMyMaterialById,
				),
			)...,
		)
		materialRoutes.PATCH(
			"/restoreMyMaterialsByIds",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "restoreMyMaterialsByIds"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.Material.RestoreMyMaterialsByIds,
					),
				},
				defaultMiddlewares,
				materialModule.Binder.BindRestoreMyMaterialsByIds(
					materialModule.Controller.RestoreMyMaterialsByIds,
				),
			)...,
		)
		materialRoutes.DELETE(
			"/deleteMyMaterialById",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "deleteMyMaterialById"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.Material.DeleteMyMaterialById,
					),
				},
				defaultMiddlewares,
				materialModule.Binder.BindDeleteMyMaterialById(
					materialModule.Controller.DeleteMyMaterialById,
				),
			)...,
		)
		materialRoutes.DELETE(
			"/deleteMyMaterialsByIds",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "deleteMyMaterialsByIds"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.Material.DeleteMyMaterialsByIds,
					),
				},
				defaultMiddlewares,
				materialModule.Binder.BindDeleteMyMaterialsByIds(
					materialModule.Controller.DeleteMyMaterialsByIds,
				),
			)...,
		)
	}
}
