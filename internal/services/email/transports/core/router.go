package core

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func NewRouter(sender Sender, validator *validator.Validate) *gin.Engine {
	router := gin.New()
	router.GET("/healthz", func(ctx *gin.Context) {
		ctx.Status(http.StatusOK)
	})
	router.GET("/readyz", func(ctx *gin.Context) {
		ctx.Status(http.StatusOK)
	})
	endpoint := NewEndpoint(sender, validator)
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
