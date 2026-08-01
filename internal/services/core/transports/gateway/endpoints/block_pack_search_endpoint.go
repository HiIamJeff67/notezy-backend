package endpoints

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	blockpacksdto "github.com/HiIamJeff67/notezy-backend/contracts/api/v1/block-packs"
	core "github.com/HiIamJeff67/notezy-backend/contracts/core/v1"
	contexts "github.com/HiIamJeff67/notezy-backend/internal/services/core/contexts"
)

func (t *BlockPackEndpoint) SearchBlockPacks(ctx *gin.Context) {
	request := &core.Request[blockpacksdto.SearchBlockPacksRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	userId, exception := contexts.GetActorUserId(ctx.Request.Context())
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

	responseDto, exception := t.blockPackService.SearchPrivateBlockPacks(ctx.Request.Context(), userId, request.Dto)
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

	ctx.JSON(http.StatusOK, core.Response[blockpacksdto.SearchBlockPacksResponseDto]{
		Version: core.Version,
		Metadata: core.ResponseMetadata{
			RequestId:   request.Metadata.RequestId,
			RespondedAt: time.Now(),
		},
		Data: *responseDto,
	})
}
