package developmentroutes

import (
	binders "notezy-backend/app/binders"
	controllers "notezy-backend/app/controllers"
	interceptors "notezy-backend/app/interceptors"
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
			authBinder.BindRegister(
				authController.Register,
			),
		)
		authRoutes.POST(
			"/login",
			authBinder.BindLogin(
				authController.Login,
			),
		)
		authRoutes.POST(
			"/logout",
			middlewares.AuthMiddleware(),
			middlewares.AuthorizedRateLimitMiddleware(),
			authBinder.BindLogout(
				authController.Logout,
			),
		)
		authRoutes.POST(
			"/sendAuthCode",

			authBinder.BindSendAuthCode(
				authController.SendAuthCode,
			),
		)
		authRoutes.PUT(
			"/validateEmail",
			middlewares.AuthMiddleware(),
			middlewares.AuthorizedRateLimitMiddleware(),
			middlewares.CSRFMiddleware(),
			interceptors.RefreshAccessTokenInterceptor(),
			authBinder.BindValidateEmail(
				authController.ValidateEmail,
			),
		)
		authRoutes.PUT(
			"/resetEmail",
			middlewares.AuthMiddleware(),
			middlewares.AuthorizedRateLimitMiddleware(),
			middlewares.UserRoleMiddleware(enums.UserRole_Normal),
			middlewares.CSRFMiddleware(),
			interceptors.RefreshAccessTokenInterceptor(),
			authBinder.BindResetEmail(
				authController.ResetEmail,
			),
		)
		authRoutes.PUT(
			"/forgetPassword",
			authBinder.BindForgetPassword(
				authController.ForgetPassword,
			),
		)
		authRoutes.DELETE(
			"/deleteMe",
			middlewares.AuthMiddleware(),
			middlewares.AuthorizedRateLimitMiddleware(),
			middlewares.CSRFMiddleware(),
			authBinder.BindDeleteMe(
				authController.DeleteMe,
			),
		)
	}
}
