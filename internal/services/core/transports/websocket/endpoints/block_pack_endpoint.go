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
	sharedtypes "github.com/HiIamJeff67/notezy-backend/internal/shared/types"
)

type BlockPackEndpoint struct {
	realtimeService services.RealtimeServiceInterface
}

func NewBlockPackEndpoint(realtimeService services.RealtimeServiceInterface) BlockPackEndpoint {
	return BlockPackEndpoint{
		realtimeService: realtimeService,
	}
}

func (e BlockPackEndpoint) ValidateBlockPackChannelPermission(ctx *gin.Context) {
	request := &core.Request[websocketdto.ValidateBlockPackChannelPermissionRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		exception := exceptions.New(
			"InvalidRequest",
			"WebSocket",
			websocketdto.ValidateBlockPackChannelPermissionOperation,
			"The WebSocket channel permission request is invalid",
			http.StatusBadRequest,
		).WithOrigin(err)
		ctx.JSON(exception.HTTPStatusCode(), core.Response[websocketdto.ValidateBlockPackChannelPermissionResponseDto]{
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
			websocketdto.ValidateBlockPackChannelPermissionOperation,
			"The WebSocket channel permission request is invalid",
			http.StatusBadRequest,
		).WithOrigin(err)
		ctx.JSON(exception.HTTPStatusCode(), core.Response[websocketdto.ValidateBlockPackChannelPermissionResponseDto]{
			Version: core.Version,
			Metadata: core.ResponseMetadata{
				RequestId:   request.Metadata.RequestId,
				RespondedAt: time.Now(),
			},
			Exception: exception,
		})
		return
	}

	errorCode, err := e.realtimeService.ValidateBlockPackChannelPermission(
		ctx.Request.Context(),
		request.Dto.UserPublicId,
		request.Dto.BlockPackId,
		sharedtypes.ChannelPermission(request.Dto.Permission),
	)
	if err != nil {
		ctx.JSON(http.StatusOK, core.Response[websocketdto.ValidateBlockPackChannelPermissionResponseDto]{
			Version: core.Version,
			Metadata: core.ResponseMetadata{
				RequestId:   request.Metadata.RequestId,
				RespondedAt: time.Now(),
			},
			Data: websocketdto.ValidateBlockPackChannelPermissionResponseDto{
				Permitted: false,
				ErrorCode: string(errorCode),
			},
		})
		return
	}

	ctx.JSON(http.StatusOK, core.Response[websocketdto.ValidateBlockPackChannelPermissionResponseDto]{
		Version: core.Version,
		Metadata: core.ResponseMetadata{
			RequestId:   request.Metadata.RequestId,
			RespondedAt: time.Now(),
		},
		Data: websocketdto.ValidateBlockPackChannelPermissionResponseDto{
			Permitted: true,
		},
	})
}
