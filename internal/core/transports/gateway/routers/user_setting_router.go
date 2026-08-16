package routers

import (
	"github.com/gin-gonic/gin"

	apicontract "github.com/HiIamJeff67/notegic-backend/contracts/core/v1/api/user-settings"

	userservices "github.com/HiIamJeff67/notegic-backend/internal/core/services/user"
	endpoints "github.com/HiIamJeff67/notegic-backend/internal/core/transports/gateway/endpoints"
	middlewares "github.com/HiIamJeff67/notegic-backend/internal/core/transports/gateway/middlewares"
)

type UserSettingRouterDependencies struct {
	Service        userservices.UserSettingServiceInterface
	AuthMiddleware gin.HandlerFunc
}

func configureUserSettingRoutes(
	router *gin.RouterGroup,
	deps UserSettingRouterDependencies,
) {
	authMiddleware := deps.AuthMiddleware
	endpoint := endpoints.NewUserSettingEndpoint(deps.Service)
	userSettingRoutes := router.Group("/user-settings")
	{
		userSettingRoutes.POST(
			"/get",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.GetMySettingOperation,
			),
			authMiddleware,
			endpoint.GetMySetting,
		)
		userSettingRoutes.POST(
			"/update",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.UpdateMySettingOperation,
			),
			authMiddleware,
			endpoint.UpdateMySetting,
		)
	}
}
