package test

import (
	"encoding/json"
	"os"
	"testing"
)

func LoadTestCases[T any](t *testing.T, filename string) []T {
	data, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("failed to read testdata: %v", err)
	}
	var cases []T
	if err := json.Unmarshal(data, &cases); err != nil {
		t.Fatalf("failed to unmarshal testdata: %v", err)
	}
	return cases
}
