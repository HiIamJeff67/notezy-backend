package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	contexts "notezy-backend/app/contexts"
	cookies "notezy-backend/app/cookies"
	dtos "notezy-backend/app/dtos"
	exceptions "notezy-backend/app/exceptions"
	services "notezy-backend/app/services"
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

type authController struct {
	authService services.AuthServiceInterface
}

var AuthController AuthControllerInterface = &authController{
	authService: services.AuthService,
}

/* ============================== Controllers ============================== */

func (c *authController) Register(ctx *gin.Context) {
	var reqDto dtos.RegisterReqDto
	reqDto.UserAgent = ctx.GetHeader("User-Agent")
	if err := ctx.ShouldBindJSON(&reqDto); err != nil {
		exception := exceptions.Auth.InvalidDto().WithError(err)
		ctx.JSON(exception.HTTPStatusCode, exception.GetGinH())
		return
	}

	resDto, exception := c.authService.Register(&reqDto)
	if exception != nil {
		exception.Log()
		if !exceptions.CompareCommonExceptions(exceptions.Auth.InvalidDto(), exception, false) {
			exception = exceptions.Auth.InternalServerWentWrong(exception)
		}
		ctx.JSON(exception.HTTPStatusCode, exception.GetGinH())
		return
	}

	cookies.AccessToken.SetCookie(ctx, resDto.AccessToken)
	cookies.RefreshToken.SetCookie(ctx, resDto.RefreshToken)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data": gin.H{ // make sure we don't response with the refresh token
			"accessToken": resDto.AccessToken,
			"createdAt":   resDto.CreatedAt,
		},
	})
}

func (c *authController) Login(ctx *gin.Context) {
	var reqDto dtos.LoginReqDto
	reqDto.UserAgent = ctx.GetHeader("User-Agent")
	if err := ctx.ShouldBindJSON(&reqDto); err != nil {
		exception := exceptions.Auth.InvalidDto().WithError(err)
		ctx.JSON(exception.HTTPStatusCode, exception.GetGinH())
		return
	}

	resDto, exception := c.authService.Login(&reqDto)
	if exception != nil {
		exception.Log()
		if !exceptions.CompareCommonExceptions(exceptions.Auth.InvalidDto(), exception, false) {
			exception = exceptions.Auth.InternalServerWentWrong(exception)
		}
		ctx.JSON(exception.HTTPStatusCode, exception.GetGinH())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    resDto,
	})
}

// with AuthMiddleware()
func (c *authController) Logout(ctx *gin.Context) {
	var reqDto dtos.LogoutReqDto
	userId, exception := contexts.FetchAndConvertContextFieldToUUID(ctx, "userId")
	if exception != nil {
		exception.Log()
		exception = exceptions.Auth.InternalServerWentWrong(exception)
		ctx.JSON(exception.HTTPStatusCode, exception.GetGinH())
		return
	}
	reqDto.UserId = *userId

	resDto, exception := c.authService.Logout(&reqDto)
	if exception != nil {
		exception.Log()
		if !exceptions.CompareCommonExceptions(exceptions.Auth.InvalidDto(), exception, false) {
			exception = exceptions.Auth.InternalServerWentWrong(exception)
		}
		ctx.JSON(exception.HTTPStatusCode, exception.GetGinH())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    resDto,
	})
}

func (c *authController) SendAuthCode(ctx *gin.Context) {
	var reqDto dtos.SendAuthCodeReqDto
	reqDto.UserAgent = ctx.GetHeader("User-Agent")
	if err := ctx.ShouldBindJSON(&reqDto); err != nil {
		exception := exceptions.Auth.InvalidDto().WithError(err)
		ctx.JSON(exception.HTTPStatusCode, exception.GetGinH())
		return
	}

	resDto, exception := c.authService.SendAuthCode(&reqDto)
	if exception != nil {
		exception.Log()
		if !exceptions.CompareCommonExceptions(exceptions.Auth.InvalidDto(), exception, false) {
			exception = exceptions.Auth.InternalServerWentWrong(exception)
		}
		ctx.JSON(exception.HTTPStatusCode, exception.GetGinH())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    resDto,
	})
}

// with AuthMiddleware()
func (c *authController) ValidateEmail(ctx *gin.Context) {
	var reqDto dtos.ValidateEmailReqDto
	userId, exception := contexts.FetchAndConvertContextFieldToUUID(ctx, "userId")
	if exception != nil {
		exception.Log()
		exception = exceptions.Auth.InternalServerWentWrong(exception)
		ctx.JSON(exception.HTTPStatusCode, exception.GetGinH())
		return
	}
	reqDto.UserId = *userId
	if err := ctx.ShouldBindJSON(&reqDto); err != nil {
		exception := exceptions.Auth.InvalidDto().WithError(err)
		ctx.JSON(exception.HTTPStatusCode, exception.GetGinH())
		return
	}

	resDto, exception := c.authService.ValidateEmail(&reqDto)
	if exception != nil {
		exception.Log()
		if !exceptions.CompareCommonExceptions(exceptions.Auth.InvalidDto(), exception, false) {
			exception = exceptions.Auth.InternalServerWentWrong(exception)
		}
		ctx.JSON(exception.HTTPStatusCode, exception.GetGinH())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    resDto,
	})
}

// with AuthMiddleware()
func (c *authController) ResetEmail(ctx *gin.Context) {
	var reqDto dtos.ResetEmailReqDto
	userId, exception := contexts.FetchAndConvertContextFieldToUUID(ctx, "userId")
	if exception != nil {
		exception.Log()
		exception = exceptions.Auth.InternalServerWentWrong(exception)
		ctx.JSON(exception.HTTPStatusCode, exception.GetGinH())
		return
	}
	reqDto.UserId = *userId
	if err := ctx.ShouldBindJSON(&reqDto); err != nil {
		exception := exceptions.Auth.InvalidDto().WithError(err)
		ctx.JSON(exception.HTTPStatusCode, exception.GetGinH())
		return
	}

	resDto, exception := c.authService.ResetEmail(&reqDto)
	if exception != nil {
		exception.Log()
		if !exceptions.CompareCommonExceptions(exceptions.Auth.InvalidDto(), exception, false) {
			exception = exceptions.Auth.InternalServerWentWrong(exception)
		}
		ctx.JSON(exception.HTTPStatusCode, exception.GetGinH())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    resDto,
	})
}

// ! this should not use any middleware, bcs we want the user to set it by providing the account
func (c *authController) ForgetPassword(ctx *gin.Context) {
	var reqDto dtos.ForgetPasswordReqDto
	reqDto.UserAgent = ctx.GetHeader("User-Agent")
	if err := ctx.ShouldBindJSON(&reqDto); err != nil {
		exception := exceptions.Auth.InvalidDto().WithError(err)
		ctx.JSON(exception.HTTPStatusCode, exception.GetGinH())
		return
	}

	resDto, exception := c.authService.ForgetPassword(&reqDto)
	if exception != nil {
		exception.Log()
		if !exceptions.CompareCommonExceptions(exceptions.Auth.InvalidDto(), exception, false) {
			exception = exceptions.Auth.InternalServerWentWrong(exception)
		}
		ctx.JSON(exception.HTTPStatusCode, exception.GetGinH())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    resDto,
	})
}

// with AuthMiddleware()
func (c *authController) DeleteMe(ctx *gin.Context) {
	var reqDto dtos.DeleteMeReqDto
	userId, exception := contexts.FetchAndConvertContextFieldToUUID(ctx, "userId")
	if exception != nil {
		exception.Log()
		exception = exceptions.Auth.InternalServerWentWrong(exception)
		ctx.JSON(exception.HTTPStatusCode, exception.GetGinH())
		return
	}
	reqDto.UserId = *userId
	if err := ctx.ShouldBindJSON(&reqDto); err != nil {
		exception := exceptions.Auth.InvalidDto().WithError(err)
		ctx.JSON(exception.HTTPStatusCode, exception.GetGinH())
		return
	}

	resDto, exception := c.authService.DeleteMe(&reqDto)
	if exception != nil {
		exception.Log()
		if !exceptions.CompareCommonExceptions(exceptions.Auth.InvalidDto(), exception, false) {
			exception = exceptions.Auth.InternalServerWentWrong(exception)
		}
		ctx.JSON(exception.HTTPStatusCode, exception.GetGinH())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    resDto,
	})
}
