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

type UserControllerInterface interface {
	GetMe(ctx *gin.Context)
	GetAllUsers(ctx *gin.Context)
	UpdateMe(ctx *gin.Context)
}

type userController struct {
	userService services.UserServiceInterface
}

var UserController UserControllerInterface = &userController{
	userService: services.UserService,
}

/* ============================== Controllers ============================== */

// with AuthMiddleware()
func (c *userController) GetMe(ctx *gin.Context) {
	var reqDto dtos.GetMeReqDto
	userId, exception := contexts.FetchAndConvertContextFieldToUUID(ctx, "userId")
	if exception != nil {
		exception.Log()
		exception = exceptions.User.InternalServerWentWrong(exception)
		ctx.JSON(exception.HTTPStatusCode, gin.H{
			"success":   false,
			"data":      nil,
			"exception": exception.GetGinH(),
		})
		return
	}
	reqDto.UserId = *userId

	resDto, exception := c.userService.GetMe(&reqDto)
	if exception != nil {
		exception.Log()
		if !exceptions.CompareCommonExceptions(exceptions.User.InvalidDto(), exception, false) {
			exception = exceptions.User.InternalServerWentWrong(exception)
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
func (c *userController) GetAllUsers(ctx *gin.Context) {
	resDto, exception := c.userService.GetAllUsers()
	if exception != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": exception.Log().Error})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      resDto,
		"exception": nil,
	})
}

// with AuthMiddleware()
func (c *userController) UpdateMe(ctx *gin.Context) {
	var reqDto dtos.UpdateMeReqDto
	userId, exception := contexts.FetchAndConvertContextFieldToUUID(ctx, "userId")
	if exception != nil {
		exception.Log()
		exception = exceptions.User.InternalServerWentWrong(exception)
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
		if !exceptions.CompareCommonExceptions(exceptions.User.InvalidDto(), exception, false) {
			exception = exceptions.User.InternalServerWentWrong(exception)
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
