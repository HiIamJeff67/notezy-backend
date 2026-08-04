package testroutes

import (
	enumcontract "github.com/HiIamJeff67/notezy-backend/contracts/types/enums"
	graphql "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/api/graphql"
	middlewares "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/api/middlewares"
	coreadapters "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/core/adapters"
	cookies "github.com/HiIamJeff67/notezy-backend/shared/cookies"

	"github.com/gin-gonic/gin"
)

func ConfigureTestGraphQLRoutes(
	routerGroup *gin.RouterGroup,
	coreClient *coreadapters.CoreClient,
	accessTokenCookieHandler *cookies.CookieHandler,
	refreshTokenCookieHandler *cookies.CookieHandler,
) {
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
