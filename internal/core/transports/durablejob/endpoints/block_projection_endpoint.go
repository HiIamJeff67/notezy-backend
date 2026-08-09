package endpoints

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	exceptions "github.com/HiIamJeff67/notezy-backend/contracts/types/exceptions"

	apicontract "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api/blocks"
	durablejobdto "github.com/HiIamJeff67/notezy-backend/contracts/durable-job/v1"
	gatewaycontract "github.com/HiIamJeff67/notezy-backend/contracts/gateway/v1"

	blockservices "github.com/HiIamJeff67/notezy-backend/internal/core/services/blocks"
)

type BlockProjectionEndpoint struct {
	blockService blockservices.BlockServiceInterface
}

func NewBlockProjectionEndpoint(blockService blockservices.BlockServiceInterface) BlockProjectionEndpoint {
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

	documents := make([]apicontract.ApplyBlockProjectionDocumentRequestDto, len(request.Dto.Documents))
	for index, document := range request.Dto.Documents {
		documents[index] = apicontract.ApplyBlockProjectionDocumentRequestDto{
			BlockPackId: document.BlockPackId,
			Projection: apicontract.ApplyBlockProjectionRequestDto{
				SchemaId:          document.Projection.SchemaId,
				SchemaVersion:     document.Projection.SchemaVersion,
				ProjectedSequence: document.Projection.ProjectedSequence,
				Blocks:            document.Projection.Blocks,
			},
		}
	}

	responseDto, err := e.blockService.ApplyMany(
		ctx.Request.Context(),
		documents,
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
