package app

import (
	"context"
	"fmt"
	"time"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"

	caches "notezy-backend/app/caches"
	configs "notezy-backend/app/configs"
	models "notezy-backend/app/models"
	developmentroutes "notezy-backend/app/routes/development_routes"
	util "notezy-backend/app/util"
	constants "notezy-backend/shared/constants"
)

func StartApplication() {
	models.NotezyDB = models.ConnectToDatabase(configs.PostgresDatabaseConfig)
	defer models.DisconnectToDatabase(models.NotezyDB)

	caches.ConnectToAllRedis()
	defer caches.DisconnectToAllRedis()
	reloadRedisLibraries()

	ctx := context.Background()
	shutdown, err := initOTel(ctx)
	if err != nil {
		fmt.Println("Failed to initialize OpenTelemetry: ", err)
		return
	}
	defer shutdown()

	developmentroutes.DevelopmentRouter = gin.Default()
	developmentroutes.ConfigureDevelopmentRoutes()
	ginAddr := util.GetEnv("GIN_DOMAIN", "") + ":" + util.GetEnv("GIN_PORT", "7777")
	if err := endless.ListenAndServe(ginAddr, developmentroutes.DevelopmentRouter); err != nil {
		fmt.Println("Failed to connect to the server")
	}
}

func reloadRedisLibraries() {
	if exception := caches.FlushCacheLibraries(); exception != nil {
		exception.Log()
	}
	if exception := caches.LoadRateLimitRecordCacheLibraries(); exception != nil {
		exception.Log()
	}
	if exception := caches.LoadUserQuotaCacheLibraries(); exception != nil {
		exception.Log()
	}
	// reload other more redis libraries here...
}

func initOTel(ctx context.Context) (func(), error) {
	response, _ := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(constants.ServiceName),
			semconv.ServiceVersion(constants.DevelopmentCompleteVersion),
		),
	)

	traceExporter, _ := otlptracegrpc.New(
		ctx,
		otlptracegrpc.WithEndpoint(
			util.GetEnv("DOCKER_OTEL_COLLECTOR_SERVICE_NAME", "notezy-otel-collector")+
				":"+
				util.GetEnv("DOCKER_OTEL_COLLECTOR_GRPC_PORT", "4317"),
		),
		otlptracegrpc.WithInsecure(),
	)
	traceProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithResource(response),
	)
	otel.SetTracerProvider(traceProvider)

	metricExporter, _ := otlpmetricgrpc.New(
		ctx,
		otlpmetricgrpc.WithEndpoint(
			util.GetEnv("DOCKER_OTEL_COLLECTOR_SERVICE_NAME", "notezy-otel-collector")+
				":"+
				util.GetEnv("DOCKER_OTEL_COLLECTOR_GRPC_PORT", "4317"),
		),
		otlpmetricgrpc.WithInsecure(),
	)
	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(
			sdkmetric.NewPeriodicReader(
				metricExporter,
				sdkmetric.WithInterval(15*time.Second),
			),
		),
		sdkmetric.WithResource(response),
	)
	otel.SetMeterProvider(meterProvider)

	logExporter, _ := otlploggrpc.New(
		ctx,
		otlploggrpc.WithEndpoint(
			util.GetEnv("DOCKER_OTEL_COLLECTOR_SERVICE_NAME", "notezy-otel-collector")+
				":"+
				util.GetEnv("DOCKER_OTEL_COLLECTOR_GRPC_PORT", "4317"),
		),
		otlploggrpc.WithInsecure(),
	)
	logProvider := sdklog.NewLoggerProvider(
		sdklog.WithProcessor(sdklog.NewBatchProcessor(logExporter)),
		sdklog.WithResource(response),
	)

	return func() {
		traceProvider.Shutdown(ctx)
		meterProvider.Shutdown(ctx)
		logProvider.Shutdown(ctx)
	}, nil
}
