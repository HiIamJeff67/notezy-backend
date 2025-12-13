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
	BindGetMyBlockGroupAndItsBlocksById(controllerFunc types.ControllerFunc[*dtos.GetMyBlockGroupAndItsBlocksByIdReqDto]) gin.HandlerFunc
	BindCreateBlockGroupAndItsBlocksByBlockPackId(controllerFunc types.ControllerFunc[*dtos.CreateBlockGroupAndItsBlocksByBlockPackIdReqDto]) gin.HandlerFunc
}

type BlockGroupBinder struct{}

func NewBlockGroupBinder() BlockGroupBinderInterface {
	return &BlockGroupBinder{}
}

/* ============================== Implementations ============================== */

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

func (b *BlockGroupBinder) BindCreateBlockGroupAndItsBlocksByBlockPackId(controllerFunc types.ControllerFunc[*dtos.CreateBlockGroupAndItsBlocksByBlockPackIdReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.CreateBlockGroupAndItsBlocksByBlockPackIdReqDto

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
