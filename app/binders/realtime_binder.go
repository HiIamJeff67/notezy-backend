package binders

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	contexts "github.com/HiIamJeff67/notezy-backend/app/contexts"
	dtos "github.com/HiIamJeff67/notezy-backend/app/dtos"
	exceptions "github.com/HiIamJeff67/notezy-backend/app/exceptions"
	types "github.com/HiIamJeff67/notezy-backend/shared/types"
)

type RealtimeBinderInterface interface {
	BindGetMyBlockPackRealtimeParticipants(controllerFunc types.ControllerFunc[*dtos.GetMyBlockPackRealtimeParticipantsReqDto]) gin.HandlerFunc
	BindCreateMyRealtimeConnectionTicket(controllerFunc types.ControllerFunc[*dtos.CreateMyRealtimeConnectionTicketReqDto]) gin.HandlerFunc
	BindCreateMyBlockPackChannelTicket(controllerFunc types.ControllerFunc[*dtos.CreateMyBlockPackChannelTicketReqDto]) gin.HandlerFunc
}

type RealtimeBinder struct{}

func NewRealtimeBinder() RealtimeBinderInterface {
	return &RealtimeBinder{}
}

func (b *RealtimeBinder) BindGetMyBlockPackRealtimeParticipants(
	controllerFunc types.ControllerFunc[*dtos.GetMyBlockPackRealtimeParticipantsReqDto],
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.GetMyBlockPackRealtimeParticipantsReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		blockPackId, err := uuid.Parse(ctx.Param("blockPackId"))
		if err != nil {
			exceptions.BlockPack.InvalidInput().WithOrigin(err).Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.Param.BlockPackId = blockPackId

		controllerFunc(ctx, &reqDto)
	}
}

func (b *RealtimeBinder) BindCreateMyRealtimeConnectionTicket(
	controllerFunc types.ControllerFunc[*dtos.CreateMyRealtimeConnectionTicketReqDto],
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.CreateMyRealtimeConnectionTicketReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userPublicId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_PublicId)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserPublicId = *userPublicId

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

		userPublicId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_PublicId)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserPublicId = *userPublicId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exceptions.BlockPack.InvalidDto().WithOrigin(err).Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}
