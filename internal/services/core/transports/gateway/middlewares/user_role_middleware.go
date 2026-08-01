package middlewares

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	core "github.com/HiIamJeff67/notezy-backend/contracts/core/v1"
	exceptions "github.com/HiIamJeff67/notezy-backend/internal/exceptions"
	enums "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/schemas/enums"
	types "github.com/HiIamJeff67/notezy-backend/internal/shared/types"
)

func UserRoleMiddleware(atLeastUserRole enums.UserRole) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		currentUserRoleValue, exists := ctx.Get(types.ContextFieldName_User_Role.String())
		currentUserRole, ok := currentUserRoleValue.(enums.UserRole)
		if !exists || !ok {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, core.Response[struct{}]{
				Version: core.Version,
				Metadata: core.ResponseMetadata{
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
				ctx.AbortWithStatusJSON(http.StatusForbidden, core.Response[struct{}]{
					Version: core.Version,
					Metadata: core.ResponseMetadata{
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

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, core.Response[struct{}]{
			Version: core.Version,
			Metadata: core.ResponseMetadata{
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
