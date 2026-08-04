package middlewares

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	gatewaycontract "github.com/HiIamJeff67/notezy-backend/contracts/gateway/v1"
	exceptions "github.com/HiIamJeff67/notezy-backend/internal/exceptions"
	enums "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/database/schemas/enums"
	sharedcontexts "github.com/HiIamJeff67/notezy-backend/shared/lib/contexts"
)

func UserRoleMiddleware(atLeastUserRole enums.UserRole) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		currentUserRoleValue, exists := ctx.Get(sharedcontexts.ContextFieldName_User_Role.String())
		currentUserRole, ok := currentUserRoleValue.(enums.UserRole)
		if !exists || !ok {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gatewaycontract.Response[struct{}]{
				Version: gatewaycontract.Version,
				Metadata: gatewaycontract.ResponseMetadata{
					RequestId:   ctx.GetHeader("X-Request-Id"),
					RespondedAt: time.Now(),
				},
				Data: struct{}{},
				Exception: exceptions.New(
					"AuthenticationContextInvalid",
					"Core",
					"AuthorizeRequest",
					"the authenticated user role is unavailable",
					http.StatusInternalServerError,
					true,
				),
			})
			return
		}

		if currentUserRole == atLeastUserRole {
			ctx.Next()
			return
		}
		for _, userRole := range enums.AllUserRoles {
			if userRole == currentUserRole {
				ctx.Next()
				return
			}
			if userRole == atLeastUserRole {
				ctx.AbortWithStatusJSON(http.StatusForbidden, gatewaycontract.Response[struct{}]{
					Version: gatewaycontract.Version,
					Metadata: gatewaycontract.ResponseMetadata{
						RequestId:   ctx.GetHeader("X-Request-Id"),
						RespondedAt: time.Now(),
					},
					Data: struct{}{},
					Exception: exceptions.New(
						"PermissionDeniedDueToUserRole",
						"Auth",
						"AuthorizeRequest",
						fmt.Sprintf("the current user role of %v cannot access this operation", currentUserRole),
						http.StatusForbidden,
					),
				})
				return
			}
		}

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gatewaycontract.Response[struct{}]{
			Version: gatewaycontract.Version,
			Metadata: gatewaycontract.ResponseMetadata{
				RequestId:   ctx.GetHeader("X-Request-Id"),
				RespondedAt: time.Now(),
			},
			Data: struct{}{},
			Exception: exceptions.New(
				"AuthenticationContextInvalid",
				"Core",
				"AuthorizeRequest",
				"the authenticated user role is invalid",
				http.StatusInternalServerError,
				true,
			),
		})
	}
}
