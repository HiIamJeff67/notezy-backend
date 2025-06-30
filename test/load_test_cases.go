package test

import (
	"encoding/json"
	"os"
	"testing"
)

func LoadTestCase[TestCaseType any](t *testing.T, relativePath string) TestCaseType {
	data, err := os.ReadFile(relativePath)
	if err != nil {
		t.Fatalf("failed to read testdata: %v", err)
	}
	var testCase TestCaseType
	if err := json.Unmarshal(data, &testCase); err != nil {
		t.Fatalf("failed to unmarshal testdata: %v", err)
	}
	return testCase
}

func LoadTestCases[TestCaseType any](t *testing.T, relativePath string) []TestCaseType {
	data, err := os.ReadFile(relativePath)
	if err != nil {
		t.Fatalf("failed to read testdata: %v", err)
	}
	var testCases []TestCaseType
	if err := json.Unmarshal(data, &testCases); err != nil {
		t.Fatalf("failed to unmarshal testdata: %v", err)
	}
	return testCases
}
