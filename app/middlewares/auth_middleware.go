package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	caches "notezy-backend/app/caches"
	cookies "notezy-backend/app/cookies"
	exceptions "notezy-backend/app/exceptions"
	repositories "notezy-backend/app/models/repositories"
	schemas "notezy-backend/app/models/schemas"
	tokens "notezy-backend/app/tokens"
	types "notezy-backend/shared/types"
)

func _extractAccessToken(ctx *gin.Context) (string, *exceptions.Exception) {
	accessToken, exception := cookies.AccessToken.GetCookie(ctx)
	if exception != nil || len(strings.ReplaceAll(accessToken, " ", "")) == 0 {
		authHeader := ctx.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return "", exceptions.Auth.FailedToExtractOrValidateAccessToken()
		}
		accessToken = strings.TrimPrefix(authHeader, "Bearer ")
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

func _validateAccessTokenAndUserAgent(accessToken string) (*types.Claims, *caches.UserDataCache, *exceptions.Exception) {
	claims, exception := tokens.ParseAccessToken(accessToken)
	if exception != nil { // if failed to parse the accessToken
		return nil, nil, exception
	}

	userId, err := uuid.Parse(claims.Id)
	if err != nil { // if the id is invalid somehow
		return nil, nil, exceptions.Util.FailedToParseAccessToken().WithError(err)
	}

	userDataCache, exception := caches.GetUserDataCache(userId)
	if exception != nil { // if there's no user cache storing its accessToken, in this way, we're impossible to validate its accessToken
		return nil, nil, exception.Log()
	}

	if accessToken != userDataCache.AccessToken { // if failed to compare and validate the accessToken as the correct token storing in the cache
		return nil, nil, exceptions.Auth.WrongAccessToken()
	}

	return claims, userDataCache, nil
}

func _validateRefreshToken(refreshToken string) (*schemas.User, *exceptions.Exception) {
	claims, exception := tokens.ParseRefreshToken(refreshToken)
	if exception != nil { // if failed to parse the refreshToken
		return nil, exception
	}

	userId, err := uuid.Parse(claims.Id)
	if err != nil { // if the id is invalid somehow
		return nil, exceptions.Util.FailedToParseAccessToken().WithError(err)
	}

	userRepository := repositories.NewUserRepository(nil)
	user, exception := userRepository.GetOneById(userId, nil)
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
		// clear all the previous field first for security
		ctx.Set("userId", "")
		ctx.Set("publicId", "")
		ctx.Set("name", "")
		ctx.Set("displayName", "")
		ctx.Set("email", "")
		ctx.Set("accessToken", "")

		// nest if statement bcs we will skip the accessToken validation if it failed
		if accessToken, exception := _extractAccessToken(ctx); exception == nil { // if extract the accessToken successfully
			if claims, userDataCache, exception := _validateAccessTokenAndUserAgent(accessToken); exception == nil { // if validate the accessToken successfully
				if currentUserAgent := ctx.GetHeader("User-Agent"); currentUserAgent == claims.UserAgent { // if the userAgent is matched
					// if everything above is all fine, we should get the valid userDataCache and claims
					ctx.Set("userId", claims.Id)
					ctx.Set("publicId", userDataCache.PublicId)
					ctx.Set("name", userDataCache.Name)
					ctx.Set("displayName", userDataCache.DisplayName)
					ctx.Set("email", userDataCache.Email)
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
			exception.Log()
			ctx.AbortWithStatusJSON(exception.HTTPStatusCode, exception.GetGinH())
			return
		}

		_user, exception := _validateRefreshToken(refreshToken)
		if exception != nil {
			exception.Log()
			ctx.AbortWithStatusJSON(exception.HTTPStatusCode, exception.GetGinH())
			return
		}

		// if we can't check the userAgent in accessToken, then we check it in our database
		if currentUserAgent := ctx.GetHeader("User-Agent"); currentUserAgent != _user.UserAgent {
			exception.Log()
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, exceptions.Auth.WrongUserAgent().GetGinH())
			return
		}

		// if we failed to validate the accessToken, but we have validated the refreshToken
		// then we need to generate the new accessToken, and storing it in the cache, and regarding the entire validation as successful
		newAccessToken, exception := tokens.GenerateAccessToken(_user.Id.String(), _user.Name, _user.Email, _user.UserAgent)
		if exception != nil {
			exception.Log()
			ctx.AbortWithStatusJSON(exception.HTTPStatusCode, exception.GetGinH())
			return
		}

		// at this stage, make sure we update the cache of the user data
		exception = caches.UpdateUserDataCache(_user.Id, caches.UpdateUserDataCacheDto{AccessToken: newAccessToken})
		if exception != nil {
			exception.Log()
			ctx.AbortWithStatusJSON(exception.HTTPStatusCode, exception.GetGinH())
		}

		ctx.Set("userId", _user.Id.String())
		ctx.Set("publicId", _user.PublicId)
		ctx.Set("name", _user.Name)
		ctx.Set("displayName", _user.DisplayName)
		ctx.Set("email", _user.Email)
		ctx.Set("accessToken", newAccessToken)
		ctx.Next()
	}
}
