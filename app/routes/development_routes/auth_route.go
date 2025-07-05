package developmentroutes

import (
	controllers "notezy-backend/app/controllers"
	middlewares "notezy-backend/app/middlewares"
	models "notezy-backend/app/models"
	enums "notezy-backend/app/models/schemas/enums"
	services "notezy-backend/app/services"
)

func configureDevelopmentAuthRoutes() {
	authController := controllers.NewAuthController(
		services.NewAuthService(
			models.NotezyDB,
		),
	)

	authRoutes := DevelopmentRouterGroup.Group("/auth")
	{
		authRoutes.POST(
			"/register",
			authController.Register,
		)
		authRoutes.POST(
			"/login",
			authController.Login,
		)
		authRoutes.POST(
			"/logout",
			middlewares.AuthMiddleware(),
			authController.Logout,
		)
		authRoutes.GET(
			"/sendAuthCode",
			authController.SendAuthCode,
		)
		authRoutes.PUT(
			"/validateEmail",
			middlewares.AuthMiddleware(),
			authController.ValidateEmail,
		)
		authRoutes.PUT(
			"/resetEmail",
			middlewares.AuthMiddleware(),
			middlewares.UserRoleMiddleware(enums.UserRole_Normal),
			authController.ResetEmail,
		)
		authRoutes.PUT(
			"/forgetPassword",
			authController.ForgetPassword,
		)
		authRoutes.DELETE(
			"/deleteMe",
			middlewares.AuthMiddleware(),
			authController.DeleteMe,
		)
	}
}
