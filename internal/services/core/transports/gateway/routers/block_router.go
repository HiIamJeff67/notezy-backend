package routers

import (
	"github.com/gin-gonic/gin"

	blocksdto "github.com/HiIamJeff67/notezy-backend/contracts/gateway/v1/api/blocks"
	endpoints "github.com/HiIamJeff67/notezy-backend/internal/services/core/transports/gateway/endpoints"
	middlewares "github.com/HiIamJeff67/notezy-backend/internal/services/core/transports/gateway/middlewares"
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
				blocksdto.GetMyBlockByIdOperation,
			),
			authMiddleware,
			endpoint.GetMyBlockById,
		)
		blockRoutes.POST(
			"/get-by-ids",
			middlewares.DelegationAuthenticatedMiddleware(
				blocksdto.GetMyBlocksByIdsOperation,
			),
			authMiddleware,
			endpoint.GetMyBlocksByIds,
		)
		blockRoutes.POST(
			"/get-by-block-pack-id",
			middlewares.DelegationAuthenticatedMiddleware(
				blocksdto.GetMyBlocksByBlockPackIdOperation,
			),
			authMiddleware,
			endpoint.GetMyBlocksByBlockPackId,
		)
		blockRoutes.POST(
			"/graphql/search",
			middlewares.DelegationAuthenticatedMiddleware(
				blocksdto.SearchBlocksOperation,
			),
			authMiddleware,
			endpoint.SearchBlocks,
		)
	}
}
