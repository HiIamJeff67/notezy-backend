package routers

import (
	"github.com/gin-gonic/gin"

	useraccountsdto "github.com/HiIamJeff67/notezy-backend/contracts/api/v1/user-accounts"
	endpoints "github.com/HiIamJeff67/notezy-backend/internal/services/core/transports/gateway/endpoints"
	middlewares "github.com/HiIamJeff67/notezy-backend/internal/services/core/transports/gateway/middlewares"
)

func configureUserAccountRoutes(
	router *gin.RouterGroup,
	authMiddleware gin.HandlerFunc,
	endpoint endpoints.UserAccountEndpointInterface,
) {
	userAccountRoutes := router.Group("/user-accounts")
	{
		userAccountRoutes.POST(
			"/get",
			middlewares.DelegationAuthenticatedMiddleware(
				useraccountsdto.GetMyAccountOperation,
			),
			authMiddleware,
			endpoint.GetMyAccount,
		)
		userAccountRoutes.POST(
			"/update",
			middlewares.DelegationAuthenticatedMiddleware(
				useraccountsdto.UpdateMyAccountOperation,
			),
			authMiddleware,
			endpoint.UpdateMyAccount,
		)
		userAccountRoutes.POST(
			"/google/bind",
			middlewares.DelegationAuthenticatedMiddleware(
				useraccountsdto.BindGoogleAccountOperation,
			),
			authMiddleware,
			endpoint.BindGoogleAccount,
		)
		userAccountRoutes.POST(
			"/google/unbind",
			middlewares.DelegationAuthenticatedMiddleware(
				useraccountsdto.UnbindGoogleAccountOperation,
			),
			authMiddleware,
			endpoint.UnbindGoogleAccount,
		)
	}
}
