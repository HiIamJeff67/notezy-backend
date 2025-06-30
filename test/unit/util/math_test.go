package utilunittest

import (
	"testing"

	"github.com/stretchr/testify/assert"

	util "notezy-backend/app/util"
	test "notezy-backend/test"
)

/* ============================== Test GetMinInMap() ============================== */

type GetMinInMapArgType = struct {
	Map map[string]int
}
type GetMinInMapReturnType = int
type GetMinInMapTestCase = test.UnitTestCase[
	GetMinInMapArgType,
	GetMinInMapReturnType,
]

func TestGetMinInMap(t *testing.T) {
	cases := test.LoadTestCases[GetMinInMapTestCase](
		t, "testdata/math_testdata/get_min_in_map_testdata.json",
	)
	for _, c := range cases {
		got := util.GetMinInMap(c.Args.Map)
		assert.Equal(t, c.Returns, got)
	}
}

/* ============================== Test GetMaxInMap() ============================== */

type GetMaxInMapArgType = struct {
	Map map[string]int
}
type GetMaxInMapReturnType = int
type GetMaxInMapTestCase = test.UnitTestCase[
	GetMaxInMapArgType,
	GetMaxInMapReturnType,
]

func TestGetMaxInMap(t *testing.T) {
	cases := test.LoadTestCases[GetMaxInMapTestCase](
		t, "testdata/math_testdata/get_max_in_map_testdata.json",
	)
	for _, c := range cases {
		got := util.GetMaxInMap(c.Args.Map)
		assert.Equal(t, c.Returns, got)
	}
}
