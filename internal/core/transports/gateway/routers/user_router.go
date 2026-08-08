package routers

import (
	"github.com/gin-gonic/gin"

	apicontract "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api/users"

	endpoints "github.com/HiIamJeff67/notezy-backend/internal/core/transports/gateway/endpoints"
	middlewares "github.com/HiIamJeff67/notezy-backend/internal/core/transports/gateway/middlewares"
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
				apicontract.GetUserDataOperation,
			),
			authMiddleware,
			endpoint.GetUserData,
		)
		routes.POST(
			"/me",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.GetMeOperation,
			),
			authMiddleware,
			endpoint.GetMe,
		)
		routes.POST(
			"/me/update",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.UpdateMeOperation,
			),
			authMiddleware,
			endpoint.UpdateMe,
		)
		routes.POST(
			"/graphql/search",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.SearchUsersOperation,
			),
			authMiddleware,
			endpoint.SearchUsers,
		)
		routes.POST(
			"/graphql/load-theme-authors",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.LoadThemeAuthorsOperation,
			),
			authMiddleware,
			endpoint.LoadThemeAuthors,
		)
	}
}
