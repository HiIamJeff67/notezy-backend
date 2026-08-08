package infrastructure_test

import (
	"context"
	"os"
	"os/exec"
	"testing"
	"time"

	nat "github.com/docker/go-connections/nat"
	testcontainers "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestInfrastructureContainers(t *testing.T) {
	defer func() {
		if recovered := recover(); recovered != nil {
			t.Skipf("Docker provider is unavailable: %v", recovered)
		}
	}()

	if os.Getenv("NOTEZY_RUN_INTEGRATION") != "1" {
		t.Skip("set NOTEZY_RUN_INTEGRATION=1 to run Docker-backed integration tests")
	}
	if err := exec.Command("docker", "info", "--format", "{{.ServerVersion}}").Run(); err != nil {
		t.Skipf("Docker daemon is unavailable: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	t.Cleanup(cancel)

	postgres := runContainer(t, ctx, "postgres:16", map[string]string{
		"POSTGRES_DB":       "notezy_test",
		"POSTGRES_USER":     "notezy",
		"POSTGRES_PASSWORD": "notezy",
	}, "5432/tcp")
	redis := runContainer(t, ctx, "redis:7.2.3-alpine", nil, "6379/tcp")
	kafka := runKafkaContainer(t, ctx)

	t.Run("postgres", func(t *testing.T) {
		assertContainerPort(t, ctx, postgres, "5432/tcp")
	})
	t.Run("redis", func(t *testing.T) {
		assertContainerPort(t, ctx, redis, "6379/tcp")
	})
	t.Run("kafka", func(t *testing.T) {
		assertContainerPort(t, ctx, kafka, "9092/tcp")
	})
}

func runContainer(t *testing.T, ctx context.Context, image string, environment map[string]string, port string) testcontainers.Container {
	t.Helper()

	container, err := testcontainers.Run(
		ctx,
		image,
		testcontainers.WithEnv(environment),
		testcontainers.WithExposedPorts(port),
		testcontainers.WithWaitStrategy(wait.ForListeningPort(nat.Port(port)).WithStartupTimeout(2*time.Minute)),
	)
	if err != nil {
		t.Fatalf("start %s container: %v", image, err)
	}
	t.Cleanup(func() {
		if err := testcontainers.TerminateContainer(container); err != nil {
			t.Errorf("terminate %s container: %v", image, err)
		}
	})

	return container
}

func runKafkaContainer(t *testing.T, ctx context.Context) testcontainers.Container {
	t.Helper()

	return runContainer(t, ctx, "apache/kafka:4.0.0", map[string]string{
		"KAFKA_CLUSTER_ID":                               "MkU3OEVBNTcwNTJENDM2Qk",
		"KAFKA_NODE_ID":                                  "1",
		"KAFKA_PROCESS_ROLES":                            "broker,controller",
		"KAFKA_CONTROLLER_QUORUM_VOTERS":                 "1@localhost:9093",
		"KAFKA_LISTENERS":                                "PLAINTEXT://:9092,CONTROLLER://:9093",
		"KAFKA_ADVERTISED_LISTENERS":                     "PLAINTEXT://localhost:9092",
		"KAFKA_CONTROLLER_LISTENER_NAMES":                "CONTROLLER",
		"KAFKA_INTER_BROKER_LISTENER_NAME":               "PLAINTEXT",
		"KAFKA_LISTENER_SECURITY_PROTOCOL_MAP":           "CONTROLLER:PLAINTEXT,PLAINTEXT:PLAINTEXT",
		"KAFKA_AUTO_CREATE_TOPICS_ENABLE":                "true",
		"KAFKA_NUM_PARTITIONS":                           "3",
		"KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR":         "1",
		"KAFKA_TRANSACTION_STATE_LOG_REPLICATION_FACTOR": "1",
		"KAFKA_TRANSACTION_STATE_LOG_MIN_ISR":            "1",
	}, "9092/tcp")
}

func assertContainerPort(t *testing.T, ctx context.Context, container testcontainers.Container, port string) {
	t.Helper()

	host, err := container.Host(ctx)
	if err != nil {
		t.Fatalf("resolve container host: %v", err)
	}
	mappedPort, err := container.MappedPort(ctx, nat.Port(port))
	if err != nil {
		t.Fatalf("resolve mapped %s port on %s: %v", port, host, err)
	}
	if mappedPort.Port() == "" {
		t.Fatalf("mapped %s port is empty", port)
	}
}
