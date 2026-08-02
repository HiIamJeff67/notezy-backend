package binders

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	subshelvesdto "github.com/HiIamJeff67/notezy-backend/contracts/gateway/v1/api/sub-shelves"
	exceptions "github.com/HiIamJeff67/notezy-backend/internal/exceptions"
	responsewriter "github.com/HiIamJeff67/notezy-backend/internal/shared/responsewriter"
	controllers "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/api/controllers"
)

type SubShelfBinderInterface interface {
	BindGetMySubShelfById(controllerFunc controllers.Func[*subshelvesdto.GetMySubShelfByIdRequestDto]) gin.HandlerFunc
	BindGetMySubShelvesByPrevSubShelfId(controllerFunc controllers.Func[*subshelvesdto.GetMySubShelvesByPrevSubShelfIdRequestDto]) gin.HandlerFunc
	BindGetAllMySubShelvesByRootShelfId(controllerFunc controllers.Func[*subshelvesdto.GetAllMySubShelvesByRootShelfIdRequestDto]) gin.HandlerFunc
	BindGetMySubShelvesAndItemsByPrevSubShelfId(controllerFunc controllers.Func[*subshelvesdto.GetMySubShelvesAndItemsByPrevSubShelfIdRequestDto]) gin.HandlerFunc
	BindCreateSubShelfByRootShelfId(controllerFunc controllers.Func[*subshelvesdto.CreateSubShelfByRootShelfIdRequestDto]) gin.HandlerFunc
	BindCreateSubShelvesByRootShelfIds(controllerFunc controllers.Func[*subshelvesdto.CreateSubShelvesByRootShelfIdsRequestDto]) gin.HandlerFunc
	BindUpdateMySubShelfById(controllerFunc controllers.Func[*subshelvesdto.UpdateMySubShelfByIdRequestDto]) gin.HandlerFunc
	BindUpdateMySubShelvesByIds(controllerFunc controllers.Func[*subshelvesdto.UpdateMySubShelvesByIdsRequestDto]) gin.HandlerFunc
	BindMoveMySubShelfByRootShelfId(controllerFunc controllers.Func[*subshelvesdto.MoveMySubShelfByRootShelfIdRequestDto]) gin.HandlerFunc
	BindMoveMySubShelvesByRootShelfId(controllerFunc controllers.Func[*subshelvesdto.MoveMySubShelvesByRootShelfIdRequestDto]) gin.HandlerFunc
	BindMoveMySubShelvesByRootShelfIds(controllerFunc controllers.Func[*subshelvesdto.MoveMySubShelvesByRootShelfIdsRequestDto]) gin.HandlerFunc
	BindRestoreMySubShelfById(controllerFunc controllers.Func[*subshelvesdto.RestoreMySubShelfByIdRequestDto]) gin.HandlerFunc
	BindRestoreMySubShelvesByIds(controllerFunc controllers.Func[*subshelvesdto.RestoreMySubShelvesByIdsRequestDto]) gin.HandlerFunc
	BindDeleteMySubShelfById(controllerFunc controllers.Func[*subshelvesdto.DeleteMySubShelfByIdRequestDto]) gin.HandlerFunc
	BindDeleteMySubShelvesByIds(controllerFunc controllers.Func[*subshelvesdto.DeleteMySubShelvesByIdsRequestDto]) gin.HandlerFunc
}

type SubShelfBinder struct{}

func NewSubShelfBinder() SubShelfBinderInterface {
	return &SubShelfBinder{}
}

func (b *SubShelfBinder) BindGetMySubShelfById(controllerFunc controllers.Func[*subshelvesdto.GetMySubShelfByIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestDto subshelvesdto.GetMySubShelfByIdRequestDto

		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		isDeletedString := ctx.Query("isDeleted")
		if isDeletedString != "" {
			isDeleted, err := strconv.ParseBool(isDeletedString)
			if err != nil {
				responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Shelf").WithOrigin(err), ctx)
				return
			}
			requestDto.Param.IsDeleted = &isDeleted
		}

		subShelfId, err := uuid.Parse(ctx.Param("subShelfId"))
		if err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Shelf").WithOrigin(err), ctx)
			return
		}
		requestDto.Param.SubShelfId = subShelfId

		controllerFunc(ctx, &requestDto)
	}
}

func (b *SubShelfBinder) BindGetMySubShelvesByPrevSubShelfId(controllerFunc controllers.Func[*subshelvesdto.GetMySubShelvesByPrevSubShelfIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestDto subshelvesdto.GetMySubShelvesByPrevSubShelfIdRequestDto

		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		areDeletedString := ctx.Query("areDeleted")
		if areDeletedString != "" {
			areDeleted, err := strconv.ParseBool(areDeletedString)
			if err != nil {
				responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Shelf").WithOrigin(err), ctx)
				return
			}
			requestDto.Param.AreDeleted = &areDeleted
		}

		prevSubShelfId, err := uuid.Parse(ctx.Param("prevSubShelfId"))
		if err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Shelf").WithOrigin(err), ctx)
			return
		}
		requestDto.Param.PrevSubShelfId = prevSubShelfId

		controllerFunc(ctx, &requestDto)
	}
}

func (b *SubShelfBinder) BindGetAllMySubShelvesByRootShelfId(controllerFunc controllers.Func[*subshelvesdto.GetAllMySubShelvesByRootShelfIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestDto subshelvesdto.GetAllMySubShelvesByRootShelfIdRequestDto

		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		areDeletedString := ctx.Query("areDeleted")
		if areDeletedString != "" {
			areDeleted, err := strconv.ParseBool(areDeletedString)
			if err != nil {
				responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Shelf").WithOrigin(err), ctx)
				return
			}
			requestDto.Param.AreDeleted = &areDeleted
		}

		rootShelfId, err := uuid.Parse(ctx.Param("rootShelfId"))
		if err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Shelf").WithOrigin(err), ctx)
			return
		}
		requestDto.Param.RootShelfId = rootShelfId

		controllerFunc(ctx, &requestDto)
	}
}

func (b *SubShelfBinder) BindGetMySubShelvesAndItemsByPrevSubShelfId(controllerFunc controllers.Func[*subshelvesdto.GetMySubShelvesAndItemsByPrevSubShelfIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestDto subshelvesdto.GetMySubShelvesAndItemsByPrevSubShelfIdRequestDto

		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		areDeletedString := ctx.Query("areDeleted")
		if areDeletedString != "" {
			areDeleted, err := strconv.ParseBool(areDeletedString)
			if err != nil {
				responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Shelf").WithOrigin(err), ctx)
				return
			}
			requestDto.Param.AreDeleted = &areDeleted
		}

		prevSubShelfId, err := uuid.Parse(ctx.Param("prevSubShelfId"))
		if err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Shelf").WithOrigin(err), ctx)
			return
		}
		requestDto.Param.PrevSubShelfId = prevSubShelfId

		controllerFunc(ctx, &requestDto)
	}
}

func (b *SubShelfBinder) BindCreateSubShelfByRootShelfId(controllerFunc controllers.Func[*subshelvesdto.CreateSubShelfByRootShelfIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestDto subshelvesdto.CreateSubShelfByRootShelfIdRequestDto

		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			exception := exceptions.InvalidDto("Shelf").WithOrigin(err)
			responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
			return
		}

		rootShelfId, err := uuid.Parse(ctx.Param("rootShelfId"))
		if err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Shelf").WithOrigin(err), ctx)
			return
		}
		requestDto.Body.RootShelfId = rootShelfId

		controllerFunc(ctx, &requestDto)
	}
}

func (b *SubShelfBinder) BindCreateSubShelvesByRootShelfIds(controllerFunc controllers.Func[*subshelvesdto.CreateSubShelvesByRootShelfIdsRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestDto subshelvesdto.CreateSubShelvesByRootShelfIdsRequestDto

		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			exception := exceptions.InvalidDto("Shelf").WithOrigin(err)
			responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
			return
		}

		controllerFunc(ctx, &requestDto)
	}
}

func (b *SubShelfBinder) BindUpdateMySubShelfById(controllerFunc controllers.Func[*subshelvesdto.UpdateMySubShelfByIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestDto subshelvesdto.UpdateMySubShelfByIdRequestDto

		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			exception := exceptions.InvalidDto("Shelf").WithOrigin(err)
			responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
			return
		}

		subShelfId, err := uuid.Parse(ctx.Param("subShelfId"))
		if err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Shelf").WithOrigin(err), ctx)
			return
		}
		requestDto.Param.SubShelfId = subShelfId

		controllerFunc(ctx, &requestDto)
	}
}

func (b *SubShelfBinder) BindUpdateMySubShelvesByIds(controllerFunc controllers.Func[*subshelvesdto.UpdateMySubShelvesByIdsRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestDto subshelvesdto.UpdateMySubShelvesByIdsRequestDto

		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			exception := exceptions.InvalidDto("Shelf").WithOrigin(err)
			responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
			return
		}

		controllerFunc(ctx, &requestDto)
	}
}

func (b *SubShelfBinder) BindMoveMySubShelfByRootShelfId(controllerFunc controllers.Func[*subshelvesdto.MoveMySubShelfByRootShelfIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestDto subshelvesdto.MoveMySubShelfByRootShelfIdRequestDto

		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		subShelfId, err := uuid.Parse(ctx.Param("subShelfId"))
		if err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Shelf").WithOrigin(err), ctx)
			return
		}
		requestDto.Body.SourceSubShelfId = subShelfId

		controllerFunc(ctx, &requestDto)
	}
}

func (b *SubShelfBinder) BindMoveMySubShelvesByRootShelfId(controllerFunc controllers.Func[*subshelvesdto.MoveMySubShelvesByRootShelfIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestDto subshelvesdto.MoveMySubShelvesByRootShelfIdRequestDto

		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			exception := exceptions.InvalidDto("Shelf").WithOrigin(err)
			responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
			return
		}

		controllerFunc(ctx, &requestDto)
	}
}

func (b *SubShelfBinder) BindMoveMySubShelvesByRootShelfIds(controllerFunc controllers.Func[*subshelvesdto.MoveMySubShelvesByRootShelfIdsRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestDto subshelvesdto.MoveMySubShelvesByRootShelfIdsRequestDto

		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			exception := exceptions.InvalidDto("Shelf").WithOrigin(err)
			responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
			return
		}

		controllerFunc(ctx, &requestDto)
	}
}

func (b *SubShelfBinder) BindRestoreMySubShelfById(controllerFunc controllers.Func[*subshelvesdto.RestoreMySubShelfByIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestDto subshelvesdto.RestoreMySubShelfByIdRequestDto

		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		subShelfId, err := uuid.Parse(ctx.Param("subShelfId"))
		if err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Shelf").WithOrigin(err), ctx)
			return
		}
		requestDto.Param.SubShelfId = subShelfId

		controllerFunc(ctx, &requestDto)
	}
}

func (b *SubShelfBinder) BindRestoreMySubShelvesByIds(controllerFunc controllers.Func[*subshelvesdto.RestoreMySubShelvesByIdsRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestDto subshelvesdto.RestoreMySubShelvesByIdsRequestDto

		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			exception := exceptions.InvalidDto("Shelf").WithOrigin(err)
			responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
			return
		}

		controllerFunc(ctx, &requestDto)
	}
}

func (b *SubShelfBinder) BindDeleteMySubShelfById(controllerFunc controllers.Func[*subshelvesdto.DeleteMySubShelfByIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestDto subshelvesdto.DeleteMySubShelfByIdRequestDto

		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		subShelfId, err := uuid.Parse(ctx.Param("subShelfId"))
		if err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Shelf").WithOrigin(err), ctx)
			return
		}
		requestDto.Param.SubShelfId = subShelfId

		controllerFunc(ctx, &requestDto)
	}
}

func (b *SubShelfBinder) BindDeleteMySubShelvesByIds(controllerFunc controllers.Func[*subshelvesdto.DeleteMySubShelvesByIdsRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestDto subshelvesdto.DeleteMySubShelvesByIdsRequestDto

		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			exception := exceptions.InvalidDto("Shelf").WithOrigin(err)
			responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
			return
		}

		controllerFunc(ctx, &requestDto)
	}
}
