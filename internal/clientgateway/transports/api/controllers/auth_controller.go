package controllers

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	cookies "github.com/HiIamJeff67/notegic-backend/shared/cookies"

	exceptionwriter "github.com/HiIamJeff67/notegic-backend/shared/util/exceptionwriter"

	apicontract "github.com/HiIamJeff67/notegic-backend/contracts/core/v1/api/auth"

	coreadapters "github.com/HiIamJeff67/notegic-backend/internal/clientgateway/transports/core/adapters"
)

type AuthControllerInterface interface {
	Register(ctx *gin.Context, requestDto *apicontract.RegisterRequestDto)
	RegisterViaGoogle(ctx *gin.Context, requestDto *apicontract.RegisterViaGoogleRequestDto)
	Login(ctx *gin.Context, requestDto *apicontract.LoginRequestDto)
	LoginViaGoogle(ctx *gin.Context, requestDto *apicontract.LoginViaGoogleRequestDto)
	Logout(ctx *gin.Context, requestDto *apicontract.LogoutRequestDto)
	SendAuthCode(ctx *gin.Context, requestDto *apicontract.SendAuthCodeRequestDto)
	ValidateEmail(ctx *gin.Context, requestDto *apicontract.ValidateEmailRequestDto)
	ResetEmail(ctx *gin.Context, requestDto *apicontract.ResetEmailRequestDto)
	ForgetPassword(ctx *gin.Context, requestDto *apicontract.ForgetPasswordRequestDto)
	ResetMe(ctx *gin.Context, requestDto *apicontract.ResetMeRequestDto)
	DeleteMe(ctx *gin.Context, requestDto *apicontract.DeleteMeRequestDto)
}

type AuthController struct {
	coreAdapter               *coreadapters.CoreAdapter
	accessTokenCookieHandler  *cookies.CookieHandler
	refreshTokenCookieHandler *cookies.CookieHandler
}

func NewAuthController(
	coreAdapter *coreadapters.CoreAdapter,
	accessTokenCookieHandler *cookies.CookieHandler,
	refreshTokenCookieHandler *cookies.CookieHandler,
) AuthControllerInterface {
	return &AuthController{
		coreAdapter:               coreAdapter,
		accessTokenCookieHandler:  accessTokenCookieHandler,
		refreshTokenCookieHandler: refreshTokenCookieHandler,
	}
}

func (c *AuthController) Register(ctx *gin.Context, requestDto *apicontract.RegisterRequestDto) {
	c.accessTokenCookieHandler.Delete(ctx)
	c.refreshTokenCookieHandler.Delete(ctx)

	response, exception := coreadapters.Call[apicontract.RegisterRequestDto, apicontract.RegisterResponseDto](
		ctx,
		c.coreAdapter,
		requestDto,
		apicontract.RegisterOperation,
		"/core/v1/auth/register",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	c.accessTokenCookieHandler.Set(ctx, response.Data.AccessToken)
	c.refreshTokenCookieHandler.Set(ctx, response.Data.RefreshToken)
	writeClientResponse(ctx, struct {
		PublicId    uuid.UUID `json:"publicId"`
		Name        string    `json:"name"`
		DisplayName string    `json:"displayName"`
		Email       string    `json:"email"`
		CSRFToken   string    `json:"csrfToken"`
		CreatedAt   time.Time `json:"createdAt"`
	}{
		PublicId:    response.Data.PublicId,
		Name:        response.Data.Name,
		DisplayName: response.Data.DisplayName,
		Email:       response.Data.Email,
		CSRFToken:   response.Data.CSRFToken,
		CreatedAt:   response.Data.CreatedAt,
	})
}

func (c *AuthController) RegisterViaGoogle(ctx *gin.Context, requestDto *apicontract.RegisterViaGoogleRequestDto) {
	c.accessTokenCookieHandler.Delete(ctx)
	c.refreshTokenCookieHandler.Delete(ctx)

	response, exception := coreadapters.Call[apicontract.RegisterViaGoogleRequestDto, apicontract.RegisterViaGoogleResponseDto](
		ctx,
		c.coreAdapter,
		requestDto,
		apicontract.RegisterViaGoogleOperation,
		"/core/v1/auth/register/google",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	c.accessTokenCookieHandler.Set(ctx, response.Data.AccessToken)
	c.refreshTokenCookieHandler.Set(ctx, response.Data.RefreshToken)
	writeClientResponse(ctx, struct {
		PublicId    uuid.UUID `json:"publicId"`
		Name        string    `json:"name"`
		DisplayName string    `json:"displayName"`
		Email       string    `json:"email"`
		CSRFToken   string    `json:"csrfToken"`
		CreatedAt   time.Time `json:"createdAt"`
	}{
		PublicId:    response.Data.PublicId,
		Name:        response.Data.Name,
		DisplayName: response.Data.DisplayName,
		Email:       response.Data.Email,
		CSRFToken:   response.Data.CSRFToken,
		CreatedAt:   response.Data.CreatedAt,
	})
}

func (c *AuthController) Login(ctx *gin.Context, requestDto *apicontract.LoginRequestDto) {
	c.accessTokenCookieHandler.Delete(ctx)
	c.refreshTokenCookieHandler.Delete(ctx)

	response, exception := coreadapters.Call[apicontract.LoginRequestDto, apicontract.LoginResponseDto](
		ctx,
		c.coreAdapter,
		requestDto,
		apicontract.LoginOperation,
		"/core/v1/auth/login",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	c.accessTokenCookieHandler.Set(ctx, response.Data.AccessToken)
	c.refreshTokenCookieHandler.Set(ctx, response.Data.RefreshToken)
	updatedAt := response.Data.UpdatedAt
	writeClientResponse(ctx, struct {
		PublicId    uuid.UUID  `json:"publicId"`
		Name        string     `json:"name"`
		DisplayName string     `json:"displayName"`
		Email       string     `json:"email"`
		CSRFToken   string     `json:"csrfToken"`
		UpdatedAt   *time.Time `json:"updatedAt,omitempty"`
		CreatedAt   time.Time  `json:"createdAt"`
	}{
		PublicId:    response.Data.PublicId,
		Name:        response.Data.Name,
		DisplayName: response.Data.DisplayName,
		Email:       response.Data.Email,
		CSRFToken:   response.Data.CSRFToken,
		UpdatedAt:   &updatedAt,
		CreatedAt:   response.Data.CreatedAt,
	})
}

func (c *AuthController) LoginViaGoogle(ctx *gin.Context, requestDto *apicontract.LoginViaGoogleRequestDto) {
	c.accessTokenCookieHandler.Delete(ctx)
	c.refreshTokenCookieHandler.Delete(ctx)

	response, exception := coreadapters.Call[apicontract.LoginViaGoogleRequestDto, apicontract.LoginViaGoogleResponseDto](
		ctx,
		c.coreAdapter,
		requestDto,
		apicontract.LoginViaGoogleOperation,
		"/core/v1/auth/login/google",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	c.accessTokenCookieHandler.Set(ctx, response.Data.AccessToken)
	c.refreshTokenCookieHandler.Set(ctx, response.Data.RefreshToken)
	updatedAt := response.Data.UpdatedAt
	writeClientResponse(ctx, struct {
		PublicId    uuid.UUID  `json:"publicId"`
		Name        string     `json:"name"`
		DisplayName string     `json:"displayName"`
		Email       string     `json:"email"`
		CSRFToken   string     `json:"csrfToken"`
		UpdatedAt   *time.Time `json:"updatedAt,omitempty"`
		CreatedAt   time.Time  `json:"createdAt"`
	}{
		PublicId:    response.Data.PublicId,
		Name:        response.Data.Name,
		DisplayName: response.Data.DisplayName,
		Email:       response.Data.Email,
		CSRFToken:   response.Data.CSRFToken,
		UpdatedAt:   &updatedAt,
		CreatedAt:   response.Data.CreatedAt,
	})
}

func (c *AuthController) Logout(ctx *gin.Context, requestDto *apicontract.LogoutRequestDto) {
	response, exception := coreadapters.CallSecurly[apicontract.LogoutRequestDto, apicontract.LogoutResponseDto](
		ctx,
		c.coreAdapter,
		requestDto,
		apicontract.LogoutOperation,
		"/core/v1/auth/logout",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	c.accessTokenCookieHandler.Delete(ctx)
	c.refreshTokenCookieHandler.Delete(ctx)
	writeClientResponse(ctx, response.Data)
}

func (c *AuthController) SendAuthCode(ctx *gin.Context, requestDto *apicontract.SendAuthCodeRequestDto) {
	response, exception := coreadapters.Call[apicontract.SendAuthCodeRequestDto, apicontract.SendAuthCodeResponseDto](ctx, c.coreAdapter, requestDto, apicontract.SendAuthCodeOperation, "/core/v1/auth/email/code")
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *AuthController) ValidateEmail(ctx *gin.Context, requestDto *apicontract.ValidateEmailRequestDto) {
	response, exception := coreadapters.CallSecurly[apicontract.ValidateEmailRequestDto, apicontract.ValidateEmailResponseDto](ctx, c.coreAdapter, requestDto, apicontract.ValidateEmailOperation, "/core/v1/auth/email/validate")
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *AuthController) ResetEmail(ctx *gin.Context, requestDto *apicontract.ResetEmailRequestDto) {
	response, exception := coreadapters.CallSecurly[apicontract.ResetEmailRequestDto, apicontract.ResetEmailResponseDto](ctx, c.coreAdapter, requestDto, apicontract.ResetEmailOperation, "/core/v1/auth/email/reset")
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *AuthController) ForgetPassword(ctx *gin.Context, requestDto *apicontract.ForgetPasswordRequestDto) {
	response, exception := coreadapters.Call[apicontract.ForgetPasswordRequestDto, apicontract.ForgetPasswordResponseDto](ctx, c.coreAdapter, requestDto, apicontract.ForgetPasswordOperation, "/core/v1/auth/password/forget")
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *AuthController) ResetMe(ctx *gin.Context, requestDto *apicontract.ResetMeRequestDto) {
	response, exception := coreadapters.CallSecurly[apicontract.ResetMeRequestDto, apicontract.ResetMeResponseDto](ctx, c.coreAdapter, requestDto, apicontract.ResetMeOperation, "/core/v1/auth/me/reset")
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *AuthController) DeleteMe(ctx *gin.Context, requestDto *apicontract.DeleteMeRequestDto) {
	response, exception := coreadapters.CallSecurly[apicontract.DeleteMeRequestDto, apicontract.DeleteMeResponseDto](ctx, c.coreAdapter, requestDto, apicontract.DeleteMeOperation, "/core/v1/auth/me/delete")
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}
