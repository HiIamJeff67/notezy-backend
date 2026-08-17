package logs

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestNewGinRouterColorsDevelopmentStatusCodes(t *testing.T) {
	t.Setenv("NOTEGIC_LOG_FORMAT", "console")
	previousWriter := gin.DefaultWriter
	gin.DefaultWriter = &bytes.Buffer{}
	defer func() {
		gin.DefaultWriter = previousWriter
		gin.DisableConsoleColor()
	}()

	router := WithGinLogger(gin.New())
	router.GET("/ok", func(ctx *gin.Context) { ctx.Status(http.StatusOK) })
	router.GET("/failed", func(ctx *gin.Context) { ctx.Status(http.StatusInternalServerError) })

	for _, path := range []string{"/ok", "/failed"} {
		request := httptest.NewRequest(http.MethodGet, path, nil)
		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)
	}

	output := gin.DefaultWriter.(*bytes.Buffer).String()
	if !strings.Contains(output, "\x1b[") || !strings.Contains(output, "200") || !strings.Contains(output, "500") {
		t.Fatalf("expected colored 200/500 access logs, got %q", output)
	}
}
