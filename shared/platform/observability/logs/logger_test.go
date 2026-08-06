package logs

import (
	"context"
	"io"
	"log/slog"
	"testing"
)

func TestLoggerJSON(t *testing.T) {
	logger := &Logger{
		logger: slog.New(slog.NewJSONHandler(io.Discard, nil)),
	}
	if err := logger.JSON(
		context.Background(),
		slog.LevelError,
		"log payload",
		map[string]string{
			"reason": "NotFound",
		},
	); err != nil {
		t.Fatalf("marshal JSON payload: %v", err)
	}

	if err := logger.JSON(
		context.Background(),
		slog.LevelError,
		"log payload",
		func() {},
	); err == nil {
		t.Fatal("expected marshal error")
	}
}
