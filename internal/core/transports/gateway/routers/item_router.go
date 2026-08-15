package routers

import (
	"github.com/gin-gonic/gin"

	apicontract "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api/items"

	shelfservices "github.com/HiIamJeff67/notezy-backend/internal/core/services/shelves"
	endpoints "github.com/HiIamJeff67/notezy-backend/internal/core/transports/gateway/endpoints"
	middlewares "github.com/HiIamJeff67/notezy-backend/internal/core/transports/gateway/middlewares"
)

type ItemRouterDependencies struct {
	Service        shelfservices.ItemServiceInterface
	AuthMiddleware gin.HandlerFunc
}

func configureItemRoutes(
	router *gin.RouterGroup,
	deps ItemRouterDependencies,
) {
	authMiddleware := deps.AuthMiddleware
	endpoint := endpoints.NewItemEndpoint(deps.Service)
	itemRoutes := router.Group("/items")
	{
		itemRoutes.POST(
			"/graphql/search",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.SearchItemsOperation,
			),
			authMiddleware,
			endpoint.SearchItems,
		)
	}
}
