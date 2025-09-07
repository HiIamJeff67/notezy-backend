package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	dtos "notezy-backend/app/dtos"
	services "notezy-backend/app/services"
)

/* ============================== Interface & Instance ============================== */

type UserAccountControllerInterface interface {
	GetMyAccount(ctx *gin.Context, reqDto *dtos.GetMyAccountReqDto)
	UpdateMyAccount(ctx *gin.Context, reqDto *dtos.UpdateMyAccountReqDto)
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
func (c *UserAccountController) GetMyAccount(ctx *gin.Context, reqDto *dtos.GetMyAccountReqDto) {
	resDto, exception := c.userAccountService.GetMyAccount(reqDto)
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
func (c *UserAccountController) UpdateMyAccount(ctx *gin.Context, reqDto *dtos.UpdateMyAccountReqDto) {
	resDto, exception := c.userAccountService.UpdateMyAccount(reqDto)
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
