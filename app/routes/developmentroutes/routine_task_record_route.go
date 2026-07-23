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
		routineTaskRecordRoutes.GET(
			"/visualizeMyRoutineTaskRecordStatusCount",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("visualizeMyRoutineTaskRecordStatusCount"),
					middlewares.ApplyMeterMiddleware("server.requests.routineTaskRecord.visualizeMyRoutineTaskRecordStatusCount"),
				},
				append(
					defaultMiddlewares,
					middlewares.AllowedPermissionsAbove(enums.AccessControlPermission_Read),
				),
				routineTaskRecordModule.Binder.BindVisualizeMyRoutineTaskRecordStatusCount(
					routineTaskRecordModule.Controller.VisualizeMyRoutineTaskRecordStatusCount,
				),
			)...,
		)
		routineTaskRecordRoutes.GET(
			"/visualizeMyRoutineTaskRecordPurposeCount",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("visualizeMyRoutineTaskRecordPurposeCount"),
					middlewares.ApplyMeterMiddleware("server.requests.routineTaskRecord.visualizeMyRoutineTaskRecordPurposeCount"),
				},
				append(
					defaultMiddlewares,
					middlewares.AllowedPermissionsAbove(enums.AccessControlPermission_Read),
				),
				routineTaskRecordModule.Binder.BindVisualizeMyRoutineTaskRecordPurposeCount(
					routineTaskRecordModule.Controller.VisualizeMyRoutineTaskRecordPurposeCount,
				),
			)...,
		)
		routineTaskRecordRoutes.GET(
			"/visualizeMyRoutineTaskRecordScheduledAtCount",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("visualizeMyRoutineTaskRecordScheduledAtCount"),
					middlewares.ApplyMeterMiddleware("server.requests.routineTaskRecord.visualizeMyRoutineTaskRecordScheduledAtCount"),
				},
				append(
					defaultMiddlewares,
					middlewares.AllowedPermissionsAbove(enums.AccessControlPermission_Read),
				),
				routineTaskRecordModule.Binder.BindVisualizeMyRoutineTaskRecordScheduledAtCount(
					routineTaskRecordModule.Controller.VisualizeMyRoutineTaskRecordScheduledAtCount,
				),
			)...,
		)
		routineTaskRecordRoutes.GET(
			"/visualizeMyRoutineTaskRecordActualStartedAtCount",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("visualizeMyRoutineTaskRecordActualStartedAtCount"),
					middlewares.ApplyMeterMiddleware("server.requests.routineTaskRecord.visualizeMyRoutineTaskRecordActualStartedAtCount"),
				},
				append(
					defaultMiddlewares,
					middlewares.AllowedPermissionsAbove(enums.AccessControlPermission_Read),
				),
				routineTaskRecordModule.Binder.BindVisualizeMyRoutineTaskRecordActualStartedAtCount(
					routineTaskRecordModule.Controller.VisualizeMyRoutineTaskRecordActualStartedAtCount,
				),
			)...,
		)
		routineTaskRecordRoutes.GET(
			"/visualizeMyRoutineTaskRecordActualEndedAtCount",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("visualizeMyRoutineTaskRecordActualEndedAtCount"),
					middlewares.ApplyMeterMiddleware("server.requests.routineTaskRecord.visualizeMyRoutineTaskRecordActualEndedAtCount"),
				},
				append(
					defaultMiddlewares,
					middlewares.AllowedPermissionsAbove(enums.AccessControlPermission_Read),
				),
				routineTaskRecordModule.Binder.BindVisualizeMyRoutineTaskRecordActualEndedAtCount(
					routineTaskRecordModule.Controller.VisualizeMyRoutineTaskRecordActualEndedAtCount,
				),
			)...,
		)
	}
}
