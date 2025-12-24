package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	dtos "notezy-backend/app/dtos"
	services "notezy-backend/app/services"
)

/* ============================== Interface & Instance ============================== */

type BlockGroupControllerInterface interface {
	GetMyBlockGroupById(ctx *gin.Context, reqDto *dtos.GetMyBlockGroupByIdReqDto)
	GetMyBlockGroupAndItsBlocksById(ctx *gin.Context, reqDto *dtos.GetMyBlockGroupAndItsBlocksByIdReqDto)
	GetMyBlockGroupsAndTheirBlocksByBlockPackId(ctx *gin.Context, reqDto *dtos.GetMyBlockGroupsAndTheirBlocksByBlockPackIdReqDto)
	GetMyBlockGroupsByPrevBlockGroupId(ctx *gin.Context, reqDto *dtos.GetMyBlockGroupsByPrevBlockGroupIdReqDto)
	GetAllMyBlockGroupsByBlockPackId(ctx *gin.Context, reqDto *dtos.GetAllMyBlockGroupsByBlockPackIdReqDto)
	InsertBlockGroupByBlockPackId(ctx *gin.Context, reqDto *dtos.InsertBlockGroupByBlockPackIdReqDto)
	InsertBlockGroupAndItsBlocksByBlockPackId(ctx *gin.Context, reqDto *dtos.InsertBlockGroupAndItsBlocksByBlockPackIdReqDto)
	InsertBlockGroupsAndTheirBlocksByBlockPackId(ctx *gin.Context, reqDto *dtos.InsertBlockGroupsAndTheirBlocksByBlockPackIdReqDto)
	InsertSequentialBlockGroupsAndTheirBlocksByBlockPackId(ctx *gin.Context, reqDto *dtos.InsertSequentialBlockGroupsAndTheirBlocksByBlockPackIdReqDto)
	MoveMyBlockGroupsByIds(ctx *gin.Context, reqDto *dtos.MoveMyBlockGroupsByIdsReqDto)
	RestoreMyBlockGroupById(ctx *gin.Context, reqDto *dtos.RestoreMyBlockGroupByIdReqDto)
	RestoreMyBlockGroupsByIds(ctx *gin.Context, reqDto *dtos.RestoreMyBlockGroupsByIdsReqDto)
	DeleteMyBlockGroupById(ctx *gin.Context, reqDto *dtos.DeleteMyBlockGroupByIdReqDto)
	DeleteMyBlockGroupsByIds(ctx *gin.Context, reqDto *dtos.DeleteMyBlockGroupsByIdsReqDto)
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

func (c *BlockGroupController) GetMyBlockGroupById(ctx *gin.Context, reqDto *dtos.GetMyBlockGroupByIdReqDto) {
	resDto, exception := c.blockGroupService.GetMyBlockGroupById(ctx.Request.Context(), reqDto)
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

func (c *BlockGroupController) GetMyBlockGroupsAndTheirBlocksByBlockPackId(ctx *gin.Context, reqDto *dtos.GetMyBlockGroupsAndTheirBlocksByBlockPackIdReqDto) {
	resDto, exception := c.blockGroupService.GetMyBlockGroupsAndTheirBlocksByBlockPackId(ctx.Request.Context(), reqDto)
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

func (c *BlockGroupController) GetMyBlockGroupsByPrevBlockGroupId(ctx *gin.Context, reqDto *dtos.GetMyBlockGroupsByPrevBlockGroupIdReqDto) {
	resDto, exception := c.blockGroupService.GetMyBlockGroupsByPrevBlockGroupId(ctx.Request.Context(), reqDto)
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

func (c *BlockGroupController) GetAllMyBlockGroupsByBlockPackId(ctx *gin.Context, reqDto *dtos.GetAllMyBlockGroupsByBlockPackIdReqDto) {
	resDto, exception := c.blockGroupService.GetAllMyBlockGroupsByBlockPackId(ctx.Request.Context(), reqDto)
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

func (c *BlockGroupController) InsertBlockGroupByBlockPackId(ctx *gin.Context, reqDto *dtos.InsertBlockGroupByBlockPackIdReqDto) {
	resDto, exception := c.blockGroupService.InsertBlockGroupByBlockPackId(ctx.Request.Context(), reqDto)
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

func (c *BlockGroupController) InsertBlockGroupAndItsBlocksByBlockPackId(ctx *gin.Context, reqDto *dtos.InsertBlockGroupAndItsBlocksByBlockPackIdReqDto) {
	resDto, exception := c.blockGroupService.InsertBlockGroupAndItsBlocksByBlockPackId(ctx.Request.Context(), reqDto)
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

func (c *BlockGroupController) InsertBlockGroupsAndTheirBlocksByBlockPackId(ctx *gin.Context, reqDto *dtos.InsertBlockGroupsAndTheirBlocksByBlockPackIdReqDto) {
	resDto, exception := c.blockGroupService.InsertBlockGroupsAndTheirBlocksByBlockPackId(ctx.Request.Context(), reqDto)
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

func (c *BlockGroupController) InsertSequentialBlockGroupsAndTheirBlocksByBlockPackId(ctx *gin.Context, reqDto *dtos.InsertSequentialBlockGroupsAndTheirBlocksByBlockPackIdReqDto) {
	resDto, exception := c.blockGroupService.InsertSequentialBlockGroupsAndTheirBlocksByBlockPackId(ctx.Request.Context(), reqDto)
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

func (c *BlockGroupController) MoveMyBlockGroupsByIds(ctx *gin.Context, reqDto *dtos.MoveMyBlockGroupsByIdsReqDto) {
	resDto, exception := c.blockGroupService.MoveMyBlockGroupsByIds(ctx.Request.Context(), reqDto)
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

func (c *BlockGroupController) RestoreMyBlockGroupById(ctx *gin.Context, reqDto *dtos.RestoreMyBlockGroupByIdReqDto) {
	resDto, exception := c.blockGroupService.RestoreMyBlockGroupById(ctx.Request.Context(), reqDto)
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

func (c *BlockGroupController) RestoreMyBlockGroupsByIds(ctx *gin.Context, reqDto *dtos.RestoreMyBlockGroupsByIdsReqDto) {
	resDto, exception := c.blockGroupService.RestoreMyBlockGroupsByIds(ctx.Request.Context(), reqDto)
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

func (c *BlockGroupController) DeleteMyBlockGroupById(ctx *gin.Context, reqDto *dtos.DeleteMyBlockGroupByIdReqDto) {
	resDto, exception := c.blockGroupService.DeleteMyBlockGroupById(ctx.Request.Context(), reqDto)
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

func (c *BlockGroupController) DeleteMyBlockGroupsByIds(ctx *gin.Context, reqDto *dtos.DeleteMyBlockGroupsByIdsReqDto) {
	resDto, exception := c.blockGroupService.DeleteMyBlockGroupsByIds(ctx.Request.Context(), reqDto)
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
