package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"notezy-backend/app/contexts"
	cookies "notezy-backend/app/cookies"
	dtos "notezy-backend/app/dtos"
	exceptions "notezy-backend/app/exceptions"
	services "notezy-backend/app/services"
)

func Register(ctx *gin.Context) {
	var reqDto dtos.RegisterReqDto
	reqDto.UserAgent = ctx.GetHeader("User-Agent")
	if err := ctx.ShouldBindJSON(&reqDto); err != nil {
		ctx.JSON(http.StatusBadRequest, exceptions.Auth.InvalidDto().WithError(err).GetGinH())
		return
	}

	resDto, exception := services.Register(&reqDto)
	if exception != nil {
		ctx.JSON(exception.HTTPStatusCode, exception.GetGinH())
		return
	}

	cookies.AccessToken.SetCookie(ctx, resDto.AccessToken)
	cookies.RefreshToken.SetCookie(ctx, resDto.RefreshToken)

	ctx.JSON(http.StatusOK, gin.H{
		"data": gin.H{ // make sure we don't response with the refresh token
			"accessToken": resDto.AccessToken,
			"createdAt":   resDto.CreatedAt,
		},
	})
}

func Login(ctx *gin.Context) {
	var reqDto dtos.LoginReqDto
	reqDto.UserAgent = ctx.GetHeader("User-Agent")
	if err := ctx.ShouldBindJSON(&reqDto); err != nil {
		ctx.JSON(
			exceptions.Auth.InvalidDto().HTTPStatusCode,
			exceptions.Auth.InvalidDto().WithError(err).GetGinH(),
		)
		return
	}

	resDto, exception := services.Login(&reqDto)
	if exception != nil {
		ctx.JSON(
			exception.HTTPStatusCode,
			exception.GetGinH(),
		)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": resDto,
	})
}

// with AuthMiddleware()
func Logout(ctx *gin.Context) {
	var reqDto dtos.LogoutReqDto
	userId, exception := contexts.FetchAndConvertContextFieldToUUID(ctx, "userId")
	if exception != nil {
		exception.Log() // only log the error without throwing them back to client
		ctx.JSON(
			exceptions.Auth.InvalidDto().HTTPStatusCode,
			exceptions.Auth.InvalidDto().GetGinH(),
		)
		return
	}
	reqDto.UserId = *userId

	resDto, exception := services.Logout(&reqDto)
	if exception != nil {
		ctx.JSON(
			exception.HTTPStatusCode,
			exception.GetGinH(),
		)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": resDto,
	})
}

func SendAuthCode(ctx *gin.Context) {
	var reqDto dtos.SendAuthCodeReqDto
	reqDto.UserAgent = ctx.GetHeader("User-Agent")
	if err := ctx.ShouldBindJSON(&reqDto); err != nil {
		exception := exceptions.Auth.InvalidDto().WithError(err)
		ctx.JSON(exception.HTTPStatusCode, exception.GetGinH())
		return
	}

	resDto, exception := services.SendAuthCode(&reqDto)
	if exception != nil {
		ctx.JSON(exception.HTTPStatusCode, exception.GetGinH())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": resDto,
	})
}

// with AuthMiddleware()
func ValidateEmail(ctx *gin.Context) {
	var reqDto dtos.ValidateEmailReqDto
	userId, exception := contexts.FetchAndConvertContextFieldToUUID(ctx, "userId")
	if exception != nil {
		ctx.JSON(exception.HTTPStatusCode, exception.GetGinH())
		return
	}
	reqDto.UserId = *userId

	resDto, exception := services.ValidateEmail(&reqDto)
	if exception != nil {
		ctx.JSON(exception.HTTPStatusCode, exception.GetGinH())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": resDto,
	})
}

// with AuthMiddleware()
func ResetEmail(ctx *gin.Context) {
	var reqDto dtos.ResetEmailReqDto
	userId, exception := contexts.FetchAndConvertContextFieldToUUID(ctx, "userId")
	if exception != nil {
		ctx.JSON(exception.HTTPStatusCode, exception.GetGinH())
		return
	}
	reqDto.UserId = *userId
	if err := ctx.ShouldBindJSON(&reqDto); err != nil {
		exception := exceptions.Auth.InvalidDto().WithError(err)
		ctx.JSON(exception.HTTPStatusCode, exception.GetGinH())
		return
	}

	resDto, exception := services.ResetEmail(&reqDto)
	if exception != nil {
		ctx.JSON(exception.HTTPStatusCode, exception.GetGinH())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": resDto,
	})
}

// ! this should not use any middleware, bcs we want the user to set it by providing the account
func ForgetPassword(ctx *gin.Context) {
	var reqDto dtos.ForgetPasswordReqDto
	reqDto.UserAgent = ctx.GetHeader("User-Agent")
	if err := ctx.ShouldBindJSON(&reqDto); err != nil {
		exception := exceptions.Auth.InvalidDto().WithError(err)
		ctx.JSON(exception.HTTPStatusCode, exception.GetGinH())
		return
	}

	resDto, exception := services.ForgetPassword(&reqDto)
	if exception != nil {
		ctx.JSON(exception.HTTPStatusCode, exception.GetGinH())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": resDto,
	})
}
