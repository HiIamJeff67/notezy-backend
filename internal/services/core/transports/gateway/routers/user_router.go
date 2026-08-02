package routers

import (
	"github.com/gin-gonic/gin"

	usersdto "github.com/HiIamJeff67/notezy-backend/contracts/gateway/v1/api/users"
	endpoints "github.com/HiIamJeff67/notezy-backend/internal/services/core/transports/gateway/endpoints"
	middlewares "github.com/HiIamJeff67/notezy-backend/internal/services/core/transports/gateway/middlewares"
)

func configureUserRoutes(
	router *gin.RouterGroup,
	authMiddleware gin.HandlerFunc,
	endpoint endpoints.UserEndpointInterface,
) {
	routes := router.Group("/users")
	{
		routes.POST(
			"/data",
			middlewares.DelegationAuthenticatedMiddleware(
				usersdto.GetUserDataOperation,
			),
			authMiddleware,
			endpoint.GetUserData,
		)
		routes.POST(
			"/me",
			middlewares.DelegationAuthenticatedMiddleware(
				usersdto.GetMeOperation,
			),
			authMiddleware,
			endpoint.GetMe,
		)
		routes.POST(
			"/me/update",
			middlewares.DelegationAuthenticatedMiddleware(
				usersdto.UpdateMeOperation,
			),
			authMiddleware,
			endpoint.UpdateMe,
		)
		routes.POST(
			"/graphql/search",
			middlewares.DelegationAuthenticatedMiddleware(
				usersdto.SearchUsersOperation,
			),
			authMiddleware,
			endpoint.SearchUsers,
		)
		routes.POST(
			"/graphql/load-theme-authors",
			middlewares.DelegationAuthenticatedMiddleware(
				usersdto.LoadThemeAuthorsOperation,
			),
			authMiddleware,
			endpoint.LoadThemeAuthors,
		)
	}
}
