package controllers

import (
	"net/http"
	"notezy-backend/app/services"

	"github.com/gin-gonic/gin"
)

func FindAllUsers(ctx *gin.Context) {
	resDto, exception := services.FindAllUsers()
	if exception != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": exception.Log().Error})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data": gin.H{
			"users": resDto,
		},
	})
}
