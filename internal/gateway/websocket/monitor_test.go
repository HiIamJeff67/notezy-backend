package websocket

import (
	"os"
	"testing"

	"go.opentelemetry.io/otel"

	logs "github.com/HiIamJeff67/notezy-backend/internal/platform/observability/logs"
	metrics "github.com/HiIamJeff67/notezy-backend/internal/platform/observability/metrics"
	traces "github.com/HiIamJeff67/notezy-backend/internal/platform/observability/traces"
)

func TestMain(m *testing.M) {
	logs.NotezyLogger = logs.NewLogger(true)
	metrics.NotezyMeter = metrics.NewMeter(otel.Meter("realtime.test"))
	traces.NotezyTracer = traces.NewTracer(otel.Tracer("realtime.test"))

	os.Exit(m.Run())
}
