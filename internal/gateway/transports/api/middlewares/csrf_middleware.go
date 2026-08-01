package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	caches "github.com/HiIamJeff67/notezy-backend/internal/caches"
	cacheinputs "github.com/HiIamJeff67/notezy-backend/internal/caches/inputs"
	exceptions "github.com/HiIamJeff67/notezy-backend/internal/exceptions"
	contexts "github.com/HiIamJeff67/notezy-backend/internal/gateway/contexts"
	responsewriter "github.com/HiIamJeff67/notezy-backend/internal/gateway/responsewriter"
	sharedtokens "github.com/HiIamJeff67/notezy-backend/internal/shared/tokens"
	types "github.com/HiIamJeff67/notezy-backend/internal/shared/types"
)

/*
A Middleware to provider CSRF token validation which should be placed after AuthMiddleware
*/
func CSRFMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userName, exception := contexts.GetAndConvertContextFieldToString(ctx, types.ContextFieldName_User_Name)
		if exception != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.New(
				"WrongMiddlewareOrder",
				"Context",
				"Middleware",
				"Cannot find the userPlan, "+
					"please make sure the AuthMiddleware() is placing before the CSRFMiddleware()",
				http.StatusInternalServerError,
				true,
			), ctx)
			return
		}

		csrfToken := ctx.GetHeader("X-CSRF-Token")
		if len(strings.TrimSpace(csrfToken)) <= 0 {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.New(
				"InvalidCSRFToken",
				"Token",
				"ValidateCSRFToken",
				"CSRF token is missing or invalid",
				http.StatusUnauthorized,
			), ctx)
			return
		}

		userDataCache, exception := caches.UserDataStore.Get(*userName)
		if exception != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
			return
		}

		claims, err := sharedtokens.ValidateCSRFToken(csrfToken, userDataCache.CSRFToken)
		if err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.New(
				"InvalidCSRFToken",
				"Token",
				"ValidateCSRFToken",
				"CSRF token is invalid",
				http.StatusUnauthorized,
			).WithOrigin(err), ctx)
			return
		}

		if sharedtokens.IsCSRFTokenExpiringSoon(claims) {
			newToken, err := sharedtokens.GenerateCSRFToken(sharedtokens.CSRFTokenClaims{})
			if err != nil {
				responsewriter.SafelyAbortAndResponseWithJSON(exceptions.New(
					"GenerationFailed",
					"Token",
					"GenerateCSRFToken",
					"Failed to generate the CSRF token",
					http.StatusInternalServerError,
					true,
				).WithOrigin(err), ctx)
				return
			}

			input := cacheinputs.UpdateUserDataCacheInput{
				CSRFToken: newToken,
			}
			caches.UserDataStore.Update(*userName, input)

			ctx.Header("X-CSRF-Token", *newToken)

			ctx.Set(types.ContextFieldName_IsNewTokens.String(), true)
			ctx.Set(types.ContextFieldName_CSRFToken.String(), *newToken)
		}

		ctx.Next()
	}
}

// eyJzaWduYXR1cmUiOiJmWkZ5MkFMS2o5U2ptMmozRnhZRVM4Q2JJSnNvLzNMMGVQWitDQ3RLOXA0PSIsImV4cGlyZXNBdCI6IjIwMjYtMDQtMjlUMTU6Mzc6NDQuNTU3Mzg5ODM5WiIsImlzc3VlZEF0IjoiMjAyNi0wNC0yMlQxNTozNzo0NC41NTczODk4MzlaIn0=

// eyJzaWduYXR1cmUiOiJmWkZ5MkFMS2o5U2ptMmozRnhZRVM4Q2JJSnNvLzNMMGVQWitDQ3RLOXA0PSIsImV4cGlyZXNBdCI6IjIwMjYtMDQtMjlUMTU6Mzc6NDQuNTU3Mzg5ODM5WiIsImlzc3VlZEF0IjoiMjAyNi0wNC0yMlQxNTozNzo0NC41NTczODk4MzlaIn0=
