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

type BlockPackBinderInterface interface {
	BindGetMyBlockPackById(controllerFunc types.ControllerFunc[*dtos.GetMyBlockPackByIdReqDto]) gin.HandlerFunc
	BindGetMyBlockPackAndItsParentById(controllerFunc types.ControllerFunc[*dtos.GetMyBlockPackAndItsParentByIdReqDto]) gin.HandlerFunc
	BindGetMyBlockPacksByParentSubShelfId(controllerFunc types.ControllerFunc[*dtos.GetMyBlockPacksByParentSubShelfIdReqDto]) gin.HandlerFunc
	BindGetAllMyBlockPacksByRootShelfId(controllerFunc types.ControllerFunc[*dtos.GetAllMyBlockPacksByRootShelfIdReqDto]) gin.HandlerFunc
	BindCreateBlockPack(controllerFunc types.ControllerFunc[*dtos.CreateBlockPackReqDto]) gin.HandlerFunc
	BindCreateBlockPacks(controllerFunc types.ControllerFunc[*dtos.CreateBlockPacksReqDto]) gin.HandlerFunc
	BindUpdateMyBlockPackById(controllerFunc types.ControllerFunc[*dtos.UpdateMyBlockPackByIdReqDto]) gin.HandlerFunc
	BindUpdateMyBlockPacksByIds(controllerFunc types.ControllerFunc[*dtos.UpdateMyBlockPacksByIdsReqDto]) gin.HandlerFunc
	BindMoveMyBlockPackById(controllerFunc types.ControllerFunc[*dtos.MoveMyBlockPackByIdReqDto]) gin.HandlerFunc
	BindMoveMyBlockPacksByIds(controllerFunc types.ControllerFunc[*dtos.MoveMyBlockPacksByIdsReqDto]) gin.HandlerFunc
	BindBatchMoveMyBlockPacksByIds(controllerFunc types.ControllerFunc[*dtos.BatchMoveMyBlockPacksByIdsReqDto]) gin.HandlerFunc
	BindRestoreMyBlockPackById(controllerFunc types.ControllerFunc[*dtos.RestoreMyBlockPackByIdReqDto]) gin.HandlerFunc
	BindRestoreMyBlockPacksByIds(controllerFunc types.ControllerFunc[*dtos.RestoreMyBlockPacksByIdsReqDto]) gin.HandlerFunc
	BindDeleteMyBlockPackById(controllerFunc types.ControllerFunc[*dtos.DeleteMyBlockPackByIdReqDto]) gin.HandlerFunc
	BindDeleteMyBlockPacksByIds(controllerFunc types.ControllerFunc[*dtos.DeleteMyBlockPacksByIdsReqDto]) gin.HandlerFunc
}

type BlockPackBinder struct{}

func NewBlockPackBinder() BlockPackBinderInterface {
	return &BlockPackBinder{}
}

func (b *BlockPackBinder) BindGetMyBlockPackById(controllerFunc types.ControllerFunc[*dtos.GetMyBlockPackByIdReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.GetMyBlockPackByIdReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		blockPackIdString := ctx.Query("blockPackId")
		if blockPackIdString == "" {
			exceptions.Shelf.InvalidInput().WithOrigin(fmt.Errorf("blockPackId is required")).Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		blockPackId, err := uuid.Parse(blockPackIdString)
		if err != nil {
			exceptions.Shelf.InvalidInput().WithOrigin(err).Log().SafelyAbortAndResponseWithJSON(ctx)
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

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		blockPackIdString := ctx.Query("blockPackId")
		if blockPackIdString == "" {
			exceptions.Shelf.InvalidInput().WithOrigin(fmt.Errorf("blockPackId is required")).Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		blockPackId, err := uuid.Parse(blockPackIdString)
		if err != nil {
			exceptions.Shelf.InvalidInput().WithOrigin(err).Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.Param.BlockPackId = blockPackId

		controllerFunc(ctx, &reqDto)
	}
}

func (b *BlockPackBinder) BindGetMyBlockPacksByParentSubShelfId(controllerFunc types.ControllerFunc[*dtos.GetMyBlockPacksByParentSubShelfIdReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.GetMyBlockPacksByParentSubShelfIdReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		parentSubShelfIdString := ctx.Query("parentSubShelfId")
		if parentSubShelfIdString == "" {
			exceptions.Shelf.InvalidInput().WithOrigin(fmt.Errorf("parentSubShelfId is required")).Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		parentSubShelfId, err := uuid.Parse(parentSubShelfIdString)
		if err != nil {
			exceptions.Shelf.InvalidInput().WithOrigin(err).Log().SafelyAbortAndResponseWithJSON(ctx)
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

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		rootShelfIdString := ctx.Query("rootShelfId")
		if rootShelfIdString == "" {
			exceptions.Shelf.InvalidInput().WithOrigin(fmt.Errorf("rootShelfId is required")).Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		rootShelfId, err := uuid.Parse(rootShelfIdString)
		if err != nil {
			exceptions.Shelf.InvalidInput().WithOrigin(err).Log().SafelyAbortAndResponseWithJSON(ctx)
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

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.BlockPack.InvalidDto().WithOrigin(err)
			exception.SafelyAbortAndResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *BlockPackBinder) BindCreateBlockPacks(controllerFunc types.ControllerFunc[*dtos.CreateBlockPacksReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.CreateBlockPacksReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.BlockPack.InvalidDto().WithOrigin(err)
			exception.SafelyAbortAndResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *BlockPackBinder) BindUpdateMyBlockPackById(controllerFunc types.ControllerFunc[*dtos.UpdateMyBlockPackByIdReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.UpdateMyBlockPackByIdReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.BlockPack.InvalidDto().WithOrigin(err)
			exception.SafelyAbortAndResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *BlockPackBinder) BindUpdateMyBlockPacksByIds(controllerFunc types.ControllerFunc[*dtos.UpdateMyBlockPacksByIdsReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.UpdateMyBlockPacksByIdsReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.BlockPack.InvalidDto().WithOrigin(err)
			exception.SafelyAbortAndResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *BlockPackBinder) BindMoveMyBlockPackById(controllerFunc types.ControllerFunc[*dtos.MoveMyBlockPackByIdReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.MoveMyBlockPackByIdReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.BlockPack.InvalidDto().WithOrigin(err)
			exception.SafelyAbortAndResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *BlockPackBinder) BindMoveMyBlockPacksByIds(controllerFunc types.ControllerFunc[*dtos.MoveMyBlockPacksByIdsReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.MoveMyBlockPacksByIdsReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.BlockPack.InvalidDto().WithOrigin(err)
			exception.SafelyAbortAndResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *BlockPackBinder) BindBatchMoveMyBlockPacksByIds(controllerFunc types.ControllerFunc[*dtos.BatchMoveMyBlockPacksByIdsReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.BatchMoveMyBlockPacksByIdsReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.BlockPack.InvalidDto().WithOrigin(err)
			exception.SafelyAbortAndResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *BlockPackBinder) BindRestoreMyBlockPackById(controllerFunc types.ControllerFunc[*dtos.RestoreMyBlockPackByIdReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.RestoreMyBlockPackByIdReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.BlockPack.InvalidDto().WithOrigin(err)
			exception.SafelyAbortAndResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *BlockPackBinder) BindRestoreMyBlockPacksByIds(controllerFunc types.ControllerFunc[*dtos.RestoreMyBlockPacksByIdsReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.RestoreMyBlockPacksByIdsReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.BlockPack.InvalidDto().WithOrigin(err)
			exception.SafelyAbortAndResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *BlockPackBinder) BindDeleteMyBlockPackById(controllerFunc types.ControllerFunc[*dtos.DeleteMyBlockPackByIdReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.DeleteMyBlockPackByIdReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.BlockPack.InvalidDto().WithOrigin(err)
			exception.SafelyAbortAndResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *BlockPackBinder) BindDeleteMyBlockPacksByIds(controllerFunc types.ControllerFunc[*dtos.DeleteMyBlockPacksByIdsReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.DeleteMyBlockPacksByIdsReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.BlockPack.InvalidDto().WithOrigin(err)
			exception.SafelyAbortAndResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}
