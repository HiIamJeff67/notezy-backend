package routers

import (
	"github.com/gin-gonic/gin"

	apicontract "github.com/HiIamJeff67/notegic-backend/contracts/core/v1/api/block-packs"

	contexts "github.com/HiIamJeff67/notegic-backend/internal/core/contexts"
	blockservices "github.com/HiIamJeff67/notegic-backend/internal/core/services/blocks"
	endpoints "github.com/HiIamJeff67/notegic-backend/internal/core/transports/gateway/endpoints"
	middlewares "github.com/HiIamJeff67/notegic-backend/internal/core/transports/gateway/middlewares"
)

type BlockPackRouterDependencies struct {
	Service          blockservices.BlockPackServiceInterface
	AuthMiddleware   gin.HandlerFunc
	APIKeyMiddleware gin.HandlerFunc
}

func configureBlockPackRoutes(
	router *gin.RouterGroup,
	deps BlockPackRouterDependencies,
) {
	authMiddleware := deps.AuthMiddleware
	apiKeyMiddleware := deps.APIKeyMiddleware
	endpoint := endpoints.NewBlockPackEndpoint(deps.Service)
	apiCompatibleAuthMiddleware := middlewares.EitherMiddleware(
		[]gin.HandlerFunc{authMiddleware},
		[]gin.HandlerFunc{apiKeyMiddleware},
		func(ctx *gin.Context) bool { return contexts.IsClientGateway(ctx.Request.Context()) },
	)[0]

	blockPackRoutes := router.Group("/block-packs")
	{
		blockPackRoutes.POST(
			"/get-by-id",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.GetMyBlockPackByIdOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.GetMyBlockPackById,
		)
		blockPackRoutes.POST(
			"/get-and-parent-by-id",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.GetMyBlockPackAndItsParentByIdOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.GetMyBlockPackAndItsParentById,
		)
		blockPackRoutes.POST(
			"/get-by-parent-sub-shelf-id",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.GetMyBlockPacksByParentSubShelfIdOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.GetMyBlockPacksByParentSubShelfId,
		)
		blockPackRoutes.POST(
			"/get-all-by-root-shelf-id",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.GetAllMyBlockPacksByRootShelfIdOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.GetAllMyBlockPacksByRootShelfId,
		)
		blockPackRoutes.POST(
			"/create",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.CreateBlockPackOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.CreateBlockPack,
		)
		blockPackRoutes.POST(
			"/create-many",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.CreateBlockPacksOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.CreateBlockPacks,
		)
		blockPackRoutes.POST(
			"/update",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.UpdateMyBlockPackByIdOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.UpdateMyBlockPackById,
		)
		blockPackRoutes.POST(
			"/update-many",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.UpdateMyBlockPacksByIdsOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.UpdateMyBlockPacksByIds,
		)
		blockPackRoutes.POST(
			"/move",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.MoveMyBlockPackByParentSubShelfIdOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.MoveMyBlockPackByParentSubShelfId,
		)
		blockPackRoutes.POST(
			"/move-many",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.MoveMyBlockPacksByParentSubShelfIdOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.MoveMyBlockPacksByParentSubShelfId,
		)
		blockPackRoutes.POST(
			"/move-many-by-parent-sub-shelves",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.MoveMyBlockPacksByParentSubShelfIdsOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.MoveMyBlockPacksByParentSubShelfIds,
		)
		blockPackRoutes.POST(
			"/restore",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.RestoreMyBlockPackByIdOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.RestoreMyBlockPackById,
		)
		blockPackRoutes.POST(
			"/restore-many",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.RestoreMyBlockPacksByIdsOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.RestoreMyBlockPacksByIds,
		)
		blockPackRoutes.POST(
			"/delete",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.DeleteMyBlockPackByIdOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.DeleteMyBlockPackById,
		)
		blockPackRoutes.POST(
			"/delete-many",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.DeleteMyBlockPacksByIdsOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.DeleteMyBlockPacksByIds,
		)
		blockPackRoutes.POST(
			"/graphql/search",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.SearchBlockPacksOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.SearchBlockPacks,
		)
	}
}
