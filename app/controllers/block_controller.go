package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	dtos "notezy-backend/app/dtos"
	services "notezy-backend/app/services"
)

type BlockControllerInterface interface {
	GetMyBlockById(ctx *gin.Context, reqDto *dtos.GetMyBlockByIdReqDto)
	GetMyBlocksByIds(ctx *gin.Context, reqDto *dtos.GetMyBlocksByIdsReqDto)
	GetMyBlocksByBlockGroupId(ctx *gin.Context, reqDto *dtos.GetMyBlocksByBlockGroupIdReqDto)
	GetMyBlocksByBlockGroupIds(ctx *gin.Context, reqDto *dtos.GetMyBlocksByBlockGroupIdsReqDto)
	GetMyBlocksByBlockPackId(ctx *gin.Context, reqDto *dtos.GetMyBlocksByBlockPackIdReqDto)
	GetAllMyBlocks(ctx *gin.Context, reqDto *dtos.GetAllMyBlocksReqDto)
	InsertBlock(ctx *gin.Context, reqDto *dtos.InsertBlockReqDto)
	InsertBlocks(ctx *gin.Context, reqDto *dtos.InsertBlocksReqDto)
	UpdateMyBlockById(ctx *gin.Context, reqDto *dtos.UpdateMyBlockByIdReqDto)
	UpdateMyBlocksByIds(ctx *gin.Context, reqDto *dtos.UpdateMyBlocksByIdsReqDto)
	RestoreMyBlockById(ctx *gin.Context, reqDto *dtos.RestoreMyBlockByIdReqDto)
	RestoreMyBlocksByIds(ctx *gin.Context, reqDto *dtos.RestoreMyBlocksByIdsReqDto)
	DeleteMyBlockById(ctx *gin.Context, reqDto *dtos.DeleteMyBlockByIdReqDto)
	DeleteMyBlocksByIds(ctx *gin.Context, reqDto *dtos.DeleteMyBlocksByIdsReqDto)
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

func (c *BlockController) GetMyBlocksByBlockGroupId(ctx *gin.Context, reqDto *dtos.GetMyBlocksByBlockGroupIdReqDto) {
	resDto, exception := c.blockService.GetMyBlocksByBlockGroupId(ctx.Request.Context(), reqDto)
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

func (c *BlockController) GetMyBlocksByBlockGroupIds(ctx *gin.Context, reqDto *dtos.GetMyBlocksByBlockGroupIdsReqDto) {
	resDto, exception := c.blockService.GetMyBlocksByBlockGroupIds(ctx.Request.Context(), reqDto)
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

func (c *BlockController) InsertBlock(ctx *gin.Context, reqDto *dtos.InsertBlockReqDto) {
	resDto, exception := c.blockService.InsertBlock(ctx.Request.Context(), reqDto)
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

func (c *BlockController) InsertBlocks(ctx *gin.Context, reqDto *dtos.InsertBlocksReqDto) {
	resDto, exception := c.blockService.InsertBlocks(ctx.Request.Context(), reqDto)
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

func (c *BlockController) UpdateMyBlockById(ctx *gin.Context, reqDto *dtos.UpdateMyBlockByIdReqDto) {
	resDto, exception := c.blockService.UpdateMyBlockById(ctx.Request.Context(), reqDto)
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

func (c *BlockController) UpdateMyBlocksByIds(ctx *gin.Context, reqDto *dtos.UpdateMyBlocksByIdsReqDto) {
	resDto, exception := c.blockService.UpdateMyBlocksByIds(ctx.Request.Context(), reqDto)
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

func (c *BlockController) RestoreMyBlockById(ctx *gin.Context, reqDto *dtos.RestoreMyBlockByIdReqDto) {
	resDto, exception := c.blockService.RestoreMyBlockById(ctx.Request.Context(), reqDto)
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

func (c *BlockController) RestoreMyBlocksByIds(ctx *gin.Context, reqDto *dtos.RestoreMyBlocksByIdsReqDto) {
	resDto, exception := c.blockService.RestoreMyBlocksByIds(ctx.Request.Context(), reqDto)
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

func (c *BlockController) DeleteMyBlockById(ctx *gin.Context, reqDto *dtos.DeleteMyBlockByIdReqDto) {
	resDto, exception := c.blockService.DeleteMyBlockById(ctx.Request.Context(), reqDto)
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

func (c *BlockController) DeleteMyBlocksByIds(ctx *gin.Context, reqDto *dtos.DeleteMyBlocksByIdsReqDto) {
	resDto, exception := c.blockService.DeleteMyBlocksByIds(ctx.Request.Context(), reqDto)
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
