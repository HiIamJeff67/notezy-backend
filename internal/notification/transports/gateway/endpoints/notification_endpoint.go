package endpoints

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	gatewaycontract "github.com/HiIamJeff67/notegic-backend/contracts/gateway/v1"
	notificationscontract "github.com/HiIamJeff67/notegic-backend/contracts/notification/v1/api"
	exceptions "github.com/HiIamJeff67/notegic-backend/contracts/types/exceptions"
	sharedcontexts "github.com/HiIamJeff67/notegic-backend/shared/lib/contexts"

	services "github.com/HiIamJeff67/notegic-backend/internal/notification/services"
)

type NotificationEndpoint struct {
	service services.NotificationServiceInterface
}

func NewNotificationEndpoint(service services.NotificationServiceInterface) *NotificationEndpoint {
	return &NotificationEndpoint{service: service}
}

func (e *NotificationEndpoint) Search(ctx *gin.Context) {
	request := &gatewaycontract.Request[notificationscontract.SearchPrivateNotificationsRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	request.Dto.RecipientUserPublicId = ctx.MustGet(sharedcontexts.ContextFieldName_User_PublicId.String()).(uuid.UUID)
	responseDto, err := e.service.SearchPrivateNotifications(ctx.Request.Context(), &request.Dto)
	if err != nil {
		responseException := exceptions.New(
			"NotificationSearchFailed",
			"Notification",
			"SearchPrivateNotifications",
			"Failed to search private notifications",
			http.StatusInternalServerError,
			true,
		).WithOrigin(err)
		var serviceException *exceptions.Exception
		if errors.As(err, &serviceException) {
			responseException = serviceException
		}
		publicException := responseException.ToPublic()

		ctx.JSON(publicException.HTTPStatusCode(), gatewaycontract.Response[struct{}]{
			Version: gatewaycontract.Version,
			Metadata: gatewaycontract.ResponseMetadata{
				RequestId:   request.Metadata.RequestId,
				RespondedAt: time.Now(),
			},
			Data:      struct{}{},
			Exception: publicException,
		})
		return
	}
	ctx.JSON(http.StatusOK, gatewaycontract.Response[notificationscontract.SearchPrivateNotificationsResponseDto]{
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
	responseDto, err := e.service.CountMyUnreadNotifications(ctx.Request.Context(), &request.Dto)
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
	responseDto, err := e.service.MarkMyNotificationsRead(ctx.Request.Context(), &request.Dto)
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
	responseDto, err := e.service.SoftDeleteMyNotifications(ctx.Request.Context(), &request.Dto)
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
