package binders

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	routinetaskrecordsdto "github.com/HiIamJeff67/notezy-backend/contracts/api/v1/routine-task-records"
	exceptions "github.com/HiIamJeff67/notezy-backend/internal/exceptions"
	responsewriter "github.com/HiIamJeff67/notezy-backend/internal/gateway/responsewriter"
	apitransport "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/api"
)

type RoutineTaskRecordBinderInterface interface {
	BindGetAllMyRoutineTaskRecordsByRoutineTaskId(controllerFunc apitransport.ControllerFunc[*routinetaskrecordsdto.GetAllMyRoutineTaskRecordsByRoutineTaskIdRequestDto]) gin.HandlerFunc

	/* ============================== Visualization Methods ============================== */
	BindVisualizeMyRoutineTaskRecordStatusCount(controllerFunc apitransport.ControllerFunc[*routinetaskrecordsdto.VisualizeMyRoutineTaskRecordStatusCountRequestDto]) gin.HandlerFunc
	BindVisualizeMyRoutineTaskRecordPurposeCount(controllerFunc apitransport.ControllerFunc[*routinetaskrecordsdto.VisualizeMyRoutineTaskRecordPurposeCountRequestDto]) gin.HandlerFunc
	BindVisualizeMyRoutineTaskRecordScheduledAtCount(controllerFunc apitransport.ControllerFunc[*routinetaskrecordsdto.VisualizeMyRoutineTaskRecordScheduledAtCountRequestDto]) gin.HandlerFunc
	BindVisualizeMyRoutineTaskRecordActualStartedAtCount(controllerFunc apitransport.ControllerFunc[*routinetaskrecordsdto.VisualizeMyRoutineTaskRecordActualStartedAtCountRequestDto]) gin.HandlerFunc
	BindVisualizeMyRoutineTaskRecordActualEndedAtCount(controllerFunc apitransport.ControllerFunc[*routinetaskrecordsdto.VisualizeMyRoutineTaskRecordActualEndedAtCountRequestDto]) gin.HandlerFunc
}

type RoutineTaskRecordBinder struct{}

func NewRoutineTaskRecordBinder() RoutineTaskRecordBinderInterface {
	return &RoutineTaskRecordBinder{}
}

/* ============================== Auxiliary Functions ============================== */

type routineTaskRecordVisualizationParams struct {
	permission     string
	routineTaskIds []uuid.UUID
}

func parseRoutineTaskRecordVisualizationParams(ctx *gin.Context) (*routineTaskRecordVisualizationParams, *exceptions.Exception) {
	permission := ctx.Query("permission")
	if permission == "" {
		return nil, exceptions.InvalidInput("RoutineTask").WithOrigin(fmt.Errorf("permission is required"))
	}

	params := &routineTaskRecordVisualizationParams{
		permission: permission,
	}
	for _, routineTaskIdsValue := range ctx.QueryArray("routineTaskIds") {
		for _, routineTaskIdValue := range strings.Split(routineTaskIdsValue, ",") {
			trimmedRoutineTaskIdValue := strings.TrimSpace(routineTaskIdValue)
			if trimmedRoutineTaskIdValue == "" {
				continue
			}
			routineTaskId, err := uuid.Parse(trimmedRoutineTaskIdValue)
			if err != nil {
				return nil, exceptions.InvalidInput("RoutineTask").WithOrigin(err)
			}
			params.routineTaskIds = append(params.routineTaskIds, routineTaskId)
		}
	}

	return params, nil
}

func parseRoutineTaskRecordVisualizationTimeRange(ctx *gin.Context) (int, time.Time, time.Time, *exceptions.Exception) {
	timeHourUnitString := ctx.Query("timeHourUnit")
	if timeHourUnitString == "" {
		return 0, time.Time{}, time.Time{}, exceptions.InvalidInput("RoutineTask").WithOrigin(fmt.Errorf("timeHourUnit is required"))
	}
	timeHourUnit, err := strconv.Atoi(timeHourUnitString)
	if err != nil {
		return 0, time.Time{}, time.Time{}, exceptions.InvalidInput("RoutineTask").WithOrigin(err)
	}

	queryRangeStartedAt, err := time.Parse(time.RFC3339, ctx.Query("queryRangeStartedAt"))
	if err != nil {
		return 0, time.Time{}, time.Time{}, exceptions.InvalidInput("RoutineTask").WithOrigin(fmt.Errorf("queryRangeStartedAt must be an RFC3339 timestamp: %w", err))
	}
	queryRangeEndedAt, err := time.Parse(time.RFC3339, ctx.Query("queryRangeEndedAt"))
	if err != nil {
		return 0, time.Time{}, time.Time{}, exceptions.InvalidInput("RoutineTask").WithOrigin(fmt.Errorf("queryRangeEndedAt must be an RFC3339 timestamp: %w", err))
	}

	return timeHourUnit, queryRangeStartedAt, queryRangeEndedAt, nil
}

/* ============================== Service Methods for RoutineTaskRecord ============================== */

func (b *RoutineTaskRecordBinder) BindGetAllMyRoutineTaskRecordsByRoutineTaskId(controllerFunc apitransport.ControllerFunc[*routinetaskrecordsdto.GetAllMyRoutineTaskRecordsByRoutineTaskIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &routinetaskrecordsdto.GetAllMyRoutineTaskRecordsByRoutineTaskIdRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		routineTaskId, err := uuid.Parse(ctx.Param("routineTaskId"))
		if err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("RoutineTask").WithOrigin(err), ctx)
			return
		}
		requestDto.Param.RoutineTaskId = routineTaskId

		limitString := ctx.Query("limit")
		if limitString != "" {
			limit, err := strconv.Atoi(limitString)
			if err != nil {
				responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("RoutineTask").WithOrigin(err), ctx)
				return
			}
			requestDto.Param.Limit = limit
		}

		controllerFunc(ctx, requestDto)
	}
}

/* ============================== Visualization Methods ============================== */

func (b *RoutineTaskRecordBinder) BindVisualizeMyRoutineTaskRecordStatusCount(controllerFunc apitransport.ControllerFunc[*routinetaskrecordsdto.VisualizeMyRoutineTaskRecordStatusCountRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &routinetaskrecordsdto.VisualizeMyRoutineTaskRecordStatusCountRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")
		params, exception := parseRoutineTaskRecordVisualizationParams(ctx)
		if exception != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
			return
		}
		requestDto.Param.Permission = params.permission
		requestDto.Param.RoutineTaskIds = params.routineTaskIds
		controllerFunc(ctx, requestDto)
	}
}

func (b *RoutineTaskRecordBinder) BindVisualizeMyRoutineTaskRecordPurposeCount(controllerFunc apitransport.ControllerFunc[*routinetaskrecordsdto.VisualizeMyRoutineTaskRecordPurposeCountRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &routinetaskrecordsdto.VisualizeMyRoutineTaskRecordPurposeCountRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")
		params, exception := parseRoutineTaskRecordVisualizationParams(ctx)
		if exception != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
			return
		}
		requestDto.Param.Permission = params.permission
		requestDto.Param.RoutineTaskIds = params.routineTaskIds
		controllerFunc(ctx, requestDto)
	}
}

func (b *RoutineTaskRecordBinder) BindVisualizeMyRoutineTaskRecordScheduledAtCount(controllerFunc apitransport.ControllerFunc[*routinetaskrecordsdto.VisualizeMyRoutineTaskRecordScheduledAtCountRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &routinetaskrecordsdto.VisualizeMyRoutineTaskRecordScheduledAtCountRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")
		params, exception := parseRoutineTaskRecordVisualizationParams(ctx)
		if exception != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
			return
		}
		timeHourUnit, queryRangeStartedAt, queryRangeEndedAt, exception := parseRoutineTaskRecordVisualizationTimeRange(ctx)
		if exception != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
			return
		}
		requestDto.Param.Permission = params.permission
		requestDto.Param.RoutineTaskIds = params.routineTaskIds
		requestDto.Param.TimeHourUnit = timeHourUnit
		requestDto.Param.QueryRangeStartedAt = queryRangeStartedAt
		requestDto.Param.QueryRangeEndedAt = queryRangeEndedAt
		controllerFunc(ctx, requestDto)
	}
}

func (b *RoutineTaskRecordBinder) BindVisualizeMyRoutineTaskRecordActualStartedAtCount(controllerFunc apitransport.ControllerFunc[*routinetaskrecordsdto.VisualizeMyRoutineTaskRecordActualStartedAtCountRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &routinetaskrecordsdto.VisualizeMyRoutineTaskRecordActualStartedAtCountRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")
		params, exception := parseRoutineTaskRecordVisualizationParams(ctx)
		if exception != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
			return
		}
		timeHourUnit, queryRangeStartedAt, queryRangeEndedAt, exception := parseRoutineTaskRecordVisualizationTimeRange(ctx)
		if exception != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
			return
		}
		requestDto.Param.Permission = params.permission
		requestDto.Param.RoutineTaskIds = params.routineTaskIds
		requestDto.Param.TimeHourUnit = timeHourUnit
		requestDto.Param.QueryRangeStartedAt = queryRangeStartedAt
		requestDto.Param.QueryRangeEndedAt = queryRangeEndedAt
		controllerFunc(ctx, requestDto)
	}
}

func (b *RoutineTaskRecordBinder) BindVisualizeMyRoutineTaskRecordActualEndedAtCount(controllerFunc apitransport.ControllerFunc[*routinetaskrecordsdto.VisualizeMyRoutineTaskRecordActualEndedAtCountRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &routinetaskrecordsdto.VisualizeMyRoutineTaskRecordActualEndedAtCountRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")
		params, exception := parseRoutineTaskRecordVisualizationParams(ctx)
		if exception != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
			return
		}
		timeHourUnit, queryRangeStartedAt, queryRangeEndedAt, exception := parseRoutineTaskRecordVisualizationTimeRange(ctx)
		if exception != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
			return
		}
		requestDto.Param.Permission = params.permission
		requestDto.Param.RoutineTaskIds = params.routineTaskIds
		requestDto.Param.TimeHourUnit = timeHourUnit
		requestDto.Param.QueryRangeStartedAt = queryRangeStartedAt
		requestDto.Param.QueryRangeEndedAt = queryRangeEndedAt
		controllerFunc(ctx, requestDto)
	}
}
