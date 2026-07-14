package app

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
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

	caches "github.com/HiIamJeff67/notezy-backend/app/caches"
	commands "github.com/HiIamJeff67/notezy-backend/app/commands"
	configs "github.com/HiIamJeff67/notezy-backend/app/configs"
	routinetask "github.com/HiIamJeff67/notezy-backend/app/durablejobs/routinetask"
	yjsmaintenance "github.com/HiIamJeff67/notezy-backend/app/durablejobs/yjsmaintenance"
	models "github.com/HiIamJeff67/notezy-backend/app/models"
	logs "github.com/HiIamJeff67/notezy-backend/app/monitor/logs"
	metrics "github.com/HiIamJeff67/notezy-backend/app/monitor/metrics"
	traces "github.com/HiIamJeff67/notezy-backend/app/monitor/traces"
	developmentroutes "github.com/HiIamJeff67/notezy-backend/app/routes/developmentroutes"
	util "github.com/HiIamJeff67/notezy-backend/app/util"
	constants "github.com/HiIamJeff67/notezy-backend/shared/constants"
)

func StartApplication(isCLI bool) {
	if isCLI {
		commands.Execute()
		return
	}

	shutdownObservability := initializeObservability(context.Background())
	defer shutdownObservability()

	models.NotezyDB = models.ConnectToDatabase(configs.PostgresDatabaseConfig)
	defer models.DisconnectToDatabase(models.NotezyDB)

	caches.ConnectToAllRedis()
	defer caches.DisconnectToAllRedis()
	reloadRedisLibraries()

	routineTaskEngine := routinetask.NewEngine(models.NotezyDB)
	shutdownRoutineTaskEngine := routineTaskEngine.Start(context.Background())
	defer shutdownRoutineTaskEngine()

	yjsMaintenanceEngine := yjsmaintenance.NewEngine(models.NotezyDB)
	shutdownYjsMaintenanceEngine := yjsMaintenanceEngine.Start(context.Background())
	defer shutdownYjsMaintenanceEngine()

	developmentroutes.DevelopmentRouter = gin.Default()
	proxies := strings.Split(util.GetEnv("GIN_TRUSTED_PROXIES", ""), ",")
	if err := developmentroutes.DevelopmentRouter.SetTrustedProxies(proxies); err != nil {
		fmt.Println("Failed to set trusted proxies for router: ", err)
		return
	}
	developmentroutes.ConfigureDevelopmentRoutes()
	ginAddr := util.GetEnv("GIN_DOMAIN", "") + ":" + util.GetEnv("GIN_PORT", "7777")
	if err := endless.ListenAndServe(ginAddr, developmentroutes.DevelopmentRouter); err != nil {
		fmt.Println("Failed to connect to the server")
		return
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

func initializeObservability(ctx context.Context) func() {
	serviceName := util.GetEnv("OTEL_SERVICE_NAME", constants.ServiceName)
	serviceVersion := util.GetEnv("OTEL_SERVICE_VERSION", constants.DevelopmentCompleteVersion)
	deploymentEnvironment := util.GetEnv("OTEL_DEPLOYMENT_ENVIRONMENT", util.GetEnv("NODE_ENV", "development"))
	serviceInstanceId := util.GetEnv("OTEL_SERVICE_INSTANCE_ID", "")
	if serviceInstanceId == "" {
		serviceInstanceId, _ = os.Hostname()
	}
	collectorEndpoint := util.GetEnv(
		"OTEL_EXPORTER_OTLP_GRPC_ENDPOINT",
		util.GetEnv("DOCKER_OTEL_COLLECTOR_SERVICE_NAME", "notezy-otel-collector")+":"+
			util.GetEnv("DOCKER_OTEL_COLLECTOR_GRPC_PORT", "4317"),
	)

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
	metrics.NotezyMeter = metrics.NewMeter(otel.Meter(constants.ServiceName))
	traces.NotezyTracer = traces.NewTracer(otel.Tracer(constants.ServiceName))

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
