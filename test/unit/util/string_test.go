package unit_test_util

import (
	"encoding/json"
	"os"
	"testing"

	util "notezy-backend/app/util"
	types "notezy-backend/shared/types"

	"github.com/stretchr/testify/assert"
)

// JoinValues
type JoinValuesArgType = struct {
	Values []string
}
type JoinValuesReturnType = string
type JoinValuesTestCase = types.TestCase[JoinValuesArgType, JoinValuesReturnType]

// ConvertCamelCaseToSentenceCase
type ConvertCamelCaseToSentenceCaseArgType = struct {
	Input string
}
type ConvertCamelCaseToSentenceCaseReturnType = string
type ConvertCamelCaseToSentenceCaseTestCase = types.TestCase[ConvertCamelCaseToSentenceCaseArgType, ConvertCamelCaseToSentenceCaseReturnType]

// IsStringIn
type IsStringInArgType = struct {
	S    string
	Strs []string
}
type IsStringInReturnType = bool
type IsStringInTestCase = types.TestCase[IsStringInArgType, IsStringInReturnType]

// IsEmailString
type IsEmailStringArgType = struct {
	S string
}
type IsEmailStringReturnType = bool
type IsEmailStringTestCase = types.TestCase[IsEmailStringArgType, IsEmailStringReturnType]

// IsAlphaNumberString
type IsAlphaNumberStringArgType = struct {
	S string
}
type IsAlphaNumberStringReturnType = bool
type IsAlphaNumberStringTestCase = types.TestCase[IsAlphaNumberStringArgType, IsAlphaNumberStringReturnType]

func loadTestCases[T any](t *testing.T, filename string) []T {
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

func TestJoinValues(t *testing.T) {
	cases := loadTestCases[JoinValuesTestCase](
		t, "testdata/string_testdata/join_values_testdata.json",
	)
	for _, c := range cases {
		got := util.JoinValues(c.Args.Values)
		assert.Equal(t, c.Returns, got)
	}
}

func TestConvertCamelCaseToSentenceCase(t *testing.T) {
	cases := loadTestCases[ConvertCamelCaseToSentenceCaseTestCase](
		t, "testdata/string_testdata/convert_camel_case_to_sentence_case_testdata.json",
	)
	for _, c := range cases {
		got := util.ConvertCamelCaseToSentenceCase(c.Args.Input)
		assert.Equal(t, c.Returns, got)
	}
}

func TestIsStringIn(t *testing.T) {
	cases := loadTestCases[IsStringInTestCase](
		t, "testdata/string_testdata/is_string_in_testdata.json",
	)
	for _, c := range cases {
		got := util.IsStringIn(c.Args.S, c.Args.Strs)
		assert.Equal(t, c.Returns, got)
	}
}

func TestIsEmailString(t *testing.T) {
	cases := loadTestCases[IsEmailStringTestCase](
		t, "testdata/string_testdata/is_email_string_testdata.json",
	)
	for _, c := range cases {
		got := util.IsEmailString(c.Args.S)
		assert.Equal(t, c.Returns, got)
	}
}

func TestIsAlphaNumberString(t *testing.T) {
	cases := loadTestCases[IsAlphaNumberStringTestCase](
		t, "testdata/string_testdata/is_alpha_number_string_testdata.json",
	)
	for _, c := range cases {
		got := util.IsAlphaNumberString(c.Args.S)
		assert.Equal(t, c.Returns, got)
	}
}
