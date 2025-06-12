package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func configureTestRoutes() {
	testRoutes := RouterGroup.Group("/test")
	{
		testRoutes.GET("hi", func(ctx *gin.Context) { ctx.JSON(http.StatusOK, gin.H{"message": "Hello"}) })
	}
}
