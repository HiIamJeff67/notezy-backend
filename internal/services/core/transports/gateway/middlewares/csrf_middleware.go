package middlewares

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	core "github.com/HiIamJeff67/notezy-backend/contracts/core/v1"
	caches "github.com/HiIamJeff67/notezy-backend/internal/caches"
	cacheinputs "github.com/HiIamJeff67/notezy-backend/internal/caches/inputs"
	exceptions "github.com/HiIamJeff67/notezy-backend/internal/exceptions"
	contexts "github.com/HiIamJeff67/notezy-backend/internal/services/core/contexts"
	sharedtokens "github.com/HiIamJeff67/notezy-backend/internal/shared/tokens"
	types "github.com/HiIamJeff67/notezy-backend/internal/shared/types"
)

func CSRFMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		actorUserName, exception := contexts.GetActorUserName(ctx.Request.Context())
		if exception != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, core.Response[struct{}]{
				Version:   core.Version,
				Metadata:  core.ResponseMetadata{RequestId: ctx.GetHeader("X-Request-Id"), RespondedAt: time.Now()},
				Data:      struct{}{},
				Exception: exception.ToPublic(),
			})
			return
		}
		csrfToken := ctx.GetHeader("X-CSRF-Token")
		if strings.TrimSpace(csrfToken) == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, core.Response[struct{}]{
				Version:   core.Version,
				Metadata:  core.ResponseMetadata{RequestId: ctx.GetHeader("X-Request-Id"), RespondedAt: time.Now()},
				Data:      struct{}{},
				Exception: exceptions.New("InvalidCSRFToken", "Token", "ValidateCSRFToken", "the CSRF token is missing or invalid", http.StatusUnauthorized),
			})
			return
		}
		userDataCache, exception := caches.UserDataStore.Get(actorUserName)
		if exception != nil {
			ctx.AbortWithStatusJSON(exception.HTTPStatusCode(), core.Response[struct{}]{
				Version:   core.Version,
				Metadata:  core.ResponseMetadata{RequestId: ctx.GetHeader("X-Request-Id"), RespondedAt: time.Now()},
				Data:      struct{}{},
				Exception: exception.ToPublic(),
			})
			return
		}
		claims, err := sharedtokens.ValidateCSRFToken(csrfToken, userDataCache.CSRFToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, core.Response[struct{}]{
				Version:   core.Version,
				Metadata:  core.ResponseMetadata{RequestId: ctx.GetHeader("X-Request-Id"), RespondedAt: time.Now()},
				Data:      struct{}{},
				Exception: exceptions.New("InvalidCSRFToken", "Token", "ValidateCSRFToken", "the CSRF token is invalid", http.StatusUnauthorized).WithOrigin(err),
			})
			return
		}
		if sharedtokens.IsCSRFTokenExpiringSoon(claims) {
			newCSRFToken, err := sharedtokens.GenerateCSRFToken(sharedtokens.CSRFTokenClaims{})
			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, core.Response[struct{}]{
					Version:   core.Version,
					Metadata:  core.ResponseMetadata{RequestId: ctx.GetHeader("X-Request-Id"), RespondedAt: time.Now()},
					Data:      struct{}{},
					Exception: exceptions.New("GenerationFailed", "Token", "GenerateCSRFToken", "failed to generate a CSRF token", http.StatusInternalServerError, true).WithOrigin(err),
				})
				return
			}
			if exception := caches.UserDataStore.Update(actorUserName, cacheinputs.UpdateUserDataCacheInput{CSRFToken: newCSRFToken}); exception != nil {
				ctx.AbortWithStatusJSON(exception.HTTPStatusCode(), core.Response[struct{}]{
					Version:   core.Version,
					Metadata:  core.ResponseMetadata{RequestId: ctx.GetHeader("X-Request-Id"), RespondedAt: time.Now()},
					Data:      struct{}{},
					Exception: exception.ToPublic(),
				})
				return
			}
			ctx.Header(core.AuthRefreshedHeader, "true")
			if accessToken, accessTokenExists := cookieValue(ctx, types.ValidCookieName_AccessToken.String()); accessTokenExists {
				ctx.Header(core.SetAccessTokenHeader, accessToken)
			}
			ctx.Header(core.SetCSRFTokenHeader, *newCSRFToken)
		}

		ctx.Next()
	}
}
