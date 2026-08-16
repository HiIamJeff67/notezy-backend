package controllers

import (
	"github.com/gin-gonic/gin"

	exceptionwriter "github.com/HiIamJeff67/notegic-backend/shared/util/exceptionwriter"

	apicontract "github.com/HiIamJeff67/notegic-backend/contracts/core/v1/api/block-packs"

	coreadapters "github.com/HiIamJeff67/notegic-backend/internal/clientgateway/transports/core/adapters"
)

type BlockPackControllerInterface interface {
	GetMyBlockPackById(ctx *gin.Context, requestDto *apicontract.GetMyBlockPackByIdRequestDto)
	GetMyBlockPackAndItsParentById(ctx *gin.Context, requestDto *apicontract.GetMyBlockPackAndItsParentByIdRequestDto)
	GetMyBlockPacksByParentSubShelfId(ctx *gin.Context, requestDto *apicontract.GetMyBlockPacksByParentSubShelfIdRequestDto)
	GetAllMyBlockPacksByRootShelfId(ctx *gin.Context, requestDto *apicontract.GetAllMyBlockPacksByRootShelfIdRequestDto)
	CreateBlockPack(ctx *gin.Context, requestDto *apicontract.CreateBlockPackRequestDto)
	CreateBlockPacks(ctx *gin.Context, requestDto *apicontract.CreateBlockPacksRequestDto)
	UpdateMyBlockPackById(ctx *gin.Context, requestDto *apicontract.UpdateMyBlockPackByIdRequestDto)
	UpdateMyBlockPacksByIds(ctx *gin.Context, requestDto *apicontract.UpdateMyBlockPacksByIdsRequestDto)
	MoveMyBlockPackByParentSubShelfId(ctx *gin.Context, requestDto *apicontract.MoveMyBlockPackByParentSubShelfIdRequestDto)
	MoveMyBlockPacksByParentSubShelfId(ctx *gin.Context, requestDto *apicontract.MoveMyBlockPacksByParentSubShelfIdRequestDto)
	MoveMyBlockPacksByParentSubShelfIds(ctx *gin.Context, requestDto *apicontract.MoveMyBlockPacksByParentSubShelfIdsRequestDto)
	RestoreMyBlockPackById(ctx *gin.Context, requestDto *apicontract.RestoreMyBlockPackByIdRequestDto)
	RestoreMyBlockPacksByIds(ctx *gin.Context, requestDto *apicontract.RestoreMyBlockPacksByIdsRequestDto)
	DeleteMyBlockPackById(ctx *gin.Context, requestDto *apicontract.DeleteMyBlockPackByIdRequestDto)
	DeleteMyBlockPacksByIds(ctx *gin.Context, requestDto *apicontract.DeleteMyBlockPacksByIdsRequestDto)
}

type BlockPackController struct {
	coreAdapter *coreadapters.CoreAdapter
}

func NewBlockPackController(coreAdapter *coreadapters.CoreAdapter) BlockPackControllerInterface {
	return &BlockPackController{
		coreAdapter: coreAdapter,
	}
}

func (c *BlockPackController) GetMyBlockPackById(ctx *gin.Context, requestDto *apicontract.GetMyBlockPackByIdRequestDto) {
	response, exception := coreadapters.CallSecurly[
		apicontract.GetMyBlockPackByIdRequestDto,
		apicontract.GetMyBlockPackByIdResponseDto,
	](
		ctx,
		c.coreAdapter,
		requestDto,
		apicontract.GetMyBlockPackByIdOperation,
		"/core/v1/block-packs/get-by-id",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *BlockPackController) GetMyBlockPackAndItsParentById(ctx *gin.Context, requestDto *apicontract.GetMyBlockPackAndItsParentByIdRequestDto) {
	response, exception := coreadapters.CallSecurly[
		apicontract.GetMyBlockPackAndItsParentByIdRequestDto,
		apicontract.GetMyBlockPackAndItsParentByIdResponseDto,
	](
		ctx,
		c.coreAdapter,
		requestDto,
		apicontract.GetMyBlockPackAndItsParentByIdOperation,
		"/core/v1/block-packs/get-and-parent-by-id",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *BlockPackController) GetMyBlockPacksByParentSubShelfId(ctx *gin.Context, requestDto *apicontract.GetMyBlockPacksByParentSubShelfIdRequestDto) {
	response, exception := coreadapters.CallSecurly[
		apicontract.GetMyBlockPacksByParentSubShelfIdRequestDto,
		apicontract.GetMyBlockPacksByParentSubShelfIdResponseDto,
	](
		ctx,
		c.coreAdapter,
		requestDto,
		apicontract.GetMyBlockPacksByParentSubShelfIdOperation,
		"/core/v1/block-packs/get-by-parent-sub-shelf-id",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *BlockPackController) GetAllMyBlockPacksByRootShelfId(ctx *gin.Context, requestDto *apicontract.GetAllMyBlockPacksByRootShelfIdRequestDto) {
	response, exception := coreadapters.CallSecurly[
		apicontract.GetAllMyBlockPacksByRootShelfIdRequestDto,
		apicontract.GetAllMyBlockPacksByRootShelfIdResponseDto,
	](
		ctx,
		c.coreAdapter,
		requestDto,
		apicontract.GetAllMyBlockPacksByRootShelfIdOperation,
		"/core/v1/block-packs/get-all-by-root-shelf-id",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *BlockPackController) CreateBlockPack(ctx *gin.Context, requestDto *apicontract.CreateBlockPackRequestDto) {
	response, exception := coreadapters.CallSecurly[
		apicontract.CreateBlockPackRequestDto,
		apicontract.CreateBlockPackResponseDto,
	](
		ctx,
		c.coreAdapter,
		requestDto,
		apicontract.CreateBlockPackOperation,
		"/core/v1/block-packs/create",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeCreatedClientResponse(ctx, response.Data)
}

func (c *BlockPackController) CreateBlockPacks(ctx *gin.Context, requestDto *apicontract.CreateBlockPacksRequestDto) {
	response, exception := coreadapters.CallSecurly[
		apicontract.CreateBlockPacksRequestDto,
		apicontract.CreateBlockPacksResponseDto,
	](
		ctx,
		c.coreAdapter,
		requestDto,
		apicontract.CreateBlockPacksOperation,
		"/core/v1/block-packs/create-many",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeCreatedClientResponse(ctx, response.Data)
}

func (c *BlockPackController) UpdateMyBlockPackById(ctx *gin.Context, requestDto *apicontract.UpdateMyBlockPackByIdRequestDto) {
	response, exception := coreadapters.CallSecurly[
		apicontract.UpdateMyBlockPackByIdRequestDto,
		apicontract.UpdateMyBlockPackByIdResponseDto,
	](
		ctx,
		c.coreAdapter,
		requestDto,
		apicontract.UpdateMyBlockPackByIdOperation,
		"/core/v1/block-packs/update",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *BlockPackController) UpdateMyBlockPacksByIds(ctx *gin.Context, requestDto *apicontract.UpdateMyBlockPacksByIdsRequestDto) {
	response, exception := coreadapters.CallSecurly[
		apicontract.UpdateMyBlockPacksByIdsRequestDto,
		apicontract.UpdateMyBlockPacksByIdsResponseDto,
	](
		ctx,
		c.coreAdapter,
		requestDto,
		apicontract.UpdateMyBlockPacksByIdsOperation,
		"/core/v1/block-packs/update-many",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *BlockPackController) MoveMyBlockPackByParentSubShelfId(ctx *gin.Context, requestDto *apicontract.MoveMyBlockPackByParentSubShelfIdRequestDto) {
	response, exception := coreadapters.CallSecurly[
		apicontract.MoveMyBlockPackByParentSubShelfIdRequestDto,
		apicontract.MoveMyBlockPackByParentSubShelfIdResponseDto,
	](
		ctx,
		c.coreAdapter,
		requestDto,
		apicontract.MoveMyBlockPackByParentSubShelfIdOperation,
		"/core/v1/block-packs/move",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *BlockPackController) MoveMyBlockPacksByParentSubShelfId(ctx *gin.Context, requestDto *apicontract.MoveMyBlockPacksByParentSubShelfIdRequestDto) {
	response, exception := coreadapters.CallSecurly[
		apicontract.MoveMyBlockPacksByParentSubShelfIdRequestDto,
		apicontract.MoveMyBlockPacksByParentSubShelfIdResponseDto,
	](
		ctx,
		c.coreAdapter,
		requestDto,
		apicontract.MoveMyBlockPacksByParentSubShelfIdOperation,
		"/core/v1/block-packs/move-many",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *BlockPackController) MoveMyBlockPacksByParentSubShelfIds(ctx *gin.Context, requestDto *apicontract.MoveMyBlockPacksByParentSubShelfIdsRequestDto) {
	response, exception := coreadapters.CallSecurly[
		apicontract.MoveMyBlockPacksByParentSubShelfIdsRequestDto,
		apicontract.MoveMyBlockPacksByParentSubShelfIdsResponseDto,
	](
		ctx,
		c.coreAdapter,
		requestDto,
		apicontract.MoveMyBlockPacksByParentSubShelfIdsOperation,
		"/core/v1/block-packs/move-many-by-parent-sub-shelves",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *BlockPackController) RestoreMyBlockPackById(ctx *gin.Context, requestDto *apicontract.RestoreMyBlockPackByIdRequestDto) {
	response, exception := coreadapters.CallSecurly[
		apicontract.RestoreMyBlockPackByIdRequestDto,
		apicontract.RestoreMyBlockPackByIdResponseDto,
	](
		ctx,
		c.coreAdapter,
		requestDto,
		apicontract.RestoreMyBlockPackByIdOperation,
		"/core/v1/block-packs/restore",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *BlockPackController) RestoreMyBlockPacksByIds(ctx *gin.Context, requestDto *apicontract.RestoreMyBlockPacksByIdsRequestDto) {
	response, exception := coreadapters.CallSecurly[
		apicontract.RestoreMyBlockPacksByIdsRequestDto,
		apicontract.RestoreMyBlockPacksByIdsResponseDto,
	](
		ctx,
		c.coreAdapter,
		requestDto,
		apicontract.RestoreMyBlockPacksByIdsOperation,
		"/core/v1/block-packs/restore-many",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *BlockPackController) DeleteMyBlockPackById(ctx *gin.Context, requestDto *apicontract.DeleteMyBlockPackByIdRequestDto) {
	response, exception := coreadapters.CallSecurly[
		apicontract.DeleteMyBlockPackByIdRequestDto,
		apicontract.DeleteMyBlockPackByIdResponseDto,
	](
		ctx,
		c.coreAdapter,
		requestDto,
		apicontract.DeleteMyBlockPackByIdOperation,
		"/core/v1/block-packs/delete",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *BlockPackController) DeleteMyBlockPacksByIds(ctx *gin.Context, requestDto *apicontract.DeleteMyBlockPacksByIdsRequestDto) {
	response, exception := coreadapters.CallSecurly[
		apicontract.DeleteMyBlockPacksByIdsRequestDto,
		apicontract.DeleteMyBlockPacksByIdsResponseDto,
	](
		ctx,
		c.coreAdapter,
		requestDto,
		apicontract.DeleteMyBlockPacksByIdsOperation,
		"/core/v1/block-packs/delete-many",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}
