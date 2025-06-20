package controllers

import (
	"net/http"
	"notezy-backend/app/contexts"
	"notezy-backend/app/dtos"
	services "notezy-backend/app/services"

	"github.com/gin-gonic/gin"
)

// with AuthMiddleware()
func GetMyInfo(ctx *gin.Context) {
	var reqDto dtos.GetMyInfoReqDto
	userId, exception := contexts.FetchAndConvertContextFieldToUUID(ctx, "userId")
	if exception != nil {
		ctx.JSON(exception.HTTPStatusCode, exception.GetGinH())
		return
	}
	reqDto.UserId = *userId

	resDto, exception := services.GetMyInfo(&reqDto)
	if exception != nil {
		ctx.JSON(
			exception.HTTPStatusCode,
			exception.GetGinH(),
		)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": resDto,
	})
}
