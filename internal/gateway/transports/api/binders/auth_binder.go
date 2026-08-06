package binders

import (
	"github.com/gin-gonic/gin"

	exceptions "github.com/HiIamJeff67/notezy-backend/shared/exceptions"

	exceptionwriter "github.com/HiIamJeff67/notezy-backend/shared/util/exceptionwriter"

	authdto "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api/auth"

	controllers "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/api/controllers"
)

type AuthBinderInterface interface {
	BindRegister(controllerFunc controllers.Func[*authdto.RegisterRequestDto]) gin.HandlerFunc
	BindRegisterViaGoogle(controllerFunc controllers.Func[*authdto.RegisterViaGoogleRequestDto]) gin.HandlerFunc
	BindLogin(controllerFunc controllers.Func[*authdto.LoginRequestDto]) gin.HandlerFunc
	BindLoginViaGoogle(controllerFunc controllers.Func[*authdto.LoginViaGoogleRequestDto]) gin.HandlerFunc
	BindLogout(controllerFunc controllers.Func[*authdto.LogoutRequestDto]) gin.HandlerFunc
	BindSendAuthCode(controllerFunc controllers.Func[*authdto.SendAuthCodeRequestDto]) gin.HandlerFunc
	BindValidateEmail(controllerFunc controllers.Func[*authdto.ValidateEmailRequestDto]) gin.HandlerFunc
	BindResetEmail(controllerFunc controllers.Func[*authdto.ResetEmailRequestDto]) gin.HandlerFunc
	BindForgetPassword(controllerFunc controllers.Func[*authdto.ForgetPasswordRequestDto]) gin.HandlerFunc
	BindResetMe(controllerFunc controllers.Func[*authdto.ResetMeRequestDto]) gin.HandlerFunc
	BindDeleteMe(controllerFunc controllers.Func[*authdto.DeleteMeRequestDto]) gin.HandlerFunc
}

type AuthBinder struct{}

func NewAuthBinder() AuthBinderInterface {
	return &AuthBinder{}
}

func (b *AuthBinder) BindRegister(controllerFunc controllers.Func[*authdto.RegisterRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &authdto.RegisterRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")
		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("Auth").WithOrigin(err), ctx)
			return
		}

		controllerFunc(ctx, requestDto)
	}
}

func (b *AuthBinder) BindRegisterViaGoogle(controllerFunc controllers.Func[*authdto.RegisterViaGoogleRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &authdto.RegisterViaGoogleRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")
		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("Auth").WithOrigin(err), ctx)
			return
		}

		controllerFunc(ctx, requestDto)
	}
}

func (b *AuthBinder) BindLogin(controllerFunc controllers.Func[*authdto.LoginRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &authdto.LoginRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")
		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("Auth").WithOrigin(err), ctx)
			return
		}

		controllerFunc(ctx, requestDto)
	}
}

func (b *AuthBinder) BindLoginViaGoogle(controllerFunc controllers.Func[*authdto.LoginViaGoogleRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &authdto.LoginViaGoogleRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")
		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("Auth").WithOrigin(err), ctx)
			return
		}

		controllerFunc(ctx, requestDto)
	}
}

func (b *AuthBinder) BindLogout(controllerFunc controllers.Func[*authdto.LogoutRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &authdto.LogoutRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		controllerFunc(ctx, requestDto)
	}
}

func (b *AuthBinder) BindSendAuthCode(controllerFunc controllers.Func[*authdto.SendAuthCodeRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &authdto.SendAuthCodeRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")
		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("Auth").WithOrigin(err), ctx)
			return
		}

		controllerFunc(ctx, requestDto)
	}
}

func (b *AuthBinder) BindValidateEmail(controllerFunc controllers.Func[*authdto.ValidateEmailRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &authdto.ValidateEmailRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")
		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("Auth").WithOrigin(err), ctx)
			return
		}

		controllerFunc(ctx, requestDto)
	}
}

func (b *AuthBinder) BindResetEmail(controllerFunc controllers.Func[*authdto.ResetEmailRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &authdto.ResetEmailRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")
		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("Auth").WithOrigin(err), ctx)
			return
		}

		controllerFunc(ctx, requestDto)
	}
}

func (b *AuthBinder) BindForgetPassword(controllerFunc controllers.Func[*authdto.ForgetPasswordRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &authdto.ForgetPasswordRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")
		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("Auth").WithOrigin(err), ctx)
			return
		}

		controllerFunc(ctx, requestDto)
	}
}

func (b *AuthBinder) BindResetMe(controllerFunc controllers.Func[*authdto.ResetMeRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &authdto.ResetMeRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")
		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("Auth").WithOrigin(err), ctx)
			return
		}

		controllerFunc(ctx, requestDto)
	}
}

func (b *AuthBinder) BindDeleteMe(controllerFunc controllers.Func[*authdto.DeleteMeRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &authdto.DeleteMeRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")
		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("Auth").WithOrigin(err), ctx)
			return
		}

		controllerFunc(ctx, requestDto)
	}
}
