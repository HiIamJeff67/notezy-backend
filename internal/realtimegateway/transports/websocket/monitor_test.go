package websocket

import (
	"go.opentelemetry.io/otel"
	"os"
	"testing"

	logs "github.com/HiIamJeff67/notegic-backend/shared/platform/observability/logs"
	metrics "github.com/HiIamJeff67/notegic-backend/shared/platform/observability/metrics"
	traces "github.com/HiIamJeff67/notegic-backend/shared/platform/observability/traces"
)

func TestMain(m *testing.M) {
	logs.NotegicLogger = logs.NewLogger(true)
	metrics.NotegicMeter = metrics.NewMeter(otel.Meter("realtime.test"))
	traces.NotegicTracer = traces.NewTracer(otel.Tracer("realtime.test"))

	os.Exit(m.Run())
}
