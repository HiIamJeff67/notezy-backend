package controllers

import (
	"net/http"
	"notezy-backend/app/dtos"
	"notezy-backend/app/services"

	"github.com/gin-gonic/gin"
)

/* ============================== Interface & Instance ============================== */

type BlockControllerInterface interface {
	GetMyBlockById(ctx *gin.Context, reqDto *dtos.GetMyBlockByIdReqDto)
	GetAllMyBlocks(ctx *gin.Context, reqDto *dtos.GetAllMyBlocksReqDto)
}

type BlockController struct {
	blockService services.BlockServiceInterface
}

func NewBlockController(blockService services.BlockServiceInterface) BlockControllerInterface {
	return &BlockController{
		blockService: blockService,
	}
}

/* ============================== Implementations ============================== */

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

func (c *BlockController) GetAllMyBlocks(ctx *gin.Context, reqDto *dtos.GetAllMyBlocksReqDto) {
	resDto, exception := c.blockService.GetAllMyBlocks(ctx.Request.Context(), reqDto)
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
