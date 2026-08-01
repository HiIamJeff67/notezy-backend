package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	exceptions "github.com/HiIamJeff67/notezy-backend/internal/exceptions"
	cookies "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/api/cookies"
	sharedtokens "github.com/HiIamJeff67/notezy-backend/internal/shared/tokens"
	types "github.com/HiIamJeff67/notezy-backend/internal/shared/types"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set(types.ContextFieldName_User_Id.String(), nil)
		ctx.Set(types.ContextFieldName_User_PublicId.String(), nil)
		ctx.Set(types.ContextFieldName_User_Name.String(), nil)
		ctx.Set(types.ContextFieldName_User_Email.String(), nil)
		ctx.Set(types.ContextFieldName_AccessToken.String(), nil)
		ctx.Set(types.ContextFieldName_CSRFToken.String(), nil)
		ctx.Set(types.ContextFieldName_IsNewTokens.String(), false)

		accessToken, _ := extractAccessToken(ctx)
		if accessToken != "" {
			if claims, err := sharedtokens.ParseAccessToken(accessToken); err == nil {
				if _, err := uuid.Parse(claims.Subject); err == nil {
					ctx.Set(types.ContextFieldName_User_PublicId.String(), claims.Subject)
					ctx.Set(types.ContextFieldName_User_Name.String(), claims.Name)
					ctx.Set(types.ContextFieldName_User_Email.String(), claims.Email)
					ctx.Set(types.ContextFieldName_AccessToken.String(), accessToken)
					ctx.Next()
					return
				}
			}
		}

		refreshToken, _ := extractRefreshToken(ctx)
		if refreshToken != "" {
			if claims, err := sharedtokens.ParseRefreshToken(refreshToken); err == nil {
				if _, err := uuid.Parse(claims.Subject); err == nil {
					ctx.Set(types.ContextFieldName_User_PublicId.String(), claims.Subject)
					ctx.Set(types.ContextFieldName_User_Name.String(), claims.Name)
					ctx.Set(types.ContextFieldName_User_Email.String(), claims.Email)
				}
			}
		}

		ctx.Next()
	}
}

func extractAccessToken(ctx *gin.Context) (string, *exceptions.Exception) {
	accessToken, exception := cookies.AccessTokenCookieHandler.Get(ctx)
	if exception == nil && strings.TrimSpace(accessToken) != "" {
		return accessToken, nil
	}

	authorizationHeader := ctx.GetHeader("Authorization")
	if !strings.HasPrefix(authorizationHeader, "Bearer ") {
		return "", exceptions.New(
			"InvalidAccessToken",
			"Token",
			"ExtractAccessToken",
			"access token is missing or invalid",
			http.StatusUnauthorized,
		)
	}

	return strings.TrimPrefix(authorizationHeader, "Bearer "), nil
}

func extractRefreshToken(ctx *gin.Context) (string, *exceptions.Exception) {
	refreshToken, exception := cookies.RefreshTokenCookieHandler.Get(ctx)
	if exception != nil || strings.TrimSpace(refreshToken) == "" {
		return "", exceptions.New(
			"InvalidRefreshToken",
			"Token",
			"ExtractRefreshToken",
			"refresh token is missing or invalid",
			http.StatusUnauthorized,
		)
	}

	return refreshToken, nil
}
