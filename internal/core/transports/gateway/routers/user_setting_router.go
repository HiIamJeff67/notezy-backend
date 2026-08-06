package routers

import (
	"github.com/gin-gonic/gin"

	usersettingsdto "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api/user-settings"

	endpoints "github.com/HiIamJeff67/notezy-backend/internal/core/transports/gateway/endpoints"
	middlewares "github.com/HiIamJeff67/notezy-backend/internal/core/transports/gateway/middlewares"
)

func configureUserSettingRoutes(
	router *gin.RouterGroup,
	authMiddleware gin.HandlerFunc,
	endpoint endpoints.UserSettingEndpointInterface,
) {
	userSettingRoutes := router.Group("/user-settings")
	{
		userSettingRoutes.POST(
			"/get",
			middlewares.DelegationAuthenticatedMiddleware(
				usersettingsdto.GetMySettingOperation,
			),
			authMiddleware,
			endpoint.GetMySetting,
		)
		userSettingRoutes.POST(
			"/update",
			middlewares.DelegationAuthenticatedMiddleware(
				usersettingsdto.UpdateMySettingOperation,
			),
			authMiddleware,
			endpoint.UpdateMySetting,
		)
	}
}
