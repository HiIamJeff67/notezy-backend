package middlewares

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	core "github.com/HiIamJeff67/notezy-backend/contracts/core/v1"
	exceptions "github.com/HiIamJeff67/notezy-backend/internal/exceptions"
	contexts "github.com/HiIamJeff67/notezy-backend/internal/services/core/contexts"
	enums "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/schemas/enums"
	sharedtokens "github.com/HiIamJeff67/notezy-backend/internal/shared/tokens"
	"github.com/gin-gonic/gin"
)

func DelegationMiddleware(expectedOperation string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		delegationClaims, err := sharedtokens.ParseDelegationToken(strings.TrimPrefix(ctx.GetHeader("Authorization"), "Bearer "))
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, core.Response[struct{}]{
				Version: core.Version,
				Metadata: core.ResponseMetadata{
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
			ctx.AbortWithStatusJSON(http.StatusBadRequest, core.Response[struct{}]{
				Version: core.Version,
				Metadata: core.ResponseMetadata{
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

		request := &core.Request[json.RawMessage]{}
		if err := json.Unmarshal(payload, request); err != nil ||
			request.GetVersion() != core.Version ||
			(expectedOperation != "" && request.GetOperation() != expectedOperation) ||
			delegationClaims.Operation != request.GetOperation() ||
			delegationClaims.RequestId != request.GetMetadata().RequestId {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, core.Response[struct{}]{
				Version: core.Version,
				Metadata: core.ResponseMetadata{
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
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, core.Response[struct{}]{
					Version: core.Version,
					Metadata: core.ResponseMetadata{
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

		ctx.Request = ctx.Request.WithContext(
			contexts.WithAllowedPermissions(ctx.Request.Context(), permissions),
		)
		ctx.Next()
	}
}
