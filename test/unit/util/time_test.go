package utilunittest

import (
	util "notezy-backend/app/util"
	test "notezy-backend/test"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

/* ============================== Test IsTimeWithinDelta() ============================== */

type IsTimeWithinDeltaArgType = struct {
	T1    time.Time
	T2    time.Time
	Delta time.Duration
}
type IsTimeWithinDeltaReturnType = bool
type IsTimeWithinDeltaTestCase = test.UnitTestCase[
	IsTimeWithinDeltaArgType,
	IsTimeWithinDeltaReturnType,
]

func TestIsTimeWithinDelta(t *testing.T) {
	cases := test.LoadTestCases[IsTimeWithinDeltaTestCase](
		t, "testdata/string_testdata/join_values_testdata.json",
	)
	for _, c := range cases {
		got := util.IsTimeWithinDelta(c.Args.T1, c.Args.T2, c.Args.Delta)
		assert.Equal(t, c.Returns, got)
	}
}
