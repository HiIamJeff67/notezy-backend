package binders

import (
	"github.com/gin-gonic/gin"

	exceptions "github.com/HiIamJeff67/notezy-backend/contracts/types/exceptions"

	exceptionwriter "github.com/HiIamJeff67/notezy-backend/shared/util/exceptionwriter"

	apicontract "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api/users"

	controllers "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/api/controllers"
)

type UserBinderInterface interface {
	BindGetUserData(controllers.Func[*apicontract.GetUserDataRequestDto]) gin.HandlerFunc
	BindGetMe(controllers.Func[*apicontract.GetMeRequestDto]) gin.HandlerFunc
	BindUpdateMe(controllers.Func[*apicontract.UpdateMeRequestDto]) gin.HandlerFunc
}
type UserBinder struct{}

func NewUserBinder() UserBinderInterface { return &UserBinder{} }
func (b *UserBinder) BindGetUserData(controllerFunc controllers.Func[*apicontract.GetUserDataRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &apicontract.GetUserDataRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")
		controllerFunc(ctx, requestDto)
	}
}
func (b *UserBinder) BindGetMe(controllerFunc controllers.Func[*apicontract.GetMeRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &apicontract.GetMeRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")
		controllerFunc(ctx, requestDto)
	}
}
func (b *UserBinder) BindUpdateMe(controllerFunc controllers.Func[*apicontract.UpdateMeRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &apicontract.UpdateMeRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")
		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("User").WithOrigin(err), ctx)
			return
		}
		controllerFunc(ctx, requestDto)
	}
}
