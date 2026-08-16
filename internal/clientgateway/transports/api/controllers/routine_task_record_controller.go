package controllers

import (
	"github.com/gin-gonic/gin"

	exceptionwriter "github.com/HiIamJeff67/notegic-backend/shared/util/exceptionwriter"

	apicontract "github.com/HiIamJeff67/notegic-backend/contracts/core/v1/api/routine-task-records"

	coreadapters "github.com/HiIamJeff67/notegic-backend/internal/clientgateway/transports/core/adapters"
)

type RoutineTaskRecordControllerInterface interface {
	GetAllMyRoutineTaskRecordsByRoutineTaskId(ctx *gin.Context, requestDto *apicontract.GetAllMyRoutineTaskRecordsByRoutineTaskIdRequestDto)

	/* ============================== Visualization Methods ============================== */
	VisualizeMyRoutineTaskRecordStatusCount(ctx *gin.Context, requestDto *apicontract.VisualizeMyRoutineTaskRecordStatusCountRequestDto)
	VisualizeMyRoutineTaskRecordPurposeCount(ctx *gin.Context, requestDto *apicontract.VisualizeMyRoutineTaskRecordPurposeCountRequestDto)
	VisualizeMyRoutineTaskRecordScheduledAtCount(ctx *gin.Context, requestDto *apicontract.VisualizeMyRoutineTaskRecordScheduledAtCountRequestDto)
	VisualizeMyRoutineTaskRecordActualStartedAtCount(ctx *gin.Context, requestDto *apicontract.VisualizeMyRoutineTaskRecordActualStartedAtCountRequestDto)
	VisualizeMyRoutineTaskRecordActualEndedAtCount(ctx *gin.Context, requestDto *apicontract.VisualizeMyRoutineTaskRecordActualEndedAtCountRequestDto)
}

type RoutineTaskRecordController struct {
	coreAdapter *coreadapters.CoreAdapter
}

func NewRoutineTaskRecordController(coreAdapter *coreadapters.CoreAdapter) RoutineTaskRecordControllerInterface {
	return &RoutineTaskRecordController{
		coreAdapter: coreAdapter,
	}
}

func (c *RoutineTaskRecordController) GetAllMyRoutineTaskRecordsByRoutineTaskId(ctx *gin.Context, requestDto *apicontract.GetAllMyRoutineTaskRecordsByRoutineTaskIdRequestDto) {
	response, exception := coreadapters.CallSecurly[apicontract.GetAllMyRoutineTaskRecordsByRoutineTaskIdRequestDto, apicontract.GetAllMyRoutineTaskRecordsByRoutineTaskIdResponseDto](ctx, c.coreAdapter, requestDto, apicontract.GetAllMyRoutineTaskRecordsByRoutineTaskIdOperation, "/core/v1/routine-task-records/get-all-by-routine-task-id")
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}
	writeClientResponse(ctx, response.Data)
}

/* ============================== Visualization Methods ============================== */

func (c *RoutineTaskRecordController) VisualizeMyRoutineTaskRecordStatusCount(ctx *gin.Context, requestDto *apicontract.VisualizeMyRoutineTaskRecordStatusCountRequestDto) {
	response, exception := coreadapters.CallSecurly[apicontract.VisualizeMyRoutineTaskRecordStatusCountRequestDto, apicontract.VisualizeMyRoutineTaskRecordStatusCountResponseDto](ctx, c.coreAdapter, requestDto, apicontract.VisualizeMyRoutineTaskRecordStatusCountOperation, "/core/v1/routine-task-records/visualizations/status-count")
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}
	writeClientResponse(ctx, response.Data)
}

func (c *RoutineTaskRecordController) VisualizeMyRoutineTaskRecordPurposeCount(ctx *gin.Context, requestDto *apicontract.VisualizeMyRoutineTaskRecordPurposeCountRequestDto) {
	response, exception := coreadapters.CallSecurly[apicontract.VisualizeMyRoutineTaskRecordPurposeCountRequestDto, apicontract.VisualizeMyRoutineTaskRecordPurposeCountResponseDto](ctx, c.coreAdapter, requestDto, apicontract.VisualizeMyRoutineTaskRecordPurposeCountOperation, "/core/v1/routine-task-records/visualizations/purpose-count")
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}
	writeClientResponse(ctx, response.Data)
}

func (c *RoutineTaskRecordController) VisualizeMyRoutineTaskRecordScheduledAtCount(ctx *gin.Context, requestDto *apicontract.VisualizeMyRoutineTaskRecordScheduledAtCountRequestDto) {
	response, exception := coreadapters.CallSecurly[apicontract.VisualizeMyRoutineTaskRecordScheduledAtCountRequestDto, apicontract.VisualizeMyRoutineTaskRecordScheduledAtCountResponseDto](ctx, c.coreAdapter, requestDto, apicontract.VisualizeMyRoutineTaskRecordScheduledAtCountOperation, "/core/v1/routine-task-records/visualizations/scheduled-at-count")
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}
	writeClientResponse(ctx, response.Data)
}

func (c *RoutineTaskRecordController) VisualizeMyRoutineTaskRecordActualStartedAtCount(ctx *gin.Context, requestDto *apicontract.VisualizeMyRoutineTaskRecordActualStartedAtCountRequestDto) {
	response, exception := coreadapters.CallSecurly[apicontract.VisualizeMyRoutineTaskRecordActualStartedAtCountRequestDto, apicontract.VisualizeMyRoutineTaskRecordActualStartedAtCountResponseDto](ctx, c.coreAdapter, requestDto, apicontract.VisualizeMyRoutineTaskRecordActualStartedAtCountOperation, "/core/v1/routine-task-records/visualizations/actual-started-at-count")
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}
	writeClientResponse(ctx, response.Data)
}

func (c *RoutineTaskRecordController) VisualizeMyRoutineTaskRecordActualEndedAtCount(ctx *gin.Context, requestDto *apicontract.VisualizeMyRoutineTaskRecordActualEndedAtCountRequestDto) {
	response, exception := coreadapters.CallSecurly[apicontract.VisualizeMyRoutineTaskRecordActualEndedAtCountRequestDto, apicontract.VisualizeMyRoutineTaskRecordActualEndedAtCountResponseDto](ctx, c.coreAdapter, requestDto, apicontract.VisualizeMyRoutineTaskRecordActualEndedAtCountOperation, "/core/v1/routine-task-records/visualizations/actual-ended-at-count")
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}
	writeClientResponse(ctx, response.Data)
}
