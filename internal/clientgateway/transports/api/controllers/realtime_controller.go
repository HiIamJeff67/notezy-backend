package controllers

import (
	"github.com/gin-gonic/gin"

	exceptionwriter "github.com/HiIamJeff67/notegic-backend/shared/util/exceptionwriter"

	apicontract "github.com/HiIamJeff67/notegic-backend/contracts/core/v1/api/realtime"

	coreadapters "github.com/HiIamJeff67/notegic-backend/internal/clientgateway/transports/core/adapters"
)

type RealtimeControllerInterface interface {
	CreateMyRealtimeConnectionTicket(ctx *gin.Context, requestDto *apicontract.CreateMyRealtimeConnectionTicketRequestDto)
	CreateMyBlockPackChannelTicket(ctx *gin.Context, requestDto *apicontract.CreateMyBlockPackChannelTicketRequestDto)
}

type RealtimeController struct {
	coreAdapter *coreadapters.CoreAdapter
}

func NewRealtimeController(
	coreAdapter *coreadapters.CoreAdapter,
) RealtimeControllerInterface {
	return &RealtimeController{
		coreAdapter: coreAdapter,
	}
}

func (c *RealtimeController) CreateMyRealtimeConnectionTicket(
	ctx *gin.Context,
	requestDto *apicontract.CreateMyRealtimeConnectionTicketRequestDto,
) {
	response, exception := coreadapters.CallSecurly[
		apicontract.CreateMyRealtimeConnectionTicketRequestDto,
		apicontract.CreateMyRealtimeConnectionTicketResponseDto,
	](
		ctx,
		c.coreAdapter,
		requestDto,
		apicontract.CreateMyRealtimeConnectionTicketOperation,
		"/core/v1/realtime/connection-ticket/create",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeCreatedClientResponse(ctx, response.Data)
}

func (c *RealtimeController) CreateMyBlockPackChannelTicket(
	ctx *gin.Context,
	requestDto *apicontract.CreateMyBlockPackChannelTicketRequestDto,
) {
	response, exception := coreadapters.CallSecurly[
		apicontract.CreateMyBlockPackChannelTicketRequestDto,
		apicontract.CreateMyBlockPackChannelTicketResponseDto,
	](
		ctx,
		c.coreAdapter,
		requestDto,
		apicontract.CreateMyBlockPackChannelTicketOperation,
		"/core/v1/realtime/block-pack-channel-ticket/create",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeCreatedClientResponse(ctx, response.Data)
}
