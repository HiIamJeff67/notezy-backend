package routers

import (
	"github.com/gin-gonic/gin"

	notificationscontract "github.com/HiIamJeff67/notezy-backend/contracts/notification/v1/api"

	endpoints "github.com/HiIamJeff67/notezy-backend/internal/notification/transports/gateway/endpoints"
	middlewares "github.com/HiIamJeff67/notezy-backend/internal/notification/transports/gateway/middlewares"
)

func ConfigureNotificationRoutes(
	router *gin.RouterGroup,
	endpoint *endpoints.NotificationEndpoint,
) {
	notificationRoutes := router.Group("/notifications")
	notificationRoutes.POST(
		"/search",
		middlewares.DelegationAuthenticatedMiddleware(notificationscontract.SearchPrivateNotificationsOperation),
		endpoint.Search,
	)
	notificationRoutes.POST(
		"/unread-count",
		middlewares.DelegationAuthenticatedMiddleware(notificationscontract.CountMyUnreadNotificationsOperation),
		endpoint.CountUnread,
	)
	notificationRoutes.POST(
		"/read",
		middlewares.DelegationAuthenticatedMiddleware(notificationscontract.MarkMyNotificationsReadOperation),
		endpoint.MarkRead,
	)
	notificationRoutes.POST(
		"/delete",
		middlewares.DelegationAuthenticatedMiddleware(notificationscontract.DeleteMyNotificationsOperation),
		endpoint.Delete,
	)
}
