package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	contexts "notezy-backend/app/contexts"
	dtos "notezy-backend/app/dtos"
	exceptions "notezy-backend/app/exceptions"
	services "notezy-backend/app/services"
	constants "notezy-backend/shared/constants"
)

/* ============================== Interface & Instance ============================== */

type UserInfoControllerInterface interface {
	GetMyInfo(ctx *gin.Context)
	UpdateMyInfo(ctx *gin.Context)
}

type UserInfoController struct {
	userInfoService services.UserInfoServiceInterface
}

func NewUserInfoController(service services.UserInfoServiceInterface) UserInfoControllerInterface {
	return &UserInfoController{
		userInfoService: service,
	}
}

/* ============================== Controllers ============================== */

// with AuthMiddleware()
func (c *UserInfoController) GetMyInfo(ctx *gin.Context) {
	var reqDto dtos.GetMyInfoReqDto
	userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
	if exception != nil {
		exception.Log().SafelyResponseWithJSON(ctx)
		return
	}
	reqDto.ContextFields.UserId = *userId

	resDto, exception := c.userInfoService.GetMyInfo(&reqDto)
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

// with AuthMiddleware()
func (c *UserInfoController) UpdateMyInfo(ctx *gin.Context) {
	var reqDto dtos.UpdateMyInfoReqDto
	userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
	if exception != nil {
		exception.Log().SafelyResponseWithJSON(ctx)
		return
	}
	reqDto.ContextFields.UserId = *userId
	if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
		exception := exceptions.UserInfo.InvalidDto().WithError(err)
		exception.ResponseWithJSON(ctx)
		return
	}

	resDto, exception := c.userInfoService.UpdateMyInfo(&reqDto)
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
