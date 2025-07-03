package testroutes

import (
	"notezy-backend/app/controllers"
	"notezy-backend/app/middlewares"
	"notezy-backend/app/models/enums"
	"notezy-backend/app/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ConfigureTestAuthRoutes(db *gorm.DB, routerGroup *gin.RouterGroup) {
	if routerGroup == nil {
		routerGroup = TestRouterGroup
	}

	AuthController := controllers.NewAuthController(
		services.NewAuthService(
			db,
		),
	)

	authRoutes := routerGroup.Group("/auth")
	{
		authRoutes.POST(
			"/register",
			AuthController.Register,
		)
		authRoutes.POST(
			"/login",
			AuthController.Login,
		)
		authRoutes.POST(
			"/logout",
			middlewares.AuthMiddleware(),
			AuthController.Logout,
		)
		authRoutes.GET(
			"/sendAuthCode",
			AuthController.SendAuthCode,
		)
		authRoutes.PUT(
			"/validateEmail",
			middlewares.AuthMiddleware(),
			AuthController.ValidateEmail,
		)
		authRoutes.PUT(
			"/resetEmail",
			middlewares.AuthMiddleware(),
			middlewares.UserRoleMiddleware(enums.UserRole_Normal),
			AuthController.ResetEmail,
		)
		authRoutes.PUT(
			"/forgetPassword",
			AuthController.ForgetPassword,
		)
		authRoutes.DELETE(
			"/deleteMe",
			middlewares.AuthMiddleware(),
			AuthController.DeleteMe,
		)
	}
}
