package observability

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	logglobal "go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/propagation"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"

	logs "github.com/HiIamJeff67/notezy-backend/internal/platform/observability/logs"
	metrics "github.com/HiIamJeff67/notezy-backend/internal/platform/observability/metrics"
	traces "github.com/HiIamJeff67/notezy-backend/internal/platform/observability/traces"
	constants "github.com/HiIamJeff67/notezy-backend/shared/constants"
)

func Initialize(ctx context.Context) func() {
	serviceName := os.Getenv("OTEL_SERVICE_NAME")
	if serviceName == "" {
		serviceName = constants.ServiceName
	}
	serviceVersion := os.Getenv("OTEL_SERVICE_VERSION")
	deploymentEnvironment := os.Getenv("OTEL_DEPLOYMENT_ENVIRONMENT")
	serviceInstanceId := os.Getenv("OTEL_SERVICE_INSTANCE_ID")
	if serviceInstanceId == "" {
		serviceInstanceId, _ = os.Hostname()
	}
	collectorEndpoint := os.Getenv("OTEL_EXPORTER_OTLP_GRPC_ENDPOINT")

	response, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
			semconv.ServiceVersion(serviceVersion),
			semconv.DeploymentEnvironment(deploymentEnvironment),
			semconv.ServiceInstanceID(serviceInstanceId),
		),
	)
	if err != nil {
		fmt.Println("Failed to create OpenTelemetry resource: ", err)
		response = resource.Default()
	}

	traceProviderOptions := []sdktrace.TracerProviderOption{
		sdktrace.WithResource(response),
	}
	traceExporter, err := otlptracegrpc.New(
		ctx,
		otlptracegrpc.WithEndpoint(collectorEndpoint),
		otlptracegrpc.WithInsecure(),
	)
	if err != nil {
		fmt.Println("Failed to create OpenTelemetry trace exporter: ", err)
	} else {
		traceProviderOptions = append(traceProviderOptions, sdktrace.WithBatcher(traceExporter))
	}
	traceProvider := sdktrace.NewTracerProvider(traceProviderOptions...)
	otel.SetTracerProvider(traceProvider)

	meterProviderOptions := []sdkmetric.Option{
		sdkmetric.WithResource(response),
	}
	metricExporter, err := otlpmetricgrpc.New(
		ctx,
		otlpmetricgrpc.WithEndpoint(collectorEndpoint),
		otlpmetricgrpc.WithInsecure(),
	)
	if err != nil {
		fmt.Println("Failed to create OpenTelemetry metric exporter: ", err)
	} else {
		meterProviderOptions = append(meterProviderOptions, sdkmetric.WithReader(
			sdkmetric.NewPeriodicReader(
				metricExporter,
				sdkmetric.WithInterval(15*time.Second),
			),
		))
	}
	meterProvider := sdkmetric.NewMeterProvider(meterProviderOptions...)
	otel.SetMeterProvider(meterProvider)

	logProviderOptions := []sdklog.LoggerProviderOption{
		sdklog.WithResource(response),
	}
	logExporter, err := otlploggrpc.New(
		ctx,
		otlploggrpc.WithEndpoint(collectorEndpoint),
		otlploggrpc.WithInsecure(),
	)
	if err != nil {
		fmt.Println("Failed to create OpenTelemetry log exporter: ", err)
	} else {
		logProviderOptions = append(logProviderOptions, sdklog.WithProcessor(sdklog.NewBatchProcessor(logExporter)))
	}
	logProvider := sdklog.NewLoggerProvider(logProviderOptions...)
	logglobal.SetLoggerProvider(logProvider)

	otel.SetTextMapPropagator(propagation.TraceContext{})

	logs.NotezyLogger = logs.NewLogger(true)
	metrics.NotezyMeter = metrics.NewMeter(otel.Meter(serviceName))
	traces.NotezyTracer = traces.NewTracer(otel.Tracer(serviceName))

	return func() {
		if err := traceProvider.Shutdown(ctx); err != nil {
			fmt.Println("Failed to shutdown OpenTelemetry trace provider: ", err)
		}
		if err := meterProvider.Shutdown(ctx); err != nil {
			fmt.Println("Failed to shutdown OpenTelemetry meter provider: ", err)
		}
		if err := logProvider.Shutdown(ctx); err != nil {
			fmt.Println("Failed to shutdown OpenTelemetry log provider: ", err)
		}
	}
}
