package controllers

import (
	"github.com/gin-gonic/gin"

	exceptionwriter "github.com/HiIamJeff67/notezy-backend/shared/util/exceptionwriter"

	apicontract "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api/routine-tasks"

	coreadapters "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/core/adapters"
)

type RoutineTaskControllerInterface interface {
	GetMyRoutineTaskById(ctx *gin.Context, requestDto *apicontract.GetMyRoutineTaskByIdRequestDto)
	GetAllMyRoutineTasksByRoutineIds(ctx *gin.Context, requestDto *apicontract.GetAllMyRoutineTasksByRoutineIdsRequestDto)
	GetAllMyRoutineTasks(ctx *gin.Context, requestDto *apicontract.GetAllMyRoutineTasksRequestDto)
	CreateRoutineTaskByRoutineId(ctx *gin.Context, requestDto *apicontract.CreateRoutineTaskByRoutineIdRequestDto)
	UpdateMyRoutineTaskById(ctx *gin.Context, requestDto *apicontract.UpdateMyRoutineTaskByIdRequestDto)
	PauseMyRoutineTaskById(ctx *gin.Context, requestDto *apicontract.PauseMyRoutineTaskByIdRequestDto)
	ResumeMyRoutineTaskById(ctx *gin.Context, requestDto *apicontract.ResumeMyRoutineTaskByIdRequestDto)
	HardDeleteMyRoutineTaskById(ctx *gin.Context, requestDto *apicontract.HardDeleteMyRoutineTaskByIdRequestDto)
	HardDeleteMyRoutineTasksByIds(ctx *gin.Context, requestDto *apicontract.HardDeleteMyRoutineTasksByIdsRequestDto)

	/* ============================== Visualization Methods ============================== */
	VisualizeMyRoutineTaskStatusCount(ctx *gin.Context, requestDto *apicontract.VisualizeMyRoutineTaskStatusCountRequestDto)
	VisualizeMyRoutineTaskPurposeCount(ctx *gin.Context, requestDto *apicontract.VisualizeMyRoutineTaskPurposeCountRequestDto)
	VisualizeMyRoutineTaskScheduledAtCount(ctx *gin.Context, requestDto *apicontract.VisualizeMyRoutineTaskScheduledAtCountRequestDto)
	VisualizeMyRoutineTaskActualStartedAtCount(ctx *gin.Context, requestDto *apicontract.VisualizeMyRoutineTaskActualStartedAtCountRequestDto)
	VisualizeMyRoutineTaskActualEndedAtCount(ctx *gin.Context, requestDto *apicontract.VisualizeMyRoutineTaskActualEndedAtCountRequestDto)
}

type RoutineTaskController struct {
	coreClient *coreadapters.CoreAdapter
}

func NewRoutineTaskController(coreClient *coreadapters.CoreAdapter) RoutineTaskControllerInterface {
	return &RoutineTaskController{coreClient: coreClient}
}

func (c *RoutineTaskController) GetMyRoutineTaskById(ctx *gin.Context, requestDto *apicontract.GetMyRoutineTaskByIdRequestDto) {
	response, exception := coreadapters.CallSecurly[apicontract.GetMyRoutineTaskByIdRequestDto, apicontract.GetMyRoutineTaskByIdResponseDto](
		ctx,
		c.coreClient,
		requestDto,
		apicontract.GetMyRoutineTaskByIdOperation,
		"/core/v1/routine-tasks/get-by-id",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *RoutineTaskController) GetAllMyRoutineTasksByRoutineIds(ctx *gin.Context, requestDto *apicontract.GetAllMyRoutineTasksByRoutineIdsRequestDto) {
	response, exception := coreadapters.CallSecurly[apicontract.GetAllMyRoutineTasksByRoutineIdsRequestDto, apicontract.GetAllMyRoutineTasksByRoutineIdsResponseDto](
		ctx,
		c.coreClient,
		requestDto,
		apicontract.GetAllMyRoutineTasksByRoutineIdsOperation,
		"/core/v1/routine-tasks/get-all-by-routine-ids",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *RoutineTaskController) GetAllMyRoutineTasks(ctx *gin.Context, requestDto *apicontract.GetAllMyRoutineTasksRequestDto) {
	response, exception := coreadapters.CallSecurly[apicontract.GetAllMyRoutineTasksRequestDto, apicontract.GetAllMyRoutineTasksResponseDto](
		ctx,
		c.coreClient,
		requestDto,
		apicontract.GetAllMyRoutineTasksOperation,
		"/core/v1/routine-tasks/get-all",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *RoutineTaskController) CreateRoutineTaskByRoutineId(ctx *gin.Context, requestDto *apicontract.CreateRoutineTaskByRoutineIdRequestDto) {
	response, exception := coreadapters.CallSecurly[apicontract.CreateRoutineTaskByRoutineIdRequestDto, apicontract.CreateRoutineTaskByRoutineIdResponseDto](
		ctx,
		c.coreClient,
		requestDto,
		apicontract.CreateRoutineTaskByRoutineIdOperation,
		"/core/v1/routine-tasks/create-by-routine-id",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeCreatedClientResponse(ctx, response.Data)
}

func (c *RoutineTaskController) UpdateMyRoutineTaskById(ctx *gin.Context, requestDto *apicontract.UpdateMyRoutineTaskByIdRequestDto) {
	response, exception := coreadapters.CallSecurly[apicontract.UpdateMyRoutineTaskByIdRequestDto, apicontract.UpdateMyRoutineTaskByIdResponseDto](
		ctx,
		c.coreClient,
		requestDto,
		apicontract.UpdateMyRoutineTaskByIdOperation,
		"/core/v1/routine-tasks/update",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *RoutineTaskController) PauseMyRoutineTaskById(ctx *gin.Context, requestDto *apicontract.PauseMyRoutineTaskByIdRequestDto) {
	response, exception := coreadapters.CallSecurly[apicontract.PauseMyRoutineTaskByIdRequestDto, apicontract.PauseMyRoutineTaskByIdResponseDto](
		ctx,
		c.coreClient,
		requestDto,
		apicontract.PauseMyRoutineTaskByIdOperation,
		"/core/v1/routine-tasks/pause",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *RoutineTaskController) ResumeMyRoutineTaskById(ctx *gin.Context, requestDto *apicontract.ResumeMyRoutineTaskByIdRequestDto) {
	response, exception := coreadapters.CallSecurly[apicontract.ResumeMyRoutineTaskByIdRequestDto, apicontract.ResumeMyRoutineTaskByIdResponseDto](
		ctx,
		c.coreClient,
		requestDto,
		apicontract.ResumeMyRoutineTaskByIdOperation,
		"/core/v1/routine-tasks/resume",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *RoutineTaskController) HardDeleteMyRoutineTaskById(ctx *gin.Context, requestDto *apicontract.HardDeleteMyRoutineTaskByIdRequestDto) {
	response, exception := coreadapters.CallSecurly[apicontract.HardDeleteMyRoutineTaskByIdRequestDto, apicontract.HardDeleteMyRoutineTaskByIdResponseDto](
		ctx,
		c.coreClient,
		requestDto,
		apicontract.HardDeleteMyRoutineTaskByIdOperation,
		"/core/v1/routine-tasks/hard-delete",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *RoutineTaskController) HardDeleteMyRoutineTasksByIds(ctx *gin.Context, requestDto *apicontract.HardDeleteMyRoutineTasksByIdsRequestDto) {
	response, exception := coreadapters.CallSecurly[apicontract.HardDeleteMyRoutineTasksByIdsRequestDto, apicontract.HardDeleteMyRoutineTasksByIdsResponseDto](
		ctx,
		c.coreClient,
		requestDto,
		apicontract.HardDeleteMyRoutineTasksByIdsOperation,
		"/core/v1/routine-tasks/hard-delete-many",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *RoutineTaskController) VisualizeMyRoutineTaskStatusCount(ctx *gin.Context, requestDto *apicontract.VisualizeMyRoutineTaskStatusCountRequestDto) {
	response, exception := coreadapters.CallSecurly[apicontract.VisualizeMyRoutineTaskStatusCountRequestDto, apicontract.VisualizeMyRoutineTaskStatusCountResponseDto](
		ctx,
		c.coreClient,
		requestDto,
		apicontract.VisualizeMyRoutineTaskStatusCountOperation,
		"/core/v1/routine-tasks/visualize-status-count",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *RoutineTaskController) VisualizeMyRoutineTaskPurposeCount(ctx *gin.Context, requestDto *apicontract.VisualizeMyRoutineTaskPurposeCountRequestDto) {
	response, exception := coreadapters.CallSecurly[apicontract.VisualizeMyRoutineTaskPurposeCountRequestDto, apicontract.VisualizeMyRoutineTaskPurposeCountResponseDto](
		ctx,
		c.coreClient,
		requestDto,
		apicontract.VisualizeMyRoutineTaskPurposeCountOperation,
		"/core/v1/routine-tasks/visualize-purpose-count",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *RoutineTaskController) VisualizeMyRoutineTaskScheduledAtCount(ctx *gin.Context, requestDto *apicontract.VisualizeMyRoutineTaskScheduledAtCountRequestDto) {
	response, exception := coreadapters.CallSecurly[apicontract.VisualizeMyRoutineTaskScheduledAtCountRequestDto, apicontract.VisualizeMyRoutineTaskScheduledAtCountResponseDto](
		ctx,
		c.coreClient,
		requestDto,
		apicontract.VisualizeMyRoutineTaskScheduledAtCountOperation,
		"/core/v1/routine-tasks/visualize-scheduled-at-count",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *RoutineTaskController) VisualizeMyRoutineTaskActualStartedAtCount(ctx *gin.Context, requestDto *apicontract.VisualizeMyRoutineTaskActualStartedAtCountRequestDto) {
	response, exception := coreadapters.CallSecurly[apicontract.VisualizeMyRoutineTaskActualStartedAtCountRequestDto, apicontract.VisualizeMyRoutineTaskActualStartedAtCountResponseDto](
		ctx,
		c.coreClient,
		requestDto,
		apicontract.VisualizeMyRoutineTaskActualStartedAtCountOperation,
		"/core/v1/routine-tasks/visualize-actual-started-at-count",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *RoutineTaskController) VisualizeMyRoutineTaskActualEndedAtCount(ctx *gin.Context, requestDto *apicontract.VisualizeMyRoutineTaskActualEndedAtCountRequestDto) {
	response, exception := coreadapters.CallSecurly[apicontract.VisualizeMyRoutineTaskActualEndedAtCountRequestDto, apicontract.VisualizeMyRoutineTaskActualEndedAtCountResponseDto](
		ctx,
		c.coreClient,
		requestDto,
		apicontract.VisualizeMyRoutineTaskActualEndedAtCountOperation,
		"/core/v1/routine-tasks/visualize-actual-ended-at-count",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}
