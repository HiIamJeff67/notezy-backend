package endpoints

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	realtimedto "github.com/HiIamJeff67/notezy-backend/contracts/api/v1/realtime"
	core "github.com/HiIamJeff67/notezy-backend/contracts/core/v1"
	services "github.com/HiIamJeff67/notezy-backend/internal/services/core/services"
)

type RealtimeEndpointInterface interface {
	GetMyBlockPackRealtimeParticipants(ctx *gin.Context)
	CreateMyRealtimeConnectionTicket(ctx *gin.Context)
	CreateMyBlockPackChannelTicket(ctx *gin.Context)
}

type RealtimeEndpoint struct {
	realtimeService services.RealtimeServiceInterface
}

func NewRealtimeEndpoint(
	realtimeService services.RealtimeServiceInterface,
) RealtimeEndpointInterface {
	return &RealtimeEndpoint{
		realtimeService: realtimeService,
	}
}

func (t *RealtimeEndpoint) GetMyBlockPackRealtimeParticipants(ctx *gin.Context) {
	request := &core.Request[realtimedto.GetMyBlockPackRealtimeParticipantsRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.realtimeService.GetMyBlockPackRealtimeParticipants(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), core.Response[struct{}]{
			Version: core.Version,
			Metadata: core.ResponseMetadata{
				RequestId:   request.Metadata.RequestId,
				RespondedAt: time.Now(),
			},
			Data:      struct{}{},
			Exception: publicException,
		})
		return
	}

	ctx.JSON(http.StatusOK, core.Response[realtimedto.GetMyBlockPackRealtimeParticipantsResponseDto]{
		Version: core.Version,
		Metadata: core.ResponseMetadata{
			RequestId:   request.Metadata.RequestId,
			RespondedAt: time.Now(),
		},
		Data: *responseDto,
	})
}

func (t *RealtimeEndpoint) CreateMyRealtimeConnectionTicket(ctx *gin.Context) {
	request := &core.Request[realtimedto.CreateMyRealtimeConnectionTicketRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.realtimeService.CreateMyRealtimeConnectionTicket(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), core.Response[struct{}]{
			Version: core.Version,
			Metadata: core.ResponseMetadata{
				RequestId:   request.Metadata.RequestId,
				RespondedAt: time.Now(),
			},
			Data:      struct{}{},
			Exception: publicException,
		})
		return
	}

	ctx.JSON(http.StatusOK, core.Response[realtimedto.CreateMyRealtimeConnectionTicketResponseDto]{
		Version: core.Version,
		Metadata: core.ResponseMetadata{
			RequestId:   request.Metadata.RequestId,
			RespondedAt: time.Now(),
		},
		Data: *responseDto,
	})
}

func (t *RealtimeEndpoint) CreateMyBlockPackChannelTicket(ctx *gin.Context) {
	request := &core.Request[realtimedto.CreateMyBlockPackChannelTicketRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.realtimeService.CreateMyBlockPackChannelTicket(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), core.Response[struct{}]{
			Version: core.Version,
			Metadata: core.ResponseMetadata{
				RequestId:   request.Metadata.RequestId,
				RespondedAt: time.Now(),
			},
			Data:      struct{}{},
			Exception: publicException,
		})
		return
	}

	ctx.JSON(http.StatusOK, core.Response[realtimedto.CreateMyBlockPackChannelTicketResponseDto]{
		Version: core.Version,
		Metadata: core.ResponseMetadata{
			RequestId:   request.Metadata.RequestId,
			RespondedAt: time.Now(),
		},
		Data: *responseDto,
	})
}
