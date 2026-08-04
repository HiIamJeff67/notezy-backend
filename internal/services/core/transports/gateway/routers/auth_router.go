package routers

import (
	"github.com/gin-gonic/gin"

	authdto "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api/auth"
	enums "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/database/schemas/enums"
	endpoints "github.com/HiIamJeff67/notezy-backend/internal/services/core/transports/gateway/endpoints"
	middlewares "github.com/HiIamJeff67/notezy-backend/internal/services/core/transports/gateway/middlewares"
)

func configureAnonymousAuthRoutes(
	router *gin.RouterGroup,
	endpoint endpoints.AuthEndpointInterface,
) {
	authRoutes := router.Group("/auth")
	{
		authRoutes.POST(
			"/register",
			middlewares.DelegationMiddleware(
				authdto.RegisterOperation,
			),
			endpoint.Register,
		)
		authRoutes.POST(
			"/register/google",
			middlewares.DelegationMiddleware(
				authdto.RegisterViaGoogleOperation,
			),
			endpoint.RegisterViaGoogle,
		)
		authRoutes.POST(
			"/login",
			middlewares.DelegationMiddleware(
				authdto.LoginOperation,
			),
			endpoint.Login,
		)
		authRoutes.POST(
			"/login/google",
			middlewares.DelegationMiddleware(
				authdto.LoginViaGoogleOperation,
			),
			endpoint.LoginViaGoogle,
		)
		authRoutes.POST(
			"/email/code",
			middlewares.DelegationMiddleware(
				authdto.SendAuthCodeOperation,
			),
			endpoint.SendAuthCode,
		)
		authRoutes.PUT(
			"/password/forget",
			middlewares.DelegationMiddleware(
				authdto.ForgetPasswordOperation,
			),
			endpoint.ForgetPassword,
		)
	}
}

func configureAuthenticatedAuthRoutes(
	router *gin.RouterGroup,
	authMiddleware gin.HandlerFunc,
	endpoint endpoints.AuthEndpointInterface,
) {
	authRoutes := router.Group("/auth")
	{
		authRoutes.POST(
			"/logout",
			middlewares.DelegationAuthenticatedMiddleware(
				authdto.LogoutOperation,
			),
			authMiddleware,
			endpoint.Logout,
		)
		authRoutes.PUT(
			"/email/validate",
			middlewares.DelegationAuthenticatedMiddleware(
				authdto.ValidateEmailOperation,
			),
			authMiddleware,
			middlewares.CSRFMiddleware(),
			endpoint.ValidateEmail,
		)
		authRoutes.PUT(
			"/email/reset",
			middlewares.DelegationAuthenticatedMiddleware(
				authdto.ResetEmailOperation,
			),
			authMiddleware,
			middlewares.UserRoleMiddleware(enums.UserRole_Normal),
			middlewares.CSRFMiddleware(),
			endpoint.ResetEmail,
		)
		authRoutes.PUT(
			"/me/reset",
			middlewares.DelegationAuthenticatedMiddleware(
				authdto.ResetMeOperation,
			),
			authMiddleware,
			middlewares.CSRFMiddleware(),
			endpoint.ResetMe,
		)
		authRoutes.DELETE(
			"/me/delete",
			middlewares.DelegationAuthenticatedMiddleware(
				authdto.DeleteMeOperation,
			),
			authMiddleware,
			middlewares.CSRFMiddleware(),
			endpoint.DeleteMe,
		)
	}
}
