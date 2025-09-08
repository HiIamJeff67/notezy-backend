package binders

import (
	"github.com/gin-gonic/gin"

	contexts "notezy-backend/app/contexts"
	dtos "notezy-backend/app/dtos"
	exceptions "notezy-backend/app/exceptions"
	constants "notezy-backend/shared/constants"
	types "notezy-backend/shared/types"
)

/* ============================== Interface & Instance ============================== */

type ShelfBinderInterface interface {
	BindGetMyShelfById(controllerFunc types.ControllerFunc[*dtos.GetMyShelfByIdReqDto]) gin.HandlerFunc
	BindSearchRecentShelves(controllerFunc types.ControllerFunc[*dtos.SearchRecentShelvesReqDto]) gin.HandlerFunc
	BindCreateShelf(controllerFunc types.ControllerFunc[*dtos.CreateShelfReqDto]) gin.HandlerFunc
	BindSynchronizeShelves(controllerFunc types.ControllerFunc[*dtos.SynchronizeShelvesReqDto]) gin.HandlerFunc
	BindRestoreMyShelfById(controllerFunc types.ControllerFunc[*dtos.RestoreMyShelfByIdReqDto]) gin.HandlerFunc
	BindRestoreMyShelvesByIds(controllerFunc types.ControllerFunc[*dtos.RestoreMyShelvesByIdsReqDto]) gin.HandlerFunc
	BindDeleteMyShelfById(controllerFunc types.ControllerFunc[*dtos.DeleteMyShelfByIdReqDto]) gin.HandlerFunc
	BindDeleteMyShelvesByIds(controllerFunc types.ControllerFunc[*dtos.DeleteMyShelvesByIdsReqDto]) gin.HandlerFunc
}

type ShelfBinder struct{}

func NewShelfBinder() ShelfBinderInterface {
	return &ShelfBinder{}
}

/* ============================== Binder ============================== */

func (b *ShelfBinder) BindGetMyShelfById(controllerFunc types.ControllerFunc[*dtos.GetMyShelfByIdReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.GetMyShelfByIdReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.OwnerId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.Shelf.InvalidDto().WithError(err)
			exception.ResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *ShelfBinder) BindSearchRecentShelves(controllerFunc types.ControllerFunc[*dtos.SearchRecentShelvesReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.SearchRecentShelvesReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.OwnerId = *userId

		if err := ctx.ShouldBindQuery(&reqDto.Param); err != nil {
			exception.Log()
			exceptions.User.InvalidInput().WithError(err).ResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *ShelfBinder) BindCreateShelf(controllerFunc types.ControllerFunc[*dtos.CreateShelfReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.CreateShelfReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.OwnerId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.Shelf.InvalidDto().WithError(err)
			exception.ResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *ShelfBinder) BindSynchronizeShelves(controllerFunc types.ControllerFunc[*dtos.SynchronizeShelvesReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.SynchronizeShelvesReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.OwnerId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.Shelf.InvalidDto().WithError(err)
			exception.ResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *ShelfBinder) BindRestoreMyShelfById(controllerFunc types.ControllerFunc[*dtos.RestoreMyShelfByIdReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.RestoreMyShelfByIdReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.OwnerId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.Shelf.InvalidDto().WithError(err)
			exception.ResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *ShelfBinder) BindRestoreMyShelvesByIds(controllerFunc types.ControllerFunc[*dtos.RestoreMyShelvesByIdsReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.RestoreMyShelvesByIdsReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.OwnerId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.Shelf.InvalidDto().WithError(err)
			exception.ResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *ShelfBinder) BindDeleteMyShelfById(controllerFunc types.ControllerFunc[*dtos.DeleteMyShelfByIdReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.DeleteMyShelfByIdReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.OwnerId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.Shelf.InvalidDto().WithError(err)
			exception.ResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *ShelfBinder) BindDeleteMyShelvesByIds(controllerFunc types.ControllerFunc[*dtos.DeleteMyShelvesByIdsReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.DeleteMyShelvesByIdsReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.OwnerId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.Shelf.InvalidDto().WithError(err)
			exception.ResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}
