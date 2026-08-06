package observability

import (
	"os"
	"strings"
)

type Config struct {
	ServiceName           string
	ServiceVersion        string
	DeploymentEnvironment string
	ServiceInstanceId     string
	CollectorEndpoint     string
}

func LoadConfig(defaultServiceName string) Config {
	serviceName := strings.TrimSpace(os.Getenv("OTEL_SERVICE_NAME"))
	if serviceName == "" {
		serviceName = defaultServiceName
	}
	serviceInstanceId := strings.TrimSpace(os.Getenv("OTEL_SERVICE_INSTANCE_ID"))
	if serviceInstanceId == "" {
		serviceInstanceId, _ = os.Hostname()
	}

	return Config{
		ServiceName:           serviceName,
		ServiceVersion:        strings.TrimSpace(os.Getenv("OTEL_SERVICE_VERSION")),
		DeploymentEnvironment: strings.TrimSpace(os.Getenv("OTEL_DEPLOYMENT_ENVIRONMENT")),
		ServiceInstanceId:     serviceInstanceId,
		CollectorEndpoint:     strings.TrimSpace(os.Getenv("OTEL_EXPORTER_OTLP_GRPC_ENDPOINT")),
	}
}
