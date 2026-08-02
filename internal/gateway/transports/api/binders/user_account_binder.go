package binders

import (
	"github.com/gin-gonic/gin"

	useraccountsdto "github.com/HiIamJeff67/notezy-backend/contracts/gateway/v1/api/user-accounts"
	exceptions "github.com/HiIamJeff67/notezy-backend/internal/exceptions"
	responsewriter "github.com/HiIamJeff67/notezy-backend/internal/shared/responsewriter"
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
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("UserAccount").WithOrigin(err), ctx)
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
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("UserAccount").WithOrigin(err), ctx)
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
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("UserAccount").WithOrigin(err), ctx)
			return
		}

		controllerFunc(ctx, requestDto)
	}
}
