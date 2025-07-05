package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	contexts "notezy-backend/app/contexts"
	dtos "notezy-backend/app/dtos"
	exceptions "notezy-backend/app/exceptions"
	services "notezy-backend/app/services"
)

/* ============================== Interface & Instance ============================== */

type UserSettingControllerInterface interface{}

type UserSettingController struct {
	userSettingService services.UserSettingServiceInterface
}

func NewUserSettingController(service services.UserSettingServiceInterface) UserSettingControllerInterface {
	return &UserSettingController{
		userSettingService: service,
	}
}

/* ============================== Controllers ============================== */

func (c *UserSettingController) GetMySetting(ctx *gin.Context) {
	var reqDto dtos.GetMySettingReqDto
	userId, exception := contexts.FetchAndConvertContextFieldToUUID(ctx, "userId")
	if exception != nil {
		exception.Log()
		exception = exceptions.UserSetting.InternalServerWentWrong(exception)
		ctx.JSON(exception.HTTPStatusCode, gin.H{
			"success":   false,
			"data":      nil,
			"exception": exception.GetGinH(),
		})
		return
	}
	reqDto.UserId = *userId

	resDto, exception := c.userSettingService.GetMySetting(&reqDto)
	if exception != nil {
		exception.Log()
		if !exceptions.CompareCommonExceptions(exceptions.UserSetting.InvalidDto(), exception, false) {
			exception = exceptions.UserSetting.InternalServerWentWrong(exception)
		}
		ctx.JSON(exception.HTTPStatusCode, gin.H{
			"success":   false,
			"data":      nil,
			"exception": exception.GetGinH(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      resDto,
		"exception": nil,
	})
}
