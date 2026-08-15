package routers

import (
	"github.com/gin-gonic/gin"

	apicontract "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api/user-infos"

	userservices "github.com/HiIamJeff67/notezy-backend/internal/core/services/user"
	endpoints "github.com/HiIamJeff67/notezy-backend/internal/core/transports/gateway/endpoints"
	middlewares "github.com/HiIamJeff67/notezy-backend/internal/core/transports/gateway/middlewares"
)

type UserInfoRouterDependencies struct {
	Service        userservices.UserInfoServiceInterface
	AuthMiddleware gin.HandlerFunc
}

func configureUserInfoRoutes(
	router *gin.RouterGroup,
	deps UserInfoRouterDependencies,
) {
	authMiddleware := deps.AuthMiddleware
	endpoint := endpoints.NewUserInfoEndpoint(deps.Service)
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
