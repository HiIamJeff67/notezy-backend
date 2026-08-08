package routers

import (
	"github.com/gin-gonic/gin"

	apicontract "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api/routine-tasks"

	endpoints "github.com/HiIamJeff67/notezy-backend/internal/core/transports/gateway/endpoints"
	middlewares "github.com/HiIamJeff67/notezy-backend/internal/core/transports/gateway/middlewares"
)

func configureRoutineTaskRoutes(
	router *gin.RouterGroup,
	authMiddleware gin.HandlerFunc,
	endpoint endpoints.RoutineTaskEndpointInterface,
) {
	routineTaskRoutes := router.Group("/routine-tasks")
	{
		routineTaskRoutes.POST(
			"/get-by-id",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.GetMyRoutineTaskByIdOperation,
			),
			authMiddleware,
			endpoint.GetMyRoutineTaskById,
		)
		routineTaskRoutes.POST(
			"/get-all-by-routine-ids",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.GetAllMyRoutineTasksByRoutineIdsOperation,
			),
			authMiddleware,
			endpoint.GetAllMyRoutineTasksByRoutineIds,
		)
		routineTaskRoutes.POST(
			"/get-all",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.GetAllMyRoutineTasksOperation,
			),
			authMiddleware,
			endpoint.GetAllMyRoutineTasks,
		)
		routineTaskRoutes.POST(
			"/create-by-routine-id",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.CreateRoutineTaskByRoutineIdOperation,
			),
			authMiddleware,
			endpoint.CreateRoutineTaskByRoutineId,
		)
		routineTaskRoutes.POST(
			"/update",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.UpdateMyRoutineTaskByIdOperation,
			),
			authMiddleware,
			endpoint.UpdateMyRoutineTaskById,
		)
		routineTaskRoutes.POST(
			"/pause",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.PauseMyRoutineTaskByIdOperation,
			),
			authMiddleware,
			endpoint.PauseMyRoutineTaskById,
		)
		routineTaskRoutes.POST(
			"/resume",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.ResumeMyRoutineTaskByIdOperation,
			),
			authMiddleware,
			endpoint.ResumeMyRoutineTaskById,
		)
		routineTaskRoutes.POST(
			"/hard-delete",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.HardDeleteMyRoutineTaskByIdOperation,
			),
			authMiddleware,
			endpoint.HardDeleteMyRoutineTaskById,
		)
		routineTaskRoutes.POST(
			"/hard-delete-many",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.HardDeleteMyRoutineTasksByIdsOperation,
			),
			authMiddleware,
			endpoint.HardDeleteMyRoutineTasksByIds,
		)
	}
	visualizationRoutes := router.Group("/routine-tasks/visualizations")
	{
		visualizationRoutes.POST(
			"/visualize-status-count",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.VisualizeMyRoutineTaskStatusCountOperation,
			),
			authMiddleware,
			endpoint.VisualizeMyRoutineTaskStatusCount,
		)
		visualizationRoutes.POST(
			"/visualize-purpose-count",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.VisualizeMyRoutineTaskPurposeCountOperation,
			),
			authMiddleware,
			endpoint.VisualizeMyRoutineTaskPurposeCount,
		)
		visualizationRoutes.POST(
			"/visualize-scheduled-at-count",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.VisualizeMyRoutineTaskScheduledAtCountOperation,
			),
			authMiddleware,
			endpoint.VisualizeMyRoutineTaskScheduledAtCount,
		)
		visualizationRoutes.POST(
			"/visualize-actual-started-at-count",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.VisualizeMyRoutineTaskActualStartedAtCountOperation,
			),
			authMiddleware,
			endpoint.VisualizeMyRoutineTaskActualStartedAtCount,
		)
		visualizationRoutes.POST(
			"/visualize-actual-ended-at-count",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.VisualizeMyRoutineTaskActualEndedAtCountOperation,
			),
			authMiddleware,
			endpoint.VisualizeMyRoutineTaskActualEndedAtCount,
		)
	}
	router.POST(
		"/routine-tasks/graphql/search",
		middlewares.DelegationAuthenticatedMiddleware(
			apicontract.SearchRoutineTasksOperation,
		),
		authMiddleware,
		endpoint.SearchRoutineTasks,
	)
}
