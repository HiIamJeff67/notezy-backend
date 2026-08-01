package middlewares

import (
	"fmt"
	"slices"

	"github.com/gin-gonic/gin"

	contexts "github.com/HiIamJeff67/notezy-backend/internal/gateway/contexts"
	sharedtypes "github.com/HiIamJeff67/notezy-backend/internal/shared/types"
)

func AllowedPermissionsAbove(permission sharedtypes.AccessControlPermission) gin.HandlerFunc {
	index := slices.Index(sharedtypes.AllAccessControlPermissions, permission)
	if index < 0 {
		panic(fmt.Sprintf("invalid access control permission: %s", permission))
	}

	return AllowedPermissionsWithin(sharedtypes.AllAccessControlPermissions[index:]...)
}

func AllowedPermissionsBelow(permission sharedtypes.AccessControlPermission) gin.HandlerFunc {
	index := slices.Index(sharedtypes.AllAccessControlPermissions, permission)
	if index < 0 {
		panic(fmt.Sprintf("invalid access control permission: %s", permission))
	}

	return AllowedPermissionsWithin(sharedtypes.AllAccessControlPermissions[:index+1]...)
}

func AllowedPermissionsWithin(allowedPermissions ...sharedtypes.AccessControlPermission) gin.HandlerFunc {
	if len(allowedPermissions) == 0 {
		panic("allowed permissions are required")
	}
	for _, permission := range allowedPermissions {
		if !slices.Contains(sharedtypes.AllAccessControlPermissions, permission) {
			panic(fmt.Sprintf("invalid access control permission: %s", permission))
		}
	}

	return func(ctx *gin.Context) {
		ctx.Request = ctx.Request.WithContext(
			contexts.WithAllowedPermissions(ctx.Request.Context(), allowedPermissions),
		)

		ctx.Next()
	}
}
