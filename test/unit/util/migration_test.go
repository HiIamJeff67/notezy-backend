package utilunittest

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"

	util "notezy-backend/app/util"
	test "notezy-backend/test"
)

/* ============================== Test GenerateMigrationFileName() ============================== */

type GenerateMigrationFileNameArgType = struct {
	DBName string
}
type GenerateMigrationFileNameReturnType = string
type GenerateMigrationFileNameTestCase = test.UnitTestCase[
	GenerateMigrationFileNameArgType,
	GenerateMigrationFileNameReturnType,
]

func TestGenerateMigrationFileName(t *testing.T) {
	cases := test.LoadTestCases[GenerateMigrationFileNameTestCase](
		t, "testdata/migration_testdata/generate_migration_file_name_testdata.json",
	)
	for _, c := range cases {
		got := util.GenerateMigrationFileName(c.Args.DBName)
		// only validate the beginning, since the content is randomly generated
		assert.True(t, regexp.MustCompile("^"+regexp.QuoteMeta(c.Args.DBName)+"_").MatchString(got))
	}
}
