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

type BlockPackBinderInterface interface {
	BindGetMyBlockPackById(controllerFunc types.ControllerFunc[*dtos.GetMyBlockPackByIdReqDto]) gin.HandlerFunc
	BindGetMyBlockPackAndItsParentById(controllerFunc types.ControllerFunc[*dtos.GetMyBlockPackAndItsParentByIdReqDto]) gin.HandlerFunc
	BindGetAllMyBlockPacksByParentSubShelfId(controllerFunc types.ControllerFunc[*dtos.GetAllMyBlockPacksByParentSubShelfIdReqDto]) gin.HandlerFunc
	BindGetAllMyBlockPacksByRootShelfId(controllerFunc types.ControllerFunc[*dtos.GetAllMyBlockPacksByRootShelfIdReqDto]) gin.HandlerFunc
	BindCreateBlockPack(controllerFunc types.ControllerFunc[*dtos.CreateBlockPackReqDto]) gin.HandlerFunc
	BindUpdateMyBlockPackById(controllerFunc types.ControllerFunc[*dtos.UpdateMyBlockPackByIdReqDto]) gin.HandlerFunc
	BindMoveMyBlockPackById(controllerFunc types.ControllerFunc[*dtos.MoveMyBlockPackByIdReqDto]) gin.HandlerFunc
	BindMoveMyBlockPacksByIds(controllerFunc types.ControllerFunc[*dtos.MoveMyBlockPacksByIdsReqDto]) gin.HandlerFunc
	BindRestoreMyBlockPackById(controllerFunc types.ControllerFunc[*dtos.RestoreMyBlockPackByIdReqDto]) gin.HandlerFunc
	BindRestoreMyBlockPacksByIds(controllerFunc types.ControllerFunc[*dtos.RestoreMyBlockPacksByIdsReqDto]) gin.HandlerFunc
	BindDeleteMyBlockPackById(controllerFunc types.ControllerFunc[*dtos.DeleteMyBlockPackByIdReqDto]) gin.HandlerFunc
	BindDeleteMyBlockPacksByIds(controllerFunc types.ControllerFunc[*dtos.DeleteMyBlockPacksByIdsReqDto]) gin.HandlerFunc
}

type BlockPackBinder struct{}

func NewBlockPackBinder() BlockPackBinderInterface {
	return &BlockPackBinder{}
}

/* ============================== Implementations ============================== */

func (b *BlockPackBinder) BindGetMyBlockPackById(controllerFunc types.ControllerFunc[*dtos.GetMyBlockPackByIdReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.GetMyBlockPackByIdReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		blockPackIdString := ctx.Query("blockPackId")
		if blockPackIdString == "" {
			exceptions.Shelf.InvalidInput().WithError(fmt.Errorf("blockPackId is required")).Log().ResponseWithJSON(ctx)
			return
		}
		blockPackId, err := uuid.Parse(blockPackIdString)
		if err != nil {
			exceptions.Shelf.InvalidInput().WithError(err).Log().ResponseWithJSON(ctx)
			return
		}
		reqDto.Param.BlockPackId = blockPackId

		controllerFunc(ctx, &reqDto)
	}
}

func (b *BlockPackBinder) BindGetMyBlockPackAndItsParentById(controllerFunc types.ControllerFunc[*dtos.GetMyBlockPackAndItsParentByIdReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.GetMyBlockPackAndItsParentByIdReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		blockPackIdString := ctx.Query("blockPackId")
		if blockPackIdString == "" {
			exceptions.Shelf.InvalidInput().WithError(fmt.Errorf("blockPackId is required")).Log().ResponseWithJSON(ctx)
			return
		}
		blockPackId, err := uuid.Parse(blockPackIdString)
		if err != nil {
			exceptions.Shelf.InvalidInput().WithError(err).Log().ResponseWithJSON(ctx)
			return
		}
		reqDto.Param.BlockPackId = blockPackId

		controllerFunc(ctx, &reqDto)
	}
}

func (b *BlockPackBinder) BindGetAllMyBlockPacksByParentSubShelfId(controllerFunc types.ControllerFunc[*dtos.GetAllMyBlockPacksByParentSubShelfIdReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.GetAllMyBlockPacksByParentSubShelfIdReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		parentSubShelfIdString := ctx.Query("parentSubShelfId")
		if parentSubShelfIdString == "" {
			exceptions.Shelf.InvalidInput().WithError(fmt.Errorf("parentSubShelfId is required")).Log().ResponseWithJSON(ctx)
			return
		}
		parentSubShelfId, err := uuid.Parse(parentSubShelfIdString)
		if err != nil {
			exceptions.Shelf.InvalidInput().WithError(err).Log().ResponseWithJSON(ctx)
			return
		}
		reqDto.Param.ParentSubShelfId = parentSubShelfId

		controllerFunc(ctx, &reqDto)
	}
}

func (b *BlockPackBinder) BindGetAllMyBlockPacksByRootShelfId(controllerFunc types.ControllerFunc[*dtos.GetAllMyBlockPacksByRootShelfIdReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.GetAllMyBlockPacksByRootShelfIdReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		rootShelfIdString := ctx.Query("rootShelfId")
		if rootShelfIdString == "" {
			exceptions.Shelf.InvalidInput().WithError(fmt.Errorf("rootShelfId is required")).Log().ResponseWithJSON(ctx)
			return
		}
		rootShelfId, err := uuid.Parse(rootShelfIdString)
		if err != nil {
			exceptions.Shelf.InvalidInput().WithError(err).Log().ResponseWithJSON(ctx)
			return
		}
		reqDto.Param.RootShelfId = rootShelfId

		controllerFunc(ctx, &reqDto)
	}
}

func (b *BlockPackBinder) BindCreateBlockPack(controllerFunc types.ControllerFunc[*dtos.CreateBlockPackReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.CreateBlockPackReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.BlockPack.InvalidDto().WithError(err)
			exception.ResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *BlockPackBinder) BindUpdateMyBlockPackById(controllerFunc types.ControllerFunc[*dtos.UpdateMyBlockPackByIdReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.UpdateMyBlockPackByIdReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.BlockPack.InvalidDto().WithError(err)
			exception.ResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *BlockPackBinder) BindMoveMyBlockPackById(controllerFunc types.ControllerFunc[*dtos.MoveMyBlockPackByIdReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.MoveMyBlockPackByIdReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.BlockPack.InvalidDto().WithError(err)
			exception.ResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *BlockPackBinder) BindMoveMyBlockPacksByIds(controllerFunc types.ControllerFunc[*dtos.MoveMyBlockPacksByIdsReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.MoveMyBlockPacksByIdsReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.BlockPack.InvalidDto().WithError(err)
			exception.ResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *BlockPackBinder) BindRestoreMyBlockPackById(controllerFunc types.ControllerFunc[*dtos.RestoreMyBlockPackByIdReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.RestoreMyBlockPackByIdReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.BlockPack.InvalidDto().WithError(err)
			exception.ResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *BlockPackBinder) BindRestoreMyBlockPacksByIds(controllerFunc types.ControllerFunc[*dtos.RestoreMyBlockPacksByIdsReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.RestoreMyBlockPacksByIdsReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.BlockPack.InvalidDto().WithError(err)
			exception.ResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *BlockPackBinder) BindDeleteMyBlockPackById(controllerFunc types.ControllerFunc[*dtos.DeleteMyBlockPackByIdReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.DeleteMyBlockPackByIdReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.BlockPack.InvalidDto().WithError(err)
			exception.ResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *BlockPackBinder) BindDeleteMyBlockPacksByIds(controllerFunc types.ControllerFunc[*dtos.DeleteMyBlockPacksByIdsReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.DeleteMyBlockPacksByIdsReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.BlockPack.InvalidDto().WithError(err)
			exception.ResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}
