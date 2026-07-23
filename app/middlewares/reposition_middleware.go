package middlewares

import "github.com/gin-gonic/gin"

func RepositionMiddleware(
	fronts []gin.HandlerFunc,
	middles []gin.HandlerFunc,
	handler gin.HandlerFunc,
	backs ...gin.HandlerFunc,
) []gin.HandlerFunc {
	handlers := make([]gin.HandlerFunc, 0, len(fronts)+len(middles)+1+len(backs))
	return append(append(append(append(handlers, fronts...), middles...), handler), backs...)
}
