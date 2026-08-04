package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	exceptions "github.com/HiIamJeff67/notezy-backend/internal/exceptions"
	cookies "github.com/HiIamJeff67/notezy-backend/shared/cookies"
	sharedcontexts "github.com/HiIamJeff67/notezy-backend/shared/lib/contexts"
	sharedtokens "github.com/HiIamJeff67/notezy-backend/shared/tokens"
)

func JWTMiddleware(accessTokenCookieHandler, refreshTokenCookieHandler *cookies.CookieHandler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set(sharedcontexts.ContextFieldName_User_Id.String(), nil)
		ctx.Set(sharedcontexts.ContextFieldName_User_PublicId.String(), nil)
		ctx.Set(sharedcontexts.ContextFieldName_User_Name.String(), nil)
		ctx.Set(sharedcontexts.ContextFieldName_User_Email.String(), nil)
		ctx.Set(sharedcontexts.ContextFieldName_AccessToken.String(), nil)
		ctx.Set(sharedcontexts.ContextFieldName_CSRFToken.String(), nil)
		ctx.Set(sharedcontexts.ContextFieldName_IsNewTokens.String(), false)

		accessToken, _ := extractAccessToken(ctx, accessTokenCookieHandler)
		if accessToken != "" {
			if claims, err := sharedtokens.ParseAccessToken(accessToken); err == nil {
				if _, err := uuid.Parse(claims.Subject); err == nil {
					ctx.Set(sharedcontexts.ContextFieldName_User_PublicId.String(), claims.Subject)
					ctx.Set(sharedcontexts.ContextFieldName_User_Name.String(), claims.Name)
					ctx.Set(sharedcontexts.ContextFieldName_User_Email.String(), claims.Email)
					ctx.Set(sharedcontexts.ContextFieldName_AccessToken.String(), accessToken)
					ctx.Next()
					return
				}
			}
		}

		refreshToken, _ := extractRefreshToken(ctx, refreshTokenCookieHandler)
		if refreshToken != "" {
			if claims, err := sharedtokens.ParseRefreshToken(refreshToken); err == nil {
				if _, err := uuid.Parse(claims.Subject); err == nil {
					ctx.Set(sharedcontexts.ContextFieldName_User_PublicId.String(), claims.Subject)
					ctx.Set(sharedcontexts.ContextFieldName_User_Name.String(), claims.Name)
					ctx.Set(sharedcontexts.ContextFieldName_User_Email.String(), claims.Email)
				}
			}
		}

		ctx.Next()
	}
}

func extractAccessToken(ctx *gin.Context, accessTokenCookieHandler *cookies.CookieHandler) (string, *exceptions.Exception) {
	accessToken, err := accessTokenCookieHandler.Get(ctx)
	if err == nil && strings.TrimSpace(accessToken) != "" {
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

func extractRefreshToken(ctx *gin.Context, refreshTokenCookieHandler *cookies.CookieHandler) (string, *exceptions.Exception) {
	refreshToken, err := refreshTokenCookieHandler.Get(ctx)
	if err != nil || strings.TrimSpace(refreshToken) == "" {
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
