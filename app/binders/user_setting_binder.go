package binders

import (
	"github.com/gin-gonic/gin"

	contexts "github.com/HiIamJeff67/notezy-backend/app/contexts"
	dtos "github.com/HiIamJeff67/notezy-backend/app/dtos"
	types "github.com/HiIamJeff67/notezy-backend/shared/types"
)

type UserSettingBinderInterface interface {
	BindGetMySetting(controllerFunc types.ControllerFunc[*dtos.GetMySettingReqDto]) gin.HandlerFunc
}

type UserSettingBinder struct{}

func NewUserSettingBinder() UserSettingBinderInterface {
	return &UserSettingBinder{}
}

func (b *UserSettingBinder) BindGetMySetting(controllerFunc types.ControllerFunc[*dtos.GetMySettingReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.GetMySettingReqDto

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
