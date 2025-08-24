package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	contexts "notezy-backend/app/contexts"
	dtos "notezy-backend/app/dtos"
	services "notezy-backend/app/services"
	constants "notezy-backend/shared/constants"
)

/* ============================== Interface & Instance ============================== */

type UserSettingControllerInterface interface {
	GetMySetting(ctx *gin.Context)
}

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
	userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
	if exception != nil {
		exception.Log().SafelyResponseWithJSON(ctx)
		return
	}
	reqDto.ContextFields.UserId = *userId

	resDto, exception := c.userSettingService.GetMySetting(&reqDto)
	if exception != nil {
		exception.Log().SafelyResponseWithJSON(ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      resDto,
		"exception": nil,
	})
}
