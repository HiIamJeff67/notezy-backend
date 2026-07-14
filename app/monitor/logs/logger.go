package logs

import (
	"context"
	"fmt"
	stdlog "log"
	"log/slog"
	"os"
	"strings"
	"time"

	colog "github.com/comail/colog"
	"go.opentelemetry.io/otel/attribute"
	otellog "go.opentelemetry.io/otel/log"
	logglobal "go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/trace"

	traces "github.com/HiIamJeff67/notezy-backend/app/monitor/traces"
)

type LoggerInterface interface {
	Debug(ctx context.Context, message string, attributes ...attribute.KeyValue)
	Info(ctx context.Context, message string, attributes ...attribute.KeyValue)
	Warn(ctx context.Context, message string, attributes ...attribute.KeyValue)
	Error(ctx context.Context, err error, message string, attributes ...attribute.KeyValue)
	Alert(ctx context.Context, err error, message string, attributes ...attribute.KeyValue)
}

type Logger struct {
	logger   *slog.Logger
	emitOTel bool
}

type CommandLineInterfaceLogger struct {
	logger *stdlog.Logger
}

func NewLogger(emitOtel bool) LoggerInterface {
	return &Logger{
		logger:   slog.New(slog.NewJSONHandler(os.Stdout, nil)),
		emitOTel: emitOtel,
	}
}

func NewCommandLineInterfaceLogger() LoggerInterface {
	coloredLogger := colog.NewCoLog(os.Stderr, "", stdlog.Ltime)
	coloredLogger.SetFormatter(&colog.StdFormatter{
		Flag:   stdlog.Ltime,
		Colors: true,
	})

	return &CommandLineInterfaceLogger{
		logger: stdlog.New(coloredLogger, "", 0),
	}
}

var (
	NotezyLogger LoggerInterface
)

func (l *Logger) Debug(ctx context.Context, message string, attributes ...attribute.KeyValue) {
	l.write(ctx, slog.LevelDebug, otellog.SeverityDebug, nil, message, attributes...)
}

func (l *Logger) Info(ctx context.Context, message string, attributes ...attribute.KeyValue) {
	l.write(ctx, slog.LevelInfo, otellog.SeverityInfo, nil, message, attributes...)
}

func (l *Logger) Warn(ctx context.Context, message string, attributes ...attribute.KeyValue) {
	l.write(ctx, slog.LevelWarn, otellog.SeverityWarn, nil, message, attributes...)
}

func (l *Logger) Error(ctx context.Context, err error, message string, attributes ...attribute.KeyValue) {
	l.write(ctx, slog.LevelError, otellog.SeverityError, err, message, attributes...)
}

func (l *Logger) Alert(ctx context.Context, err error, message string, attributes ...attribute.KeyValue) {
	l.write(ctx, slog.LevelError, otellog.SeverityFatal, err, message, attributes...)
}

func (l *CommandLineInterfaceLogger) Debug(ctx context.Context, message string, attributes ...attribute.KeyValue) {
	l.write("debug", nil, message, attributes...)
}

func (l *CommandLineInterfaceLogger) Info(ctx context.Context, message string, attributes ...attribute.KeyValue) {
	l.write("info", nil, message, attributes...)
}

func (l *CommandLineInterfaceLogger) Warn(ctx context.Context, message string, attributes ...attribute.KeyValue) {
	l.write("warn", nil, message, attributes...)
}

func (l *CommandLineInterfaceLogger) Error(ctx context.Context, err error, message string, attributes ...attribute.KeyValue) {
	l.write("error", err, message, attributes...)
}

func (l *CommandLineInterfaceLogger) Alert(ctx context.Context, err error, message string, attributes ...attribute.KeyValue) {
	l.write("alert", err, message, attributes...)
}

func (l *Logger) write(
	ctx context.Context,
	level slog.Level,
	severity otellog.Severity,
	err error,
	message string,
	attributes ...attribute.KeyValue,
) {
	if ctx == nil {
		ctx = context.Background()
	}

	caller := traces.GetTrace(2)
	attributes = append(attributes,
		attribute.String("code.filepath", caller.File),
		attribute.Int("code.lineno", caller.Line),
		attribute.String("code.function", caller.Function),
	)
	if err != nil {
		attributes = append(attributes,
			attribute.String("error.type", typeName(err)),
			attribute.String("error.message", err.Error()),
		)
	}

	spanContext := trace.SpanContextFromContext(ctx)
	if spanContext.IsValid() {
		attributes = append(attributes,
			attribute.String("trace.id", spanContext.TraceID().String()),
			attribute.String("span.id", spanContext.SpanID().String()),
		)
	}

	slogAttributes := make([]slog.Attr, 0, len(attributes))
	otelAttributes := make([]otellog.KeyValue, 0, len(attributes))
	for _, item := range attributes {
		slogAttributes = append(slogAttributes, slog.Any(string(item.Key), item.Value.AsInterface()))
		otelAttributes = append(otelAttributes, otellog.KeyValueFromAttribute(item))
	}

	l.logger.LogAttrs(ctx, level, message, slogAttributes...)
	if !l.emitOTel {
		return
	}

	var record otellog.Record
	record.SetTimestamp(time.Now())
	record.SetSeverity(severity)
	record.SetSeverityText(severity.String())
	record.SetBody(otellog.StringValue(message))
	record.AddAttributes(otelAttributes...)
	logglobal.Logger("notezy").Emit(ctx, record)
}

func (l *CommandLineInterfaceLogger) write(
	level string,
	err error,
	message string,
	attributes ...attribute.KeyValue,
) {
	if err != nil {
		message += ": " + err.Error()
	}
	if len(attributes) > 0 {
		items := make([]string, 0, len(attributes))
		for _, item := range attributes {
			items = append(items, fmt.Sprintf("%s=%v", item.Key, item.Value.AsInterface()))
		}
		message += " " + strings.Join(items, " ")
	}

	l.logger.Printf("%s: %s", level, message)
}

func typeName(err error) string {
	return fmt.Sprintf("%T", err)
}
