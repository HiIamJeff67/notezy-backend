package middlewares

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	gatewaycontract "github.com/HiIamJeff67/notezy-backend/contracts/gateway/v1"
	exceptions "github.com/HiIamJeff67/notezy-backend/internal/exceptions"
	contexts "github.com/HiIamJeff67/notezy-backend/internal/services/core/contexts"
	userdata "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/cache/userdata"
	cacheinputs "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/cache/userdata/inputs"
	cookies "github.com/HiIamJeff67/notezy-backend/shared/cookies"
	sharedtokens "github.com/HiIamJeff67/notezy-backend/shared/tokens"
)

func CSRFMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		actorUserName, exception := contexts.GetActorUserName(ctx.Request.Context())
		if exception != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gatewaycontract.Response[struct{}]{
				Version:   gatewaycontract.Version,
				Metadata:  gatewaycontract.ResponseMetadata{RequestId: ctx.GetHeader("X-Request-Id"), RespondedAt: time.Now()},
				Data:      struct{}{},
				Exception: exception.ToPublic(),
			})
			return
		}
		csrfToken := ctx.GetHeader("X-CSRF-Token")
		if strings.TrimSpace(csrfToken) == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gatewaycontract.Response[struct{}]{
				Version:   gatewaycontract.Version,
				Metadata:  gatewaycontract.ResponseMetadata{RequestId: ctx.GetHeader("X-Request-Id"), RespondedAt: time.Now()},
				Data:      struct{}{},
				Exception: exceptions.New("InvalidCSRFToken", "Token", "ValidateCSRFToken", "the CSRF token is missing or invalid", http.StatusUnauthorized),
			})
			return
		}
		userDataCache, exception := userdata.NewUserDataCacheClient().Get(actorUserName)
		if exception != nil {
			ctx.AbortWithStatusJSON(exception.HTTPStatusCode(), gatewaycontract.Response[struct{}]{
				Version:   gatewaycontract.Version,
				Metadata:  gatewaycontract.ResponseMetadata{RequestId: ctx.GetHeader("X-Request-Id"), RespondedAt: time.Now()},
				Data:      struct{}{},
				Exception: exception.ToPublic(),
			})
			return
		}
		claims, err := sharedtokens.ValidateCSRFToken(csrfToken, userDataCache.CSRFToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gatewaycontract.Response[struct{}]{
				Version:   gatewaycontract.Version,
				Metadata:  gatewaycontract.ResponseMetadata{RequestId: ctx.GetHeader("X-Request-Id"), RespondedAt: time.Now()},
				Data:      struct{}{},
				Exception: exceptions.New("InvalidCSRFToken", "Token", "ValidateCSRFToken", "the CSRF token is invalid", http.StatusUnauthorized).WithOrigin(err),
			})
			return
		}
		if sharedtokens.IsCSRFTokenExpiringSoon(claims) {
			newCSRFToken, err := sharedtokens.GenerateCSRFToken(sharedtokens.CSRFTokenClaims{})
			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, gatewaycontract.Response[struct{}]{
					Version:   gatewaycontract.Version,
					Metadata:  gatewaycontract.ResponseMetadata{RequestId: ctx.GetHeader("X-Request-Id"), RespondedAt: time.Now()},
					Data:      struct{}{},
					Exception: exceptions.New("GenerationFailed", "Token", "GenerateCSRFToken", "failed to generate a CSRF token", http.StatusInternalServerError, true).WithOrigin(err),
				})
				return
			}
			if exception := userdata.NewUserDataCacheClient().Update(actorUserName, cacheinputs.UpdateUserDataCacheInput{CSRFToken: newCSRFToken}); exception != nil {
				ctx.AbortWithStatusJSON(exception.HTTPStatusCode(), gatewaycontract.Response[struct{}]{
					Version:   gatewaycontract.Version,
					Metadata:  gatewaycontract.ResponseMetadata{RequestId: ctx.GetHeader("X-Request-Id"), RespondedAt: time.Now()},
					Data:      struct{}{},
					Exception: exception.ToPublic(),
				})
				return
			}
			ctx.Header(gatewaycontract.CoreAuthRefreshed.String(), "true")
			if cookie, err := ctx.Request.Cookie(cookies.ValidCookieName_AccessToken.String()); err == nil && strings.TrimSpace(cookie.Value) != "" {
				ctx.Header(gatewaycontract.CoreSetAccessToken.String(), cookie.Value)
			}
			ctx.Header(gatewaycontract.CoreSetCSRFToken.String(), *newCSRFToken)
		}

		ctx.Next()
	}
}
