package endpoints

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	apicontract "github.com/HiIamJeff67/notegic-backend/contracts/core/v1/api/realtime"
	gatewaycontract "github.com/HiIamJeff67/notegic-backend/contracts/gateway/v1"

	realtimeservices "github.com/HiIamJeff67/notegic-backend/internal/core/services/realtime"
)

type RealtimeEndpointInterface interface {
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

func (t *RealtimeEndpoint) CreateMyRealtimeConnectionTicket(ctx *gin.Context) {
	request := &gatewaycontract.Request[apicontract.CreateMyRealtimeConnectionTicketRequestDto]{}
	if err := ctx.ShouldBindBodyWithJSON(request); err != nil {
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

	ctx.JSON(http.StatusOK, gatewaycontract.Response[apicontract.CreateMyRealtimeConnectionTicketResponseDto]{
		Version: gatewaycontract.Version,
		Metadata: gatewaycontract.ResponseMetadata{
			RequestId:   request.Metadata.RequestId,
			RespondedAt: time.Now(),
		},
		Data: *responseDto,
	})
}

func (t *RealtimeEndpoint) CreateMyBlockPackChannelTicket(ctx *gin.Context) {
	request := &gatewaycontract.Request[apicontract.CreateMyBlockPackChannelTicketRequestDto]{}
	if err := ctx.ShouldBindBodyWithJSON(request); err != nil {
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

	ctx.JSON(http.StatusOK, gatewaycontract.Response[apicontract.CreateMyBlockPackChannelTicketResponseDto]{
		Version: gatewaycontract.Version,
		Metadata: gatewaycontract.ResponseMetadata{
			RequestId:   request.Metadata.RequestId,
			RespondedAt: time.Now(),
		},
		Data: *responseDto,
	})
}
