package routers

import (
	"github.com/gin-gonic/gin"

	apicontract "github.com/HiIamJeff67/notegic-backend/contracts/core/v1/api/sub-shelves"

	contexts "github.com/HiIamJeff67/notegic-backend/internal/core/contexts"
	shelfservices "github.com/HiIamJeff67/notegic-backend/internal/core/services/shelves"
	endpoints "github.com/HiIamJeff67/notegic-backend/internal/core/transports/gateway/endpoints"
	middlewares "github.com/HiIamJeff67/notegic-backend/internal/core/transports/gateway/middlewares"
)

type SubShelfRouterDependencies struct {
	Service          shelfservices.SubShelfServiceInterface
	AuthMiddleware   gin.HandlerFunc
	APIKeyMiddleware gin.HandlerFunc
}

func configureSubShelfRoutes(
	router *gin.RouterGroup,
	deps SubShelfRouterDependencies,
) {
	authMiddleware := deps.AuthMiddleware
	apiKeyMiddleware := deps.APIKeyMiddleware
	endpoint := endpoints.NewSubShelfEndpoint(deps.Service)
	apiCompatibleAuthMiddleware := middlewares.EitherMiddleware(
		[]gin.HandlerFunc{authMiddleware},
		[]gin.HandlerFunc{apiKeyMiddleware},
		func(ctx *gin.Context) bool { return contexts.IsClientGateway(ctx.Request.Context()) },
	)[0]

	subShelfRoutes := router.Group("/sub-shelves")
	{
		subShelfRoutes.POST(
			"/get-by-id",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.GetMySubShelfByIdOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.GetMySubShelfById,
		)
		subShelfRoutes.POST(
			"/get-by-prev-sub-shelf-id",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.GetMySubShelvesByPrevSubShelfIdOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.GetMySubShelvesByPrevSubShelfId,
		)
		subShelfRoutes.POST(
			"/get-all-by-root-shelf-id",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.GetAllMySubShelvesByRootShelfIdOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.GetAllMySubShelvesByRootShelfId,
		)
		subShelfRoutes.POST(
			"/get-and-items-by-prev-sub-shelf-id",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.GetMySubShelvesAndItemsByPrevSubShelfIdOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.GetMySubShelvesAndItemsByPrevSubShelfId,
		)
		subShelfRoutes.POST(
			"/create",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.CreateSubShelfByRootShelfIdOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.CreateSubShelfByRootShelfId,
		)
		subShelfRoutes.POST(
			"/create-many",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.CreateSubShelvesByRootShelfIdsOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.CreateSubShelvesByRootShelfIds,
		)
		subShelfRoutes.POST(
			"/update",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.UpdateMySubShelfByIdOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.UpdateMySubShelfById,
		)
		subShelfRoutes.POST(
			"/update-many",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.UpdateMySubShelvesByIdsOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.UpdateMySubShelvesByIds,
		)
		subShelfRoutes.POST(
			"/move",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.MoveMySubShelfByRootShelfIdOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.MoveMySubShelfByRootShelfId,
		)
		subShelfRoutes.POST(
			"/move-many",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.MoveMySubShelvesByRootShelfIdOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.MoveMySubShelvesByRootShelfId,
		)
		subShelfRoutes.POST(
			"/move-many-by-root-shelves",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.MoveMySubShelvesByRootShelfIdsOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.MoveMySubShelvesByRootShelfIds,
		)
		subShelfRoutes.POST(
			"/restore",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.RestoreMySubShelfByIdOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.RestoreMySubShelfById,
		)
		subShelfRoutes.POST(
			"/restore-many",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.RestoreMySubShelvesByIdsOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.RestoreMySubShelvesByIds,
		)
		subShelfRoutes.POST(
			"/delete",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.DeleteMySubShelfByIdOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.DeleteMySubShelfById,
		)
		subShelfRoutes.POST(
			"/delete-many",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.DeleteMySubShelvesByIdsOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.DeleteMySubShelvesByIds,
		)
		subShelfRoutes.POST(
			"/graphql/search",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.SearchSubShelvesOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.SearchSubShelves,
		)
	}
}
