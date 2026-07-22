package binders

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	dtos "github.com/HiIamJeff67/notezy-backend/app/dtos"
	types "github.com/HiIamJeff67/notezy-backend/shared/types"
)

func TestBindUpsertMyRootShelfPermissionParsesURIUUIDs(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userId := uuid.New()
	rootShelfId := uuid.New()
	userPublicId := uuid.New()

	var capturedReqDto *dtos.UpsertMyRootShelfPermissionReqDto

	router := gin.New()
	router.Use(func(ctx *gin.Context) {
		ctx.Set(types.ContextFieldName_User_Id.String(), userId.String())
		ctx.Next()
	})
	router.PUT(
		"/rootShelf/:rootShelfId/permissions/:userPublicId",
		NewRootShelfBinder().BindUpsertMyRootShelfPermission(
			func(ctx *gin.Context, reqDto *dtos.UpsertMyRootShelfPermissionReqDto) {
				capturedReqDto = reqDto
				ctx.Status(http.StatusNoContent)
			},
		),
	)

	request := httptest.NewRequest(
		http.MethodPut,
		fmt.Sprintf("/rootShelf/%s/permissions/%s", rootShelfId, userPublicId),
		strings.NewReader(`{"permission":"Admin"}`),
	)
	request.Header.Set("Content-Type", "application/json")

	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d: %s", http.StatusNoContent, responseRecorder.Code, responseRecorder.Body.String())
	}
	if capturedReqDto == nil {
		t.Fatal("expected binder to invoke the controller")
	}
	if capturedReqDto.ContextFields.UserId != userId {
		t.Fatalf("expected user ID %s, got %s", userId, capturedReqDto.ContextFields.UserId)
	}
	if capturedReqDto.Param.RootShelfId != rootShelfId {
		t.Fatalf("expected root shelf ID %s, got %s", rootShelfId, capturedReqDto.Param.RootShelfId)
	}
	if capturedReqDto.Param.UserPublicId != userPublicId {
		t.Fatalf("expected user public ID %s, got %s", userPublicId, capturedReqDto.Param.UserPublicId)
	}
}
