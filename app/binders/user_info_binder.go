package binders

import (
	"github.com/gin-gonic/gin"

	contexts "notezy-backend/app/contexts"
	dtos "notezy-backend/app/dtos"
	exceptions "notezy-backend/app/exceptions"
	constants "notezy-backend/shared/constants"
	types "notezy-backend/shared/types"
)

/* ============================== Interface & Instance ============================== */

type UserInfoBinderInterface interface {
	BindGetMyInfo(controllerFunc types.ControllerFunc[*dtos.GetMyInfoReqDto]) gin.HandlerFunc
	BindUpdateMyInfo(controllerFunc types.ControllerFunc[*dtos.UpdateMyInfoReqDto]) gin.HandlerFunc
}

type UserInfoBinder struct{}

func NewUserInfoBinder() UserInfoBinderInterface {
	return &UserInfoBinder{}
}

/* ============================== Binder ============================== */

func (b *UserInfoBinder) BindGetMyInfo(controllerFunc types.ControllerFunc[*dtos.GetMyInfoReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.GetMyInfoReqDto
		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		controllerFunc(ctx, &reqDto)
	}
}

func (b *UserInfoBinder) BindUpdateMyInfo(controllerFunc types.ControllerFunc[*dtos.UpdateMyInfoReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.UpdateMyInfoReqDto
		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId
		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.UserInfo.InvalidDto().WithError(err)
			exception.ResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}
