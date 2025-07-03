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
	TestRegisterValidUserAccount(t *testing.T)
	TestRegisterNoName(t *testing.T)
	TestRegisterNameWithoutNumber(t *testing.T)
	TestRegisterShortName(t *testing.T)
	TestRegisterInvalidEmail(t *testing.T)
	TestRegisterShortPassword(t *testing.T)
	TestRegisterPasswordWithoutLowerCaseLetter(t *testing.T)
	TestRegisterPasswordWithoutUpperCaseLetter(t *testing.T)
	TestRegisterPasswordWithoutNumber(t *testing.T)
	TestRegisterPasswordWithoutSign(t *testing.T)
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

/* ============================== Auxiliary Functions ============================== */

func (et *registerE2ETester) getRegisterTestdataAndResponse(
	t *testing.T,
	method string,
	testdataPath string,
) (
	w *httptest.ResponseRecorder,
	testCase RegisterE2ETestCase,
	res RegisterResponseType,
	cookieMap map[string]string,
) {
	if et == nil || et.router == nil {
		t.Fatal("registerE2ETester or router is nil")
	}

	testCase = test.LoadTestCase[RegisterE2ETestCase](
		t, testdataPath,
	)

	jsonBody, _ := json.Marshal(testCase.Request.Body)
	req, err := http.NewRequest(
		method,
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

/* ============================== Test Procedures ============================== */

func (et *registerE2ETester) TestRegisterValidTestAccount(t *testing.T) {
	if et.router == nil {
		return
	}

	w, testCase, res, cookieMap := et.getRegisterTestdataAndResponse(
		t, "POST", testdataPath+"valid_test_account_testdata.json",
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
	if !util.IsTimeWithinDelta(res.Body.Data.CreatedAt, now, 10*time.Second) {
		t.Errorf("expected body.data.createdAt to be %v (within tolerable time duration of %v), got %v", testCase.Response.Body.Data.CreatedAt, 10*time.Second, now)
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

func (et *registerE2ETester) TestRegisterValidUserAccount(t *testing.T) {
	if et.router == nil {
		return
	}

	w, testCase, res, cookieMap := et.getRegisterTestdataAndResponse(
		t, "POST", testdataPath+"valid_user_account_testdata.json",
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
		t.Errorf("expected response data to be exist, got nil")
	}
	if len(strings.ReplaceAll(res.Body.Data.AccessToken, " ", "")) == 0 {
		t.Errorf("expected body.data.accessToken to be exist, got nil")
	}
	now := time.Now()
	if !util.IsTimeWithinDelta(res.Body.Data.CreatedAt, now, 10*time.Second) {
		t.Errorf("expected body.data.createdAt to be %v (within tolerable time duration of %v), got %v", testCase.Response.Body.Data.CreatedAt, 10*time.Second, now)
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

func (et *registerE2ETester) TestRegisterNoName(t *testing.T) {
	if et.router == nil {
		return
	}

	w, testCase, res, cookieMap := et.getRegisterTestdataAndResponse(
		t, "POST", testdataPath+"no_name_testdata.json",
	)

	// check status code
	if w.Code != testCase.Response.HTTPStatusCode {
		t.Errorf("expected http status code to be %d, got %d", testCase.Response.HTTPStatusCode, w.Code)
	}

	// check the body
	if err := json.Unmarshal(w.Body.Bytes(), &res.Body); err != nil {
		t.Errorf("failed to unmarshal response body: %v, body: %s", err, w.Body.String())
	}

	if res.Body.Success {
		t.Errorf("expected body.success to be false, got true")
	}

	if res.Body.Data != nil {
		t.Errorf("expected response data to be nil, got %v", res.Body.Data)
	}

	if res.Body.Exception == nil {
		t.Errorf("expected body.exception to be exist, got nil")
	}

	// check the accessToken in cookies
	if val, ok := cookieMap["accessToken"]; ok {
		t.Errorf("expected cookie.accessToken to be not set, got %v", val)
	}

	// check the refreshToken in cookies
	if val, ok := cookieMap["refreshToken"]; ok {
		t.Errorf("expected cookie.refreshToken to be not set, got %v", val)
	}
}

func (et *registerE2ETester) TestRegisterNameWithoutNumber(t *testing.T) {
	if et.router == nil {
		return
	}

	w, testCase, res, cookieMap := et.getRegisterTestdataAndResponse(
		t, "POST", testdataPath+"name_without_number_testdata.json",
	)

	// check status code
	if w.Code != testCase.Response.HTTPStatusCode {
		t.Errorf("expected http status code to be %d, got %d", testCase.Response.HTTPStatusCode, w.Code)
	}

	// check the body
	if err := json.Unmarshal(w.Body.Bytes(), &res.Body); err != nil {
		t.Errorf("failed to unmarshal response body: %v, body: %s", err, w.Body.String())
	}

	if res.Body.Success {
		t.Errorf("expected body.success to be false, got true")
	}

	if res.Body.Data != nil {
		t.Errorf("expected response data to be nil, got %v", res.Body.Data)
	}

	if res.Body.Exception == nil {
		t.Errorf("expected body.exception to be exist, got nil")
	}

	// check the accessToken in cookies
	if val, ok := cookieMap["accessToken"]; ok {
		t.Errorf("expected cookie.accessToken to be not set, got %v", val)
	}

	// check the refreshToken in cookies
	if val, ok := cookieMap["refreshToken"]; ok {
		t.Errorf("expected cookie.refreshToken to be not set, got %v", val)
	}
}

func (et *registerE2ETester) TestRegisterShortName(t *testing.T) {
	if et.router == nil {
		return
	}

	w, testCase, res, cookieMap := et.getRegisterTestdataAndResponse(
		t, "POST", testdataPath+"short_name_testdata.json",
	)

	// check status code
	if w.Code != testCase.Response.HTTPStatusCode {
		t.Errorf("expected http status code to be %d, got %d", testCase.Response.HTTPStatusCode, w.Code)
	}

	// check the body
	if err := json.Unmarshal(w.Body.Bytes(), &res.Body); err != nil {
		t.Errorf("failed to unmarshal response body: %v, body: %s", err, w.Body.String())
	}

	if res.Body.Success {
		t.Errorf("expected body.success to be false, got true")
	}

	if res.Body.Data != nil {
		t.Errorf("expected response data to be nil, got %v", res.Body.Data)
	}

	if res.Body.Exception == nil {
		t.Errorf("expected body.exception to be exist, got nil")
	}

	// check the accessToken in cookies
	if val, ok := cookieMap["accessToken"]; ok {
		t.Errorf("expected cookie.accessToken to be not set, got %v", val)
	}

	// check the refreshToken in cookies
	if val, ok := cookieMap["refreshToken"]; ok {
		t.Errorf("expected cookie.refreshToken to be not set, got %v", val)
	}
}

func (et *registerE2ETester) TestRegisterInvalidEmail(t *testing.T) {
	if et.router == nil {
		return
	}

	w, testCase, res, cookieMap := et.getRegisterTestdataAndResponse(
		t, "POST", testdataPath+"invalid_email_testdata.json",
	)

	// check status code
	if w.Code != testCase.Response.HTTPStatusCode {
		t.Errorf("expected http status code to be %d, got %d", testCase.Response.HTTPStatusCode, w.Code)
	}

	// check the body
	if err := json.Unmarshal(w.Body.Bytes(), &res.Body); err != nil {
		t.Errorf("failed to unmarshal response body: %v, body: %s", err, w.Body.String())
	}

	if res.Body.Success {
		t.Errorf("expected body.success to be false, got true")
	}

	if res.Body.Data != nil {
		t.Errorf("expected response data to be nil, got %v", res.Body.Data)
	}

	if res.Body.Exception == nil {
		t.Errorf("expected body.exception to be exist, got nil")
	}

	// check the accessToken in cookies
	if val, ok := cookieMap["accessToken"]; ok {
		t.Errorf("expected cookie.accessToken to be not set, got %v", val)
	}

	// check the refreshToken in cookies
	if val, ok := cookieMap["refreshToken"]; ok {
		t.Errorf("expected cookie.refreshToken to be not set, got %v", val)
	}
}

func (et *registerE2ETester) TestRegisterShortPassword(t *testing.T) {
	if et.router == nil {
		return
	}

	w, testCase, res, cookieMap := et.getRegisterTestdataAndResponse(
		t, "POST", testdataPath+"short_password_testdata.json",
	)

	// check status code
	if w.Code != testCase.Response.HTTPStatusCode {
		t.Errorf("expected http status code to be %d, got %d", testCase.Response.HTTPStatusCode, w.Code)
	}

	// check the body
	if err := json.Unmarshal(w.Body.Bytes(), &res.Body); err != nil {
		t.Errorf("failed to unmarshal response body: %v, body: %s", err, w.Body.String())
	}

	if res.Body.Success {
		t.Errorf("expected body.success to be false, got true")
	}

	if res.Body.Data != nil {
		t.Errorf("expected response data to be nil, got %v", res.Body.Data)
	}

	if res.Body.Exception == nil {
		t.Errorf("expected body.exception to be exist, got nil")
	}

	// check the accessToken in cookies
	if val, ok := cookieMap["accessToken"]; ok {
		t.Errorf("expected cookie.accessToken to be not set, got %v", val)
	}

	// check the refreshToken in cookies
	if val, ok := cookieMap["refreshToken"]; ok {
		t.Errorf("expected cookie.refreshToken to be not set, got %v", val)
	}
}

func (et *registerE2ETester) TestRegisterPasswordWithoutLowerCaseLetter(t *testing.T) {
	if et.router == nil {
		return
	}

	w, testCase, res, cookieMap := et.getRegisterTestdataAndResponse(
		t, "POST", testdataPath+"password_without_lower_case_letter_testdata.json",
	)

	// check status code
	if w.Code != testCase.Response.HTTPStatusCode {
		t.Errorf("expected http status code to be %d, got %d", testCase.Response.HTTPStatusCode, w.Code)
	}

	// check the body
	if err := json.Unmarshal(w.Body.Bytes(), &res.Body); err != nil {
		t.Errorf("failed to unmarshal response body: %v, body: %s", err, w.Body.String())
	}

	if res.Body.Success {
		t.Errorf("expected body.success to be false, got true")
	}

	if res.Body.Data != nil {
		t.Errorf("expected response data to be nil, got %v", res.Body.Data)
	}

	if res.Body.Exception == nil {
		t.Errorf("expected body.exception to be exist, got nil")
	}

	// check the accessToken in cookies
	if val, ok := cookieMap["accessToken"]; ok {
		t.Errorf("expected cookie.accessToken to be not set, got %v", val)
	}

	// check the refreshToken in cookies
	if val, ok := cookieMap["refreshToken"]; ok {
		t.Errorf("expected cookie.refreshToken to be not set, got %v", val)
	}
}

func (et *registerE2ETester) TestRegisterPasswordWithoutUpperCaseLetter(t *testing.T) {
	if et.router == nil {
		return
	}

	w, testCase, res, cookieMap := et.getRegisterTestdataAndResponse(
		t, "POST", testdataPath+"password_without_upper_case_letter_testdata.json",
	)

	// check status code
	if w.Code != testCase.Response.HTTPStatusCode {
		t.Errorf("expected http status code to be %d, got %d", testCase.Response.HTTPStatusCode, w.Code)
	}

	// check the body
	if err := json.Unmarshal(w.Body.Bytes(), &res.Body); err != nil {
		t.Errorf("failed to unmarshal response body: %v, body: %s", err, w.Body.String())
	}

	if res.Body.Success {
		t.Errorf("expected body.success to be false, got true")
	}

	if res.Body.Data != nil {
		t.Errorf("expected response data to be nil, got %v", res.Body.Data)
	}

	if res.Body.Exception == nil {
		t.Errorf("expected body.exception to be exist, got nil")
	}

	// check the accessToken in cookies
	if val, ok := cookieMap["accessToken"]; ok {
		t.Errorf("expected cookie.accessToken to be not set, got %v", val)
	}

	// check the refreshToken in cookies
	if val, ok := cookieMap["refreshToken"]; ok {
		t.Errorf("expected cookie.refreshToken to be not set, got %v", val)
	}
}

func (et *registerE2ETester) TestRegisterPasswordWithoutNumber(t *testing.T) {
	if et.router == nil {
		return
	}

	w, testCase, res, cookieMap := et.getRegisterTestdataAndResponse(
		t, "POST", testdataPath+"password_without_number_testdata.json",
	)

	// check status code
	if w.Code != testCase.Response.HTTPStatusCode {
		t.Errorf("expected http status code to be %d, got %d", testCase.Response.HTTPStatusCode, w.Code)
	}

	// check the body
	if err := json.Unmarshal(w.Body.Bytes(), &res.Body); err != nil {
		t.Errorf("failed to unmarshal response body: %v, body: %s", err, w.Body.String())
	}

	if res.Body.Success {
		t.Errorf("expected body.success to be false, got true")
	}

	if res.Body.Data != nil {
		t.Errorf("expected response data to be nil, got %v", res.Body.Data)
	}

	if res.Body.Exception == nil {
		t.Errorf("expected body.exception to be exist, got nil")
	}

	// check the accessToken in cookies
	if val, ok := cookieMap["accessToken"]; ok {
		t.Errorf("expected cookie.accessToken to be not set, got %v", val)
	}

	// check the refreshToken in cookies
	if val, ok := cookieMap["refreshToken"]; ok {
		t.Errorf("expected cookie.refreshToken to be not set, got %v", val)
	}
}

func (et *registerE2ETester) TestRegisterPasswordWithoutSign(t *testing.T) {
	if et.router == nil {
		return
	}

	w, testCase, res, cookieMap := et.getRegisterTestdataAndResponse(
		t, "POST", testdataPath+"password_without_sign_testdata.json",
	)

	// check status code
	if w.Code != testCase.Response.HTTPStatusCode {
		t.Errorf("expected http status code to be %d, got %d", testCase.Response.HTTPStatusCode, w.Code)
	}

	// check the body
	if err := json.Unmarshal(w.Body.Bytes(), &res.Body); err != nil {
		t.Errorf("failed to unmarshal response body: %v, body: %s", err, w.Body.String())
	}

	if res.Body.Success {
		t.Errorf("expected body.success to be false, got true")
	}

	if res.Body.Data != nil {
		t.Errorf("expected response data to be nil, got %v", res.Body.Data)
	}

	if res.Body.Exception == nil {
		t.Errorf("expected body.exception to be exist, got nil")
	}

	// check the accessToken in cookies
	if val, ok := cookieMap["accessToken"]; ok {
		t.Errorf("expected cookie.accessToken to be not set, got %v", val)
	}

	// check the refreshToken in cookies
	if val, ok := cookieMap["refreshToken"]; ok {
		t.Errorf("expected cookie.refreshToken to be not set, got %v", val)
	}
}
