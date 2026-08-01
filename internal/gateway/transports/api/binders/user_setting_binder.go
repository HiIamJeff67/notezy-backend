package binders

import (
	"github.com/gin-gonic/gin"

	usersettingsdto "github.com/HiIamJeff67/notezy-backend/contracts/api/v1/user-settings"
	exceptions "github.com/HiIamJeff67/notezy-backend/internal/exceptions"
	responsewriter "github.com/HiIamJeff67/notezy-backend/internal/gateway/responsewriter"
	apitransport "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/api"
)

type UserSettingBinderInterface interface {
	BindGetMySetting(controllerFunc apitransport.ControllerFunc[*usersettingsdto.GetMySettingRequestDto]) gin.HandlerFunc
	BindUpdateMySetting(controllerFunc apitransport.ControllerFunc[*usersettingsdto.UpdateMySettingRequestDto]) gin.HandlerFunc
}

type UserSettingBinder struct{}

func NewUserSettingBinder() UserSettingBinderInterface {
	return &UserSettingBinder{}
}

func (b *UserSettingBinder) BindGetMySetting(controllerFunc apitransport.ControllerFunc[*usersettingsdto.GetMySettingRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &usersettingsdto.GetMySettingRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		controllerFunc(ctx, requestDto)
	}
}

func (b *UserSettingBinder) BindUpdateMySetting(controllerFunc apitransport.ControllerFunc[*usersettingsdto.UpdateMySettingRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &usersettingsdto.UpdateMySettingRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")
		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("UserSetting").WithOrigin(err), ctx)
			return
		}

		controllerFunc(ctx, requestDto)
	}
}
