package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"notezy-backend/app/dtos"
	exceptions "notezy-backend/app/exceptions"
	services "notezy-backend/app/services"
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

func UpdateMe(ctx *gin.Context) {
	var reqDto dtos.UpdateMeReqDto
	accessTokenFromCookie, exists := ctx.Get("accessToken")
	if !exists {
		ctx.JSON(
			exceptions.Auth.InvalidDto().HTTPStatusCode,
			exceptions.Auth.InvalidDto().GetGinH(),
		)
		return
	}

	tokenStr, ok := accessTokenFromCookie.(string)
	if !ok {
		ctx.JSON(
			exceptions.Auth.InvalidDto().HTTPStatusCode,
			exceptions.Auth.InvalidDto().GetGinH(),
		)
		return
	}

	reqDto.AccessToken = tokenStr
	if err := ctx.ShouldBindJSON(&reqDto); err != nil {
		ctx.JSON(
			exceptions.Auth.InvalidDto().HTTPStatusCode,
			exceptions.Auth.InvalidDto().WithError(err).GetGinH(),
		)
		return
	}

	resDto, exception := services.UpdateMe(&reqDto)
	if exception != nil {
		ctx.JSON(exception.HTTPStatusCode, exception.GetGinH())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"updatedAt": resDto.UpdatedAt,
		},
	})
}
