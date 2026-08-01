package binders

import (
	"github.com/gin-gonic/gin"

	authdto "github.com/HiIamJeff67/notezy-backend/contracts/api/v1/auth"
	exceptions "github.com/HiIamJeff67/notezy-backend/internal/exceptions"
	responsewriter "github.com/HiIamJeff67/notezy-backend/internal/gateway/responsewriter"
	apitransport "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/api"
)

type AuthBinderInterface interface {
	BindRegister(controllerFunc apitransport.ControllerFunc[*authdto.RegisterRequestDto]) gin.HandlerFunc
	BindRegisterViaGoogle(controllerFunc apitransport.ControllerFunc[*authdto.RegisterViaGoogleRequestDto]) gin.HandlerFunc
	BindLogin(controllerFunc apitransport.ControllerFunc[*authdto.LoginRequestDto]) gin.HandlerFunc
	BindLoginViaGoogle(controllerFunc apitransport.ControllerFunc[*authdto.LoginViaGoogleRequestDto]) gin.HandlerFunc
	BindLogout(controllerFunc apitransport.ControllerFunc[*authdto.LogoutRequestDto]) gin.HandlerFunc
	BindSendAuthCode(controllerFunc apitransport.ControllerFunc[*authdto.SendAuthCodeRequestDto]) gin.HandlerFunc
	BindValidateEmail(controllerFunc apitransport.ControllerFunc[*authdto.ValidateEmailRequestDto]) gin.HandlerFunc
	BindResetEmail(controllerFunc apitransport.ControllerFunc[*authdto.ResetEmailRequestDto]) gin.HandlerFunc
	BindForgetPassword(controllerFunc apitransport.ControllerFunc[*authdto.ForgetPasswordRequestDto]) gin.HandlerFunc
	BindResetMe(controllerFunc apitransport.ControllerFunc[*authdto.ResetMeRequestDto]) gin.HandlerFunc
	BindDeleteMe(controllerFunc apitransport.ControllerFunc[*authdto.DeleteMeRequestDto]) gin.HandlerFunc
}

type AuthBinder struct{}

func NewAuthBinder() AuthBinderInterface {
	return &AuthBinder{}
}

func (b *AuthBinder) BindRegister(controllerFunc apitransport.ControllerFunc[*authdto.RegisterRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &authdto.RegisterRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")
		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("Auth").WithOrigin(err), ctx)
			return
		}

		controllerFunc(ctx, requestDto)
	}
}

func (b *AuthBinder) BindRegisterViaGoogle(controllerFunc apitransport.ControllerFunc[*authdto.RegisterViaGoogleRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &authdto.RegisterViaGoogleRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")
		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("Auth").WithOrigin(err), ctx)
			return
		}

		controllerFunc(ctx, requestDto)
	}
}

func (b *AuthBinder) BindLogin(controllerFunc apitransport.ControllerFunc[*authdto.LoginRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &authdto.LoginRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")
		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("Auth").WithOrigin(err), ctx)
			return
		}

		controllerFunc(ctx, requestDto)
	}
}

func (b *AuthBinder) BindLoginViaGoogle(controllerFunc apitransport.ControllerFunc[*authdto.LoginViaGoogleRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &authdto.LoginViaGoogleRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")
		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("Auth").WithOrigin(err), ctx)
			return
		}

		controllerFunc(ctx, requestDto)
	}
}

func (b *AuthBinder) BindLogout(controllerFunc apitransport.ControllerFunc[*authdto.LogoutRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &authdto.LogoutRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		controllerFunc(ctx, requestDto)
	}
}

func (b *AuthBinder) BindSendAuthCode(controllerFunc apitransport.ControllerFunc[*authdto.SendAuthCodeRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &authdto.SendAuthCodeRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")
		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("Auth").WithOrigin(err), ctx)
			return
		}

		controllerFunc(ctx, requestDto)
	}
}

func (b *AuthBinder) BindValidateEmail(controllerFunc apitransport.ControllerFunc[*authdto.ValidateEmailRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &authdto.ValidateEmailRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")
		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("Auth").WithOrigin(err), ctx)
			return
		}

		controllerFunc(ctx, requestDto)
	}
}

func (b *AuthBinder) BindResetEmail(controllerFunc apitransport.ControllerFunc[*authdto.ResetEmailRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &authdto.ResetEmailRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")
		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("Auth").WithOrigin(err), ctx)
			return
		}

		controllerFunc(ctx, requestDto)
	}
}

func (b *AuthBinder) BindForgetPassword(controllerFunc apitransport.ControllerFunc[*authdto.ForgetPasswordRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &authdto.ForgetPasswordRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")
		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("Auth").WithOrigin(err), ctx)
			return
		}

		controllerFunc(ctx, requestDto)
	}
}

func (b *AuthBinder) BindResetMe(controllerFunc apitransport.ControllerFunc[*authdto.ResetMeRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &authdto.ResetMeRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")
		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("Auth").WithOrigin(err), ctx)
			return
		}

		controllerFunc(ctx, requestDto)
	}
}

func (b *AuthBinder) BindDeleteMe(controllerFunc apitransport.ControllerFunc[*authdto.DeleteMeRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &authdto.DeleteMeRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")
		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("Auth").WithOrigin(err), ctx)
			return
		}

		controllerFunc(ctx, requestDto)
	}
}
