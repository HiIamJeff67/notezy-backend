package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	sharedtokens "github.com/HiIamJeff67/notezy-backend/internal/shared/tokens"
	types "github.com/HiIamJeff67/notezy-backend/internal/shared/types"
)

func TestAuthMiddlewareExtractsAccessTokenSubjectWithoutDataStore(t *testing.T) {
	t.Setenv("JWT_ACCESS_TOKEN_SECRET_KEY", "test-access-secret")
	accessToken, err := sharedtokens.GenerateAccessToken(
		"83bdeac1-02de-42fe-a7a8-4e1a83174866",
		sharedtokens.AccessTokenClaims{
			Name:      "notezy",
			Email:     "notezy@example.com",
			UserAgent: "test-agent",
		},
	)
	if err != nil {
		t.Fatalf("generate access token: %v", err)
	}

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(AuthMiddleware())
	router.GET("/", func(ctx *gin.Context) {
		value, exists := ctx.Get(types.ContextFieldName_User_PublicId.String())
		if !exists || value != "83bdeac1-02de-42fe-a7a8-4e1a83174866" {
			t.Fatalf("unexpected user public ID context: %#v", value)
		}
		ctx.Status(http.StatusNoContent)
	})

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.Header.Set("User-Agent", "test-agent")
	request.AddCookie(&http.Cookie{Name: "accessToken", Value: *accessToken})
	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, request)
	if responseRecorder.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, responseRecorder.Code)
	}
}

func TestAuthMiddlewareFallsBackToRefreshTokenSubject(t *testing.T) {
	t.Setenv("JWT_REFRESH_TOKEN_SECRET_KEY", "test-refresh-secret")
	refreshToken, err := sharedtokens.GenerateRefreshToken(
		"83bdeac1-02de-42fe-a7a8-4e1a83174866",
		sharedtokens.RefreshTokenClaims{
			Name:      "notezy",
			Email:     "notezy@example.com",
			UserAgent: "test-agent",
		},
	)
	if err != nil {
		t.Fatalf("generate refresh token: %v", err)
	}

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(AuthMiddleware())
	router.GET("/", func(ctx *gin.Context) {
		value, exists := ctx.Get(types.ContextFieldName_User_PublicId.String())
		if !exists || value != "83bdeac1-02de-42fe-a7a8-4e1a83174866" {
			t.Fatalf("unexpected user public ID context: %#v", value)
		}
		ctx.Status(http.StatusNoContent)
	})

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.Header.Set("User-Agent", "test-agent")
	request.AddCookie(&http.Cookie{Name: "refreshToken", Value: *refreshToken})
	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, request)
	if responseRecorder.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, responseRecorder.Code)
	}
}
