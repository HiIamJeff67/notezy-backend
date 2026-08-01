package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	realtimedto "github.com/HiIamJeff67/notezy-backend/contracts/api/v1/realtime"
	responsewriter "github.com/HiIamJeff67/notezy-backend/internal/gateway/responsewriter"
	coreadapters "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/core/adapters"
)

type RealtimeControllerInterface interface {
	GetMyBlockPackRealtimeParticipants(ctx *gin.Context, requestDto *realtimedto.GetMyBlockPackRealtimeParticipantsRequestDto)
	CreateMyRealtimeConnectionTicket(ctx *gin.Context, requestDto *realtimedto.CreateMyRealtimeConnectionTicketRequestDto)
	CreateMyBlockPackChannelTicket(ctx *gin.Context, requestDto *realtimedto.CreateMyBlockPackChannelTicketRequestDto)
}

type RealtimeController struct {
	coreClient *coreadapters.CoreClient
}

func NewRealtimeController(coreClient *coreadapters.CoreClient) RealtimeControllerInterface {
	return &RealtimeController{
		coreClient: coreClient,
	}
}

func (c *RealtimeController) GetMyBlockPackRealtimeParticipants(
	ctx *gin.Context, requestDto *realtimedto.GetMyBlockPackRealtimeParticipantsRequestDto,
) {
	response, exception := coreadapters.CallSecurly[
		realtimedto.GetMyBlockPackRealtimeParticipantsRequestDto,
		realtimedto.GetMyBlockPackRealtimeParticipantsResponseDto,
	](
		ctx,
		c.coreClient,
		requestDto,
		realtimedto.GetMyBlockPackRealtimeParticipantsOperation,
		"/core/v1/realtime/block-pack-participants/get",
	)
	if exception != nil {
		responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      response.Data,
		"exception": nil,
	})
}

func (c *RealtimeController) CreateMyRealtimeConnectionTicket(
	ctx *gin.Context,
	requestDto *realtimedto.CreateMyRealtimeConnectionTicketRequestDto,
) {
	response, exception := coreadapters.CallSecurly[
		realtimedto.CreateMyRealtimeConnectionTicketRequestDto,
		realtimedto.CreateMyRealtimeConnectionTicketResponseDto,
	](
		ctx,
		c.coreClient,
		requestDto,
		realtimedto.CreateMyRealtimeConnectionTicketOperation,
		"/core/v1/realtime/connection-ticket/create",
	)
	if exception != nil {
		responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      response.Data,
		"exception": nil,
	})
}

func (c *RealtimeController) CreateMyBlockPackChannelTicket(
	ctx *gin.Context,
	requestDto *realtimedto.CreateMyBlockPackChannelTicketRequestDto,
) {
	response, exception := coreadapters.CallSecurly[
		realtimedto.CreateMyBlockPackChannelTicketRequestDto,
		realtimedto.CreateMyBlockPackChannelTicketResponseDto,
	](
		ctx,
		c.coreClient,
		requestDto,
		realtimedto.CreateMyBlockPackChannelTicketOperation,
		"/core/v1/realtime/block-pack-channel-ticket/create",
	)
	if exception != nil {
		responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      response.Data,
		"exception": nil,
	})
}
