package middlewares

import (
	"strings"

	"github.com/gin-gonic/gin"

	caches "notezy-backend/app/caches"
	cookies "notezy-backend/app/cookies"
	exceptions "notezy-backend/app/exceptions"
	repositories "notezy-backend/app/models/repositories"
	schemas "notezy-backend/app/models/schemas"
	logs "notezy-backend/app/monitor/logs"
	metrics "notezy-backend/app/monitor/metrics"
	traces "notezy-backend/app/monitor/traces"
	tokens "notezy-backend/app/tokens"
	types "notezy-backend/shared/types"
)

func _extractAccessToken(ctx *gin.Context) (string, *exceptions.Exception) {
	accessToken, exception := cookies.AccessTokenCookieHandler.Get(ctx)
	if exception != nil || len(strings.ReplaceAll(accessToken, " ", "")) == 0 {
		authHeader := ctx.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return "", exceptions.Token.FailedToExtractOrValidateAccessToken()
		}
		accessToken = strings.TrimPrefix(authHeader, "Bearer ")
	}
	return accessToken, nil
}

func _extractRefreshToken(ctx *gin.Context) (string, *exceptions.Exception) {
	refreshToken, exception := cookies.RefreshTokenCookieHandler.Get(ctx)
	if exception != nil || strings.ReplaceAll(refreshToken, " ", "") == "" {
		return "", exceptions.Token.FailedToExtractOrValidateRefreshToken()
	}
	return refreshToken, nil
}

func _validateAccessTokenAndUserAgent(accessToken string) (*types.JWTClaims, *caches.UserDataCache, *exceptions.Exception) {
	claims, exception := tokens.ParseAccessToken(accessToken)
	if exception != nil { // if failed to parse the accessToken
		return nil, nil, exception
	}

	userDataCache, exception := caches.GetUserDataCache(claims.Name)
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

	userRepository := repositories.NewUserRepository()
	user, exception := userRepository.GetOneByName(
		claims.Name,
		[]schemas.UserRelation{
			schemas.UserRelation_UserInfo,
			schemas.UserRelation_UserSetting,
		})
	if exception != nil { // if there's not such user with the parsed id
		return nil, exception
	}

	if refreshToken != user.RefreshToken { // if failed to compare and validate the refreshToken as the correct token storing in the database
		logs.FInfo(traces.GetTrace(0).FileLineString(), "current: %s, backend: %s", refreshToken, user.RefreshToken)
		return nil, exceptions.Auth.WrongRefreshToken()
	}

	return user, nil
}

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// clear all the previous field first for security
		ctx.Set(types.ContextFieldName_User_Id.String(), nil)
		ctx.Set(types.ContextFieldName_User_PublicId.String(), nil)
		ctx.Set(types.ContextFieldName_User_Name.String(), nil)
		ctx.Set(types.ContextFieldName_User_DisplayName.String(), nil)
		ctx.Set(types.ContextFieldName_User_Email.String(), nil)
		ctx.Set(types.ContextFieldName_IsNewTokens.String(), nil)
		ctx.Set(types.ContextFieldName_AccessToken.String(), nil)
		ctx.Set(types.ContextFieldName_User_Role.String(), nil)
		ctx.Set(types.ContextFieldName_User_Plan.String(), nil)

		// nest if statement bcs we will skip the accessToken validation if it failed
		if accessToken, exception := _extractAccessToken(ctx); exception == nil { // if extract the accessToken successfully
			if claims, userDataCache, exception := _validateAccessTokenAndUserAgent(accessToken); exception == nil { // if validate the accessToken successfully
				if currentUserAgent := ctx.GetHeader("User-Agent"); currentUserAgent == claims.UserAgent { // if the userAgent is matched
					// if everything above is all fine, we should get the valid userDataCache and claims
					ctx.Set(types.ContextFieldName_User_Id.String(), userDataCache.Id.String()) // remain that all the context values should be type of string
					ctx.Set(types.ContextFieldName_User_PublicId.String(), userDataCache.PublicId)
					ctx.Set(types.ContextFieldName_User_Name.String(), userDataCache.Name)
					ctx.Set(types.ContextFieldName_User_DisplayName.String(), userDataCache.DisplayName)
					ctx.Set(types.ContextFieldName_User_Email.String(), userDataCache.Email)
					ctx.Set(types.ContextFieldName_IsNewTokens.String(), false)
					ctx.Set(types.ContextFieldName_AccessToken.String(), accessToken)
					ctx.Set(types.ContextFieldName_User_Role.String(), userDataCache.Role)
					ctx.Set(types.ContextFieldName_User_Plan.String(), userDataCache.Plan)
					// also extend the ttl of the user data cache
					if exception := caches.ExtendUserDataCacheTTL(userDataCache.Name); exception != nil {
						exception.Log()
					}
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
			exception.Log().SafelyAbortAndResponseWithJSON(ctx, metrics.MetricNames.Server.Responses.Failed.Unauthorized)
			return
		}

		_user, exception := _validateRefreshToken(refreshToken)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx, metrics.MetricNames.Server.Responses.Failed.Unauthorized)
			return
		}

		// if we can't check the userAgent in accessToken, then we check it in our database
		if currentUserAgent := ctx.GetHeader("User-Agent"); currentUserAgent != _user.UserAgent {
			exceptions.Auth.WrongUserAgent().Log().SafelyAbortAndResponseWithJSON(ctx, metrics.MetricNames.Server.Responses.Failed.Unauthorized)
			return
		}

		// if we failed to validate the accessToken, but we have validated the refreshToken
		// then we need to generate the new accessToken, and storing it in the cache, and regarding the entire validation as successful
		newAccessToken, exception := tokens.GenerateAccessToken(_user.Name, _user.Email, _user.UserAgent)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx, metrics.MetricNames.Server.Responses.Failed.Unauthorized)
			return
		}
		// also generate a new CSRF token, and storing it in the cache
		newCSRFToken, exception := tokens.GenerateCSRFToken()
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx, metrics.MetricNames.Server.Responses.Failed.Unauthorized)
			return
		}

		// at this stage, make sure we update the cache of the user data
		exception = caches.UpdateUserDataCache(
			_user.Name,
			caches.UpdateUserDataCacheDto{
				AccessToken: newAccessToken,
				CSRFToken:   newCSRFToken,
			},
		)
		if exception != nil {
			exception.WithDetails("trying to set the new user data instead").Log()
			newUserDataCache := caches.UserDataCache{
				Id:                 _user.Id,
				PublicId:           _user.PublicId,
				Name:               _user.Name,
				DisplayName:        _user.DisplayName,
				Email:              _user.Email,
				AccessToken:        *newAccessToken,
				CSRFToken:          *newCSRFToken,
				Role:               _user.Role,
				Plan:               _user.Plan,
				Status:             _user.Status,
				Language:           _user.UserSetting.Language,
				GeneralSettingCode: _user.UserSetting.GeneralSettingCode,
				PrivacySettingCode: _user.UserSetting.PrivacySettingCode,
				CreatedAt:          _user.CreatedAt,
				UpdatedAt:          _user.UpdatedAt,
			}
			if _user.UserInfo.AvatarURL != nil {
				newUserDataCache.AvatarURL = *_user.UserInfo.AvatarURL
			}

			if exception = caches.SetUserDataCache(_user.Name, newUserDataCache); exception != nil {
				exception.Log().SafelyAbortAndResponseWithJSON(ctx, metrics.MetricNames.Server.Responses.Failed.Unauthorized)
				return
			}
		} else {
			if exception := caches.ExtendUserDataCacheTTL(_user.Name); exception != nil {
				exception.Log()
			}
		}

		ctx.Set(types.ContextFieldName_User_Id.String(), _user.Id.String())
		ctx.Set(types.ContextFieldName_User_PublicId.String(), _user.PublicId)
		ctx.Set(types.ContextFieldName_User_Name.String(), _user.Name)
		ctx.Set(types.ContextFieldName_User_DisplayName.String(), _user.DisplayName)
		ctx.Set(types.ContextFieldName_User_Email.String(), _user.Email)
		ctx.Set(types.ContextFieldName_IsNewTokens.String(), true)
		ctx.Set(types.ContextFieldName_AccessToken.String(), *newAccessToken)
		ctx.Set(types.ContextFieldName_CSRFToken.String(), *newCSRFToken)
		ctx.Set(types.ContextFieldName_User_Role.String(), _user.Role)
		ctx.Set(types.ContextFieldName_User_Plan.String(), _user.Plan)
		ctx.Next()
	}
}
