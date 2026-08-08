package routers

import (
	"github.com/gin-gonic/gin"

	apicontract "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api/user-infos"

	endpoints "github.com/HiIamJeff67/notezy-backend/internal/core/transports/gateway/endpoints"
	middlewares "github.com/HiIamJeff67/notezy-backend/internal/core/transports/gateway/middlewares"
)

func configureUserInfoRoutes(
	router *gin.RouterGroup,
	authMiddleware gin.HandlerFunc,
	endpoint endpoints.UserInfoEndpointInterface,
) {
	userInfoRoutes := router.Group("/user-infos")
	{
		userInfoRoutes.POST(
			"/get",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.GetMyInfoOperation,
			),
			authMiddleware,
			endpoint.GetMyInfo,
		)
		userInfoRoutes.POST(
			"/update",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.UpdateMyInfoOperation,
			),
			authMiddleware,
			endpoint.UpdateMyInfo,
		)
		userInfoRoutes.POST(
			"/graphql/load",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.LoadUserInfosOperation,
			),
			authMiddleware,
			endpoint.LoadUserInfos,
		)
	}
}
