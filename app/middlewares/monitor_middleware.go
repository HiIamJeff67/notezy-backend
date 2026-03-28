package middlewares

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"

	exceptions "notezy-backend/app/exceptions"
)

func WithTracerMiddleware(tracer trace.Tracer, spanName string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		newCtx, span := tracer.Start(ctx.Request.Context(), spanName)
		defer span.End()
		ctx.Request = ctx.Request.WithContext(newCtx)
		ctx.Next()
	}
}

func WithMeterMiddleware(meter metric.Meter) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestCounter, err := meter.Int64Counter("http.server.requests.total")
		if err != nil {
			exceptions.Monitor.FailedToInitializeRequestCounter().
				Log().SafelyAbortAndResponseWithJSON(ctx)
		}
		requestCounter.Add(ctx, 1)
		ctx.Next()
	}
}
