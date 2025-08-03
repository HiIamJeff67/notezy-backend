package utilunittest

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	util "notezy-backend/app/util"
	test "notezy-backend/test"
)

/* ============================== Test GetLoginBlockedUntilByLoginCount() ============================== */

type GetLoginBlockedUntilByLoginCountArgType = struct {
	LoginCount int32
}
type GetLoginBlockedUntilByLoginCountReturnType = struct {
	ShouldBlock bool
}
type GetLoginBlockedUntilByLoginCountTestCase = test.UnitTestCase[
	GetLoginBlockedUntilByLoginCountArgType,
	GetLoginBlockedUntilByLoginCountReturnType,
]

func TestGetLoginBlockedUntilByLoginCount(t *testing.T) {
	cases := test.LoadTestCases[GetLoginBlockedUntilByLoginCountTestCase](
		t, "testdata/get_blocked_time_testdata/get_login_blocked_until_by_login_count_testdata.json",
	)
	for _, c := range cases {
		got, _ := util.GetLoginBlockedUntilByLoginCount(c.Args.LoginCount)
		if c.Returns.ShouldBlock {
			assert.NotNil(t, got)
			assert.True(t, got.After(time.Now().Add(-1*time.Second)))
		} else {
			assert.Nil(t, got)
		}
	}
}

/* ============================== Test ShouldBlockLogin() ============================== */

type ShouldBlockLoginArgType = struct {
	LoginCount int32
}
type ShouldBlockLoginReturnType = bool
type ShouldBlockLoginTestCase = test.UnitTestCase[
	ShouldBlockLoginArgType,
	ShouldBlockLoginReturnType,
]

func TestShouldBlockLogin(t *testing.T) {
	cases := test.LoadTestCases[ShouldBlockLoginTestCase](
		t, "testdata/get_blocked_time_testdata/should_block_login_testdata.json",
	)
	for _, c := range cases {
		got := util.ShouldBlockLogin(c.Args.LoginCount)
		assert.Equal(t, c.Returns, got)
	}
}

/* ============================== Test GetNextBlockThreshold() ============================== */

type GetNextBlockThresholdArgType = struct {
	LoginCount int32
}
type GetNextBlockThresholdReturnType = int32
type GetNextBlockThresholdTestCase = test.UnitTestCase[
	GetNextBlockThresholdArgType,
	GetNextBlockThresholdReturnType,
]

func TestGetNextBlockThreshold(t *testing.T) {
	cases := test.LoadTestCases[GetNextBlockThresholdTestCase](
		t, "testdata/get_blocked_time_testdata/get_next_block_threshold_testdata.json",
	)
	for _, c := range cases {
		got := util.GetNextBlockThreshold(c.Args.LoginCount)
		assert.Equal(t, c.Returns, got)
	}
}
