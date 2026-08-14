package authe2etest

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	cookies "github.com/HiIamJeff67/notezy-backend/shared/cookies"
	platformdatabase "github.com/HiIamJeff67/notezy-backend/shared/platform/database"

	testroutes "github.com/HiIamJeff67/notezy-backend/internal/clientgateway/transports/api/routes/testroutes"
	coreadapters "github.com/HiIamJeff67/notezy-backend/internal/clientgateway/transports/core/adapters"
	"github.com/gin-gonic/gin"
)

const testAuthRouteNamespace = "/testRoute/auth"

func TestAuthE2E(t *testing.T) {
	databaseConfig, err := platformdatabase.LoadConfig()
	if err != nil {
		t.Skipf("auth E2E test requires database configuration: %v", err)
	}
	db, err := platformdatabase.Connect(databaseConfig)
	if err != nil {
		t.Skipf("auth E2E test requires an available database: %v", err)
	}
	t.Cleanup(func() {
		_ = platformdatabase.Disconnect(db)
	})

	gin.SetMode(gin.TestMode)
	router := gin.New()
	testroutes.ConfigureTestAuthRoutes(
		router.Group(testAuthRouteNamespace),
		coreadapters.NewCoreAdapter("http://127.0.0.1:7778", 10*time.Second),
		cookies.New(cookies.Config{
			Name:     cookies.ValidCookieName_AccessToken,
			Path:     "/",
			Duration: 30 * time.Minute,
			HTTPOnly: true,
			SameSite: http.SameSiteLaxMode,
		}),
		cookies.New(cookies.Config{
			Name:     cookies.ValidCookieName_RefreshToken,
			Path:     "/",
			Duration: 14 * 24 * time.Hour,
			HTTPOnly: true,
			SameSite: http.SameSiteStrictMode,
		}),
		nil,
	)
	server := httptest.NewServer(router)
	t.Cleanup(server.Close)
	client := server.Client()

	t.Run("register", func(t *testing.T) {
		registerE2ETester := NewRegisterE2ETester(server.URL, client)
		if registerE2ETester == nil {
			t.Fatal("NewRegisterE2ETester returned nil, router may be nil")
		}

		t.Run("valid_test_account", func(t *testing.T) {
			registerE2ETester.TestRegisterValidTestAccount(t)
		})
		t.Run("valid_user_account", func(t *testing.T) {
			registerE2ETester.TestRegisterValidUserAccount(t)
		})
		t.Run("no_name", func(t *testing.T) {
			registerE2ETester.TestRegisterNoName(t)
		})
		t.Run("name_without_number", func(t *testing.T) {
			registerE2ETester.TestRegisterNameWithoutNumber(t)
		})
		t.Run("short_name", func(t *testing.T) {
			registerE2ETester.TestRegisterShortName(t)
		})
		t.Run("invalid_email", func(t *testing.T) {
			registerE2ETester.TestRegisterInvalidEmail(t)
		})
		t.Run("short_password", func(t *testing.T) {
			registerE2ETester.TestRegisterShortPassword(t)
		})
		t.Run("password_without_lower_case_letter", func(t *testing.T) {
			registerE2ETester.TestRegisterPasswordWithoutLowerCaseLetter(t)
		})
		t.Run("password_without_upper_case_letter", func(t *testing.T) {
			registerE2ETester.TestRegisterPasswordWithoutUpperCaseLetter(t)
		})
		t.Run("password_without_number", func(t *testing.T) {
			registerE2ETester.TestRegisterPasswordWithoutNumber(t)
		})
		t.Run("password_without_sign", func(t *testing.T) {
			registerE2ETester.TestRegisterPasswordWithoutSign(t)
		})
	})

	t.Run("login", func(t *testing.T) {
		loginE2ETester := NewLoginE2ETester(server.URL, client)
		if loginE2ETester == nil {
			t.Fatal("NewLoginE2ETester returned nil, router may be nil")
		}

		t.Run("valid_test_account_by_name", func(t *testing.T) {
			loginE2ETester.TestLoginValidTestAccountByName(t)
		})
		t.Run("valid_test_account_by_email", func(t *testing.T) {
			loginE2ETester.TestLoginValidTestAccountByEmail(t)
		})
	})
}
