package middlewares

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	gatewaycontract "github.com/HiIamJeff67/notezy-backend/contracts/gateway/v1"
	exceptions "github.com/HiIamJeff67/notezy-backend/internal/exceptions"
	contexts "github.com/HiIamJeff67/notezy-backend/internal/services/core/contexts"
	enums "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/database/schemas/enums"
	sharedtokens "github.com/HiIamJeff67/notezy-backend/shared/tokens"
)

func DelegationAuthenticatedMiddleware(expectedOperation string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		delegationClaims, err := sharedtokens.ParseDelegationToken(strings.TrimPrefix(ctx.GetHeader("Authorization"), "Bearer "))
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gatewaycontract.Response[struct{}]{
				Version: gatewaycontract.Version,
				Metadata: gatewaycontract.ResponseMetadata{
					RequestId:   ctx.GetHeader("X-Request-Id"),
					RespondedAt: time.Now(),
				},
				Data: struct{}{},
				Exception: exceptions.New(
					"Unauthorized",
					"Core",
					"VerifyDelegation",
					"invalid internal delegation credential",
					http.StatusUnauthorized,
				),
			})
			return
		}

		payload, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gatewaycontract.Response[struct{}]{
				Version: gatewaycontract.Version,
				Metadata: gatewaycontract.ResponseMetadata{
					RequestId:   ctx.GetHeader("X-Request-Id"),
					RespondedAt: time.Now(),
				},
				Data: struct{}{},
				Exception: exceptions.New(
					"InvalidRequest",
					"Core",
					"VerifyDelegation",
					"failed to read the Core service request",
					http.StatusBadRequest,
				),
			})
			return
		}

		request := &gatewaycontract.Request[json.RawMessage]{}
		if err := json.Unmarshal(payload, request); err != nil ||
			request.GetVersion() != gatewaycontract.Version ||
			(expectedOperation != "" && request.GetOperation() != expectedOperation) ||
			delegationClaims.Operation != request.GetOperation() ||
			delegationClaims.RequestId != request.GetMetadata().RequestId {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gatewaycontract.Response[struct{}]{
				Version: gatewaycontract.Version,
				Metadata: gatewaycontract.ResponseMetadata{
					RequestId:   ctx.GetHeader("X-Request-Id"),
					RespondedAt: time.Now(),
				},
				Data: struct{}{},
				Exception: exceptions.New(
					"InvalidDelegation",
					"Core",
					"VerifyDelegation",
					"delegation credential does not match the request",
					http.StatusUnauthorized,
				),
			})
			return
		}
		ctx.Request.Body = io.NopCloser(bytes.NewReader(payload))

		permissions := make([]enums.AccessControlPermission, 0, len(delegationClaims.AllowedPermissions))
		for _, permissionString := range delegationClaims.AllowedPermissions {
			permission, err := enums.ConvertStringToAccessControlPermission(permissionString)
			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gatewaycontract.Response[struct{}]{
					Version: gatewaycontract.Version,
					Metadata: gatewaycontract.ResponseMetadata{
						RequestId:   ctx.GetHeader("X-Request-Id"),
						RespondedAt: time.Now(),
					},
					Data: struct{}{},
					Exception: exceptions.New(
						"InvalidDelegation",
						"Core",
						"VerifyDelegation",
						"invalid delegated permission",
						http.StatusUnauthorized,
					),
				})
				return
			}
			permissions = append(permissions, *permission)
		}

		if delegationClaims.UserSubject == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gatewaycontract.Response[struct{}]{
				Version: gatewaycontract.Version,
				Metadata: gatewaycontract.ResponseMetadata{
					RequestId:   ctx.GetHeader("X-Request-Id"),
					RespondedAt: time.Now(),
				},
				Data: struct{}{},
				Exception: exceptions.New(
					"InvalidDelegation",
					"Core",
					"VerifyDelegation",
					"delegation user subject is required",
					http.StatusUnauthorized,
				),
			})
			return
		}
		userSubject, err := uuid.Parse(delegationClaims.UserSubject)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gatewaycontract.Response[struct{}]{
				Version: gatewaycontract.Version,
				Metadata: gatewaycontract.ResponseMetadata{
					RequestId:   ctx.GetHeader("X-Request-Id"),
					RespondedAt: time.Now(),
				},
				Data: struct{}{},
				Exception: exceptions.New(
					"InvalidDelegation",
					"Core",
					"VerifyDelegation",
					"delegation user subject is invalid",
					http.StatusUnauthorized,
				),
			})
			return
		}

		ctx.Request = ctx.Request.WithContext(
			contexts.WithActorUserPublicId(ctx.Request.Context(), userSubject),
		)
		ctx.Request = ctx.Request.WithContext(
			contexts.WithAllowedPermissions(ctx.Request.Context(), permissions),
		)
		ctx.Next()
	}
}
