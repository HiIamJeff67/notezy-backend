package binders

import (
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	exceptions "github.com/HiIamJeff67/notezy-backend/shared/exceptions"

	exceptionwriter "github.com/HiIamJeff67/notezy-backend/shared/util/exceptionwriter"

	routinesdto "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api/routines"
	enumcontract "github.com/HiIamJeff67/notezy-backend/contracts/types/enums"

	controllers "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/api/controllers"
)

type RoutineBinderInterface interface {
	BindGetMyRoutineById(controllerFunc controllers.Func[*routinesdto.GetMyRoutineByIdRequestDto]) gin.HandlerFunc
	BindGetMyRoutinesByStationId(controllerFunc controllers.Func[*routinesdto.GetMyRoutinesByStationIdRequestDto]) gin.HandlerFunc
	BindGetAllMyRoutinesByTimeRange(controllerFunc controllers.Func[*routinesdto.GetAllMyRoutinesByTimeRangeRequestDto]) gin.HandlerFunc
	BindCreateRoutineByStationId(controllerFunc controllers.Func[*routinesdto.CreateRoutineByStationIdRequestDto]) gin.HandlerFunc
	BindCreateRoutinesByStationIds(controllerFunc controllers.Func[*routinesdto.CreateRoutinesByStationIdsRequestDto]) gin.HandlerFunc
	BindUpdateMyRoutineById(controllerFunc controllers.Func[*routinesdto.UpdateMyRoutineByIdRequestDto]) gin.HandlerFunc
	BindUpdateMyRoutinesByIds(controllerFunc controllers.Func[*routinesdto.UpdateMyRoutinesByIdsRequestDto]) gin.HandlerFunc
	BindLinkRoutineTagById(controllerFunc controllers.Func[*routinesdto.LinkRoutineTagByIdRequestDto]) gin.HandlerFunc
	BindLinkRoutineTagsByIds(controllerFunc controllers.Func[*routinesdto.LinkRoutineTagsByIdsRequestDto]) gin.HandlerFunc
	BindLinkRoutineItemById(controllerFunc controllers.Func[*routinesdto.LinkRoutineItemByIdRequestDto]) gin.HandlerFunc
	BindLinkRoutineItemsByIds(controllerFunc controllers.Func[*routinesdto.LinkRoutineItemsByIdsRequestDto]) gin.HandlerFunc
	BindRestoreMyRoutineById(controllerFunc controllers.Func[*routinesdto.RestoreMyRoutineByIdRequestDto]) gin.HandlerFunc
	BindRestoreMyRoutinesByIds(controllerFunc controllers.Func[*routinesdto.RestoreMyRoutinesByIdsRequestDto]) gin.HandlerFunc
	BindDeleteMyRoutineById(controllerFunc controllers.Func[*routinesdto.DeleteMyRoutineByIdRequestDto]) gin.HandlerFunc
	BindDeleteMyRoutinesByIds(controllerFunc controllers.Func[*routinesdto.DeleteMyRoutinesByIdsRequestDto]) gin.HandlerFunc
	BindHardDeleteMyRoutineById(controllerFunc controllers.Func[*routinesdto.HardDeleteMyRoutineByIdRequestDto]) gin.HandlerFunc
	BindHardDeleteMyRoutinesByIds(controllerFunc controllers.Func[*routinesdto.HardDeleteMyRoutinesByIdsRequestDto]) gin.HandlerFunc

	/* ============================== Visualization Methods ============================== */
	BindVisualizeMyRoutineStatusCount(controllerFunc controllers.Func[*routinesdto.VisualizeMyRoutineStatusCountRequestDto]) gin.HandlerFunc
	BindVisualizeMyRoutinePeriodCount(controllerFunc controllers.Func[*routinesdto.VisualizeMyRoutinePeriodCountRequestDto]) gin.HandlerFunc
	BindVisualizeMyRoutineScheduledStartAtCount(controllerFunc controllers.Func[*routinesdto.VisualizeMyRoutineScheduledStartAtCountRequestDto]) gin.HandlerFunc
	BindVisualizeMyRoutineScheduledEndAtCount(controllerFunc controllers.Func[*routinesdto.VisualizeMyRoutineScheduledEndAtCountRequestDto]) gin.HandlerFunc
}

type RoutineBinder struct{}

func NewRoutineBinder() RoutineBinderInterface { return &RoutineBinder{} }

func bindRoutineJSON[T any](ctx *gin.Context, requestDto *T, body any, controllerFunc controllers.Func[*T]) {
	if err := ctx.ShouldBindJSON(body); err != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("Routine").WithOrigin(err), ctx)
		return
	}
	controllerFunc(ctx, requestDto)
}

func parseRoutineUUID(ctx *gin.Context, name string) (uuid.UUID, bool) {
	value, err := uuid.Parse(ctx.Param(name))
	if err != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Routine").WithOrigin(err), ctx)
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
		exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Routine").WithOrigin(err), ctx)
		return nil, false
	}
	return &value, true
}

func parseRoutinePermission(ctx *gin.Context) (enumcontract.AccessControlPermission, bool) {
	permission := enumcontract.AccessControlPermission(ctx.Query("permission"))
	switch permission {
	case enumcontract.AccessControlPermission_Read,
		enumcontract.AccessControlPermission_Write,
		enumcontract.AccessControlPermission_Admin,
		enumcontract.AccessControlPermission_Owner:
		return permission, true
	default:
		exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Routine"), ctx)
		return "", false
	}
}

func parseRoutineTime(ctx *gin.Context, name string) (time.Time, bool) {
	value, err := time.Parse(time.RFC3339, ctx.Query(name))
	if err != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Routine").WithOrigin(err), ctx)
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
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Routine").WithOrigin(err), ctx)
			return nil, false
		}
		ids[index] = parsed
	}
	return ids, true
}

func (b *RoutineBinder) BindGetMyRoutineById(controllerFunc controllers.Func[*routinesdto.GetMyRoutineByIdRequestDto]) gin.HandlerFunc {
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

func (b *RoutineBinder) BindGetMyRoutinesByStationId(controllerFunc controllers.Func[*routinesdto.GetMyRoutinesByStationIdRequestDto]) gin.HandlerFunc {
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

func (b *RoutineBinder) BindGetAllMyRoutinesByTimeRange(controllerFunc controllers.Func[*routinesdto.GetAllMyRoutinesByTimeRangeRequestDto]) gin.HandlerFunc {
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

func (b *RoutineBinder) BindCreateRoutineByStationId(controllerFunc controllers.Func[*routinesdto.CreateRoutineByStationIdRequestDto]) gin.HandlerFunc {
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

func (b *RoutineBinder) BindCreateRoutinesByStationIds(controllerFunc controllers.Func[*routinesdto.CreateRoutinesByStationIdsRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &routinesdto.CreateRoutinesByStationIdsRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")
		bindRoutineJSON(ctx, requestDto, &requestDto.Body, controllerFunc)
		return
	}
}

func (b *RoutineBinder) BindUpdateMyRoutineById(controllerFunc controllers.Func[*routinesdto.UpdateMyRoutineByIdRequestDto]) gin.HandlerFunc {
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

func (b *RoutineBinder) BindUpdateMyRoutinesByIds(controllerFunc controllers.Func[*routinesdto.UpdateMyRoutinesByIdsRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &routinesdto.UpdateMyRoutinesByIdsRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")
		bindRoutineJSON(ctx, requestDto, &requestDto.Body, controllerFunc)
		return
	}
}

func (b *RoutineBinder) BindLinkRoutineTagById(controllerFunc controllers.Func[*routinesdto.LinkRoutineTagByIdRequestDto]) gin.HandlerFunc {
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

func (b *RoutineBinder) BindLinkRoutineTagsByIds(controllerFunc controllers.Func[*routinesdto.LinkRoutineTagsByIdsRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &routinesdto.LinkRoutineTagsByIdsRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")
		bindRoutineJSON(ctx, requestDto, &requestDto.Body, controllerFunc)
		return
	}
}

func (b *RoutineBinder) BindLinkRoutineItemById(controllerFunc controllers.Func[*routinesdto.LinkRoutineItemByIdRequestDto]) gin.HandlerFunc {
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
		requestDto.Body.ItemType = enumcontract.ItemType(ctx.Query("itemType"))
		requestDto.Body.IsUnlink = ctx.Query("isUnlink") == "true"
		controllerFunc(ctx, requestDto)
		return
	}
}

func (b *RoutineBinder) BindLinkRoutineItemsByIds(controllerFunc controllers.Func[*routinesdto.LinkRoutineItemsByIdsRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &routinesdto.LinkRoutineItemsByIdsRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")
		bindRoutineJSON(ctx, requestDto, &requestDto.Body, controllerFunc)
		return
	}
}

func (b *RoutineBinder) BindRestoreMyRoutineById(controllerFunc controllers.Func[*routinesdto.RestoreMyRoutineByIdRequestDto]) gin.HandlerFunc {
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

func (b *RoutineBinder) BindRestoreMyRoutinesByIds(controllerFunc controllers.Func[*routinesdto.RestoreMyRoutinesByIdsRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &routinesdto.RestoreMyRoutinesByIdsRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")
		bindRoutineJSON(ctx, requestDto, &requestDto.Body, controllerFunc)
		return
	}
}

func (b *RoutineBinder) BindDeleteMyRoutineById(controllerFunc controllers.Func[*routinesdto.DeleteMyRoutineByIdRequestDto]) gin.HandlerFunc {
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

func (b *RoutineBinder) BindDeleteMyRoutinesByIds(controllerFunc controllers.Func[*routinesdto.DeleteMyRoutinesByIdsRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &routinesdto.DeleteMyRoutinesByIdsRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")
		bindRoutineJSON(ctx, requestDto, &requestDto.Body, controllerFunc)
		return
	}
}

func (b *RoutineBinder) BindHardDeleteMyRoutineById(controllerFunc controllers.Func[*routinesdto.HardDeleteMyRoutineByIdRequestDto]) gin.HandlerFunc {
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

func (b *RoutineBinder) BindHardDeleteMyRoutinesByIds(controllerFunc controllers.Func[*routinesdto.HardDeleteMyRoutinesByIdsRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &routinesdto.HardDeleteMyRoutinesByIdsRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")
		bindRoutineJSON(ctx, requestDto, &requestDto.Body, controllerFunc)
		return
	}
}

func (b *RoutineBinder) BindVisualizeMyRoutineStatusCount(controllerFunc controllers.Func[*routinesdto.VisualizeMyRoutineStatusCountRequestDto]) gin.HandlerFunc {
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

func (b *RoutineBinder) BindVisualizeMyRoutinePeriodCount(controllerFunc controllers.Func[*routinesdto.VisualizeMyRoutinePeriodCountRequestDto]) gin.HandlerFunc {
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

func (b *RoutineBinder) BindVisualizeMyRoutineScheduledStartAtCount(controllerFunc controllers.Func[*routinesdto.VisualizeMyRoutineScheduledStartAtCountRequestDto]) gin.HandlerFunc {
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

func (b *RoutineBinder) BindVisualizeMyRoutineScheduledEndAtCount(controllerFunc controllers.Func[*routinesdto.VisualizeMyRoutineScheduledEndAtCountRequestDto]) gin.HandlerFunc {
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
