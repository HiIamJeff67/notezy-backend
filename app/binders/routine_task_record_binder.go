package binders

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	contexts "github.com/HiIamJeff67/notezy-backend/app/contexts"
	dtos "github.com/HiIamJeff67/notezy-backend/app/dtos"
	exceptions "github.com/HiIamJeff67/notezy-backend/app/exceptions"
	enums "github.com/HiIamJeff67/notezy-backend/app/models/schemas/enums"
	types "github.com/HiIamJeff67/notezy-backend/shared/types"
)

type RoutineTaskRecordBinderInterface interface {
	BindGetAllMyRoutineTaskRecordsByRoutineTaskId(controllerFunc types.ControllerFunc[*dtos.GetAllMyRoutineTaskRecordsByRoutineTaskIdReqDto]) gin.HandlerFunc
	BindVisualizeMyRoutineTaskRecordStatusCount(controllerFunc types.ControllerFunc[*dtos.VisualizeMyRoutineTaskRecordStatusCountReqDto]) gin.HandlerFunc
	BindVisualizeMyRoutineTaskRecordPurposeCount(controllerFunc types.ControllerFunc[*dtos.VisualizeMyRoutineTaskRecordPurposeCountReqDto]) gin.HandlerFunc
	BindVisualizeMyRoutineTaskRecordScheduledAtCount(controllerFunc types.ControllerFunc[*dtos.VisualizeMyRoutineTaskRecordScheduledAtCountReqDto]) gin.HandlerFunc
	BindVisualizeMyRoutineTaskRecordActualStartedAtCount(controllerFunc types.ControllerFunc[*dtos.VisualizeMyRoutineTaskRecordActualStartedAtCountReqDto]) gin.HandlerFunc
	BindVisualizeMyRoutineTaskRecordActualEndedAtCount(controllerFunc types.ControllerFunc[*dtos.VisualizeMyRoutineTaskRecordActualEndedAtCountReqDto]) gin.HandlerFunc
}

type RoutineTaskRecordBinder struct{}

func NewRoutineTaskRecordBinder() RoutineTaskRecordBinderInterface {
	return &RoutineTaskRecordBinder{}
}

func (b *RoutineTaskRecordBinder) BindGetAllMyRoutineTaskRecordsByRoutineTaskId(
	controllerFunc types.ControllerFunc[*dtos.GetAllMyRoutineTaskRecordsByRoutineTaskIdReqDto],
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.GetAllMyRoutineTaskRecordsByRoutineTaskIdReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		routineTaskIdString := ctx.Query("routineTaskId")
		if routineTaskIdString == "" {
			exceptions.RoutineTask.InvalidInput().WithOrigin(fmt.Errorf("routineTaskId is required")).SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		routineTaskId, err := uuid.Parse(routineTaskIdString)
		if err != nil {
			exceptions.RoutineTask.InvalidInput().WithOrigin(err).SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.Param.RoutineTaskId = routineTaskId

		limitString := ctx.Query("limit")
		if limitString != "" {
			limit, err := strconv.Atoi(limitString)
			if err != nil {
				exceptions.RoutineTask.InvalidInput().WithOrigin(err).SafelyAbortAndResponseWithJSON(ctx)
				return
			}
			reqDto.Param.Limit = limit
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *RoutineTaskRecordBinder) BindVisualizeMyRoutineTaskRecordStatusCount(
	controllerFunc types.ControllerFunc[*dtos.VisualizeMyRoutineTaskRecordStatusCountReqDto],
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.VisualizeMyRoutineTaskRecordStatusCountReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		permissionString := ctx.Query("permission")
		if permissionString == "" {
			exceptions.RoutineTask.InvalidInput().WithOrigin(fmt.Errorf("permission is required")).SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		permission, err := enums.ConvertStringToAccessControlPermission(permissionString)
		if err != nil {
			exceptions.RoutineTask.InvalidInput().WithOrigin(err).SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.Param.Permission = *permission

		for _, routineTaskIdsValue := range ctx.QueryArray("routineTaskIds") {
			for _, routineTaskIdValue := range strings.Split(routineTaskIdsValue, ",") {
				trimmedRoutineTaskIdValue := strings.TrimSpace(routineTaskIdValue)
				if trimmedRoutineTaskIdValue == "" {
					continue
				}
				routineTaskId, err := uuid.Parse(trimmedRoutineTaskIdValue)
				if err != nil {
					exceptions.RoutineTask.InvalidInput().WithOrigin(err).SafelyAbortAndResponseWithJSON(ctx)
					return
				}
				reqDto.Param.RoutineTaskIds = append(reqDto.Param.RoutineTaskIds, routineTaskId)
			}
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *RoutineTaskRecordBinder) BindVisualizeMyRoutineTaskRecordPurposeCount(
	controllerFunc types.ControllerFunc[*dtos.VisualizeMyRoutineTaskRecordPurposeCountReqDto],
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.VisualizeMyRoutineTaskRecordPurposeCountReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		permissionString := ctx.Query("permission")
		if permissionString == "" {
			exceptions.RoutineTask.InvalidInput().WithOrigin(fmt.Errorf("permission is required")).SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		permission, err := enums.ConvertStringToAccessControlPermission(permissionString)
		if err != nil {
			exceptions.RoutineTask.InvalidInput().WithOrigin(err).SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.Param.Permission = *permission

		for _, routineTaskIdsValue := range ctx.QueryArray("routineTaskIds") {
			for _, routineTaskIdValue := range strings.Split(routineTaskIdsValue, ",") {
				trimmedRoutineTaskIdValue := strings.TrimSpace(routineTaskIdValue)
				if trimmedRoutineTaskIdValue == "" {
					continue
				}
				routineTaskId, err := uuid.Parse(trimmedRoutineTaskIdValue)
				if err != nil {
					exceptions.RoutineTask.InvalidInput().WithOrigin(err).SafelyAbortAndResponseWithJSON(ctx)
					return
				}
				reqDto.Param.RoutineTaskIds = append(reqDto.Param.RoutineTaskIds, routineTaskId)
			}
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *RoutineTaskRecordBinder) BindVisualizeMyRoutineTaskRecordScheduledAtCount(
	controllerFunc types.ControllerFunc[*dtos.VisualizeMyRoutineTaskRecordScheduledAtCountReqDto],
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.VisualizeMyRoutineTaskRecordScheduledAtCountReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		permissionString := ctx.Query("permission")
		if permissionString == "" {
			exceptions.RoutineTask.InvalidInput().WithOrigin(fmt.Errorf("permission is required")).SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		permission, err := enums.ConvertStringToAccessControlPermission(permissionString)
		if err != nil {
			exceptions.RoutineTask.InvalidInput().WithOrigin(err).SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.Param.Permission = *permission

		for _, routineTaskIdsValue := range ctx.QueryArray("routineTaskIds") {
			for _, routineTaskIdValue := range strings.Split(routineTaskIdsValue, ",") {
				trimmedRoutineTaskIdValue := strings.TrimSpace(routineTaskIdValue)
				if trimmedRoutineTaskIdValue == "" {
					continue
				}
				routineTaskId, err := uuid.Parse(trimmedRoutineTaskIdValue)
				if err != nil {
					exceptions.RoutineTask.InvalidInput().WithOrigin(err).SafelyAbortAndResponseWithJSON(ctx)
					return
				}
				reqDto.Param.RoutineTaskIds = append(reqDto.Param.RoutineTaskIds, routineTaskId)
			}
		}

		timeHourUnitString := ctx.Query("timeHourUnit")
		if timeHourUnitString == "" {
			exceptions.RoutineTask.InvalidInput().WithOrigin(fmt.Errorf("timeHourUnit is required")).SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		timeHourUnit, err := strconv.Atoi(timeHourUnitString)
		if err != nil {
			exceptions.RoutineTask.InvalidInput().WithOrigin(err).SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.Param.TimeHourUnit = timeHourUnit

		queryRangeStartedAt, err := time.Parse(time.RFC3339, ctx.Query("queryRangeStartedAt"))
		if err != nil {
			exceptions.RoutineTask.InvalidInput().WithOrigin(fmt.Errorf("queryRangeStartedAt must be an RFC3339 timestamp: %w", err)).SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		queryRangeEndedAt, err := time.Parse(time.RFC3339, ctx.Query("queryRangeEndedAt"))
		if err != nil {
			exceptions.RoutineTask.InvalidInput().WithOrigin(fmt.Errorf("queryRangeEndedAt must be an RFC3339 timestamp: %w", err)).SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.Param.QueryRangeStartedAt = queryRangeStartedAt
		reqDto.Param.QueryRangeEndedAt = queryRangeEndedAt

		controllerFunc(ctx, &reqDto)
	}
}

func (b *RoutineTaskRecordBinder) BindVisualizeMyRoutineTaskRecordActualStartedAtCount(
	controllerFunc types.ControllerFunc[*dtos.VisualizeMyRoutineTaskRecordActualStartedAtCountReqDto],
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.VisualizeMyRoutineTaskRecordActualStartedAtCountReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		permissionString := ctx.Query("permission")
		if permissionString == "" {
			exceptions.RoutineTask.InvalidInput().WithOrigin(fmt.Errorf("permission is required")).SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		permission, err := enums.ConvertStringToAccessControlPermission(permissionString)
		if err != nil {
			exceptions.RoutineTask.InvalidInput().WithOrigin(err).SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.Param.Permission = *permission

		for _, routineTaskIdsValue := range ctx.QueryArray("routineTaskIds") {
			for _, routineTaskIdValue := range strings.Split(routineTaskIdsValue, ",") {
				trimmedRoutineTaskIdValue := strings.TrimSpace(routineTaskIdValue)
				if trimmedRoutineTaskIdValue == "" {
					continue
				}
				routineTaskId, err := uuid.Parse(trimmedRoutineTaskIdValue)
				if err != nil {
					exceptions.RoutineTask.InvalidInput().WithOrigin(err).SafelyAbortAndResponseWithJSON(ctx)
					return
				}
				reqDto.Param.RoutineTaskIds = append(reqDto.Param.RoutineTaskIds, routineTaskId)
			}
		}

		timeHourUnitString := ctx.Query("timeHourUnit")
		if timeHourUnitString == "" {
			exceptions.RoutineTask.InvalidInput().WithOrigin(fmt.Errorf("timeHourUnit is required")).SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		timeHourUnit, err := strconv.Atoi(timeHourUnitString)
		if err != nil {
			exceptions.RoutineTask.InvalidInput().WithOrigin(err).SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.Param.TimeHourUnit = timeHourUnit

		queryRangeStartedAt, err := time.Parse(time.RFC3339, ctx.Query("queryRangeStartedAt"))
		if err != nil {
			exceptions.RoutineTask.InvalidInput().WithOrigin(fmt.Errorf("queryRangeStartedAt must be an RFC3339 timestamp: %w", err)).SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		queryRangeEndedAt, err := time.Parse(time.RFC3339, ctx.Query("queryRangeEndedAt"))
		if err != nil {
			exceptions.RoutineTask.InvalidInput().WithOrigin(fmt.Errorf("queryRangeEndedAt must be an RFC3339 timestamp: %w", err)).SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.Param.QueryRangeStartedAt = queryRangeStartedAt
		reqDto.Param.QueryRangeEndedAt = queryRangeEndedAt

		controllerFunc(ctx, &reqDto)
	}
}

func (b *RoutineTaskRecordBinder) BindVisualizeMyRoutineTaskRecordActualEndedAtCount(
	controllerFunc types.ControllerFunc[*dtos.VisualizeMyRoutineTaskRecordActualEndedAtCountReqDto],
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.VisualizeMyRoutineTaskRecordActualEndedAtCountReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		permissionString := ctx.Query("permission")
		if permissionString == "" {
			exceptions.RoutineTask.InvalidInput().WithOrigin(fmt.Errorf("permission is required")).SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		permission, err := enums.ConvertStringToAccessControlPermission(permissionString)
		if err != nil {
			exceptions.RoutineTask.InvalidInput().WithOrigin(err).SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.Param.Permission = *permission

		for _, routineTaskIdsValue := range ctx.QueryArray("routineTaskIds") {
			for _, routineTaskIdValue := range strings.Split(routineTaskIdsValue, ",") {
				trimmedRoutineTaskIdValue := strings.TrimSpace(routineTaskIdValue)
				if trimmedRoutineTaskIdValue == "" {
					continue
				}
				routineTaskId, err := uuid.Parse(trimmedRoutineTaskIdValue)
				if err != nil {
					exceptions.RoutineTask.InvalidInput().WithOrigin(err).SafelyAbortAndResponseWithJSON(ctx)
					return
				}
				reqDto.Param.RoutineTaskIds = append(reqDto.Param.RoutineTaskIds, routineTaskId)
			}
		}

		timeHourUnitString := ctx.Query("timeHourUnit")
		if timeHourUnitString == "" {
			exceptions.RoutineTask.InvalidInput().WithOrigin(fmt.Errorf("timeHourUnit is required")).SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		timeHourUnit, err := strconv.Atoi(timeHourUnitString)
		if err != nil {
			exceptions.RoutineTask.InvalidInput().WithOrigin(err).SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.Param.TimeHourUnit = timeHourUnit

		queryRangeStartedAt, err := time.Parse(time.RFC3339, ctx.Query("queryRangeStartedAt"))
		if err != nil {
			exceptions.RoutineTask.InvalidInput().WithOrigin(fmt.Errorf("queryRangeStartedAt must be an RFC3339 timestamp: %w", err)).SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		queryRangeEndedAt, err := time.Parse(time.RFC3339, ctx.Query("queryRangeEndedAt"))
		if err != nil {
			exceptions.RoutineTask.InvalidInput().WithOrigin(fmt.Errorf("queryRangeEndedAt must be an RFC3339 timestamp: %w", err)).SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.Param.QueryRangeStartedAt = queryRangeStartedAt
		reqDto.Param.QueryRangeEndedAt = queryRangeEndedAt

		controllerFunc(ctx, &reqDto)
	}
}
