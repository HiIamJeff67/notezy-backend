package routers

import (
	"github.com/gin-gonic/gin"

	routinetagsdto "github.com/HiIamJeff67/notezy-backend/contracts/api/v1/routine-tags"
	endpoints "github.com/HiIamJeff67/notezy-backend/internal/services/core/transports/gateway/endpoints"
	middlewares "github.com/HiIamJeff67/notezy-backend/internal/services/core/transports/gateway/middlewares"
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
				routinetagsdto.GetMyRoutineTagByIdOperation,
			),
			authMiddleware,
			endpoint.GetMyRoutineTagById,
		)
		routineTagRoutes.POST(
			"/get-all",
			middlewares.DelegationAuthenticatedMiddleware(
				routinetagsdto.GetAllMyRoutineTagsOperation,
			),
			authMiddleware,
			endpoint.GetAllMyRoutineTags,
		)
		routineTagRoutes.POST(
			"/create",
			middlewares.DelegationAuthenticatedMiddleware(
				routinetagsdto.CreateRoutineTagOperation,
			),
			authMiddleware,
			endpoint.CreateRoutineTag,
		)
		routineTagRoutes.POST(
			"/create-many",
			middlewares.DelegationAuthenticatedMiddleware(
				routinetagsdto.CreateRoutineTagsOperation,
			),
			authMiddleware,
			endpoint.CreateRoutineTags,
		)
		routineTagRoutes.POST(
			"/update",
			middlewares.DelegationAuthenticatedMiddleware(
				routinetagsdto.UpdateMyRoutineTagByIdOperation,
			),
			authMiddleware,
			endpoint.UpdateMyRoutineTagById,
		)
		routineTagRoutes.POST(
			"/update-many",
			middlewares.DelegationAuthenticatedMiddleware(
				routinetagsdto.UpdateMyRoutineTagsByIdsOperation,
			),
			authMiddleware,
			endpoint.UpdateMyRoutineTagsByIds,
		)
		routineTagRoutes.POST(
			"/hard-delete",
			middlewares.DelegationAuthenticatedMiddleware(
				routinetagsdto.HardDeleteMyRoutineTagByIdOperation,
			),
			authMiddleware,
			endpoint.HardDeleteMyRoutineTagById,
		)
		routineTagRoutes.POST(
			"/hard-delete-many",
			middlewares.DelegationAuthenticatedMiddleware(
				routinetagsdto.HardDeleteMyRoutineTagsByIdsOperation,
			),
			authMiddleware,
			endpoint.HardDeleteMyRoutineTagsByIds,
		)
		routineTagRoutes.POST(
			"/graphql/search",
			middlewares.DelegationAuthenticatedMiddleware(
				routinetagsdto.SearchRoutineTagsOperation,
			),
			authMiddleware,
			endpoint.SearchRoutineTags,
		)
	}
}
