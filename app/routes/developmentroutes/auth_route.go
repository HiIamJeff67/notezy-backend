package developmentroutes

import (
	"time"

	"github.com/gin-gonic/gin"

	interceptors "github.com/HiIamJeff67/notezy-backend/app/interceptors"
	middlewares "github.com/HiIamJeff67/notezy-backend/app/middlewares"
	enums "github.com/HiIamJeff67/notezy-backend/app/models/schemas/enums"
	modules "github.com/HiIamJeff67/notezy-backend/app/modules"
)

func configureDevelopmentAuthRoutes(router *gin.RouterGroup) {
	if router == nil {
		router = DevelopmentAPIRouterGroup
	}

	authModule := modules.NewAuthModule()

	authRoutes := router.Group("/auth")
	{
		authRoutes.POST(
			"/register",
			middlewares.ApplyTracerMiddleware("register"),
			middlewares.ApplyMeterMiddleware("server.requests.auth.register"),
			middlewares.UnauthorizedRateLimitMiddleware(),
			middlewares.TimeoutMiddleware(5*time.Second),
			authModule.Binder.BindRegister(
				authModule.Controller.Register,
			),
		)
		authRoutes.POST(
			"/registerViaGoogle",
			middlewares.ApplyTracerMiddleware("registerViaGoogle"),
			middlewares.ApplyMeterMiddleware("server.requests.auth.registerViaGoogle"),
			middlewares.UnauthorizedRateLimitMiddleware(),
			middlewares.TimeoutMiddleware(5*time.Second),
			authModule.Binder.BindRegisterViaGoogle(
				authModule.Controller.RegisterViaGoogle,
			),
		)
		authRoutes.POST(
			"/login",
			middlewares.ApplyTracerMiddleware("login"),
			middlewares.ApplyMeterMiddleware("server.requests.auth.login"),
			middlewares.UnauthorizedRateLimitMiddleware(),
			middlewares.TimeoutMiddleware(3*time.Second),
			authModule.Binder.BindLogin(
				authModule.Controller.Login,
			),
		)
		authRoutes.POST(
			"/loginViaGoogle",
			middlewares.ApplyTracerMiddleware("loginViaGoogle"),
			middlewares.ApplyMeterMiddleware("server.requests.auth.loginViaGoogle"),
			middlewares.UnauthorizedRateLimitMiddleware(),
			middlewares.TimeoutMiddleware(3*time.Second),
			authModule.Binder.BindLoginViaGoogle(
				authModule.Controller.LoginViaGoogle,
			),
		)
		authRoutes.POST(
			"/logout",
			middlewares.ApplyTracerMiddleware("logout"),
			middlewares.ApplyMeterMiddleware("server.requests.auth.logout"),
			middlewares.UnauthorizedRateLimitMiddleware(),
			middlewares.TimeoutMiddleware(3*time.Second),
			middlewares.AuthMiddleware(),
			interceptors.ShareableResponseWriterInterceptor(
				interceptors.EmbeddedInterceptor,
			),
			authModule.Binder.BindLogout(
				authModule.Controller.Logout,
			),
		)
		authRoutes.POST(
			"/sendAuthCode",
			middlewares.ApplyTracerMiddleware("sendAuthCode"),
			middlewares.ApplyMeterMiddleware("server.requests.auth.sendAuthCode"),
			middlewares.UnauthorizedRateLimitMiddleware(),
			middlewares.TimeoutMiddleware(3*time.Second),
			authModule.Binder.BindSendAuthCode(
				authModule.Controller.SendAuthCode,
			),
		)
		authRoutes.PUT(
			"/validateEmail",
			middlewares.ApplyTracerMiddleware("validateEmail"),
			middlewares.ApplyMeterMiddleware("server.requests.auth.validateEmail"),
			middlewares.UnauthorizedRateLimitMiddleware(),
			middlewares.TimeoutMiddleware(3*time.Second),
			middlewares.AuthMiddleware(),
			middlewares.CSRFMiddleware(),
			interceptors.ShareableResponseWriterInterceptor(
				interceptors.RefreshTokenInterceptor,
				interceptors.EmbeddedInterceptor,
			),
			authModule.Binder.BindValidateEmail(
				authModule.Controller.ValidateEmail,
			),
		)
		authRoutes.PUT(
			"/resetEmail",
			middlewares.ApplyTracerMiddleware("resetEmail"),
			middlewares.ApplyMeterMiddleware("server.requests.auth.resetEmail"),
			middlewares.UnauthorizedRateLimitMiddleware(),
			middlewares.TimeoutMiddleware(3*time.Second),
			middlewares.AuthMiddleware(),
			middlewares.UserRoleMiddleware(enums.UserRole_Normal),
			middlewares.CSRFMiddleware(),
			interceptors.ShareableResponseWriterInterceptor(
				interceptors.RefreshTokenInterceptor,
				interceptors.EmbeddedInterceptor,
			),
			authModule.Binder.BindResetEmail(
				authModule.Controller.ResetEmail,
			),
		)
		authRoutes.PUT(
			"/forgetPassword",
			middlewares.ApplyTracerMiddleware("forgetPassword"),
			middlewares.ApplyMeterMiddleware("server.requests.auth.forgetPassword"),
			middlewares.UnauthorizedRateLimitMiddleware(),
			middlewares.TimeoutMiddleware(3*time.Second),
			authModule.Binder.BindForgetPassword(
				authModule.Controller.ForgetPassword,
			),
		)
		authRoutes.PUT(
			"/resetMe",
			middlewares.ApplyTracerMiddleware("resetMe"),
			middlewares.ApplyMeterMiddleware("server.requests.auth.resetMe"),
			middlewares.UnauthorizedRateLimitMiddleware(),
			middlewares.TimeoutMiddleware(3*time.Second),
			middlewares.AuthMiddleware(),
			middlewares.CSRFMiddleware(),
			interceptors.ShareableResponseWriterInterceptor(
				interceptors.RefreshTokenInterceptor,
				interceptors.EmbeddedInterceptor,
			),
			authModule.Binder.BindResetMe(
				authModule.Controller.ResetMe,
			),
		)
		authRoutes.DELETE(
			"/deleteMe",
			middlewares.ApplyTracerMiddleware("deleteMe"),
			middlewares.ApplyMeterMiddleware("server.requests.auth.deleteMe"),
			middlewares.UnauthorizedRateLimitMiddleware(),
			middlewares.TimeoutMiddleware(5*time.Second),
			middlewares.AuthMiddleware(),
			middlewares.CSRFMiddleware(),
			interceptors.ShareableResponseWriterInterceptor(
				interceptors.EmbeddedInterceptor,
			),
			authModule.Binder.BindDeleteMe(
				authModule.Controller.DeleteMe,
			),
		)
	}
}
