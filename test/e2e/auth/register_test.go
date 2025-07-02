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

/* ============================== Auxilary Functions ============================== */

func (et *registerE2ETester) getRegisterTestdataAndResponse(
	t *testing.T,
	method string,
	testdataPath string,
) (*httptest.ResponseRecorder, RegisterE2ETestCase, RegisterResponseType) {
	if et.router == nil {
		var zeroTestCase RegisterE2ETestCase
		var zeroResponse RegisterResponseType
		return nil, zeroTestCase, zeroResponse
	}

	testCase := test.LoadTestCase[RegisterE2ETestCase](
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

	w := httptest.NewRecorder()
	et.router.ServeHTTP(w, req)

	var res RegisterResponseType
	if err := json.Unmarshal(w.Body.Bytes(), &res.Body); err != nil {
		t.Errorf("failed to unmarshal response body: %v, body: %s", err, w.Body.String())
	}

	return w, testCase, res
}

/* ============================== Test Procedures ============================== */

func (et *registerE2ETester) TestRegisterValidTestAccount(t *testing.T) {
	if et.router == nil {
		return
	}

	w, testCase, res := et.getRegisterTestdataAndResponse(
		t, "POST", testdataPath+"valid_test_account_testdata.json",
	)

	// check status code
	if w.Code != testCase.Response.HTTPStatusCode {
		t.Errorf("expected http status code to be %d, got %d, body: %s", testCase.Response.HTTPStatusCode, w.Code, w.Body.String())
	}

	// check the body
	if err := json.Unmarshal(w.Body.Bytes(), &res.Body); err != nil {
		t.Errorf("failed to unmarshal response body: %v, body: %s", err, w.Body.String())
	}

	if !res.Body.Success {
		t.Errorf("expected body.success to be true, got false")
	}

	if res.Body.Data == nil {
		t.Errorf("response data does not exist")
	}
	if len(strings.ReplaceAll(res.Body.Data.AccessToken, " ", "")) > 0 {
		t.Errorf("expected body.data.accessToken to be %s, got %s", testCase.Response.Body.Data.AccessToken, res.Body.Data.AccessToken)
	}
	if !util.IsTimeWithinDelta(res.Body.Data.CreatedAt, testCase.Response.Body.Data.CreatedAt, 10*time.Second) {
		t.Errorf("expected body.data.createdAt to be %v (within tolerable time duration of %v), got %v", testCase.Response.Body.Data.CreatedAt, 10*time.Second, res.Body.Data.CreatedAt)
	}

	if res.Body.Exception != nil {
		t.Errorf("expected body.exception to be nil, got not %v", res.Body.Exception)
	}

	cookies := w.Result().Cookies()
	cookieMap := make(map[string]string)
	for _, c := range cookies {
		cookieMap[c.Name] = c.Value
	}
	// check the accessToken in cookies
	if _, ok := cookieMap["accessToken"]; !ok {
		t.Errorf("expected cookie.accessToken to be set")
	}

	// check the refreshToken in cookies
	if _, ok := cookieMap["refreshToken"]; !ok {
		t.Errorf("expected cookie.refreshToken to be set")
	}
}

func (et *registerE2ETester) TestRegisterValidUserAccount(t *testing.T) {

}
