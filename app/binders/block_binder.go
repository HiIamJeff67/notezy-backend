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

type BlockBinderInterface interface {
	BindGetMyBlockById(controllerFunc types.ControllerFunc[*dtos.GetMyBlockByIdReqDto]) gin.HandlerFunc
	BindGetMyBlocksByIds(controllerFunc types.ControllerFunc[*dtos.GetMyBlocksByIdsReqDto]) gin.HandlerFunc
	BindGetMyBlocksByBlockGroupId(controllerFunc types.ControllerFunc[*dtos.GetMyBlocksByBlockGroupIdReqDto]) gin.HandlerFunc
	BindGetMyBlocksByBlockGroupIds(controllerFunc types.ControllerFunc[*dtos.GetMyBlocksByBlockGroupIdsReqDto]) gin.HandlerFunc
	BindGetMyBlocksByBlockPackId(controllerFunc types.ControllerFunc[*dtos.GetMyBlocksByBlockPackIdReqDto]) gin.HandlerFunc
	BindGetAllMyBlocks(controllerFunc types.ControllerFunc[*dtos.GetAllMyBlocksReqDto]) gin.HandlerFunc
	BindInsertBlock(controllerFunc types.ControllerFunc[*dtos.InsertBlockReqDto]) gin.HandlerFunc
	BindInsertBlocks(controllerFunc types.ControllerFunc[*dtos.InsertBlocksReqDto]) gin.HandlerFunc
	BindUpdateMyBlockById(controllerFunc types.ControllerFunc[*dtos.UpdateMyBlockByIdReqDto]) gin.HandlerFunc
	BindUpdateMyBlocksByIds(controllerFunc types.ControllerFunc[*dtos.UpdateMyBlocksByIdsReqDto]) gin.HandlerFunc
	BindRestoreMyBlockById(controllerFunc types.ControllerFunc[*dtos.RestoreMyBlockByIdReqDto]) gin.HandlerFunc
	BindRestoreMyBlocksByIds(controllerFunc types.ControllerFunc[*dtos.RestoreMyBlocksByIdsReqDto]) gin.HandlerFunc
	BindDeleteMyBlockById(controllerFunc types.ControllerFunc[*dtos.DeleteMyBlockByIdReqDto]) gin.HandlerFunc
	BindDeleteMyBlocksByIds(controllerFunc types.ControllerFunc[*dtos.DeleteMyBlocksByIdsReqDto]) gin.HandlerFunc
}

type BlockBinder struct{}

func NewBlockBinder() BlockBinderInterface {
	return &BlockBinder{}
}

/* ============================== Implementations ============================== */

func (b *BlockBinder) BindGetMyBlockById(controllerFunc types.ControllerFunc[*dtos.GetMyBlockByIdReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.GetMyBlockByIdReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		blockIdString := ctx.Query("blockId")
		if blockIdString == "" {
			exceptions.Shelf.InvalidDto().WithError(fmt.Errorf("blockId is required")).Log().ResponseWithJSON(ctx)
			return
		}
		blockId, err := uuid.Parse(blockIdString)
		if err != nil {
			exceptions.Shelf.InvalidDto().WithError(err).Log().ResponseWithJSON(ctx)
			return
		}
		reqDto.Param.BlockId = blockId

		controllerFunc(ctx, &reqDto)
	}
}

func (b *BlockBinder) BindGetMyBlocksByIds(controllerFunc types.ControllerFunc[*dtos.GetMyBlocksByIdsReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.GetMyBlocksByIdsReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindQuery(&reqDto.Param); err != nil {
			exceptions.Shelf.InvalidDto().WithError(err).Log().ResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *BlockBinder) BindGetMyBlocksByBlockGroupId(controllerFunc types.ControllerFunc[*dtos.GetMyBlocksByBlockGroupIdReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.GetMyBlocksByBlockGroupIdReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		blockGroupIdString := ctx.Query("blockGroupId")
		if blockGroupIdString == "" {
			exceptions.Shelf.InvalidDto().WithError(fmt.Errorf("blockGroupId is required")).Log().ResponseWithJSON(ctx)
			return
		}
		blockGroupId, err := uuid.Parse(blockGroupIdString)
		if err != nil {
			exceptions.Shelf.InvalidDto().WithError(err).Log().ResponseWithJSON(ctx)
			return
		}
		reqDto.Param.BlockGroupId = blockGroupId

		controllerFunc(ctx, &reqDto)
	}
}

func (b *BlockBinder) BindGetMyBlocksByBlockGroupIds(controllerFunc types.ControllerFunc[*dtos.GetMyBlocksByBlockGroupIdsReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.GetMyBlocksByBlockGroupIdsReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindQuery(&reqDto.Param); err != nil {
			exceptions.Shelf.InvalidDto().WithError(err).Log().ResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *BlockBinder) BindGetMyBlocksByBlockPackId(controllerFunc types.ControllerFunc[*dtos.GetMyBlocksByBlockPackIdReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.GetMyBlocksByBlockPackIdReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		blockPackIdString := ctx.Query("blockPackId")
		if blockPackIdString == "" {
			exceptions.Shelf.InvalidDto().WithError(fmt.Errorf("blockPackId is required")).Log().ResponseWithJSON(ctx)
			return
		}
		blockPackId, err := uuid.Parse(blockPackIdString)
		if err != nil {
			exceptions.Shelf.InvalidDto().WithError(err).Log().ResponseWithJSON(ctx)
			return
		}
		reqDto.Param.BlockPackId = blockPackId

		controllerFunc(ctx, &reqDto)
	}
}

func (b *BlockBinder) BindGetAllMyBlocks(controllerFunc types.ControllerFunc[*dtos.GetAllMyBlocksReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.GetAllMyBlocksReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		controllerFunc(ctx, &reqDto)
	}
}

func (b *BlockBinder) BindInsertBlock(controllerFunc types.ControllerFunc[*dtos.InsertBlockReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.InsertBlockReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exceptions.Shelf.InvalidDto().WithError(err).Log().ResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *BlockBinder) BindInsertBlocks(controllerFunc types.ControllerFunc[*dtos.InsertBlocksReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.InsertBlocksReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exceptions.Shelf.InvalidDto().WithError(err).Log().ResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *BlockBinder) BindUpdateMyBlockById(controllerFunc types.ControllerFunc[*dtos.UpdateMyBlockByIdReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.UpdateMyBlockByIdReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exceptions.Shelf.InvalidDto().WithError(err).Log().ResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *BlockBinder) BindUpdateMyBlocksByIds(controllerFunc types.ControllerFunc[*dtos.UpdateMyBlocksByIdsReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.UpdateMyBlocksByIdsReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exceptions.Shelf.InvalidDto().WithError(err).Log().ResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *BlockBinder) BindRestoreMyBlockById(controllerFunc types.ControllerFunc[*dtos.RestoreMyBlockByIdReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.RestoreMyBlockByIdReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exceptions.Shelf.InvalidDto().WithError(err).Log().ResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *BlockBinder) BindRestoreMyBlocksByIds(controllerFunc types.ControllerFunc[*dtos.RestoreMyBlocksByIdsReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.RestoreMyBlocksByIdsReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exceptions.Shelf.InvalidDto().WithError(err).Log().ResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *BlockBinder) BindDeleteMyBlockById(controllerFunc types.ControllerFunc[*dtos.DeleteMyBlockByIdReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.DeleteMyBlockByIdReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exceptions.Shelf.InvalidDto().WithError(err).Log().ResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *BlockBinder) BindDeleteMyBlocksByIds(controllerFunc types.ControllerFunc[*dtos.DeleteMyBlocksByIdsReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.DeleteMyBlocksByIdsReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exceptions.Shelf.InvalidDto().WithError(err).Log().ResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}
