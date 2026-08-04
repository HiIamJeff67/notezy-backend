package endpoints

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	realtimegatewaycontract "github.com/HiIamJeff67/notezy-backend/contracts/realtimegateway/v1"
	exceptions "github.com/HiIamJeff67/notezy-backend/internal/exceptions"
	realtimelease "github.com/HiIamJeff67/notezy-backend/internal/realtimegateway/data/cache/realtimelease"
)

type BlockPackEndpoint struct {
	realtimeLeaseCache *realtimelease.RealtimeLeaseCacheClient
}

func NewBlockPackEndpoint(realtimeLeaseCache *realtimelease.RealtimeLeaseCacheClient) BlockPackEndpoint {
	return BlockPackEndpoint{
		realtimeLeaseCache: realtimeLeaseCache,
	}
}

func (e BlockPackEndpoint) GetParticipants(ctx *gin.Context) {
	request := &realtimegatewaycontract.Request[realtimegatewaycontract.GetBlockPackParticipantsRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil || request.Dto.BlockPackId == uuid.Nil {
		exception := exceptions.New(
			"InvalidRequest",
			"RealtimeGateway",
			realtimegatewaycontract.GetBlockPackParticipantsOperation,
			"The RealtimeGateway participant request is invalid",
			http.StatusBadRequest,
		)
		if err != nil {
			exception = exception.WithOrigin(err)
		}
		ctx.JSON(exception.HTTPStatusCode(), realtimegatewaycontract.Response[realtimegatewaycontract.GetBlockPackParticipantsResponseDto]{
			Version: realtimegatewaycontract.Version,
			Metadata: realtimegatewaycontract.ResponseMetadata{
				RequestId:   request.Metadata.RequestId,
				RespondedAt: time.Now(),
			},
			Data: realtimegatewaycontract.GetBlockPackParticipantsResponseDto{
				Participants: []realtimegatewaycontract.BlockPackParticipantResponseDto{},
			},
			Exception: exception,
		})
		return
	}
	if e.realtimeLeaseCache == nil {
		exception := exceptions.New(
			"RealtimeLeaseCacheRequired",
			"RealtimeGateway",
			realtimegatewaycontract.GetBlockPackParticipantsOperation,
			"The RealtimeGateway lease cache is unavailable",
			http.StatusServiceUnavailable,
			true,
		)
		ctx.JSON(exception.HTTPStatusCode(), realtimegatewaycontract.Response[realtimegatewaycontract.GetBlockPackParticipantsResponseDto]{
			Version: realtimegatewaycontract.Version,
			Metadata: realtimegatewaycontract.ResponseMetadata{
				RequestId:   request.Metadata.RequestId,
				RespondedAt: time.Now(),
			},
			Data: realtimegatewaycontract.GetBlockPackParticipantsResponseDto{
				Participants: []realtimegatewaycontract.BlockPackParticipantResponseDto{},
			},
			Exception: exception,
		})
		return
	}

	participants, err := e.realtimeLeaseCache.GetBlockPackParticipants(request.Dto.BlockPackId)
	if err != nil {
		exception := exceptions.New(
			"RealtimePresenceUnavailable",
			"RealtimeGateway",
			realtimegatewaycontract.GetBlockPackParticipantsOperation,
			"Realtime participant presence is unavailable",
			http.StatusServiceUnavailable,
			true,
		).WithOrigin(err)
		ctx.JSON(exception.HTTPStatusCode(), realtimegatewaycontract.Response[realtimegatewaycontract.GetBlockPackParticipantsResponseDto]{
			Version: realtimegatewaycontract.Version,
			Metadata: realtimegatewaycontract.ResponseMetadata{
				RequestId:   request.Metadata.RequestId,
				RespondedAt: time.Now(),
			},
			Data: realtimegatewaycontract.GetBlockPackParticipantsResponseDto{
				Participants: []realtimegatewaycontract.BlockPackParticipantResponseDto{},
			},
			Exception: exception,
		})
		return
	}

	responseDto := realtimegatewaycontract.GetBlockPackParticipantsResponseDto{
		Participants: make([]realtimegatewaycontract.BlockPackParticipantResponseDto, len(participants)),
	}
	for index, participant := range participants {
		responseDto.Participants[index] = realtimegatewaycontract.BlockPackParticipantResponseDto{
			UserPublicId:      participant.UserPublicId,
			ChannelPermission: participant.ChannelPermission,
			ConnectionCount:   participant.ConnectionCount,
		}
	}

	ctx.JSON(http.StatusOK, realtimegatewaycontract.Response[realtimegatewaycontract.GetBlockPackParticipantsResponseDto]{
		Version: realtimegatewaycontract.Version,
		Metadata: realtimegatewaycontract.ResponseMetadata{
			RequestId:   request.Metadata.RequestId,
			RespondedAt: time.Now(),
		},
		Data: responseDto,
	})
}
