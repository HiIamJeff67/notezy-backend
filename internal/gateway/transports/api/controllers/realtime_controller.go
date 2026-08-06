package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	exceptionwriter "github.com/HiIamJeff67/notezy-backend/shared/util/exceptionwriter"

	realtimedto "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api/realtime"
	realtimegatewaycontract "github.com/HiIamJeff67/notezy-backend/contracts/realtime-gateway/v1"

	coreadapters "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/core/adapters"
	realtimegatewayadapters "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/realtimegateway/adapters"
)

type RealtimeControllerInterface interface {
	GetMyBlockPackRealtimeParticipants(ctx *gin.Context, requestDto *realtimedto.GetMyBlockPackRealtimeParticipantsRequestDto)
	CreateMyRealtimeConnectionTicket(ctx *gin.Context, requestDto *realtimedto.CreateMyRealtimeConnectionTicketRequestDto)
	CreateMyBlockPackChannelTicket(ctx *gin.Context, requestDto *realtimedto.CreateMyBlockPackChannelTicketRequestDto)
}

type RealtimeController struct {
	coreClient            *coreadapters.CoreClient
	realtimeGatewayClient *realtimegatewayadapters.RealtimeGatewayClient
}

func NewRealtimeController(
	coreClient *coreadapters.CoreClient,
	realtimeGatewayClient *realtimegatewayadapters.RealtimeGatewayClient,
) RealtimeControllerInterface {
	return &RealtimeController{
		coreClient:            coreClient,
		realtimeGatewayClient: realtimeGatewayClient,
	}
}

func (c *RealtimeController) GetMyBlockPackRealtimeParticipants(
	ctx *gin.Context, requestDto *realtimedto.GetMyBlockPackRealtimeParticipantsRequestDto,
) {
	participantsResponseDto, exception := c.realtimeGatewayClient.GetBlockPackParticipants(
		ctx,
		&realtimegatewaycontract.GetBlockPackParticipantsRequestDto{
			BlockPackId: requestDto.Param.BlockPackId,
		},
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	requestDto.Body.Participants = make(
		[]realtimedto.RealtimeBlockPackParticipantRequestDto,
		len(participantsResponseDto.Participants),
	)
	for index, participant := range participantsResponseDto.Participants {
		requestDto.Body.Participants[index] = realtimedto.RealtimeBlockPackParticipantRequestDto{
			UserPublicId:      participant.UserPublicId,
			ChannelPermission: participant.ChannelPermission,
			ConnectionCount:   participant.ConnectionCount,
		}
	}

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
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
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
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
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
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      response.Data,
		"exception": nil,
	})
}
