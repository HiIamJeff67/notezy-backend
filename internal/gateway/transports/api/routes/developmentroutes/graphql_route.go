package developmentroutes

import (
	"time"

	"github.com/gin-gonic/gin"

	enumcontract "github.com/HiIamJeff67/notezy-backend/contracts/types/enums"
	graphql "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/api/graphql"
	interceptors "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/api/interceptors"
	middlewares "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/api/middlewares"
	coreadapters "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/core/adapters"
	cookies "github.com/HiIamJeff67/notezy-backend/shared/cookies"
)

func configureDevelopmentGraphQLRoutes(
	router *gin.RouterGroup,
	coreClient *coreadapters.CoreClient,
	accessTokenCookieHandler *cookies.CookieHandler,
	refreshTokenCookieHandler *cookies.CookieHandler,
) {
	if router == nil {
		router = DevelopmentAPIRouterGroup
	}

	graphqlRoutes := router.Group("/graphql")

	graphqlRoutes.Use(
		middlewares.UnauthorizedRateLimitMiddleware(),
		middlewares.TimeoutMiddleware(3*time.Second),
		middlewares.JWTMiddleware(accessTokenCookieHandler, refreshTokenCookieHandler),
		middlewares.AllowedPermissionsAbove(enumcontract.AccessControlPermission_Read),
		interceptors.ShareableResponseWriterInterceptor(
			interceptors.RefreshTokenInterceptor(accessTokenCookieHandler),
			interceptors.EmbeddedInterceptor,
		),
	)
	{
		graphqlRoutes.POST("/", graphql.GraphQLHandler(coreClient))
		if gin.Mode() == gin.DebugMode {
			graphqlRoutes.GET("/", graphql.PlaygroundHandler())
		}
	}
}
