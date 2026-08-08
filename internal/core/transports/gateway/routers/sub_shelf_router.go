package routers

import (
	"github.com/gin-gonic/gin"

	apicontract "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api/sub-shelves"

	endpoints "github.com/HiIamJeff67/notezy-backend/internal/core/transports/gateway/endpoints"
	middlewares "github.com/HiIamJeff67/notezy-backend/internal/core/transports/gateway/middlewares"
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
				apicontract.GetMySubShelfByIdOperation,
			),
			authMiddleware,
			endpoint.GetMySubShelfById,
		)
		subShelfRoutes.POST(
			"/get-by-prev-sub-shelf-id",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.GetMySubShelvesByPrevSubShelfIdOperation,
			),
			authMiddleware,
			endpoint.GetMySubShelvesByPrevSubShelfId,
		)
		subShelfRoutes.POST(
			"/get-all-by-root-shelf-id",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.GetAllMySubShelvesByRootShelfIdOperation,
			),
			authMiddleware,
			endpoint.GetAllMySubShelvesByRootShelfId,
		)
		subShelfRoutes.POST(
			"/get-and-items-by-prev-sub-shelf-id",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.GetMySubShelvesAndItemsByPrevSubShelfIdOperation,
			),
			authMiddleware,
			endpoint.GetMySubShelvesAndItemsByPrevSubShelfId,
		)
		subShelfRoutes.POST(
			"/create",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.CreateSubShelfByRootShelfIdOperation,
			),
			authMiddleware,
			endpoint.CreateSubShelfByRootShelfId,
		)
		subShelfRoutes.POST(
			"/create-many",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.CreateSubShelvesByRootShelfIdsOperation,
			),
			authMiddleware,
			endpoint.CreateSubShelvesByRootShelfIds,
		)
		subShelfRoutes.POST(
			"/update",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.UpdateMySubShelfByIdOperation,
			),
			authMiddleware,
			endpoint.UpdateMySubShelfById,
		)
		subShelfRoutes.POST(
			"/update-many",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.UpdateMySubShelvesByIdsOperation,
			),
			authMiddleware,
			endpoint.UpdateMySubShelvesByIds,
		)
		subShelfRoutes.POST(
			"/move",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.MoveMySubShelfByRootShelfIdOperation,
			),
			authMiddleware,
			endpoint.MoveMySubShelfByRootShelfId,
		)
		subShelfRoutes.POST(
			"/move-many",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.MoveMySubShelvesByRootShelfIdOperation,
			),
			authMiddleware,
			endpoint.MoveMySubShelvesByRootShelfId,
		)
		subShelfRoutes.POST(
			"/move-many-by-root-shelves",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.MoveMySubShelvesByRootShelfIdsOperation,
			),
			authMiddleware,
			endpoint.MoveMySubShelvesByRootShelfIds,
		)
		subShelfRoutes.POST(
			"/restore",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.RestoreMySubShelfByIdOperation,
			),
			authMiddleware,
			endpoint.RestoreMySubShelfById,
		)
		subShelfRoutes.POST(
			"/restore-many",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.RestoreMySubShelvesByIdsOperation,
			),
			authMiddleware,
			endpoint.RestoreMySubShelvesByIds,
		)
		subShelfRoutes.POST(
			"/delete",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.DeleteMySubShelfByIdOperation,
			),
			authMiddleware,
			endpoint.DeleteMySubShelfById,
		)
		subShelfRoutes.POST(
			"/delete-many",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.DeleteMySubShelvesByIdsOperation,
			),
			authMiddleware,
			endpoint.DeleteMySubShelvesByIds,
		)
		subShelfRoutes.POST(
			"/graphql/search",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.SearchSubShelvesOperation,
			),
			authMiddleware,
			endpoint.SearchSubShelves,
		)
	}
}
