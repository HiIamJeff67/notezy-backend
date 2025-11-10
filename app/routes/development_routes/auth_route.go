package developmentroutes

import (
	"time"

	interceptors "notezy-backend/app/interceptors"
	middlewares "notezy-backend/app/middlewares"
	enums "notezy-backend/app/models/schemas/enums"
	modules "notezy-backend/app/modules"
)

func configureDevelopmentAuthRoutes() {
	authModule := modules.NewAuthModule()

	authRoutes := DevelopmentRouterGroup.Group("/auth")
	authRoutes.Use(
		middlewares.TimeoutMiddleware(3 * time.Second),
	)
	{
		authRoutes.POST(
			"/register",
			authModule.Binder.BindRegister(
				authModule.Controller.Register,
			),
		)
		authRoutes.POST(
			"/login",
			authModule.Binder.BindLogin(
				authModule.Controller.Login,
			),
		)
		authRoutes.POST(
			"/logout",
			middlewares.AuthMiddleware(),
			middlewares.AuthorizedRateLimitMiddleware(),
			authModule.Binder.BindLogout(
				authModule.Controller.Logout,
			),
		)
		authRoutes.POST(
			"/sendAuthCode",

			authModule.Binder.BindSendAuthCode(
				authModule.Controller.SendAuthCode,
			),
		)
		authRoutes.PUT(
			"/validateEmail",
			middlewares.AuthMiddleware(),
			middlewares.AuthorizedRateLimitMiddleware(),
			middlewares.CSRFMiddleware(),
			interceptors.RefreshAccessTokenInterceptor(),
			authModule.Binder.BindValidateEmail(
				authModule.Controller.ValidateEmail,
			),
		)
		authRoutes.PUT(
			"/resetEmail",
			middlewares.AuthMiddleware(),
			middlewares.AuthorizedRateLimitMiddleware(),
			middlewares.UserRoleMiddleware(enums.UserRole_Normal),
			middlewares.CSRFMiddleware(),
			interceptors.RefreshAccessTokenInterceptor(),
			authModule.Binder.BindResetEmail(
				authModule.Controller.ResetEmail,
			),
		)
		authRoutes.PUT(
			"/forgetPassword",
			authModule.Binder.BindForgetPassword(
				authModule.Controller.ForgetPassword,
			),
		)
		authRoutes.DELETE(
			"/deleteMe",
			middlewares.AuthMiddleware(),
			middlewares.AuthorizedRateLimitMiddleware(),
			middlewares.CSRFMiddleware(),
			authModule.Binder.BindDeleteMe(
				authModule.Controller.DeleteMe,
			),
		)
	}
}
