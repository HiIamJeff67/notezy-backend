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

	util "notezy-backend/app/util"
	test "notezy-backend/test"
)

/* ============================== Test Case Types ============================== */
type LoginRequestType = test.CommonRequestType[
	struct {
		Account  string `json:"account"`
		Password string `json:"password"`
	},
	test.CommonCookiesType,
]
type LoginResponseType = test.CommonResponseType[
	struct {
		AccessToken string    `json:"accessToken"`
		UpdatedAt   time.Time `json:"updatedAt"`
	},
	test.CommonCookiesType,
]
type LoginE2ETestCase = test.E2ETestCase[
	LoginRequestType,
	LoginResponseType,
]

/* ============================== Test Data Path & Some Constants ============================== */

const (
	loginTestdataPath = "testdata/login_testdata/"
	loginRoute        = testAuthRouteNamespace + "/login"
)

/* ============================== Interface & Instance ============================== */

type LoginE2ETesterInterface interface {
	getLoginTestDataAndResponse(
		t *testing.T,
		method string,
		loginTestdataPath string,
	) (
		w *httptest.ResponseRecorder,
		testCase LoginE2ETestCase,
		res LoginResponseType,
		cookieMap map[string]string,
	)
	TestLoginValidTestAccountByName(t *testing.T)
	TestLoginValidTestAccountByEmail(t *testing.T)
}

type LoginE2ETester struct {
	router *gin.Engine
}

func NewLoginE2ETester(router *gin.Engine) LoginE2ETesterInterface {
	if router == nil {
		return nil
	}
	return &LoginE2ETester{
		router: router,
	}
}

/* ============================== Auxiliary Functions ============================== */

func (et *LoginE2ETester) getLoginTestDataAndResponse(
	t *testing.T,
	method string,
	loginTestdataPath string,
) (
	w *httptest.ResponseRecorder,
	testCase LoginE2ETestCase,
	res LoginResponseType,
	cookieMap map[string]string,
) {
	if et == nil || et.router == nil {
		t.Fatalf("loginE2ETester or router is nil")
	}

	testCase = test.LoadTestCase[LoginE2ETestCase](
		t, loginTestdataPath,
	)

	jsonBody, _ := json.Marshal(testCase.Request.Body)
	req, err := http.NewRequest(
		method,
		loginRoute,
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		t.Errorf("failed to marshal json body, maybe something went wrong in testdata")
	}

	req.Header.Set("Content-Type", "application/json")
	if ua := testCase.Request.Header.UserAgent; ua != nil {
		req.Header.Set("User-Agent", *ua)
	}

	w = httptest.NewRecorder()
	et.router.ServeHTTP(w, req)
	if err := json.Unmarshal(w.Body.Bytes(), &res.Body); err != nil {
		t.Errorf("failed to unmarshal response body: %v, body: %s", err, w.Body.String())
	}

	cookies := w.Result().Cookies()
	cookieMap = make(map[string]string)
	for _, c := range cookies {
		cookieMap[c.Name] = c.Value
	}
	return w, testCase, res, cookieMap
}

/* ============================== Test Cases ============================== */

func (et *LoginE2ETester) TestLoginValidTestAccountByName(t *testing.T) {
	if et.router == nil {
		return
	}

	w, testCase, res, cookieMap := et.getLoginTestDataAndResponse(
		t, "POST", loginTestdataPath+"valid_test_account_by_name_testdata.json",
	)

	// check status code
	if w.Code != testCase.Response.HTTPStatusCode {
		t.Errorf("expected http status code to be %d, got %d", testCase.Response.HTTPStatusCode, w.Code)
	}

	// check the body
	if err := json.Unmarshal(w.Body.Bytes(), &res.Body); err != nil {
		t.Errorf("failed to unmarshal response body: %v, body: %s", err, w.Body.String())
	}

	if !res.Body.Success {
		t.Errorf("expected body.success to be true, got false")
	}

	if res.Body.Data == nil {
		t.Errorf("expected response data to be not nil, got nil")
	}
	if len(strings.ReplaceAll(res.Body.Data.AccessToken, " ", "")) == 0 {
		t.Errorf("expected body.data.accessToken to be exist, got nil")
	}

	now := time.Now()
	if !util.IsTimeWithinDelta(res.Body.Data.UpdatedAt, now, 10*time.Second) {
		t.Errorf("expected body.data.createdAt to be %v (within tolerable time duration of %v), got %v", testCase.Response.Body.Data.UpdatedAt, 10*time.Second, now)
	}

	if res.Body.Exception != nil {
		t.Errorf("expected body.exception to be nil, got not %v", res.Body.Exception)
	}

	// check the accessToken in cookies
	if _, ok := cookieMap["accessToken"]; !ok {
		t.Errorf("expected cookie.accessToken to be set, got nil")
	}

	// check the refreshToken in cookies
	if _, ok := cookieMap["refreshToken"]; !ok {
		t.Errorf("expected cookie.refreshToken to be set, got nil")
	}
}

func (et *LoginE2ETester) TestLoginValidTestAccountByEmail(t *testing.T) {
	if et.router == nil {
		return
	}

	w, testCase, res, cookieMap := et.getLoginTestDataAndResponse(
		t, "POST", loginTestdataPath+"valid_test_account_by_email_testdata.json",
	)

	// check status code
	if w.Code != testCase.Response.HTTPStatusCode {
		t.Errorf("expected http status code to be %d, got %d", testCase.Response.HTTPStatusCode, w.Code)
	}

	// check the body
	if err := json.Unmarshal(w.Body.Bytes(), &res.Body); err != nil {
		t.Errorf("failed to unmarshal response body: %v, body: %s", err, w.Body.String())
	}

	if !res.Body.Success {
		t.Errorf("expected body.success to be true, got false")
	}

	if res.Body.Data == nil {
		t.Errorf("expected response data to be not nil, got nil")
	}
	if len(strings.ReplaceAll(res.Body.Data.AccessToken, " ", "")) == 0 {
		t.Errorf("expected body.data.accessToken to be exist, got nil")
	}

	now := time.Now()
	if !util.IsTimeWithinDelta(res.Body.Data.UpdatedAt, now, 10*time.Second) {
		t.Errorf("expected body.data.createdAt to be %v (within tolerable time duration of %v), got %v", testCase.Response.Body.Data.UpdatedAt, 10*time.Second, now)
	}

	if res.Body.Exception != nil {
		t.Errorf("expected body.exception to be nil, got not %v", res.Body.Exception)
	}

	// check the accessToken in cookies
	if _, ok := cookieMap["accessToken"]; !ok {
		t.Errorf("expected cookie.accessToken to be set, got nil")
	}

	// check the refreshToken in cookies
	if _, ok := cookieMap["refreshToken"]; !ok {
		t.Errorf("expected cookie.refreshToken to be set, got nil")
	}
}

// may test using access token or refresh token to login
