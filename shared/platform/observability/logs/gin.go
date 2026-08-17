package logs

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
)

func WithGinLogger(router *gin.Engine) *gin.Engine {
	if ConsoleLoggingEnabled() {
		gin.ForceConsoleColor()
		router.Use(gin.Logger(), gin.Recovery())
	} else {
		gin.DisableConsoleColor()
		router.Use(StructuredAccessLogger(), gin.Recovery())
	}
	return router
}

func StructuredAccessLogger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		startedAt := time.Now()
		ctx.Next()

		if NotegicLogger == nil {
			return
		}

		status := ctx.Writer.Status()
		message := fmt.Sprintf("%s %s", ctx.Request.Method, ctx.Request.URL.Path)
		attributes := []attribute.KeyValue{
			attribute.String("http.request.method", ctx.Request.Method),
			attribute.String("http.route", ctx.FullPath()),
			attribute.Int("http.response.status_code", status),
			attribute.Int64("http.response.duration_ms", time.Since(startedAt).Milliseconds()),
		}
		if ctx.FullPath() == "" {
			attributes[1] = attribute.String("http.route", "<unmatched>")
		}

		switch {
		case status >= http.StatusInternalServerError:
			NotegicLogger.Error(ctx.Request.Context(), nil, message, attributes...)
		case status >= http.StatusBadRequest:
			NotegicLogger.Warn(ctx.Request.Context(), message, attributes...)
		default:
			NotegicLogger.Info(ctx.Request.Context(), message, attributes...)
		}
	}
}
