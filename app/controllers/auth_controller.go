package controllers

import (
	"net/http"
	"notezy-backend/app/dtos"
	"notezy-backend/app/exceptions"
	"notezy-backend/app/services"

	"github.com/gin-gonic/gin"
)

/* ============================== Controller ============================== */
func Register(ctx *gin.Context) {
	var reqDto dtos.RegisterReqDto
	if err := ctx.ShouldBindJSON(&reqDto); err != nil {
		ctx.JSON(http.StatusBadRequest, exceptions.Auth.InvalidDto().WithError(err).GetGinH())
		return
	}

	resDto, exception := services.Register(&reqDto)
	if exception != nil {
		ctx.JSON(exception.HTTPStatusCode, exception.GetGinH())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"AccessToken": resDto.AccessToken,
			"createdAt":   resDto.CreatedAt,
		},
	})
}

func Login(ctx *gin.Context) {
	var reqDto dtos.LoginReqDto
	if err := ctx.ShouldBindJSON(&reqDto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resDto, exception := services.Login(&reqDto)
	if exception != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": exception.Log().Error})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"AccessToken": resDto.AccessToken,
			"createdAt":   resDto.CreatedAt,
		},
	})
}
