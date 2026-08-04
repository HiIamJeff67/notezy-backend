package routers

import (
	"github.com/gin-gonic/gin"

	itemsdto "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api/items"
	endpoints "github.com/HiIamJeff67/notezy-backend/internal/services/core/transports/gateway/endpoints"
	middlewares "github.com/HiIamJeff67/notezy-backend/internal/services/core/transports/gateway/middlewares"
)

func configureItemRoutes(
	router *gin.RouterGroup,
	authMiddleware gin.HandlerFunc,
	endpoint endpoints.ItemEndpointInterface,
) {
	itemRoutes := router.Group("/items")
	{
		itemRoutes.POST(
			"/graphql/search",
			middlewares.DelegationAuthenticatedMiddleware(
				itemsdto.SearchItemsOperation,
			),
			authMiddleware,
			endpoint.SearchItems,
		)
	}
}
