package binders

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	contexts "notezy-backend/app/contexts"
	dtos "notezy-backend/app/dtos"
	exceptions "notezy-backend/app/exceptions"
	types "notezy-backend/shared/types"
)

type RoutineTagBinderInterface interface {
	BindGetMyRoutineTagById(controllerFunc types.ControllerFunc[*dtos.GetMyRoutineTagByIdReqDto]) gin.HandlerFunc
	BindCreateRoutineTag(controllerFunc types.ControllerFunc[*dtos.CreateRoutineTagReqDto]) gin.HandlerFunc
	BindCreateRoutineTags(controllerFunc types.ControllerFunc[*dtos.CreateRoutineTagsReqDto]) gin.HandlerFunc
	BindUpdateMyRoutineTagById(controllerFunc types.ControllerFunc[*dtos.UpdateMyRoutineTagByIdReqDto]) gin.HandlerFunc
	BindUpdateMyRoutineTagsByIds(controllerFunc types.ControllerFunc[*dtos.UpdateMyRoutineTagsByIdsReqDto]) gin.HandlerFunc
	BindHardDeleteMyRoutineTagById(controllerFunc types.ControllerFunc[*dtos.HardDeleteMyRoutineTagByIdReqDto]) gin.HandlerFunc
	BindHardDeleteMyRoutineTagsByIds(controllerFunc types.ControllerFunc[*dtos.HardDeleteMyRoutineTagsByIdsReqDto]) gin.HandlerFunc
}

type RoutineTagBinder struct{}

func NewRoutineTagBinder() RoutineTagBinderInterface {
	return &RoutineTagBinder{}
}

func (b *RoutineTagBinder) BindGetMyRoutineTagById(controllerFunc types.ControllerFunc[*dtos.GetMyRoutineTagByIdReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.GetMyRoutineTagByIdReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		routineTagIdString := ctx.Query("routineTagId")
		if routineTagIdString == "" {
			exceptions.RoutineTag.InvalidInput().WithOrigin(fmt.Errorf("routineTagId is required")).SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		routineTagId, err := uuid.Parse(routineTagIdString)
		if err != nil {
			exceptions.RoutineTag.InvalidInput().WithOrigin(err).SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.Param.RoutineTagId = routineTagId

		controllerFunc(ctx, &reqDto)
	}
}

func (b *RoutineTagBinder) BindCreateRoutineTag(controllerFunc types.ControllerFunc[*dtos.CreateRoutineTagReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.CreateRoutineTagReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.RoutineTag.InvalidDto().WithOrigin(err)
			exception.SafelyAbortAndResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *RoutineTagBinder) BindCreateRoutineTags(controllerFunc types.ControllerFunc[*dtos.CreateRoutineTagsReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.CreateRoutineTagsReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.RoutineTag.InvalidDto().WithOrigin(err)
			exception.SafelyAbortAndResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *RoutineTagBinder) BindUpdateMyRoutineTagById(controllerFunc types.ControllerFunc[*dtos.UpdateMyRoutineTagByIdReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.UpdateMyRoutineTagByIdReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.RoutineTag.InvalidDto().WithOrigin(err)
			exception.SafelyAbortAndResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *RoutineTagBinder) BindUpdateMyRoutineTagsByIds(controllerFunc types.ControllerFunc[*dtos.UpdateMyRoutineTagsByIdsReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.UpdateMyRoutineTagsByIdsReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.RoutineTag.InvalidDto().WithOrigin(err)
			exception.SafelyAbortAndResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *RoutineTagBinder) BindHardDeleteMyRoutineTagById(controllerFunc types.ControllerFunc[*dtos.HardDeleteMyRoutineTagByIdReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.HardDeleteMyRoutineTagByIdReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.RoutineTag.InvalidDto().WithOrigin(err)
			exception.SafelyAbortAndResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *RoutineTagBinder) BindHardDeleteMyRoutineTagsByIds(controllerFunc types.ControllerFunc[*dtos.HardDeleteMyRoutineTagsByIdsReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.HardDeleteMyRoutineTagsByIdsReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.RoutineTag.InvalidDto().WithOrigin(err)
			exception.SafelyAbortAndResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}
