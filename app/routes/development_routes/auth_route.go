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
			middlewares.UnauthorizedRateLimitMiddleware(1),
			authController.Register,
		)
		authRoutes.POST(
			"/login",
			middlewares.UnauthorizedRateLimitMiddleware(1),
			authController.Login,
		)
		authRoutes.POST(
			"/logout",
			middlewares.AuthMiddleware(),
			middlewares.RateLimitMiddleware(1),
			authController.Logout,
		)
		authRoutes.POST(
			"/sendAuthCode",
			middlewares.UnauthorizedRateLimitMiddleware(1), // may implement a block middleware to block user using this route within 1 minute
			authController.SendAuthCode,
		)
		authRoutes.PUT(
			"/validateEmail",
			middlewares.AuthMiddleware(),
			middlewares.RateLimitMiddleware(1),
			authController.ValidateEmail,
		)
		authRoutes.PUT(
			"/resetEmail",
			middlewares.AuthMiddleware(),
			middlewares.UserRoleMiddleware(enums.UserRole_Normal),
			middlewares.RateLimitMiddleware(1),
			authController.ResetEmail,
		)
		authRoutes.PUT(
			"/forgetPassword",
			middlewares.UnauthorizedRateLimitMiddleware(1),
			authController.ForgetPassword,
		)
		authRoutes.DELETE(
			"/deleteMe",
			middlewares.AuthMiddleware(),
			middlewares.RateLimitMiddleware(1),
			authController.DeleteMe,
		)
	}
}
