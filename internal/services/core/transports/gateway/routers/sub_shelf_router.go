package routers

import (
	"github.com/gin-gonic/gin"

	subshelvesdto "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api/sub-shelves"
	endpoints "github.com/HiIamJeff67/notezy-backend/internal/services/core/transports/gateway/endpoints"
	middlewares "github.com/HiIamJeff67/notezy-backend/internal/services/core/transports/gateway/middlewares"
)

func configureSubShelfRoutes(
	router *gin.RouterGroup,
	authMiddleware gin.HandlerFunc,
	endpoint endpoints.SubShelfEndpointInterface,
) {
	subShelfRoutes := router.Group("/sub-shelves")
	{
		subShelfRoutes.POST(
			"/get-by-id",
			middlewares.DelegationAuthenticatedMiddleware(
				subshelvesdto.GetMySubShelfByIdOperation,
			),
			authMiddleware,
			endpoint.GetMySubShelfById,
		)
		subShelfRoutes.POST(
			"/get-by-prev-sub-shelf-id",
			middlewares.DelegationAuthenticatedMiddleware(
				subshelvesdto.GetMySubShelvesByPrevSubShelfIdOperation,
			),
			authMiddleware,
			endpoint.GetMySubShelvesByPrevSubShelfId,
		)
		subShelfRoutes.POST(
			"/get-all-by-root-shelf-id",
			middlewares.DelegationAuthenticatedMiddleware(
				subshelvesdto.GetAllMySubShelvesByRootShelfIdOperation,
			),
			authMiddleware,
			endpoint.GetAllMySubShelvesByRootShelfId,
		)
		subShelfRoutes.POST(
			"/get-and-items-by-prev-sub-shelf-id",
			middlewares.DelegationAuthenticatedMiddleware(
				subshelvesdto.GetMySubShelvesAndItemsByPrevSubShelfIdOperation,
			),
			authMiddleware,
			endpoint.GetMySubShelvesAndItemsByPrevSubShelfId,
		)
		subShelfRoutes.POST(
			"/create",
			middlewares.DelegationAuthenticatedMiddleware(
				subshelvesdto.CreateSubShelfByRootShelfIdOperation,
			),
			authMiddleware,
			endpoint.CreateSubShelfByRootShelfId,
		)
		subShelfRoutes.POST(
			"/create-many",
			middlewares.DelegationAuthenticatedMiddleware(
				subshelvesdto.CreateSubShelvesByRootShelfIdsOperation,
			),
			authMiddleware,
			endpoint.CreateSubShelvesByRootShelfIds,
		)
		subShelfRoutes.POST(
			"/update",
			middlewares.DelegationAuthenticatedMiddleware(
				subshelvesdto.UpdateMySubShelfByIdOperation,
			),
			authMiddleware,
			endpoint.UpdateMySubShelfById,
		)
		subShelfRoutes.POST(
			"/update-many",
			middlewares.DelegationAuthenticatedMiddleware(
				subshelvesdto.UpdateMySubShelvesByIdsOperation,
			),
			authMiddleware,
			endpoint.UpdateMySubShelvesByIds,
		)
		subShelfRoutes.POST(
			"/move",
			middlewares.DelegationAuthenticatedMiddleware(
				subshelvesdto.MoveMySubShelfByRootShelfIdOperation,
			),
			authMiddleware,
			endpoint.MoveMySubShelfByRootShelfId,
		)
		subShelfRoutes.POST(
			"/move-many",
			middlewares.DelegationAuthenticatedMiddleware(
				subshelvesdto.MoveMySubShelvesByRootShelfIdOperation,
			),
			authMiddleware,
			endpoint.MoveMySubShelvesByRootShelfId,
		)
		subShelfRoutes.POST(
			"/move-many-by-root-shelves",
			middlewares.DelegationAuthenticatedMiddleware(
				subshelvesdto.MoveMySubShelvesByRootShelfIdsOperation,
			),
			authMiddleware,
			endpoint.MoveMySubShelvesByRootShelfIds,
		)
		subShelfRoutes.POST(
			"/restore",
			middlewares.DelegationAuthenticatedMiddleware(
				subshelvesdto.RestoreMySubShelfByIdOperation,
			),
			authMiddleware,
			endpoint.RestoreMySubShelfById,
		)
		subShelfRoutes.POST(
			"/restore-many",
			middlewares.DelegationAuthenticatedMiddleware(
				subshelvesdto.RestoreMySubShelvesByIdsOperation,
			),
			authMiddleware,
			endpoint.RestoreMySubShelvesByIds,
		)
		subShelfRoutes.POST(
			"/delete",
			middlewares.DelegationAuthenticatedMiddleware(
				subshelvesdto.DeleteMySubShelfByIdOperation,
			),
			authMiddleware,
			endpoint.DeleteMySubShelfById,
		)
		subShelfRoutes.POST(
			"/delete-many",
			middlewares.DelegationAuthenticatedMiddleware(
				subshelvesdto.DeleteMySubShelvesByIdsOperation,
			),
			authMiddleware,
			endpoint.DeleteMySubShelvesByIds,
		)
		subShelfRoutes.POST(
			"/graphql/search",
			middlewares.DelegationAuthenticatedMiddleware(
				subshelvesdto.SearchSubShelvesOperation,
			),
			authMiddleware,
			endpoint.SearchSubShelves,
		)
	}
}
