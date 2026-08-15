package routers

import (
	"github.com/gin-gonic/gin"

	apicontract "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api/routines"

	contexts "github.com/HiIamJeff67/notezy-backend/internal/core/contexts"
	routineservices "github.com/HiIamJeff67/notezy-backend/internal/core/services/routines"
	endpoints "github.com/HiIamJeff67/notezy-backend/internal/core/transports/gateway/endpoints"
	middlewares "github.com/HiIamJeff67/notezy-backend/internal/core/transports/gateway/middlewares"
)

type RoutineRouterDependencies struct {
	Service          routineservices.RoutineServiceInterface
	AuthMiddleware   gin.HandlerFunc
	APIKeyMiddleware gin.HandlerFunc
}

func configureRoutineRoutes(
	router *gin.RouterGroup,
	deps RoutineRouterDependencies,
) {
	authMiddleware := deps.AuthMiddleware
	apiKeyMiddleware := deps.APIKeyMiddleware
	endpoint := endpoints.NewRoutineEndpoint(deps.Service)
	apiCompatibleAuthMiddleware := middlewares.EitherMiddleware(
		[]gin.HandlerFunc{authMiddleware},
		[]gin.HandlerFunc{apiKeyMiddleware},
		func(ctx *gin.Context) bool { return contexts.IsClientGateway(ctx.Request.Context()) },
	)[0]

	routineRoutes := router.Group("/routines")
	{
		routineRoutes.POST(
			"/get-by-id",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.GetMyRoutineByIdOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.GetMyRoutineById,
		)
		routineRoutes.POST(
			"/get-by-station-id",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.GetMyRoutinesByStationIdOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.GetMyRoutinesByStationId,
		)
		routineRoutes.POST(
			"/get-all-by-time-range",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.GetAllMyRoutinesByTimeRangeOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.GetAllMyRoutinesByTimeRange,
		)
		routineRoutes.POST(
			"/create-by-station-id",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.CreateRoutineByStationIdOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.CreateRoutineByStationId,
		)
		routineRoutes.POST(
			"/create-many-by-station-ids",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.CreateRoutinesByStationIdsOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.CreateRoutinesByStationIds,
		)
		routineRoutes.POST(
			"/update",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.UpdateMyRoutineByIdOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.UpdateMyRoutineById,
		)
		routineRoutes.POST(
			"/update-many",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.UpdateMyRoutinesByIdsOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.UpdateMyRoutinesByIds,
		)
		routineRoutes.POST(
			"/link-tag",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.LinkRoutineTagByIdOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.LinkRoutineTagById,
		)
		routineRoutes.POST(
			"/link-tags",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.LinkRoutineTagsByIdsOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.LinkRoutineTagsByIds,
		)
		routineRoutes.POST(
			"/link-item",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.LinkRoutineItemByIdOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.LinkRoutineItemById,
		)
		routineRoutes.POST(
			"/link-items",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.LinkRoutineItemsByIdsOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.LinkRoutineItemsByIds,
		)
		routineRoutes.POST(
			"/restore",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.RestoreMyRoutineByIdOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.RestoreMyRoutineById,
		)
		routineRoutes.POST(
			"/restore-many",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.RestoreMyRoutinesByIdsOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.RestoreMyRoutinesByIds,
		)
		routineRoutes.POST(
			"/delete",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.DeleteMyRoutineByIdOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.DeleteMyRoutineById,
		)
		routineRoutes.POST(
			"/delete-many",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.DeleteMyRoutinesByIdsOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.DeleteMyRoutinesByIds,
		)
		routineRoutes.POST(
			"/hard-delete",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.HardDeleteMyRoutineByIdOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.HardDeleteMyRoutineById,
		)
		routineRoutes.POST(
			"/hard-delete-many",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.HardDeleteMyRoutinesByIdsOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.HardDeleteMyRoutinesByIds,
		)
	}
	visualizationRoutes := router.Group("/routines/visualizations")
	{
		visualizationRoutes.POST(
			"/status-count",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.VisualizeMyRoutineStatusCountOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.VisualizeMyRoutineStatusCount,
		)
		visualizationRoutes.POST(
			"/period-count",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.VisualizeMyRoutinePeriodCountOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.VisualizeMyRoutinePeriodCount,
		)
		visualizationRoutes.POST(
			"/scheduled-start-at-count",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.VisualizeMyRoutineScheduledStartAtCountOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.VisualizeMyRoutineScheduledStartAtCount,
		)
		visualizationRoutes.POST(
			"/scheduled-end-at-count",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.VisualizeMyRoutineScheduledEndAtCountOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.VisualizeMyRoutineScheduledEndAtCount,
		)
	}
	router.POST(
		"/routines/graphql/search",
		middlewares.DelegationAuthenticatedMiddleware(
			apicontract.SearchRoutinesOperation,
		),
		apiCompatibleAuthMiddleware,
		endpoint.SearchRoutines,
	)
}
