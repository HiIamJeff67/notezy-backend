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
			authBinder.BindValidateEmail(
				authController.ValidateEmail,
			),
		)
		authRoutes.PUT(
			"/resetEmail",
			middlewares.AuthMiddleware(),
			middlewares.UserRoleMiddleware(enums.UserRole_Normal),
			middlewares.AuthorizedRateLimitMiddleware(),
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
			authBinder.BindDeleteMe(
				authController.DeleteMe,
			),
		)
	}
}
