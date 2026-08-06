package endpoints

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	realtimedto "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api/realtime"
	gatewaycontract "github.com/HiIamJeff67/notezy-backend/contracts/gateway/v1"

	realtimeservices "github.com/HiIamJeff67/notezy-backend/internal/core/services/realtime"
)

type RealtimeEndpointInterface interface {
	GetMyBlockPackRealtimeParticipants(ctx *gin.Context)
	CreateMyRealtimeConnectionTicket(ctx *gin.Context)
	CreateMyBlockPackChannelTicket(ctx *gin.Context)
}

type RealtimeEndpoint struct {
	realtimeService realtimeservices.RealtimeServiceInterface
}

func NewRealtimeEndpoint(
	realtimeService realtimeservices.RealtimeServiceInterface,
) RealtimeEndpointInterface {
	return &RealtimeEndpoint{
		realtimeService: realtimeService,
	}
}

func (t *RealtimeEndpoint) GetMyBlockPackRealtimeParticipants(ctx *gin.Context) {
	request := &gatewaycontract.Request[realtimedto.GetMyBlockPackRealtimeParticipantsRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.realtimeService.GetMyBlockPackRealtimeParticipants(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), gatewaycontract.Response[struct{}]{
			Version: gatewaycontract.Version,
			Metadata: gatewaycontract.ResponseMetadata{
				RequestId:   request.Metadata.RequestId,
				RespondedAt: time.Now(),
			},
			Data:      struct{}{},
			Exception: publicException,
		})
		return
	}

	ctx.JSON(http.StatusOK, gatewaycontract.Response[realtimedto.GetMyBlockPackRealtimeParticipantsResponseDto]{
		Version: gatewaycontract.Version,
		Metadata: gatewaycontract.ResponseMetadata{
			RequestId:   request.Metadata.RequestId,
			RespondedAt: time.Now(),
		},
		Data: *responseDto,
	})
}

func (t *RealtimeEndpoint) CreateMyRealtimeConnectionTicket(ctx *gin.Context) {
	request := &gatewaycontract.Request[realtimedto.CreateMyRealtimeConnectionTicketRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.realtimeService.CreateMyRealtimeConnectionTicket(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), gatewaycontract.Response[struct{}]{
			Version: gatewaycontract.Version,
			Metadata: gatewaycontract.ResponseMetadata{
				RequestId:   request.Metadata.RequestId,
				RespondedAt: time.Now(),
			},
			Data:      struct{}{},
			Exception: publicException,
		})
		return
	}

	ctx.JSON(http.StatusOK, gatewaycontract.Response[realtimedto.CreateMyRealtimeConnectionTicketResponseDto]{
		Version: gatewaycontract.Version,
		Metadata: gatewaycontract.ResponseMetadata{
			RequestId:   request.Metadata.RequestId,
			RespondedAt: time.Now(),
		},
		Data: *responseDto,
	})
}

func (t *RealtimeEndpoint) CreateMyBlockPackChannelTicket(ctx *gin.Context) {
	request := &gatewaycontract.Request[realtimedto.CreateMyBlockPackChannelTicketRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.realtimeService.CreateMyBlockPackChannelTicket(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), gatewaycontract.Response[struct{}]{
			Version: gatewaycontract.Version,
			Metadata: gatewaycontract.ResponseMetadata{
				RequestId:   request.Metadata.RequestId,
				RespondedAt: time.Now(),
			},
			Data:      struct{}{},
			Exception: publicException,
		})
		return
	}

	ctx.JSON(http.StatusOK, gatewaycontract.Response[realtimedto.CreateMyBlockPackChannelTicketResponseDto]{
		Version: gatewaycontract.Version,
		Metadata: gatewaycontract.ResponseMetadata{
			RequestId:   request.Metadata.RequestId,
			RespondedAt: time.Now(),
		},
		Data: *responseDto,
	})
}
