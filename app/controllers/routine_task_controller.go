package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	dtos "github.com/HiIamJeff67/notezy-backend/app/dtos"
	services "github.com/HiIamJeff67/notezy-backend/app/services"
)

type RoutineTaskControllerInterface interface {
	GetMyRoutineTaskById(ctx *gin.Context, reqDto *dtos.GetMyRoutineTaskByIdReqDto)
	GetAllMyRoutineTasksByRoutineIds(ctx *gin.Context, reqDto *dtos.GetAllMyRoutineTasksByRoutineIdsReqDto)
	GetAllMyRoutineTasks(ctx *gin.Context, reqDto *dtos.GetAllMyRoutineTasksReqDto)
	CreateRoutineTaskByRoutineId(ctx *gin.Context, reqDto *dtos.CreateRoutineTaskByRoutineIdReqDto)
	UpdateMyRoutineTaskById(ctx *gin.Context, reqDto *dtos.UpdateMyRoutineTaskByIdReqDto)
	PauseMyRoutineTaskById(ctx *gin.Context, reqDto *dtos.PauseMyRoutineTaskByIdReqDto)
	ResumeMyRoutineTaskById(ctx *gin.Context, reqDto *dtos.ResumeMyRoutineTaskByIdReqDto)
	HardDeleteMyRoutineTaskById(ctx *gin.Context, reqDto *dtos.HardDeleteMyRoutineTaskByIdReqDto)
	HardDeleteMyRoutineTasksByIds(ctx *gin.Context, reqDto *dtos.HardDeleteMyRoutineTasksByIdsReqDto)
	VisualizeMyRoutineTaskStatusCount(ctx *gin.Context, reqDto *dtos.VisualizeMyRoutineTaskStatusCountReqDto)
	VisualizeMyRoutineTaskPurposeCount(ctx *gin.Context, reqDto *dtos.VisualizeMyRoutineTaskPurposeCountReqDto)
	VisualizeMyRoutineTaskScheduledAtCount(ctx *gin.Context, reqDto *dtos.VisualizeMyRoutineTaskScheduledAtCountReqDto)
	VisualizeMyRoutineTaskActualStartedAtCount(ctx *gin.Context, reqDto *dtos.VisualizeMyRoutineTaskActualStartedAtCountReqDto)
	VisualizeMyRoutineTaskActualEndedAtCount(ctx *gin.Context, reqDto *dtos.VisualizeMyRoutineTaskActualEndedAtCountReqDto)
}

type RoutineTaskController struct {
	routineTaskService services.RoutineTaskServiceInterface
}

func NewRoutineTaskController(service services.RoutineTaskServiceInterface) RoutineTaskControllerInterface {
	return &RoutineTaskController{
		routineTaskService: service,
	}
}

func (c *RoutineTaskController) GetMyRoutineTaskById(ctx *gin.Context, reqDto *dtos.GetMyRoutineTaskByIdReqDto) {
	resDto, exception := c.routineTaskService.GetMyRoutineTaskById(ctx.Request.Context(), reqDto)
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

func (c *RoutineTaskController) GetAllMyRoutineTasksByRoutineIds(
	ctx *gin.Context,
	reqDto *dtos.GetAllMyRoutineTasksByRoutineIdsReqDto,
) {
	resDto, exception := c.routineTaskService.GetAllMyRoutineTasksByRoutineIds(ctx.Request.Context(), reqDto)
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

func (c *RoutineTaskController) GetAllMyRoutineTasks(ctx *gin.Context, reqDto *dtos.GetAllMyRoutineTasksReqDto) {
	resDto, exception := c.routineTaskService.GetAllMyRoutineTasks(ctx.Request.Context(), reqDto)
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

func (c *RoutineTaskController) CreateRoutineTaskByRoutineId(ctx *gin.Context, reqDto *dtos.CreateRoutineTaskByRoutineIdReqDto) {
	resDto, exception := c.routineTaskService.CreateRoutineTaskByRoutineId(ctx.Request.Context(), reqDto)
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

func (c *RoutineTaskController) UpdateMyRoutineTaskById(ctx *gin.Context, reqDto *dtos.UpdateMyRoutineTaskByIdReqDto) {
	resDto, exception := c.routineTaskService.UpdateMyRoutineTaskById(ctx.Request.Context(), reqDto)
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

func (c *RoutineTaskController) PauseMyRoutineTaskById(ctx *gin.Context, reqDto *dtos.PauseMyRoutineTaskByIdReqDto) {
	resDto, exception := c.routineTaskService.PauseMyRoutineTaskById(ctx.Request.Context(), reqDto)
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

func (c *RoutineTaskController) ResumeMyRoutineTaskById(ctx *gin.Context, reqDto *dtos.ResumeMyRoutineTaskByIdReqDto) {
	resDto, exception := c.routineTaskService.ResumeMyRoutineTaskById(ctx.Request.Context(), reqDto)
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

func (c *RoutineTaskController) HardDeleteMyRoutineTaskById(ctx *gin.Context, reqDto *dtos.HardDeleteMyRoutineTaskByIdReqDto) {
	resDto, exception := c.routineTaskService.HardDeleteMyRoutineTaskById(ctx.Request.Context(), reqDto)
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

func (c *RoutineTaskController) HardDeleteMyRoutineTasksByIds(ctx *gin.Context, reqDto *dtos.HardDeleteMyRoutineTasksByIdsReqDto) {
	resDto, exception := c.routineTaskService.HardDeleteMyRoutineTasksByIds(ctx.Request.Context(), reqDto)
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

func (c *RoutineTaskController) VisualizeMyRoutineTaskStatusCount(
	ctx *gin.Context,
	reqDto *dtos.VisualizeMyRoutineTaskStatusCountReqDto,
) {
	resDto, exception := c.routineTaskService.VisualizeMyRoutineTaskStatusCount(ctx.Request.Context(), reqDto)
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

func (c *RoutineTaskController) VisualizeMyRoutineTaskPurposeCount(
	ctx *gin.Context,
	reqDto *dtos.VisualizeMyRoutineTaskPurposeCountReqDto,
) {
	resDto, exception := c.routineTaskService.VisualizeMyRoutineTaskPurposeCount(ctx.Request.Context(), reqDto)
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

func (c *RoutineTaskController) VisualizeMyRoutineTaskScheduledAtCount(
	ctx *gin.Context,
	reqDto *dtos.VisualizeMyRoutineTaskScheduledAtCountReqDto,
) {
	resDto, exception := c.routineTaskService.VisualizeMyRoutineTaskScheduledAtCount(ctx.Request.Context(), reqDto)
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

func (c *RoutineTaskController) VisualizeMyRoutineTaskActualStartedAtCount(
	ctx *gin.Context,
	reqDto *dtos.VisualizeMyRoutineTaskActualStartedAtCountReqDto,
) {
	resDto, exception := c.routineTaskService.VisualizeMyRoutineTaskActualStartedAtCount(ctx.Request.Context(), reqDto)
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

func (c *RoutineTaskController) VisualizeMyRoutineTaskActualEndedAtCount(
	ctx *gin.Context,
	reqDto *dtos.VisualizeMyRoutineTaskActualEndedAtCountReqDto,
) {
	resDto, exception := c.routineTaskService.VisualizeMyRoutineTaskActualEndedAtCount(ctx.Request.Context(), reqDto)
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
