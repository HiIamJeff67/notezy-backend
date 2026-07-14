package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	dtos "github.com/HiIamJeff67/notezy-backend/app/dtos"
	services "github.com/HiIamJeff67/notezy-backend/app/services"
)

type BlockControllerInterface interface {
	GetMyBlockById(ctx *gin.Context, reqDto *dtos.GetMyBlockByIdReqDto)
	GetMyBlocksByIds(ctx *gin.Context, reqDto *dtos.GetMyBlocksByIdsReqDto)
	GetMyBlocksByBlockPackId(ctx *gin.Context, reqDto *dtos.GetMyBlocksByBlockPackIdReqDto)
}

type BlockController struct {
	blockService services.BlockServiceInterface
}

func NewBlockController(blockService services.BlockServiceInterface) BlockControllerInterface {
	return &BlockController{
		blockService: blockService,
	}
}

func (c *BlockController) GetMyBlockById(ctx *gin.Context, reqDto *dtos.GetMyBlockByIdReqDto) {
	resDto, exception := c.blockService.GetMyBlockById(ctx.Request.Context(), reqDto)
	if exception != nil {
		exception.Log().SafelyAbortAndResponseWithJSON(ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      resDto,
		"exception": nil,
	})
}

func (c *BlockController) GetMyBlocksByIds(ctx *gin.Context, reqDto *dtos.GetMyBlocksByIdsReqDto) {
	resDto, exception := c.blockService.GetMyBlocksByIds(ctx.Request.Context(), reqDto)
	if exception != nil {
		exception.Log().SafelyAbortAndResponseWithJSON(ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      resDto,
		"exception": nil,
	})
}

func (c *BlockController) GetMyBlocksByBlockPackId(ctx *gin.Context, reqDto *dtos.GetMyBlocksByBlockPackIdReqDto) {
	resDto, exception := c.blockService.GetMyBlocksByBlockPackId(ctx.Request.Context(), reqDto)
	if exception != nil {
		exception.Log().SafelyAbortAndResponseWithJSON(ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      resDto,
		"exception": nil,
	})
}
