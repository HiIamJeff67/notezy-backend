package testroutes

import (
	"notezy-backend/app/controllers"
	"notezy-backend/app/middlewares"
	"notezy-backend/app/models/enums"

	"github.com/gin-gonic/gin"
)

func ConfigureTestAuthRoutes(routerGroup *gin.RouterGroup) {
	if routerGroup == nil {
		routerGroup = TestRouterGroup
	}

	authRoutes := routerGroup.Group("/auth")
	{
		authRoutes.POST(
			"/register",
			controllers.AuthController.Register,
		)
		authRoutes.POST(
			"/login",
			controllers.AuthController.Login,
		)
		authRoutes.POST(
			"/logout",
			middlewares.AuthMiddleware(),
			controllers.AuthController.Logout,
		)
		authRoutes.GET(
			"/sendAuthCode",
			controllers.AuthController.SendAuthCode,
		)
		authRoutes.PUT(
			"/validateEmail",
			middlewares.AuthMiddleware(),
			controllers.AuthController.ValidateEmail,
		)
		authRoutes.PUT(
			"/resetEmail",
			middlewares.AuthMiddleware(),
			middlewares.UserRoleMiddleware(enums.UserRole_Normal),
			controllers.AuthController.ResetEmail,
		)
		authRoutes.PUT(
			"/forgetPassword",
			controllers.AuthController.ForgetPassword,
		)
		authRoutes.DELETE(
			"/deleteMe",
			middlewares.AuthMiddleware(),
			controllers.AuthController.DeleteMe,
		)
	}
}
