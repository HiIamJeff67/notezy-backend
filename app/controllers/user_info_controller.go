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

type UserInfoControllerInterface interface {
	GetMyInfo(ctx *gin.Context)
	UpdateMyInfo(ctx *gin.Context)
}

type userInfoController struct {
	userInfoService services.UserInfoServiceInterface
}

var UserInfoController UserInfoControllerInterface = &userInfoController{
	userInfoService: services.UserInfoService,
}

/* ============================== Controllers ============================== */

// with AuthMiddleware()
func (c *userInfoController) GetMyInfo(ctx *gin.Context) {
	var reqDto dtos.GetMyInfoReqDto
	userId, exception := contexts.FetchAndConvertContextFieldToUUID(ctx, "userId")
	if exception != nil {
		ctx.JSON(exception.HTTPStatusCode, exception.GetGinH())
		return
	}
	reqDto.UserId = *userId

	resDto, exception := c.userInfoService.GetMyInfo(&reqDto)
	if exception != nil {
		ctx.JSON(exception.HTTPStatusCode, exception.GetGinH())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": resDto,
	})
}

// with AuthMiddleware()
func (c *userInfoController) UpdateMyInfo(ctx *gin.Context) {
	var reqDto dtos.UpdateMyInfoReqDto
	userId, exception := contexts.FetchAndConvertContextFieldToUUID(ctx, "userId")
	if exception != nil {
		ctx.JSON(exception.HTTPStatusCode, exception.GetGinH())
		return
	}
	reqDto.UserId = *userId
	if err := ctx.ShouldBindJSON(&reqDto); err != nil {
		exception := exceptions.UserInfo.InvalidDto().WithError(err)
		ctx.JSON(exception.HTTPStatusCode, exception.GetGinH())
		return
	}

	resDto, exception := c.userInfoService.UpdateMyInfo(&reqDto)
	if exception != nil {
		ctx.JSON(exception.HTTPStatusCode, exception.GetGinH())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": resDto,
	})
}
