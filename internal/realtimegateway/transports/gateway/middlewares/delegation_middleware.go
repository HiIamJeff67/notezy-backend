package middlewares

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	realtimegatewaycontract "github.com/HiIamJeff67/notezy-backend/contracts/realtime-gateway/v1"
	exceptions "github.com/HiIamJeff67/notezy-backend/internal/exceptions"
	sharedtokens "github.com/HiIamJeff67/notezy-backend/shared/tokens"
)

func DelegationMiddleware(expectedOperation string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		delegationClaims, err := sharedtokens.ParseDelegationToken(strings.TrimPrefix(ctx.GetHeader("Authorization"), "Bearer "))
		if err != nil || delegationClaims.Actor != "gateway" {
			exception := exceptions.New(
				"Unauthorized",
				"RealtimeGateway",
				"VerifyDelegation",
				"The internal delegation credential is invalid",
				http.StatusUnauthorized,
			)
			ctx.AbortWithStatusJSON(exception.HTTPStatusCode(), realtimegatewaycontract.Response[struct{}]{
				Version: realtimegatewaycontract.Version,
				Metadata: realtimegatewaycontract.ResponseMetadata{
					RequestId:   ctx.GetHeader("X-Request-Id"),
					RespondedAt: time.Now(),
				},
				Data:      struct{}{},
				Exception: exception,
			})
			return
		}

		payload, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			exception := exceptions.New(
				"InvalidRequest",
				"RealtimeGateway",
				"VerifyDelegation",
				"Failed to read the RealtimeGateway request",
				http.StatusBadRequest,
			).WithOrigin(err)
			ctx.AbortWithStatusJSON(exception.HTTPStatusCode(), realtimegatewaycontract.Response[struct{}]{
				Version: realtimegatewaycontract.Version,
				Metadata: realtimegatewaycontract.ResponseMetadata{
					RequestId:   ctx.GetHeader("X-Request-Id"),
					RespondedAt: time.Now(),
				},
				Data:      struct{}{},
				Exception: exception,
			})
			return
		}

		request := &realtimegatewaycontract.Request[json.RawMessage]{}
		if err := json.Unmarshal(payload, request); err != nil ||
			request.GetVersion() != realtimegatewaycontract.Version ||
			(expectedOperation != "" && request.GetOperation() != expectedOperation) ||
			delegationClaims.Operation != request.GetOperation() ||
			delegationClaims.RequestId != request.GetMetadata().RequestId {
			exception := exceptions.New(
				"InvalidDelegation",
				"RealtimeGateway",
				"VerifyDelegation",
				"The delegation credential does not match the RealtimeGateway request",
				http.StatusUnauthorized,
			)
			ctx.AbortWithStatusJSON(exception.HTTPStatusCode(), realtimegatewaycontract.Response[struct{}]{
				Version: realtimegatewaycontract.Version,
				Metadata: realtimegatewaycontract.ResponseMetadata{
					RequestId:   ctx.GetHeader("X-Request-Id"),
					RespondedAt: time.Now(),
				},
				Data:      struct{}{},
				Exception: exception,
			})
			return
		}

		ctx.Request.Body = io.NopCloser(bytes.NewReader(payload))
		ctx.Next()
	}
}
