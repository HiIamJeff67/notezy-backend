package binders

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	contexts "github.com/HiIamJeff67/notezy-backend/app/contexts"
	dtos "github.com/HiIamJeff67/notezy-backend/app/dtos"
	exceptions "github.com/HiIamJeff67/notezy-backend/app/exceptions"
	types "github.com/HiIamJeff67/notezy-backend/shared/types"
)

type BlockBinderInterface interface {
	BindGetMyBlockById(controllerFunc types.ControllerFunc[*dtos.GetMyBlockByIdReqDto]) gin.HandlerFunc
	BindGetMyBlocksByIds(controllerFunc types.ControllerFunc[*dtos.GetMyBlocksByIdsReqDto]) gin.HandlerFunc
	BindGetMyBlocksByBlockPackId(controllerFunc types.ControllerFunc[*dtos.GetMyBlocksByBlockPackIdReqDto]) gin.HandlerFunc
}

type BlockBinder struct{}

func NewBlockBinder() BlockBinderInterface {
	return &BlockBinder{}
}

func (b *BlockBinder) BindGetMyBlockById(controllerFunc types.ControllerFunc[*dtos.GetMyBlockByIdReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.GetMyBlockByIdReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		blockId, err := uuid.Parse(ctx.Param("blockId"))
		if err != nil {
			exceptions.Block.InvalidDto().WithOrigin(err).Log().SafelyAbortAndResponseWithJSON(ctx)
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

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindQuery(&reqDto.Param); err != nil {
			exceptions.Block.InvalidDto().WithOrigin(err).Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *BlockBinder) BindGetMyBlocksByBlockPackId(controllerFunc types.ControllerFunc[*dtos.GetMyBlocksByBlockPackIdReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.GetMyBlocksByBlockPackIdReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		blockPackId, err := uuid.Parse(ctx.Param("blockPackId"))
		if err != nil {
			exceptions.Block.InvalidDto().WithOrigin(err).Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.Param.BlockPackId = blockPackId

		controllerFunc(ctx, &reqDto)
	}
}
