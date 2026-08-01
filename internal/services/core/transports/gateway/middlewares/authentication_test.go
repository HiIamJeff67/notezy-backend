package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	coreadapters "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/core/adapters"
	sharedtokens "github.com/HiIamJeff67/notezy-backend/internal/shared/tokens"
)

func TestAuthenticationMiddlewareValidatesForwardedAccessToken(t *testing.T) {
	t.Setenv("CORE_DELEGATION_AUDIENCE", "notezy-core-test")
	t.Setenv("CORE_DELEGATION_ISSUER", "notezy-gateway-test")
	t.Setenv("CORE_DELEGATION_SECRET", "test-delegation-secret")
	t.Setenv("JWT_ACCESS_TOKEN_SECRET_KEY", "test-access-secret")

	const userSubject = "83bdeac1-02de-42fe-a7a8-4e1a83174866"
	delegationToken, err := coreadapters.IssueDelegationToken("gateway", userSubject, nil, "station.get", "request-id")
	if err != nil {
		t.Fatalf("issue delegation token: %v", err)
	}
	accessToken, err := sharedtokens.GenerateAccessToken(
		userSubject,
		sharedtokens.AccessTokenClaims{UserAgent: "test-agent"},
	)
	if err != nil {
		t.Fatalf("generate access token: %v", err)
	}

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(DelegationAuthenticatedMiddleware(""), AuthMiddleware(nil))
	router.GET("/", func(ctx *gin.Context) {
		ctx.Status(http.StatusNoContent)
	})

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.Header.Set("Authorization", "Bearer "+delegationToken)
	request.Header.Set("User-Agent", "test-agent")
	request.AddCookie(&http.Cookie{Name: "accessToken", Value: *accessToken})
	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, request)
	if responseRecorder.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, responseRecorder.Code)
	}
}
