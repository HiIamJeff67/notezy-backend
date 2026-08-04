package middlewares

import (
	"fmt"
	"slices"

	"github.com/gin-gonic/gin"

	enumcontract "github.com/HiIamJeff67/notezy-backend/contracts/types/enums"
	contexts "github.com/HiIamJeff67/notezy-backend/internal/gateway/contexts"
)

var accessControlPermissions = []enumcontract.AccessControlPermission{
	enumcontract.AccessControlPermission_Read,
	enumcontract.AccessControlPermission_Write,
	enumcontract.AccessControlPermission_Admin,
	enumcontract.AccessControlPermission_Owner,
}

func AllowedPermissionsAbove(permission enumcontract.AccessControlPermission) gin.HandlerFunc {
	index := slices.Index(accessControlPermissions, permission)
	if index < 0 {
		panic(fmt.Sprintf("invalid access control permission: %s", permission))
	}

	return AllowedPermissionsWithin(accessControlPermissions[index:]...)
}

func AllowedPermissionsBelow(permission enumcontract.AccessControlPermission) gin.HandlerFunc {
	index := slices.Index(accessControlPermissions, permission)
	if index < 0 {
		panic(fmt.Sprintf("invalid access control permission: %s", permission))
	}

	return AllowedPermissionsWithin(accessControlPermissions[:index+1]...)
}

func AllowedPermissionsWithin(allowedPermissions ...enumcontract.AccessControlPermission) gin.HandlerFunc {
	if len(allowedPermissions) == 0 {
		panic("allowed permissions are required")
	}
	for _, permission := range allowedPermissions {
		if !slices.Contains(accessControlPermissions, permission) {
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
