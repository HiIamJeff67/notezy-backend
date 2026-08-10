package binders

import (
	"github.com/gin-gonic/gin"

	notificationscontract "github.com/HiIamJeff67/notezy-backend/contracts/notification/v1/api"
	exceptions "github.com/HiIamJeff67/notezy-backend/contracts/types/exceptions"

	exceptionwriter "github.com/HiIamJeff67/notezy-backend/shared/util/exceptionwriter"

	controllers "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/api/controllers"
)

type NotificationBinderInterface interface {
	BindSearch(controllers.Func[*notificationscontract.SearchPrivateNotificationsRequestDto]) gin.HandlerFunc
	BindCountUnread(controllers.Func[*notificationscontract.CountUnreadNotificationsRequestDto]) gin.HandlerFunc
	BindMarkRead(controllers.Func[*notificationscontract.MarkNotificationsReadRequestDto]) gin.HandlerFunc
	BindDelete(controllers.Func[*notificationscontract.DeleteNotificationsRequestDto]) gin.HandlerFunc
}

type NotificationBinder struct{}

func NewNotificationBinder() NotificationBinderInterface {
	return &NotificationBinder{}
}

func (b *NotificationBinder) BindSearch(
	controllerFunc controllers.Func[*notificationscontract.SearchPrivateNotificationsRequestDto],
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &notificationscontract.SearchPrivateNotificationsRequestDto{}
		if err := ctx.ShouldBindQuery(requestDto); err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("Notification").WithOrigin(err), ctx)
			return
		}
		controllerFunc(ctx, requestDto)
	}
}

func (b *NotificationBinder) BindCountUnread(
	controllerFunc controllers.Func[*notificationscontract.CountUnreadNotificationsRequestDto],
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		controllerFunc(ctx, &notificationscontract.CountUnreadNotificationsRequestDto{})
	}
}

func (b *NotificationBinder) BindMarkRead(
	controllerFunc controllers.Func[*notificationscontract.MarkNotificationsReadRequestDto],
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &notificationscontract.MarkNotificationsReadRequestDto{}
		if err := ctx.ShouldBindJSON(requestDto); err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("Notification").WithOrigin(err), ctx)
			return
		}
		controllerFunc(ctx, requestDto)
	}
}

func (b *NotificationBinder) BindDelete(
	controllerFunc controllers.Func[*notificationscontract.DeleteNotificationsRequestDto],
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &notificationscontract.DeleteNotificationsRequestDto{}
		if err := ctx.ShouldBindJSON(requestDto); err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("Notification").WithOrigin(err), ctx)
			return
		}
		controllerFunc(ctx, requestDto)
	}
}
