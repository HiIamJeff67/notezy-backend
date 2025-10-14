package middlewares

import (
	"github.com/gin-gonic/gin"

	exceptions "notezy-backend/app/exceptions"
	enums "notezy-backend/app/models/schemas/enums"
	constants "notezy-backend/shared/constants"
)

// This UserPlanMiddleware() MUST be processed AFTER the AuthMiddleware()
// so that it can parse the existing accessToken
func UserPlanMiddleware(atLeastUserPlan enums.UserPlan) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		currentUserPlanValue, exists := ctx.Get(constants.ContextFieldName_User_Plan.String())
		if !exists {
			exceptions.Auth.MissPlacingOrWrongMiddlewareOrder(
				"Cannot find the userPlan, " +
					"please make sure the AuthMiddleware() is placing before the UserPlanMiddleware()",
			).Log().SafelyResponseWithJSON(ctx)
			return
		}
		currentUserPlan, ok := currentUserPlanValue.(enums.UserPlan)
		if !ok {
			exceptions.User.InvalidType("the userPlan is not in the correct enum type").
				Log().
				SafelyResponseWithJSON(ctx)
			return
		}

		// iterate the AllUserRole from the highest permission to the lowest
		// if we find the atLeastUserRole first, then the currentUserPlan is under the atLeastUserRole
		// 	=> the current user does have access to do the following
		// else if we find the currentUserPlan first, then the atLeastUserRole is under it
		//  => the current user doest not have access to do the following
		// else if they are the same, then we just pass the below iteration check
		if currentUserPlan == atLeastUserPlan {
			ctx.Next()
			return
		}
		// from high level plans to low level plans
		for _, enum := range enums.AllUserPlans {
			if enum == currentUserPlan {
				ctx.Next()
				return
			} else if enum == atLeastUserPlan {
				exceptions.Auth.PermissionDeniedDueToUserPlan(currentUserPlan).
					Log().
					SafelyResponseWithJSON(ctx)
				return
			}
		}

		// if some how we can't find the currentUserPlan or atLeastUserPlan
		// then we raise an undefined error at the end
		exceptions.UndefinedError(
			"Cannot find atLeastUserPlan or currentUserPlan in UserRoleMiddleware",
		).Log().SafelyResponseWithJSON(ctx)
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
		currentUserPlanValue, exists := ctx.Get(constants.ContextFieldName_User_Plan.String())
		if !exists {
			exceptions.Auth.MissPlacingOrWrongMiddlewareOrder(
				"Cannot find the userPlan, " +
					"please make sure the AuthMiddleware() is placing before the AllowedUserPlanMiddleware()",
			).Log().SafelyResponseWithJSON(ctx)
			return
		}
		currentUserPlan, ok := currentUserPlanValue.(enums.UserPlan)
		if !ok {
			exceptions.User.InvalidType("the userPlan is not in the correct enum type").Log().SafelyResponseWithJSON(ctx)
			return
		}

		if len(allowedPlan) == 0 {
			ctx.Next()
			return
		}
		for _, enum := range allowedPlan {
			if enum == currentUserPlan {
				ctx.Next()
				return
			}
		}

		exceptions.Auth.PermissionDeniedDueToUserRole(currentUserPlan).Log().SafelyResponseWithJSON(ctx)
	}
}
