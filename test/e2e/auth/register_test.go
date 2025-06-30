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

	"notezy-backend/app/util"
	test "notezy-backend/test"
)

/* ============================== Test Case Type ============================== */

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

/* ============================== Test Data Path & Some Constants ============================== */

const (
	testdataPath   = "testdata/register_testdata/"
	routeNamespace = "/testRegisterRoute"
	registerRoute  = routeNamespace + "/auth/register"
)

/* ============================== Interface & Instance ============================== */

type RegisterE2ETesterInterface interface {
	TestRegisterValidTestAccount(t *testing.T)
}

type registerE2ETester struct {
	router *gin.Engine
}

func NewRegisterE2ETester(router *gin.Engine) RegisterE2ETesterInterface {
	if router == nil {
		return nil
	}
	return &registerE2ETester{
		router: router,
	}
}

/* ============================== Test Procedures ============================== */

func (et *registerE2ETester) TestRegisterValidTestAccount(t *testing.T) {
	if et.router == nil {
		return
	}

	testCase := test.LoadTestCase[RegisterE2ETestCase](
		t, testdataPath+"valid_test_account_testdata.json",
	)

	jsonBody, _ := json.Marshal(testCase.Request.Body)
	req, err := http.NewRequest(
		"POST",
		registerRoute,
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		t.Errorf("failed to marshal json body, maybe something went wrong in testdata")
	}

	req.Header.Set("Content-Type", "application/json")
	if ua := testCase.Request.Header.UserAgent; ua != nil {
		req.Header.Set("User-Agent", *ua)
	}

	w := httptest.NewRecorder()
	et.router.ServeHTTP(w, req)

	// check status code
	if w.Code != testCase.Response.HTTPStatusCode {
		t.Errorf("expected status %d, got %d, body: %s", testCase.Response.HTTPStatusCode, w.Code, w.Body.String())
	}

	var res RegisterResponseType
	if err := json.Unmarshal(w.Body.Bytes(), &res.Body); err != nil {
		t.Errorf("failed to unmarshal response body: %v, body: %s", err, w.Body.String())
	}

	if res.Body.Data == nil {
		t.Errorf("response data does not exist")
	}

	if res.Body.Data != nil {
		if !res.Body.Success {
			t.Errorf("expected body/success to be true, got false")
		}
		if len(strings.ReplaceAll(res.Body.Data.AccessToken, " ", "")) > 0 {
			t.Errorf("expected accessToken %s, got %s", testCase.Response.Body.Data.AccessToken, res.Body.Data.AccessToken)
		}
		if !util.IsTimeWithinDelta(res.Body.Data.CreatedAt, testCase.Response.Body.Data.CreatedAt, 10*time.Second) {
			t.Errorf("expected createdAt %v (within tolerable time duration of %v), got %v", testCase.Response.Body.Data.CreatedAt, 10*time.Second, res.Body.Data.CreatedAt)
		}
	}

}
