package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	routinetasksdto "github.com/HiIamJeff67/notezy-backend/contracts/gateway/v1/api/routine-tasks"
	responsewriter "github.com/HiIamJeff67/notezy-backend/internal/shared/responsewriter"
	coreadapters "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/core/adapters"
)

type RoutineTaskControllerInterface interface {
	GetMyRoutineTaskById(ctx *gin.Context, requestDto *routinetasksdto.GetMyRoutineTaskByIdRequestDto)
	GetAllMyRoutineTasksByRoutineIds(ctx *gin.Context, requestDto *routinetasksdto.GetAllMyRoutineTasksByRoutineIdsRequestDto)
	GetAllMyRoutineTasks(ctx *gin.Context, requestDto *routinetasksdto.GetAllMyRoutineTasksRequestDto)
	CreateRoutineTaskByRoutineId(ctx *gin.Context, requestDto *routinetasksdto.CreateRoutineTaskByRoutineIdRequestDto)
	UpdateMyRoutineTaskById(ctx *gin.Context, requestDto *routinetasksdto.UpdateMyRoutineTaskByIdRequestDto)
	PauseMyRoutineTaskById(ctx *gin.Context, requestDto *routinetasksdto.PauseMyRoutineTaskByIdRequestDto)
	ResumeMyRoutineTaskById(ctx *gin.Context, requestDto *routinetasksdto.ResumeMyRoutineTaskByIdRequestDto)
	HardDeleteMyRoutineTaskById(ctx *gin.Context, requestDto *routinetasksdto.HardDeleteMyRoutineTaskByIdRequestDto)
	HardDeleteMyRoutineTasksByIds(ctx *gin.Context, requestDto *routinetasksdto.HardDeleteMyRoutineTasksByIdsRequestDto)

	/* ============================== Visualization Methods ============================== */
	VisualizeMyRoutineTaskStatusCount(ctx *gin.Context, requestDto *routinetasksdto.VisualizeMyRoutineTaskStatusCountRequestDto)
	VisualizeMyRoutineTaskPurposeCount(ctx *gin.Context, requestDto *routinetasksdto.VisualizeMyRoutineTaskPurposeCountRequestDto)
	VisualizeMyRoutineTaskScheduledAtCount(ctx *gin.Context, requestDto *routinetasksdto.VisualizeMyRoutineTaskScheduledAtCountRequestDto)
	VisualizeMyRoutineTaskActualStartedAtCount(ctx *gin.Context, requestDto *routinetasksdto.VisualizeMyRoutineTaskActualStartedAtCountRequestDto)
	VisualizeMyRoutineTaskActualEndedAtCount(ctx *gin.Context, requestDto *routinetasksdto.VisualizeMyRoutineTaskActualEndedAtCountRequestDto)
}

type RoutineTaskController struct {
	coreClient *coreadapters.CoreClient
}

func NewRoutineTaskController(coreClient *coreadapters.CoreClient) RoutineTaskControllerInterface {
	return &RoutineTaskController{coreClient: coreClient}
}

func (c *RoutineTaskController) GetMyRoutineTaskById(ctx *gin.Context, requestDto *routinetasksdto.GetMyRoutineTaskByIdRequestDto) {
	response, exception := coreadapters.CallSecurly[routinetasksdto.GetMyRoutineTaskByIdRequestDto, routinetasksdto.GetMyRoutineTaskByIdResponseDto](
		ctx,
		c.coreClient,
		requestDto,
		routinetasksdto.GetMyRoutineTaskByIdOperation,
		"/core/v1/routine-tasks/get-by-id",
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

func (c *RoutineTaskController) GetAllMyRoutineTasksByRoutineIds(ctx *gin.Context, requestDto *routinetasksdto.GetAllMyRoutineTasksByRoutineIdsRequestDto) {
	response, exception := coreadapters.CallSecurly[routinetasksdto.GetAllMyRoutineTasksByRoutineIdsRequestDto, routinetasksdto.GetAllMyRoutineTasksByRoutineIdsResponseDto](
		ctx,
		c.coreClient,
		requestDto,
		routinetasksdto.GetAllMyRoutineTasksByRoutineIdsOperation,
		"/core/v1/routine-tasks/get-all-by-routine-ids",
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

func (c *RoutineTaskController) GetAllMyRoutineTasks(ctx *gin.Context, requestDto *routinetasksdto.GetAllMyRoutineTasksRequestDto) {
	response, exception := coreadapters.CallSecurly[routinetasksdto.GetAllMyRoutineTasksRequestDto, routinetasksdto.GetAllMyRoutineTasksResponseDto](
		ctx,
		c.coreClient,
		requestDto,
		routinetasksdto.GetAllMyRoutineTasksOperation,
		"/core/v1/routine-tasks/get-all",
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

func (c *RoutineTaskController) CreateRoutineTaskByRoutineId(ctx *gin.Context, requestDto *routinetasksdto.CreateRoutineTaskByRoutineIdRequestDto) {
	response, exception := coreadapters.CallSecurly[routinetasksdto.CreateRoutineTaskByRoutineIdRequestDto, routinetasksdto.CreateRoutineTaskByRoutineIdResponseDto](
		ctx,
		c.coreClient,
		requestDto,
		routinetasksdto.CreateRoutineTaskByRoutineIdOperation,
		"/core/v1/routine-tasks/create-by-routine-id",
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

func (c *RoutineTaskController) UpdateMyRoutineTaskById(ctx *gin.Context, requestDto *routinetasksdto.UpdateMyRoutineTaskByIdRequestDto) {
	response, exception := coreadapters.CallSecurly[routinetasksdto.UpdateMyRoutineTaskByIdRequestDto, routinetasksdto.UpdateMyRoutineTaskByIdResponseDto](
		ctx,
		c.coreClient,
		requestDto,
		routinetasksdto.UpdateMyRoutineTaskByIdOperation,
		"/core/v1/routine-tasks/update",
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

func (c *RoutineTaskController) PauseMyRoutineTaskById(ctx *gin.Context, requestDto *routinetasksdto.PauseMyRoutineTaskByIdRequestDto) {
	response, exception := coreadapters.CallSecurly[routinetasksdto.PauseMyRoutineTaskByIdRequestDto, routinetasksdto.PauseMyRoutineTaskByIdResponseDto](
		ctx,
		c.coreClient,
		requestDto,
		routinetasksdto.PauseMyRoutineTaskByIdOperation,
		"/core/v1/routine-tasks/pause",
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

func (c *RoutineTaskController) ResumeMyRoutineTaskById(ctx *gin.Context, requestDto *routinetasksdto.ResumeMyRoutineTaskByIdRequestDto) {
	response, exception := coreadapters.CallSecurly[routinetasksdto.ResumeMyRoutineTaskByIdRequestDto, routinetasksdto.ResumeMyRoutineTaskByIdResponseDto](
		ctx,
		c.coreClient,
		requestDto,
		routinetasksdto.ResumeMyRoutineTaskByIdOperation,
		"/core/v1/routine-tasks/resume",
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

func (c *RoutineTaskController) HardDeleteMyRoutineTaskById(ctx *gin.Context, requestDto *routinetasksdto.HardDeleteMyRoutineTaskByIdRequestDto) {
	response, exception := coreadapters.CallSecurly[routinetasksdto.HardDeleteMyRoutineTaskByIdRequestDto, routinetasksdto.HardDeleteMyRoutineTaskByIdResponseDto](
		ctx,
		c.coreClient,
		requestDto,
		routinetasksdto.HardDeleteMyRoutineTaskByIdOperation,
		"/core/v1/routine-tasks/hard-delete",
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

func (c *RoutineTaskController) HardDeleteMyRoutineTasksByIds(ctx *gin.Context, requestDto *routinetasksdto.HardDeleteMyRoutineTasksByIdsRequestDto) {
	response, exception := coreadapters.CallSecurly[routinetasksdto.HardDeleteMyRoutineTasksByIdsRequestDto, routinetasksdto.HardDeleteMyRoutineTasksByIdsResponseDto](
		ctx,
		c.coreClient,
		requestDto,
		routinetasksdto.HardDeleteMyRoutineTasksByIdsOperation,
		"/core/v1/routine-tasks/hard-delete-many",
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

func (c *RoutineTaskController) VisualizeMyRoutineTaskStatusCount(ctx *gin.Context, requestDto *routinetasksdto.VisualizeMyRoutineTaskStatusCountRequestDto) {
	response, exception := coreadapters.CallSecurly[routinetasksdto.VisualizeMyRoutineTaskStatusCountRequestDto, routinetasksdto.VisualizeMyRoutineTaskStatusCountResponseDto](
		ctx,
		c.coreClient,
		requestDto,
		routinetasksdto.VisualizeMyRoutineTaskStatusCountOperation,
		"/core/v1/routine-tasks/visualize-status-count",
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

func (c *RoutineTaskController) VisualizeMyRoutineTaskPurposeCount(ctx *gin.Context, requestDto *routinetasksdto.VisualizeMyRoutineTaskPurposeCountRequestDto) {
	response, exception := coreadapters.CallSecurly[routinetasksdto.VisualizeMyRoutineTaskPurposeCountRequestDto, routinetasksdto.VisualizeMyRoutineTaskPurposeCountResponseDto](
		ctx,
		c.coreClient,
		requestDto,
		routinetasksdto.VisualizeMyRoutineTaskPurposeCountOperation,
		"/core/v1/routine-tasks/visualize-purpose-count",
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

func (c *RoutineTaskController) VisualizeMyRoutineTaskScheduledAtCount(ctx *gin.Context, requestDto *routinetasksdto.VisualizeMyRoutineTaskScheduledAtCountRequestDto) {
	response, exception := coreadapters.CallSecurly[routinetasksdto.VisualizeMyRoutineTaskScheduledAtCountRequestDto, routinetasksdto.VisualizeMyRoutineTaskScheduledAtCountResponseDto](
		ctx,
		c.coreClient,
		requestDto,
		routinetasksdto.VisualizeMyRoutineTaskScheduledAtCountOperation,
		"/core/v1/routine-tasks/visualize-scheduled-at-count",
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

func (c *RoutineTaskController) VisualizeMyRoutineTaskActualStartedAtCount(ctx *gin.Context, requestDto *routinetasksdto.VisualizeMyRoutineTaskActualStartedAtCountRequestDto) {
	response, exception := coreadapters.CallSecurly[routinetasksdto.VisualizeMyRoutineTaskActualStartedAtCountRequestDto, routinetasksdto.VisualizeMyRoutineTaskActualStartedAtCountResponseDto](
		ctx,
		c.coreClient,
		requestDto,
		routinetasksdto.VisualizeMyRoutineTaskActualStartedAtCountOperation,
		"/core/v1/routine-tasks/visualize-actual-started-at-count",
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

func (c *RoutineTaskController) VisualizeMyRoutineTaskActualEndedAtCount(ctx *gin.Context, requestDto *routinetasksdto.VisualizeMyRoutineTaskActualEndedAtCountRequestDto) {
	response, exception := coreadapters.CallSecurly[routinetasksdto.VisualizeMyRoutineTaskActualEndedAtCountRequestDto, routinetasksdto.VisualizeMyRoutineTaskActualEndedAtCountResponseDto](
		ctx,
		c.coreClient,
		requestDto,
		routinetasksdto.VisualizeMyRoutineTaskActualEndedAtCountOperation,
		"/core/v1/routine-tasks/visualize-actual-ended-at-count",
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
