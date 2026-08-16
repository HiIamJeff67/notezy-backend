package middlewares

import (
	"go.opentelemetry.io/otel/attribute"

	"github.com/gin-gonic/gin"

	metrics "github.com/HiIamJeff67/notegic-backend/shared/platform/observability/metrics"
	traces "github.com/HiIamJeff67/notegic-backend/shared/platform/observability/traces"
)

func ApplyTracerMiddleware(spanName string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		newCtx, span := traces.NotegicTracer.Start(ctx.Request.Context(), "http."+spanName)
		span.SetAttributes(
			attribute.String("http.request.method", ctx.Request.Method),
			attribute.String("http.route", ctx.FullPath()),
			attribute.String("gateway.surface", "client-gateway"),
			attribute.String("gateway.auth_method", "jwt"),
			attribute.String("gateway.operation", spanName),
		)
		defer func() {
			span.SetAttributes(attribute.Int("http.response.status_code", ctx.Writer.Status()))
			traces.NotegicTracer.End(span, nil)
		}()

		ctx.Request = ctx.Request.WithContext(newCtx)
		ctx.Next()
	}
}

func ApplyMeterMiddleware(names ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
		status := ctx.Writer.Status()
		isTotalCounted := false
		for _, name := range names {
			if name == "server.requests.total" {
				isTotalCounted = true
			}
			metrics.NotegicMeter.Count(ctx, name, 1,
				attribute.String("gateway.surface", "client-gateway"),
				attribute.String("gateway.auth_method", "jwt"),
				attribute.String("gateway.operation", name),
				attribute.Int("http.response.status_code", status),
			)
		}
		if !isTotalCounted {
			metrics.NotegicMeter.Count(ctx, "server.requests.total", 1,
				attribute.String("gateway.surface", "client-gateway"),
				attribute.String("gateway.auth_method", "jwt"),
				attribute.String("gateway.operation", "server.requests.total"),
				attribute.Int("http.response.status_code", status),
			)
		}
	}
}
