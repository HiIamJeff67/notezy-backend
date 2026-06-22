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

type RoutineTaskBinderInterface interface {
	BindGetMyRoutineTaskById(controllerFunc types.ControllerFunc[*dtos.GetMyRoutineTaskByIdReqDto]) gin.HandlerFunc
	BindGetAllMyRoutineTasksByStationIds(controllerFunc types.ControllerFunc[*dtos.GetAllMyRoutineTasksByStationIdsReqDto]) gin.HandlerFunc
	BindGetAllMyRoutineTasks(controllerFunc types.ControllerFunc[*dtos.GetAllMyRoutineTasksReqDto]) gin.HandlerFunc
	BindCreateRoutineTaskByStationId(controllerFunc types.ControllerFunc[*dtos.CreateRoutineTaskByStationIdReqDto]) gin.HandlerFunc
	BindUpdateMyRoutineTaskById(controllerFunc types.ControllerFunc[*dtos.UpdateMyRoutineTaskByIdReqDto]) gin.HandlerFunc
	BindHardDeleteMyRoutineTaskById(controllerFunc types.ControllerFunc[*dtos.HardDeleteMyRoutineTaskByIdReqDto]) gin.HandlerFunc
	BindHardDeleteMyRoutineTasksByIds(controllerFunc types.ControllerFunc[*dtos.HardDeleteMyRoutineTasksByIdsReqDto]) gin.HandlerFunc
	BindVisualizeMyRoutineTaskStatusCount(controllerFunc types.ControllerFunc[*dtos.VisualizeMyRoutineTaskStatusCountReqDto]) gin.HandlerFunc
	BindVisualizeMyRoutineTaskPurposeCount(controllerFunc types.ControllerFunc[*dtos.VisualizeMyRoutineTaskPurposeCountReqDto]) gin.HandlerFunc
	BindVisualizeMyRoutineTaskScheduledAtCount(controllerFunc types.ControllerFunc[*dtos.VisualizeMyRoutineTaskScheduledAtCountReqDto]) gin.HandlerFunc
	BindVisualizeMyRoutineTaskActualStartedAtCount(controllerFunc types.ControllerFunc[*dtos.VisualizeMyRoutineTaskActualStartedAtCountReqDto]) gin.HandlerFunc
	BindVisualizeMyRoutineTaskActualEndedAtCount(controllerFunc types.ControllerFunc[*dtos.VisualizeMyRoutineTaskActualEndedAtCountReqDto]) gin.HandlerFunc
}

type RoutineTaskBinder struct{}

func NewRoutineTaskBinder() RoutineTaskBinderInterface {
	return &RoutineTaskBinder{}
}

func (b *RoutineTaskBinder) BindGetMyRoutineTaskById(controllerFunc types.ControllerFunc[*dtos.GetMyRoutineTaskByIdReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.GetMyRoutineTaskByIdReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		isDeletedString := ctx.Query("isDeleted")
		if isDeletedString != "" {
			isDeleted, err := strconv.ParseBool(isDeletedString)
			if err != nil {
				exceptions.RoutineTask.InvalidInput().WithOrigin(err).SafelyAbortAndResponseWithJSON(ctx)
				return
			}
			reqDto.Param.IsDeleted = &isDeleted
		}

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

		controllerFunc(ctx, &reqDto)
	}
}

func (b *RoutineTaskBinder) BindGetAllMyRoutineTasksByStationIds(
	controllerFunc types.ControllerFunc[*dtos.GetAllMyRoutineTasksByStationIdsReqDto],
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.GetAllMyRoutineTasksByStationIdsReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		areDeletedString := ctx.Query("areDeleted")
		if areDeletedString != "" {
			areDeleted, err := strconv.ParseBool(areDeletedString)
			if err != nil {
				exceptions.RoutineTask.InvalidInput().WithOrigin(err).SafelyAbortAndResponseWithJSON(ctx)
				return
			}
			reqDto.Param.AreDeleted = &areDeleted
		}

		for _, stationIdsValue := range ctx.QueryArray("stationIds") {
			for _, stationIdValue := range strings.Split(stationIdsValue, ",") {
				stationId, err := uuid.Parse(strings.TrimSpace(stationIdValue))
				if err != nil {
					exceptions.RoutineTask.InvalidInput().WithOrigin(err).SafelyAbortAndResponseWithJSON(ctx)
					return
				}
				reqDto.Param.StationIds = append(reqDto.Param.StationIds, stationId)
			}
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *RoutineTaskBinder) BindGetAllMyRoutineTasks(controllerFunc types.ControllerFunc[*dtos.GetAllMyRoutineTasksReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.GetAllMyRoutineTasksReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		areDeletedString := ctx.Query("areDeleted")
		if areDeletedString != "" {
			areDeleted, err := strconv.ParseBool(areDeletedString)
			if err != nil {
				exceptions.RoutineTask.InvalidInput().WithOrigin(err).SafelyAbortAndResponseWithJSON(ctx)
				return
			}
			reqDto.Param.AreDeleted = &areDeleted
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *RoutineTaskBinder) BindCreateRoutineTaskByStationId(controllerFunc types.ControllerFunc[*dtos.CreateRoutineTaskByStationIdReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.CreateRoutineTaskByStationIdReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.RoutineTask.InvalidDto().WithOrigin(err)
			exception.SafelyAbortAndResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *RoutineTaskBinder) BindUpdateMyRoutineTaskById(controllerFunc types.ControllerFunc[*dtos.UpdateMyRoutineTaskByIdReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.UpdateMyRoutineTaskByIdReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.RoutineTask.InvalidDto().WithOrigin(err)
			exception.SafelyAbortAndResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *RoutineTaskBinder) BindHardDeleteMyRoutineTaskById(controllerFunc types.ControllerFunc[*dtos.HardDeleteMyRoutineTaskByIdReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.HardDeleteMyRoutineTaskByIdReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.RoutineTask.InvalidDto().WithOrigin(err)
			exception.SafelyAbortAndResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *RoutineTaskBinder) BindHardDeleteMyRoutineTasksByIds(controllerFunc types.ControllerFunc[*dtos.HardDeleteMyRoutineTasksByIdsReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.HardDeleteMyRoutineTasksByIdsReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.RoutineTask.InvalidDto().WithOrigin(err)
			exception.SafelyAbortAndResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *RoutineTaskBinder) BindVisualizeMyRoutineTaskStatusCount(
	controllerFunc types.ControllerFunc[*dtos.VisualizeMyRoutineTaskStatusCountReqDto],
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.VisualizeMyRoutineTaskStatusCountReqDto

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

		controllerFunc(ctx, &reqDto)
	}
}

func (b *RoutineTaskBinder) BindVisualizeMyRoutineTaskPurposeCount(
	controllerFunc types.ControllerFunc[*dtos.VisualizeMyRoutineTaskPurposeCountReqDto],
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.VisualizeMyRoutineTaskPurposeCountReqDto

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

		controllerFunc(ctx, &reqDto)
	}
}

func (b *RoutineTaskBinder) BindVisualizeMyRoutineTaskScheduledAtCount(
	controllerFunc types.ControllerFunc[*dtos.VisualizeMyRoutineTaskScheduledAtCountReqDto],
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.VisualizeMyRoutineTaskScheduledAtCountReqDto

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

func (b *RoutineTaskBinder) BindVisualizeMyRoutineTaskActualStartedAtCount(
	controllerFunc types.ControllerFunc[*dtos.VisualizeMyRoutineTaskActualStartedAtCountReqDto],
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.VisualizeMyRoutineTaskActualStartedAtCountReqDto

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

func (b *RoutineTaskBinder) BindVisualizeMyRoutineTaskActualEndedAtCount(
	controllerFunc types.ControllerFunc[*dtos.VisualizeMyRoutineTaskActualEndedAtCountReqDto],
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.VisualizeMyRoutineTaskActualEndedAtCountReqDto

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
