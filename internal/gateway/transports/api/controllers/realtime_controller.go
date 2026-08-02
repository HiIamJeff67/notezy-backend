package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	realtimedto "github.com/HiIamJeff67/notezy-backend/contracts/gateway/v1/api/realtime"
	exceptions "github.com/HiIamJeff67/notezy-backend/internal/exceptions"
	responsewriter "github.com/HiIamJeff67/notezy-backend/internal/shared/responsewriter"
	coreadapters "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/core/adapters"
	realtimelease "github.com/HiIamJeff67/notezy-backend/internal/realtimegateway/data/cache/realtimelease"
)

type RealtimeControllerInterface interface {
	GetMyBlockPackRealtimeParticipants(ctx *gin.Context, requestDto *realtimedto.GetMyBlockPackRealtimeParticipantsRequestDto)
	CreateMyRealtimeConnectionTicket(ctx *gin.Context, requestDto *realtimedto.CreateMyRealtimeConnectionTicketRequestDto)
	CreateMyBlockPackChannelTicket(ctx *gin.Context, requestDto *realtimedto.CreateMyBlockPackChannelTicketRequestDto)
}

type RealtimeController struct {
	coreClient         *coreadapters.CoreClient
	realtimeLeaseCache *realtimelease.RealtimeLeaseCacheClient
}

func NewRealtimeController(coreClient *coreadapters.CoreClient) RealtimeControllerInterface {
	return &RealtimeController{
		coreClient:         coreClient,
		realtimeLeaseCache: realtimelease.NewRealtimeLeaseCacheClient(),
	}
}

func (c *RealtimeController) GetMyBlockPackRealtimeParticipants(
	ctx *gin.Context, requestDto *realtimedto.GetMyBlockPackRealtimeParticipantsRequestDto,
) {
	participants, err := c.realtimeLeaseCache.GetBlockPackParticipants(requestDto.Param.BlockPackId)
	if err != nil {
		responsewriter.SafelyAbortAndResponseWithJSON(
			exceptions.New(
				"Unavailable",
				"Realtime",
				"GetMyBlockPackRealtimeParticipants",
				"Realtime participant presence is unavailable",
				http.StatusServiceUnavailable,
			).WithOrigin(err),
			ctx,
		)
		return
	}

	requestDto.Body.Participants = make(
		[]realtimedto.RealtimeBlockPackParticipantRequestDto,
		len(participants),
	)
	for index, participant := range participants {
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
