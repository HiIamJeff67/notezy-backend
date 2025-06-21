package routes

import (
	controllers "notezy-backend/app/controllers"
	middlewares "notezy-backend/app/middlewares"
	enums "notezy-backend/app/models/enums"
)

func configureAuthRoutes() {
	authRoutes := RouterGroup.Group("/auth")
	{
		authRoutes.POST(
			"/register",
			controllers.AuthController.Register,
		)
		authRoutes.POST(
			"/login",
			controllers.AuthController.Login,
		)
		// only protected the logout route
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
