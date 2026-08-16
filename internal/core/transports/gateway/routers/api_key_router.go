package routers

import (
	"github.com/gin-gonic/gin"

	apicontract "github.com/HiIamJeff67/notegic-backend/contracts/core/v1/api/api-keys"

	apikeyservices "github.com/HiIamJeff67/notegic-backend/internal/core/services/apikey"
	endpoints "github.com/HiIamJeff67/notegic-backend/internal/core/transports/gateway/endpoints"
	middlewares "github.com/HiIamJeff67/notegic-backend/internal/core/transports/gateway/middlewares"
)

type APIKeyRouterDependencies struct {
	Service        apikeyservices.APIKeyServiceInterface
	AuthMiddleware gin.HandlerFunc
}

func configureAPIKeyRoutes(
	router *gin.RouterGroup,
	deps APIKeyRouterDependencies,
) {
	authMiddleware := deps.AuthMiddleware
	endpoint := endpoints.NewAPIKeyEndpoint(deps.Service)
	routes := router.Group("/api-keys")
	{
		routes.POST(
			"/create",
			middlewares.DelegationAuthenticatedMiddleware(apicontract.CreateMyAPIKeyOperation),
			authMiddleware,
			endpoint.CreateMyAPIKey,
		)
		routes.POST(
			"/list",
			middlewares.DelegationAuthenticatedMiddleware(apicontract.ListMyAPIKeysOperation),
			authMiddleware,
			endpoint.ListMyAPIKeys,
		)
		routes.POST(
			"/revoke",
			middlewares.DelegationAuthenticatedMiddleware(apicontract.RevokeMyAPIKeyOperation),
			authMiddleware,
			endpoint.RevokeMyAPIKey,
		)
	}
}
