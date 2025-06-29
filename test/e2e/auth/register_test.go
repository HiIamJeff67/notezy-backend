package authe2etest

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	testroutes "notezy-backend/app/routes/test_routes"
	test "notezy-backend/test"
)

type RegisterRequestType = test.CommonRequestType[
	struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	},
	test.CommonCookiesType,
]
type RegisterResponseType = test.CommonResponseType[
	struct {
		AccessToken string    `json:"accessToken"`
		CreatedAt   time.Time `json:"createdAt"`
	},
	test.CommonCookiesType,
]
type RegisterE2ETestCase = test.E2ETestCase[
	RegisterRequestType,
	RegisterResponseType,
]

func TestRegisterRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	testRouterGroup := router.Group("/testRegisterRoute")
	testroutes.ConfigureTestAuthRoutes(testRouterGroup)

	// 準備 request body
	body := map[string]string{
		"email":    "test@example.com",
		"password": "test123!",
	}
	jsonBody, _ := json.Marshal(body)

	// 建立 request
	req, _ := http.NewRequest("POST", "/testRegisterRoute/auth/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// 建立 recorder
	w := httptest.NewRecorder()

	// 執行
	router.ServeHTTP(w, req)

	// 驗證
	if w.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d, body: %s", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), "success") {
		t.Errorf("unexpected response body: %s", w.Body.String())
	}
}
