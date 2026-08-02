package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	blockpacksdto "github.com/HiIamJeff67/notezy-backend/contracts/gateway/v1/api/block-packs"
	responsewriter "github.com/HiIamJeff67/notezy-backend/internal/shared/responsewriter"
	coreadapters "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/core/adapters"
)

type BlockPackControllerInterface interface {
	GetMyBlockPackById(ctx *gin.Context, requestDto *blockpacksdto.GetMyBlockPackByIdRequestDto)
	GetMyBlockPackAndItsParentById(ctx *gin.Context, requestDto *blockpacksdto.GetMyBlockPackAndItsParentByIdRequestDto)
	GetMyBlockPacksByParentSubShelfId(ctx *gin.Context, requestDto *blockpacksdto.GetMyBlockPacksByParentSubShelfIdRequestDto)
	GetAllMyBlockPacksByRootShelfId(ctx *gin.Context, requestDto *blockpacksdto.GetAllMyBlockPacksByRootShelfIdRequestDto)
	CreateBlockPack(ctx *gin.Context, requestDto *blockpacksdto.CreateBlockPackRequestDto)
	CreateBlockPacks(ctx *gin.Context, requestDto *blockpacksdto.CreateBlockPacksRequestDto)
	UpdateMyBlockPackById(ctx *gin.Context, requestDto *blockpacksdto.UpdateMyBlockPackByIdRequestDto)
	UpdateMyBlockPacksByIds(ctx *gin.Context, requestDto *blockpacksdto.UpdateMyBlockPacksByIdsRequestDto)
	MoveMyBlockPackByParentSubShelfId(ctx *gin.Context, requestDto *blockpacksdto.MoveMyBlockPackByParentSubShelfIdRequestDto)
	MoveMyBlockPacksByParentSubShelfId(ctx *gin.Context, requestDto *blockpacksdto.MoveMyBlockPacksByParentSubShelfIdRequestDto)
	MoveMyBlockPacksByParentSubShelfIds(ctx *gin.Context, requestDto *blockpacksdto.MoveMyBlockPacksByParentSubShelfIdsRequestDto)
	RestoreMyBlockPackById(ctx *gin.Context, requestDto *blockpacksdto.RestoreMyBlockPackByIdRequestDto)
	RestoreMyBlockPacksByIds(ctx *gin.Context, requestDto *blockpacksdto.RestoreMyBlockPacksByIdsRequestDto)
	DeleteMyBlockPackById(ctx *gin.Context, requestDto *blockpacksdto.DeleteMyBlockPackByIdRequestDto)
	DeleteMyBlockPacksByIds(ctx *gin.Context, requestDto *blockpacksdto.DeleteMyBlockPacksByIdsRequestDto)
}

type BlockPackController struct {
	coreClient *coreadapters.CoreClient
}

func NewBlockPackController(coreClient *coreadapters.CoreClient) BlockPackControllerInterface {
	return &BlockPackController{
		coreClient: coreClient,
	}
}

func (c *BlockPackController) GetMyBlockPackById(ctx *gin.Context, requestDto *blockpacksdto.GetMyBlockPackByIdRequestDto) {
	response, exception := coreadapters.CallSecurly[
		blockpacksdto.GetMyBlockPackByIdRequestDto,
		blockpacksdto.GetMyBlockPackByIdResponseDto,
	](
		ctx,
		c.coreClient,
		requestDto,
		blockpacksdto.GetMyBlockPackByIdOperation,
		"/core/v1/block-packs/get-by-id",
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

func (c *BlockPackController) GetMyBlockPackAndItsParentById(ctx *gin.Context, requestDto *blockpacksdto.GetMyBlockPackAndItsParentByIdRequestDto) {
	response, exception := coreadapters.CallSecurly[
		blockpacksdto.GetMyBlockPackAndItsParentByIdRequestDto,
		blockpacksdto.GetMyBlockPackAndItsParentByIdResponseDto,
	](
		ctx,
		c.coreClient,
		requestDto,
		blockpacksdto.GetMyBlockPackAndItsParentByIdOperation,
		"/core/v1/block-packs/get-and-parent-by-id",
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

func (c *BlockPackController) GetMyBlockPacksByParentSubShelfId(ctx *gin.Context, requestDto *blockpacksdto.GetMyBlockPacksByParentSubShelfIdRequestDto) {
	response, exception := coreadapters.CallSecurly[
		blockpacksdto.GetMyBlockPacksByParentSubShelfIdRequestDto,
		blockpacksdto.GetMyBlockPacksByParentSubShelfIdResponseDto,
	](
		ctx,
		c.coreClient,
		requestDto,
		blockpacksdto.GetMyBlockPacksByParentSubShelfIdOperation,
		"/core/v1/block-packs/get-by-parent-sub-shelf-id",
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

func (c *BlockPackController) GetAllMyBlockPacksByRootShelfId(ctx *gin.Context, requestDto *blockpacksdto.GetAllMyBlockPacksByRootShelfIdRequestDto) {
	response, exception := coreadapters.CallSecurly[
		blockpacksdto.GetAllMyBlockPacksByRootShelfIdRequestDto,
		blockpacksdto.GetAllMyBlockPacksByRootShelfIdResponseDto,
	](
		ctx,
		c.coreClient,
		requestDto,
		blockpacksdto.GetAllMyBlockPacksByRootShelfIdOperation,
		"/core/v1/block-packs/get-all-by-root-shelf-id",
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

func (c *BlockPackController) CreateBlockPack(ctx *gin.Context, requestDto *blockpacksdto.CreateBlockPackRequestDto) {
	response, exception := coreadapters.CallSecurly[
		blockpacksdto.CreateBlockPackRequestDto,
		blockpacksdto.CreateBlockPackResponseDto,
	](
		ctx,
		c.coreClient,
		requestDto,
		blockpacksdto.CreateBlockPackOperation,
		"/core/v1/block-packs/create",
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

func (c *BlockPackController) CreateBlockPacks(ctx *gin.Context, requestDto *blockpacksdto.CreateBlockPacksRequestDto) {
	response, exception := coreadapters.CallSecurly[
		blockpacksdto.CreateBlockPacksRequestDto,
		blockpacksdto.CreateBlockPacksResponseDto,
	](
		ctx,
		c.coreClient,
		requestDto,
		blockpacksdto.CreateBlockPacksOperation,
		"/core/v1/block-packs/create-many",
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

func (c *BlockPackController) UpdateMyBlockPackById(ctx *gin.Context, requestDto *blockpacksdto.UpdateMyBlockPackByIdRequestDto) {
	response, exception := coreadapters.CallSecurly[
		blockpacksdto.UpdateMyBlockPackByIdRequestDto,
		blockpacksdto.UpdateMyBlockPackByIdResponseDto,
	](
		ctx,
		c.coreClient,
		requestDto,
		blockpacksdto.UpdateMyBlockPackByIdOperation,
		"/core/v1/block-packs/update",
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

func (c *BlockPackController) UpdateMyBlockPacksByIds(ctx *gin.Context, requestDto *blockpacksdto.UpdateMyBlockPacksByIdsRequestDto) {
	response, exception := coreadapters.CallSecurly[
		blockpacksdto.UpdateMyBlockPacksByIdsRequestDto,
		blockpacksdto.UpdateMyBlockPacksByIdsResponseDto,
	](
		ctx,
		c.coreClient,
		requestDto,
		blockpacksdto.UpdateMyBlockPacksByIdsOperation,
		"/core/v1/block-packs/update-many",
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

func (c *BlockPackController) MoveMyBlockPackByParentSubShelfId(ctx *gin.Context, requestDto *blockpacksdto.MoveMyBlockPackByParentSubShelfIdRequestDto) {
	response, exception := coreadapters.CallSecurly[
		blockpacksdto.MoveMyBlockPackByParentSubShelfIdRequestDto,
		blockpacksdto.MoveMyBlockPackByParentSubShelfIdResponseDto,
	](
		ctx,
		c.coreClient,
		requestDto,
		blockpacksdto.MoveMyBlockPackByParentSubShelfIdOperation,
		"/core/v1/block-packs/move",
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

func (c *BlockPackController) MoveMyBlockPacksByParentSubShelfId(ctx *gin.Context, requestDto *blockpacksdto.MoveMyBlockPacksByParentSubShelfIdRequestDto) {
	response, exception := coreadapters.CallSecurly[
		blockpacksdto.MoveMyBlockPacksByParentSubShelfIdRequestDto,
		blockpacksdto.MoveMyBlockPacksByParentSubShelfIdResponseDto,
	](
		ctx,
		c.coreClient,
		requestDto,
		blockpacksdto.MoveMyBlockPacksByParentSubShelfIdOperation,
		"/core/v1/block-packs/move-many",
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

func (c *BlockPackController) MoveMyBlockPacksByParentSubShelfIds(ctx *gin.Context, requestDto *blockpacksdto.MoveMyBlockPacksByParentSubShelfIdsRequestDto) {
	response, exception := coreadapters.CallSecurly[
		blockpacksdto.MoveMyBlockPacksByParentSubShelfIdsRequestDto,
		blockpacksdto.MoveMyBlockPacksByParentSubShelfIdsResponseDto,
	](
		ctx,
		c.coreClient,
		requestDto,
		blockpacksdto.MoveMyBlockPacksByParentSubShelfIdsOperation,
		"/core/v1/block-packs/move-many-by-parent-sub-shelves",
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

func (c *BlockPackController) RestoreMyBlockPackById(ctx *gin.Context, requestDto *blockpacksdto.RestoreMyBlockPackByIdRequestDto) {
	response, exception := coreadapters.CallSecurly[
		blockpacksdto.RestoreMyBlockPackByIdRequestDto,
		blockpacksdto.RestoreMyBlockPackByIdResponseDto,
	](
		ctx,
		c.coreClient,
		requestDto,
		blockpacksdto.RestoreMyBlockPackByIdOperation,
		"/core/v1/block-packs/restore",
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

func (c *BlockPackController) RestoreMyBlockPacksByIds(ctx *gin.Context, requestDto *blockpacksdto.RestoreMyBlockPacksByIdsRequestDto) {
	response, exception := coreadapters.CallSecurly[
		blockpacksdto.RestoreMyBlockPacksByIdsRequestDto,
		blockpacksdto.RestoreMyBlockPacksByIdsResponseDto,
	](
		ctx,
		c.coreClient,
		requestDto,
		blockpacksdto.RestoreMyBlockPacksByIdsOperation,
		"/core/v1/block-packs/restore-many",
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

func (c *BlockPackController) DeleteMyBlockPackById(ctx *gin.Context, requestDto *blockpacksdto.DeleteMyBlockPackByIdRequestDto) {
	response, exception := coreadapters.CallSecurly[
		blockpacksdto.DeleteMyBlockPackByIdRequestDto,
		blockpacksdto.DeleteMyBlockPackByIdResponseDto,
	](
		ctx,
		c.coreClient,
		requestDto,
		blockpacksdto.DeleteMyBlockPackByIdOperation,
		"/core/v1/block-packs/delete",
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

func (c *BlockPackController) DeleteMyBlockPacksByIds(ctx *gin.Context, requestDto *blockpacksdto.DeleteMyBlockPacksByIdsRequestDto) {
	response, exception := coreadapters.CallSecurly[
		blockpacksdto.DeleteMyBlockPacksByIdsRequestDto,
		blockpacksdto.DeleteMyBlockPacksByIdsResponseDto,
	](
		ctx,
		c.coreClient,
		requestDto,
		blockpacksdto.DeleteMyBlockPacksByIdsOperation,
		"/core/v1/block-packs/delete-many",
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
