package binders

import (
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	contexts "github.com/HiIamJeff67/notezy-backend/app/contexts"
	dtos "github.com/HiIamJeff67/notezy-backend/app/dtos"
	exceptions "github.com/HiIamJeff67/notezy-backend/app/exceptions"
	types "github.com/HiIamJeff67/notezy-backend/shared/types"
)

type RoutineBinderInterface interface {
	BindGetMyRoutineById(controllerFunc types.ControllerFunc[*dtos.GetMyRoutineByIdReqDto]) gin.HandlerFunc
	BindGetAllMyRoutinesByTimeRange(controllerFunc types.ControllerFunc[*dtos.GetAllMyRoutinesByTimeRangeReqDto]) gin.HandlerFunc
	BindCreateRoutineByStationId(controllerFunc types.ControllerFunc[*dtos.CreateRoutineByStationIdReqDto]) gin.HandlerFunc
	BindCreateRoutinesByStationIds(controllerFunc types.ControllerFunc[*dtos.CreateRoutinesByStationIdsReqDto]) gin.HandlerFunc
	BindUpdateMyRoutineById(controllerFunc types.ControllerFunc[*dtos.UpdateMyRoutineByIdReqDto]) gin.HandlerFunc
	BindUpdateMyRoutinesByIds(controllerFunc types.ControllerFunc[*dtos.UpdateMyRoutinesByIdsReqDto]) gin.HandlerFunc
	BindLinkRoutineTagById(controllerFunc types.ControllerFunc[*dtos.LinkRoutineTagByIdReqDto]) gin.HandlerFunc
	BindBulkLinkRoutineTagsByIds(controllerFunc types.ControllerFunc[*dtos.BulkLinkRoutineTagsByIdsReqDto]) gin.HandlerFunc
	BindLinkRoutineTaskById(controllerFunc types.ControllerFunc[*dtos.LinkRoutineTaskByIdReqDto]) gin.HandlerFunc
	BindBulkLinkRoutineTasksByIds(controllerFunc types.ControllerFunc[*dtos.BulkLinkRoutineTasksByIdsReqDto]) gin.HandlerFunc
	BindLinkRoutineItemById(controllerFunc types.ControllerFunc[*dtos.LinkRoutineItemByIdReqDto]) gin.HandlerFunc
	BindBulkLinkRoutineItemsByIds(controllerFunc types.ControllerFunc[*dtos.BulkLinkRoutineItemsByIdsReqDto]) gin.HandlerFunc
	BindRestoreMyRoutineById(controllerFunc types.ControllerFunc[*dtos.RestoreMyRoutineByIdReqDto]) gin.HandlerFunc
	BindRestoreMyRoutinesByIds(controllerFunc types.ControllerFunc[*dtos.RestoreMyRoutinesByIdsReqDto]) gin.HandlerFunc
	BindDeleteMyRoutineById(controllerFunc types.ControllerFunc[*dtos.DeleteMyRoutineByIdReqDto]) gin.HandlerFunc
	BindDeleteMyRoutinesByIds(controllerFunc types.ControllerFunc[*dtos.DeleteMyRoutinesByIdsReqDto]) gin.HandlerFunc
	BindHardDeleteMyRoutineById(controllerFunc types.ControllerFunc[*dtos.HardDeleteMyRoutineByIdReqDto]) gin.HandlerFunc
	BindHardDeleteMyRoutinesByIds(controllerFunc types.ControllerFunc[*dtos.HardDeleteMyRoutinesByIdsReqDto]) gin.HandlerFunc
}

type RoutineBinder struct{}

func NewRoutineBinder() RoutineBinderInterface {
	return &RoutineBinder{}
}

func (b *RoutineBinder) BindGetMyRoutineById(controllerFunc types.ControllerFunc[*dtos.GetMyRoutineByIdReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.GetMyRoutineByIdReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		routineIdString := ctx.Query("routineId")
		if routineIdString == "" {
			exceptions.Routine.InvalidInput().WithOrigin(fmt.Errorf("routineId is required")).SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		routineId, err := uuid.Parse(routineIdString)
		if err != nil {
			exceptions.Routine.InvalidInput().WithOrigin(err).SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.Param.RoutineId = routineId

		controllerFunc(ctx, &reqDto)
	}
}

func (b *RoutineBinder) BindGetAllMyRoutinesByTimeRange(
	controllerFunc types.ControllerFunc[*dtos.GetAllMyRoutinesByTimeRangeReqDto],
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.GetAllMyRoutinesByTimeRangeReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		from, err := time.Parse(time.RFC3339, ctx.Query("from"))
		if err != nil {
			exceptions.Routine.InvalidInput().WithOrigin(fmt.Errorf("from must be an RFC3339 timestamp: %w", err)).SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		to, err := time.Parse(time.RFC3339, ctx.Query("to"))
		if err != nil {
			exceptions.Routine.InvalidInput().WithOrigin(fmt.Errorf("to must be an RFC3339 timestamp: %w", err)).SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.Param.From = from
		reqDto.Param.To = to

		for _, stationIdsValue := range ctx.QueryArray("stationIds") {
			for _, stationIdValue := range strings.Split(stationIdsValue, ",") {
				stationId, err := uuid.Parse(strings.TrimSpace(stationIdValue))
				if err != nil {
					exceptions.Routine.InvalidInput().WithOrigin(err).SafelyAbortAndResponseWithJSON(ctx)
					return
				}
				reqDto.Param.StationIds = append(reqDto.Param.StationIds, stationId)
			}
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *RoutineBinder) BindCreateRoutineByStationId(controllerFunc types.ControllerFunc[*dtos.CreateRoutineByStationIdReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.CreateRoutineByStationIdReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.Routine.InvalidDto().WithOrigin(err)
			exception.SafelyAbortAndResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *RoutineBinder) BindCreateRoutinesByStationIds(controllerFunc types.ControllerFunc[*dtos.CreateRoutinesByStationIdsReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.CreateRoutinesByStationIdsReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.Routine.InvalidDto().WithOrigin(err)
			exception.SafelyAbortAndResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *RoutineBinder) BindUpdateMyRoutineById(controllerFunc types.ControllerFunc[*dtos.UpdateMyRoutineByIdReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.UpdateMyRoutineByIdReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.Routine.InvalidDto().WithOrigin(err)
			exception.SafelyAbortAndResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *RoutineBinder) BindUpdateMyRoutinesByIds(controllerFunc types.ControllerFunc[*dtos.UpdateMyRoutinesByIdsReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.UpdateMyRoutinesByIdsReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.Routine.InvalidDto().WithOrigin(err)
			exception.SafelyAbortAndResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *RoutineBinder) BindLinkRoutineTagById(controllerFunc types.ControllerFunc[*dtos.LinkRoutineTagByIdReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.LinkRoutineTagByIdReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.Routine.InvalidDto().WithOrigin(err)
			exception.SafelyAbortAndResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *RoutineBinder) BindBulkLinkRoutineTagsByIds(controllerFunc types.ControllerFunc[*dtos.BulkLinkRoutineTagsByIdsReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.BulkLinkRoutineTagsByIdsReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.Routine.InvalidDto().WithOrigin(err)
			exception.SafelyAbortAndResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *RoutineBinder) BindLinkRoutineTaskById(controllerFunc types.ControllerFunc[*dtos.LinkRoutineTaskByIdReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.LinkRoutineTaskByIdReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.Routine.InvalidDto().WithOrigin(err)
			exception.SafelyAbortAndResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *RoutineBinder) BindBulkLinkRoutineTasksByIds(controllerFunc types.ControllerFunc[*dtos.BulkLinkRoutineTasksByIdsReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.BulkLinkRoutineTasksByIdsReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.Routine.InvalidDto().WithOrigin(err)
			exception.SafelyAbortAndResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *RoutineBinder) BindLinkRoutineItemById(controllerFunc types.ControllerFunc[*dtos.LinkRoutineItemByIdReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.LinkRoutineItemByIdReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.Routine.InvalidDto().WithOrigin(err)
			exception.SafelyAbortAndResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *RoutineBinder) BindBulkLinkRoutineItemsByIds(controllerFunc types.ControllerFunc[*dtos.BulkLinkRoutineItemsByIdsReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.BulkLinkRoutineItemsByIdsReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.Routine.InvalidDto().WithOrigin(err)
			exception.SafelyAbortAndResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *RoutineBinder) BindRestoreMyRoutineById(controllerFunc types.ControllerFunc[*dtos.RestoreMyRoutineByIdReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.RestoreMyRoutineByIdReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.Routine.InvalidDto().WithOrigin(err)
			exception.SafelyAbortAndResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *RoutineBinder) BindRestoreMyRoutinesByIds(controllerFunc types.ControllerFunc[*dtos.RestoreMyRoutinesByIdsReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.RestoreMyRoutinesByIdsReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.Routine.InvalidDto().WithOrigin(err)
			exception.SafelyAbortAndResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *RoutineBinder) BindDeleteMyRoutineById(controllerFunc types.ControllerFunc[*dtos.DeleteMyRoutineByIdReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.DeleteMyRoutineByIdReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.Routine.InvalidDto().WithOrigin(err)
			exception.SafelyAbortAndResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *RoutineBinder) BindDeleteMyRoutinesByIds(controllerFunc types.ControllerFunc[*dtos.DeleteMyRoutinesByIdsReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.DeleteMyRoutinesByIdsReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.Routine.InvalidDto().WithOrigin(err)
			exception.SafelyAbortAndResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *RoutineBinder) BindHardDeleteMyRoutineById(controllerFunc types.ControllerFunc[*dtos.HardDeleteMyRoutineByIdReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.HardDeleteMyRoutineByIdReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.Routine.InvalidDto().WithOrigin(err)
			exception.SafelyAbortAndResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *RoutineBinder) BindHardDeleteMyRoutinesByIds(controllerFunc types.ControllerFunc[*dtos.HardDeleteMyRoutinesByIdsReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.HardDeleteMyRoutinesByIdsReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.Routine.InvalidDto().WithOrigin(err)
			exception.SafelyAbortAndResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}
