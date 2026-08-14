package controllers

import (
	"github.com/gin-gonic/gin"

	exceptionwriter "github.com/HiIamJeff67/notezy-backend/shared/util/exceptionwriter"

	apicontract "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api/blocks"

	coreadapters "github.com/HiIamJeff67/notezy-backend/internal/apigateway/transports/core/adapters"
)

type BlockControllerInterface interface {
	GetMyBlockById(ctx *gin.Context, requestDto *apicontract.GetMyBlockByIdRequestDto)
	GetMyBlocksByIds(ctx *gin.Context, requestDto *apicontract.GetMyBlocksByIdsRequestDto)
	GetMyBlocksByBlockPackId(ctx *gin.Context, requestDto *apicontract.GetMyBlocksByBlockPackIdRequestDto)
}

type BlockController struct {
	coreClient *coreadapters.CoreAdapter
}

func NewBlockController(coreClient *coreadapters.CoreAdapter) BlockControllerInterface {
	return &BlockController{
		coreClient: coreClient,
	}
}

func (c *BlockController) GetMyBlockById(ctx *gin.Context, requestDto *apicontract.GetMyBlockByIdRequestDto) {
	response, exception := coreadapters.CallSecurly[
		apicontract.GetMyBlockByIdRequestDto,
		apicontract.GetMyBlockByIdResponseDto,
	](
		ctx,
		c.coreClient,
		requestDto,
		apicontract.GetMyBlockByIdOperation,
		"/core/v1/blocks/get-by-id",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *BlockController) GetMyBlocksByIds(ctx *gin.Context, requestDto *apicontract.GetMyBlocksByIdsRequestDto) {
	response, exception := coreadapters.CallSecurly[
		apicontract.GetMyBlocksByIdsRequestDto,
		apicontract.GetMyBlocksByIdsResponseDto,
	](
		ctx,
		c.coreClient,
		requestDto,
		apicontract.GetMyBlocksByIdsOperation,
		"/core/v1/blocks/get-by-ids",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *BlockController) GetMyBlocksByBlockPackId(ctx *gin.Context, requestDto *apicontract.GetMyBlocksByBlockPackIdRequestDto) {
	response, exception := coreadapters.CallSecurly[
		apicontract.GetMyBlocksByBlockPackIdRequestDto,
		apicontract.GetMyBlocksByBlockPackIdResponseDto,
	](
		ctx,
		c.coreClient,
		requestDto,
		apicontract.GetMyBlocksByBlockPackIdOperation,
		"/core/v1/blocks/get-by-block-pack-id",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}
