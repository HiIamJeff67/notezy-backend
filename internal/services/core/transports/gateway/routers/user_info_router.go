package routers

import (
	"github.com/gin-gonic/gin"

	userinfosdto "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api/user-infos"
	endpoints "github.com/HiIamJeff67/notezy-backend/internal/services/core/transports/gateway/endpoints"
	middlewares "github.com/HiIamJeff67/notezy-backend/internal/services/core/transports/gateway/middlewares"
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
				userinfosdto.GetMyInfoOperation,
			),
			authMiddleware,
			endpoint.GetMyInfo,
		)
		userInfoRoutes.POST(
			"/update",
			middlewares.DelegationAuthenticatedMiddleware(
				userinfosdto.UpdateMyInfoOperation,
			),
			authMiddleware,
			endpoint.UpdateMyInfo,
		)
		userInfoRoutes.POST(
			"/graphql/load",
			middlewares.DelegationAuthenticatedMiddleware(
				userinfosdto.LoadUserInfosOperation,
			),
			authMiddleware,
			endpoint.LoadUserInfos,
		)
	}
}
