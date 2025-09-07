package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	dtos "notezy-backend/app/dtos"
	services "notezy-backend/app/services"
)

/* ============================== Interface & Instance ============================== */

type UserSettingControllerInterface interface {
	GetMySetting(ctx *gin.Context, reqDto *dtos.GetMySettingReqDto)
}

type UserSettingController struct {
	userSettingService services.UserSettingServiceInterface
}

func NewUserSettingController(service services.UserSettingServiceInterface) UserSettingControllerInterface {
	return &UserSettingController{
		userSettingService: service,
	}
}

/* ============================== Controller ============================== */

func (c *UserSettingController) GetMySetting(ctx *gin.Context, reqDto *dtos.GetMySettingReqDto) {
	resDto, exception := c.userSettingService.GetMySetting(reqDto)
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
