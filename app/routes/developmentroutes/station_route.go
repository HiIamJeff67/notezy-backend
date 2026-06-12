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

func configureDevelopmentStationRoutes() {
	stationModule := modules.NewStationModule()

	stationRoutes := DevelopmentRouterGroup.Group("/station")
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
		stationRoutes.GET(
			"/getMyStationById",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "getMyStationById"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.Station.GetMyStationById,
					),
				},
				defaultMiddlewares,
				stationModule.Binder.BindGetMyStationById(
					stationModule.Controller.GetMyStationById,
				),
			)...,
		)
		stationRoutes.GET(
			"/getAllMyStations",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "getAllMyStations"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.Station.GetAllMyStations,
					),
				},
				defaultMiddlewares,
				stationModule.Binder.BindGetAllMyStations(
					stationModule.Controller.GetAllMyStations,
				),
			)...,
		)
		stationRoutes.POST(
			"/createStation",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "createStation"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.Station.CreateStation,
					),
				},
				defaultMiddlewares,
				stationModule.Binder.BindCreateStation(
					stationModule.Controller.CreateStation,
				),
			)...,
		)
		stationRoutes.POST(
			"/createStations",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "createStations"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.Station.CreateStations,
					),
				},
				defaultMiddlewares,
				stationModule.Binder.BindCreateStations(
					stationModule.Controller.CreateStations,
				),
			)...,
		)
		stationRoutes.PUT(
			"/updateMyStationById",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "updateMyStationById"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.Station.UpdateMyStationById,
					),
				},
				defaultMiddlewares,
				stationModule.Binder.BindUpdateMyStationById(
					stationModule.Controller.UpdateMyStationById,
				),
			)...,
		)
		stationRoutes.PUT(
			"/updateMyStationsByIds",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "updateMyStationsByIds"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.Station.UpdateMyStationsByIds,
					),
				},
				defaultMiddlewares,
				stationModule.Binder.BindUpdateMyStationsByIds(
					stationModule.Controller.UpdateMyStationsByIds,
				),
			)...,
		)
		stationRoutes.PATCH(
			"/restoreMyStationById",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "restoreMyStationById"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.Station.RestoreMyStationById,
					),
				},
				defaultMiddlewares,
				stationModule.Binder.BindRestoreMyStationById(
					stationModule.Controller.RestoreMyStationById,
				),
			)...,
		)
		stationRoutes.PATCH(
			"/restoreMyStationsByIds",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "restoreMyStationsByIds"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.Station.RestoreMyStationsByIds,
					),
				},
				defaultMiddlewares,
				stationModule.Binder.BindRestoreMyStationsByIds(
					stationModule.Controller.RestoreMyStationsByIds,
				),
			)...,
		)
		stationRoutes.DELETE(
			"/deleteMyStationById",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "deleteMyStationById"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.Station.DeleteMyStationById,
					),
				},
				defaultMiddlewares,
				stationModule.Binder.BindDeleteMyStationById(
					stationModule.Controller.DeleteMyStationById,
				),
			)...,
		)
		stationRoutes.DELETE(
			"/deleteMyStationsByIds",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "deleteMyStationsByIds"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.Station.DeleteMyStationsByIds,
					),
				},
				defaultMiddlewares,
				stationModule.Binder.BindDeleteMyStationsByIds(
					stationModule.Controller.DeleteMyStationsByIds,
				),
			)...,
		)
		stationRoutes.DELETE(
			"/hardDeleteMyStationById",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "hardDeleteMyStationById"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.Station.HardDeleteMyStationById,
					),
				},
				defaultMiddlewares,
				stationModule.Binder.BindHardDeleteMyStationById(
					stationModule.Controller.HardDeleteMyStationById,
				),
			)...,
		)
		stationRoutes.DELETE(
			"/hardDeleteMyStationsByIds",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "hardDeleteMyStationsByIds"),
					middlewares.ApplyMeterMiddleware(
						otel.Meter(constants.ServiceName),
						metrics.MetricNames.Server.Requests.Station.HardDeleteMyStationsByIds,
					),
				},
				defaultMiddlewares,
				stationModule.Binder.BindHardDeleteMyStationsByIds(
					stationModule.Controller.HardDeleteMyStationsByIds,
				),
			)...,
		)
	}
}
