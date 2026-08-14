package middlewares

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	gatewaycontract "github.com/HiIamJeff67/notezy-backend/contracts/gateway/v1"
	exceptions "github.com/HiIamJeff67/notezy-backend/contracts/types/exceptions"
	gatewaycontexts "github.com/HiIamJeff67/notezy-backend/internal/apigateway/contexts"
	sharedtokens "github.com/HiIamJeff67/notezy-backend/shared/tokens"
)

// KeyMiddleware is the APIGateway edge check. It intentionally verifies only
// presence and format; ownership and revocation are authoritative in Core's
// APIKeyMiddleware after the delegation credential is verified.
func KeyMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.Request.Method == http.MethodOptions {
			ctx.Next()
			return
		}
		key := strings.TrimSpace(ctx.GetHeader("X-API-Key"))
		if key == "" || sharedtokens.ValidateAPIKeyFormat(key) != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gatewaycontract.Response[struct{}]{
				Version: gatewaycontract.Version,
				Metadata: gatewaycontract.ResponseMetadata{
					RequestId:   ctx.GetHeader("X-Request-Id"),
					RespondedAt: time.Now(),
				},
				Data: struct{}{},
				Exception: exceptions.New(
					"Unauthorized",
					"Gateway",
					"AuthenticateAPIKey",
					"a valid API key is required",
					http.StatusUnauthorized,
				),
			})
			return
		}
		gatewaycontexts.SetGatewaySource(ctx, sharedtokens.GatewaySourceAPI)
		ctx.Next()
	}
}
