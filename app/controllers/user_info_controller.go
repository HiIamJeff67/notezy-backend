package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	dtos "notezy-backend/app/dtos"
	services "notezy-backend/app/services"
)

/* ============================== Interface & Instance ============================== */

type UserInfoControllerInterface interface {
	GetMyInfo(ctx *gin.Context, reqDto *dtos.GetMyInfoReqDto)
	UpdateMyInfo(ctx *gin.Context, reqDto *dtos.UpdateMyInfoReqDto)
}

type UserInfoController struct {
	userInfoService services.UserInfoServiceInterface
}

func NewUserInfoController(service services.UserInfoServiceInterface) UserInfoControllerInterface {
	return &UserInfoController{
		userInfoService: service,
	}
}

/* ============================== Controller ============================== */

// with AuthMiddleware
func (c *UserInfoController) GetMyInfo(ctx *gin.Context, reqDto *dtos.GetMyInfoReqDto) {
	resDto, exception := c.userInfoService.GetMyInfo(ctx.Request.Context(), reqDto)
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

// with AuthMiddleware
func (c *UserInfoController) UpdateMyInfo(ctx *gin.Context, reqDto *dtos.UpdateMyInfoReqDto) {
	resDto, exception := c.userInfoService.UpdateMyInfo(ctx.Request.Context(), reqDto)
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
