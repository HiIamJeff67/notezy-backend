package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	routinetaskrecordsdto "github.com/HiIamJeff67/notezy-backend/contracts/gateway/v1/api/routine-task-records"
	responsewriter "github.com/HiIamJeff67/notezy-backend/internal/shared/responsewriter"
	coreadapters "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/core/adapters"
)

type RoutineTaskRecordControllerInterface interface {
	GetAllMyRoutineTaskRecordsByRoutineTaskId(ctx *gin.Context, requestDto *routinetaskrecordsdto.GetAllMyRoutineTaskRecordsByRoutineTaskIdRequestDto)

	/* ============================== Visualization Methods ============================== */
	VisualizeMyRoutineTaskRecordStatusCount(ctx *gin.Context, requestDto *routinetaskrecordsdto.VisualizeMyRoutineTaskRecordStatusCountRequestDto)
	VisualizeMyRoutineTaskRecordPurposeCount(ctx *gin.Context, requestDto *routinetaskrecordsdto.VisualizeMyRoutineTaskRecordPurposeCountRequestDto)
	VisualizeMyRoutineTaskRecordScheduledAtCount(ctx *gin.Context, requestDto *routinetaskrecordsdto.VisualizeMyRoutineTaskRecordScheduledAtCountRequestDto)
	VisualizeMyRoutineTaskRecordActualStartedAtCount(ctx *gin.Context, requestDto *routinetaskrecordsdto.VisualizeMyRoutineTaskRecordActualStartedAtCountRequestDto)
	VisualizeMyRoutineTaskRecordActualEndedAtCount(ctx *gin.Context, requestDto *routinetaskrecordsdto.VisualizeMyRoutineTaskRecordActualEndedAtCountRequestDto)
}

type RoutineTaskRecordController struct {
	coreClient *coreadapters.CoreClient
}

func NewRoutineTaskRecordController(coreClient *coreadapters.CoreClient) RoutineTaskRecordControllerInterface {
	return &RoutineTaskRecordController{
		coreClient: coreClient,
	}
}

func (c *RoutineTaskRecordController) GetAllMyRoutineTaskRecordsByRoutineTaskId(ctx *gin.Context, requestDto *routinetaskrecordsdto.GetAllMyRoutineTaskRecordsByRoutineTaskIdRequestDto) {
	response, exception := coreadapters.CallSecurly[routinetaskrecordsdto.GetAllMyRoutineTaskRecordsByRoutineTaskIdRequestDto, routinetaskrecordsdto.GetAllMyRoutineTaskRecordsByRoutineTaskIdResponseDto](ctx, c.coreClient, requestDto, routinetaskrecordsdto.GetAllMyRoutineTaskRecordsByRoutineTaskIdOperation, "/core/v1/routine-task-records/get-all-by-routine-task-id")
	if exception != nil {
		responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "data": response.Data, "exception": nil})
}

/* ============================== Visualization Methods ============================== */

func (c *RoutineTaskRecordController) VisualizeMyRoutineTaskRecordStatusCount(ctx *gin.Context, requestDto *routinetaskrecordsdto.VisualizeMyRoutineTaskRecordStatusCountRequestDto) {
	response, exception := coreadapters.CallSecurly[routinetaskrecordsdto.VisualizeMyRoutineTaskRecordStatusCountRequestDto, routinetaskrecordsdto.VisualizeMyRoutineTaskRecordStatusCountResponseDto](ctx, c.coreClient, requestDto, routinetaskrecordsdto.VisualizeMyRoutineTaskRecordStatusCountOperation, "/core/v1/routine-task-records/visualizations/status-count")
	if exception != nil {
		responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "data": response.Data, "exception": nil})
}

func (c *RoutineTaskRecordController) VisualizeMyRoutineTaskRecordPurposeCount(ctx *gin.Context, requestDto *routinetaskrecordsdto.VisualizeMyRoutineTaskRecordPurposeCountRequestDto) {
	response, exception := coreadapters.CallSecurly[routinetaskrecordsdto.VisualizeMyRoutineTaskRecordPurposeCountRequestDto, routinetaskrecordsdto.VisualizeMyRoutineTaskRecordPurposeCountResponseDto](ctx, c.coreClient, requestDto, routinetaskrecordsdto.VisualizeMyRoutineTaskRecordPurposeCountOperation, "/core/v1/routine-task-records/visualizations/purpose-count")
	if exception != nil {
		responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "data": response.Data, "exception": nil})
}

func (c *RoutineTaskRecordController) VisualizeMyRoutineTaskRecordScheduledAtCount(ctx *gin.Context, requestDto *routinetaskrecordsdto.VisualizeMyRoutineTaskRecordScheduledAtCountRequestDto) {
	response, exception := coreadapters.CallSecurly[routinetaskrecordsdto.VisualizeMyRoutineTaskRecordScheduledAtCountRequestDto, routinetaskrecordsdto.VisualizeMyRoutineTaskRecordScheduledAtCountResponseDto](ctx, c.coreClient, requestDto, routinetaskrecordsdto.VisualizeMyRoutineTaskRecordScheduledAtCountOperation, "/core/v1/routine-task-records/visualizations/scheduled-at-count")
	if exception != nil {
		responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "data": response.Data, "exception": nil})
}

func (c *RoutineTaskRecordController) VisualizeMyRoutineTaskRecordActualStartedAtCount(ctx *gin.Context, requestDto *routinetaskrecordsdto.VisualizeMyRoutineTaskRecordActualStartedAtCountRequestDto) {
	response, exception := coreadapters.CallSecurly[routinetaskrecordsdto.VisualizeMyRoutineTaskRecordActualStartedAtCountRequestDto, routinetaskrecordsdto.VisualizeMyRoutineTaskRecordActualStartedAtCountResponseDto](ctx, c.coreClient, requestDto, routinetaskrecordsdto.VisualizeMyRoutineTaskRecordActualStartedAtCountOperation, "/core/v1/routine-task-records/visualizations/actual-started-at-count")
	if exception != nil {
		responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "data": response.Data, "exception": nil})
}

func (c *RoutineTaskRecordController) VisualizeMyRoutineTaskRecordActualEndedAtCount(ctx *gin.Context, requestDto *routinetaskrecordsdto.VisualizeMyRoutineTaskRecordActualEndedAtCountRequestDto) {
	response, exception := coreadapters.CallSecurly[routinetaskrecordsdto.VisualizeMyRoutineTaskRecordActualEndedAtCountRequestDto, routinetaskrecordsdto.VisualizeMyRoutineTaskRecordActualEndedAtCountResponseDto](ctx, c.coreClient, requestDto, routinetaskrecordsdto.VisualizeMyRoutineTaskRecordActualEndedAtCountOperation, "/core/v1/routine-task-records/visualizations/actual-ended-at-count")
	if exception != nil {
		responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "data": response.Data, "exception": nil})
}
