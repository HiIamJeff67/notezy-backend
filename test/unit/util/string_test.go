package utilunittest

import (
	"testing"

	"github.com/stretchr/testify/assert"

	util "notezy-backend/app/util"
	test "notezy-backend/test"
)

/* ============================== Test JoinValues ============================== */

type JoinValuesArgType = struct {
	Values []string
}
type JoinValuesReturnType = string
type JoinValuesTestCase = test.UnitTestCase[
	JoinValuesArgType,
	JoinValuesReturnType,
]

func TestJoinValues(t *testing.T) {
	cases := test.LoadTestCases[JoinValuesTestCase](
		t, "testdata/string_testdata/join_values_testdata.json",
	)
	for _, c := range cases {
		got := util.JoinValues(c.Args.Values)
		assert.Equal(t, c.Returns, got)
	}
}

/* ============================== Test ConvertCamelCaseToSenctenceCase ============================== */

type ConvertCamelCaseToSentenceCaseArgType = struct {
	Input string
}
type ConvertCamelCaseToSentenceCaseReturnType = string
type ConvertCamelCaseToSentenceCaseTestCase = test.UnitTestCase[
	ConvertCamelCaseToSentenceCaseArgType,
	ConvertCamelCaseToSentenceCaseReturnType,
]

func TestConvertCamelCaseToSentenceCase(t *testing.T) {
	cases := test.LoadTestCases[ConvertCamelCaseToSentenceCaseTestCase](
		t, "testdata/string_testdata/convert_camel_case_to_sentence_case_testdata.json",
	)
	for _, c := range cases {
		got := util.ConvertCamelCaseToSentenceCase(c.Args.Input)
		assert.Equal(t, c.Returns, got)
	}
}

/* ============================== Test IsStringIn ============================== */

type IsStringInArgType = struct {
	S    string
	Strs []string
}
type IsStringInReturnType = bool
type IsStringInTestCase = test.UnitTestCase[
	IsStringInArgType,
	IsStringInReturnType,
]

func TestIsStringIn(t *testing.T) {
	cases := test.LoadTestCases[IsStringInTestCase](
		t, "testdata/string_testdata/is_string_in_testdata.json",
	)
	for _, c := range cases {
		got := util.IsStringIn(c.Args.S, c.Args.Strs)
		assert.Equal(t, c.Returns, got)
	}
}

/* ============================== Test IsEmailString ============================== */

type IsEmailStringArgType = struct {
	S string
}
type IsEmailStringReturnType = bool
type IsEmailStringTestCase = test.UnitTestCase[
	IsEmailStringArgType,
	IsEmailStringReturnType,
]

func TestIsEmailString(t *testing.T) {
	cases := test.LoadTestCases[IsEmailStringTestCase](
		t, "testdata/string_testdata/is_email_string_testdata.json",
	)
	for _, c := range cases {
		got := util.IsEmailString(c.Args.S)
		assert.Equal(t, c.Returns, got)
	}
}

/* ============================== Test IsAlphaNumberString ============================== */
type IsAlphaNumberStringArgType = struct {
	S string
}
type IsAlphaNumberStringReturnType = bool
type IsAlphaNumberStringTestCase = test.UnitTestCase[
	IsAlphaNumberStringArgType,
	IsAlphaNumberStringReturnType,
]

func TestIsAlphaNumberString(t *testing.T) {
	cases := test.LoadTestCases[IsAlphaNumberStringTestCase](
		t, "testdata/string_testdata/is_alpha_number_string_testdata.json",
	)
	for _, c := range cases {
		got := util.IsAlphaNumberString(c.Args.S)
		assert.Equal(t, c.Returns, got)
	}
}
