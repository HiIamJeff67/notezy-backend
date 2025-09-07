package binders

import (
	contexts "notezy-backend/app/contexts"
	"notezy-backend/app/dtos"
	constants "notezy-backend/shared/constants"
	types "notezy-backend/shared/types"

	"github.com/gin-gonic/gin"
)

/* ============================== Interface & Instance ============================== */

type UserSettingBinderInterface interface {
	BindGetMySetting(controllerFunc types.ControllerFunc[*dtos.GetMySettingReqDto]) gin.HandlerFunc
}

type UserSettingBinder struct{}

func NewUserSettingBinder() UserSettingBinderInterface {
	return &UserSettingBinder{}
}

/* ============================== Binder ============================== */

func (b *UserSettingBinder) BindGetMySetting(controllerFunc types.ControllerFunc[*dtos.GetMySettingReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.GetMySettingReqDto
		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		controllerFunc(ctx, &reqDto)
	}
}
