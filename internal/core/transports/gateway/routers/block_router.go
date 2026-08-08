package routers

import (
	"github.com/gin-gonic/gin"

	apicontract "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api/blocks"

	endpoints "github.com/HiIamJeff67/notezy-backend/internal/core/transports/gateway/endpoints"
	middlewares "github.com/HiIamJeff67/notezy-backend/internal/core/transports/gateway/middlewares"
)

func configureBlockRoutes(
	router *gin.RouterGroup,
	authMiddleware gin.HandlerFunc,
	endpoint endpoints.BlockEndpointInterface,
) {
	blockRoutes := router.Group("/blocks")
	{
		blockRoutes.POST(
			"/get-by-id",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.GetMyBlockByIdOperation,
			),
			authMiddleware,
			endpoint.GetMyBlockById,
		)
		blockRoutes.POST(
			"/get-by-ids",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.GetMyBlocksByIdsOperation,
			),
			authMiddleware,
			endpoint.GetMyBlocksByIds,
		)
		blockRoutes.POST(
			"/get-by-block-pack-id",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.GetMyBlocksByBlockPackIdOperation,
			),
			authMiddleware,
			endpoint.GetMyBlocksByBlockPackId,
		)
		blockRoutes.POST(
			"/graphql/search",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.SearchBlocksOperation,
			),
			authMiddleware,
			endpoint.SearchBlocks,
		)
	}
}
