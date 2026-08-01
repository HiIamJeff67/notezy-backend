package routers

import (
	"github.com/gin-gonic/gin"

	routinetasksdto "github.com/HiIamJeff67/notezy-backend/contracts/api/v1/routine-tasks"
	endpoints "github.com/HiIamJeff67/notezy-backend/internal/services/core/transports/gateway/endpoints"
	middlewares "github.com/HiIamJeff67/notezy-backend/internal/services/core/transports/gateway/middlewares"
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
				routinetasksdto.GetMyRoutineTaskByIdOperation,
			),
			authMiddleware,
			endpoint.GetMyRoutineTaskById,
		)
		routineTaskRoutes.POST(
			"/get-all-by-routine-ids",
			middlewares.DelegationAuthenticatedMiddleware(
				routinetasksdto.GetAllMyRoutineTasksByRoutineIdsOperation,
			),
			authMiddleware,
			endpoint.GetAllMyRoutineTasksByRoutineIds,
		)
		routineTaskRoutes.POST(
			"/get-all",
			middlewares.DelegationAuthenticatedMiddleware(
				routinetasksdto.GetAllMyRoutineTasksOperation,
			),
			authMiddleware,
			endpoint.GetAllMyRoutineTasks,
		)
		routineTaskRoutes.POST(
			"/create-by-routine-id",
			middlewares.DelegationAuthenticatedMiddleware(
				routinetasksdto.CreateRoutineTaskByRoutineIdOperation,
			),
			authMiddleware,
			endpoint.CreateRoutineTaskByRoutineId,
		)
		routineTaskRoutes.POST(
			"/update",
			middlewares.DelegationAuthenticatedMiddleware(
				routinetasksdto.UpdateMyRoutineTaskByIdOperation,
			),
			authMiddleware,
			endpoint.UpdateMyRoutineTaskById,
		)
		routineTaskRoutes.POST(
			"/pause",
			middlewares.DelegationAuthenticatedMiddleware(
				routinetasksdto.PauseMyRoutineTaskByIdOperation,
			),
			authMiddleware,
			endpoint.PauseMyRoutineTaskById,
		)
		routineTaskRoutes.POST(
			"/resume",
			middlewares.DelegationAuthenticatedMiddleware(
				routinetasksdto.ResumeMyRoutineTaskByIdOperation,
			),
			authMiddleware,
			endpoint.ResumeMyRoutineTaskById,
		)
		routineTaskRoutes.POST(
			"/hard-delete",
			middlewares.DelegationAuthenticatedMiddleware(
				routinetasksdto.HardDeleteMyRoutineTaskByIdOperation,
			),
			authMiddleware,
			endpoint.HardDeleteMyRoutineTaskById,
		)
		routineTaskRoutes.POST(
			"/hard-delete-many",
			middlewares.DelegationAuthenticatedMiddleware(
				routinetasksdto.HardDeleteMyRoutineTasksByIdsOperation,
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
				routinetasksdto.VisualizeMyRoutineTaskStatusCountOperation,
			),
			authMiddleware,
			endpoint.VisualizeMyRoutineTaskStatusCount,
		)
		visualizationRoutes.POST(
			"/visualize-purpose-count",
			middlewares.DelegationAuthenticatedMiddleware(
				routinetasksdto.VisualizeMyRoutineTaskPurposeCountOperation,
			),
			authMiddleware,
			endpoint.VisualizeMyRoutineTaskPurposeCount,
		)
		visualizationRoutes.POST(
			"/visualize-scheduled-at-count",
			middlewares.DelegationAuthenticatedMiddleware(
				routinetasksdto.VisualizeMyRoutineTaskScheduledAtCountOperation,
			),
			authMiddleware,
			endpoint.VisualizeMyRoutineTaskScheduledAtCount,
		)
		visualizationRoutes.POST(
			"/visualize-actual-started-at-count",
			middlewares.DelegationAuthenticatedMiddleware(
				routinetasksdto.VisualizeMyRoutineTaskActualStartedAtCountOperation,
			),
			authMiddleware,
			endpoint.VisualizeMyRoutineTaskActualStartedAtCount,
		)
		visualizationRoutes.POST(
			"/visualize-actual-ended-at-count",
			middlewares.DelegationAuthenticatedMiddleware(
				routinetasksdto.VisualizeMyRoutineTaskActualEndedAtCountOperation,
			),
			authMiddleware,
			endpoint.VisualizeMyRoutineTaskActualEndedAtCount,
		)
	}
	router.POST(
		"/routine-tasks/graphql/search",
		middlewares.DelegationAuthenticatedMiddleware(
			routinetasksdto.SearchRoutineTasksOperation,
		),
		authMiddleware,
		endpoint.SearchRoutineTasks,
	)
}
