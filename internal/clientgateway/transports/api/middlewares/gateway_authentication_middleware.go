package middlewares

import (
	"github.com/gin-gonic/gin"

	cookies "github.com/HiIamJeff67/notegic-backend/shared/cookies"
)

func GatewayAuthenticationMiddleware(accessTokenCookieHandler, refreshTokenCookieHandler *cookies.CookieHandler) gin.HandlerFunc {
	return JWTMiddleware(accessTokenCookieHandler, refreshTokenCookieHandler)
}
