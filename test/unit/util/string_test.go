package utilunittest

import (
	"testing"

	"github.com/stretchr/testify/assert"

	strings "github.com/HiIamJeff67/notezy-backend/shared/lib/strings"
	validators "github.com/HiIamJeff67/notezy-backend/shared/validations/validators"
	test "github.com/HiIamJeff67/notezy-backend/test"
)

/* ============================== Test ConvertCamelCaseToSenctenceCase() ============================== */

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
		got := strings.ConvertCamelCaseToSentenceCase(c.Args.Input)
		assert.Equal(t, c.Returns, got)
	}
}

/* ============================== Test IsEmailString() ============================== */

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
		got := validators.IsEmailString(c.Args.S)
		assert.Equal(t, c.Returns, got)
	}
}

/* ============================== Test IsAlphaNumberString() ============================== */
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
		got := validators.IsAlphaOrNumberString(c.Args.S)
		assert.Equal(t, c.Returns, got)
	}
}
