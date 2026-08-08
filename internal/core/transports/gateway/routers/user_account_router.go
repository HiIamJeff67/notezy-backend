package routers

import (
	"github.com/gin-gonic/gin"

	apicontract "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api/user-accounts"

	endpoints "github.com/HiIamJeff67/notezy-backend/internal/core/transports/gateway/endpoints"
	middlewares "github.com/HiIamJeff67/notezy-backend/internal/core/transports/gateway/middlewares"
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
