package middlewares

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	exceptions "github.com/HiIamJeff67/notezy-backend/internal/exceptions"
	responsewriter "github.com/HiIamJeff67/notezy-backend/internal/shared/responsewriter"
	types "github.com/HiIamJeff67/notezy-backend/internal/shared/types"
)

func MaxContextSizeMiddleware(limitBytes int64, unit types.ByteType) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.Request.ContentLength > limitBytes*int64(unit) {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.New(
				"MaxContextBodySizeExceeded",
				"Context",
				"Validate",
				fmt.Sprintf("The request body size of %d bytes exceeds the maximum of %d bytes", ctx.Request.ContentLength, limitBytes*unit.ToInt64()),
				http.StatusRequestEntityTooLarge,
			), ctx)
			return
		}

		ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, limitBytes*int64(unit))
		ctx.Next()
	}
}
