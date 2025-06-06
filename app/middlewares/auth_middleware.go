package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	caches "notezy-backend/app/caches"
	cookies "notezy-backend/app/cookies"
	exceptions "notezy-backend/app/exceptions"
	operations "notezy-backend/app/models/operations"
	schemas "notezy-backend/app/models/schemas"
	util "notezy-backend/app/util"
	types "notezy-backend/global/types"
)

func _extractAccessToken(ctx *gin.Context) (string, *exceptions.Exception) {
	accessToken, exception := cookies.AccessToken.GetCookie(ctx)
	if exception != nil || strings.ReplaceAll(accessToken, " ", "") == "" {
		authHeader := ctx.GetHeader("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			accessToken = strings.TrimPrefix(authHeader, "Bearer ")
		}
	}
	if accessToken == "" {
		return "", exceptions.Auth.FailedToExtractOrValidateAccessToken()
	}
	return accessToken, nil
}

func _extractRefreshToken(ctx *gin.Context) (string, *exceptions.Exception) {
	refreshToken, exception := cookies.RefreshToken.GetCookie(ctx)
	if exception != nil || strings.ReplaceAll(refreshToken, " ", "") == "" {
		return "", exceptions.Auth.FailedToExtractOrValidateRefreshToken()
	}
	return refreshToken, nil
}

func _validateAccessTokenAndUserAgent(ctx *gin.Context, accessToken string) (*types.Claims, *exceptions.Exception) {
	claims, exception := util.ParseAccessToken(accessToken)
	if exception != nil { // if failed to parse the accessToken
		return nil, exception
	}

	userId, err := uuid.Parse(claims.Id)
	if err != nil { // if the id is invalid somehow
		return nil, exceptions.Util.FailedToParseAccessToken().WithError(err)
	}

	userDataCache, exception := caches.GetUserDataCache(userId)
	if exception != nil { // if there's no user cache storing its accessToken, in this way, we're impossible to validate its accessToken
		return nil, exception
	}

	if accessToken != userDataCache.AccessToken { // if failed to compare and validate the accessToken as the correct token storing in the cache
		return nil, exceptions.Auth.WrongAccessToken()
	}

	ctx.Set("userRole", userDataCache.Role)
	ctx.Set("userPlan", userDataCache.Plan)
	return claims, nil
}

func _validateRefreshToken(ctx *gin.Context, refreshToken string) (*schemas.User, *exceptions.Exception) {
	claims, exception := util.ParseRefreshToken(refreshToken)
	if exception != nil { // if failed to parse the refreshToken
		return nil, exception
	}

	userId, err := uuid.Parse(claims.Id)
	if err != nil { // if the id is invalid somehow
		return nil, exceptions.Util.FailedToParseAccessToken().WithError(err)
	}

	user, exception := operations.GetUserById(nil, userId)
	if exception != nil { // if there's not such user with the parsed id
		return nil, exception
	}

	if refreshToken != user.RefreshToken { // if failed to compare and validate the refreshToken as the correct token storing in the database
		return nil, exceptions.Auth.WrongRefreshToken()
	}

	ctx.Set("userRole", user.Role)
	ctx.Set("userPlan", user.Plan)
	return user, nil
}

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// nest if statement bcs we will skip the accessToken validation if it failed
		if accessToken, exception := _extractAccessToken(ctx); exception == nil { // if extract the accessToken successfully
			if claims, exception := _validateAccessTokenAndUserAgent(ctx, accessToken); exception == nil { // if validate the accessToken successfully
				if currentUserAgent := ctx.GetHeader("User-Agent"); currentUserAgent == claims.UserAgent { // if the userAgent is matched
					ctx.Set("accessToken", accessToken)
					ctx.Next()
					return
				}
			}
		}

		// if the above procedures to validating accessToken is failed,
		// we now try to extract and validate the refreshToken
		// this means the old accessToken can no longer get any data of the user
		refreshToken, exception := _extractRefreshToken(ctx)
		if exception != nil { // if failed to extract the refreshToken
			ctx.AbortWithStatusJSON(exception.HTTPStatusCode, exception.GetGinH())
			return
		}

		_user, exception := _validateRefreshToken(ctx, refreshToken)
		if exception != nil {
			ctx.AbortWithStatusJSON(exception.HTTPStatusCode, exception.GetGinH())
			return
		}

		// if we can't check the userAgent in accessToken, then we check it in our database
		if currentUserAgent := ctx.GetHeader("User-Agent"); currentUserAgent != _user.UserAgent {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, exceptions.Auth.WrongUserAgent().GetGinH())
			return
		}

		// if we failed to validate the accessToken, but we have validated the refreshToken
		// then we need to generate the new accessToken, and storing it in the cache, and regarding the entire validation as successful
		newAccessToken, exception := util.GenerateAccessToken(_user.Id.String(), _user.Name, _user.Email, _user.UserAgent)
		if exception != nil {
			ctx.AbortWithStatusJSON(exception.HTTPStatusCode, exception.GetGinH())
			return
		}

		// at this stage, make sure we update the cache of the user data
		exception = caches.UpdateUserDataCache(_user.Id, caches.UpdateUserDataCacheDto{AccessToken: newAccessToken})
		if exception != nil {
			ctx.AbortWithStatusJSON(exception.HTTPStatusCode, exception.GetGinH())
		}

		ctx.Set("accessToken", newAccessToken)
		ctx.Next()
	}
}
