package testroutes

import (
	"github.com/gin-gonic/gin"

	cookies "github.com/HiIamJeff67/notezy-backend/shared/cookies"

	enumcontract "github.com/HiIamJeff67/notezy-backend/contracts/types/enums"

	graphql "github.com/HiIamJeff67/notezy-backend/internal/clientgateway/transports/api/graphql"
	middlewares "github.com/HiIamJeff67/notezy-backend/internal/clientgateway/transports/api/middlewares"
	coreadapters "github.com/HiIamJeff67/notezy-backend/internal/clientgateway/transports/core/adapters"
)

type GraphQLRouteDependencies struct {
	CoreClient                *coreadapters.CoreAdapter
	AccessTokenCookieHandler  *cookies.CookieHandler
	RefreshTokenCookieHandler *cookies.CookieHandler
}

func ConfigureTestGraphQLRoutes(
	routerGroup *gin.RouterGroup,
	deps GraphQLRouteDependencies,
) {
	coreClient, accessTokenCookieHandler, refreshTokenCookieHandler := deps.CoreClient, deps.AccessTokenCookieHandler, deps.RefreshTokenCookieHandler
	graphqlRoutes := routerGroup.Group("/graphql")

	graphqlRoutes.Use(
		middlewares.JWTMiddleware(accessTokenCookieHandler, refreshTokenCookieHandler),
		middlewares.AllowedPermissionsAbove(enumcontract.AccessControlPermission_Read),
	)
	{
		graphqlRoutes.POST("/", graphql.GraphQLHandler(coreClient))
		if gin.Mode() == gin.DebugMode {
			graphqlRoutes.GET("/", graphql.PlaygroundHandler())
		}
	}
}
