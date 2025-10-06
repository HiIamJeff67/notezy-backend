package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	cookies "notezy-backend/app/cookies"
	dtos "notezy-backend/app/dtos"
	services "notezy-backend/app/services"
)

/* ============================== Interface & Instance ============================== */

type AuthControllerInterface interface {
	Register(ctx *gin.Context, reqDto *dtos.RegisterReqDto)
	Login(ctx *gin.Context, reqDto *dtos.LoginReqDto)
	Logout(ctx *gin.Context, reqDto *dtos.LogoutReqDto)
	SendAuthCode(ctx *gin.Context, reqDto *dtos.SendAuthCodeReqDto)
	ValidateEmail(ctx *gin.Context, reqDto *dtos.ValidateEmailReqDto)
	ResetEmail(ctx *gin.Context, reqDto *dtos.ResetEmailReqDto)
	ForgetPassword(ctx *gin.Context, reqDto *dtos.ForgetPasswordReqDto)
	DeleteMe(ctx *gin.Context, reqDto *dtos.DeleteMeReqDto)
}

type AuthController struct {
	authService services.AuthServiceInterface
}

func NewAuthController(service services.AuthServiceInterface) AuthControllerInterface {
	return &AuthController{
		authService: service,
	}
}

/* ============================== Controller ============================== */

func (c *AuthController) Register(ctx *gin.Context, reqDto *dtos.RegisterReqDto) {
	cookies.AccessToken.DeleteCookie(ctx)
	cookies.RefreshToken.DeleteCookie(ctx)

	resDto, exception := c.authService.Register(reqDto)
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

func (c *AuthController) Login(ctx *gin.Context, reqDto *dtos.LoginReqDto) {
	cookies.AccessToken.DeleteCookie(ctx)
	cookies.RefreshToken.DeleteCookie(ctx)

	resDto, exception := c.authService.Login(reqDto)
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

// with AuthMiddleware
func (c *AuthController) Logout(ctx *gin.Context, reqDto *dtos.LogoutReqDto) {
	resDto, exception := c.authService.Logout(reqDto)
	if exception != nil {
		exception.Log().SafelyResponseWithJSON(ctx)
		return
	}

	cookies.AccessToken.DeleteCookie(ctx)
	cookies.RefreshToken.DeleteCookie(ctx)

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      resDto,
		"exceptoin": nil,
	})
}

func (c *AuthController) SendAuthCode(ctx *gin.Context, reqDto *dtos.SendAuthCodeReqDto) {
	resDto, exception := c.authService.SendAuthCode(reqDto)
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
func (c *AuthController) ValidateEmail(ctx *gin.Context, reqDto *dtos.ValidateEmailReqDto) {
	resDto, exception := c.authService.ValidateEmail(reqDto)
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
func (c *AuthController) ResetEmail(ctx *gin.Context, reqDto *dtos.ResetEmailReqDto) {
	resDto, exception := c.authService.ResetEmail(reqDto)
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
func (c *AuthController) ForgetPassword(ctx *gin.Context, reqDto *dtos.ForgetPasswordReqDto) {
	resDto, exception := c.authService.ForgetPassword(reqDto)
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
func (c *AuthController) DeleteMe(ctx *gin.Context, reqDto *dtos.DeleteMeReqDto) {
	resDto, exception := c.authService.DeleteMe(reqDto)
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
