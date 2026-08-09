package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	notificationscontract "github.com/HiIamJeff67/notezy-backend/contracts/notification/v1/api"
	sharedcontexts "github.com/HiIamJeff67/notezy-backend/shared/lib/contexts"

	exceptionwriter "github.com/HiIamJeff67/notezy-backend/shared/util/exceptionwriter"

	notificationadapters "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/notification/adapters"
)

type NotificationControllerInterface interface {
	List(ctx *gin.Context, requestDto *notificationscontract.ListNotificationsRequestDto)
	CountUnread(ctx *gin.Context, requestDto *notificationscontract.CountUnreadNotificationsRequestDto)
	MarkRead(ctx *gin.Context, requestDto *notificationscontract.MarkNotificationsReadRequestDto)
	Delete(ctx *gin.Context, requestDto *notificationscontract.DeleteNotificationsRequestDto)
}

type NotificationController struct {
	notificationClient *notificationadapters.NotificationAdapter
}

func NewNotificationController(
	notificationClient *notificationadapters.NotificationAdapter,
) NotificationControllerInterface {
	return &NotificationController{notificationClient: notificationClient}
}

func (c *NotificationController) List(
	ctx *gin.Context,
	requestDto *notificationscontract.ListNotificationsRequestDto,
) {
	requestDto.RecipientUserPublicId = ctx.MustGet(sharedcontexts.ContextFieldName_User_PublicId.String()).(uuid.UUID)
	response, exception := notificationadapters.CallSecurly[
		notificationscontract.ListNotificationsRequestDto,
		notificationscontract.ListNotificationsResponseDto,
	](ctx, c.notificationClient, requestDto, notificationscontract.ListMyNotificationsOperation, "/internal/v1/notifications/list")
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}
	writeClientResponse(ctx, response.Data)
}

func (c *NotificationController) CountUnread(
	ctx *gin.Context,
	requestDto *notificationscontract.CountUnreadNotificationsRequestDto,
) {
	requestDto.RecipientUserPublicId = ctx.MustGet(sharedcontexts.ContextFieldName_User_PublicId.String()).(uuid.UUID)
	response, exception := notificationadapters.CallSecurly[
		notificationscontract.CountUnreadNotificationsRequestDto,
		notificationscontract.CountUnreadNotificationsResponseDto,
	](ctx, c.notificationClient, requestDto, notificationscontract.CountMyUnreadNotificationsOperation, "/internal/v1/notifications/unread-count")
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}
	writeClientResponse(ctx, response.Data)
}

func (c *NotificationController) MarkRead(
	ctx *gin.Context,
	requestDto *notificationscontract.MarkNotificationsReadRequestDto,
) {
	requestDto.RecipientUserPublicId = ctx.MustGet(sharedcontexts.ContextFieldName_User_PublicId.String()).(uuid.UUID)
	response, exception := notificationadapters.CallSecurly[
		notificationscontract.MarkNotificationsReadRequestDto,
		notificationscontract.MarkNotificationsReadResponseDto,
	](ctx, c.notificationClient, requestDto, notificationscontract.MarkMyNotificationsReadOperation, "/internal/v1/notifications/read")
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}
	writeClientResponse(ctx, response.Data)
}

func (c *NotificationController) Delete(
	ctx *gin.Context,
	requestDto *notificationscontract.DeleteNotificationsRequestDto,
) {
	requestDto.RecipientUserPublicId = ctx.MustGet(sharedcontexts.ContextFieldName_User_PublicId.String()).(uuid.UUID)
	response, exception := notificationadapters.CallSecurly[
		notificationscontract.DeleteNotificationsRequestDto,
		notificationscontract.DeleteNotificationsResponseDto,
	](ctx, c.notificationClient, requestDto, notificationscontract.DeleteMyNotificationsOperation, "/internal/v1/notifications/delete")
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}
	writeClientResponse(ctx, response.Data)
}
