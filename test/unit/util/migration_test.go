package unit_test_util

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"

	util "notezy-backend/app/util"
	types "notezy-backend/shared/types"
	test "notezy-backend/test"
)

type GenerateMigrationFileNameArgType = struct {
	DBName string
}
type GenerateMigrationFileNameReturnType = string
type GenerateMigrationFileNameTestCase = types.TestCase[GenerateMigrationFileNameArgType, GenerateMigrationFileNameReturnType]

func TestGenerateMigrationFileName(t *testing.T) {
	cases := test.LoadTestCases[GenerateMigrationFileNameTestCase](
		t, "testdata/migration_testdata/generate_migration_file_name_testdata.json",
	)
	for _, c := range cases {
		got := util.GenerateMigrationFileName(c.Args.DBName)
		// 只驗證開頭是否正確，因為內容是隨機的
		assert.True(t, regexp.MustCompile("^"+regexp.QuoteMeta(c.Args.DBName)+"_").MatchString(got))
	}
}
