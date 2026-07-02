package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	dtos "github.com/HiIamJeff67/notezy-backend/app/dtos"
	services "github.com/HiIamJeff67/notezy-backend/app/services"
)

type RoutineTaskRecordControllerInterface interface {
	GetAllMyRoutineTaskRecordsByRoutineTaskId(ctx *gin.Context, reqDto *dtos.GetAllMyRoutineTaskRecordsByRoutineTaskIdReqDto)
	VisualizeMyRoutineTaskRecordStatusCount(ctx *gin.Context, reqDto *dtos.VisualizeMyRoutineTaskRecordStatusCountReqDto)
	VisualizeMyRoutineTaskRecordPurposeCount(ctx *gin.Context, reqDto *dtos.VisualizeMyRoutineTaskRecordPurposeCountReqDto)
	VisualizeMyRoutineTaskRecordScheduledAtCount(ctx *gin.Context, reqDto *dtos.VisualizeMyRoutineTaskRecordScheduledAtCountReqDto)
	VisualizeMyRoutineTaskRecordActualStartedAtCount(ctx *gin.Context, reqDto *dtos.VisualizeMyRoutineTaskRecordActualStartedAtCountReqDto)
	VisualizeMyRoutineTaskRecordActualEndedAtCount(ctx *gin.Context, reqDto *dtos.VisualizeMyRoutineTaskRecordActualEndedAtCountReqDto)
}

type RoutineTaskRecordController struct {
	routineTaskRecordService services.RoutineTaskRecordServiceInterface
}

func NewRoutineTaskRecordController(
	routineTaskRecordService services.RoutineTaskRecordServiceInterface,
) RoutineTaskRecordControllerInterface {
	return &RoutineTaskRecordController{
		routineTaskRecordService: routineTaskRecordService,
	}
}

func (c *RoutineTaskRecordController) GetAllMyRoutineTaskRecordsByRoutineTaskId(ctx *gin.Context, reqDto *dtos.GetAllMyRoutineTaskRecordsByRoutineTaskIdReqDto) {
	resDto, exception := c.routineTaskRecordService.GetAllMyRoutineTaskRecordsByRoutineTaskId(ctx.Request.Context(), reqDto)
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

func (c *RoutineTaskRecordController) VisualizeMyRoutineTaskRecordStatusCount(ctx *gin.Context, reqDto *dtos.VisualizeMyRoutineTaskRecordStatusCountReqDto) {
	resDto, exception := c.routineTaskRecordService.VisualizeMyRoutineTaskRecordStatusCount(ctx.Request.Context(), reqDto)
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

func (c *RoutineTaskRecordController) VisualizeMyRoutineTaskRecordPurposeCount(ctx *gin.Context, reqDto *dtos.VisualizeMyRoutineTaskRecordPurposeCountReqDto) {
	resDto, exception := c.routineTaskRecordService.VisualizeMyRoutineTaskRecordPurposeCount(ctx.Request.Context(), reqDto)
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

func (c *RoutineTaskRecordController) VisualizeMyRoutineTaskRecordScheduledAtCount(ctx *gin.Context, reqDto *dtos.VisualizeMyRoutineTaskRecordScheduledAtCountReqDto) {
	resDto, exception := c.routineTaskRecordService.VisualizeMyRoutineTaskRecordScheduledAtCount(ctx.Request.Context(), reqDto)
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

func (c *RoutineTaskRecordController) VisualizeMyRoutineTaskRecordActualStartedAtCount(
	ctx *gin.Context,
	reqDto *dtos.VisualizeMyRoutineTaskRecordActualStartedAtCountReqDto,
) {
	resDto, exception := c.routineTaskRecordService.VisualizeMyRoutineTaskRecordActualStartedAtCount(ctx.Request.Context(), reqDto)
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

func (c *RoutineTaskRecordController) VisualizeMyRoutineTaskRecordActualEndedAtCount(ctx *gin.Context, reqDto *dtos.VisualizeMyRoutineTaskRecordActualEndedAtCountReqDto) {
	resDto, exception := c.routineTaskRecordService.VisualizeMyRoutineTaskRecordActualEndedAtCount(ctx.Request.Context(), reqDto)
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
