package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	contexts "github.com/HiIamJeff67/notezy-backend/internal/core/contexts"
	coreadapters "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/core/adapters"
)

func TestDelegationMiddlewareSetsRoutePermissions(t *testing.T) {
	t.Setenv("CORE_DELEGATION_AUDIENCE", "notezy-api-test")
	t.Setenv("CORE_DELEGATION_ISSUER", "notezy-gateway-test")
	t.Setenv("CORE_DELEGATION_SECRET", "test-delegation-secret")

	tokenString, err := coreadapters.IssueDelegationToken(
		"gateway",
		"",
		[]string{"Read", "Write"},
		"station.update",
		"request-id",
	)
	if err != nil {
		t.Fatalf("issue delegation token: %v", err)
	}

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(DelegationMiddleware(""))
	router.GET("/", func(ctx *gin.Context) {
		allowedPermissions, exception := contexts.GetAllowedPermissions(ctx.Request.Context())
		if exception != nil || len(allowedPermissions) != 2 {
			t.Fatal("expected delegated route permissions")
		}
		ctx.Status(http.StatusNoContent)
	})

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.Header.Set("Authorization", "Bearer "+tokenString)
	request.Header.Set("X-Request-Id", "request-id")
	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, request)
	if responseRecorder.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, responseRecorder.Code)
	}
}

func TestDelegationAuthenticatedMiddlewareRequiresUserSubject(t *testing.T) {
	t.Setenv("CORE_DELEGATION_AUDIENCE", "notezy-api-test")
	t.Setenv("CORE_DELEGATION_ISSUER", "notezy-gateway-test")
	t.Setenv("CORE_DELEGATION_SECRET", "test-delegation-secret")

	tokenString, err := coreadapters.IssueDelegationToken(
		"gateway",
		"",
		nil,
		"station.update",
		"request-id",
	)
	if err != nil {
		t.Fatalf("issue delegation token: %v", err)
	}

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(DelegationAuthenticatedMiddleware(""))
	router.GET("/", func(ctx *gin.Context) {
		ctx.Status(http.StatusNoContent)
	})

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.Header.Set("Authorization", "Bearer "+tokenString)
	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, request)
	if responseRecorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, responseRecorder.Code)
	}
}
