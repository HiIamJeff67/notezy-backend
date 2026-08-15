package routers

import (
	"github.com/gin-gonic/gin"

	apicontract "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api/user-accounts"

	userservices "github.com/HiIamJeff67/notezy-backend/internal/core/services/user"
	endpoints "github.com/HiIamJeff67/notezy-backend/internal/core/transports/gateway/endpoints"
	middlewares "github.com/HiIamJeff67/notezy-backend/internal/core/transports/gateway/middlewares"
)

type UserAccountRouterDependencies struct {
	Service        userservices.UserAccountServiceInterface
	AuthMiddleware gin.HandlerFunc
}

func configureUserAccountRoutes(
	router *gin.RouterGroup,
	deps UserAccountRouterDependencies,
) {
	authMiddleware := deps.AuthMiddleware
	endpoint := endpoints.NewUserAccountEndpoint(deps.Service)
	userAccountRoutes := router.Group("/user-accounts")
	{
		userAccountRoutes.POST(
			"/get",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.GetMyAccountOperation,
			),
			authMiddleware,
			endpoint.GetMyAccount,
		)
		userAccountRoutes.POST(
			"/update",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.UpdateMyAccountOperation,
			),
			authMiddleware,
			endpoint.UpdateMyAccount,
		)
		userAccountRoutes.POST(
			"/google/bind",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.BindGoogleAccountOperation,
			),
			authMiddleware,
			endpoint.BindGoogleAccount,
		)
		userAccountRoutes.POST(
			"/google/unbind",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.UnbindGoogleAccountOperation,
			),
			authMiddleware,
			endpoint.UnbindGoogleAccount,
		)
	}
}
