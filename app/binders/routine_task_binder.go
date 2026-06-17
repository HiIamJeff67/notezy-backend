package binders

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	contexts "github.com/HiIamJeff67/notezy-backend/app/contexts"
	dtos "github.com/HiIamJeff67/notezy-backend/app/dtos"
	exceptions "github.com/HiIamJeff67/notezy-backend/app/exceptions"
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
