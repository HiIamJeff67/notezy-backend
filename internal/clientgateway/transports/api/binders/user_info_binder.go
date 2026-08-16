package binders

import (
	"github.com/gin-gonic/gin"

	exceptions "github.com/HiIamJeff67/notegic-backend/contracts/types/exceptions"

	exceptionwriter "github.com/HiIamJeff67/notegic-backend/shared/util/exceptionwriter"

	apicontract "github.com/HiIamJeff67/notegic-backend/contracts/core/v1/api/user-infos"

	controllers "github.com/HiIamJeff67/notegic-backend/internal/clientgateway/transports/api/controllers"
)

type UserInfoBinderInterface interface {
	BindGetMyInfo(controllerFunc controllers.Func[*apicontract.GetMyInfoRequestDto]) gin.HandlerFunc
	BindUpdateMyInfo(controllerFunc controllers.Func[*apicontract.UpdateMyInfoRequestDto]) gin.HandlerFunc
}

type UserInfoBinder struct{}

func NewUserInfoBinder() UserInfoBinderInterface {
	return &UserInfoBinder{}
}

func (b *UserInfoBinder) BindGetMyInfo(controllerFunc controllers.Func[*apicontract.GetMyInfoRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &apicontract.GetMyInfoRequestDto{}
		request.Header.UserAgent = ctx.GetHeader("User-Agent")
		controllerFunc(ctx, request)
	}
}

func (b *UserInfoBinder) BindUpdateMyInfo(controllerFunc controllers.Func[*apicontract.UpdateMyInfoRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &apicontract.UpdateMyInfoRequestDto{}
		request.Header.UserAgent = ctx.GetHeader("User-Agent")
		if err := ctx.ShouldBindJSON(&request.Body); err != nil {
			exception := exceptions.InvalidDto("UserInfo").WithOrigin(err)
			exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
			return
		}

		controllerFunc(ctx, request)
	}
}
