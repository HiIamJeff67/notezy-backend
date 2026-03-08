package utilunittest

import (
	"regexp"
	"testing"

	util "notezy-backend/app/util"
)

/* ============================== Test GenerateRandomFakeDisplayName() ============================== */

func TestGenerateRandomFakeDisplayName(t *testing.T) {
	got := util.GenerateRandomFakeDisplayName()
	// 檢查格式: 字母開頭，結尾為6位數字
	assert := regexp.MustCompile(`^[A-Za-z]+[0-9]{6}$`).MatchString
	if !assert(got) {
		t.Errorf("unexpected format: %s", got)
	}
}
