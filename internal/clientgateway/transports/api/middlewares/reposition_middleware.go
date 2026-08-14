package middlewares

import "github.com/gin-gonic/gin"

func Reposition(
	fronts []gin.HandlerFunc,
	middles []gin.HandlerFunc,
	backs ...gin.HandlerFunc,
) []gin.HandlerFunc {
	handlers := make([]gin.HandlerFunc, 0, len(fronts)+len(middles)+len(backs))
	return append(
		append(
			append(handlers, fronts...),
			middles...,
		), backs...,
	)
}

// RepositionMiddleware is kept as a compatibility alias while route files
// migrate to the shorter Reposition name.
func RepositionMiddleware(
	fronts []gin.HandlerFunc,
	middles []gin.HandlerFunc,
	backs ...gin.HandlerFunc,
) []gin.HandlerFunc {
	return Reposition(fronts, middles, backs...)
}
