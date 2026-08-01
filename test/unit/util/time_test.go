package utilunittest

import (
	"testing"
	"time"

	test "github.com/HiIamJeff67/notezy-backend/test"
	testutil "github.com/HiIamJeff67/notezy-backend/test/util"

	"github.com/stretchr/testify/assert"
)

/* ============================== Test IsTimeWithin() ============================== */

type IsTimeWithinDeltaArgType = struct {
	T1    time.Time
	T2    time.Time
	Delta int64
}
type IsTimeWithinDeltaReturnType = bool
type IsTimeWithinDeltaTestCase = test.UnitTestCase[
	IsTimeWithinDeltaArgType,
	IsTimeWithinDeltaReturnType,
]

func TestIsTimeWithinDelta(t *testing.T) {
	cases := test.LoadTestCases[IsTimeWithinDeltaTestCase](
		t, "testdata/time_testdata/is_time_within_delta_testdata.json",
	)
	for _, c := range cases {
		got := testutil.IsTimeWithin(
			c.Args.T1,
			c.Args.T2,
			time.Duration(c.Args.Delta)*time.Second,
		)
		assert.Equal(t, c.Returns, got)
	}
}
