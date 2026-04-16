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

type SubShelfBinderInterface interface {
	BindGetMySubShelfById(controllerFunc types.ControllerFunc[*dtos.GetMySubShelfByIdReqDto]) gin.HandlerFunc
	BindGetMySubShelvesByPrevSubShelfId(controllerFunc types.ControllerFunc[*dtos.GetMySubShelvesByPrevSubShelfIdReqDto]) gin.HandlerFunc
	BindGetAllMySubShelvesByRootShelfId(controllerFunc types.ControllerFunc[*dtos.GetAllMySubShelvesByRootShelfIdReqDto]) gin.HandlerFunc
	BindGetMySubShelvesAndItemsByPrevSubShelfId(controllerFunc types.ControllerFunc[*dtos.GetMySubShelvesAndItemsByPrevSubShelfIdReqDto]) gin.HandlerFunc
	BindCreateSubShelfByRootShelfId(controllerFunc types.ControllerFunc[*dtos.CreateSubShelfByRootShelfIdReqDto]) gin.HandlerFunc
	BindCreateSubShelvesByRootShelfIds(controllerFunc types.ControllerFunc[*dtos.CreateSubShelvesByRootShelfIdsReqDto]) gin.HandlerFunc
	BindUpdateMySubShelfById(controllerFunc types.ControllerFunc[*dtos.UpdateMySubShelfByIdReqDto]) gin.HandlerFunc
	BindUpdateMySubShelvesByIds(controllerFunc types.ControllerFunc[*dtos.UpdateMySubShelvesByIdsReqDto]) gin.HandlerFunc
	BindMoveMySubShelf(controllerFunc types.ControllerFunc[*dtos.MoveMySubShelfReqDto]) gin.HandlerFunc
	BindMoveMySubShelves(controllerFunc types.ControllerFunc[*dtos.MoveMySubShelvesReqDto]) gin.HandlerFunc
	BindBatchMoveMySubShelves(controllerFunc types.ControllerFunc[*dtos.BatchMoveMySubShelvesReqDto]) gin.HandlerFunc
	BindRestoreMySubShelfById(controllerFunc types.ControllerFunc[*dtos.RestoreMySubShelfByIdReqDto]) gin.HandlerFunc
	BindRestoreMySubShelvesByIds(controllerFunc types.ControllerFunc[*dtos.RestoreMySubShelvesByIdsReqDto]) gin.HandlerFunc
	BindDeleteMySubShelfById(controllerFunc types.ControllerFunc[*dtos.DeleteMySubShelfByIdReqDto]) gin.HandlerFunc
	BindDeleteMySubShelvesByIds(controllerFunc types.ControllerFunc[*dtos.DeleteMySubShelvesByIdsReqDto]) gin.HandlerFunc
}

type SubShelfBinder struct{}

func NewSubShelfBinder() SubShelfBinderInterface {
	return &SubShelfBinder{}
}

func (b *SubShelfBinder) BindGetMySubShelfById(controllerFunc types.ControllerFunc[*dtos.GetMySubShelfByIdReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.GetMySubShelfByIdReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		subShelfIdString := ctx.Query("subShelfId")
		if subShelfIdString == "" {
			exceptions.Shelf.InvalidInput().WithOrigin(fmt.Errorf("subShelfId is required")).Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		subShelfId, err := uuid.Parse(subShelfIdString)
		if err != nil {
			exceptions.Shelf.InvalidInput().WithOrigin(err).Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.Param.SubShelfId = subShelfId

		controllerFunc(ctx, &reqDto)
	}
}

func (b *SubShelfBinder) BindGetMySubShelvesByPrevSubShelfId(controllerFunc types.ControllerFunc[*dtos.GetMySubShelvesByPrevSubShelfIdReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.GetMySubShelvesByPrevSubShelfIdReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		prevSubShelfIdString := ctx.Query("prevSubShelfId")
		if prevSubShelfIdString == "" {
			exceptions.Shelf.InvalidInput().WithOrigin(fmt.Errorf("prevSubShelfId is required")).SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		prevSubShelfId, err := uuid.Parse(prevSubShelfIdString)
		if err != nil {
			exceptions.Shelf.InvalidInput().WithOrigin(err).SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.Param.PrevSubShelfId = prevSubShelfId

		controllerFunc(ctx, &reqDto)
	}
}

func (b *SubShelfBinder) BindGetAllMySubShelvesByRootShelfId(controllerFunc types.ControllerFunc[*dtos.GetAllMySubShelvesByRootShelfIdReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.GetAllMySubShelvesByRootShelfIdReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		rootShelfIdString := ctx.Query("rootShelfId")
		if rootShelfIdString == "" {
			exceptions.Shelf.InvalidInput().WithOrigin(fmt.Errorf("rootShelfId is required")).SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		rootShelfId, err := uuid.Parse(rootShelfIdString)
		if err != nil {
			exceptions.Shelf.InvalidInput().WithOrigin(err).SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.Param.RootShelfId = rootShelfId

		controllerFunc(ctx, &reqDto)
	}
}

func (b *SubShelfBinder) BindGetMySubShelvesAndItemsByPrevSubShelfId(controllerFunc types.ControllerFunc[*dtos.GetMySubShelvesAndItemsByPrevSubShelfIdReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.GetMySubShelvesAndItemsByPrevSubShelfIdReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		prevSubShelfIdString := ctx.Query("prevSubShelfId")
		if prevSubShelfIdString == "" {
			exceptions.Shelf.InvalidInput().WithOrigin(fmt.Errorf("prevSubShelfId is required")).SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		prevSubShelfId, err := uuid.Parse(prevSubShelfIdString)
		if err != nil {
			exceptions.Shelf.InvalidInput().WithOrigin(err).SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.Param.PrevSubShelfId = prevSubShelfId

		controllerFunc(ctx, &reqDto)
	}
}

func (b *SubShelfBinder) BindCreateSubShelfByRootShelfId(controllerFunc types.ControllerFunc[*dtos.CreateSubShelfByRootShelfIdReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.CreateSubShelfByRootShelfIdReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.Shelf.InvalidDto().WithOrigin(err)
			exception.SafelyAbortAndResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *SubShelfBinder) BindCreateSubShelvesByRootShelfIds(controllerFunc types.ControllerFunc[*dtos.CreateSubShelvesByRootShelfIdsReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.CreateSubShelvesByRootShelfIdsReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.Shelf.InvalidDto().WithOrigin(err)
			exception.SafelyAbortAndResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *SubShelfBinder) BindUpdateMySubShelfById(controllerFunc types.ControllerFunc[*dtos.UpdateMySubShelfByIdReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.UpdateMySubShelfByIdReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.Shelf.InvalidDto().WithOrigin(err)
			exception.SafelyAbortAndResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *SubShelfBinder) BindUpdateMySubShelvesByIds(controllerFunc types.ControllerFunc[*dtos.UpdateMySubShelvesByIdsReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.UpdateMySubShelvesByIdsReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.Shelf.InvalidDto().WithOrigin(err)
			exception.SafelyAbortAndResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *SubShelfBinder) BindMoveMySubShelf(controllerFunc types.ControllerFunc[*dtos.MoveMySubShelfReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.MoveMySubShelfReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.Shelf.InvalidDto().WithOrigin(err)
			exception.SafelyAbortAndResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *SubShelfBinder) BindMoveMySubShelves(controllerFunc types.ControllerFunc[*dtos.MoveMySubShelvesReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.MoveMySubShelvesReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.Shelf.InvalidDto().WithOrigin(err)
			exception.SafelyAbortAndResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *SubShelfBinder) BindBatchMoveMySubShelves(controllerFunc types.ControllerFunc[*dtos.BatchMoveMySubShelvesReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.BatchMoveMySubShelvesReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.Shelf.InvalidDto().WithOrigin(err)
			exception.SafelyAbortAndResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *SubShelfBinder) BindRestoreMySubShelfById(controllerFunc types.ControllerFunc[*dtos.RestoreMySubShelfByIdReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.RestoreMySubShelfByIdReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.Shelf.InvalidDto().WithOrigin(err)
			exception.SafelyAbortAndResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}
func (b *SubShelfBinder) BindRestoreMySubShelvesByIds(controllerFunc types.ControllerFunc[*dtos.RestoreMySubShelvesByIdsReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.RestoreMySubShelvesByIdsReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.Shelf.InvalidDto().WithOrigin(err)
			exception.SafelyAbortAndResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *SubShelfBinder) BindDeleteMySubShelfById(controllerFunc types.ControllerFunc[*dtos.DeleteMySubShelfByIdReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.DeleteMySubShelfByIdReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.Shelf.InvalidDto().WithOrigin(err)
			exception.SafelyAbortAndResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *SubShelfBinder) BindDeleteMySubShelvesByIds(controllerFunc types.ControllerFunc[*dtos.DeleteMySubShelvesByIdsReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.DeleteMySubShelvesByIdsReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.Shelf.InvalidDto().WithOrigin(err)
			exception.SafelyAbortAndResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}
