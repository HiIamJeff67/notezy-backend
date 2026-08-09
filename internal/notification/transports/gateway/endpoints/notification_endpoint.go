package endpoints

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	gatewaycontract "github.com/HiIamJeff67/notezy-backend/contracts/gateway/v1"
	notificationscontract "github.com/HiIamJeff67/notezy-backend/contracts/notification/v1/api"
	exceptions "github.com/HiIamJeff67/notezy-backend/contracts/types/exceptions"
	sharedcontexts "github.com/HiIamJeff67/notezy-backend/shared/lib/contexts"

	services "github.com/HiIamJeff67/notezy-backend/internal/notification/services"
)

type NotificationEndpoint struct {
	service *services.NotificationService
}

func NewNotificationEndpoint(service *services.NotificationService) *NotificationEndpoint {
	return &NotificationEndpoint{service: service}
}

func (e *NotificationEndpoint) List(ctx *gin.Context) {
	request := &gatewaycontract.Request[notificationscontract.ListNotificationsRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	request.Dto.RecipientUserPublicId = ctx.MustGet(sharedcontexts.ContextFieldName_User_PublicId.String()).(uuid.UUID)
	responseDto, err := e.service.List(ctx.Request.Context(), &request.Dto)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gatewaycontract.Response[struct{}]{
			Version: gatewaycontract.Version,
			Metadata: gatewaycontract.ResponseMetadata{
				RequestId:   request.Metadata.RequestId,
				RespondedAt: time.Now(),
			},
			Data: struct{}{},
			Exception: exceptions.New(
				"NotificationListFailed",
				"Notification",
				"List",
				"Failed to list notifications",
				http.StatusInternalServerError,
				true,
			).WithOrigin(err),
		})
		return
	}
	ctx.JSON(http.StatusOK, gatewaycontract.Response[notificationscontract.ListNotificationsResponseDto]{
		Version: gatewaycontract.Version,
		Metadata: gatewaycontract.ResponseMetadata{
			RequestId:   request.Metadata.RequestId,
			RespondedAt: time.Now(),
		},
		Data: *responseDto,
	})
}

func (e *NotificationEndpoint) CountUnread(ctx *gin.Context) {
	request := &gatewaycontract.Request[notificationscontract.CountUnreadNotificationsRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	request.Dto.RecipientUserPublicId = ctx.MustGet(sharedcontexts.ContextFieldName_User_PublicId.String()).(uuid.UUID)
	responseDto, err := e.service.CountUnread(ctx.Request.Context(), &request.Dto)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gatewaycontract.Response[struct{}]{
			Version: gatewaycontract.Version,
			Metadata: gatewaycontract.ResponseMetadata{
				RequestId:   request.Metadata.RequestId,
				RespondedAt: time.Now(),
			},
			Data: struct{}{},
			Exception: exceptions.New(
				"NotificationUnreadCountFailed",
				"Notification",
				"CountUnread",
				"Failed to count unread notifications",
				http.StatusInternalServerError,
				true,
			).WithOrigin(err),
		})
		return
	}
	ctx.JSON(http.StatusOK, gatewaycontract.Response[notificationscontract.CountUnreadNotificationsResponseDto]{
		Version: gatewaycontract.Version,
		Metadata: gatewaycontract.ResponseMetadata{
			RequestId:   request.Metadata.RequestId,
			RespondedAt: time.Now(),
		},
		Data: *responseDto,
	})
}

func (e *NotificationEndpoint) MarkRead(ctx *gin.Context) {
	request := &gatewaycontract.Request[notificationscontract.MarkNotificationsReadRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	request.Dto.RecipientUserPublicId = ctx.MustGet(sharedcontexts.ContextFieldName_User_PublicId.String()).(uuid.UUID)
	responseDto, err := e.service.MarkRead(ctx.Request.Context(), &request.Dto)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gatewaycontract.Response[struct{}]{
			Version: gatewaycontract.Version,
			Metadata: gatewaycontract.ResponseMetadata{
				RequestId:   request.Metadata.RequestId,
				RespondedAt: time.Now(),
			},
			Data: struct{}{},
			Exception: exceptions.New(
				"NotificationReadFailed",
				"Notification",
				"MarkRead",
				"Failed to mark notifications as read",
				http.StatusInternalServerError,
				true,
			).WithOrigin(err),
		})
		return
	}
	ctx.JSON(http.StatusOK, gatewaycontract.Response[notificationscontract.MarkNotificationsReadResponseDto]{
		Version: gatewaycontract.Version,
		Metadata: gatewaycontract.ResponseMetadata{
			RequestId:   request.Metadata.RequestId,
			RespondedAt: time.Now(),
		},
		Data: *responseDto,
	})
}

func (e *NotificationEndpoint) Delete(ctx *gin.Context) {
	request := &gatewaycontract.Request[notificationscontract.DeleteNotificationsRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	request.Dto.RecipientUserPublicId = ctx.MustGet(sharedcontexts.ContextFieldName_User_PublicId.String()).(uuid.UUID)
	responseDto, err := e.service.SoftDelete(ctx.Request.Context(), &request.Dto)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gatewaycontract.Response[struct{}]{
			Version: gatewaycontract.Version,
			Metadata: gatewaycontract.ResponseMetadata{
				RequestId:   request.Metadata.RequestId,
				RespondedAt: time.Now(),
			},
			Data: struct{}{},
			Exception: exceptions.New(
				"NotificationDeleteFailed",
				"Notification",
				"Delete",
				"Failed to delete notifications",
				http.StatusInternalServerError,
				true,
			).WithOrigin(err),
		})
		return
	}
	ctx.JSON(http.StatusOK, gatewaycontract.Response[notificationscontract.DeleteNotificationsResponseDto]{
		Version: gatewaycontract.Version,
		Metadata: gatewaycontract.ResponseMetadata{
			RequestId:   request.Metadata.RequestId,
			RespondedAt: time.Now(),
		},
		Data: *responseDto,
	})
}
