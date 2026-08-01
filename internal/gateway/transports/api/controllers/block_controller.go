package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	blocksdto "github.com/HiIamJeff67/notezy-backend/contracts/api/v1/blocks"
	responsewriter "github.com/HiIamJeff67/notezy-backend/internal/gateway/responsewriter"
	coreadapters "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/core/adapters"
)

type BlockControllerInterface interface {
	GetMyBlockById(ctx *gin.Context, requestDto *blocksdto.GetMyBlockByIdRequestDto)
	GetMyBlocksByIds(ctx *gin.Context, requestDto *blocksdto.GetMyBlocksByIdsRequestDto)
	GetMyBlocksByBlockPackId(ctx *gin.Context, requestDto *blocksdto.GetMyBlocksByBlockPackIdRequestDto)
}

type BlockController struct {
	coreClient *coreadapters.CoreClient
}

func NewBlockController(coreClient *coreadapters.CoreClient) BlockControllerInterface {
	return &BlockController{
		coreClient: coreClient,
	}
}

func (c *BlockController) GetMyBlockById(ctx *gin.Context, requestDto *blocksdto.GetMyBlockByIdRequestDto) {
	response, exception := coreadapters.CallSecurly[
		blocksdto.GetMyBlockByIdRequestDto,
		blocksdto.GetMyBlockByIdResponseDto,
	](
		ctx,
		c.coreClient,
		requestDto,
		blocksdto.GetMyBlockByIdOperation,
		"/core/v1/blocks/get-by-id",
	)
	if exception != nil {
		responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      response.Data,
		"exception": nil,
	})
}

func (c *BlockController) GetMyBlocksByIds(ctx *gin.Context, requestDto *blocksdto.GetMyBlocksByIdsRequestDto) {
	response, exception := coreadapters.CallSecurly[
		blocksdto.GetMyBlocksByIdsRequestDto,
		blocksdto.GetMyBlocksByIdsResponseDto,
	](
		ctx,
		c.coreClient,
		requestDto,
		blocksdto.GetMyBlocksByIdsOperation,
		"/core/v1/blocks/get-by-ids",
	)
	if exception != nil {
		responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      response.Data,
		"exception": nil,
	})
}

func (c *BlockController) GetMyBlocksByBlockPackId(ctx *gin.Context, requestDto *blocksdto.GetMyBlocksByBlockPackIdRequestDto) {
	response, exception := coreadapters.CallSecurly[
		blocksdto.GetMyBlocksByBlockPackIdRequestDto,
		blocksdto.GetMyBlocksByBlockPackIdResponseDto,
	](
		ctx,
		c.coreClient,
		requestDto,
		blocksdto.GetMyBlocksByBlockPackIdOperation,
		"/core/v1/blocks/get-by-block-pack-id",
	)
	if exception != nil {
		responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      response.Data,
		"exception": nil,
	})
}
