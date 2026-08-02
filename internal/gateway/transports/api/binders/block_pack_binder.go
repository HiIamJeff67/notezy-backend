package binders

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	blockpacksdto "github.com/HiIamJeff67/notezy-backend/contracts/gateway/v1/api/block-packs"
	exceptions "github.com/HiIamJeff67/notezy-backend/internal/exceptions"
	responsewriter "github.com/HiIamJeff67/notezy-backend/internal/shared/responsewriter"
	controllers "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/api/controllers"
)

type BlockPackBinderInterface interface {
	BindGetMyBlockPackById(controllerFunc controllers.Func[*blockpacksdto.GetMyBlockPackByIdRequestDto]) gin.HandlerFunc
	BindGetMyBlockPackAndItsParentById(controllerFunc controllers.Func[*blockpacksdto.GetMyBlockPackAndItsParentByIdRequestDto]) gin.HandlerFunc
	BindGetMyBlockPacksByParentSubShelfId(controllerFunc controllers.Func[*blockpacksdto.GetMyBlockPacksByParentSubShelfIdRequestDto]) gin.HandlerFunc
	BindGetAllMyBlockPacksByRootShelfId(controllerFunc controllers.Func[*blockpacksdto.GetAllMyBlockPacksByRootShelfIdRequestDto]) gin.HandlerFunc
	BindCreateBlockPack(controllerFunc controllers.Func[*blockpacksdto.CreateBlockPackRequestDto]) gin.HandlerFunc
	BindCreateBlockPacks(controllerFunc controllers.Func[*blockpacksdto.CreateBlockPacksRequestDto]) gin.HandlerFunc
	BindUpdateMyBlockPackById(controllerFunc controllers.Func[*blockpacksdto.UpdateMyBlockPackByIdRequestDto]) gin.HandlerFunc
	BindUpdateMyBlockPacksByIds(controllerFunc controllers.Func[*blockpacksdto.UpdateMyBlockPacksByIdsRequestDto]) gin.HandlerFunc
	BindMoveMyBlockPackByParentSubShelfId(controllerFunc controllers.Func[*blockpacksdto.MoveMyBlockPackByParentSubShelfIdRequestDto]) gin.HandlerFunc
	BindMoveMyBlockPacksByParentSubShelfId(controllerFunc controllers.Func[*blockpacksdto.MoveMyBlockPacksByParentSubShelfIdRequestDto]) gin.HandlerFunc
	BindMoveMyBlockPacksByParentSubShelfIds(controllerFunc controllers.Func[*blockpacksdto.MoveMyBlockPacksByParentSubShelfIdsRequestDto]) gin.HandlerFunc
	BindRestoreMyBlockPackById(controllerFunc controllers.Func[*blockpacksdto.RestoreMyBlockPackByIdRequestDto]) gin.HandlerFunc
	BindRestoreMyBlockPacksByIds(controllerFunc controllers.Func[*blockpacksdto.RestoreMyBlockPacksByIdsRequestDto]) gin.HandlerFunc
	BindDeleteMyBlockPackById(controllerFunc controllers.Func[*blockpacksdto.DeleteMyBlockPackByIdRequestDto]) gin.HandlerFunc
	BindDeleteMyBlockPacksByIds(controllerFunc controllers.Func[*blockpacksdto.DeleteMyBlockPacksByIdsRequestDto]) gin.HandlerFunc
}

type BlockPackBinder struct{}

func NewBlockPackBinder() BlockPackBinderInterface {
	return &BlockPackBinder{}
}

func (b *BlockPackBinder) BindGetMyBlockPackById(controllerFunc controllers.Func[*blockpacksdto.GetMyBlockPackByIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestDto blockpacksdto.GetMyBlockPackByIdRequestDto

		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		isDeletedString := ctx.Query("isDeleted")
		if isDeletedString != "" {
			value, err := strconv.ParseBool(isDeletedString)
			if err != nil {
				responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("BlockPack").WithOrigin(err), ctx)
				return
			}
			requestDto.Param.IsDeleted = &value
		}

		value, err := uuid.Parse(ctx.Param("blockPackId"))
		if err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("BlockPack").WithOrigin(err), ctx)
			return
		}
		requestDto.Param.BlockPackId = value

		controllerFunc(ctx, &requestDto)
	}
}

func (b *BlockPackBinder) BindGetMyBlockPackAndItsParentById(controllerFunc controllers.Func[*blockpacksdto.GetMyBlockPackAndItsParentByIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestDto blockpacksdto.GetMyBlockPackAndItsParentByIdRequestDto

		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		isDeletedString := ctx.Query("isDeleted")
		if isDeletedString != "" {
			value, err := strconv.ParseBool(isDeletedString)
			if err != nil {
				responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("BlockPack").WithOrigin(err), ctx)
				return
			}
			requestDto.Param.IsDeleted = &value
		}

		value, err := uuid.Parse(ctx.Param("blockPackId"))
		if err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("BlockPack").WithOrigin(err), ctx)
			return
		}
		requestDto.Param.BlockPackId = value

		controllerFunc(ctx, &requestDto)
	}
}

func (b *BlockPackBinder) BindGetMyBlockPacksByParentSubShelfId(controllerFunc controllers.Func[*blockpacksdto.GetMyBlockPacksByParentSubShelfIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestDto blockpacksdto.GetMyBlockPacksByParentSubShelfIdRequestDto

		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		areDeletedString := ctx.Query("areDeleted")
		if areDeletedString != "" {
			value, err := strconv.ParseBool(areDeletedString)
			if err != nil {
				responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("BlockPack").WithOrigin(err), ctx)
				return
			}
			requestDto.Param.AreDeleted = &value
		}

		value, err := uuid.Parse(ctx.Param("parentSubShelfId"))
		if err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("BlockPack").WithOrigin(err), ctx)
			return
		}
		requestDto.Param.ParentSubShelfId = value

		controllerFunc(ctx, &requestDto)
	}
}

func (b *BlockPackBinder) BindGetAllMyBlockPacksByRootShelfId(controllerFunc controllers.Func[*blockpacksdto.GetAllMyBlockPacksByRootShelfIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestDto blockpacksdto.GetAllMyBlockPacksByRootShelfIdRequestDto

		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		areDeletedString := ctx.Query("areDeleted")
		if areDeletedString != "" {
			value, err := strconv.ParseBool(areDeletedString)
			if err != nil {
				responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("BlockPack").WithOrigin(err), ctx)
				return
			}
			requestDto.Param.AreDeleted = &value
		}

		value, err := uuid.Parse(ctx.Param("rootShelfId"))
		if err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("BlockPack").WithOrigin(err), ctx)
			return
		}
		requestDto.Param.RootShelfId = value

		controllerFunc(ctx, &requestDto)
	}
}

func (b *BlockPackBinder) BindCreateBlockPack(controllerFunc controllers.Func[*blockpacksdto.CreateBlockPackRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestDto blockpacksdto.CreateBlockPackRequestDto

		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			exception := exceptions.InvalidDto("BlockPack").WithOrigin(err)
			responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
			return
		}

		value, err := uuid.Parse(ctx.Param("parentSubShelfId"))
		if err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("BlockPack").WithOrigin(err), ctx)
			return
		}
		requestDto.Body.ParentSubShelfId = value

		controllerFunc(ctx, &requestDto)
	}
}

func (b *BlockPackBinder) BindCreateBlockPacks(controllerFunc controllers.Func[*blockpacksdto.CreateBlockPacksRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestDto blockpacksdto.CreateBlockPacksRequestDto

		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			exception := exceptions.InvalidDto("BlockPack").WithOrigin(err)
			responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
			return
		}

		controllerFunc(ctx, &requestDto)
	}
}

func (b *BlockPackBinder) BindUpdateMyBlockPackById(controllerFunc controllers.Func[*blockpacksdto.UpdateMyBlockPackByIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestDto blockpacksdto.UpdateMyBlockPackByIdRequestDto

		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			exception := exceptions.InvalidDto("BlockPack").WithOrigin(err)
			responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
			return
		}

		value, err := uuid.Parse(ctx.Param("blockPackId"))
		if err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("BlockPack").WithOrigin(err), ctx)
			return
		}
		requestDto.Param.BlockPackId = value

		controllerFunc(ctx, &requestDto)
	}
}

func (b *BlockPackBinder) BindUpdateMyBlockPacksByIds(controllerFunc controllers.Func[*blockpacksdto.UpdateMyBlockPacksByIdsRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestDto blockpacksdto.UpdateMyBlockPacksByIdsRequestDto

		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			exception := exceptions.InvalidDto("BlockPack").WithOrigin(err)
			responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
			return
		}

		controllerFunc(ctx, &requestDto)
	}
}

func (b *BlockPackBinder) BindMoveMyBlockPackByParentSubShelfId(controllerFunc controllers.Func[*blockpacksdto.MoveMyBlockPackByParentSubShelfIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestDto blockpacksdto.MoveMyBlockPackByParentSubShelfIdRequestDto

		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			exception := exceptions.InvalidDto("BlockPack").WithOrigin(err)
			responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
			return
		}

		value, err := uuid.Parse(ctx.Param("blockPackId"))
		if err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("BlockPack").WithOrigin(err), ctx)
			return
		}
		requestDto.Body.BlockPackId = value

		controllerFunc(ctx, &requestDto)
	}
}

func (b *BlockPackBinder) BindMoveMyBlockPacksByParentSubShelfId(controllerFunc controllers.Func[*blockpacksdto.MoveMyBlockPacksByParentSubShelfIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestDto blockpacksdto.MoveMyBlockPacksByParentSubShelfIdRequestDto

		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			exception := exceptions.InvalidDto("BlockPack").WithOrigin(err)
			responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
			return
		}

		controllerFunc(ctx, &requestDto)
	}
}

func (b *BlockPackBinder) BindMoveMyBlockPacksByParentSubShelfIds(controllerFunc controllers.Func[*blockpacksdto.MoveMyBlockPacksByParentSubShelfIdsRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestDto blockpacksdto.MoveMyBlockPacksByParentSubShelfIdsRequestDto

		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			exception := exceptions.InvalidDto("BlockPack").WithOrigin(err)
			responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
			return
		}

		controllerFunc(ctx, &requestDto)
	}
}

func (b *BlockPackBinder) BindRestoreMyBlockPackById(controllerFunc controllers.Func[*blockpacksdto.RestoreMyBlockPackByIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestDto blockpacksdto.RestoreMyBlockPackByIdRequestDto

		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		value, err := uuid.Parse(ctx.Param("blockPackId"))
		if err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("BlockPack").WithOrigin(err), ctx)
			return
		}
		requestDto.Param.BlockPackId = value

		controllerFunc(ctx, &requestDto)
	}
}

func (b *BlockPackBinder) BindRestoreMyBlockPacksByIds(controllerFunc controllers.Func[*blockpacksdto.RestoreMyBlockPacksByIdsRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestDto blockpacksdto.RestoreMyBlockPacksByIdsRequestDto

		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			exception := exceptions.InvalidDto("BlockPack").WithOrigin(err)
			responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
			return
		}

		controllerFunc(ctx, &requestDto)
	}
}

func (b *BlockPackBinder) BindDeleteMyBlockPackById(controllerFunc controllers.Func[*blockpacksdto.DeleteMyBlockPackByIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestDto blockpacksdto.DeleteMyBlockPackByIdRequestDto

		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		value, err := uuid.Parse(ctx.Param("blockPackId"))
		if err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("BlockPack").WithOrigin(err), ctx)
			return
		}
		requestDto.Param.BlockPackId = value

		controllerFunc(ctx, &requestDto)
	}
}

func (b *BlockPackBinder) BindDeleteMyBlockPacksByIds(controllerFunc controllers.Func[*blockpacksdto.DeleteMyBlockPacksByIdsRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestDto blockpacksdto.DeleteMyBlockPacksByIdsRequestDto

		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			exception := exceptions.InvalidDto("BlockPack").WithOrigin(err)
			responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
			return
		}

		controllerFunc(ctx, &requestDto)
	}
}
