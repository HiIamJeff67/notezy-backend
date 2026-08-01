package binders

import (
	"github.com/gin-gonic/gin"

	userinfosdto "github.com/HiIamJeff67/notezy-backend/contracts/api/v1/user-infos"
	exceptions "github.com/HiIamJeff67/notezy-backend/internal/exceptions"
	responsewriter "github.com/HiIamJeff67/notezy-backend/internal/gateway/responsewriter"
	apitransport "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/api"
)

type UserInfoBinderInterface interface {
	BindGetMyInfo(controllerFunc apitransport.ControllerFunc[*userinfosdto.GetMyInfoRequestDto]) gin.HandlerFunc
	BindUpdateMyInfo(controllerFunc apitransport.ControllerFunc[*userinfosdto.UpdateMyInfoRequestDto]) gin.HandlerFunc
}

type UserInfoBinder struct{}

func NewUserInfoBinder() UserInfoBinderInterface {
	return &UserInfoBinder{}
}

func (b *UserInfoBinder) BindGetMyInfo(controllerFunc apitransport.ControllerFunc[*userinfosdto.GetMyInfoRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &userinfosdto.GetMyInfoRequestDto{}
		request.Header.UserAgent = ctx.GetHeader("User-Agent")
		controllerFunc(ctx, request)
	}
}

func (b *UserInfoBinder) BindUpdateMyInfo(controllerFunc apitransport.ControllerFunc[*userinfosdto.UpdateMyInfoRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &userinfosdto.UpdateMyInfoRequestDto{}
		request.Header.UserAgent = ctx.GetHeader("User-Agent")
		if err := ctx.ShouldBindJSON(&request.Body); err != nil {
			exception := exceptions.InvalidDto("UserInfo").WithOrigin(err)
			responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
			return
		}

		controllerFunc(ctx, request)
	}
}
