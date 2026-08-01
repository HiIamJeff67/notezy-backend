package endpoints

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	blocksdto "github.com/HiIamJeff67/notezy-backend/contracts/api/v1/blocks"
	core "github.com/HiIamJeff67/notezy-backend/contracts/core/v1"
	services "github.com/HiIamJeff67/notezy-backend/internal/services/core/services"
)

type BlockEndpointInterface interface {
	GetMyBlockById(ctx *gin.Context)
	GetMyBlocksByIds(ctx *gin.Context)
	GetMyBlocksByBlockPackId(ctx *gin.Context)

	/* ============================== GraphQL Methods ============================== */
	SearchBlocks(ctx *gin.Context)
}

type BlockEndpoint struct {
	blockService services.BlockServiceInterface
}

func NewBlockEndpoint(
	blockService services.BlockServiceInterface,
) BlockEndpointInterface {
	return &BlockEndpoint{
		blockService: blockService,
	}
}

func (t *BlockEndpoint) GetMyBlockById(ctx *gin.Context) {
	request := &core.Request[blocksdto.GetMyBlockByIdRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.blockService.GetMyBlockById(ctx.Request.Context(), &request.Dto)
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

	ctx.JSON(http.StatusOK, core.Response[blocksdto.GetMyBlockByIdResponseDto]{
		Version: core.Version,
		Metadata: core.ResponseMetadata{
			RequestId:   request.Metadata.RequestId,
			RespondedAt: time.Now(),
		},
		Data: *responseDto,
	})
}

func (t *BlockEndpoint) GetMyBlocksByIds(ctx *gin.Context) {
	request := &core.Request[blocksdto.GetMyBlocksByIdsRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.blockService.GetMyBlocksByIds(ctx.Request.Context(), &request.Dto)
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

	ctx.JSON(http.StatusOK, core.Response[blocksdto.GetMyBlocksByIdsResponseDto]{
		Version: core.Version,
		Metadata: core.ResponseMetadata{
			RequestId:   request.Metadata.RequestId,
			RespondedAt: time.Now(),
		},
		Data: *responseDto,
	})
}

func (t *BlockEndpoint) GetMyBlocksByBlockPackId(ctx *gin.Context) {
	request := &core.Request[blocksdto.GetMyBlocksByBlockPackIdRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.blockService.GetMyBlocksByBlockPackId(ctx.Request.Context(), &request.Dto)
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

	ctx.JSON(http.StatusOK, core.Response[blocksdto.GetMyBlocksByBlockPackIdResponseDto]{
		Version: core.Version,
		Metadata: core.ResponseMetadata{
			RequestId:   request.Metadata.RequestId,
			RespondedAt: time.Now(),
		},
		Data: *responseDto,
	})
}
