package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	dtos "notezy-backend/app/dtos"
	services "notezy-backend/app/services"
)

/* ============================== Interface & Instance ============================== */

type BlockPackControllerInterface interface {
	GetMyBlockPackById(ctx *gin.Context, reqDto *dtos.GetMyBlockPackByIdReqDto)
	GetMyBlockPackAndItsParentById(ctx *gin.Context, reqDto *dtos.GetMyBlockPackAndItsParentByIdReqDto)
	GetAllMyBlockPacksByParentSubShelfId(ctx *gin.Context, reqDto *dtos.GetAllMyBlockPacksByParentSubShelfIdReqDto)
	GetAllMyBlockPacksByRootShelfId(ctx *gin.Context, reqDto *dtos.GetAllMyBlockPacksByRootShelfIdReqDto)
	CreateBlockPack(ctx *gin.Context, reqDto *dtos.CreateBlockPackReqDto)
	UpdateMyBlockPackById(ctx *gin.Context, reqDto *dtos.UpdateMyBlockPackByIdReqDto)
	MoveMyBlockPackById(ctx *gin.Context, reqDto *dtos.MoveMyBlockPackByIdReqDto)
	MoveMyBlockPacksByIds(ctx *gin.Context, reqDto *dtos.MoveMyBlockPacksByIdsReqDto)
	RestoreMyBlockPackById(ctx *gin.Context, reqDto *dtos.RestoreMyBlockPackByIdReqDto)
	RestoreMyBlockPacksByIds(ctx *gin.Context, reqDto *dtos.RestoreMyBlockPacksByIdsReqDto)
	DeleteMyBlockPackById(ctx *gin.Context, reqDto *dtos.DeleteMyBlockPackByIdReqDto)
	DeleteMyBlockPacksByIds(ctx *gin.Context, reqDto *dtos.DeleteMyBlockPacksByIdsReqDto)
}

type BlockPackController struct {
	blockPackService services.BlockPackServiceInterface
}

func NewBlockPackController(service services.BlockPackServiceInterface) BlockPackControllerInterface {
	return &BlockPackController{
		blockPackService: service,
	}
}

/* ============================== Implementations ============================== */

func (c *BlockPackController) GetMyBlockPackById(ctx *gin.Context, reqDto *dtos.GetMyBlockPackByIdReqDto) {
	resDto, exception := c.blockPackService.GetMyBlockPackById(ctx.Request.Context(), reqDto)
	if exception != nil {
		exception.Log().SafelyResponseWithJSON(ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      resDto,
		"exception": nil,
	})
}

func (c *BlockPackController) GetMyBlockPackAndItsParentById(ctx *gin.Context, reqDto *dtos.GetMyBlockPackAndItsParentByIdReqDto) {
	resDto, exception := c.blockPackService.GetMyBlockPackAndItsParentById(ctx.Request.Context(), reqDto)
	if exception != nil {
		exception.Log().SafelyResponseWithJSON(ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      resDto,
		"exception": nil,
	})
}

func (c *BlockPackController) GetAllMyBlockPacksByParentSubShelfId(ctx *gin.Context, reqDto *dtos.GetAllMyBlockPacksByParentSubShelfIdReqDto) {
	resDto, exception := c.blockPackService.GetAllMyBlockPacksByParentSubShelfId(ctx.Request.Context(), reqDto)
	if exception != nil {
		exception.Log().SafelyResponseWithJSON(ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      resDto,
		"exception": nil,
	})
}

func (c *BlockPackController) GetAllMyBlockPacksByRootShelfId(ctx *gin.Context, reqDto *dtos.GetAllMyBlockPacksByRootShelfIdReqDto) {
	resDto, exception := c.blockPackService.GetAllMyBlockPacksByRootShelfId(ctx.Request.Context(), reqDto)
	if exception != nil {
		exception.Log().SafelyResponseWithJSON(ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      resDto,
		"exception": nil,
	})
}

func (c *BlockPackController) CreateBlockPack(ctx *gin.Context, reqDto *dtos.CreateBlockPackReqDto) {
	resDto, exception := c.blockPackService.CreateBlockPack(ctx.Request.Context(), reqDto)
	if exception != nil {
		exception.Log().SafelyResponseWithJSON(ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      resDto,
		"exception": nil,
	})
}

func (c *BlockPackController) UpdateMyBlockPackById(ctx *gin.Context, reqDto *dtos.UpdateMyBlockPackByIdReqDto) {
	resDto, exception := c.blockPackService.UpdateMyBlockPackById(ctx.Request.Context(), reqDto)
	if exception != nil {
		exception.Log().SafelyResponseWithJSON(ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      resDto,
		"exception": nil,
	})
}

func (c *BlockPackController) MoveMyBlockPackById(ctx *gin.Context, reqDto *dtos.MoveMyBlockPackByIdReqDto) {
	resDto, exception := c.blockPackService.MoveMyBlockPackById(ctx.Request.Context(), reqDto)
	if exception != nil {
		exception.Log().SafelyResponseWithJSON(ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      resDto,
		"exception": nil,
	})
}

func (c *BlockPackController) MoveMyBlockPacksByIds(ctx *gin.Context, reqDto *dtos.MoveMyBlockPacksByIdsReqDto) {
	resDto, exception := c.blockPackService.MoveMyBlockPacksByIds(ctx.Request.Context(), reqDto)
	if exception != nil {
		exception.Log().SafelyResponseWithJSON(ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      resDto,
		"exception": nil,
	})
}

func (c *BlockPackController) RestoreMyBlockPackById(ctx *gin.Context, reqDto *dtos.RestoreMyBlockPackByIdReqDto) {
	resDto, exception := c.blockPackService.RestoreMyBlockPackById(ctx.Request.Context(), reqDto)
	if exception != nil {
		exception.Log().SafelyResponseWithJSON(ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      resDto,
		"exception": nil,
	})
}

func (c *BlockPackController) RestoreMyBlockPacksByIds(ctx *gin.Context, reqDto *dtos.RestoreMyBlockPacksByIdsReqDto) {
	resDto, exception := c.blockPackService.RestoreMyBlockPacksByIds(ctx.Request.Context(), reqDto)
	if exception != nil {
		exception.Log().SafelyResponseWithJSON(ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      resDto,
		"exception": nil,
	})
}

func (c *BlockPackController) DeleteMyBlockPackById(ctx *gin.Context, reqDto *dtos.DeleteMyBlockPackByIdReqDto) {
	resDto, exception := c.blockPackService.DeleteMyBlockPackById(ctx.Request.Context(), reqDto)
	if exception != nil {
		exception.Log().SafelyResponseWithJSON(ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      resDto,
		"exception": nil,
	})
}

func (c *BlockPackController) DeleteMyBlockPacksByIds(ctx *gin.Context, reqDto *dtos.DeleteMyBlockPacksByIdsReqDto) {
	resDto, exception := c.blockPackService.DeleteMyBlockPacksByIds(ctx.Request.Context(), reqDto)
	if exception != nil {
		exception.Log().SafelyResponseWithJSON(ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      resDto,
		"exception": nil,
	})
}
