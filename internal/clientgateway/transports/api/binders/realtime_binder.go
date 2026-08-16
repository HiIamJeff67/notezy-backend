package binders

import (
	"github.com/gin-gonic/gin"

	exceptions "github.com/HiIamJeff67/notegic-backend/contracts/types/exceptions"

	exceptionwriter "github.com/HiIamJeff67/notegic-backend/shared/util/exceptionwriter"

	apicontract "github.com/HiIamJeff67/notegic-backend/contracts/core/v1/api/realtime"

	controllers "github.com/HiIamJeff67/notegic-backend/internal/clientgateway/transports/api/controllers"
)

type RealtimeBinderInterface interface {
	BindCreateMyRealtimeConnectionTicket(controllerFunc controllers.Func[*apicontract.CreateMyRealtimeConnectionTicketRequestDto]) gin.HandlerFunc
	BindCreateMyBlockPackChannelTicket(controllerFunc controllers.Func[*apicontract.CreateMyBlockPackChannelTicketRequestDto]) gin.HandlerFunc
}

type RealtimeBinder struct{}

func NewRealtimeBinder() RealtimeBinderInterface {
	return &RealtimeBinder{}
}

func (b *RealtimeBinder) BindCreateMyRealtimeConnectionTicket(controllerFunc controllers.Func[*apicontract.CreateMyRealtimeConnectionTicketRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &apicontract.CreateMyRealtimeConnectionTicketRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		controllerFunc(ctx, requestDto)
	}
}

func (b *RealtimeBinder) BindCreateMyBlockPackChannelTicket(controllerFunc controllers.Func[*apicontract.CreateMyBlockPackChannelTicketRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &apicontract.CreateMyBlockPackChannelTicketRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("BlockPack").WithOrigin(err), ctx)
			return
		}

		controllerFunc(ctx, requestDto)
	}
}
