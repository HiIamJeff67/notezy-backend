package infrastructure_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestInfrastructureComposeServices(t *testing.T) {
	if os.Getenv("NOTEZY_RUN_INTEGRATION") != "1" {
		t.Skip("set NOTEZY_RUN_INTEGRATION=1 to run Docker-backed integration tests")
	}

	repositoryRoot, err := os.Getwd()
	if err != nil {
		t.Fatalf("get working directory: %v", err)
	}
	repositoryRoot = filepath.Clean(filepath.Join(repositoryRoot, "..", "..", ".."))
	composeFile := filepath.Join(repositoryRoot, "infra", "docker", "docker-compose.integration.yaml")

	command := exec.Command(
		"docker",
		"compose",
		"--project-name", "notezy-integration",
		"--file", composeFile,
		"ps",
		"--status", "running",
		"--services",
	)
	output, err := command.Output()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			t.Fatalf("inspect integration Compose services: %v\n%s", err, exitError.Stderr)
		}
		t.Fatalf("inspect integration Compose services: %v", err)
	}

	runningServices := make(map[string]struct{})
	for _, service := range strings.Fields(string(output)) {
		runningServices[service] = struct{}{}
	}

	for _, service := range []string{
		"notezy-integration-db",
		"notezy-integration-notification-db",
		"notezy-integration-redis",
		"notezy-integration-kafka",
	} {
		if _, running := runningServices[service]; !running {
			t.Errorf("integration Compose service %q is not running", service)
		}
	}
}
