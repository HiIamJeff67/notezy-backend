package binders

import (
	"github.com/gin-gonic/gin"

	exceptions "github.com/HiIamJeff67/notezy-backend/contracts/types/exceptions"

	exceptionwriter "github.com/HiIamJeff67/notezy-backend/shared/util/exceptionwriter"

	apicontract "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api/user-settings"

	controllers "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/api/controllers"
)

type UserSettingBinderInterface interface {
	BindGetMySetting(controllerFunc controllers.Func[*apicontract.GetMySettingRequestDto]) gin.HandlerFunc
	BindUpdateMySetting(controllerFunc controllers.Func[*apicontract.UpdateMySettingRequestDto]) gin.HandlerFunc
}

type UserSettingBinder struct{}

func NewUserSettingBinder() UserSettingBinderInterface {
	return &UserSettingBinder{}
}

func (b *UserSettingBinder) BindGetMySetting(controllerFunc controllers.Func[*apicontract.GetMySettingRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &apicontract.GetMySettingRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		controllerFunc(ctx, requestDto)
	}
}

func (b *UserSettingBinder) BindUpdateMySetting(controllerFunc controllers.Func[*apicontract.UpdateMySettingRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &apicontract.UpdateMySettingRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")
		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("UserSetting").WithOrigin(err), ctx)
			return
		}

		controllerFunc(ctx, requestDto)
	}
}
