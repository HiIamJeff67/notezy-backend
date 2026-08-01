package binders

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	realtimedto "github.com/HiIamJeff67/notezy-backend/contracts/api/v1/realtime"
	exceptions "github.com/HiIamJeff67/notezy-backend/internal/exceptions"
	responsewriter "github.com/HiIamJeff67/notezy-backend/internal/gateway/responsewriter"
	apitransport "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/api"
)

type RealtimeBinderInterface interface {
	BindGetMyBlockPackRealtimeParticipants(controllerFunc apitransport.ControllerFunc[*realtimedto.GetMyBlockPackRealtimeParticipantsRequestDto]) gin.HandlerFunc
	BindCreateMyRealtimeConnectionTicket(controllerFunc apitransport.ControllerFunc[*realtimedto.CreateMyRealtimeConnectionTicketRequestDto]) gin.HandlerFunc
	BindCreateMyBlockPackChannelTicket(controllerFunc apitransport.ControllerFunc[*realtimedto.CreateMyBlockPackChannelTicketRequestDto]) gin.HandlerFunc
}

type RealtimeBinder struct{}

func NewRealtimeBinder() RealtimeBinderInterface {
	return &RealtimeBinder{}
}

func (b *RealtimeBinder) BindGetMyBlockPackRealtimeParticipants(controllerFunc apitransport.ControllerFunc[*realtimedto.GetMyBlockPackRealtimeParticipantsRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &realtimedto.GetMyBlockPackRealtimeParticipantsRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		blockPackId, err := uuid.Parse(ctx.Param("blockPackId"))
		if err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("BlockPack").WithOrigin(err), ctx)
			return
		}
		requestDto.Param.BlockPackId = blockPackId

		controllerFunc(ctx, requestDto)
	}
}

func (b *RealtimeBinder) BindCreateMyRealtimeConnectionTicket(controllerFunc apitransport.ControllerFunc[*realtimedto.CreateMyRealtimeConnectionTicketRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &realtimedto.CreateMyRealtimeConnectionTicketRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		controllerFunc(ctx, requestDto)
	}
}

func (b *RealtimeBinder) BindCreateMyBlockPackChannelTicket(controllerFunc apitransport.ControllerFunc[*realtimedto.CreateMyBlockPackChannelTicketRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &realtimedto.CreateMyBlockPackChannelTicketRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("BlockPack").WithOrigin(err), ctx)
			return
		}

		controllerFunc(ctx, requestDto)
	}
}
