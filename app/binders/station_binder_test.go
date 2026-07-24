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
	enums "github.com/HiIamJeff67/notezy-backend/app/models/schemas/enums"
	types "github.com/HiIamJeff67/notezy-backend/shared/types"
)

func TestBindCreateMyStationPermissionParsesPermissionRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)
	userId, stationId, userPublicId := uuid.New(), uuid.New(), uuid.New()

	var capturedReqDto *dtos.CreateMyStationPermissionReqDto
	router := gin.New()
	router.Use(func(ctx *gin.Context) {
		ctx.Set(types.ContextFieldName_User_Id.String(), userId.String())
		ctx.Next()
	})
	router.POST("/station/:stationId/permissions/:userPublicId", NewStationBinder().BindCreateMyStationPermission(func(ctx *gin.Context, reqDto *dtos.CreateMyStationPermissionReqDto) {
		capturedReqDto = reqDto
		ctx.Status(http.StatusNoContent)
	}))

	request := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/station/%s/permissions/%s", stationId, userPublicId), strings.NewReader(`{"permission":"Write"}`))
	request.Header.Set("Content-Type", "application/json")
	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d: %s", http.StatusNoContent, responseRecorder.Code, responseRecorder.Body.String())
	}
	if capturedReqDto == nil || capturedReqDto.ContextFields.UserId != userId || capturedReqDto.Param.StationId != stationId || capturedReqDto.Param.UserPublicId != userPublicId || capturedReqDto.Body.Permission != enums.AccessControlPermission_Write {
		t.Fatal("expected station permission binder to populate the request")
	}
}

func TestBindLeaveMyStationParsesOwnerTransferTarget(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userId, stationId, targetUserPublicId := uuid.New(), uuid.New(), uuid.New()
	var capturedReqDto *dtos.LeaveMyStationReqDto

	router := gin.New()
	router.Use(func(ctx *gin.Context) {
		ctx.Set(types.ContextFieldName_User_Id.String(), userId.String())
		ctx.Next()
	})
	router.DELETE("/station/:stationId/leave", NewStationBinder().BindLeaveMyStation(func(ctx *gin.Context, reqDto *dtos.LeaveMyStationReqDto) {
		capturedReqDto = reqDto
		ctx.Status(http.StatusNoContent)
	}))

	request := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/station/%s/leave", stationId), strings.NewReader(fmt.Sprintf(`{"targetUserPublicId":"%s"}`, targetUserPublicId)))
	request.Header.Set("Content-Type", "application/json")
	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != http.StatusNoContent || capturedReqDto == nil || capturedReqDto.ContextFields.UserId != userId || capturedReqDto.Param.StationId != stationId || capturedReqDto.Body.TargetUserPublicId == nil || *capturedReqDto.Body.TargetUserPublicId != targetUserPublicId {
		t.Fatal("expected leave binder to populate the authenticated user, station, and owner transfer target")
	}
}
