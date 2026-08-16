package routers

import (
	"github.com/gin-gonic/gin"

	apicontract "github.com/HiIamJeff67/notegic-backend/contracts/core/v1/api/blocks"

	contexts "github.com/HiIamJeff67/notegic-backend/internal/core/contexts"
	blockservices "github.com/HiIamJeff67/notegic-backend/internal/core/services/blocks"
	endpoints "github.com/HiIamJeff67/notegic-backend/internal/core/transports/gateway/endpoints"
	middlewares "github.com/HiIamJeff67/notegic-backend/internal/core/transports/gateway/middlewares"
)

type BlockRouterDependencies struct {
	Service          blockservices.BlockServiceInterface
	AuthMiddleware   gin.HandlerFunc
	APIKeyMiddleware gin.HandlerFunc
}

func configureBlockRoutes(
	router *gin.RouterGroup,
	deps BlockRouterDependencies,
) {
	authMiddleware := deps.AuthMiddleware
	apiKeyMiddleware := deps.APIKeyMiddleware
	endpoint := endpoints.NewBlockEndpoint(deps.Service)
	apiCompatibleAuthMiddleware := middlewares.EitherMiddleware(
		[]gin.HandlerFunc{authMiddleware},
		[]gin.HandlerFunc{apiKeyMiddleware},
		func(ctx *gin.Context) bool { return contexts.IsClientGateway(ctx.Request.Context()) },
	)[0]

	blockRoutes := router.Group("/blocks")
	{
		blockRoutes.POST(
			"/get-by-id",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.GetMyBlockByIdOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.GetMyBlockById,
		)
		blockRoutes.POST(
			"/get-by-ids",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.GetMyBlocksByIdsOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.GetMyBlocksByIds,
		)
		blockRoutes.POST(
			"/get-by-block-pack-id",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.GetMyBlocksByBlockPackIdOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.GetMyBlocksByBlockPackId,
		)
		blockRoutes.POST(
			"/graphql/search",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.SearchBlocksOperation,
			),
			apiCompatibleAuthMiddleware,
			endpoint.SearchBlocks,
		)
	}
}
