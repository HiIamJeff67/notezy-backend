package binders

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	contexts "github.com/HiIamJeff67/notezy-backend/app/contexts"
	dtos "github.com/HiIamJeff67/notezy-backend/app/dtos"
	exceptions "github.com/HiIamJeff67/notezy-backend/app/exceptions"
	types "github.com/HiIamJeff67/notezy-backend/shared/types"
)

type RoutineBinderInterface interface {
	BindGetMyRoutineById(controllerFunc types.ControllerFunc[*dtos.GetMyRoutineByIdReqDto]) gin.HandlerFunc
	BindCreateRoutineByStationId(controllerFunc types.ControllerFunc[*dtos.CreateRoutineByStationIdReqDto]) gin.HandlerFunc
	BindCreateRoutinesByStationIds(controllerFunc types.ControllerFunc[*dtos.CreateRoutinesByStationIdsReqDto]) gin.HandlerFunc
	BindUpdateMyRoutineById(controllerFunc types.ControllerFunc[*dtos.UpdateMyRoutineByIdReqDto]) gin.HandlerFunc
	BindUpdateMyRoutinesByIds(controllerFunc types.ControllerFunc[*dtos.UpdateMyRoutinesByIdsReqDto]) gin.HandlerFunc
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

		if err := ctx.ShouldBindQuery(&reqDto.Param); err != nil {
			exceptions.Routine.InvalidInput().WithOrigin(err).SafelyAbortAndResponseWithJSON(ctx)
			return
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
