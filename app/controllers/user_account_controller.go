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

type UserAccountControllerInterface interface {
	GetMyAccount(ctx *gin.Context)
	UpdateMyAccount(ctx *gin.Context)
}

type UserAccountController struct {
	userAccountService services.UserAccountServiceInterface
}

func NewUserAccountController(service services.UserAccountServiceInterface) UserAccountControllerInterface {
	return &UserAccountController{
		userAccountService: service,
	}
}

/* ============================== Controllers ============================== */

// with AuthMiddleware
func (c *UserAccountController) GetMyAccount(ctx *gin.Context) {
	var reqDto dtos.GetMyAccountReqDto
	userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
	if exception != nil {
		exception.Log().SafelyResponseWithJSON(ctx)
		return
	}
	reqDto.ContextFields.UserId = *userId

	resDto, exception := c.userAccountService.GetMyAccount(&reqDto)
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
func (c *UserAccountController) UpdateMyAccount(ctx *gin.Context) {
	var reqDto dtos.UpdateMyAccountReqDto
	userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
	if exception != nil {
		exception.Log().SafelyResponseWithJSON(ctx)
		return
	}
	reqDto.ContextFields.UserId = *userId
	if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
		exception := exceptions.UserAccount.InvalidDto().WithError(err)
		exception.ResponseWithJSON(ctx)
		return
	}

	resDto, exception := c.userAccountService.UpdateMyAccount(&reqDto)
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
