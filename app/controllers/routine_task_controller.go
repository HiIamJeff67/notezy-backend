package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	dtos "github.com/HiIamJeff67/notezy-backend/app/dtos"
	services "github.com/HiIamJeff67/notezy-backend/app/services"
)

type RoutineTaskControllerInterface interface {
	GetMyRoutineTaskById(ctx *gin.Context, reqDto *dtos.GetMyRoutineTaskByIdReqDto)
	GetAllMyRoutineTasksByStationIds(ctx *gin.Context, reqDto *dtos.GetAllMyRoutineTasksByStationIdsReqDto)
	GetAllMyRoutineTasks(ctx *gin.Context, reqDto *dtos.GetAllMyRoutineTasksReqDto)
	CreateRoutineTaskByStationId(ctx *gin.Context, reqDto *dtos.CreateRoutineTaskByStationIdReqDto)
	UpdateMyRoutineTaskById(ctx *gin.Context, reqDto *dtos.UpdateMyRoutineTaskByIdReqDto)
	HardDeleteMyRoutineTaskById(ctx *gin.Context, reqDto *dtos.HardDeleteMyRoutineTaskByIdReqDto)
	HardDeleteMyRoutineTasksByIds(ctx *gin.Context, reqDto *dtos.HardDeleteMyRoutineTasksByIdsReqDto)
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

func (c *RoutineTaskController) GetAllMyRoutineTasksByStationIds(
	ctx *gin.Context,
	reqDto *dtos.GetAllMyRoutineTasksByStationIdsReqDto,
) {
	resDto, exception := c.routineTaskService.GetAllMyRoutineTasksByStationIds(ctx.Request.Context(), reqDto)
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

func (c *RoutineTaskController) CreateRoutineTaskByStationId(ctx *gin.Context, reqDto *dtos.CreateRoutineTaskByStationIdReqDto) {
	resDto, exception := c.routineTaskService.CreateRoutineTaskByStationId(ctx.Request.Context(), reqDto)
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
