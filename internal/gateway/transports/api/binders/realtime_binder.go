package binders

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	realtimedto "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api/realtime"
	exceptions "github.com/HiIamJeff67/notezy-backend/internal/exceptions"
	controllers "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/api/controllers"
	exceptionwriter "github.com/HiIamJeff67/notezy-backend/shared/exceptionwriter"
)

type RealtimeBinderInterface interface {
	BindGetMyBlockPackRealtimeParticipants(controllerFunc controllers.Func[*realtimedto.GetMyBlockPackRealtimeParticipantsRequestDto]) gin.HandlerFunc
	BindCreateMyRealtimeConnectionTicket(controllerFunc controllers.Func[*realtimedto.CreateMyRealtimeConnectionTicketRequestDto]) gin.HandlerFunc
	BindCreateMyBlockPackChannelTicket(controllerFunc controllers.Func[*realtimedto.CreateMyBlockPackChannelTicketRequestDto]) gin.HandlerFunc
}

type RealtimeBinder struct{}

func NewRealtimeBinder() RealtimeBinderInterface {
	return &RealtimeBinder{}
}

func (b *RealtimeBinder) BindGetMyBlockPackRealtimeParticipants(controllerFunc controllers.Func[*realtimedto.GetMyBlockPackRealtimeParticipantsRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &realtimedto.GetMyBlockPackRealtimeParticipantsRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		blockPackId, err := uuid.Parse(ctx.Param("blockPackId"))
		if err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("BlockPack").WithOrigin(err), ctx)
			return
		}
		requestDto.Param.BlockPackId = blockPackId

		controllerFunc(ctx, requestDto)
	}
}

func (b *RealtimeBinder) BindCreateMyRealtimeConnectionTicket(controllerFunc controllers.Func[*realtimedto.CreateMyRealtimeConnectionTicketRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &realtimedto.CreateMyRealtimeConnectionTicketRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		controllerFunc(ctx, requestDto)
	}
}

func (b *RealtimeBinder) BindCreateMyBlockPackChannelTicket(controllerFunc controllers.Func[*realtimedto.CreateMyBlockPackChannelTicketRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &realtimedto.CreateMyBlockPackChannelTicketRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("BlockPack").WithOrigin(err), ctx)
			return
		}

		controllerFunc(ctx, requestDto)
	}
}
