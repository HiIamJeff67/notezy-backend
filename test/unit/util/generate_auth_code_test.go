package unit_test_util

import (
	"notezy-backend/app/util"
	"notezy-backend/shared/constants"
	"strconv"
	"testing"
)

func TestGenerateAuthCode(t *testing.T) {
	code := util.GenerateAuthCode()
	if len(code) != constants.MaxLengthOfAuthCode {
		t.Errorf("unexpected code length: %s", code)
	}
	num, err := strconv.Atoi(code)
	if err != nil || num < 0 || num > constants.MaxAuthCode {
		t.Errorf("unexpected code value: %s", code)
	}
}
