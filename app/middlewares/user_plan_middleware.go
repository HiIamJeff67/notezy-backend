package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"

	exceptions "notezy-backend/app/exceptions"
	enums "notezy-backend/app/models/enums"
)

// This UserPlanMiddleware() MUST be processed AFTER the AuthMiddleware()
// so that it can parse the existing accessToken
func UserPlanMiddleware(atLeastUserPlan enums.UserPlan) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		currentUserPlanValue, exists := ctx.Get("userPlan")
		if !exists {
			exception := exceptions.Auth.MissPlacingOrWrongMiddlewareOrder(
				"Does not find the userPlan, " +
					"please make sure the AuthMiddleware() is placing before the UserPlanMiddleware()",
			)
			ctx.AbortWithStatusJSON(exception.HTTPStatusCode, exception.GetGinH())
			return
		}
		currentUserPlan, ok := currentUserPlanValue.(enums.UserPlan)
		if !ok {
			exception := exceptions.User.InvalidType("the userPlan is not in the correct enum type")
			ctx.AbortWithStatusJSON(exception.HTTPStatusCode, exception.GetGinH())
			return
		}

		// iterate the AllUserRole from the highest permission to the lowest
		// if we find the atLeastUserRole first, then the currentUserPlan is under the atLeastUserRole
		// 	=> the current user does have access to do the following
		// else if we find the currentUserPlan first, then the atLeastUserRole is under it
		//  => the current user doest not have access to do the following
		// else if they are the same, then we just pass the below iteration check
		if currentUserPlan == atLeastUserPlan {
			ctx.Set("userPlan", currentUserPlan)
			ctx.Next()
			return
		}
		for _, enum := range enums.AllUserPlans {
			if enum == atLeastUserPlan {
				ctx.Set("userPlan", currentUserPlan)
				ctx.Next()
				return
			} else if enum == currentUserPlan {
				ctx.AbortWithStatusJSON(
					http.StatusUnauthorized,
					exceptions.Auth.PermissionDeniedDueToUserPlan(currentUserPlan).GetGinH(),
				)
				return
			}
		}

		// if some how we can't find the currentUserPlan or atLeastUserPlan
		// then we raise an undefined error at the end
		exception := exceptions.Auth.UndefinedError(
			"Cannot find atLeastUserPlan or currentUserPlan in UserRoleMiddleware",
		)
		ctx.AbortWithStatusJSON(
			exception.HTTPStatusCode,
			exception.GetGinH(),
		)
	}
}

/*
A Middleware to indicate which type of UserPlan can have access to the following operation,

Args:
  - allowedPlans []enums.UserPlan : if the current user has the user plan in this arguments, this middleware will pass, else it won't

Note: If the allowedPlans is empty, all types of the UserPlan will pass
*/
func AllowedUserPlanMiddleware(allowedPlan []enums.UserPlan) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		currentUserPlanValue, exists := ctx.Get("userPlan")
		if !exists {
			exception := exceptions.Auth.MissPlacingOrWrongMiddlewareOrder(
				"Does not find the userPlan, " +
					"please make sure the AuthMiddleware() is placing before the AllowedUserPlanMiddleware()",
			)
			ctx.AbortWithStatusJSON(exception.HTTPStatusCode, exception.GetGinH())
			return
		}
		currentUserPlan, ok := currentUserPlanValue.(enums.UserPlan)
		if !ok {
			exception := exceptions.User.InvalidType("the userPlan is not in the correct enum type")
			ctx.AbortWithStatusJSON(exception.HTTPStatusCode, exception.GetGinH())
			return
		}

		if len(allowedPlan) == 0 {
			ctx.Set("userPlan", currentUserPlan)
			ctx.Next()
			return
		}
		for _, enum := range allowedPlan {
			if enum == currentUserPlan {
				ctx.Set("userPlan", currentUserPlan)
				ctx.Next()
				return
			}
		}

		exception := exceptions.Auth.PermissionDeniedDueToUserRole(currentUserPlan)
		ctx.AbortWithStatusJSON(
			exception.HTTPStatusCode,
			exception.GetGinH(),
		)
	}
}
