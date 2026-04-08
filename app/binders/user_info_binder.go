package binders

import (
	"github.com/gin-gonic/gin"

	contexts "notezy-backend/app/contexts"
	dtos "notezy-backend/app/dtos"
	exceptions "notezy-backend/app/exceptions"
	types "notezy-backend/shared/types"
)

type UserInfoBinderInterface interface {
	BindGetMyInfo(controllerFunc types.ControllerFunc[*dtos.GetMyInfoReqDto]) gin.HandlerFunc
	BindUpdateMyInfo(controllerFunc types.ControllerFunc[*dtos.UpdateMyInfoReqDto]) gin.HandlerFunc
}

type UserInfoBinder struct{}

func NewUserInfoBinder() UserInfoBinderInterface {
	return &UserInfoBinder{}
}

func (b *UserInfoBinder) BindGetMyInfo(controllerFunc types.ControllerFunc[*dtos.GetMyInfoReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.GetMyInfoReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		controllerFunc(ctx, &reqDto)
	}
}

func (b *UserInfoBinder) BindUpdateMyInfo(controllerFunc types.ControllerFunc[*dtos.UpdateMyInfoReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.UpdateMyInfoReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.UserInfo.InvalidDto().WithOrigin(err)
			exception.SafelyAbortAndResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}
