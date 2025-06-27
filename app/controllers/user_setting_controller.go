package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	contexts "notezy-backend/app/contexts"
	dtos "notezy-backend/app/dtos"
	"notezy-backend/app/exceptions"
	services "notezy-backend/app/services"
)

/* ============================== Interface & Instance ============================== */

type UserSettingControllerInterface interface{}

type userSettingController struct {
	userSettingService services.UserSettingServiceInterface
}

var UserSettingController UserSettingControllerInterface = &userSettingController{
	userSettingService: services.UserSettingService,
}

/* ============================== Controllers ============================== */

func (c *userSettingController) GetMySetting(ctx *gin.Context) {
	var reqDto dtos.GetMySettingReqDto
	userId, exception := contexts.FetchAndConvertContextFieldToUUID(ctx, "userId")
	if exception != nil {
		exception.Log()
		exception = exceptions.UserSetting.InternalServerWentWrong(exception)
		ctx.JSON(exception.HTTPStatusCode, exception.GetGinH())
		return
	}
	reqDto.UserId = *userId

	resDto, exception := c.userSettingService.GetMySetting(&reqDto)
	if exception != nil {
		exception.Log()
		if !exceptions.CompareCommonExceptions(exceptions.UserSetting.InvalidDto(), exception, false) {
			exception = exceptions.UserSetting.InternalServerWentWrong(exception)
		}
		ctx.JSON(exception.HTTPStatusCode, exception.GetGinH())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    resDto,
	})
}
