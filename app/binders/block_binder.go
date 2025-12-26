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
	BindGetAllMyBlocks(controllerFunc types.ControllerFunc[*dtos.GetAllMyBlocksReqDto]) gin.HandlerFunc
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
			exceptions.Shelf.InvalidInput().WithError(fmt.Errorf("blockId is required")).Log().ResponseWithJSON(ctx)
			return
		}
		blockId, err := uuid.Parse(blockIdString)
		if err != nil {
			exceptions.Shelf.InvalidInput().WithError(err).Log().ResponseWithJSON(ctx)
			return
		}
		reqDto.Param.BlockId = blockId

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
