package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	routinetagsdto "github.com/HiIamJeff67/notezy-backend/contracts/api/v1/routine-tags"
	responsewriter "github.com/HiIamJeff67/notezy-backend/internal/gateway/responsewriter"
	coreadapters "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/core/adapters"
)

type RoutineTagControllerInterface interface {
	GetMyRoutineTagById(ctx *gin.Context, requestDto *routinetagsdto.GetMyRoutineTagByIdRequestDto)
	GetAllMyRoutineTags(ctx *gin.Context, requestDto *routinetagsdto.GetAllMyRoutineTagsRequestDto)
	CreateRoutineTag(ctx *gin.Context, requestDto *routinetagsdto.CreateRoutineTagRequestDto)
	CreateRoutineTags(ctx *gin.Context, requestDto *routinetagsdto.CreateRoutineTagsRequestDto)
	UpdateMyRoutineTagById(ctx *gin.Context, requestDto *routinetagsdto.UpdateMyRoutineTagByIdRequestDto)
	UpdateMyRoutineTagsByIds(ctx *gin.Context, requestDto *routinetagsdto.UpdateMyRoutineTagsByIdsRequestDto)
	HardDeleteMyRoutineTagById(ctx *gin.Context, requestDto *routinetagsdto.HardDeleteMyRoutineTagByIdRequestDto)
	HardDeleteMyRoutineTagsByIds(ctx *gin.Context, requestDto *routinetagsdto.HardDeleteMyRoutineTagsByIdsRequestDto)
}

type RoutineTagController struct {
	coreClient *coreadapters.CoreClient
}

func NewRoutineTagController(coreClient *coreadapters.CoreClient) RoutineTagControllerInterface {
	return &RoutineTagController{
		coreClient: coreClient,
	}
}

func (c *RoutineTagController) GetMyRoutineTagById(ctx *gin.Context, requestDto *routinetagsdto.GetMyRoutineTagByIdRequestDto) {
	response, exception := coreadapters.CallSecurly[
		routinetagsdto.GetMyRoutineTagByIdRequestDto,
		routinetagsdto.GetMyRoutineTagByIdResponseDto,
	](
		ctx,
		c.coreClient,
		requestDto,
		routinetagsdto.GetMyRoutineTagByIdOperation,
		"/core/v1/routine-tags/get-by-id",
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

func (c *RoutineTagController) GetAllMyRoutineTags(ctx *gin.Context, requestDto *routinetagsdto.GetAllMyRoutineTagsRequestDto) {
	response, exception := coreadapters.CallSecurly[
		routinetagsdto.GetAllMyRoutineTagsRequestDto,
		routinetagsdto.GetAllMyRoutineTagsResponseDto,
	](
		ctx,
		c.coreClient,
		requestDto,
		routinetagsdto.GetAllMyRoutineTagsOperation,
		"/core/v1/routine-tags/get-all",
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

func (c *RoutineTagController) CreateRoutineTag(ctx *gin.Context, requestDto *routinetagsdto.CreateRoutineTagRequestDto) {
	response, exception := coreadapters.CallSecurly[
		routinetagsdto.CreateRoutineTagRequestDto,
		routinetagsdto.CreateRoutineTagResponseDto,
	](
		ctx,
		c.coreClient,
		requestDto,
		routinetagsdto.CreateRoutineTagOperation,
		"/core/v1/routine-tags/create",
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

func (c *RoutineTagController) CreateRoutineTags(ctx *gin.Context, requestDto *routinetagsdto.CreateRoutineTagsRequestDto) {
	response, exception := coreadapters.CallSecurly[
		routinetagsdto.CreateRoutineTagsRequestDto,
		routinetagsdto.CreateRoutineTagsResponseDto,
	](
		ctx,
		c.coreClient,
		requestDto,
		routinetagsdto.CreateRoutineTagsOperation,
		"/core/v1/routine-tags/create-many",
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

func (c *RoutineTagController) UpdateMyRoutineTagById(ctx *gin.Context, requestDto *routinetagsdto.UpdateMyRoutineTagByIdRequestDto) {
	response, exception := coreadapters.CallSecurly[
		routinetagsdto.UpdateMyRoutineTagByIdRequestDto,
		routinetagsdto.UpdateMyRoutineTagByIdResponseDto,
	](
		ctx,
		c.coreClient,
		requestDto,
		routinetagsdto.UpdateMyRoutineTagByIdOperation,
		"/core/v1/routine-tags/update",
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

func (c *RoutineTagController) UpdateMyRoutineTagsByIds(ctx *gin.Context, requestDto *routinetagsdto.UpdateMyRoutineTagsByIdsRequestDto) {
	response, exception := coreadapters.CallSecurly[
		routinetagsdto.UpdateMyRoutineTagsByIdsRequestDto,
		routinetagsdto.UpdateMyRoutineTagsByIdsResponseDto,
	](
		ctx,
		c.coreClient,
		requestDto,
		routinetagsdto.UpdateMyRoutineTagsByIdsOperation,
		"/core/v1/routine-tags/update-many",
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

func (c *RoutineTagController) HardDeleteMyRoutineTagById(ctx *gin.Context, requestDto *routinetagsdto.HardDeleteMyRoutineTagByIdRequestDto) {
	response, exception := coreadapters.CallSecurly[
		routinetagsdto.HardDeleteMyRoutineTagByIdRequestDto,
		routinetagsdto.HardDeleteMyRoutineTagByIdResponseDto,
	](
		ctx,
		c.coreClient,
		requestDto,
		routinetagsdto.HardDeleteMyRoutineTagByIdOperation,
		"/core/v1/routine-tags/hard-delete",
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

func (c *RoutineTagController) HardDeleteMyRoutineTagsByIds(ctx *gin.Context, requestDto *routinetagsdto.HardDeleteMyRoutineTagsByIdsRequestDto) {
	response, exception := coreadapters.CallSecurly[
		routinetagsdto.HardDeleteMyRoutineTagsByIdsRequestDto,
		routinetagsdto.HardDeleteMyRoutineTagsByIdsResponseDto,
	](
		ctx,
		c.coreClient,
		requestDto,
		routinetagsdto.HardDeleteMyRoutineTagsByIdsOperation,
		"/core/v1/routine-tags/hard-delete-many",
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
