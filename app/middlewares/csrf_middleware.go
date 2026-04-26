package middlewares

import (
	"strings"

	"github.com/gin-gonic/gin"

	caches "notezy-backend/app/caches"
	contexts "notezy-backend/app/contexts"
	exceptions "notezy-backend/app/exceptions"
	tokens "notezy-backend/app/tokens"
	types "notezy-backend/shared/types"
)

/*
A Middleware to provider CSRF token validation which should be placed after AuthMiddleware
*/
func CSRFMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userName, exception := contexts.GetAndConvertContextFieldToString(ctx, types.ContextFieldName_User_Name)
		if exception != nil {
			exceptions.Context.MissPlacingOrWrongMiddlewareOrder(
				"Cannot find the userPlan, " +
					"please make sure the AuthMiddleware() is placing before the CSRFMiddleware()",
			).Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}

		csrfToken := ctx.GetHeader("X-CSRF-Token")
		if len(strings.TrimSpace(csrfToken)) <= 0 {
			exceptions.Token.FailedToExtractOrValidateCSRFToken().Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}

		userDataCache, exception := caches.GetUserDataCache(*userName)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}

		claims, exception := tokens.ValidateCSRFToken(csrfToken, userDataCache.CSRFToken)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}

		if tokens.IsCSRFTokenExpiringSoon(claims) {
			newToken, exception := tokens.GenerateCSRFToken()
			if exception == nil {
				dto := caches.UpdateUserDataCacheDto{
					CSRFToken: newToken,
				}
				caches.UpdateUserDataCache(*userName, dto)

				ctx.Header("X-CSRF-Token", *newToken)

				ctx.Set(types.ContextFieldName_IsNewTokens.String(), true)
				ctx.Set(types.ContextFieldName_CSRFToken.String(), *newToken)
			}
		}

		ctx.Next()
	}
}

// eyJzaWduYXR1cmUiOiJmWkZ5MkFMS2o5U2ptMmozRnhZRVM4Q2JJSnNvLzNMMGVQWitDQ3RLOXA0PSIsImV4cGlyZXNBdCI6IjIwMjYtMDQtMjlUMTU6Mzc6NDQuNTU3Mzg5ODM5WiIsImlzc3VlZEF0IjoiMjAyNi0wNC0yMlQxNTozNzo0NC41NTczODk4MzlaIn0=

// eyJzaWduYXR1cmUiOiJmWkZ5MkFMS2o5U2ptMmozRnhZRVM4Q2JJSnNvLzNMMGVQWitDQ3RLOXA0PSIsImV4cGlyZXNBdCI6IjIwMjYtMDQtMjlUMTU6Mzc6NDQuNTU3Mzg5ODM5WiIsImlzc3VlZEF0IjoiMjAyNi0wNC0yMlQxNTozNzo0NC41NTczODk4MzlaIn0=
