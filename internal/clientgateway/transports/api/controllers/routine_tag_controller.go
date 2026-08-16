package controllers

import (
	"github.com/gin-gonic/gin"

	exceptionwriter "github.com/HiIamJeff67/notegic-backend/shared/util/exceptionwriter"

	apicontract "github.com/HiIamJeff67/notegic-backend/contracts/core/v1/api/routine-tags"

	coreadapters "github.com/HiIamJeff67/notegic-backend/internal/clientgateway/transports/core/adapters"
)

type RoutineTagControllerInterface interface {
	GetMyRoutineTagById(ctx *gin.Context, requestDto *apicontract.GetMyRoutineTagByIdRequestDto)
	GetAllMyRoutineTags(ctx *gin.Context, requestDto *apicontract.GetAllMyRoutineTagsRequestDto)
	CreateRoutineTag(ctx *gin.Context, requestDto *apicontract.CreateRoutineTagRequestDto)
	CreateRoutineTags(ctx *gin.Context, requestDto *apicontract.CreateRoutineTagsRequestDto)
	UpdateMyRoutineTagById(ctx *gin.Context, requestDto *apicontract.UpdateMyRoutineTagByIdRequestDto)
	UpdateMyRoutineTagsByIds(ctx *gin.Context, requestDto *apicontract.UpdateMyRoutineTagsByIdsRequestDto)
	HardDeleteMyRoutineTagById(ctx *gin.Context, requestDto *apicontract.HardDeleteMyRoutineTagByIdRequestDto)
	HardDeleteMyRoutineTagsByIds(ctx *gin.Context, requestDto *apicontract.HardDeleteMyRoutineTagsByIdsRequestDto)
}

type RoutineTagController struct {
	coreAdapter *coreadapters.CoreAdapter
}

func NewRoutineTagController(coreAdapter *coreadapters.CoreAdapter) RoutineTagControllerInterface {
	return &RoutineTagController{
		coreAdapter: coreAdapter,
	}
}

func (c *RoutineTagController) GetMyRoutineTagById(ctx *gin.Context, requestDto *apicontract.GetMyRoutineTagByIdRequestDto) {
	response, exception := coreadapters.CallSecurly[
		apicontract.GetMyRoutineTagByIdRequestDto,
		apicontract.GetMyRoutineTagByIdResponseDto,
	](
		ctx,
		c.coreAdapter,
		requestDto,
		apicontract.GetMyRoutineTagByIdOperation,
		"/core/v1/routine-tags/get-by-id",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *RoutineTagController) GetAllMyRoutineTags(ctx *gin.Context, requestDto *apicontract.GetAllMyRoutineTagsRequestDto) {
	response, exception := coreadapters.CallSecurly[
		apicontract.GetAllMyRoutineTagsRequestDto,
		apicontract.GetAllMyRoutineTagsResponseDto,
	](
		ctx,
		c.coreAdapter,
		requestDto,
		apicontract.GetAllMyRoutineTagsOperation,
		"/core/v1/routine-tags/get-all",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *RoutineTagController) CreateRoutineTag(ctx *gin.Context, requestDto *apicontract.CreateRoutineTagRequestDto) {
	response, exception := coreadapters.CallSecurly[
		apicontract.CreateRoutineTagRequestDto,
		apicontract.CreateRoutineTagResponseDto,
	](
		ctx,
		c.coreAdapter,
		requestDto,
		apicontract.CreateRoutineTagOperation,
		"/core/v1/routine-tags/create",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeCreatedClientResponse(ctx, response.Data)
}

func (c *RoutineTagController) CreateRoutineTags(ctx *gin.Context, requestDto *apicontract.CreateRoutineTagsRequestDto) {
	response, exception := coreadapters.CallSecurly[
		apicontract.CreateRoutineTagsRequestDto,
		apicontract.CreateRoutineTagsResponseDto,
	](
		ctx,
		c.coreAdapter,
		requestDto,
		apicontract.CreateRoutineTagsOperation,
		"/core/v1/routine-tags/create-many",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeCreatedClientResponse(ctx, response.Data)
}

func (c *RoutineTagController) UpdateMyRoutineTagById(ctx *gin.Context, requestDto *apicontract.UpdateMyRoutineTagByIdRequestDto) {
	response, exception := coreadapters.CallSecurly[
		apicontract.UpdateMyRoutineTagByIdRequestDto,
		apicontract.UpdateMyRoutineTagByIdResponseDto,
	](
		ctx,
		c.coreAdapter,
		requestDto,
		apicontract.UpdateMyRoutineTagByIdOperation,
		"/core/v1/routine-tags/update",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *RoutineTagController) UpdateMyRoutineTagsByIds(ctx *gin.Context, requestDto *apicontract.UpdateMyRoutineTagsByIdsRequestDto) {
	response, exception := coreadapters.CallSecurly[
		apicontract.UpdateMyRoutineTagsByIdsRequestDto,
		apicontract.UpdateMyRoutineTagsByIdsResponseDto,
	](
		ctx,
		c.coreAdapter,
		requestDto,
		apicontract.UpdateMyRoutineTagsByIdsOperation,
		"/core/v1/routine-tags/update-many",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *RoutineTagController) HardDeleteMyRoutineTagById(ctx *gin.Context, requestDto *apicontract.HardDeleteMyRoutineTagByIdRequestDto) {
	response, exception := coreadapters.CallSecurly[
		apicontract.HardDeleteMyRoutineTagByIdRequestDto,
		apicontract.HardDeleteMyRoutineTagByIdResponseDto,
	](
		ctx,
		c.coreAdapter,
		requestDto,
		apicontract.HardDeleteMyRoutineTagByIdOperation,
		"/core/v1/routine-tags/hard-delete",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *RoutineTagController) HardDeleteMyRoutineTagsByIds(ctx *gin.Context, requestDto *apicontract.HardDeleteMyRoutineTagsByIdsRequestDto) {
	response, exception := coreadapters.CallSecurly[
		apicontract.HardDeleteMyRoutineTagsByIdsRequestDto,
		apicontract.HardDeleteMyRoutineTagsByIdsResponseDto,
	](
		ctx,
		c.coreAdapter,
		requestDto,
		apicontract.HardDeleteMyRoutineTagsByIdsOperation,
		"/core/v1/routine-tags/hard-delete-many",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}
