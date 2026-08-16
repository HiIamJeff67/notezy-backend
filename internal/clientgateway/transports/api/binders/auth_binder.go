package binders

import (
	"github.com/gin-gonic/gin"

	exceptions "github.com/HiIamJeff67/notegic-backend/contracts/types/exceptions"

	exceptionwriter "github.com/HiIamJeff67/notegic-backend/shared/util/exceptionwriter"

	apicontract "github.com/HiIamJeff67/notegic-backend/contracts/core/v1/api/auth"

	controllers "github.com/HiIamJeff67/notegic-backend/internal/clientgateway/transports/api/controllers"
)

type AuthBinderInterface interface {
	BindRegister(controllerFunc controllers.Func[*apicontract.RegisterRequestDto]) gin.HandlerFunc
	BindRegisterViaGoogle(controllerFunc controllers.Func[*apicontract.RegisterViaGoogleRequestDto]) gin.HandlerFunc
	BindLogin(controllerFunc controllers.Func[*apicontract.LoginRequestDto]) gin.HandlerFunc
	BindLoginViaGoogle(controllerFunc controllers.Func[*apicontract.LoginViaGoogleRequestDto]) gin.HandlerFunc
	BindLogout(controllerFunc controllers.Func[*apicontract.LogoutRequestDto]) gin.HandlerFunc
	BindSendAuthCode(controllerFunc controllers.Func[*apicontract.SendAuthCodeRequestDto]) gin.HandlerFunc
	BindValidateEmail(controllerFunc controllers.Func[*apicontract.ValidateEmailRequestDto]) gin.HandlerFunc
	BindResetEmail(controllerFunc controllers.Func[*apicontract.ResetEmailRequestDto]) gin.HandlerFunc
	BindForgetPassword(controllerFunc controllers.Func[*apicontract.ForgetPasswordRequestDto]) gin.HandlerFunc
	BindResetMe(controllerFunc controllers.Func[*apicontract.ResetMeRequestDto]) gin.HandlerFunc
	BindDeleteMe(controllerFunc controllers.Func[*apicontract.DeleteMeRequestDto]) gin.HandlerFunc
}

type AuthBinder struct{}

func NewAuthBinder() AuthBinderInterface {
	return &AuthBinder{}
}

func (b *AuthBinder) BindRegister(controllerFunc controllers.Func[*apicontract.RegisterRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &apicontract.RegisterRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")
		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("Auth").WithOrigin(err), ctx)
			return
		}

		controllerFunc(ctx, requestDto)
	}
}

func (b *AuthBinder) BindRegisterViaGoogle(controllerFunc controllers.Func[*apicontract.RegisterViaGoogleRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &apicontract.RegisterViaGoogleRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")
		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("Auth").WithOrigin(err), ctx)
			return
		}

		controllerFunc(ctx, requestDto)
	}
}

func (b *AuthBinder) BindLogin(controllerFunc controllers.Func[*apicontract.LoginRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &apicontract.LoginRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")
		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("Auth").WithOrigin(err), ctx)
			return
		}

		controllerFunc(ctx, requestDto)
	}
}

func (b *AuthBinder) BindLoginViaGoogle(controllerFunc controllers.Func[*apicontract.LoginViaGoogleRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &apicontract.LoginViaGoogleRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")
		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("Auth").WithOrigin(err), ctx)
			return
		}

		controllerFunc(ctx, requestDto)
	}
}

func (b *AuthBinder) BindLogout(controllerFunc controllers.Func[*apicontract.LogoutRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &apicontract.LogoutRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		controllerFunc(ctx, requestDto)
	}
}

func (b *AuthBinder) BindSendAuthCode(controllerFunc controllers.Func[*apicontract.SendAuthCodeRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &apicontract.SendAuthCodeRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")
		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("Auth").WithOrigin(err), ctx)
			return
		}

		controllerFunc(ctx, requestDto)
	}
}

func (b *AuthBinder) BindValidateEmail(controllerFunc controllers.Func[*apicontract.ValidateEmailRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &apicontract.ValidateEmailRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")
		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("Auth").WithOrigin(err), ctx)
			return
		}

		controllerFunc(ctx, requestDto)
	}
}

func (b *AuthBinder) BindResetEmail(controllerFunc controllers.Func[*apicontract.ResetEmailRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &apicontract.ResetEmailRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")
		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("Auth").WithOrigin(err), ctx)
			return
		}

		controllerFunc(ctx, requestDto)
	}
}

func (b *AuthBinder) BindForgetPassword(controllerFunc controllers.Func[*apicontract.ForgetPasswordRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &apicontract.ForgetPasswordRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")
		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("Auth").WithOrigin(err), ctx)
			return
		}

		controllerFunc(ctx, requestDto)
	}
}

func (b *AuthBinder) BindResetMe(controllerFunc controllers.Func[*apicontract.ResetMeRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &apicontract.ResetMeRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")
		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("Auth").WithOrigin(err), ctx)
			return
		}

		controllerFunc(ctx, requestDto)
	}
}

func (b *AuthBinder) BindDeleteMe(controllerFunc controllers.Func[*apicontract.DeleteMeRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &apicontract.DeleteMeRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")
		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("Auth").WithOrigin(err), ctx)
			return
		}

		controllerFunc(ctx, requestDto)
	}
}
