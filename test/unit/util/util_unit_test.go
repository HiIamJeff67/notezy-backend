package utilunittest

import (
	"fmt"
	"testing"
)

const (
	testTargetPath = "notezy-backend/app/util/"
)

func TestObjectInParallel(t *testing.T) {
	t.Run("Unit-Test---Util", func(t *testing.T) { // object level
		t.Run(fmt.Sprintf("Test-Generate-Auth-Code-(%s):", testTargetPath+"generate_auth_code.go"), func(t *testing.T) { // file level
			t.Run("Test-GenerateAuthCode()", func(t *testing.T) { // function level
				t.Parallel()
				TestGenerateAuthCode(t)
			})
		})
		t.Run(fmt.Sprintf("Test-Generate-Random-Name-(%s):", testTargetPath+"generate_random_name.go"), func(t *testing.T) { // file level
			t.Run("Test-GenerateRandomFakeName()", func(t *testing.T) { // function level
				t.Parallel()
				TestGenerateRandomFakeName(t)
			})
		})
		t.Run(fmt.Sprintf("Test-Get-Blocked-Time-(%s):", testTargetPath+"get_blocked_time.go"), func(t *testing.T) { // file level
			t.Run("Test-GetLoginBlockedUntilByLoginCount()", func(t *testing.T) { // function level
				t.Parallel()
				TestGetLoginBlockedUntilByLoginCount(t)
			})
			t.Run("Test-ShouldBlockLogin()", func(t *testing.T) { // function level
				t.Parallel()
				TestShouldBlockLogin(t)
			})
			t.Run("Test-GetNextBlockThreshold()", func(t *testing.T) { // function level
				t.Parallel()
				TestGetNextBlockThreshold(t)
			})
		})
		t.Run(fmt.Sprintf("Test-Math-(%s):", testTargetPath+"math.go"), func(t *testing.T) { // file level
			t.Run("Test-GetMinInMap()", func(t *testing.T) { // function level
				t.Parallel()
				TestGetMinInMap(t)
			})
			t.Run("Test-GetMaxInMap()", func(t *testing.T) { // function level
				t.Parallel()
				TestGetMaxInMap(t)
			})
		})
		t.Run(fmt.Sprintf("Test-Migration-(%s):", testTargetPath+"migration.go"), func(t *testing.T) { // file level
			t.Run("Test-GenerateMigrationFileName()", func(t *testing.T) { // function level
				t.Parallel()
				TestGenerateMigrationFileName(t)
			})
		})
		t.Run(fmt.Sprintf("Test-String-(%s):", testTargetPath+"string.go"), func(t *testing.T) { // file level
			t.Run("Test-JoinValues()", func(t *testing.T) { // function level
				t.Parallel()
				TestJoinValues(t)
			})
			t.Run("Test-ConvertCamelCaseToSenctenceCase()", func(t *testing.T) { // function level
				t.Parallel()
				TestConvertCamelCaseToSentenceCase(t)
			})
			t.Run("Test-IsStringIn()", func(t *testing.T) { // function level
				t.Parallel()
				TestIsStringIn(t)
			})
			t.Run("Test-IsEmailString()", func(t *testing.T) { // function level
				t.Parallel()
				TestIsEmailString(t)
			})
			t.Run("Test-IsAlphaNumberString()", func(t *testing.T) { // function level
				t.Parallel()
				TestIsAlphaNumberString(t)
			})
		})
	})
}
