package testroutes

import (
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"

	middlewares "notezy-backend/app/middlewares"
	enums "notezy-backend/app/models/schemas/enums"
	modules "notezy-backend/app/modules"
)

// the route structure is different here, since we use these routes to do the e2e test
// like it receive a database instance and a gin router group
// and its function name also start with the upper case letter
func ConfigureTestAuthRoutes(db *gorm.DB, routerGroup *gin.RouterGroup) {
	if routerGroup == nil {
		routerGroup = TestRouterGroup
	}

	authModule := modules.NewAuthModule()

	authRoutes := routerGroup.Group("/auth")
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
			authModule.Binder.BindValidateEmail(
				authModule.Controller.ValidateEmail,
			),
		)
		authRoutes.PUT(
			"/resetEmail",
			middlewares.AuthMiddleware(),
			middlewares.UserRoleMiddleware(enums.UserRole_Normal),
			middlewares.AuthorizedRateLimitMiddleware(),
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
			authModule.Binder.BindDeleteMe(
				authModule.Controller.DeleteMe,
			),
		)
	}
}
