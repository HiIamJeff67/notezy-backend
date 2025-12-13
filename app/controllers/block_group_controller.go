package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	dtos "notezy-backend/app/dtos"
	services "notezy-backend/app/services"
)

/* ============================== Interface & Instance ============================== */

type BlockGroupControllerInterface interface {
	GetMyBlockGroupAndItsBlocksById(ctx *gin.Context, reqDto *dtos.GetMyBlockGroupAndItsBlocksByIdReqDto)
	CreateBlockGroupAndItsBlocksByBlockPackId(ctx *gin.Context, reqDto *dtos.CreateBlockGroupAndItsBlocksByBlockPackIdReqDto)
}

type BlockGroupController struct {
	blockGroupService services.BlockGroupServiceInterface
}

func NewBlockGroupController(blockGroupService services.BlockGroupServiceInterface) BlockGroupControllerInterface {
	return &BlockGroupController{
		blockGroupService: blockGroupService,
	}
}

/* ============================== Implementations ============================== */

func (c *BlockGroupController) GetMyBlockGroupAndItsBlocksById(ctx *gin.Context, reqDto *dtos.GetMyBlockGroupAndItsBlocksByIdReqDto) {
	resDto, exception := c.blockGroupService.GetMyBlockGroupAndItsBlocksById(ctx.Request.Context(), reqDto)
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

func (c *BlockGroupController) CreateBlockGroupAndItsBlocksByBlockPackId(ctx *gin.Context, reqDto *dtos.CreateBlockGroupAndItsBlocksByBlockPackIdReqDto) {
	resDto, exception := c.blockGroupService.CreateBlockGroupAndItsBlocksByBlockPackId(ctx.Request.Context(), reqDto)
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
