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

type UserAccountControllerInterface interface {
	GetMyAccount(ctx *gin.Context)
	UpdateMyAccount(ctx *gin.Context)
}

type userAccountController struct {
	userAccountService services.UserAccountServiceInterface
}

var UserAccountController UserAccountControllerInterface = &userAccountController{
	userAccountService: services.UserAccountService,
}

/* ============================== Controllers ============================== */

// with AuthMiddleware
func (c *userAccountController) GetMyAccount(ctx *gin.Context) {
	var reqDto dtos.GetMyAccountReqDto
	userId, exception := contexts.FetchAndConvertContextFieldToUUID(ctx, "userId")
	if exception != nil {
		ctx.JSON(exception.HTTPStatusCode, exception.GetGinH())
		return
	}
	reqDto.UserId = *userId

	resDto, exception := c.userAccountService.GetMyAccount(&reqDto)
	if exception != nil {
		ctx.JSON(exception.HTTPStatusCode, exception.GetGinH())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"updatedAt": resDto,
		},
	})
}

// with AuthMiddleware
func (c *userAccountController) UpdateMyAccount(ctx *gin.Context) {
	var reqDto dtos.UpdateMyAccountReqDto
	userId, exception := contexts.FetchAndConvertContextFieldToUUID(ctx, "userId")
	if exception != nil {
		ctx.JSON(exception.HTTPStatusCode, exception.GetGinH())
		return
	}
	reqDto.UserId = *userId
	if err := ctx.ShouldBindJSON(&reqDto); err != nil {
		exception := exceptions.UserAccount.InvalidDto().WithError(err)
		ctx.JSON(exception.HTTPStatusCode, exception.GetGinH())
		return
	}

	resDto, exception := c.userAccountService.UpdateMyAccount(&reqDto)
	if exception != nil {
		ctx.JSON(exception.HTTPStatusCode, exception.GetGinH())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"updatedAt": resDto,
		},
	})
}
