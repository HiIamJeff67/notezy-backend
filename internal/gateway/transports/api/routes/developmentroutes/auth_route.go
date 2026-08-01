package developmentroutes

import (
	"time"

	"github.com/gin-gonic/gin"

	binders "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/api/binders"
	controllers "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/api/controllers"
	interceptors "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/api/interceptors"
	middlewares "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/api/middlewares"
	coreadapters "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/core/adapters"
)

func configureDevelopmentAuthRoutes(router *gin.RouterGroup, coreClient *coreadapters.CoreClient) {
	if router == nil {
		router = DevelopmentAPIRouterGroup
	}

	authBinder := binders.NewAuthBinder()
	authController := controllers.NewAuthController(coreClient)

	authRoutes := router.Group("/auth")
	{
		authRoutes.POST(
			"/register",
			middlewares.ApplyTracerMiddleware("register"),
			middlewares.ApplyMeterMiddleware("server.requests.auth.register"),
			middlewares.UnauthorizedRateLimitMiddleware(),
			middlewares.TimeoutMiddleware(5*time.Second),
			authBinder.BindRegister(authController.Register),
		)
		authRoutes.POST(
			"/registerViaGoogle",
			middlewares.ApplyTracerMiddleware("registerViaGoogle"),
			middlewares.ApplyMeterMiddleware("server.requests.auth.registerViaGoogle"),
			middlewares.UnauthorizedRateLimitMiddleware(),
			middlewares.TimeoutMiddleware(5*time.Second),
			authBinder.BindRegisterViaGoogle(authController.RegisterViaGoogle),
		)
		authRoutes.POST(
			"/login",
			middlewares.ApplyTracerMiddleware("login"),
			middlewares.ApplyMeterMiddleware("server.requests.auth.login"),
			middlewares.UnauthorizedRateLimitMiddleware(),
			middlewares.TimeoutMiddleware(3*time.Second),
			authBinder.BindLogin(authController.Login),
		)
		authRoutes.POST(
			"/loginViaGoogle",
			middlewares.ApplyTracerMiddleware("loginViaGoogle"),
			middlewares.ApplyMeterMiddleware("server.requests.auth.loginViaGoogle"),
			middlewares.UnauthorizedRateLimitMiddleware(),
			middlewares.TimeoutMiddleware(3*time.Second),
			authBinder.BindLoginViaGoogle(authController.LoginViaGoogle),
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
			authBinder.BindLogout(authController.Logout),
		)
		authRoutes.POST(
			"/sendAuthCode",
			middlewares.ApplyTracerMiddleware("sendAuthCode"),
			middlewares.ApplyMeterMiddleware("server.requests.auth.sendAuthCode"),
			middlewares.UnauthorizedRateLimitMiddleware(),
			middlewares.TimeoutMiddleware(3*time.Second),
			authBinder.BindSendAuthCode(authController.SendAuthCode),
		)
		authRoutes.PUT(
			"/validateEmail",
			middlewares.ApplyTracerMiddleware("validateEmail"),
			middlewares.ApplyMeterMiddleware("server.requests.auth.validateEmail"),
			middlewares.UnauthorizedRateLimitMiddleware(),
			middlewares.TimeoutMiddleware(3*time.Second),
			middlewares.AuthMiddleware(),
			interceptors.ShareableResponseWriterInterceptor(
				interceptors.RefreshTokenInterceptor,
				interceptors.EmbeddedInterceptor,
			),
			authBinder.BindValidateEmail(authController.ValidateEmail),
		)
		authRoutes.PUT(
			"/resetEmail",
			middlewares.ApplyTracerMiddleware("resetEmail"),
			middlewares.ApplyMeterMiddleware("server.requests.auth.resetEmail"),
			middlewares.UnauthorizedRateLimitMiddleware(),
			middlewares.TimeoutMiddleware(3*time.Second),
			middlewares.AuthMiddleware(),
			interceptors.ShareableResponseWriterInterceptor(
				interceptors.RefreshTokenInterceptor,
				interceptors.EmbeddedInterceptor,
			),
			authBinder.BindResetEmail(authController.ResetEmail),
		)
		authRoutes.PUT(
			"/forgetPassword",
			middlewares.ApplyTracerMiddleware("forgetPassword"),
			middlewares.ApplyMeterMiddleware("server.requests.auth.forgetPassword"),
			middlewares.UnauthorizedRateLimitMiddleware(),
			middlewares.TimeoutMiddleware(3*time.Second),
			authBinder.BindForgetPassword(authController.ForgetPassword),
		)
		authRoutes.PUT(
			"/resetMe",
			middlewares.ApplyTracerMiddleware("resetMe"),
			middlewares.ApplyMeterMiddleware("server.requests.auth.resetMe"),
			middlewares.UnauthorizedRateLimitMiddleware(),
			middlewares.TimeoutMiddleware(3*time.Second),
			middlewares.AuthMiddleware(),
			interceptors.ShareableResponseWriterInterceptor(
				interceptors.RefreshTokenInterceptor,
				interceptors.EmbeddedInterceptor,
			),
			authBinder.BindResetMe(authController.ResetMe),
		)
		authRoutes.DELETE(
			"/deleteMe",
			middlewares.ApplyTracerMiddleware("deleteMe"),
			middlewares.ApplyMeterMiddleware("server.requests.auth.deleteMe"),
			middlewares.UnauthorizedRateLimitMiddleware(),
			middlewares.TimeoutMiddleware(5*time.Second),
			middlewares.AuthMiddleware(),
			interceptors.ShareableResponseWriterInterceptor(
				interceptors.EmbeddedInterceptor,
			),
			authBinder.BindDeleteMe(authController.DeleteMe),
		)
	}
}
