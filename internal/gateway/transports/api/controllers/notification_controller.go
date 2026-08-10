package controllers

import (
	"github.com/gin-gonic/gin"

	notificationscontract "github.com/HiIamJeff67/notezy-backend/contracts/notification/v1/api"
	sharedcontexts "github.com/HiIamJeff67/notezy-backend/shared/lib/contexts"

	exceptionwriter "github.com/HiIamJeff67/notezy-backend/shared/util/exceptionwriter"

	gatewaycontexts "github.com/HiIamJeff67/notezy-backend/internal/gateway/contexts"
	notificationadapters "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/notification/adapters"
)

type NotificationControllerInterface interface {
	Search(ctx *gin.Context, requestDto *notificationscontract.SearchPrivateNotificationsRequestDto)
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

func (c *NotificationController) Search(
	ctx *gin.Context,
	requestDto *notificationscontract.SearchPrivateNotificationsRequestDto,
) {
	recipientUserPublicId, exception := gatewaycontexts.GetAndConvertContextFieldToUUID(
		ctx,
		sharedcontexts.ContextFieldName_User_PublicId,
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}
	requestDto.RecipientUserPublicId = *recipientUserPublicId

	response, exception := notificationadapters.CallSecurly[
		notificationscontract.SearchPrivateNotificationsRequestDto,
		notificationscontract.SearchPrivateNotificationsResponseDto,
	](ctx, c.notificationClient, requestDto, notificationscontract.SearchPrivateNotificationsOperation, "/internal/v1/notifications/search")
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
	recipientUserPublicId, exception := gatewaycontexts.GetAndConvertContextFieldToUUID(
		ctx,
		sharedcontexts.ContextFieldName_User_PublicId,
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}
	requestDto.RecipientUserPublicId = *recipientUserPublicId

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
	recipientUserPublicId, exception := gatewaycontexts.GetAndConvertContextFieldToUUID(
		ctx,
		sharedcontexts.ContextFieldName_User_PublicId,
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}
	requestDto.RecipientUserPublicId = *recipientUserPublicId

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
	recipientUserPublicId, exception := gatewaycontexts.GetAndConvertContextFieldToUUID(
		ctx,
		sharedcontexts.ContextFieldName_User_PublicId,
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}
	requestDto.RecipientUserPublicId = *recipientUserPublicId

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
