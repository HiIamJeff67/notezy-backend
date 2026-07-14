package middlewares

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"

	metrics "github.com/HiIamJeff67/notezy-backend/app/monitor/metrics"
	traces "github.com/HiIamJeff67/notezy-backend/app/monitor/traces"
)

func ApplyTracerMiddleware(spanName string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		newCtx, span := traces.NotezyTracer.Start(ctx.Request.Context(), "http."+spanName)
		span.SetAttributes(
			attribute.String("http.request.method", ctx.Request.Method),
			attribute.String("http.route", ctx.FullPath()),
		)
		defer func() {
			span.SetAttributes(attribute.Int("http.response.status_code", ctx.Writer.Status()))
			traces.NotezyTracer.End(span, nil)
		}()

		ctx.Request = ctx.Request.WithContext(newCtx)
		ctx.Next()
	}
}

func ApplyMeterMiddleware(names ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		isTotalCounted := false
		for _, name := range names {
			if name == "server.requests.total" {
				isTotalCounted = true
			}
			metrics.NotezyMeter.Count(ctx, name, 1)
		}
		if !isTotalCounted {
			metrics.NotezyMeter.Count(ctx, "server.requests.total", 1)
		}
		ctx.Next()
	}
}
