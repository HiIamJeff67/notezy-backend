package binders

import (
	"github.com/gin-gonic/gin"

	contexts "github.com/HiIamJeff67/notezy-backend/app/contexts"
	dtos "github.com/HiIamJeff67/notezy-backend/app/dtos"
	exceptions "github.com/HiIamJeff67/notezy-backend/app/exceptions"
	types "github.com/HiIamJeff67/notezy-backend/shared/types"
)

type RealtimeBinderInterface interface {
	BindCreateMyRealtimeConnectionTicket(controllerFunc types.ControllerFunc[*dtos.CreateMyRealtimeConnectionTicketReqDto]) gin.HandlerFunc
	BindCreateMyBlockPackChannelTicket(controllerFunc types.ControllerFunc[*dtos.CreateMyBlockPackChannelTicketReqDto]) gin.HandlerFunc
}

type RealtimeBinder struct{}

func NewRealtimeBinder() RealtimeBinderInterface {
	return &RealtimeBinder{}
}

func (b *RealtimeBinder) BindCreateMyRealtimeConnectionTicket(
	controllerFunc types.ControllerFunc[*dtos.CreateMyRealtimeConnectionTicketReqDto],
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.CreateMyRealtimeConnectionTicketReqDto

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

func (b *RealtimeBinder) BindCreateMyBlockPackChannelTicket(
	controllerFunc types.ControllerFunc[*dtos.CreateMyBlockPackChannelTicketReqDto],
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.CreateMyBlockPackChannelTicketReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exceptions.BlockPack.InvalidDto().WithOrigin(err).Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}
