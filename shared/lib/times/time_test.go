package times

import (
	"encoding/json"
	"os"
	"testing"
	"time"
)

func TestIsTimeWithin(t *testing.T) {
	cases := loadTimeCases(t, "testdata/time/is_time_within_delta_testdata.json")
	for index, testCase := range cases {
		t.Run(string(rune('A'+index)), func(t *testing.T) {
			got := IsTimeWithin(
				testCase.Args.T1,
				testCase.Args.T2,
				time.Duration(testCase.Args.Delta)*time.Second,
			)
			if got != testCase.Returns {
				t.Fatalf("expected %t, got %t", testCase.Returns, got)
			}
		})
	}
}

type timeCase struct {
	Args struct {
		T1    time.Time `json:"t1"`
		T2    time.Time `json:"t2"`
		Delta int64     `json:"delta"`
	} `json:"args"`
	Returns bool `json:"returns"`
}

func loadTimeCases(t *testing.T, path string) []timeCase {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read testdata: %v", err)
	}

	var cases []timeCase
	if err := json.Unmarshal(data, &cases); err != nil {
		t.Fatalf("decode testdata: %v", err)
	}

	return cases
}
