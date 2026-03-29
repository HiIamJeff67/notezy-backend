package middlewares

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"

	exceptions "notezy-backend/app/exceptions"
	metrics "notezy-backend/app/monitor/metrics"
)

func ApplyTracerMiddleware(tracer trace.Tracer, spanName string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		newCtx, span := tracer.Start(ctx.Request.Context(), spanName)
		defer span.End()
		ctx.Request = ctx.Request.WithContext(newCtx)
		ctx.Next()
	}
}

// The ApplyMeterMiddleware will accept a meter and the optional field of names,
// then iterate all the names to get the corresponding request counter and increament them in int64.
// The meter is a type of metric.Meter, which should be initialized by calling otel.Meter("service-name"),
// the names is the names of the target request counter.
// Note that the label of 'server.requests.total' will always going to increament once apply this middleware
// even if its name is not passed.
func ApplyMeterMiddleware(meter metric.Meter, names ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		isTotalCounted := false
		for _, name := range names {
			if name == metrics.MetricNames.Server.Requests.Total {
				isTotalCounted = true
			}
			requestCounter, err := meter.Int64Counter(name)
			if err != nil {
				exceptions.Monitor.FailedToInitializeRequestCounter().Log()
			} else {
				requestCounter.Add(ctx, 1)
			}
		}
		if !isTotalCounted {
			requestCounter, err := meter.Int64Counter(metrics.MetricNames.Server.Requests.Total)
			if err != nil {
				exceptions.Monitor.FailedToInitializeRequestCounter().Log()
			} else {
				requestCounter.Add(ctx, 1)
			}
		}
		ctx.Next()
	}
}
