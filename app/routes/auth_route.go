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
			controllers.Register,
		)
		authRoutes.POST(
			"/login",
			controllers.Login,
		)
		// only protected the logout route
		authRoutes.POST(
			"/logout",
			middlewares.AuthMiddleware(),
			controllers.Logout,
		)
		authRoutes.GET(
			"/sendAuthCode",
			controllers.SendAuthCode,
		)
		authRoutes.PUT(
			"/validateEmail",
			middlewares.AuthMiddleware(),
			controllers.ValidateEmail,
		)
		authRoutes.PUT(
			"/resetEmail",
			middlewares.AuthMiddleware(),
			middlewares.UserRoleMiddleware(enums.UserRole_Normal),
			controllers.ResetEmail,
		)
		authRoutes.PUT(
			"/forgetPassword",
			controllers.ForgetPassword,
		)
	}
}
