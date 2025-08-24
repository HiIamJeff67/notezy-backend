package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	contexts "notezy-backend/app/contexts"
	cookies "notezy-backend/app/cookies"
	dtos "notezy-backend/app/dtos"
	exceptions "notezy-backend/app/exceptions"
	services "notezy-backend/app/services"
	constants "notezy-backend/shared/constants"
)

/* ============================== Interface & Instance ============================== */

type AuthControllerInterface interface {
	Register(ctx *gin.Context)
	Login(ctx *gin.Context)
	Logout(ctx *gin.Context)
	SendAuthCode(ctx *gin.Context)
	ValidateEmail(ctx *gin.Context)
	ResetEmail(ctx *gin.Context)
	ForgetPassword(ctx *gin.Context)
	DeleteMe(ctx *gin.Context)
}

type AuthController struct {
	authService services.AuthServiceInterface
}

func NewAuthController(service services.AuthServiceInterface) AuthControllerInterface {
	return &AuthController{
		authService: service,
	}
}

/* ============================== Controllers ============================== */

func (c *AuthController) Register(ctx *gin.Context) {
	var reqDto dtos.RegisterReqDto
	reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")
	if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
		exception := exceptions.Auth.InvalidDto().WithError(err)
		exception.ResponseWithJSON(ctx)
		return
	}

	cookies.AccessToken.DeleteCookie(ctx)
	cookies.RefreshToken.DeleteCookie(ctx)

	resDto, exception := c.authService.Register(&reqDto)
	if exception != nil {
		exception.Log().SafelyResponseWithJSON(ctx)
		return
	}

	cookies.AccessToken.SetCookie(ctx, resDto.AccessToken)
	cookies.RefreshToken.SetCookie(ctx, resDto.RefreshToken)

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{ // make sure we don't response with the refresh token
			"accessToken": resDto.AccessToken,
			"createdAt":   resDto.CreatedAt,
		},
		"exception": nil,
	})
}

func (c *AuthController) Login(ctx *gin.Context) {
	var reqDto dtos.LoginReqDto
	reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")
	if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
		exception := exceptions.Auth.InvalidDto().WithError(err)
		exception.ResponseWithJSON(ctx)
		return
	}

	cookies.AccessToken.DeleteCookie(ctx)
	cookies.RefreshToken.DeleteCookie(ctx)

	resDto, exception := c.authService.Login(&reqDto)
	if exception != nil {
		exception.Log().SafelyResponseWithJSON(ctx)
		return
	}

	cookies.AccessToken.SetCookie(ctx, resDto.AccessToken)
	cookies.RefreshToken.SetCookie(ctx, resDto.RefreshToken)

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{ // make sure we don't response with the refresh token
			"accessToken": resDto.AccessToken,
			"updatedAt":   resDto.UpdatedAt,
		},
		"exception": nil,
	})
}

// with AuthMiddleware()
func (c *AuthController) Logout(ctx *gin.Context) {
	var reqDto dtos.LogoutReqDto
	userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
	if exception != nil {
		exception.Log().SafelyResponseWithJSON(ctx)
		return
	}
	reqDto.ContextFields.UserId = *userId

	resDto, exception := c.authService.Logout(&reqDto)
	if exception != nil {
		exception.Log().SafelyResponseWithJSON(ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      resDto,
		"exceptoin": nil,
	})
}

func (c *AuthController) SendAuthCode(ctx *gin.Context) {
	var reqDto dtos.SendAuthCodeReqDto
	reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")
	if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
		exception := exceptions.Auth.InvalidDto().WithError(err)
		exception.ResponseWithJSON(ctx)
		return
	}

	resDto, exception := c.authService.SendAuthCode(&reqDto)
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
func (c *AuthController) ValidateEmail(ctx *gin.Context) {
	var reqDto dtos.ValidateEmailReqDto
	userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
	if exception != nil {
		exception.Log().SafelyResponseWithJSON(ctx)
		return
	}
	reqDto.ContextFields.UserId = *userId
	if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
		exception := exceptions.Auth.InvalidDto().WithError(err)
		exception.ResponseWithJSON(ctx)
		return
	}

	resDto, exception := c.authService.ValidateEmail(&reqDto)
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
func (c *AuthController) ResetEmail(ctx *gin.Context) {
	var reqDto dtos.ResetEmailReqDto
	userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
	if exception != nil {
		exception.Log().SafelyResponseWithJSON(ctx)
		return
	}
	reqDto.ContextFields.UserId = *userId
	if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
		exception := exceptions.Auth.InvalidDto().WithError(err)
		exception.ResponseWithJSON(ctx)
		return
	}

	resDto, exception := c.authService.ResetEmail(&reqDto)
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

// ! this should not use any middleware, bcs we want the user to set it by providing the account
func (c *AuthController) ForgetPassword(ctx *gin.Context) {
	var reqDto dtos.ForgetPasswordReqDto
	reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")
	if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
		exception := exceptions.Auth.InvalidDto().WithError(err)
		exception.ResponseWithJSON(ctx)
		return
	}

	resDto, exception := c.authService.ForgetPassword(&reqDto)
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
func (c *AuthController) DeleteMe(ctx *gin.Context) {
	var reqDto dtos.DeleteMeReqDto
	userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
	if exception != nil {
		exception.Log().SafelyResponseWithJSON(ctx)
		return
	}
	reqDto.ContextFields.UserId = *userId
	if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
		exception := exceptions.Auth.InvalidDto().WithError(err)
		exception.ResponseWithJSON(ctx)
		return
	}

	resDto, exception := c.authService.DeleteMe(&reqDto)
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
