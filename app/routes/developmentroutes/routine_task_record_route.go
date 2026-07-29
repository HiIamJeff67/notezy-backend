package developmentroutes

import (
	"time"

	"github.com/gin-gonic/gin"

	interceptors "github.com/HiIamJeff67/notezy-backend/app/interceptors"
	middlewares "github.com/HiIamJeff67/notezy-backend/app/middlewares"
	enums "github.com/HiIamJeff67/notezy-backend/app/models/schemas/enums"
	modules "github.com/HiIamJeff67/notezy-backend/app/modules"
)

func configureDevelopmentRoutineTaskRecordRoutes(router *gin.RouterGroup) {
	if router == nil {
		router = DevelopmentAPIRouterGroup
	}

	routineTaskRecordModule := modules.NewRoutineTaskRecordModule()

	routineTaskRecordRouterGroup := router.Group("/routine-task-records")
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
		routineTaskRecordRouterGroup.GET(
			"/routine-task/:routineTaskId",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("getAllMyRoutineTaskRecordsByRoutineTaskId"),
					middlewares.ApplyMeterMiddleware("server.requests.routineTaskRecord.getAllMyRoutineTaskRecordsByRoutineTaskId"),
				},
				append(
					defaultMiddlewares,
					middlewares.AllowedPermissionsAbove(enums.AccessControlPermission_Read),
				),
				routineTaskRecordModule.Binder.BindGetAllMyRoutineTaskRecordsByRoutineTaskId(
					routineTaskRecordModule.Controller.GetAllMyRoutineTaskRecordsByRoutineTaskId,
				),
			)...,
		)
	}

	/* ============================== Routes for Visualization ============================== */

	visualizationRoutes := router.Group("/routine-task-records/visualizations")
	visualizationMiddlewares := []gin.HandlerFunc{
		middlewares.UnauthorizedRateLimitMiddleware(),
		middlewares.TimeoutMiddleware(3 * time.Second),
		middlewares.AuthMiddleware(),
		interceptors.ShareableResponseWriterInterceptor(
			interceptors.RefreshTokenInterceptor,
			interceptors.EmbeddedInterceptor,
		),
	}
	{
		visualizationRoutes.GET(
			"/status-count",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("visualizeMyRoutineTaskRecordStatusCount"),
					middlewares.ApplyMeterMiddleware("server.requests.routineTaskRecord.visualizeMyRoutineTaskRecordStatusCount"),
				},
				append(
					visualizationMiddlewares,
					middlewares.AllowedPermissionsAbove(enums.AccessControlPermission_Read),
				),
				routineTaskRecordModule.Binder.BindVisualizeMyRoutineTaskRecordStatusCount(
					routineTaskRecordModule.Controller.VisualizeMyRoutineTaskRecordStatusCount,
				),
			)...,
		)
		visualizationRoutes.GET(
			"/purpose-count",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("visualizeMyRoutineTaskRecordPurposeCount"),
					middlewares.ApplyMeterMiddleware("server.requests.routineTaskRecord.visualizeMyRoutineTaskRecordPurposeCount"),
				},
				append(
					visualizationMiddlewares,
					middlewares.AllowedPermissionsAbove(enums.AccessControlPermission_Read),
				),
				routineTaskRecordModule.Binder.BindVisualizeMyRoutineTaskRecordPurposeCount(
					routineTaskRecordModule.Controller.VisualizeMyRoutineTaskRecordPurposeCount,
				),
			)...,
		)
		visualizationRoutes.GET(
			"/scheduled-at-count",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("visualizeMyRoutineTaskRecordScheduledAtCount"),
					middlewares.ApplyMeterMiddleware("server.requests.routineTaskRecord.visualizeMyRoutineTaskRecordScheduledAtCount"),
				},
				append(
					visualizationMiddlewares,
					middlewares.AllowedPermissionsAbove(enums.AccessControlPermission_Read),
				),
				routineTaskRecordModule.Binder.BindVisualizeMyRoutineTaskRecordScheduledAtCount(
					routineTaskRecordModule.Controller.VisualizeMyRoutineTaskRecordScheduledAtCount,
				),
			)...,
		)
		visualizationRoutes.GET(
			"/actual-started-at-count",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("visualizeMyRoutineTaskRecordActualStartedAtCount"),
					middlewares.ApplyMeterMiddleware("server.requests.routineTaskRecord.visualizeMyRoutineTaskRecordActualStartedAtCount"),
				},
				append(
					visualizationMiddlewares,
					middlewares.AllowedPermissionsAbove(enums.AccessControlPermission_Read),
				),
				routineTaskRecordModule.Binder.BindVisualizeMyRoutineTaskRecordActualStartedAtCount(
					routineTaskRecordModule.Controller.VisualizeMyRoutineTaskRecordActualStartedAtCount,
				),
			)...,
		)
		visualizationRoutes.GET(
			"/actual-ended-at-count",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("visualizeMyRoutineTaskRecordActualEndedAtCount"),
					middlewares.ApplyMeterMiddleware("server.requests.routineTaskRecord.visualizeMyRoutineTaskRecordActualEndedAtCount"),
				},
				append(
					visualizationMiddlewares,
					middlewares.AllowedPermissionsAbove(enums.AccessControlPermission_Read),
				),
				routineTaskRecordModule.Binder.BindVisualizeMyRoutineTaskRecordActualEndedAtCount(
					routineTaskRecordModule.Controller.VisualizeMyRoutineTaskRecordActualEndedAtCount,
				),
			)...,
		)
	}
}
