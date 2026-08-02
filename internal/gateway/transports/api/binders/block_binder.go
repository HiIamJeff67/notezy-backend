package binders

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	blocksdto "github.com/HiIamJeff67/notezy-backend/contracts/gateway/v1/api/blocks"
	exceptions "github.com/HiIamJeff67/notezy-backend/internal/exceptions"
	responsewriter "github.com/HiIamJeff67/notezy-backend/internal/shared/responsewriter"
	controllers "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/api/controllers"
)

type BlockBinderInterface interface {
	BindGetMyBlockById(controllerFunc controllers.Func[*blocksdto.GetMyBlockByIdRequestDto]) gin.HandlerFunc
	BindGetMyBlocksByIds(controllerFunc controllers.Func[*blocksdto.GetMyBlocksByIdsRequestDto]) gin.HandlerFunc
	BindGetMyBlocksByBlockPackId(controllerFunc controllers.Func[*blocksdto.GetMyBlocksByBlockPackIdRequestDto]) gin.HandlerFunc
}

type BlockBinder struct{}

func NewBlockBinder() BlockBinderInterface {
	return &BlockBinder{}
}

func (b *BlockBinder) BindGetMyBlockById(controllerFunc controllers.Func[*blocksdto.GetMyBlockByIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &blocksdto.GetMyBlockByIdRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		blockId, err := uuid.Parse(ctx.Param("blockId"))
		if err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Block").WithOrigin(err), ctx)
			return
		}
		requestDto.Param.BlockId = blockId

		controllerFunc(ctx, requestDto)
	}
}

func (b *BlockBinder) BindGetMyBlocksByIds(controllerFunc controllers.Func[*blocksdto.GetMyBlocksByIdsRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &blocksdto.GetMyBlocksByIdsRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		if err := ctx.ShouldBindQuery(&requestDto.Param); err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("Block").WithOrigin(err), ctx)
			return
		}

		controllerFunc(ctx, requestDto)
	}
}

func (b *BlockBinder) BindGetMyBlocksByBlockPackId(controllerFunc controllers.Func[*blocksdto.GetMyBlocksByBlockPackIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &blocksdto.GetMyBlocksByBlockPackIdRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		blockPackId, err := uuid.Parse(ctx.Param("blockPackId"))
		if err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Block").WithOrigin(err), ctx)
			return
		}
		requestDto.Param.BlockPackId = blockPackId

		controllerFunc(ctx, requestDto)
	}
}
