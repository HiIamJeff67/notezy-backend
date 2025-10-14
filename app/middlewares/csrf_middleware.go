package middlewares

import (
	"strings"

	"github.com/gin-gonic/gin"

	caches "notezy-backend/app/caches"
	contexts "notezy-backend/app/contexts"
	exceptions "notezy-backend/app/exceptions"
	tokens "notezy-backend/app/tokens"
	constants "notezy-backend/shared/constants"
)

/*
A Middleware to provider CSRF token validation which should be placed after AuthMiddleware
*/
func CSRFMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
		if exception != nil {
			exceptions.Auth.MissPlacingOrWrongMiddlewareOrder(
				"Cannot find the userPlan, " +
					"please make sure the AuthMiddleware() is placing before the CSRFMiddleware()",
			).Log().SafelyResponseWithJSON(ctx)
			return
		}

		csrfToken := ctx.GetHeader("X-CSRF-Token")
		if len(strings.TrimSpace(csrfToken)) <= 0 {
			exceptions.Token.FailedToExtractOrValidateCSRFToken().Log().SafelyResponseWithJSON(ctx)
			return
		}

		userDataCache, exception := caches.GetUserDataCache(*userId)
		if exception != nil {
			exception.Log().SafelyResponseWithJSON(ctx)
			return
		}

		claims, exception := tokens.ValidateCSRFToken(csrfToken, userDataCache.CSRFToken)
		if exception != nil {
			exception.Log().SafelyResponseWithJSON(ctx)
			return
		}

		if tokens.IsCSRFTokenExpiringSoon(claims) {
			newToken, exception := tokens.GenerateCSRFToken()
			if exception == nil {
				dto := caches.UpdateUserDataCacheDto{
					CSRFToken: newToken,
				}
				caches.UpdateUserDataCache(*userId, dto)

				ctx.Header("X-CSRF-Token", *newToken)
			}
		}

		ctx.Next()
	}
}
