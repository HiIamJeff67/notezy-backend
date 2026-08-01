package core

import "github.com/gin-gonic/gin"

func NewRouter(sender Sender) *gin.Engine {
	router := gin.New()
	endpoint := NewEndpoint(sender)
	emailRoutes := router.Group("/email/v1")
	{
		emailRoutes.POST(
			"/send/welcome",
			endpoint.SendWelcomeEmail,
		)
		emailRoutes.POST(
			"/send/validation",
			endpoint.SendValidationEmail,
		)
		emailRoutes.POST(
			"/send/security-alert",
			endpoint.SendSecurityAlertEmail,
		)
	}

	return router
}
