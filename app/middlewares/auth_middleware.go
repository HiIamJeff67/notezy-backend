package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	caches "notezy-backend/app/caches"
	cookies "notezy-backend/app/cookie"
	exceptions "notezy-backend/app/exceptions"
	operations "notezy-backend/app/models/operations"
	schemas "notezy-backend/app/models/schemas"
	util "notezy-backend/app/util"
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

func _validateAccessToken(accessToken string) *exceptions.Exception {
	claims, exception := util.ParseAccessToken(accessToken)
	if exception != nil { // if failed to parse the accessToken
		return exception
	}

	userId, err := uuid.Parse(claims.Id)
	if err != nil { // if the id is invalid somehow
		return exceptions.Util.FailedToParseAccessToken().WithError(err)
	}

	userDataCache, exception := caches.GetUserDataCache(userId)
	if exception != nil { // if there's no user cache storing its accessToken, in this way, we're impossible to validate its accessToken
		return exception
	}

	if accessToken != userDataCache.AccessToken { // if failed to compare and validate the accessToken as the correct token storing in the cache
		return exceptions.Auth.WrongAccessToken()
	}

	return nil
}

func _validateRefreshToken(refreshToken string) (*schemas.User, *exceptions.Exception) {
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

	return user, nil
}

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if accessToken, exception := _extractAccessToken(ctx); exception == nil { // if extract the accessToken successfully
			if exception = _validateAccessToken(accessToken); exception == nil { // if validate the accessToken successfully
				ctx.Set("accessToken", accessToken)
				ctx.Next()
				return
			}
		}

		// if the above procedures to validating accessToken is failed,
		// we now try to extract and validate the refreshToken
		refreshToken, exception := _extractRefreshToken(ctx)
		if exception != nil { // if failed to extract the refreshToken
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, exception.GetGinH())
			return
		}

		_user, exception := _validateRefreshToken(refreshToken)
		if exception != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, exception.GetGinH())
			return
		}

		// if we failed to validate the accessToken, but we have validated the refreshToken
		// then we need to generate the new accessToken, and storing it in the cache, and regarding the entire validation as successful
		newAccessToken, exception := util.GenerateAccessToken(_user.Id.String(), _user.Name, _user.Email)
		if exception != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, exception.GetGinH())
			return
		}

		// try to update the cache of the user data
		exception = caches.UpdateUserDataCache(_user.Id, caches.UpdateUserDataCacheDto{AccessToken: newAccessToken})
		if exception != nil {
			exception.Log()
		}

		ctx.Set("accessToken", newAccessToken)
		ctx.Next()
	}
}
