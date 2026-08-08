package strings

import (
	"encoding/json"
	"os"
	"testing"
)

func TestConvertCamelCaseToSentenceCase(t *testing.T) {
	cases := loadStringCases(t, "testdata/string/convert_camel_case_to_sentence_case_testdata.json")
	for index, testCase := range cases {
		t.Run(string(rune('A'+index)), func(t *testing.T) {
			if got := ConvertCamelCaseToSentenceCase(testCase.Args.Input); got != testCase.Returns {
				t.Fatalf("expected %q, got %q", testCase.Returns, got)
			}
		})
	}
}

type stringCase struct {
	Args struct {
		Input string `json:"input"`
	} `json:"args"`
	Returns string `json:"returns"`
}

func loadStringCases(t *testing.T, path string) []stringCase {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read testdata: %v", err)
	}

	var cases []stringCase
	if err := json.Unmarshal(data, &cases); err != nil {
		t.Fatalf("decode testdata: %v", err)
	}

	return cases
}
