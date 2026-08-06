package binders

import (
	"github.com/gin-gonic/gin"

	exceptions "github.com/HiIamJeff67/notezy-backend/shared/exceptions"

	exceptionwriter "github.com/HiIamJeff67/notezy-backend/shared/util/exceptionwriter"

	useraccountsdto "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api/user-accounts"

	controllers "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/api/controllers"
)

type UserAccountBinderInterface interface {
	BindGetMyAccount(controllerFunc controllers.Func[*useraccountsdto.GetMyAccountRequestDto]) gin.HandlerFunc
	BindUpdateMyAccount(controllerFunc controllers.Func[*useraccountsdto.UpdateMyAccountRequestDto]) gin.HandlerFunc
	BindBindGoogleAccount(controllerFunc controllers.Func[*useraccountsdto.BindGoogleAccountRequestDto]) gin.HandlerFunc
	BindUnbindGoogleAccount(controllerFunc controllers.Func[*useraccountsdto.UnbindGoogleAccountRequestDto]) gin.HandlerFunc
}

type UserAccountBinder struct{}

func NewUserAccountBinder() UserAccountBinderInterface {
	return &UserAccountBinder{}
}

func (b *UserAccountBinder) BindGetMyAccount(controllerFunc controllers.Func[*useraccountsdto.GetMyAccountRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &useraccountsdto.GetMyAccountRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		controllerFunc(ctx, requestDto)
	}
}

func (b *UserAccountBinder) BindUpdateMyAccount(controllerFunc controllers.Func[*useraccountsdto.UpdateMyAccountRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &useraccountsdto.UpdateMyAccountRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")
		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("UserAccount").WithOrigin(err), ctx)
			return
		}

		controllerFunc(ctx, requestDto)
	}
}

func (b *UserAccountBinder) BindBindGoogleAccount(controllerFunc controllers.Func[*useraccountsdto.BindGoogleAccountRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &useraccountsdto.BindGoogleAccountRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")
		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("UserAccount").WithOrigin(err), ctx)
			return
		}

		controllerFunc(ctx, requestDto)
	}
}

func (b *UserAccountBinder) BindUnbindGoogleAccount(controllerFunc controllers.Func[*useraccountsdto.UnbindGoogleAccountRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &useraccountsdto.UnbindGoogleAccountRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")
		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("UserAccount").WithOrigin(err), ctx)
			return
		}

		controllerFunc(ctx, requestDto)
	}
}
