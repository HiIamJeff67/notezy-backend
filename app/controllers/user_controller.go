package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	contexts "notezy-backend/app/contexts"
	dtos "notezy-backend/app/dtos"
	exceptions "notezy-backend/app/exceptions"
	logs "notezy-backend/app/logs"
	services "notezy-backend/app/services"
	constants "notezy-backend/shared/constants"
)

/* ============================== Interface & Instance ============================== */

type UserControllerInterface interface {
	GetUserData(ctx *gin.Context)
	GetMe(ctx *gin.Context)
	UpdateMe(ctx *gin.Context)
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
func (c *UserController) GetUserData(ctx *gin.Context) {
	var reqDto dtos.GetUserDataReqDto
	userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
	if exception != nil {
		exception.Log()
		exception = exceptions.User.InternalServerWentWrong(nil)
		ctx.JSON(exception.HTTPStatusCode, gin.H{
			"success":   false,
			"data":      nil,
			"exception": exception.GetGinH(),
		})
		return
	}
	reqDto.UserId = *userId

	resDto, exception := c.userService.GetUserData(&reqDto)
	if exception != nil {
		exception.Log()
		if exception.IsInternal {
			exception = exceptions.User.InternalServerWentWrong(nil)
		}
		ctx.JSON(exception.HTTPStatusCode, gin.H{
			"success":   false,
			"data":      nil,
			"exception": exception.GetGinH(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      resDto,
		"exception": nil,
	})
}

// with AuthMiddleware()
func (c *UserController) GetMe(ctx *gin.Context) {
	var reqDto dtos.GetMeReqDto
	userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
	if exception != nil {
		exception.Log()
		exception = exceptions.User.InternalServerWentWrong(nil)
		ctx.JSON(exception.HTTPStatusCode, gin.H{
			"success":   false,
			"data":      nil,
			"exception": exception.GetGinH(),
		})
		return
	}
	reqDto.UserId = *userId
	if err := ctx.ShouldBindJSON(&reqDto); err != nil {
		logs.Error(err)
		exception := exceptions.User.InvalidDto().WithError(err)
		ctx.JSON(exception.HTTPStatusCode, gin.H{
			"success":   false,
			"data":      nil,
			"exception": exception.GetGinH(),
		})
		return
	}

	resDto, exception := c.userService.GetMe(&reqDto)
	if exception != nil {
		exception.Log()
		if exception.IsInternal {
			exception = exceptions.User.InternalServerWentWrong(nil)
		}
		ctx.JSON(exception.HTTPStatusCode, gin.H{
			"success":   false,
			"data":      nil,
			"exception": exception.GetGinH(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      resDto,
		"exception": nil,
	})
}

// with AuthMiddleware()
func (c *UserController) UpdateMe(ctx *gin.Context) {
	var reqDto dtos.UpdateMeReqDto
	userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
	if exception != nil {
		exception.Log()
		exception = exceptions.User.InternalServerWentWrong(nil)
		ctx.JSON(exception.HTTPStatusCode, gin.H{
			"success":   false,
			"data":      nil,
			"exception": exception.GetGinH(),
		})
		return
	}
	reqDto.UserId = *userId
	if err := ctx.ShouldBindJSON(&reqDto); err != nil {
		exception := exceptions.User.InvalidDto().WithError(err)
		ctx.JSON(exception.HTTPStatusCode, gin.H{
			"success":   false,
			"data":      nil,
			"exception": exception.GetGinH(),
		})
		return
	}

	resDto, exception := c.userService.UpdateMe(&reqDto)
	if exception != nil {
		exception.Log()
		if exception.IsInternal {
			exception = exceptions.User.InternalServerWentWrong(nil)
		}
		ctx.JSON(exception.HTTPStatusCode, gin.H{
			"success":   false,
			"data":      nil,
			"exception": exception.GetGinH(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      resDto,
		"exception": nil,
	})
}
