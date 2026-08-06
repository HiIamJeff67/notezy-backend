package yjsworker

import (
	"go.opentelemetry.io/otel"
	"os"
	"testing"

	logs "github.com/HiIamJeff67/notezy-backend/shared/platform/observability/logs"
	metrics "github.com/HiIamJeff67/notezy-backend/shared/platform/observability/metrics"
	traces "github.com/HiIamJeff67/notezy-backend/shared/platform/observability/traces"
)

func TestMain(m *testing.M) {
	logs.NotezyLogger = logs.NewLogger(true)
	metrics.NotezyMeter = metrics.NewMeter(otel.Meter("realtime.test"))
	traces.NotezyTracer = traces.NewTracer(otel.Tracer("realtime.test"))

	os.Exit(m.Run())
}
