package developmentroutes

import (
	binders "notezy-backend/app/binders"
	controllers "notezy-backend/app/controllers"
	middlewares "notezy-backend/app/middlewares"
	models "notezy-backend/app/models"
	enums "notezy-backend/app/models/schemas/enums"
	services "notezy-backend/app/services"
)

func configureDevelopmentAuthRoutes() {
	authBinder := binders.NewAuthBinder()
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
			authBinder.BindRegister(
				authController.Register,
			),
		)
		authRoutes.POST(
			"/login",
			middlewares.UnauthorizedRateLimitMiddleware(1),
			authBinder.BindLogin(
				authController.Login,
			),
		)
		authRoutes.POST(
			"/logout",
			middlewares.AuthMiddleware(),
			middlewares.RateLimitMiddleware(1),
			authBinder.BindLogout(
				authController.Logout,
			),
		)
		authRoutes.POST(
			"/sendAuthCode",
			middlewares.UnauthorizedRateLimitMiddleware(1), // may implement a block middleware to block user using this route within 1 minute
			authBinder.BindSendAuthCode(
				authController.SendAuthCode,
			),
		)
		authRoutes.PUT(
			"/validateEmail",
			middlewares.AuthMiddleware(),
			middlewares.RateLimitMiddleware(1),
			authBinder.BindValidateEmail(
				authController.ValidateEmail,
			),
		)
		authRoutes.PUT(
			"/resetEmail",
			middlewares.AuthMiddleware(),
			middlewares.UserRoleMiddleware(enums.UserRole_Normal),
			middlewares.RateLimitMiddleware(1),
			authBinder.BindResetEmail(
				authController.ResetEmail,
			),
		)
		authRoutes.PUT(
			"/forgetPassword",
			middlewares.UnauthorizedRateLimitMiddleware(1),
			authBinder.BindForgetPassword(
				authController.ForgetPassword,
			),
		)
		authRoutes.DELETE(
			"/deleteMe",
			middlewares.AuthMiddleware(),
			middlewares.RateLimitMiddleware(1),
			authBinder.BindDeleteMe(
				authController.DeleteMe,
			),
		)
	}
}
