package endpoints

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	itemsdto "github.com/HiIamJeff67/notezy-backend/contracts/api/v1/items"
	core "github.com/HiIamJeff67/notezy-backend/contracts/core/v1"
	contexts "github.com/HiIamJeff67/notezy-backend/internal/services/core/contexts"
	services "github.com/HiIamJeff67/notezy-backend/internal/services/core/services"
)

type ItemEndpointInterface interface {
	SearchItems(ctx *gin.Context)
}

type ItemEndpoint struct {
	itemService services.ItemServiceInterface
}

func NewItemEndpoint(
	itemService services.ItemServiceInterface,
) ItemEndpointInterface {
	return &ItemEndpoint{
		itemService: itemService,
	}
}

func (t *ItemEndpoint) SearchItems(ctx *gin.Context) {
	request := &core.Request[itemsdto.SearchItemsRequestDto]{}
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

	responseDto, exception := t.itemService.SearchPrivateItems(ctx.Request.Context(), userId, request.Dto)
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

	ctx.JSON(http.StatusOK, core.Response[itemsdto.SearchItemsResponseDto]{
		Version: core.Version,
		Metadata: core.ResponseMetadata{
			RequestId:   request.Metadata.RequestId,
			RespondedAt: time.Now(),
		},
		Data: *responseDto,
	})
}
