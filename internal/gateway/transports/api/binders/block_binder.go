package binders

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	exceptions "github.com/HiIamJeff67/notezy-backend/contracts/types/exceptions"

	exceptionwriter "github.com/HiIamJeff67/notezy-backend/shared/util/exceptionwriter"

	apicontract "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api/blocks"

	controllers "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/api/controllers"
)

type BlockBinderInterface interface {
	BindGetMyBlockById(controllerFunc controllers.Func[*apicontract.GetMyBlockByIdRequestDto]) gin.HandlerFunc
	BindGetMyBlocksByIds(controllerFunc controllers.Func[*apicontract.GetMyBlocksByIdsRequestDto]) gin.HandlerFunc
	BindGetMyBlocksByBlockPackId(controllerFunc controllers.Func[*apicontract.GetMyBlocksByBlockPackIdRequestDto]) gin.HandlerFunc
}

type BlockBinder struct{}

func NewBlockBinder() BlockBinderInterface {
	return &BlockBinder{}
}

func (b *BlockBinder) BindGetMyBlockById(controllerFunc controllers.Func[*apicontract.GetMyBlockByIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &apicontract.GetMyBlockByIdRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		blockId, err := uuid.Parse(ctx.Param("block-id"))
		if err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Block").WithOrigin(err), ctx)
			return
		}
		requestDto.Param.BlockId = blockId

		controllerFunc(ctx, requestDto)
	}
}

func (b *BlockBinder) BindGetMyBlocksByIds(controllerFunc controllers.Func[*apicontract.GetMyBlocksByIdsRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &apicontract.GetMyBlocksByIdsRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		if err := ctx.ShouldBindQuery(&requestDto.Param); err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("Block").WithOrigin(err), ctx)
			return
		}

		controllerFunc(ctx, requestDto)
	}
}

func (b *BlockBinder) BindGetMyBlocksByBlockPackId(controllerFunc controllers.Func[*apicontract.GetMyBlocksByBlockPackIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &apicontract.GetMyBlocksByBlockPackIdRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		blockPackId, err := uuid.Parse(ctx.Param("block-pack-id"))
		if err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Block").WithOrigin(err), ctx)
			return
		}
		requestDto.Param.BlockPackId = blockPackId

		controllerFunc(ctx, requestDto)
	}
}
