package middlewares

import (
	"fmt"
	"slices"

	"github.com/gin-gonic/gin"

	contexts "github.com/HiIamJeff67/notezy-backend/app/contexts"
	enums "github.com/HiIamJeff67/notezy-backend/app/models/schemas/enums"
)

func AllowedPermissionsAbove(permission enums.AccessControlPermission) gin.HandlerFunc {
	index := slices.Index(enums.AllAccessControlPermissions, permission)
	if index < 0 {
		panic(fmt.Sprintf("invalid access control permission: %s", permission))
	}

	return AllowedPermissionsWithin(enums.AllAccessControlPermissions[index:]...)
}

func AllowedPermissionsBelow(permission enums.AccessControlPermission) gin.HandlerFunc {
	index := slices.Index(enums.AllAccessControlPermissions, permission)
	if index < 0 {
		panic(fmt.Sprintf("invalid access control permission: %s", permission))
	}

	return AllowedPermissionsWithin(enums.AllAccessControlPermissions[:index+1]...)
}

func AllowedPermissionsWithin(allowedPermissions ...enums.AccessControlPermission) gin.HandlerFunc {
	if len(allowedPermissions) == 0 {
		panic("allowed permissions are required")
	}
	for _, permission := range allowedPermissions {
		if !slices.Contains(enums.AllAccessControlPermissions, permission) {
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
