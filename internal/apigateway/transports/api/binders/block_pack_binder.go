package binders

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	exceptions "github.com/HiIamJeff67/notezy-backend/contracts/types/exceptions"

	exceptionwriter "github.com/HiIamJeff67/notezy-backend/shared/util/exceptionwriter"

	apicontract "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api/block-packs"

	controllers "github.com/HiIamJeff67/notezy-backend/internal/apigateway/transports/api/controllers"
)

type BlockPackBinderInterface interface {
	BindGetMyBlockPackById(controllerFunc controllers.Func[*apicontract.GetMyBlockPackByIdRequestDto]) gin.HandlerFunc
	BindGetMyBlockPackAndItsParentById(controllerFunc controllers.Func[*apicontract.GetMyBlockPackAndItsParentByIdRequestDto]) gin.HandlerFunc
	BindGetMyBlockPacksByParentSubShelfId(controllerFunc controllers.Func[*apicontract.GetMyBlockPacksByParentSubShelfIdRequestDto]) gin.HandlerFunc
	BindGetAllMyBlockPacksByRootShelfId(controllerFunc controllers.Func[*apicontract.GetAllMyBlockPacksByRootShelfIdRequestDto]) gin.HandlerFunc
	BindCreateBlockPack(controllerFunc controllers.Func[*apicontract.CreateBlockPackRequestDto]) gin.HandlerFunc
	BindCreateBlockPacks(controllerFunc controllers.Func[*apicontract.CreateBlockPacksRequestDto]) gin.HandlerFunc
	BindUpdateMyBlockPackById(controllerFunc controllers.Func[*apicontract.UpdateMyBlockPackByIdRequestDto]) gin.HandlerFunc
	BindUpdateMyBlockPacksByIds(controllerFunc controllers.Func[*apicontract.UpdateMyBlockPacksByIdsRequestDto]) gin.HandlerFunc
	BindMoveMyBlockPackByParentSubShelfId(controllerFunc controllers.Func[*apicontract.MoveMyBlockPackByParentSubShelfIdRequestDto]) gin.HandlerFunc
	BindMoveMyBlockPacksByParentSubShelfId(controllerFunc controllers.Func[*apicontract.MoveMyBlockPacksByParentSubShelfIdRequestDto]) gin.HandlerFunc
	BindMoveMyBlockPacksByParentSubShelfIds(controllerFunc controllers.Func[*apicontract.MoveMyBlockPacksByParentSubShelfIdsRequestDto]) gin.HandlerFunc
	BindRestoreMyBlockPackById(controllerFunc controllers.Func[*apicontract.RestoreMyBlockPackByIdRequestDto]) gin.HandlerFunc
	BindRestoreMyBlockPacksByIds(controllerFunc controllers.Func[*apicontract.RestoreMyBlockPacksByIdsRequestDto]) gin.HandlerFunc
	BindDeleteMyBlockPackById(controllerFunc controllers.Func[*apicontract.DeleteMyBlockPackByIdRequestDto]) gin.HandlerFunc
	BindDeleteMyBlockPacksByIds(controllerFunc controllers.Func[*apicontract.DeleteMyBlockPacksByIdsRequestDto]) gin.HandlerFunc
}

type BlockPackBinder struct{}

func NewBlockPackBinder() BlockPackBinderInterface {
	return &BlockPackBinder{}
}

func (b *BlockPackBinder) BindGetMyBlockPackById(controllerFunc controllers.Func[*apicontract.GetMyBlockPackByIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestDto apicontract.GetMyBlockPackByIdRequestDto

		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		isDeletedString := ctx.Query("isDeleted")
		if isDeletedString != "" {
			value, err := strconv.ParseBool(isDeletedString)
			if err != nil {
				exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("BlockPack").WithOrigin(err), ctx)
				return
			}
			requestDto.Param.IsDeleted = &value
		}

		value, err := uuid.Parse(ctx.Param("block-pack-id"))
		if err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("BlockPack").WithOrigin(err), ctx)
			return
		}
		requestDto.Param.BlockPackId = value

		controllerFunc(ctx, &requestDto)
	}
}

func (b *BlockPackBinder) BindGetMyBlockPackAndItsParentById(controllerFunc controllers.Func[*apicontract.GetMyBlockPackAndItsParentByIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestDto apicontract.GetMyBlockPackAndItsParentByIdRequestDto

		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		isDeletedString := ctx.Query("isDeleted")
		if isDeletedString != "" {
			value, err := strconv.ParseBool(isDeletedString)
			if err != nil {
				exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("BlockPack").WithOrigin(err), ctx)
				return
			}
			requestDto.Param.IsDeleted = &value
		}

		value, err := uuid.Parse(ctx.Param("block-pack-id"))
		if err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("BlockPack").WithOrigin(err), ctx)
			return
		}
		requestDto.Param.BlockPackId = value

		controllerFunc(ctx, &requestDto)
	}
}

func (b *BlockPackBinder) BindGetMyBlockPacksByParentSubShelfId(controllerFunc controllers.Func[*apicontract.GetMyBlockPacksByParentSubShelfIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestDto apicontract.GetMyBlockPacksByParentSubShelfIdRequestDto

		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		areDeletedString := ctx.Query("areDeleted")
		if areDeletedString != "" {
			value, err := strconv.ParseBool(areDeletedString)
			if err != nil {
				exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("BlockPack").WithOrigin(err), ctx)
				return
			}
			requestDto.Param.AreDeleted = &value
		}

		value, err := uuid.Parse(ctx.Param("parent-sub-shelf-id"))
		if err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("BlockPack").WithOrigin(err), ctx)
			return
		}
		requestDto.Param.ParentSubShelfId = value

		controllerFunc(ctx, &requestDto)
	}
}

func (b *BlockPackBinder) BindGetAllMyBlockPacksByRootShelfId(controllerFunc controllers.Func[*apicontract.GetAllMyBlockPacksByRootShelfIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestDto apicontract.GetAllMyBlockPacksByRootShelfIdRequestDto

		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		areDeletedString := ctx.Query("areDeleted")
		if areDeletedString != "" {
			value, err := strconv.ParseBool(areDeletedString)
			if err != nil {
				exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("BlockPack").WithOrigin(err), ctx)
				return
			}
			requestDto.Param.AreDeleted = &value
		}

		value, err := uuid.Parse(ctx.Param("root-shelf-id"))
		if err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("BlockPack").WithOrigin(err), ctx)
			return
		}
		requestDto.Param.RootShelfId = value

		controllerFunc(ctx, &requestDto)
	}
}

func (b *BlockPackBinder) BindCreateBlockPack(controllerFunc controllers.Func[*apicontract.CreateBlockPackRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestDto apicontract.CreateBlockPackRequestDto

		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			exception := exceptions.InvalidDto("BlockPack").WithOrigin(err)
			exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
			return
		}

		value, err := uuid.Parse(ctx.Param("parent-sub-shelf-id"))
		if err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("BlockPack").WithOrigin(err), ctx)
			return
		}
		requestDto.Body.ParentSubShelfId = value

		controllerFunc(ctx, &requestDto)
	}
}

func (b *BlockPackBinder) BindCreateBlockPacks(controllerFunc controllers.Func[*apicontract.CreateBlockPacksRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestDto apicontract.CreateBlockPacksRequestDto

		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			exception := exceptions.InvalidDto("BlockPack").WithOrigin(err)
			exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
			return
		}

		controllerFunc(ctx, &requestDto)
	}
}

func (b *BlockPackBinder) BindUpdateMyBlockPackById(controllerFunc controllers.Func[*apicontract.UpdateMyBlockPackByIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestDto apicontract.UpdateMyBlockPackByIdRequestDto

		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			exception := exceptions.InvalidDto("BlockPack").WithOrigin(err)
			exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
			return
		}

		value, err := uuid.Parse(ctx.Param("block-pack-id"))
		if err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("BlockPack").WithOrigin(err), ctx)
			return
		}
		requestDto.Param.BlockPackId = value

		controllerFunc(ctx, &requestDto)
	}
}

func (b *BlockPackBinder) BindUpdateMyBlockPacksByIds(controllerFunc controllers.Func[*apicontract.UpdateMyBlockPacksByIdsRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestDto apicontract.UpdateMyBlockPacksByIdsRequestDto

		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			exception := exceptions.InvalidDto("BlockPack").WithOrigin(err)
			exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
			return
		}

		controllerFunc(ctx, &requestDto)
	}
}

func (b *BlockPackBinder) BindMoveMyBlockPackByParentSubShelfId(controllerFunc controllers.Func[*apicontract.MoveMyBlockPackByParentSubShelfIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestDto apicontract.MoveMyBlockPackByParentSubShelfIdRequestDto

		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			exception := exceptions.InvalidDto("BlockPack").WithOrigin(err)
			exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
			return
		}

		value, err := uuid.Parse(ctx.Param("block-pack-id"))
		if err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("BlockPack").WithOrigin(err), ctx)
			return
		}
		requestDto.Body.BlockPackId = value

		controllerFunc(ctx, &requestDto)
	}
}

func (b *BlockPackBinder) BindMoveMyBlockPacksByParentSubShelfId(controllerFunc controllers.Func[*apicontract.MoveMyBlockPacksByParentSubShelfIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestDto apicontract.MoveMyBlockPacksByParentSubShelfIdRequestDto

		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			exception := exceptions.InvalidDto("BlockPack").WithOrigin(err)
			exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
			return
		}

		controllerFunc(ctx, &requestDto)
	}
}

func (b *BlockPackBinder) BindMoveMyBlockPacksByParentSubShelfIds(controllerFunc controllers.Func[*apicontract.MoveMyBlockPacksByParentSubShelfIdsRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestDto apicontract.MoveMyBlockPacksByParentSubShelfIdsRequestDto

		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			exception := exceptions.InvalidDto("BlockPack").WithOrigin(err)
			exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
			return
		}

		controllerFunc(ctx, &requestDto)
	}
}

func (b *BlockPackBinder) BindRestoreMyBlockPackById(controllerFunc controllers.Func[*apicontract.RestoreMyBlockPackByIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestDto apicontract.RestoreMyBlockPackByIdRequestDto

		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		value, err := uuid.Parse(ctx.Param("block-pack-id"))
		if err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("BlockPack").WithOrigin(err), ctx)
			return
		}
		requestDto.Param.BlockPackId = value

		controllerFunc(ctx, &requestDto)
	}
}

func (b *BlockPackBinder) BindRestoreMyBlockPacksByIds(controllerFunc controllers.Func[*apicontract.RestoreMyBlockPacksByIdsRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestDto apicontract.RestoreMyBlockPacksByIdsRequestDto

		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			exception := exceptions.InvalidDto("BlockPack").WithOrigin(err)
			exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
			return
		}

		controllerFunc(ctx, &requestDto)
	}
}

func (b *BlockPackBinder) BindDeleteMyBlockPackById(controllerFunc controllers.Func[*apicontract.DeleteMyBlockPackByIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestDto apicontract.DeleteMyBlockPackByIdRequestDto

		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		value, err := uuid.Parse(ctx.Param("block-pack-id"))
		if err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("BlockPack").WithOrigin(err), ctx)
			return
		}
		requestDto.Param.BlockPackId = value

		controllerFunc(ctx, &requestDto)
	}
}

func (b *BlockPackBinder) BindDeleteMyBlockPacksByIds(controllerFunc controllers.Func[*apicontract.DeleteMyBlockPacksByIdsRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestDto apicontract.DeleteMyBlockPacksByIdsRequestDto

		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			exception := exceptions.InvalidDto("BlockPack").WithOrigin(err)
			exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
			return
		}

		controllerFunc(ctx, &requestDto)
	}
}
