package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	dtos "notezy-backend/app/dtos"
	services "notezy-backend/app/services"
)

/* ============================== Interface & Instance ============================== */

type UserControllerInterface interface {
	GetUserData(ctx *gin.Context, reqDto *dtos.GetUserDataReqDto)
	GetMe(ctx *gin.Context, reqDto *dtos.GetMeReqDto)
	UpdateMe(ctx *gin.Context, reqDto *dtos.UpdateMeReqDto)
}

type UserController struct {
	userService services.UserServiceInterface
}

func NewUserController(service services.UserServiceInterface) UserControllerInterface {
	return &UserController{
		userService: service,
	}
}

/* ============================== Implementationss ============================== */

func (c *UserController) GetUserData(ctx *gin.Context, reqDto *dtos.GetUserDataReqDto) {
	resDto, exception := c.userService.GetUserData(ctx.Request.Context(), reqDto)
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

func (c *UserController) GetMe(ctx *gin.Context, reqDto *dtos.GetMeReqDto) {
	resDto, exception := c.userService.GetMe(ctx.Request.Context(), reqDto)
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

func (c *UserController) UpdateMe(ctx *gin.Context, reqDto *dtos.UpdateMeReqDto) {
	resDto, exception := c.userService.UpdateMe(ctx.Request.Context(), reqDto)
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
