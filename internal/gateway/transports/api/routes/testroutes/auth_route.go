package testroutes

import (
	"github.com/gin-gonic/gin"

	cookies "github.com/HiIamJeff67/notezy-backend/shared/cookies"

	binders "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/api/binders"
	controllers "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/api/controllers"
	middlewares "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/api/middlewares"
	coreadapters "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/core/adapters"
)

// the route structure is different here, since we use these routes to do the e2e test
// like it receive a database instance and a gin router group
// and its function name also start with the upper case letter
func ConfigureTestAuthRoutes(
	routerGroup *gin.RouterGroup,
	coreClient *coreadapters.CoreClient,
	accessTokenCookieHandler *cookies.CookieHandler,
	refreshTokenCookieHandler *cookies.CookieHandler,
) {
	authBinder := binders.NewAuthBinder()
	authController := controllers.NewAuthController(
		coreClient,
		accessTokenCookieHandler,
		refreshTokenCookieHandler,
	)

	authRoutes := routerGroup.Group("/auth")
	{
		authRoutes.POST(
			"/register",
			authBinder.BindRegister(authController.Register),
		)
		authRoutes.POST(
			"/registerViaGoogle",
			authBinder.BindRegisterViaGoogle(authController.RegisterViaGoogle),
		)
		authRoutes.POST(
			"/login",
			authBinder.BindLogin(authController.Login),
		)
		authRoutes.POST(
			"/loginViaGoogle",
			authBinder.BindLoginViaGoogle(authController.LoginViaGoogle),
		)
		authRoutes.POST(
			"/logout",
			middlewares.JWTMiddleware(accessTokenCookieHandler, refreshTokenCookieHandler),
			middlewares.AuthorizedRateLimitMiddleware(),
			authBinder.BindLogout(authController.Logout),
		)
		authRoutes.POST(
			"/sendAuthCode",
			authBinder.BindSendAuthCode(authController.SendAuthCode),
		)
		authRoutes.PUT(
			"/validateEmail",
			middlewares.JWTMiddleware(accessTokenCookieHandler, refreshTokenCookieHandler),
			middlewares.AuthorizedRateLimitMiddleware(),
			authBinder.BindValidateEmail(authController.ValidateEmail),
		)
		authRoutes.PUT(
			"/resetEmail",
			middlewares.JWTMiddleware(accessTokenCookieHandler, refreshTokenCookieHandler),
			middlewares.AuthorizedRateLimitMiddleware(),
			authBinder.BindResetEmail(authController.ResetEmail),
		)
		authRoutes.PUT(
			"/forgetPassword",
			authBinder.BindForgetPassword(authController.ForgetPassword),
		)
		authRoutes.PUT(
			"/resetMe",
			middlewares.JWTMiddleware(accessTokenCookieHandler, refreshTokenCookieHandler),
			middlewares.AuthorizedRateLimitMiddleware(),
			authBinder.BindResetMe(authController.ResetMe),
		)
		authRoutes.DELETE(
			"/deleteMe",
			middlewares.JWTMiddleware(accessTokenCookieHandler, refreshTokenCookieHandler),
			middlewares.AuthorizedRateLimitMiddleware(),
			authBinder.BindDeleteMe(authController.DeleteMe),
		)
	}
}
