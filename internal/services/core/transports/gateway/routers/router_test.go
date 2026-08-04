package routers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	gatewaycontract "github.com/HiIamJeff67/notezy-backend/contracts/gateway/v1"
	rootshelvesdto "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api/root-shelves"
	stationsdto "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api/stations"
	coreadapters "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/core/adapters"
	sharedtokens "github.com/HiIamJeff67/notezy-backend/shared/tokens"
)

func TestRouterValidatesRootShelfEnvelopeBeforeCallingService(t *testing.T) {
	t.Setenv("CORE_DELEGATION_AUDIENCE", "notezy-api-test")
	t.Setenv("CORE_DELEGATION_ISSUER", "notezy-gateway-test")
	t.Setenv("CORE_DELEGATION_SECRET", "test-delegation-secret")
	t.Setenv("JWT_ACCESS_TOKEN_SECRET_KEY", "test-access-secret")

	tokenString, err := coreadapters.IssueDelegationToken(
		"gateway",
		"83bdeac1-02de-42fe-a7a8-4e1a83174866",
		[]string{"Read"},
		rootshelvesdto.GetMyRootShelfByIdOperation,
		"request-id",
	)
	if err != nil {
		t.Fatalf("issue delegation token: %v", err)
	}

	requestBody := []byte(`{"version":"v1","operation":"root-shelf.get-my-root-shelf-by-id","metadata":{"requestId":"request-id"},"dto":"invalid"}`)

	router := newRouterForDelegationValidation()
	request := httptest.NewRequest(
		http.MethodPost,
		"/core/v1/root-shelves/get-by-id",
		bytes.NewReader(requestBody),
	)
	request.Header.Set("Authorization", "Bearer "+tokenString)
	accessToken, err := sharedtokens.GenerateAccessToken(
		"83bdeac1-02de-42fe-a7a8-4e1a83174866",
		sharedtokens.AccessTokenClaims{
			UserAgent: "test-agent",
		},
	)
	if err != nil {
		t.Fatalf("issue access token: %v", err)
	}
	request.AddCookie(&http.Cookie{Name: "accessToken", Value: *accessToken})
	request.Header.Set("User-Agent", "test-agent")
	request.Header.Set("Content-Type", "application/json")
	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, request)
	if responseRecorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, responseRecorder.Code)
	}
}

func TestRouterValidatesStationEnvelopeBeforeCallingService(t *testing.T) {
	t.Setenv("CORE_DELEGATION_AUDIENCE", "notezy-api-test")
	t.Setenv("CORE_DELEGATION_ISSUER", "notezy-gateway-test")
	t.Setenv("CORE_DELEGATION_SECRET", "test-delegation-secret")
	t.Setenv("JWT_ACCESS_TOKEN_SECRET_KEY", "test-access-secret")

	tokenString, err := coreadapters.IssueDelegationToken(
		"gateway",
		"83bdeac1-02de-42fe-a7a8-4e1a83174866",
		[]string{"Read"},
		stationsdto.GetMyStationByIdOperation,
		"request-id",
	)
	if err != nil {
		t.Fatalf("issue delegation token: %v", err)
	}

	requestBody := []byte(`{"version":"v1","operation":"station.get-my-station-by-id","metadata":{"requestId":"request-id"},"dto":"invalid"}`)

	router := newRouterForDelegationValidation()
	request := httptest.NewRequest(
		http.MethodPost,
		"/core/v1/stations/get-by-id",
		bytes.NewReader(requestBody),
	)
	request.Header.Set("Authorization", "Bearer "+tokenString)
	accessToken, err := sharedtokens.GenerateAccessToken(
		"83bdeac1-02de-42fe-a7a8-4e1a83174866",
		sharedtokens.AccessTokenClaims{
			UserAgent: "test-agent",
		},
	)
	if err != nil {
		t.Fatalf("issue access token: %v", err)
	}
	request.AddCookie(&http.Cookie{Name: "accessToken", Value: *accessToken})
	request.Header.Set("User-Agent", "test-agent")
	request.Header.Set("Content-Type", "application/json")
	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, request)
	if responseRecorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, responseRecorder.Code)
	}
}

func TestRouterRejectsDelegationForAnotherOperation(t *testing.T) {
	t.Setenv("CORE_DELEGATION_AUDIENCE", "notezy-api-test")
	t.Setenv("CORE_DELEGATION_ISSUER", "notezy-gateway-test")
	t.Setenv("CORE_DELEGATION_SECRET", "test-delegation-secret")
	t.Setenv("JWT_ACCESS_TOKEN_SECRET_KEY", "test-access-secret")

	tokenString, err := coreadapters.IssueDelegationToken(
		"gateway",
		"83bdeac1-02de-42fe-a7a8-4e1a83174866",
		nil,
		"root-shelf.create",
		"request-id",
	)
	if err != nil {
		t.Fatalf("issue delegation token: %v", err)
	}
	requestBody, err := json.Marshal(gatewaycontract.Request[struct{}]{
		Version:   gatewaycontract.Version,
		Operation: rootshelvesdto.GetMyRootShelfByIdOperation,
		Metadata: gatewaycontract.RequestMetadata{
			RequestId: "request-id",
		},
		Dto: struct{}{},
	})
	if err != nil {
		t.Fatalf("marshal Core service request: %v", err)
	}

	router := newRouterForDelegationValidation()
	request := httptest.NewRequest(
		http.MethodPost,
		"/core/v1/root-shelves/get-by-id",
		bytes.NewReader(requestBody),
	)
	request.Header.Set("Authorization", "Bearer "+tokenString)
	accessToken, err := sharedtokens.GenerateAccessToken(
		"83bdeac1-02de-42fe-a7a8-4e1a83174866",
		sharedtokens.AccessTokenClaims{
			UserAgent: "test-agent",
		},
	)
	if err != nil {
		t.Fatalf("issue access token: %v", err)
	}
	request.AddCookie(&http.Cookie{Name: "accessToken", Value: *accessToken})
	request.Header.Set("User-Agent", "test-agent")
	request.Header.Set("Content-Type", "application/json")
	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, request)
	if responseRecorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, responseRecorder.Code)
	}
}

func newRouterForDelegationValidation() *gin.Engine {
	return NewRouter(
		func(ctx *gin.Context) {
			ctx.Next()
		},
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
	)
}
