package core

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	router := gin.New()
	router.GET("/healthz", func(ctx *gin.Context) {
		ctx.Status(http.StatusOK)
	})
	router.GET("/readyz", func(ctx *gin.Context) {
		ctx.Status(http.StatusOK)
	})
	return router
}
