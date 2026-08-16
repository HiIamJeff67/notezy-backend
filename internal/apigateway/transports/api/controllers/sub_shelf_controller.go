package controllers

import (
	"github.com/gin-gonic/gin"

	exceptionwriter "github.com/HiIamJeff67/notegic-backend/shared/util/exceptionwriter"

	apicontract "github.com/HiIamJeff67/notegic-backend/contracts/core/v1/api/sub-shelves"

	coreadapters "github.com/HiIamJeff67/notegic-backend/internal/apigateway/transports/core/adapters"
)

type SubShelfControllerInterface interface {
	GetMySubShelfById(ctx *gin.Context, requestDto *apicontract.GetMySubShelfByIdRequestDto)
	GetMySubShelvesByPrevSubShelfId(ctx *gin.Context, requestDto *apicontract.GetMySubShelvesByPrevSubShelfIdRequestDto)
	GetAllMySubShelvesByRootShelfId(ctx *gin.Context, requestDto *apicontract.GetAllMySubShelvesByRootShelfIdRequestDto)
	GetMySubShelvesAndItemsByPrevSubShelfId(ctx *gin.Context, requestDto *apicontract.GetMySubShelvesAndItemsByPrevSubShelfIdRequestDto)
	CreateSubShelfByRootShelfId(ctx *gin.Context, requestDto *apicontract.CreateSubShelfByRootShelfIdRequestDto)
	CreateSubShelvesByRootShelfIds(ctx *gin.Context, requestDto *apicontract.CreateSubShelvesByRootShelfIdsRequestDto)
	UpdateMySubShelfById(ctx *gin.Context, requestDto *apicontract.UpdateMySubShelfByIdRequestDto)
	UpdateMySubShelvesByIds(ctx *gin.Context, requestDto *apicontract.UpdateMySubShelvesByIdsRequestDto)
	MoveMySubShelfByRootShelfId(ctx *gin.Context, requestDto *apicontract.MoveMySubShelfByRootShelfIdRequestDto)
	MoveMySubShelvesByRootShelfId(ctx *gin.Context, requestDto *apicontract.MoveMySubShelvesByRootShelfIdRequestDto)
	MoveMySubShelvesByRootShelfIds(ctx *gin.Context, requestDto *apicontract.MoveMySubShelvesByRootShelfIdsRequestDto)
	RestoreMySubShelfById(ctx *gin.Context, requestDto *apicontract.RestoreMySubShelfByIdRequestDto)
	RestoreMySubShelvesByIds(ctx *gin.Context, requestDto *apicontract.RestoreMySubShelvesByIdsRequestDto)
	DeleteMySubShelfById(ctx *gin.Context, requestDto *apicontract.DeleteMySubShelfByIdRequestDto)
	DeleteMySubShelvesByIds(ctx *gin.Context, requestDto *apicontract.DeleteMySubShelvesByIdsRequestDto)
}

type SubShelfController struct {
	coreAdapter *coreadapters.CoreAdapter
}

func NewSubShelfController(coreAdapter *coreadapters.CoreAdapter) SubShelfControllerInterface {
	return &SubShelfController{
		coreAdapter: coreAdapter,
	}
}

func (c *SubShelfController) GetMySubShelfById(ctx *gin.Context, requestDto *apicontract.GetMySubShelfByIdRequestDto) {
	response, exception := coreadapters.CallSecurly[
		apicontract.GetMySubShelfByIdRequestDto,
		apicontract.GetMySubShelfByIdResponseDto,
	](
		ctx,
		c.coreAdapter,
		requestDto,
		apicontract.GetMySubShelfByIdOperation,
		"/core/v1/sub-shelves/get-by-id",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *SubShelfController) GetMySubShelvesByPrevSubShelfId(ctx *gin.Context, requestDto *apicontract.GetMySubShelvesByPrevSubShelfIdRequestDto) {
	response, exception := coreadapters.CallSecurly[
		apicontract.GetMySubShelvesByPrevSubShelfIdRequestDto,
		apicontract.GetMySubShelvesByPrevSubShelfIdResponseDto,
	](
		ctx,
		c.coreAdapter,
		requestDto,
		apicontract.GetMySubShelvesByPrevSubShelfIdOperation,
		"/core/v1/sub-shelves/get-by-prev-sub-shelf-id",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *SubShelfController) GetAllMySubShelvesByRootShelfId(ctx *gin.Context, requestDto *apicontract.GetAllMySubShelvesByRootShelfIdRequestDto) {
	response, exception := coreadapters.CallSecurly[
		apicontract.GetAllMySubShelvesByRootShelfIdRequestDto,
		apicontract.GetAllMySubShelvesByRootShelfIdResponseDto,
	](
		ctx,
		c.coreAdapter,
		requestDto,
		apicontract.GetAllMySubShelvesByRootShelfIdOperation,
		"/core/v1/sub-shelves/get-all-by-root-shelf-id",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *SubShelfController) GetMySubShelvesAndItemsByPrevSubShelfId(ctx *gin.Context, requestDto *apicontract.GetMySubShelvesAndItemsByPrevSubShelfIdRequestDto) {
	response, exception := coreadapters.CallSecurly[
		apicontract.GetMySubShelvesAndItemsByPrevSubShelfIdRequestDto,
		apicontract.GetMySubShelvesAndItemsByPrevSubShelfIdResponseDto,
	](
		ctx,
		c.coreAdapter,
		requestDto,
		apicontract.GetMySubShelvesAndItemsByPrevSubShelfIdOperation,
		"/core/v1/sub-shelves/get-and-items-by-prev-sub-shelf-id",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *SubShelfController) CreateSubShelfByRootShelfId(ctx *gin.Context, requestDto *apicontract.CreateSubShelfByRootShelfIdRequestDto) {
	response, exception := coreadapters.CallSecurly[
		apicontract.CreateSubShelfByRootShelfIdRequestDto,
		apicontract.CreateSubShelfByRootShelfIdResponseDto,
	](
		ctx,
		c.coreAdapter,
		requestDto,
		apicontract.CreateSubShelfByRootShelfIdOperation,
		"/core/v1/sub-shelves/create",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeCreatedClientResponse(ctx, response.Data)
}

func (c *SubShelfController) CreateSubShelvesByRootShelfIds(ctx *gin.Context, requestDto *apicontract.CreateSubShelvesByRootShelfIdsRequestDto) {
	response, exception := coreadapters.CallSecurly[
		apicontract.CreateSubShelvesByRootShelfIdsRequestDto,
		apicontract.CreateSubShelvesByRootShelfIdsResponseDto,
	](
		ctx,
		c.coreAdapter,
		requestDto,
		apicontract.CreateSubShelvesByRootShelfIdsOperation,
		"/core/v1/sub-shelves/create-many",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeCreatedClientResponse(ctx, response.Data)
}

func (c *SubShelfController) UpdateMySubShelfById(ctx *gin.Context, requestDto *apicontract.UpdateMySubShelfByIdRequestDto) {
	response, exception := coreadapters.CallSecurly[
		apicontract.UpdateMySubShelfByIdRequestDto,
		apicontract.UpdateMySubShelfByIdResponseDto,
	](
		ctx,
		c.coreAdapter,
		requestDto,
		apicontract.UpdateMySubShelfByIdOperation,
		"/core/v1/sub-shelves/update",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *SubShelfController) UpdateMySubShelvesByIds(ctx *gin.Context, requestDto *apicontract.UpdateMySubShelvesByIdsRequestDto) {
	response, exception := coreadapters.CallSecurly[
		apicontract.UpdateMySubShelvesByIdsRequestDto,
		apicontract.UpdateMySubShelvesByIdsResponseDto,
	](
		ctx,
		c.coreAdapter,
		requestDto,
		apicontract.UpdateMySubShelvesByIdsOperation,
		"/core/v1/sub-shelves/update-many",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *SubShelfController) MoveMySubShelfByRootShelfId(ctx *gin.Context, requestDto *apicontract.MoveMySubShelfByRootShelfIdRequestDto) {
	response, exception := coreadapters.CallSecurly[
		apicontract.MoveMySubShelfByRootShelfIdRequestDto,
		apicontract.MoveMySubShelfByRootShelfIdResponseDto,
	](
		ctx,
		c.coreAdapter,
		requestDto,
		apicontract.MoveMySubShelfByRootShelfIdOperation,
		"/core/v1/sub-shelves/move",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *SubShelfController) MoveMySubShelvesByRootShelfId(ctx *gin.Context, requestDto *apicontract.MoveMySubShelvesByRootShelfIdRequestDto) {
	response, exception := coreadapters.CallSecurly[
		apicontract.MoveMySubShelvesByRootShelfIdRequestDto,
		apicontract.MoveMySubShelvesByRootShelfIdResponseDto,
	](
		ctx,
		c.coreAdapter,
		requestDto,
		apicontract.MoveMySubShelvesByRootShelfIdOperation,
		"/core/v1/sub-shelves/move-many",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *SubShelfController) MoveMySubShelvesByRootShelfIds(ctx *gin.Context, requestDto *apicontract.MoveMySubShelvesByRootShelfIdsRequestDto) {
	response, exception := coreadapters.CallSecurly[
		apicontract.MoveMySubShelvesByRootShelfIdsRequestDto,
		apicontract.MoveMySubShelvesByRootShelfIdsResponseDto,
	](
		ctx,
		c.coreAdapter,
		requestDto,
		apicontract.MoveMySubShelvesByRootShelfIdsOperation,
		"/core/v1/sub-shelves/move-many-by-root-shelves",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *SubShelfController) RestoreMySubShelfById(ctx *gin.Context, requestDto *apicontract.RestoreMySubShelfByIdRequestDto) {
	response, exception := coreadapters.CallSecurly[
		apicontract.RestoreMySubShelfByIdRequestDto,
		apicontract.RestoreMySubShelfByIdResponseDto,
	](
		ctx,
		c.coreAdapter,
		requestDto,
		apicontract.RestoreMySubShelfByIdOperation,
		"/core/v1/sub-shelves/restore",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *SubShelfController) RestoreMySubShelvesByIds(ctx *gin.Context, requestDto *apicontract.RestoreMySubShelvesByIdsRequestDto) {
	response, exception := coreadapters.CallSecurly[
		apicontract.RestoreMySubShelvesByIdsRequestDto,
		apicontract.RestoreMySubShelvesByIdsResponseDto,
	](
		ctx,
		c.coreAdapter,
		requestDto,
		apicontract.RestoreMySubShelvesByIdsOperation,
		"/core/v1/sub-shelves/restore-many",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *SubShelfController) DeleteMySubShelfById(ctx *gin.Context, requestDto *apicontract.DeleteMySubShelfByIdRequestDto) {
	response, exception := coreadapters.CallSecurly[
		apicontract.DeleteMySubShelfByIdRequestDto,
		apicontract.DeleteMySubShelfByIdResponseDto,
	](
		ctx,
		c.coreAdapter,
		requestDto,
		apicontract.DeleteMySubShelfByIdOperation,
		"/core/v1/sub-shelves/delete",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *SubShelfController) DeleteMySubShelvesByIds(ctx *gin.Context, requestDto *apicontract.DeleteMySubShelvesByIdsRequestDto) {
	response, exception := coreadapters.CallSecurly[
		apicontract.DeleteMySubShelvesByIdsRequestDto,
		apicontract.DeleteMySubShelvesByIdsResponseDto,
	](
		ctx,
		c.coreAdapter,
		requestDto,
		apicontract.DeleteMySubShelvesByIdsOperation,
		"/core/v1/sub-shelves/delete-many",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}
