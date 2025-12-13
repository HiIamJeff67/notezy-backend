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

type UserAccountBinderInterface interface {
	BindGetMyAccount(controllerFunc types.ControllerFunc[*dtos.GetMyAccountReqDto]) gin.HandlerFunc
	BindUpdateMyAccount(controllerFunc types.ControllerFunc[*dtos.UpdateMyAccountReqDto]) gin.HandlerFunc
}

type UserAccountBinder struct{}

func NewUserAccountBinder() UserAccountBinderInterface {
	return &UserAccountBinder{}
}

/* ============================== Implementations ============================== */

func (b *UserAccountBinder) BindGetMyAccount(controllerFunc types.ControllerFunc[*dtos.GetMyAccountReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.GetMyAccountReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		controllerFunc(ctx, &reqDto)
	}
}

func (b *UserAccountBinder) BindUpdateMyAccount(controllerFunc types.ControllerFunc[*dtos.UpdateMyAccountReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.UpdateMyAccountReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.UserAccount.InvalidDto().WithError(err)
			exception.ResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}
