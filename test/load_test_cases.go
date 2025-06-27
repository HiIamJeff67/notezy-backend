package test

import (
	"encoding/json"
	"os"
	"testing"
)

func LoadTestCases[T any](t *testing.T, relativePath string) []T {
	data, err := os.ReadFile(relativePath)
	if err != nil {
		t.Fatalf("failed to read testdata: %v", err)
	}
	var cases []T
	if err := json.Unmarshal(data, &cases); err != nil {
		t.Fatalf("failed to unmarshal testdata: %v", err)
	}
	return cases
}
