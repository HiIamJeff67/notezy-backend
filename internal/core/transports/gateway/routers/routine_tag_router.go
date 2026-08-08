package routers

import (
	"github.com/gin-gonic/gin"

	apicontract "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api/routine-tags"

	endpoints "github.com/HiIamJeff67/notezy-backend/internal/core/transports/gateway/endpoints"
	middlewares "github.com/HiIamJeff67/notezy-backend/internal/core/transports/gateway/middlewares"
)

func configureRoutineTagRoutes(
	router *gin.RouterGroup,
	authMiddleware gin.HandlerFunc,
	endpoint endpoints.RoutineTagEndpointInterface,
) {
	routineTagRoutes := router.Group("/routine-tags")
	{
		routineTagRoutes.POST(
			"/get-by-id",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.GetMyRoutineTagByIdOperation,
			),
			authMiddleware,
			endpoint.GetMyRoutineTagById,
		)
		routineTagRoutes.POST(
			"/get-all",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.GetAllMyRoutineTagsOperation,
			),
			authMiddleware,
			endpoint.GetAllMyRoutineTags,
		)
		routineTagRoutes.POST(
			"/create",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.CreateRoutineTagOperation,
			),
			authMiddleware,
			endpoint.CreateRoutineTag,
		)
		routineTagRoutes.POST(
			"/create-many",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.CreateRoutineTagsOperation,
			),
			authMiddleware,
			endpoint.CreateRoutineTags,
		)
		routineTagRoutes.POST(
			"/update",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.UpdateMyRoutineTagByIdOperation,
			),
			authMiddleware,
			endpoint.UpdateMyRoutineTagById,
		)
		routineTagRoutes.POST(
			"/update-many",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.UpdateMyRoutineTagsByIdsOperation,
			),
			authMiddleware,
			endpoint.UpdateMyRoutineTagsByIds,
		)
		routineTagRoutes.POST(
			"/hard-delete",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.HardDeleteMyRoutineTagByIdOperation,
			),
			authMiddleware,
			endpoint.HardDeleteMyRoutineTagById,
		)
		routineTagRoutes.POST(
			"/hard-delete-many",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.HardDeleteMyRoutineTagsByIdsOperation,
			),
			authMiddleware,
			endpoint.HardDeleteMyRoutineTagsByIds,
		)
		routineTagRoutes.POST(
			"/graphql/search",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.SearchRoutineTagsOperation,
			),
			authMiddleware,
			endpoint.SearchRoutineTags,
		)
	}
}
