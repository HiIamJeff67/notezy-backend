package contexts

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func GetRealClientIP(ctx *gin.Context) string {
	if xff := ctx.GetHeader("X-Forwarded-For"); xff != "" {
		if commaIndex := strings.Index(xff, ","); commaIndex > 0 {
			return strings.TrimSpace(xff[:commaIndex])
		}
		return strings.TrimSpace(xff)
	}

	if xri := ctx.GetHeader("X-Real-IP"); xri != "" {
		return strings.TrimSpace(xri)
	}

	if cfip := ctx.GetHeader("CF-Connecting-IP"); cfip != "" {
		return strings.TrimSpace(cfip)
	}

	return ctx.ClientIP()
}
