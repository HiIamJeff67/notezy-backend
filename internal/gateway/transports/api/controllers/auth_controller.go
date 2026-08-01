package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	authdto "github.com/HiIamJeff67/notezy-backend/contracts/api/v1/auth"
	responsewriter "github.com/HiIamJeff67/notezy-backend/internal/gateway/responsewriter"
	cookies "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/api/cookies"
	coreadapters "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/core/adapters"
)

type AuthControllerInterface interface {
	Register(ctx *gin.Context, requestDto *authdto.RegisterRequestDto)
	RegisterViaGoogle(ctx *gin.Context, requestDto *authdto.RegisterViaGoogleRequestDto)
	Login(ctx *gin.Context, requestDto *authdto.LoginRequestDto)
	LoginViaGoogle(ctx *gin.Context, requestDto *authdto.LoginViaGoogleRequestDto)
	Logout(ctx *gin.Context, requestDto *authdto.LogoutRequestDto)
	SendAuthCode(ctx *gin.Context, requestDto *authdto.SendAuthCodeRequestDto)
	ValidateEmail(ctx *gin.Context, requestDto *authdto.ValidateEmailRequestDto)
	ResetEmail(ctx *gin.Context, requestDto *authdto.ResetEmailRequestDto)
	ForgetPassword(ctx *gin.Context, requestDto *authdto.ForgetPasswordRequestDto)
	ResetMe(ctx *gin.Context, requestDto *authdto.ResetMeRequestDto)
	DeleteMe(ctx *gin.Context, requestDto *authdto.DeleteMeRequestDto)
}

type AuthController struct {
	coreClient *coreadapters.CoreClient
}

func NewAuthController(coreClient *coreadapters.CoreClient) AuthControllerInterface {
	return &AuthController{
		coreClient: coreClient,
	}
}

func (c *AuthController) Register(ctx *gin.Context, requestDto *authdto.RegisterRequestDto) {
	cookies.AccessTokenCookieHandler.Delete(ctx)
	cookies.RefreshTokenCookieHandler.Delete(ctx)

	response, exception := coreadapters.Call[authdto.RegisterRequestDto, authdto.RegisterResponseDto](
		ctx,
		c.coreClient,
		requestDto,
		authdto.RegisterOperation,
		"/core/v1/auth/register",
	)
	if exception != nil {
		responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	cookies.AccessTokenCookieHandler.Set(ctx, response.Data.AccessToken)
	cookies.RefreshTokenCookieHandler.Set(ctx, response.Data.RefreshToken)
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"publicId":    response.Data.PublicId,
			"name":        response.Data.Name,
			"displayName": response.Data.DisplayName,
			"email":       response.Data.Email,
			"accessToken": response.Data.AccessToken,
			"csrfToken":   response.Data.CSRFToken,
			"createdAt":   response.Data.CreatedAt,
		},
		"exception": nil,
	})
}

func (c *AuthController) RegisterViaGoogle(ctx *gin.Context, requestDto *authdto.RegisterViaGoogleRequestDto) {
	cookies.AccessTokenCookieHandler.Delete(ctx)
	cookies.RefreshTokenCookieHandler.Delete(ctx)

	response, exception := coreadapters.Call[authdto.RegisterViaGoogleRequestDto, authdto.RegisterViaGoogleResponseDto](
		ctx,
		c.coreClient,
		requestDto,
		authdto.RegisterViaGoogleOperation,
		"/core/v1/auth/register/google",
	)
	if exception != nil {
		responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	cookies.AccessTokenCookieHandler.Set(ctx, response.Data.AccessToken)
	cookies.RefreshTokenCookieHandler.Set(ctx, response.Data.RefreshToken)
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"publicId":    response.Data.PublicId,
			"name":        response.Data.Name,
			"displayName": response.Data.DisplayName,
			"email":       response.Data.Email,
			"accessToken": response.Data.AccessToken,
			"csrfToken":   response.Data.CSRFToken,
			"createdAt":   response.Data.CreatedAt,
		},
		"exception": nil,
	})
}

func (c *AuthController) Login(ctx *gin.Context, requestDto *authdto.LoginRequestDto) {
	cookies.AccessTokenCookieHandler.Delete(ctx)
	cookies.RefreshTokenCookieHandler.Delete(ctx)

	response, exception := coreadapters.Call[authdto.LoginRequestDto, authdto.LoginResponseDto](
		ctx,
		c.coreClient,
		requestDto,
		authdto.LoginOperation,
		"/core/v1/auth/login",
	)
	if exception != nil {
		responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	cookies.AccessTokenCookieHandler.Set(ctx, response.Data.AccessToken)
	cookies.RefreshTokenCookieHandler.Set(ctx, response.Data.RefreshToken)
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"publicId":    response.Data.PublicId,
			"name":        response.Data.Name,
			"displayName": response.Data.DisplayName,
			"email":       response.Data.Email,
			"accessToken": response.Data.AccessToken,
			"csrfToken":   response.Data.CSRFToken,
			"updatedAt":   response.Data.UpdatedAt,
			"createdAt":   response.Data.CreatedAt,
		},
		"exception": nil,
	})
}

func (c *AuthController) LoginViaGoogle(ctx *gin.Context, requestDto *authdto.LoginViaGoogleRequestDto) {
	cookies.AccessTokenCookieHandler.Delete(ctx)
	cookies.RefreshTokenCookieHandler.Delete(ctx)

	response, exception := coreadapters.Call[authdto.LoginViaGoogleRequestDto, authdto.LoginViaGoogleResponseDto](
		ctx,
		c.coreClient,
		requestDto,
		authdto.LoginViaGoogleOperation,
		"/core/v1/auth/login/google",
	)
	if exception != nil {
		responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	cookies.AccessTokenCookieHandler.Set(ctx, response.Data.AccessToken)
	cookies.RefreshTokenCookieHandler.Set(ctx, response.Data.RefreshToken)
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"publicId":    response.Data.PublicId,
			"name":        response.Data.Name,
			"displayName": response.Data.DisplayName,
			"email":       response.Data.Email,
			"accessToken": response.Data.AccessToken,
			"csrfToken":   response.Data.CSRFToken,
			"updatedAt":   response.Data.UpdatedAt,
			"createdAt":   response.Data.CreatedAt,
		},
		"exception": nil,
	})
}

func (c *AuthController) Logout(ctx *gin.Context, requestDto *authdto.LogoutRequestDto) {
	response, exception := coreadapters.CallSecurly[authdto.LogoutRequestDto, authdto.LogoutResponseDto](
		ctx,
		c.coreClient,
		requestDto,
		authdto.LogoutOperation,
		"/core/v1/auth/logout",
	)
	if exception != nil {
		responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	cookies.AccessTokenCookieHandler.Delete(ctx)
	cookies.RefreshTokenCookieHandler.Delete(ctx)
	ctx.JSON(http.StatusOK, gin.H{"success": true, "data": response.Data, "exception": nil})
}

func (c *AuthController) SendAuthCode(ctx *gin.Context, requestDto *authdto.SendAuthCodeRequestDto) {
	response, exception := coreadapters.Call[authdto.SendAuthCodeRequestDto, authdto.SendAuthCodeResponseDto](ctx, c.coreClient, requestDto, authdto.SendAuthCodeOperation, "/core/v1/auth/email/code")
	if exception != nil {
		responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": true, "data": response.Data, "exception": nil})
}

func (c *AuthController) ValidateEmail(ctx *gin.Context, requestDto *authdto.ValidateEmailRequestDto) {
	response, exception := coreadapters.CallSecurly[authdto.ValidateEmailRequestDto, authdto.ValidateEmailResponseDto](ctx, c.coreClient, requestDto, authdto.ValidateEmailOperation, "/core/v1/auth/email/validate")
	if exception != nil {
		responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": true, "data": response.Data, "exception": nil})
}

func (c *AuthController) ResetEmail(ctx *gin.Context, requestDto *authdto.ResetEmailRequestDto) {
	response, exception := coreadapters.CallSecurly[authdto.ResetEmailRequestDto, authdto.ResetEmailResponseDto](ctx, c.coreClient, requestDto, authdto.ResetEmailOperation, "/core/v1/auth/email/reset")
	if exception != nil {
		responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": true, "data": response.Data, "exception": nil})
}

func (c *AuthController) ForgetPassword(ctx *gin.Context, requestDto *authdto.ForgetPasswordRequestDto) {
	response, exception := coreadapters.Call[authdto.ForgetPasswordRequestDto, authdto.ForgetPasswordResponseDto](ctx, c.coreClient, requestDto, authdto.ForgetPasswordOperation, "/core/v1/auth/password/forget")
	if exception != nil {
		responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": true, "data": response.Data, "exception": nil})
}

func (c *AuthController) ResetMe(ctx *gin.Context, requestDto *authdto.ResetMeRequestDto) {
	response, exception := coreadapters.CallSecurly[authdto.ResetMeRequestDto, authdto.ResetMeResponseDto](ctx, c.coreClient, requestDto, authdto.ResetMeOperation, "/core/v1/auth/me/reset")
	if exception != nil {
		responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": true, "data": response.Data, "exception": nil})
}

func (c *AuthController) DeleteMe(ctx *gin.Context, requestDto *authdto.DeleteMeRequestDto) {
	response, exception := coreadapters.CallSecurly[authdto.DeleteMeRequestDto, authdto.DeleteMeResponseDto](ctx, c.coreClient, requestDto, authdto.DeleteMeOperation, "/core/v1/auth/me/delete")
	if exception != nil {
		responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": true, "data": response.Data, "exception": nil})
}
