package integration_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestComposeRuntimeBoundaries(t *testing.T) {
	repositoryRoot := repositoryRoot(t)
	compose, err := os.ReadFile(filepath.Join(repositoryRoot, "docker-compose.yaml"))
	if err != nil {
		t.Fatalf("read docker compose file: %v", err)
	}

	composeText := string(compose)
	for _, serviceName := range []string{
		"  notezy-client-gateway:",
		"  notezy-core:",
		"  notezy-realtime-gateway:",
		"  notezy-durable-job:",
		"  notezy-email:",
		"  notezy-yjs-worker:",
	} {
		if !strings.Contains(composeText, serviceName) {
			t.Errorf("docker compose is missing runtime service %q", strings.TrimSpace(serviceName))
		}
	}

	if !strings.Contains(composeText, "/healthz") {
		t.Error("docker compose is missing the health status route")
	}
}

func repositoryRoot(t *testing.T) string {
	t.Helper()

	workingDirectory, err := os.Getwd()
	if err != nil {
		t.Fatalf("get working directory: %v", err)
	}

	return filepath.Clean(filepath.Join(workingDirectory, "..", ".."))
}
