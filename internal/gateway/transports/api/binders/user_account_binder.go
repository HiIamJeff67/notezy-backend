package binders

import (
	"github.com/gin-gonic/gin"

	useraccountsdto "github.com/HiIamJeff67/notezy-backend/contracts/api/v1/user-accounts"
	exceptions "github.com/HiIamJeff67/notezy-backend/internal/exceptions"
	responsewriter "github.com/HiIamJeff67/notezy-backend/internal/gateway/responsewriter"
	apitransport "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/api"
)

type UserAccountBinderInterface interface {
	BindGetMyAccount(controllerFunc apitransport.ControllerFunc[*useraccountsdto.GetMyAccountRequestDto]) gin.HandlerFunc
	BindUpdateMyAccount(controllerFunc apitransport.ControllerFunc[*useraccountsdto.UpdateMyAccountRequestDto]) gin.HandlerFunc
	BindBindGoogleAccount(controllerFunc apitransport.ControllerFunc[*useraccountsdto.BindGoogleAccountRequestDto]) gin.HandlerFunc
	BindUnbindGoogleAccount(controllerFunc apitransport.ControllerFunc[*useraccountsdto.UnbindGoogleAccountRequestDto]) gin.HandlerFunc
}

type UserAccountBinder struct{}

func NewUserAccountBinder() UserAccountBinderInterface {
	return &UserAccountBinder{}
}

func (b *UserAccountBinder) BindGetMyAccount(controllerFunc apitransport.ControllerFunc[*useraccountsdto.GetMyAccountRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &useraccountsdto.GetMyAccountRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		controllerFunc(ctx, requestDto)
	}
}

func (b *UserAccountBinder) BindUpdateMyAccount(controllerFunc apitransport.ControllerFunc[*useraccountsdto.UpdateMyAccountRequestDto]) gin.HandlerFunc {
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

func (b *UserAccountBinder) BindBindGoogleAccount(controllerFunc apitransport.ControllerFunc[*useraccountsdto.BindGoogleAccountRequestDto]) gin.HandlerFunc {
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

func (b *UserAccountBinder) BindUnbindGoogleAccount(controllerFunc apitransport.ControllerFunc[*useraccountsdto.UnbindGoogleAccountRequestDto]) gin.HandlerFunc {
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
