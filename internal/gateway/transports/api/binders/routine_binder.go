package binders

import (
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	routinesdto "github.com/HiIamJeff67/notezy-backend/contracts/api/v1/routines"
	exceptions "github.com/HiIamJeff67/notezy-backend/internal/exceptions"
	responsewriter "github.com/HiIamJeff67/notezy-backend/internal/gateway/responsewriter"
	apitransport "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/api"
	enums "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/schemas/enums"
	sharedtypes "github.com/HiIamJeff67/notezy-backend/internal/shared/types"
)

type RoutineBinderInterface interface {
	BindGetMyRoutineById(controllerFunc apitransport.ControllerFunc[*routinesdto.GetMyRoutineByIdRequestDto]) gin.HandlerFunc
	BindGetMyRoutinesByStationId(controllerFunc apitransport.ControllerFunc[*routinesdto.GetMyRoutinesByStationIdRequestDto]) gin.HandlerFunc
	BindGetAllMyRoutinesByTimeRange(controllerFunc apitransport.ControllerFunc[*routinesdto.GetAllMyRoutinesByTimeRangeRequestDto]) gin.HandlerFunc
	BindCreateRoutineByStationId(controllerFunc apitransport.ControllerFunc[*routinesdto.CreateRoutineByStationIdRequestDto]) gin.HandlerFunc
	BindCreateRoutinesByStationIds(controllerFunc apitransport.ControllerFunc[*routinesdto.CreateRoutinesByStationIdsRequestDto]) gin.HandlerFunc
	BindUpdateMyRoutineById(controllerFunc apitransport.ControllerFunc[*routinesdto.UpdateMyRoutineByIdRequestDto]) gin.HandlerFunc
	BindUpdateMyRoutinesByIds(controllerFunc apitransport.ControllerFunc[*routinesdto.UpdateMyRoutinesByIdsRequestDto]) gin.HandlerFunc
	BindLinkRoutineTagById(controllerFunc apitransport.ControllerFunc[*routinesdto.LinkRoutineTagByIdRequestDto]) gin.HandlerFunc
	BindLinkRoutineTagsByIds(controllerFunc apitransport.ControllerFunc[*routinesdto.LinkRoutineTagsByIdsRequestDto]) gin.HandlerFunc
	BindLinkRoutineItemById(controllerFunc apitransport.ControllerFunc[*routinesdto.LinkRoutineItemByIdRequestDto]) gin.HandlerFunc
	BindLinkRoutineItemsByIds(controllerFunc apitransport.ControllerFunc[*routinesdto.LinkRoutineItemsByIdsRequestDto]) gin.HandlerFunc
	BindRestoreMyRoutineById(controllerFunc apitransport.ControllerFunc[*routinesdto.RestoreMyRoutineByIdRequestDto]) gin.HandlerFunc
	BindRestoreMyRoutinesByIds(controllerFunc apitransport.ControllerFunc[*routinesdto.RestoreMyRoutinesByIdsRequestDto]) gin.HandlerFunc
	BindDeleteMyRoutineById(controllerFunc apitransport.ControllerFunc[*routinesdto.DeleteMyRoutineByIdRequestDto]) gin.HandlerFunc
	BindDeleteMyRoutinesByIds(controllerFunc apitransport.ControllerFunc[*routinesdto.DeleteMyRoutinesByIdsRequestDto]) gin.HandlerFunc
	BindHardDeleteMyRoutineById(controllerFunc apitransport.ControllerFunc[*routinesdto.HardDeleteMyRoutineByIdRequestDto]) gin.HandlerFunc
	BindHardDeleteMyRoutinesByIds(controllerFunc apitransport.ControllerFunc[*routinesdto.HardDeleteMyRoutinesByIdsRequestDto]) gin.HandlerFunc

	/* ============================== Visualization Methods ============================== */
	BindVisualizeMyRoutineStatusCount(controllerFunc apitransport.ControllerFunc[*routinesdto.VisualizeMyRoutineStatusCountRequestDto]) gin.HandlerFunc
	BindVisualizeMyRoutinePeriodCount(controllerFunc apitransport.ControllerFunc[*routinesdto.VisualizeMyRoutinePeriodCountRequestDto]) gin.HandlerFunc
	BindVisualizeMyRoutineScheduledStartAtCount(controllerFunc apitransport.ControllerFunc[*routinesdto.VisualizeMyRoutineScheduledStartAtCountRequestDto]) gin.HandlerFunc
	BindVisualizeMyRoutineScheduledEndAtCount(controllerFunc apitransport.ControllerFunc[*routinesdto.VisualizeMyRoutineScheduledEndAtCountRequestDto]) gin.HandlerFunc
}

type RoutineBinder struct{}

func NewRoutineBinder() RoutineBinderInterface { return &RoutineBinder{} }

func bindRoutineJSON[T any](ctx *gin.Context, requestDto *T, body any, controllerFunc apitransport.ControllerFunc[*T]) {
	if err := ctx.ShouldBindJSON(body); err != nil {
		responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("Routine").WithOrigin(err), ctx)
		return
	}
	controllerFunc(ctx, requestDto)
}

func parseRoutineUUID(ctx *gin.Context, name string) (uuid.UUID, bool) {
	value, err := uuid.Parse(ctx.Param(name))
	if err != nil {
		responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Routine").WithOrigin(err), ctx)
		return uuid.Nil, false
	}
	return value, true
}

func parseRoutineBool(ctx *gin.Context, name string) (*bool, bool) {
	valueString := ctx.Query(name)
	if valueString == "" {
		return nil, true
	}
	value, err := strconv.ParseBool(valueString)
	if err != nil {
		responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Routine").WithOrigin(err), ctx)
		return nil, false
	}
	return &value, true
}

func parseRoutinePermission(ctx *gin.Context) (sharedtypes.AccessControlPermission, bool) {
	value, err := sharedtypes.ParseAccessControlPermission(ctx.Query("permission"))
	if err != nil {
		responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Routine").WithOrigin(err), ctx)
		return "", false
	}
	return *value, true
}

func parseRoutineTime(ctx *gin.Context, name string) (time.Time, bool) {
	value, err := time.Parse(time.RFC3339, ctx.Query(name))
	if err != nil {
		responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Routine").WithOrigin(err), ctx)
		return time.Time{}, false
	}
	return value, true
}

func parseRoutineUUIDs(ctx *gin.Context, name string) ([]uuid.UUID, bool) {
	values := ctx.QueryArray(name)
	if len(values) == 1 {
		values = strings.Split(values[0], ",")
	}
	ids := make([]uuid.UUID, len(values))
	for index, value := range values {
		parsed, err := uuid.Parse(value)
		if err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Routine").WithOrigin(err), ctx)
			return nil, false
		}
		ids[index] = parsed
	}
	return ids, true
}

func (b *RoutineBinder) BindGetMyRoutineById(controllerFunc apitransport.ControllerFunc[*routinesdto.GetMyRoutineByIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &routinesdto.GetMyRoutineByIdRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")
		isDeleted, ok := parseRoutineBool(ctx, "isDeleted")
		if !ok {
			return
		}
		requestDto.Param.IsDeleted = isDeleted
		value, ok := parseRoutineUUID(ctx, "routineId")
		if !ok {
			return
		}
		requestDto.Param.RoutineId = value
		controllerFunc(ctx, requestDto)
		return
	}
}

func (b *RoutineBinder) BindGetMyRoutinesByStationId(controllerFunc apitransport.ControllerFunc[*routinesdto.GetMyRoutinesByStationIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &routinesdto.GetMyRoutinesByStationIdRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")
		areDeleted, ok := parseRoutineBool(ctx, "areDeleted")
		if !ok {
			return
		}
		requestDto.Param.AreDeleted = areDeleted
		value, ok := parseRoutineUUID(ctx, "stationId")
		if !ok {
			return
		}
		requestDto.Param.StationId = value
		controllerFunc(ctx, requestDto)
		return
	}
}

func (b *RoutineBinder) BindGetAllMyRoutinesByTimeRange(controllerFunc apitransport.ControllerFunc[*routinesdto.GetAllMyRoutinesByTimeRangeRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &routinesdto.GetAllMyRoutinesByTimeRangeRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")
		ok := true
		requestDto.Param.AreDeleted, ok = parseRoutineBool(ctx, "areDeleted")
		if requestDto.Param.AreDeleted == nil && ctx.Query("areDeleted") != "" {
			return
		}
		requestDto.Param.From, ok = parseRoutineTime(ctx, "from")
		if !ok {
			return
		}
		requestDto.Param.To, ok = parseRoutineTime(ctx, "to")
		if !ok {
			return
		}
		requestDto.Param.StationIds, ok = parseRoutineUUIDs(ctx, "stationIds")
		if !ok {
			return
		}
		controllerFunc(ctx, requestDto)
		return
	}
}

func (b *RoutineBinder) BindCreateRoutineByStationId(controllerFunc apitransport.ControllerFunc[*routinesdto.CreateRoutineByStationIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &routinesdto.CreateRoutineByStationIdRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")
		value, ok := parseRoutineUUID(ctx, "stationId")
		if !ok {
			return
		}
		requestDto.Body.StationId = value
		bindRoutineJSON(ctx, requestDto, &requestDto.Body, controllerFunc)
		return
	}
}

func (b *RoutineBinder) BindCreateRoutinesByStationIds(controllerFunc apitransport.ControllerFunc[*routinesdto.CreateRoutinesByStationIdsRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &routinesdto.CreateRoutinesByStationIdsRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")
		bindRoutineJSON(ctx, requestDto, &requestDto.Body, controllerFunc)
		return
	}
}

func (b *RoutineBinder) BindUpdateMyRoutineById(controllerFunc apitransport.ControllerFunc[*routinesdto.UpdateMyRoutineByIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &routinesdto.UpdateMyRoutineByIdRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")
		value, ok := parseRoutineUUID(ctx, "routineId")
		if !ok {
			return
		}
		requestDto.Body.RoutineId = value
		bindRoutineJSON(ctx, requestDto, &requestDto.Body, controllerFunc)
		return
	}
}

func (b *RoutineBinder) BindUpdateMyRoutinesByIds(controllerFunc apitransport.ControllerFunc[*routinesdto.UpdateMyRoutinesByIdsRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &routinesdto.UpdateMyRoutinesByIdsRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")
		bindRoutineJSON(ctx, requestDto, &requestDto.Body, controllerFunc)
		return
	}
}

func (b *RoutineBinder) BindLinkRoutineTagById(controllerFunc apitransport.ControllerFunc[*routinesdto.LinkRoutineTagByIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &routinesdto.LinkRoutineTagByIdRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")
		value, ok := parseRoutineUUID(ctx, "routineId")
		if !ok {
			return
		}
		requestDto.Body.RoutineId = value
		value, ok = parseRoutineUUID(ctx, "routineTagId")
		if !ok {
			return
		}
		requestDto.Body.RoutineTagId = value
		requestDto.Body.IsUnlink = ctx.Query("isUnlink") == "true"
		controllerFunc(ctx, requestDto)
		return
	}
}

func (b *RoutineBinder) BindLinkRoutineTagsByIds(controllerFunc apitransport.ControllerFunc[*routinesdto.LinkRoutineTagsByIdsRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &routinesdto.LinkRoutineTagsByIdsRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")
		bindRoutineJSON(ctx, requestDto, &requestDto.Body, controllerFunc)
		return
	}
}

func (b *RoutineBinder) BindLinkRoutineItemById(controllerFunc apitransport.ControllerFunc[*routinesdto.LinkRoutineItemByIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &routinesdto.LinkRoutineItemByIdRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")
		value, ok := parseRoutineUUID(ctx, "routineId")
		if !ok {
			return
		}
		requestDto.Body.RoutineId = value
		value, ok = parseRoutineUUID(ctx, "itemId")
		if !ok {
			return
		}
		requestDto.Body.ItemId = value
		requestDto.Body.ItemType = enums.ItemType(ctx.Query("itemType"))
		requestDto.Body.IsUnlink = ctx.Query("isUnlink") == "true"
		controllerFunc(ctx, requestDto)
		return
	}
}

func (b *RoutineBinder) BindLinkRoutineItemsByIds(controllerFunc apitransport.ControllerFunc[*routinesdto.LinkRoutineItemsByIdsRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &routinesdto.LinkRoutineItemsByIdsRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")
		bindRoutineJSON(ctx, requestDto, &requestDto.Body, controllerFunc)
		return
	}
}

func (b *RoutineBinder) BindRestoreMyRoutineById(controllerFunc apitransport.ControllerFunc[*routinesdto.RestoreMyRoutineByIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &routinesdto.RestoreMyRoutineByIdRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")
		value, ok := parseRoutineUUID(ctx, "routineId")
		if !ok {
			return
		}
		requestDto.Body.RoutineId = value
		controllerFunc(ctx, requestDto)
		return
	}
}

func (b *RoutineBinder) BindRestoreMyRoutinesByIds(controllerFunc apitransport.ControllerFunc[*routinesdto.RestoreMyRoutinesByIdsRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &routinesdto.RestoreMyRoutinesByIdsRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")
		bindRoutineJSON(ctx, requestDto, &requestDto.Body, controllerFunc)
		return
	}
}

func (b *RoutineBinder) BindDeleteMyRoutineById(controllerFunc apitransport.ControllerFunc[*routinesdto.DeleteMyRoutineByIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &routinesdto.DeleteMyRoutineByIdRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")
		value, ok := parseRoutineUUID(ctx, "routineId")
		if !ok {
			return
		}
		requestDto.Body.RoutineId = value
		controllerFunc(ctx, requestDto)
		return
	}
}

func (b *RoutineBinder) BindDeleteMyRoutinesByIds(controllerFunc apitransport.ControllerFunc[*routinesdto.DeleteMyRoutinesByIdsRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &routinesdto.DeleteMyRoutinesByIdsRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")
		bindRoutineJSON(ctx, requestDto, &requestDto.Body, controllerFunc)
		return
	}
}

func (b *RoutineBinder) BindHardDeleteMyRoutineById(controllerFunc apitransport.ControllerFunc[*routinesdto.HardDeleteMyRoutineByIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &routinesdto.HardDeleteMyRoutineByIdRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")
		value, ok := parseRoutineUUID(ctx, "routineId")
		if !ok {
			return
		}
		requestDto.Body.RoutineId = value
		controllerFunc(ctx, requestDto)
		return
	}
}

func (b *RoutineBinder) BindHardDeleteMyRoutinesByIds(controllerFunc apitransport.ControllerFunc[*routinesdto.HardDeleteMyRoutinesByIdsRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &routinesdto.HardDeleteMyRoutinesByIdsRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")
		bindRoutineJSON(ctx, requestDto, &requestDto.Body, controllerFunc)
		return
	}
}

func (b *RoutineBinder) BindVisualizeMyRoutineStatusCount(controllerFunc apitransport.ControllerFunc[*routinesdto.VisualizeMyRoutineStatusCountRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &routinesdto.VisualizeMyRoutineStatusCountRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")
		ok := true
		requestDto.Param.Permission, ok = parseRoutinePermission(ctx)
		if !ok {
			return
		}
		controllerFunc(ctx, requestDto)
		return
	}
}

func (b *RoutineBinder) BindVisualizeMyRoutinePeriodCount(controllerFunc apitransport.ControllerFunc[*routinesdto.VisualizeMyRoutinePeriodCountRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &routinesdto.VisualizeMyRoutinePeriodCountRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")
		ok := true
		requestDto.Param.Permission, ok = parseRoutinePermission(ctx)
		if !ok {
			return
		}
		controllerFunc(ctx, requestDto)
		return
	}
}

func (b *RoutineBinder) BindVisualizeMyRoutineScheduledStartAtCount(controllerFunc apitransport.ControllerFunc[*routinesdto.VisualizeMyRoutineScheduledStartAtCountRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &routinesdto.VisualizeMyRoutineScheduledStartAtCountRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")
		ok := true
		requestDto.Param.Permission, ok = parseRoutinePermission(ctx)
		if !ok {
			return
		}
		requestDto.Param.TimeHourUnit, _ = strconv.Atoi(ctx.Query("timeHourUnit"))
		requestDto.Param.QueryRangeStartedAt, ok = parseRoutineTime(ctx, "queryRangeStartedAt")
		if !ok {
			return
		}
		requestDto.Param.QueryRangeEndedAt, ok = parseRoutineTime(ctx, "queryRangeEndedAt")
		if !ok {
			return
		}
		controllerFunc(ctx, requestDto)
		return
	}
}

func (b *RoutineBinder) BindVisualizeMyRoutineScheduledEndAtCount(controllerFunc apitransport.ControllerFunc[*routinesdto.VisualizeMyRoutineScheduledEndAtCountRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &routinesdto.VisualizeMyRoutineScheduledEndAtCountRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")
		ok := true
		requestDto.Param.Permission, ok = parseRoutinePermission(ctx)
		if !ok {
			return
		}
		requestDto.Param.TimeHourUnit, _ = strconv.Atoi(ctx.Query("timeHourUnit"))
		requestDto.Param.QueryRangeStartedAt, ok = parseRoutineTime(ctx, "queryRangeStartedAt")
		if !ok {
			return
		}
		requestDto.Param.QueryRangeEndedAt, ok = parseRoutineTime(ctx, "queryRangeEndedAt")
		if !ok {
			return
		}
		controllerFunc(ctx, requestDto)
		return
	}
}
