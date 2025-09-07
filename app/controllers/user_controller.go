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

/* ============================== Controllers ============================== */

// with AuthMiddleware()
func (c *UserController) GetUserData(ctx *gin.Context, reqDto *dtos.GetUserDataReqDto) {
	resDto, exception := c.userService.GetUserData(reqDto)
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
func (c *UserController) GetMe(ctx *gin.Context, reqDto *dtos.GetMeReqDto) {
	resDto, exception := c.userService.GetMe(reqDto)
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
func (c *UserController) UpdateMe(ctx *gin.Context, reqDto *dtos.UpdateMeReqDto) {
	resDto, exception := c.userService.UpdateMe(reqDto)
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
