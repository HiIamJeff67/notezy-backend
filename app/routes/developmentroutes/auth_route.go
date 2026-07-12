package developmentroutes

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"

	interceptors "github.com/HiIamJeff67/notezy-backend/app/interceptors"
	middlewares "github.com/HiIamJeff67/notezy-backend/app/middlewares"
	enums "github.com/HiIamJeff67/notezy-backend/app/models/schemas/enums"
	modules "github.com/HiIamJeff67/notezy-backend/app/modules"
	metrics "github.com/HiIamJeff67/notezy-backend/app/monitor/metrics"
	constants "github.com/HiIamJeff67/notezy-backend/shared/constants"
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
			middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "register"),
			middlewares.ApplyMeterMiddleware(
				otel.Meter(constants.ServiceName),
				metrics.MetricNames.Server.Requests.Auth.Register,
			),
			middlewares.UnauthorizedRateLimitMiddleware(),
			middlewares.TimeoutMiddleware(5*time.Second),
			authModule.Binder.BindRegister(
				authModule.Controller.Register,
			),
		)
		authRoutes.POST(
			"/registerViaGoogle",
			middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "registerViaGoogle"),
			middlewares.ApplyMeterMiddleware(
				otel.Meter(constants.ServiceName),
				metrics.MetricNames.Server.Requests.Auth.RegisterViaGoogle,
			),
			middlewares.UnauthorizedRateLimitMiddleware(),
			middlewares.TimeoutMiddleware(5*time.Second),
			authModule.Binder.BindRegisterViaGoogle(
				authModule.Controller.RegisterViaGoogle,
			),
		)
		authRoutes.POST(
			"/login",
			middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "login"),
			middlewares.ApplyMeterMiddleware(
				otel.Meter(constants.ServiceName),
				metrics.MetricNames.Server.Requests.Auth.Login,
			),
			middlewares.UnauthorizedRateLimitMiddleware(),
			middlewares.TimeoutMiddleware(3*time.Second),
			authModule.Binder.BindLogin(
				authModule.Controller.Login,
			),
		)
		authRoutes.POST(
			"/loginViaGoogle",
			middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "loginViaGoogle"),
			middlewares.ApplyMeterMiddleware(
				otel.Meter(constants.ServiceName),
				metrics.MetricNames.Server.Requests.Auth.LoginViaGoogle,
			),
			middlewares.UnauthorizedRateLimitMiddleware(),
			middlewares.TimeoutMiddleware(3*time.Second),
			authModule.Binder.BindLoginViaGoogle(
				authModule.Controller.LoginViaGoogle,
			),
		)
		authRoutes.POST(
			"/logout",
			middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "logout"),
			middlewares.ApplyMeterMiddleware(
				otel.Meter(constants.ServiceName),
				metrics.MetricNames.Server.Requests.Auth.Logout,
			),
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
			middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "sendAuthCode"),
			middlewares.ApplyMeterMiddleware(
				otel.Meter(constants.ServiceName),
				metrics.MetricNames.Server.Requests.Auth.SendAuthCode,
			),
			middlewares.UnauthorizedRateLimitMiddleware(),
			middlewares.TimeoutMiddleware(3*time.Second),
			authModule.Binder.BindSendAuthCode(
				authModule.Controller.SendAuthCode,
			),
		)
		authRoutes.PUT(
			"/validateEmail",
			middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "validateEmail"),
			middlewares.ApplyMeterMiddleware(
				otel.Meter(constants.ServiceName),
				metrics.MetricNames.Server.Requests.Auth.ValidateEmail,
			),
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
			middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "resetEmail"),
			middlewares.ApplyMeterMiddleware(
				otel.Meter(constants.ServiceName),
				metrics.MetricNames.Server.Requests.Auth.ResetEmail,
			),
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
			middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "forgetPassword"),
			middlewares.ApplyMeterMiddleware(
				otel.Meter(constants.ServiceName),
				metrics.MetricNames.Server.Requests.Auth.ForgetPassword,
			),
			middlewares.UnauthorizedRateLimitMiddleware(),
			middlewares.TimeoutMiddleware(3*time.Second),
			authModule.Binder.BindForgetPassword(
				authModule.Controller.ForgetPassword,
			),
		)
		authRoutes.PUT(
			"/resetMe",
			middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "resetMe"),
			middlewares.ApplyMeterMiddleware(
				otel.Meter(constants.ServiceName),
				metrics.MetricNames.Server.Requests.Auth.ResetMe,
			),
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
			middlewares.ApplyTracerMiddleware(otel.Tracer(constants.ServiceName), "deleteMe"),
			middlewares.ApplyMeterMiddleware(
				otel.Meter(constants.ServiceName),
				metrics.MetricNames.Server.Requests.Auth.DeleteMe,
			),
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
