package binders

import (
	usersdto "github.com/HiIamJeff67/notezy-backend/contracts/api/v1/users"
	exceptions "github.com/HiIamJeff67/notezy-backend/internal/exceptions"
	responsewriter "github.com/HiIamJeff67/notezy-backend/internal/gateway/responsewriter"
	apitransport "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/api"
	"github.com/gin-gonic/gin"
)

type UserBinderInterface interface {
	BindGetUserData(apitransport.ControllerFunc[*usersdto.GetUserDataRequestDto]) gin.HandlerFunc
	BindGetMe(apitransport.ControllerFunc[*usersdto.GetMeRequestDto]) gin.HandlerFunc
	BindUpdateMe(apitransport.ControllerFunc[*usersdto.UpdateMeRequestDto]) gin.HandlerFunc
}
type UserBinder struct{}

func NewUserBinder() UserBinderInterface { return &UserBinder{} }
func (b *UserBinder) BindGetUserData(controllerFunc apitransport.ControllerFunc[*usersdto.GetUserDataRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &usersdto.GetUserDataRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")
		controllerFunc(ctx, requestDto)
	}
}
func (b *UserBinder) BindGetMe(controllerFunc apitransport.ControllerFunc[*usersdto.GetMeRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &usersdto.GetMeRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")
		controllerFunc(ctx, requestDto)
	}
}
func (b *UserBinder) BindUpdateMe(controllerFunc apitransport.ControllerFunc[*usersdto.UpdateMeRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &usersdto.UpdateMeRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")
		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("User").WithOrigin(err), ctx)
			return
		}
		controllerFunc(ctx, requestDto)
	}
}
