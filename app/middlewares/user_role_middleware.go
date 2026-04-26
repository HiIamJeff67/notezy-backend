package middlewares

import (
	"github.com/gin-gonic/gin"

	exceptions "notezy-backend/app/exceptions"
	enums "notezy-backend/app/models/schemas/enums"
	types "notezy-backend/shared/types"
)

// This UserRoleMiddleware() MUST be processed AFTER the AuthMiddleware()
// so that it can parse the existing accessToken
func UserRoleMiddleware(atLeastUserRole enums.UserRole) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		currentUserRoleValue, exists := ctx.Get(types.ContextFieldName_User_Role.String())
		if !exists {
			exceptions.Context.MissPlacingOrWrongMiddlewareOrder(
				"Cannot find the userRole, " +
					"please make sure the AuthMiddleware() is placing before the UserRoleMiddleware()",
			).Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		currentUserRole, ok := currentUserRoleValue.(enums.UserRole)
		if !ok {
			exceptions.User.InvalidType("the userRole is not in the correct enum type").Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}

		// iterate the AllUserRole from the highest permission to the lowest
		// if we find the atLeastUserRole first, then the currentUserRole is under the atLeastUserRole
		// 	=> the current user does have access to do the following
		// else if we find the currentUserRole first, then the atLeastUserRole is under it
		//  => the current user doest not have access to do the following
		// else if they are the same, then we just pass the below iteration check
		if currentUserRole == atLeastUserRole {
			ctx.Next()
			return
		}
		// from high level roles to low level roles
		for _, enum := range enums.AllUserRoles {
			if enum == currentUserRole {
				ctx.Next()
				return
			} else if enum == atLeastUserRole {
				exceptions.Auth.PermissionDeniedDueToUserRole(currentUserRole).Log().SafelyAbortAndResponseWithJSON(ctx)
				return
			}
		}

		// if some how we can't find the userDataCache.Role or atLeastUserRole
		// then we raise an undefined error at the end
		exceptions.UndefinedError(
			"Cannot find atLeastUserRole or userDataCache.Role in UserRoleMiddleware",
		).Log().SafelyAbortAndResponseWithJSON(ctx)
	}
}

/*
A Middleware to indicate which type of UserRole can have access to the following operation,

Args:
  - allowedRoles []enums.UserRole : if the current user has the user role in this arguments, this middleware will pass, else it won't

Note: If the allowedRoles is empty, all types of the UserRole will pass
*/
func AllowedUserRolesMiddleware(allowedRoles []enums.UserRole) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		currentUserRoleValue, exists := ctx.Get(types.ContextFieldName_User_Role.String())
		if !exists {
			exception := exceptions.Context.MissPlacingOrWrongMiddlewareOrder(
				"Cannot find the userRole, " +
					"please make sure the AuthMiddleware() is placing before the UserRoleMiddleware()",
			)
			ctx.AbortWithStatusJSON(exception.HTTPStatusCode, exception.GetGinH())
			return
		}
		currentUserRole, ok := currentUserRoleValue.(enums.UserRole)
		if !ok {
			exception := exceptions.User.InvalidType("the userRole is not in the correct enum type")
			ctx.AbortWithStatusJSON(exception.HTTPStatusCode, exception.GetGinH())
			return
		}

		if len(allowedRoles) == 0 {
			ctx.Next()
			return
		}
		for _, enum := range allowedRoles {
			if enum == currentUserRole {
				ctx.Next()
				return
			}
		}

		exception := exceptions.Auth.PermissionDeniedDueToUserRole(currentUserRole)
		ctx.AbortWithStatusJSON(
			exception.HTTPStatusCode,
			exception.GetGinH(),
		)
	}
}
