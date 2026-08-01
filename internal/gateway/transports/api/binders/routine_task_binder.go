package binders

import (
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	routinetasksdto "github.com/HiIamJeff67/notezy-backend/contracts/api/v1/routine-tasks"
	exceptions "github.com/HiIamJeff67/notezy-backend/internal/exceptions"
	responsewriter "github.com/HiIamJeff67/notezy-backend/internal/gateway/responsewriter"
	apitransport "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/api"
	sharedtypes "github.com/HiIamJeff67/notezy-backend/internal/shared/types"
)

type RoutineTaskBinderInterface interface {
	BindGetMyRoutineTaskById(controllerFunc apitransport.ControllerFunc[*routinetasksdto.GetMyRoutineTaskByIdRequestDto]) gin.HandlerFunc
	BindGetAllMyRoutineTasksByRoutineIds(controllerFunc apitransport.ControllerFunc[*routinetasksdto.GetAllMyRoutineTasksByRoutineIdsRequestDto]) gin.HandlerFunc
	BindGetAllMyRoutineTasks(controllerFunc apitransport.ControllerFunc[*routinetasksdto.GetAllMyRoutineTasksRequestDto]) gin.HandlerFunc
	BindCreateRoutineTaskByRoutineId(controllerFunc apitransport.ControllerFunc[*routinetasksdto.CreateRoutineTaskByRoutineIdRequestDto]) gin.HandlerFunc
	BindUpdateMyRoutineTaskById(controllerFunc apitransport.ControllerFunc[*routinetasksdto.UpdateMyRoutineTaskByIdRequestDto]) gin.HandlerFunc
	BindPauseMyRoutineTaskById(controllerFunc apitransport.ControllerFunc[*routinetasksdto.PauseMyRoutineTaskByIdRequestDto]) gin.HandlerFunc
	BindResumeMyRoutineTaskById(controllerFunc apitransport.ControllerFunc[*routinetasksdto.ResumeMyRoutineTaskByIdRequestDto]) gin.HandlerFunc
	BindHardDeleteMyRoutineTaskById(controllerFunc apitransport.ControllerFunc[*routinetasksdto.HardDeleteMyRoutineTaskByIdRequestDto]) gin.HandlerFunc
	BindHardDeleteMyRoutineTasksByIds(controllerFunc apitransport.ControllerFunc[*routinetasksdto.HardDeleteMyRoutineTasksByIdsRequestDto]) gin.HandlerFunc

	/* ============================== Visualization Methods ============================== */
	BindVisualizeMyRoutineTaskStatusCount(controllerFunc apitransport.ControllerFunc[*routinetasksdto.VisualizeMyRoutineTaskStatusCountRequestDto]) gin.HandlerFunc
	BindVisualizeMyRoutineTaskPurposeCount(controllerFunc apitransport.ControllerFunc[*routinetasksdto.VisualizeMyRoutineTaskPurposeCountRequestDto]) gin.HandlerFunc
	BindVisualizeMyRoutineTaskScheduledAtCount(controllerFunc apitransport.ControllerFunc[*routinetasksdto.VisualizeMyRoutineTaskScheduledAtCountRequestDto]) gin.HandlerFunc
	BindVisualizeMyRoutineTaskActualStartedAtCount(controllerFunc apitransport.ControllerFunc[*routinetasksdto.VisualizeMyRoutineTaskActualStartedAtCountRequestDto]) gin.HandlerFunc
	BindVisualizeMyRoutineTaskActualEndedAtCount(controllerFunc apitransport.ControllerFunc[*routinetasksdto.VisualizeMyRoutineTaskActualEndedAtCountRequestDto]) gin.HandlerFunc
}

type RoutineTaskBinder struct{}

func NewRoutineTaskBinder() RoutineTaskBinderInterface { return &RoutineTaskBinder{} }

func bindRoutineTaskJSON[T any](ctx *gin.Context, requestDto *T, body any, controllerFunc apitransport.ControllerFunc[*T]) {
	if err := ctx.ShouldBindJSON(body); err != nil {
		responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("RoutineTask").WithOrigin(err), ctx)
		return
	}
	controllerFunc(ctx, requestDto)
}

func parseRoutineTaskUUID(ctx *gin.Context, name string) (uuid.UUID, bool) {
	value, err := uuid.Parse(ctx.Param(name))
	if err != nil {
		responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("RoutineTask").WithOrigin(err), ctx)
		return uuid.Nil, false
	}
	return value, true
}
func parseRoutineTaskBool(ctx *gin.Context, name string) (*bool, bool) {
	valueString := ctx.Query(name)
	if valueString == "" {
		return nil, true
	}
	value, err := strconv.ParseBool(valueString)
	if err != nil {
		responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("RoutineTask").WithOrigin(err), ctx)
		return nil, false
	}
	return &value, true
}
func parseRoutineTaskPermission(ctx *gin.Context) (sharedtypes.AccessControlPermission, bool) {
	value, err := sharedtypes.ParseAccessControlPermission(ctx.Query("permission"))
	if err != nil {
		responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("RoutineTask").WithOrigin(err), ctx)
		return "", false
	}
	return *value, true
}
func parseRoutineTaskTime(ctx *gin.Context, name string) (time.Time, bool) {
	value, err := time.Parse(time.RFC3339, ctx.Query(name))
	if err != nil {
		responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("RoutineTask").WithOrigin(err), ctx)
		return time.Time{}, false
	}
	return value, true
}

func parseRoutineTaskUUIDs(ctx *gin.Context, name string) ([]uuid.UUID, bool) {
	values := ctx.QueryArray(name)
	if len(values) == 1 {
		values = strings.Split(values[0], ",")
	}
	ids := make([]uuid.UUID, len(values))
	for index, value := range values {
		parsed, err := uuid.Parse(value)
		if err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("RoutineTask").WithOrigin(err), ctx)
			return nil, false
		}
		ids[index] = parsed
	}
	return ids, true
}

func (b *RoutineTaskBinder) BindGetMyRoutineTaskById(controllerFunc apitransport.ControllerFunc[*routinetasksdto.GetMyRoutineTaskByIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &routinetasksdto.GetMyRoutineTaskByIdRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")
		isDeleted, ok := parseRoutineTaskBool(ctx, "isDeleted")
		if !ok {
			return
		}
		requestDto.Param.IsDeleted = isDeleted
		value, ok := parseRoutineTaskUUID(ctx, "routineTaskId")
		if !ok {
			return
		}
		requestDto.Param.RoutineTaskId = value
		controllerFunc(ctx, requestDto)
		return
	}
}

func (b *RoutineTaskBinder) BindGetAllMyRoutineTasksByRoutineIds(controllerFunc apitransport.ControllerFunc[*routinetasksdto.GetAllMyRoutineTasksByRoutineIdsRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &routinetasksdto.GetAllMyRoutineTasksByRoutineIdsRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")
		areDeleted, ok := parseRoutineTaskBool(ctx, "areDeleted")
		if !ok {
			return
		}
		requestDto.Param.AreDeleted = areDeleted
		if ctx.Query("areDeleted") != "" && requestDto.Param.AreDeleted == nil {
			return
		}
		requestDto.Param.RoutineIds, ok = parseRoutineTaskUUIDs(ctx, "routineIds")
		if !ok {
			return
		}
		controllerFunc(ctx, requestDto)
		return
	}
}

func (b *RoutineTaskBinder) BindGetAllMyRoutineTasks(controllerFunc apitransport.ControllerFunc[*routinetasksdto.GetAllMyRoutineTasksRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &routinetasksdto.GetAllMyRoutineTasksRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")
		areDeleted, ok := parseRoutineTaskBool(ctx, "areDeleted")
		if !ok {
			return
		}
		requestDto.Param.AreDeleted = areDeleted
		controllerFunc(ctx, requestDto)
		return
	}
}

func (b *RoutineTaskBinder) BindCreateRoutineTaskByRoutineId(controllerFunc apitransport.ControllerFunc[*routinetasksdto.CreateRoutineTaskByRoutineIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &routinetasksdto.CreateRoutineTaskByRoutineIdRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")
		value, ok := parseRoutineTaskUUID(ctx, "routineId")
		if !ok {
			return
		}
		requestDto.Body.RoutineId = value
		bindRoutineTaskJSON(ctx, requestDto, &requestDto.Body, controllerFunc)
		return
	}
}

func (b *RoutineTaskBinder) BindUpdateMyRoutineTaskById(controllerFunc apitransport.ControllerFunc[*routinetasksdto.UpdateMyRoutineTaskByIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &routinetasksdto.UpdateMyRoutineTaskByIdRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")
		value, ok := parseRoutineTaskUUID(ctx, "routineTaskId")
		if !ok {
			return
		}
		requestDto.Body.RoutineTaskId = value
		bindRoutineTaskJSON(ctx, requestDto, &requestDto.Body, controllerFunc)
		return
	}
}

func (b *RoutineTaskBinder) BindPauseMyRoutineTaskById(controllerFunc apitransport.ControllerFunc[*routinetasksdto.PauseMyRoutineTaskByIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &routinetasksdto.PauseMyRoutineTaskByIdRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")
		value, ok := parseRoutineTaskUUID(ctx, "routineTaskId")
		if !ok {
			return
		}
		requestDto.Body.RoutineTaskId = value
		controllerFunc(ctx, requestDto)
		return
	}
}

func (b *RoutineTaskBinder) BindResumeMyRoutineTaskById(controllerFunc apitransport.ControllerFunc[*routinetasksdto.ResumeMyRoutineTaskByIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &routinetasksdto.ResumeMyRoutineTaskByIdRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")
		value, ok := parseRoutineTaskUUID(ctx, "routineTaskId")
		if !ok {
			return
		}
		requestDto.Body.RoutineTaskId = value
		controllerFunc(ctx, requestDto)
		return
	}
}

func (b *RoutineTaskBinder) BindHardDeleteMyRoutineTaskById(controllerFunc apitransport.ControllerFunc[*routinetasksdto.HardDeleteMyRoutineTaskByIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &routinetasksdto.HardDeleteMyRoutineTaskByIdRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")
		value, ok := parseRoutineTaskUUID(ctx, "routineTaskId")
		if !ok {
			return
		}
		requestDto.Body.RoutineTaskId = value
		controllerFunc(ctx, requestDto)
		return
	}
}

func (b *RoutineTaskBinder) BindHardDeleteMyRoutineTasksByIds(controllerFunc apitransport.ControllerFunc[*routinetasksdto.HardDeleteMyRoutineTasksByIdsRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &routinetasksdto.HardDeleteMyRoutineTasksByIdsRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")
		bindRoutineTaskJSON(ctx, requestDto, &requestDto.Body, controllerFunc)
		return
	}
}

func (b *RoutineTaskBinder) BindVisualizeMyRoutineTaskStatusCount(controllerFunc apitransport.ControllerFunc[*routinetasksdto.VisualizeMyRoutineTaskStatusCountRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &routinetasksdto.VisualizeMyRoutineTaskStatusCountRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")
		ok := true
		requestDto.Param.Permission, ok = parseRoutineTaskPermission(ctx)
		if !ok {
			return
		}
		controllerFunc(ctx, requestDto)
		return
	}
}

func (b *RoutineTaskBinder) BindVisualizeMyRoutineTaskPurposeCount(controllerFunc apitransport.ControllerFunc[*routinetasksdto.VisualizeMyRoutineTaskPurposeCountRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &routinetasksdto.VisualizeMyRoutineTaskPurposeCountRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")
		ok := true
		requestDto.Param.Permission, ok = parseRoutineTaskPermission(ctx)
		if !ok {
			return
		}
		controllerFunc(ctx, requestDto)
		return
	}
}

func (b *RoutineTaskBinder) BindVisualizeMyRoutineTaskScheduledAtCount(controllerFunc apitransport.ControllerFunc[*routinetasksdto.VisualizeMyRoutineTaskScheduledAtCountRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &routinetasksdto.VisualizeMyRoutineTaskScheduledAtCountRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")
		ok := true
		requestDto.Param.Permission, ok = parseRoutineTaskPermission(ctx)
		if !ok {
			return
		}
		requestDto.Param.TimeHourUnit, _ = strconv.Atoi(ctx.Query("timeHourUnit"))
		requestDto.Param.QueryRangeStartedAt, ok = parseRoutineTaskTime(ctx, "queryRangeStartedAt")
		if !ok {
			return
		}
		requestDto.Param.QueryRangeEndedAt, ok = parseRoutineTaskTime(ctx, "queryRangeEndedAt")
		if !ok {
			return
		}
		controllerFunc(ctx, requestDto)
		return
	}
}

func (b *RoutineTaskBinder) BindVisualizeMyRoutineTaskActualStartedAtCount(controllerFunc apitransport.ControllerFunc[*routinetasksdto.VisualizeMyRoutineTaskActualStartedAtCountRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &routinetasksdto.VisualizeMyRoutineTaskActualStartedAtCountRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")
		ok := true
		requestDto.Param.Permission, ok = parseRoutineTaskPermission(ctx)
		if !ok {
			return
		}
		requestDto.Param.TimeHourUnit, _ = strconv.Atoi(ctx.Query("timeHourUnit"))
		requestDto.Param.QueryRangeStartedAt, ok = parseRoutineTaskTime(ctx, "queryRangeStartedAt")
		if !ok {
			return
		}
		requestDto.Param.QueryRangeEndedAt, ok = parseRoutineTaskTime(ctx, "queryRangeEndedAt")
		if !ok {
			return
		}
		controllerFunc(ctx, requestDto)
		return
	}
}

func (b *RoutineTaskBinder) BindVisualizeMyRoutineTaskActualEndedAtCount(controllerFunc apitransport.ControllerFunc[*routinetasksdto.VisualizeMyRoutineTaskActualEndedAtCountRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &routinetasksdto.VisualizeMyRoutineTaskActualEndedAtCountRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")
		ok := true
		requestDto.Param.Permission, ok = parseRoutineTaskPermission(ctx)
		if !ok {
			return
		}
		requestDto.Param.TimeHourUnit, _ = strconv.Atoi(ctx.Query("timeHourUnit"))
		requestDto.Param.QueryRangeStartedAt, ok = parseRoutineTaskTime(ctx, "queryRangeStartedAt")
		if !ok {
			return
		}
		requestDto.Param.QueryRangeEndedAt, ok = parseRoutineTaskTime(ctx, "queryRangeEndedAt")
		if !ok {
			return
		}
		controllerFunc(ctx, requestDto)
		return
	}
}
