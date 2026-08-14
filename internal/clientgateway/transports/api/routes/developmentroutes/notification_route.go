package developmentroutes

import (
	"time"

	"github.com/gin-gonic/gin"

	cookies "github.com/HiIamJeff67/notezy-backend/shared/cookies"

	binder "github.com/HiIamJeff67/notezy-backend/internal/clientgateway/transports/api/binders"
	controllers "github.com/HiIamJeff67/notezy-backend/internal/clientgateway/transports/api/controllers"
	interceptors "github.com/HiIamJeff67/notezy-backend/internal/clientgateway/transports/api/interceptors"
	middlewares "github.com/HiIamJeff67/notezy-backend/internal/clientgateway/transports/api/middlewares"
	notificationadapters "github.com/HiIamJeff67/notezy-backend/internal/clientgateway/transports/notification/adapters"
)

func configureDevelopmentNotificationRoutes(
	router *gin.RouterGroup,
	notificationClient *notificationadapters.NotificationAdapter,
	accessTokenCookieHandler *cookies.CookieHandler,
	refreshTokenCookieHandler *cookies.CookieHandler,
	rateLimiters RateLimiters,
) {
	if router == nil {
		router = DevelopmentAPIRouterGroup
	}

	notificationBinder := binder.NewNotificationBinder()
	notificationController := controllers.NewNotificationController(notificationClient)
	notificationRoutes := router.Group("/notifications")
	defaultMiddlewares := []gin.HandlerFunc{
		middlewares.UnauthorizedRateLimitMiddleware(rateLimiters.Unauthorized),
		middlewares.TimeoutMiddleware(1 * time.Second),
		middlewares.GatewayAuthenticationMiddleware(accessTokenCookieHandler, refreshTokenCookieHandler),
		interceptors.ShareableResponseWriterInterceptor(
			interceptors.RefreshTokenInterceptor(accessTokenCookieHandler),
			interceptors.EmbeddedInterceptor,
		),
	}
	{
		notificationRoutes.GET(
			"/",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("searchPrivateNotifications"),
					middlewares.ApplyMeterMiddleware("server.requests.notifications.search"),
				},
				defaultMiddlewares,
				notificationBinder.BindSearch(notificationController.Search),
			)...,
		)
		notificationRoutes.GET(
			"/unread-count",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("countMyUnreadNotifications"),
					middlewares.ApplyMeterMiddleware("server.requests.notifications.unreadCount"),
				},
				defaultMiddlewares,
				notificationBinder.BindCountUnread(notificationController.CountUnread),
			)...,
		)
		notificationRoutes.PATCH(
			"/read",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("markMyNotificationsRead"),
					middlewares.ApplyMeterMiddleware("server.requests.notifications.read"),
				},
				defaultMiddlewares,
				notificationBinder.BindMarkRead(notificationController.MarkRead),
			)...,
		)
		notificationRoutes.DELETE(
			"/",
			middlewares.RepositionMiddleware(
				[]gin.HandlerFunc{
					middlewares.ApplyTracerMiddleware("deleteMyNotifications"),
					middlewares.ApplyMeterMiddleware("server.requests.notifications.delete"),
				},
				defaultMiddlewares,
				notificationBinder.BindDelete(notificationController.Delete),
			)...,
		)
	}
}
