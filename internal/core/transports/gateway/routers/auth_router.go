package routers

import (
	"github.com/gin-gonic/gin"

	apicontract "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api/auth"

	userdata "github.com/HiIamJeff67/notezy-backend/internal/core/data/cache/userdata"
	enums "github.com/HiIamJeff67/notezy-backend/internal/core/data/database/schemas/enums"
	endpoints "github.com/HiIamJeff67/notezy-backend/internal/core/transports/gateway/endpoints"
	middlewares "github.com/HiIamJeff67/notezy-backend/internal/core/transports/gateway/middlewares"
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
				apicontract.RegisterOperation,
			),
			endpoint.Register,
		)
		authRoutes.POST(
			"/register/google",
			middlewares.DelegationMiddleware(
				apicontract.RegisterViaGoogleOperation,
			),
			endpoint.RegisterViaGoogle,
		)
		authRoutes.POST(
			"/login",
			middlewares.DelegationMiddleware(
				apicontract.LoginOperation,
			),
			endpoint.Login,
		)
		authRoutes.POST(
			"/login/google",
			middlewares.DelegationMiddleware(
				apicontract.LoginViaGoogleOperation,
			),
			endpoint.LoginViaGoogle,
		)
		authRoutes.POST(
			"/email/code",
			middlewares.DelegationMiddleware(
				apicontract.SendAuthCodeOperation,
			),
			endpoint.SendAuthCode,
		)
		authRoutes.PUT(
			"/password/forget",
			middlewares.DelegationMiddleware(
				apicontract.ForgetPasswordOperation,
			),
			endpoint.ForgetPassword,
		)
	}
}

func configureAuthenticatedAuthRoutes(
	router *gin.RouterGroup,
	authMiddleware gin.HandlerFunc,
	endpoint endpoints.AuthEndpointInterface,
	userDataCacheClient *userdata.UserDataCacheClient,
) {
	authRoutes := router.Group("/auth")
	{
		authRoutes.POST(
			"/logout",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.LogoutOperation,
			),
			authMiddleware,
			endpoint.Logout,
		)
		authRoutes.PUT(
			"/email/validate",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.ValidateEmailOperation,
			),
			authMiddleware,
			middlewares.CSRFMiddleware(userDataCacheClient),
			endpoint.ValidateEmail,
		)
		authRoutes.PUT(
			"/email/reset",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.ResetEmailOperation,
			),
			authMiddleware,
			middlewares.UserRoleMiddleware(enums.UserRole_Normal),
			middlewares.CSRFMiddleware(userDataCacheClient),
			endpoint.ResetEmail,
		)
		authRoutes.PUT(
			"/me/reset",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.ResetMeOperation,
			),
			authMiddleware,
			middlewares.CSRFMiddleware(userDataCacheClient),
			endpoint.ResetMe,
		)
		authRoutes.DELETE(
			"/me/delete",
			middlewares.DelegationAuthenticatedMiddleware(
				apicontract.DeleteMeOperation,
			),
			authMiddleware,
			middlewares.CSRFMiddleware(userDataCacheClient),
			endpoint.DeleteMe,
		)
	}
}
