package utilunittest

import (
	"strconv"
	"testing"

	util "notezy-backend/app/util"
)

/* ============================== Test GenerateAuthCode() ============================== */

func TestGenerateAuthCode(t *testing.T) {
	code := util.GenerateAuthCode()
	if len(code) != util.MaxLengthOfAuthCode {
		t.Errorf("unexpected code length: %s", code)
	}
	num, err := strconv.Atoi(code)
	if err != nil || num < 0 || num > util.MaxAuthCode {
		t.Errorf("unexpected code value: %s", code)
	}
}
