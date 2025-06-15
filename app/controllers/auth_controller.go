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

/* ============================== Controller ============================== */
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
		"data": gin.H{
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
		"data": gin.H{
			"accessToken": resDto.AccessToken,
			"updatedAt":   resDto.UpdatedAt,
		},
	})
}

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
		"data": gin.H{
			"updatedAt": resDto.UpdatedAt,
		},
	})
}
