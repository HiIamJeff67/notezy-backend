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

type BlockGroupBinderInterface interface {
	BindGetMyBlockGroupById(controllerFunc types.ControllerFunc[*dtos.GetMyBlockGroupByIdReqDto]) gin.HandlerFunc
	BindGetMyBlockGroupAndItsBlocksById(controllerFunc types.ControllerFunc[*dtos.GetMyBlockGroupAndItsBlocksByIdReqDto]) gin.HandlerFunc
	BindGetMyBlockGroupsAndTheirBlocksByBlockPackId(controllerFunc types.ControllerFunc[*dtos.GetMyBlockGroupsAndTheirBlocksByBlockPackIdReqDto]) gin.HandlerFunc
	BindGetMyBlockGroupsByPrevBlockGroupId(controllerFunc types.ControllerFunc[*dtos.GetMyBlockGroupsByPrevBlockGroupIdReqDto]) gin.HandlerFunc
	BindGetAllMyBlockGroupsByBlockPackId(controllerFunc types.ControllerFunc[*dtos.GetAllMyBlockGroupsByBlockPackIdReqDto]) gin.HandlerFunc
	BindInsertBlockGroupByBlockPackId(controllerFunc types.ControllerFunc[*dtos.InsertBlockGroupByBlockPackIdReqDto]) gin.HandlerFunc
	BindInsertBlockGroupAndItsBlocksByBlockPackId(controllerFunc types.ControllerFunc[*dtos.InsertBlockGroupAndItsBlocksByBlockPackIdReqDto]) gin.HandlerFunc
	BindInsertBlockGroupsAndTheirBlocksByBlockPackId(controllerFunc types.ControllerFunc[*dtos.InsertBlockGroupsAndTheirBlocksByBlockPackIdReqDto]) gin.HandlerFunc
	BindInsertSequentialBlockGroupsAndTheirBlocksByBlockPackId(controllerFunc types.ControllerFunc[*dtos.InsertSequentialBlockGroupsAndTheirBlocksByBlockPackIdReqDto]) gin.HandlerFunc
	BindMoveMyBlockGroupsByIds(controllerFunc types.ControllerFunc[*dtos.MoveMyBlockGroupsByIdsReqDto]) gin.HandlerFunc
	BindRestoreMyBlockGroupById(controllerFunc types.ControllerFunc[*dtos.RestoreMyBlockGroupByIdReqDto]) gin.HandlerFunc
	BindRestoreMyBlockGroupsByIds(controllerFunc types.ControllerFunc[*dtos.RestoreMyBlockGroupsByIdsReqDto]) gin.HandlerFunc
	BindDeleteMyBlockGroupById(controllerFunc types.ControllerFunc[*dtos.DeleteMyBlockGroupByIdReqDto]) gin.HandlerFunc
	BindDeleteMyBlockGroupsByIds(controllerFunc types.ControllerFunc[*dtos.DeleteMyBlockGroupsByIdsReqDto]) gin.HandlerFunc
}

type BlockGroupBinder struct{}

func NewBlockGroupBinder() BlockGroupBinderInterface {
	return &BlockGroupBinder{}
}

/* ============================== Implementations ============================== */

func (b *BlockGroupBinder) BindGetMyBlockGroupById(controllerFunc types.ControllerFunc[*dtos.GetMyBlockGroupByIdReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.GetMyBlockGroupByIdReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		blockGroupIdString := ctx.Query("blockGroupId")
		if blockGroupIdString == "" {
			exceptions.Shelf.InvalidInput().WithError(fmt.Errorf("blockGroupId is required")).Log().ResponseWithJSON(ctx)
			return
		}
		blockGroupId, err := uuid.Parse(blockGroupIdString)
		if err != nil {
			exceptions.Shelf.InvalidInput().WithError(err).Log().ResponseWithJSON(ctx)
			return
		}
		reqDto.Param.BlockGroupId = blockGroupId

		controllerFunc(ctx, &reqDto)
	}
}

func (b *BlockGroupBinder) BindGetMyBlockGroupAndItsBlocksById(controllerFunc types.ControllerFunc[*dtos.GetMyBlockGroupAndItsBlocksByIdReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.GetMyBlockGroupAndItsBlocksByIdReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		blockGroupIdString := ctx.Query("blockGroupId")
		if blockGroupIdString == "" {
			exceptions.Shelf.InvalidInput().WithError(fmt.Errorf("blockGroupId is required")).Log().ResponseWithJSON(ctx)
			return
		}
		blockGroupId, err := uuid.Parse(blockGroupIdString)
		if err != nil {
			exceptions.Shelf.InvalidInput().WithError(err).Log().ResponseWithJSON(ctx)
			return
		}
		reqDto.Param.BlockGroupId = blockGroupId

		controllerFunc(ctx, &reqDto)
	}
}

func (b *BlockGroupBinder) BindGetMyBlockGroupsAndTheirBlocksByBlockPackId(controllerFunc types.ControllerFunc[*dtos.GetMyBlockGroupsAndTheirBlocksByBlockPackIdReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.GetMyBlockGroupsAndTheirBlocksByBlockPackIdReqDto

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

func (b *BlockGroupBinder) BindGetMyBlockGroupsByPrevBlockGroupId(controllerFunc types.ControllerFunc[*dtos.GetMyBlockGroupsByPrevBlockGroupIdReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.GetMyBlockGroupsByPrevBlockGroupIdReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		prevBlockGroupIdString := ctx.Query("prevBlockGroupId")
		if prevBlockGroupIdString == "" {
			exceptions.Shelf.InvalidInput().WithError(fmt.Errorf("prevBlockGroupId is required")).Log().ResponseWithJSON(ctx)
			return
		}
		prevBlockGroupId, err := uuid.Parse(prevBlockGroupIdString)
		if err != nil {
			exceptions.Shelf.InvalidInput().WithError(err).Log().ResponseWithJSON(ctx)
			return
		}
		reqDto.Param.PrevBlockGroupId = prevBlockGroupId

		controllerFunc(ctx, &reqDto)
	}
}

func (b *BlockGroupBinder) BindGetAllMyBlockGroupsByBlockPackId(controllerFunc types.ControllerFunc[*dtos.GetAllMyBlockGroupsByBlockPackIdReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.GetAllMyBlockGroupsByBlockPackIdReqDto

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

func (b *BlockGroupBinder) BindInsertBlockGroupByBlockPackId(controllerFunc types.ControllerFunc[*dtos.InsertBlockGroupByBlockPackIdReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.InsertBlockGroupByBlockPackIdReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.BlockGroup.InvalidDto().WithError(err)
			exception.ResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *BlockGroupBinder) BindInsertBlockGroupAndItsBlocksByBlockPackId(controllerFunc types.ControllerFunc[*dtos.InsertBlockGroupAndItsBlocksByBlockPackIdReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.InsertBlockGroupAndItsBlocksByBlockPackIdReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.BlockGroup.InvalidDto().WithError(err)
			exception.ResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *BlockGroupBinder) BindInsertBlockGroupsAndTheirBlocksByBlockPackId(controllerFunc types.ControllerFunc[*dtos.InsertBlockGroupsAndTheirBlocksByBlockPackIdReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.InsertBlockGroupsAndTheirBlocksByBlockPackIdReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.BlockGroup.InvalidDto().WithError(err)
			exception.ResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *BlockGroupBinder) BindInsertSequentialBlockGroupsAndTheirBlocksByBlockPackId(controllerFunc types.ControllerFunc[*dtos.InsertSequentialBlockGroupsAndTheirBlocksByBlockPackIdReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.InsertSequentialBlockGroupsAndTheirBlocksByBlockPackIdReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.BlockGroup.InvalidDto().WithError(err)
			exception.ResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *BlockGroupBinder) BindMoveMyBlockGroupsByIds(controllerFunc types.ControllerFunc[*dtos.MoveMyBlockGroupsByIdsReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.MoveMyBlockGroupsByIdsReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.BlockGroup.InvalidDto().WithError(err)
			exception.ResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *BlockGroupBinder) BindRestoreMyBlockGroupById(controllerFunc types.ControllerFunc[*dtos.RestoreMyBlockGroupByIdReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.RestoreMyBlockGroupByIdReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.BlockGroup.InvalidDto().WithError(err)
			exception.ResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *BlockGroupBinder) BindRestoreMyBlockGroupsByIds(controllerFunc types.ControllerFunc[*dtos.RestoreMyBlockGroupsByIdsReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.RestoreMyBlockGroupsByIdsReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.BlockGroup.InvalidDto().WithError(err)
			exception.ResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *BlockGroupBinder) BindDeleteMyBlockGroupById(controllerFunc types.ControllerFunc[*dtos.DeleteMyBlockGroupByIdReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.DeleteMyBlockGroupByIdReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.BlockGroup.InvalidDto().WithError(err)
			exception.ResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *BlockGroupBinder) BindDeleteMyBlockGroupsByIds(controllerFunc types.ControllerFunc[*dtos.DeleteMyBlockGroupsByIdsReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.DeleteMyBlockGroupsByIdsReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, constants.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.BlockGroup.InvalidDto().WithError(err)
			exception.ResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}
