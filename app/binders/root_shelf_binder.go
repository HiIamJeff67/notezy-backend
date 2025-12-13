package binders

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	contexts "notezy-backend/app/contexts"
	dtos "notezy-backend/app/dtos"
	exceptions "notezy-backend/app/exceptions"
	constants "notezy-backend/shared/constants"
	types "notezy-backend/shared/types"
)

/* ============================== Interface & Instance ============================== */

type RootShelfBinderInterface interface {
	BindGetMyRootShelfById(controllerFunc types.ControllerFunc[*dtos.GetMyRootShelfByIdReqDto]) gin.HandlerFunc
	BindSearchRecentRootShelves(controllerFunc types.ControllerFunc[*dtos.SearchRecentRootShelvesReqDto]) gin.HandlerFunc
	BindCreateRootShelf(controllerFunc types.ControllerFunc[*dtos.CreateRootShelfReqDto]) gin.HandlerFunc
	BindUpdateMyRootShelfById(controllerFunc types.ControllerFunc[*dtos.UpdateMyRootShelfByIdReqDto]) gin.HandlerFunc
	BindRestoreMyRootShelfById(controllerFunc types.ControllerFunc[*dtos.RestoreMyRootShelfByIdReqDto]) gin.HandlerFunc
	BindRestoreMyRootShelvesByIds(controllerFunc types.ControllerFunc[*dtos.RestoreMyRootShelvesByIdsReqDto]) gin.HandlerFunc
	BindDeleteMyRootShelfById(controllerFunc types.ControllerFunc[*dtos.DeleteMyRootShelfByIdReqDto]) gin.HandlerFunc
	BindDeleteMyRootShelvesByIds(controllerFunc types.ControllerFunc[*dtos.DeleteMyRootShelvesByIdsReqDto]) gin.HandlerFunc
}

type RootShelfBinder struct{}

func NewRootShelfBinder() RootShelfBinderInterface {
	return &RootShelfBinder{}
}

/* ============================== Implementations ============================== */

func (b *RootShelfBinder) BindGetMyRootShelfById(controllerFunc types.ControllerFunc[*dtos.GetMyRootShelfByIdReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.GetMyRootShelfByIdReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		rootShelfIdString := ctx.Query("rootShelfId")
		if rootShelfIdString == "" {
			exceptions.Shelf.InvalidInput().WithError(fmt.Errorf("rootShelfId is required")).ResponseWithJSON(ctx)
			return
		}
		rootShelfId, err := uuid.Parse(rootShelfIdString)
		if err != nil {
			exceptions.Shelf.InvalidInput().WithError(err).ResponseWithJSON(ctx)
			return
		}
		reqDto.Param.RootShelfId = rootShelfId

		controllerFunc(ctx, &reqDto)
	}
}

func (b *RootShelfBinder) BindSearchRecentRootShelves(controllerFunc types.ControllerFunc[*dtos.SearchRecentRootShelvesReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.SearchRecentRootShelvesReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindQuery(&reqDto.Param); err != nil {
			exceptions.Shelf.InvalidInput().WithError(err).Log().ResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *RootShelfBinder) BindCreateRootShelf(controllerFunc types.ControllerFunc[*dtos.CreateRootShelfReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.CreateRootShelfReqDto

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

func (b *RootShelfBinder) BindUpdateMyRootShelfById(controllerFunc types.ControllerFunc[*dtos.UpdateMyRootShelfByIdReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.UpdateMyRootShelfByIdReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.Shelf.InvalidDto().WithError(err)
			exception.ResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *RootShelfBinder) BindRestoreMyRootShelfById(controllerFunc types.ControllerFunc[*dtos.RestoreMyRootShelfByIdReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.RestoreMyRootShelfByIdReqDto

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
func (b *RootShelfBinder) BindRestoreMyRootShelvesByIds(controllerFunc types.ControllerFunc[*dtos.RestoreMyRootShelvesByIdsReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.RestoreMyRootShelvesByIdsReqDto

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

func (b *RootShelfBinder) BindDeleteMyRootShelfById(controllerFunc types.ControllerFunc[*dtos.DeleteMyRootShelfByIdReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.DeleteMyRootShelfByIdReqDto

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

func (b *RootShelfBinder) BindDeleteMyRootShelvesByIds(controllerFunc types.ControllerFunc[*dtos.DeleteMyRootShelvesByIdsReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.DeleteMyRootShelvesByIdsReqDto

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
