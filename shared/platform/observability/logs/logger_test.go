package logs

import (
	"bytes"
	"context"
	"io"
	stdlog "log"
	"log/slog"
	"strings"
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

func TestLoggerJSONUsesPrettyConsolePayload(t *testing.T) {
	var output bytes.Buffer
	logger := &Logger{
		logger: slog.New(slog.NewJSONHandler(io.Discard, nil)),
		console: &CommandLineInterfaceLogger{
			logger: stdlog.New(&output, "", 0),
		},
	}

	if err := logger.JSON(
		context.Background(),
		slog.LevelError,
		"exception",
		map[string]string{"reason": "NotFound"},
	); err != nil {
		t.Fatalf("marshal JSON payload: %v", err)
	}

	if !strings.Contains(output.String(), "exception\n") || !strings.Contains(output.String(), "\n  \"reason\": \"NotFound\"") {
		t.Fatalf("console payload was not rendered on multiple lines: %q", output.String())
	}
}

func TestConsoleLoggingEnabledPrefersExplicitFormat(t *testing.T) {
	t.Setenv("OTEL_DEPLOYMENT_ENVIRONMENT", "production")
	t.Setenv("NOTEGIC_LOG_FORMAT", "console")
	if !ConsoleLoggingEnabled() {
		t.Fatal("expected explicit console format to enable console logging")
	}

	t.Setenv("NOTEGIC_LOG_FORMAT", "json")
	if ConsoleLoggingEnabled() {
		t.Fatal("expected explicit JSON format to disable console logging")
	}
}
