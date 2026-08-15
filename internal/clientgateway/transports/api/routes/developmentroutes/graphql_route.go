package developmentroutes

import (
	"time"

	"github.com/gin-gonic/gin"

	cookies "github.com/HiIamJeff67/notezy-backend/shared/cookies"

	enumcontract "github.com/HiIamJeff67/notezy-backend/contracts/types/enums"

	graphql "github.com/HiIamJeff67/notezy-backend/internal/clientgateway/transports/api/graphql"
	interceptors "github.com/HiIamJeff67/notezy-backend/internal/clientgateway/transports/api/interceptors"
	middlewares "github.com/HiIamJeff67/notezy-backend/internal/clientgateway/transports/api/middlewares"
	coreadapters "github.com/HiIamJeff67/notezy-backend/internal/clientgateway/transports/core/adapters"
)

type GraphQLRouteDependencies struct {
	CoreClient                *coreadapters.CoreAdapter
	AccessTokenCookieHandler  *cookies.CookieHandler
	RefreshTokenCookieHandler *cookies.CookieHandler
	RateLimiters              RateLimiters
}

func configureDevelopmentGraphQLRoutes(
	router *gin.RouterGroup,
	deps GraphQLRouteDependencies,
) {
	coreClient, accessTokenCookieHandler, refreshTokenCookieHandler, rateLimiters := deps.CoreClient, deps.AccessTokenCookieHandler, deps.RefreshTokenCookieHandler, deps.RateLimiters
	if router == nil {
		router = DevelopmentAPIRouterGroup
	}

	graphqlRoutes := router.Group("/graphql")

	graphqlRoutes.Use(
		middlewares.UnauthorizedRateLimitMiddleware(rateLimiters.Unauthorized),
		middlewares.TimeoutMiddleware(3*time.Second),
		middlewares.GatewayAuthenticationMiddleware(accessTokenCookieHandler, refreshTokenCookieHandler),
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
