package testroutes

import (
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"

	"notezy-backend/app/binders"
	controllers "notezy-backend/app/controllers"
	middlewares "notezy-backend/app/middlewares"
	enums "notezy-backend/app/models/schemas/enums"
	services "notezy-backend/app/services"
)

// the route structure is different here, since we use these routes to do the e2e test
// like it receive a database instance and a gin router group
// and its function name also start with the upper case letter
func ConfigureTestAuthRoutes(db *gorm.DB, routerGroup *gin.RouterGroup) {
	if routerGroup == nil {
		routerGroup = TestRouterGroup
	}

	authBinder := binders.NewAuthBinder()
	authController := controllers.NewAuthController(
		services.NewAuthService(
			db,
		),
	)

	authRoutes := routerGroup.Group("/auth")
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
