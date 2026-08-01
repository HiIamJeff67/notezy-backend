package logs

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"time"

	"go.opentelemetry.io/otel/attribute"
	otellog "go.opentelemetry.io/otel/log"
	logglobal "go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/trace"

	traces "github.com/HiIamJeff67/notezy-backend/internal/platform/observability/traces"
)

type LoggerInterface interface {
	Debug(ctx context.Context, message string, attributes ...attribute.KeyValue)
	Info(ctx context.Context, message string, attributes ...attribute.KeyValue)
	Warn(ctx context.Context, message string, attributes ...attribute.KeyValue)
	Error(ctx context.Context, err error, message string, attributes ...attribute.KeyValue)
	Alert(ctx context.Context, err error, message string, attributes ...attribute.KeyValue)
	JSON(ctx context.Context, level slog.Level, message string, payload any, attributes ...attribute.KeyValue) error
}

type Logger struct {
	logger   *slog.Logger
	emitOTel bool
}

func NewLogger(emitOtel bool) LoggerInterface {
	return &Logger{
		logger:   slog.New(slog.NewJSONHandler(os.Stdout, nil)),
		emitOTel: emitOtel,
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

func (l *Logger) JSON(
	ctx context.Context,
	level slog.Level,
	message string,
	payload any,
	attributes ...attribute.KeyValue,
) error {
	marshaledPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	attributes = append(attributes, attribute.String("payload", string(marshaledPayload)))
	l.write(ctx, level, severityFromLevel(level), nil, message, attributes...)

	return nil
}

func severityFromLevel(level slog.Level) otellog.Severity {
	switch {
	case level <= slog.LevelDebug:
		return otellog.SeverityDebug
	case level <= slog.LevelInfo:
		return otellog.SeverityInfo
	case level <= slog.LevelWarn:
		return otellog.SeverityWarn
	default:
		return otellog.SeverityError
	}
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

func typeName(err error) string {
	return fmt.Sprintf("%T", err)
}
