package endpoints

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	core "github.com/HiIamJeff67/notezy-backend/contracts/core/v1"
	websocketdto "github.com/HiIamJeff67/notezy-backend/contracts/gateway/v1/websocket"
	exceptions "github.com/HiIamJeff67/notezy-backend/internal/exceptions"
	services "github.com/HiIamJeff67/notezy-backend/internal/services/core/services"
	validation "github.com/HiIamJeff67/notezy-backend/internal/services/core/validation"
)

type BlockProjectionEndpoint struct {
	blockService services.BlockServiceInterface
}

func NewBlockProjectionEndpoint(blockService services.BlockServiceInterface) BlockProjectionEndpoint {
	return BlockProjectionEndpoint{
		blockService: blockService,
	}
}

func (e BlockProjectionEndpoint) ApplyBlockProjection(ctx *gin.Context) {
	request := &core.Request[websocketdto.ApplyBlockProjectionRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		exception := exceptions.New(
			"InvalidRequest",
			"WebSocket",
			websocketdto.ApplyBlockProjectionOperation,
			"The block projection request is invalid",
			http.StatusBadRequest,
		).WithOrigin(err)
		ctx.JSON(exception.HTTPStatusCode(), core.Response[websocketdto.ApplyBlockProjectionResponseDto]{
			Version: core.Version,
			Metadata: core.ResponseMetadata{
				RequestId:   request.Metadata.RequestId,
				RespondedAt: time.Now(),
			},
			Exception: exception,
		})
		return
	}
	if err := validation.Validator.Struct(request.Dto); err != nil {
		exception := exceptions.New(
			"InvalidRequest",
			"WebSocket",
			websocketdto.ApplyBlockProjectionOperation,
			"The block projection request is invalid",
			http.StatusBadRequest,
		).WithOrigin(err)
		ctx.JSON(exception.HTTPStatusCode(), core.Response[websocketdto.ApplyBlockProjectionResponseDto]{
			Version: core.Version,
			Metadata: core.ResponseMetadata{
				RequestId:   request.Metadata.RequestId,
				RespondedAt: time.Now(),
			},
			Exception: exception,
		})
		return
	}

	responseDto, err := e.blockService.Apply(
		ctx.Request.Context(),
		request.Dto.BlockPackId,
		request.Dto.Projection,
	)
	if err != nil {
		exception := exceptions.New(
			"FailedToApplyProjection",
			"Block",
			websocketdto.ApplyBlockProjectionOperation,
			"Failed to apply the block projection",
			http.StatusInternalServerError,
			true,
		).WithOrigin(err)
		ctx.JSON(exception.HTTPStatusCode(), core.Response[websocketdto.ApplyBlockProjectionResponseDto]{
			Version: core.Version,
			Metadata: core.ResponseMetadata{
				RequestId:   request.Metadata.RequestId,
				RespondedAt: time.Now(),
			},
			Exception: exception,
		})
		return
	}

	ctx.JSON(http.StatusOK, core.Response[websocketdto.ApplyBlockProjectionResponseDto]{
		Version: core.Version,
		Metadata: core.ResponseMetadata{
			RequestId:   request.Metadata.RequestId,
			RespondedAt: time.Now(),
		},
		Data: websocketdto.ApplyBlockProjectionResponseDto{
			Applied:                responseDto.Applied,
			ProjectedUntilSequence: responseDto.ProjectedUntilSequence,
		},
	})
}
