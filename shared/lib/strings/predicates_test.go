package strings

import (
	"encoding/json"
	"os"
	"testing"
)

func TestIsEmailString(t *testing.T) {
	cases := loadPredicateCases(t, "testdata/string/is_email_string_testdata.json")
	for index, testCase := range cases {
		t.Run(string(rune('A'+index)), func(t *testing.T) {
			if got := IsEmailString(testCase.Args.S); got != testCase.Returns {
				t.Fatalf("expected %t, got %t", testCase.Returns, got)
			}
		})
	}
}

func TestIsAlphaOrNumberString(t *testing.T) {
	cases := loadPredicateCases(t, "testdata/string/is_alpha_number_string_testdata.json")
	for index, testCase := range cases {
		t.Run(string(rune('A'+index)), func(t *testing.T) {
			if got := IsAlphaOrNumberString(testCase.Args.S); got != testCase.Returns {
				t.Fatalf("expected %t, got %t", testCase.Returns, got)
			}
		})
	}
}

func TestIsNumberString(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{name: "digits", value: "123456", want: true},
		{name: "empty", value: "", want: true},
		{name: "contains letters", value: "123a", want: false},
		{name: "contains punctuation", value: "12.3", want: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := IsNumberString(test.value); got != test.want {
				t.Fatalf("IsNumberString(%q) = %t, want %t", test.value, got, test.want)
			}
		})
	}
}

func TestIsAlphaAndNumberString(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{name: "letters and digits", value: "Notezy123", want: true},
		{name: "letters only", value: "Notezy", want: false},
		{name: "digits only", value: "123", want: false},
		{name: "contains punctuation", value: "Notezy-123", want: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := IsAlphaAndNumberString(test.value); got != test.want {
				t.Fatalf("IsAlphaAndNumberString(%q) = %t, want %t", test.value, got, test.want)
			}
		})
	}
}

type predicateCase struct {
	Args struct {
		S string `json:"s"`
	} `json:"args"`
	Returns bool `json:"returns"`
}

func loadPredicateCases(t *testing.T, path string) []predicateCase {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read testdata: %v", err)
	}

	var cases []predicateCase
	if err := json.Unmarshal(data, &cases); err != nil {
		t.Fatalf("decode testdata: %v", err)
	}

	return cases
}
