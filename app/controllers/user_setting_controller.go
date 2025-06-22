package controllers

import (
	"net/http"
	"notezy-backend/app/contexts"
	"notezy-backend/app/dtos"
	"notezy-backend/app/services"

	"github.com/gin-gonic/gin"
)

/* ============================== Interface & Instance ============================== */

type UserSettingControllerInterface interface{}

type userSettingController struct {
	userSettingService services.UserSettingServiceInterface
}

var UserSettingController UserSettingControllerInterface = &userSettingController{}

/* ============================== Controllers ============================== */

func (c *userSettingController) GetMySetting(ctx *gin.Context) {
	var reqDto dtos.GetMySettingReqDto
	userId, exception := contexts.FetchAndConvertContextFieldToUUID(ctx, "userId")
	if exception != nil {
		ctx.JSON(exception.HTTPStatusCode, exception.GetGinH())
		return
	}
	reqDto.UserId = *userId

	resDto, exception := c.userSettingService.GetMySetting(&reqDto)
	if exception != nil {
		ctx.JSON(exception.HTTPStatusCode, exception.GetGinH())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": resDto,
	})
}
