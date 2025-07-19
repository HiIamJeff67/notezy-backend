package middlewares

import "github.com/gin-gonic/gin"

func CORSMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		origin := ctx.GetHeader("Origin")

		if origin != "" {
			ctx.Header("Access-Control-Allow-Origin", origin)
		}
		ctx.Header("Access-Control-Allow-Credentials", "true")
		ctx.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		ctx.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, User-Agent, X-Requested-With")
		ctx.Header("Access-Control-Max-Age", "86400") // 24 hours

		if ctx.Request.Method == "OPTIONS" {
			ctx.Status(200)
			return
		}

		ctx.Next()
	}
}
