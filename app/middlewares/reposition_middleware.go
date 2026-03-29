package middlewares

import "github.com/gin-gonic/gin"

func RepositionMiddleware(fronts []gin.HandlerFunc, backs []gin.HandlerFunc, handler gin.HandlerFunc) []gin.HandlerFunc {
	return append(append(fronts, backs...), handler)
}
