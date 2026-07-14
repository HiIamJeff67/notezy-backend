package realtime

import (
	"os"
	"testing"

	"go.opentelemetry.io/otel"

	logs "github.com/HiIamJeff67/notezy-backend/app/monitor/logs"
	metrics "github.com/HiIamJeff67/notezy-backend/app/monitor/metrics"
	traces "github.com/HiIamJeff67/notezy-backend/app/monitor/traces"
)

func TestMain(m *testing.M) {
	logs.NotezyLogger = logs.NewLogger(true)
	metrics.NotezyMeter = metrics.NewMeter(otel.Meter("realtime.test"))
	traces.NotezyTracer = traces.NewTracer(otel.Tracer("realtime.test"))

	os.Exit(m.Run())
}
