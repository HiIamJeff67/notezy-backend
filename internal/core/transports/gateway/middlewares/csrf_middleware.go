package middlewares

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	sharedtokens "github.com/HiIamJeff67/notezy-backend/shared/tokens"

	sharedcontexts "github.com/HiIamJeff67/notezy-backend/shared/lib/contexts"

	gatewaycontract "github.com/HiIamJeff67/notezy-backend/contracts/gateway/v1"
	exceptions "github.com/HiIamJeff67/notezy-backend/contracts/types/exceptions"

	contexts "github.com/HiIamJeff67/notezy-backend/internal/core/contexts"
	userdata "github.com/HiIamJeff67/notezy-backend/internal/core/data/cache/userdata"
	cacheinputs "github.com/HiIamJeff67/notezy-backend/internal/core/data/cache/userdata/inputs"
)

func CSRFMiddleware(userDataCacheClient *userdata.UserDataCacheClient) gin.HandlerFunc {
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
		request := &gatewaycontract.Request[json.RawMessage]{}
		if ctx.Request.ContentLength != 0 {
			_ = ctx.ShouldBindBodyWithJSON(request)
		}
		csrfToken := request.Tokens.CSRFToken
		if strings.TrimSpace(csrfToken) == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gatewaycontract.Response[struct{}]{
				Version:   gatewaycontract.Version,
				Metadata:  gatewaycontract.ResponseMetadata{RequestId: ctx.GetHeader("X-Request-Id"), RespondedAt: time.Now()},
				Data:      struct{}{},
				Exception: exceptions.New("InvalidCSRFToken", "Token", "ValidateCSRFToken", "the CSRF token is missing or invalid", http.StatusUnauthorized),
			})
			return
		}
		userDataCache, exception := userDataCacheClient.Get(actorUserName)
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
			if exception := userDataCacheClient.Update(actorUserName, cacheinputs.UpdateUserDataCacheInput{CSRFToken: newCSRFToken}); exception != nil {
				ctx.AbortWithStatusJSON(exception.HTTPStatusCode(), gatewaycontract.Response[struct{}]{
					Version:   gatewaycontract.Version,
					Metadata:  gatewaycontract.ResponseMetadata{RequestId: ctx.GetHeader("X-Request-Id"), RespondedAt: time.Now()},
					Data:      struct{}{},
					Exception: exception.ToPublic(),
				})
				return
			}
			ctx.Set(sharedcontexts.ContextFieldName_IsNewTokens.String(), true)
			if strings.TrimSpace(request.Tokens.AccessToken) != "" {
				ctx.Set(sharedcontexts.ContextFieldName_AccessToken.String(), request.Tokens.AccessToken)
			}
			ctx.Set(sharedcontexts.ContextFieldName_CSRFToken.String(), *newCSRFToken)
		}

		ctx.Next()
	}
}
