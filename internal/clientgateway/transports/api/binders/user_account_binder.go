package binders

import (
	"github.com/gin-gonic/gin"

	exceptions "github.com/HiIamJeff67/notezy-backend/contracts/types/exceptions"

	exceptionwriter "github.com/HiIamJeff67/notezy-backend/shared/util/exceptionwriter"

	apicontract "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api/user-accounts"

	controllers "github.com/HiIamJeff67/notezy-backend/internal/clientgateway/transports/api/controllers"
)

type UserAccountBinderInterface interface {
	BindGetMyAccount(controllerFunc controllers.Func[*apicontract.GetMyAccountRequestDto]) gin.HandlerFunc
	BindUpdateMyAccount(controllerFunc controllers.Func[*apicontract.UpdateMyAccountRequestDto]) gin.HandlerFunc
	BindBindGoogleAccount(controllerFunc controllers.Func[*apicontract.BindGoogleAccountRequestDto]) gin.HandlerFunc
	BindUnbindGoogleAccount(controllerFunc controllers.Func[*apicontract.UnbindGoogleAccountRequestDto]) gin.HandlerFunc
}

type UserAccountBinder struct{}

func NewUserAccountBinder() UserAccountBinderInterface {
	return &UserAccountBinder{}
}

func (b *UserAccountBinder) BindGetMyAccount(controllerFunc controllers.Func[*apicontract.GetMyAccountRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &apicontract.GetMyAccountRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		controllerFunc(ctx, requestDto)
	}
}

func (b *UserAccountBinder) BindUpdateMyAccount(controllerFunc controllers.Func[*apicontract.UpdateMyAccountRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &apicontract.UpdateMyAccountRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")
		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("UserAccount").WithOrigin(err), ctx)
			return
		}

		controllerFunc(ctx, requestDto)
	}
}

func (b *UserAccountBinder) BindBindGoogleAccount(controllerFunc controllers.Func[*apicontract.BindGoogleAccountRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &apicontract.BindGoogleAccountRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")
		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("UserAccount").WithOrigin(err), ctx)
			return
		}

		controllerFunc(ctx, requestDto)
	}
}

func (b *UserAccountBinder) BindUnbindGoogleAccount(controllerFunc controllers.Func[*apicontract.UnbindGoogleAccountRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &apicontract.UnbindGoogleAccountRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")
		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("UserAccount").WithOrigin(err), ctx)
			return
		}

		controllerFunc(ctx, requestDto)
	}
}
