package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	cookies "notezy-backend/app/cookie"
	dtos "notezy-backend/app/dtos"
	exceptions "notezy-backend/app/exceptions"
	services "notezy-backend/app/services"
)

/* ============================== Controller ============================== */
func Register(ctx *gin.Context) {
	var reqDto dtos.RegisterReqDto
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
	accessTokenFromCookie, exists := ctx.Get("accessToken")
	if !exists {
		ctx.JSON(
			exceptions.Auth.InvalidDto().HTTPStatusCode,
			exceptions.Auth.InvalidDto().GetGinH(),
		)
		return
	}

	tokenStr, ok := accessTokenFromCookie.(string)
	if !ok {
		ctx.JSON(
			exceptions.Auth.InvalidDto().HTTPStatusCode,
			exceptions.Auth.InvalidDto().GetGinH(),
		)
		return
	}

	reqDto.AccessToken = tokenStr
	// use the below code if we can't use the cookies of the user
	// if err := ctx.ShouldBindJSON(&reqDto); err != nil {
	// 	ctx.JSON(
	// 		exceptions.Auth.InvalidDto().HTTPStatusCode,
	// 		exceptions.Auth.InvalidDto().WithError(err).GetGinH(),
	// 	)
	// 	return
	// }

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
