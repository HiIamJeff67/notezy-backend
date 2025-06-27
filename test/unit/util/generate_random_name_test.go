package unit_test_util

import (
	"regexp"
	"testing"

	util "notezy-backend/app/util"
)

func TestGenerateRandomFakeName(t *testing.T) {
	got := util.GenerateRandomFakeName()
	// 檢查格式: 字母開頭，結尾為6位數字
	assert := regexp.MustCompile(`^[A-Za-z]+[0-9]{6}$`).MatchString
	if !assert(got) {
		t.Errorf("unexpected format: %s", got)
	}
}
