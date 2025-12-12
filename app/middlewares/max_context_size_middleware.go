package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"

	exceptions "notezy-backend/app/exceptions"
	types "notezy-backend/shared/types"
)

func MaxContextSizeMiddleware(limitBytes int64, unit types.ByteType) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.Request.ContentLength > limitBytes*int64(unit) {
			exceptions.Context.MaxContextBodySizeExceeded(ctx.Request.ContentLength, limitBytes*unit.ToInt64()).
				SafelyAbortAndResponseWithJSON(ctx)
			return
		}

		ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, limitBytes*int64(unit))
		ctx.Next()
	}
}
