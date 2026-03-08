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
	BindGoogleAccount(ctx *gin.Context, reqDto *dtos.BindGoogleAccountReqDto)
	UnbindGoogleAccount(ctx *gin.Context, reqDto *dtos.UnbindGoogleAccountReqDto)
}

type UserAccountController struct {
	userAccountService services.UserAccountServiceInterface
}

func NewUserAccountController(service services.UserAccountServiceInterface) UserAccountControllerInterface {
	return &UserAccountController{
		userAccountService: service,
	}
}

/* ============================== Implementationss ============================== */

func (c *UserAccountController) GetMyAccount(ctx *gin.Context, reqDto *dtos.GetMyAccountReqDto) {
	resDto, exception := c.userAccountService.GetMyAccount(ctx.Request.Context(), reqDto)
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

func (c *UserAccountController) UpdateMyAccount(ctx *gin.Context, reqDto *dtos.UpdateMyAccountReqDto) {
	resDto, exception := c.userAccountService.UpdateMyAccount(ctx.Request.Context(), reqDto)
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

func (c *UserAccountController) BindGoogleAccount(ctx *gin.Context, reqDto *dtos.BindGoogleAccountReqDto) {
	resDto, exception := c.userAccountService.BindGoogleAccount(ctx.Request.Context(), reqDto)
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

func (c *UserAccountController) UnbindGoogleAccount(ctx *gin.Context, reqDto *dtos.UnbindGoogleAccountReqDto) {
	resDto, exception := c.userAccountService.UnbindGoogleAccount(ctx.Request.Context(), reqDto)
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
