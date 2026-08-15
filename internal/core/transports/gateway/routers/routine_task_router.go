package routers

import (
	"github.com/gin-gonic/gin"

	apicontract "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api/routine-tasks"

	contexts "github.com/HiIamJeff67/notezy-backend/internal/core/contexts"
	routineservices "github.com/HiIamJeff67/notezy-backend/internal/core/services/routines"
	endpoints "github.com/HiIamJeff67/notezy-backend/internal/core/transports/gateway/endpoints"
	middlewares "github.com/HiIamJeff67/notezy-backend/internal/core/transports/gateway/middlewares"
)

type RoutineTaskRouterDependencies struct {
	Service          routineservices.RoutineTaskServiceInterface
	AuthMiddleware   gin.HandlerFunc
	APIKeyMiddleware gin.HandlerFunc
}

func configureRoutineTaskRoutes(
	router *gin.RouterGroup,
	deps RoutineTaskRouterDependencies,
) {
	authMiddleware := deps.AuthMiddleware
	apiKeyMiddleware := deps.APIKeyMiddleware
	endpoint := endpoints.NewRoutineTaskEndpoint(deps.Service)
	apiCompatibleAuthMiddleware := middlewares.EitherMiddleware(
		[]gin.HandlerFunc{authMiddleware},
		[]gin.HandlerFunc{apiKeyMiddleware},
		func(ctx *gin.Context) bool { return contexts.IsClientGateway(ctx.Request.Context()) },
	)[0]

	routineTaskRoutes := router.Group("/routine-tasks")
	{
		routineTaskRoutes.POST(
			"/get-by-id",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.GetMyRoutineTaskByIdOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.GetMyRoutineTaskById,
		)
		routineTaskRoutes.POST(
			"/get-all-by-routine-ids",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.GetAllMyRoutineTasksByRoutineIdsOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.GetAllMyRoutineTasksByRoutineIds,
		)
		routineTaskRoutes.POST(
			"/get-all",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.GetAllMyRoutineTasksOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.GetAllMyRoutineTasks,
		)
		routineTaskRoutes.POST(
			"/create-by-routine-id",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.CreateRoutineTaskByRoutineIdOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.CreateRoutineTaskByRoutineId,
		)
		routineTaskRoutes.POST(
			"/update",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.UpdateMyRoutineTaskByIdOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.UpdateMyRoutineTaskById,
		)
		routineTaskRoutes.POST(
			"/pause",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.PauseMyRoutineTaskByIdOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.PauseMyRoutineTaskById,
		)
		routineTaskRoutes.POST(
			"/resume",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.ResumeMyRoutineTaskByIdOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.ResumeMyRoutineTaskById,
		)
		routineTaskRoutes.POST(
			"/hard-delete",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.HardDeleteMyRoutineTaskByIdOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.HardDeleteMyRoutineTaskById,
		)
		routineTaskRoutes.POST(
			"/hard-delete-many",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.HardDeleteMyRoutineTasksByIdsOperation,
			),
			apiCompatibleAuthMiddleware,
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
			apiCompatibleAuthMiddleware,
			endpoint.VisualizeMyRoutineTaskStatusCount,
		)
		visualizationRoutes.POST(
			"/visualize-purpose-count",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.VisualizeMyRoutineTaskPurposeCountOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.VisualizeMyRoutineTaskPurposeCount,
		)
		visualizationRoutes.POST(
			"/visualize-scheduled-at-count",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.VisualizeMyRoutineTaskScheduledAtCountOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.VisualizeMyRoutineTaskScheduledAtCount,
		)
		visualizationRoutes.POST(
			"/visualize-actual-started-at-count",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.VisualizeMyRoutineTaskActualStartedAtCountOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.VisualizeMyRoutineTaskActualStartedAtCount,
		)
		visualizationRoutes.POST(
			"/visualize-actual-ended-at-count",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.VisualizeMyRoutineTaskActualEndedAtCountOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.VisualizeMyRoutineTaskActualEndedAtCount,
		)
	}
	router.POST(
		"/routine-tasks/graphql/search",
		middlewares.DelegationAuthenticatedMiddleware(
			apicontract.SearchRoutineTasksOperation,
		),
		apiCompatibleAuthMiddleware,
		endpoint.SearchRoutineTasks,
	)
}
