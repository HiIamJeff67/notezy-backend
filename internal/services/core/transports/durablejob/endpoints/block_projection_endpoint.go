package endpoints

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	durablejobdto "github.com/HiIamJeff67/notezy-backend/contracts/durablejob/v1"
	gatewaycontract "github.com/HiIamJeff67/notezy-backend/contracts/gateway/v1"
	exceptions "github.com/HiIamJeff67/notezy-backend/internal/exceptions"
	services "github.com/HiIamJeff67/notezy-backend/internal/services/core/services"
)

type BlockProjectionEndpoint struct {
	blockService services.BlockServiceInterface
}

func NewBlockProjectionEndpoint(blockService services.BlockServiceInterface) BlockProjectionEndpoint {
	return BlockProjectionEndpoint{blockService: blockService}
}

func (e BlockProjectionEndpoint) Apply(ctx *gin.Context) {
	request := &gatewaycontract.Request[durablejobdto.ApplyBlockProjectionRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		exception := exceptions.New(
			"InvalidRequest",
			"DurableJob",
			durablejobdto.ApplyBlockProjectionOperation,
			"The DurableJob projection request is invalid",
			http.StatusBadRequest,
		).WithOrigin(err)
		ctx.JSON(exception.HTTPStatusCode(), gatewaycontract.Response[durablejobdto.ApplyBlockProjectionResponseDto]{
			Version: gatewaycontract.Version,
			Metadata: gatewaycontract.ResponseMetadata{
				RequestId:   request.Metadata.RequestId,
				RespondedAt: time.Now(),
			},
			Exception: exception,
		})
		return
	}

	responseDto, err := e.blockService.ApplyMany(
		ctx.Request.Context(),
		request.Dto.Documents,
	)
	if err != nil {
		exception := exceptions.New(
			"FailedToApplyProjection",
			"DurableJob",
			durablejobdto.ApplyBlockProjectionOperation,
			"Failed to apply projected blocks",
			http.StatusInternalServerError,
			true,
		).WithOrigin(err)
		ctx.JSON(exception.HTTPStatusCode(), gatewaycontract.Response[durablejobdto.ApplyBlockProjectionResponseDto]{
			Version: gatewaycontract.Version,
			Metadata: gatewaycontract.ResponseMetadata{
				RequestId:   request.Metadata.RequestId,
				RespondedAt: time.Now(),
			},
			Exception: exception,
		})
		return
	}

	ctx.JSON(http.StatusOK, gatewaycontract.Response[durablejobdto.ApplyBlockProjectionResponseDto]{
		Version: gatewaycontract.Version,
		Metadata: gatewaycontract.ResponseMetadata{
			RequestId:   request.Metadata.RequestId,
			RespondedAt: time.Now(),
		},
		Data: durablejobdto.ApplyBlockProjectionResponseDto{
			AppliedBlockPackIds: []uuid.UUID(responseDto),
		},
	})
}
