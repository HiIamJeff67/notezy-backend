package binders

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	exceptions "github.com/HiIamJeff67/notegic-backend/contracts/types/exceptions"

	exceptionwriter "github.com/HiIamJeff67/notegic-backend/shared/util/exceptionwriter"

	apicontract "github.com/HiIamJeff67/notegic-backend/contracts/core/v1/api/sub-shelves"

	controllers "github.com/HiIamJeff67/notegic-backend/internal/apigateway/transports/api/controllers"
)

type SubShelfBinderInterface interface {
	BindGetMySubShelfById(controllerFunc controllers.Func[*apicontract.GetMySubShelfByIdRequestDto]) gin.HandlerFunc
	BindGetMySubShelvesByPrevSubShelfId(controllerFunc controllers.Func[*apicontract.GetMySubShelvesByPrevSubShelfIdRequestDto]) gin.HandlerFunc
	BindGetAllMySubShelvesByRootShelfId(controllerFunc controllers.Func[*apicontract.GetAllMySubShelvesByRootShelfIdRequestDto]) gin.HandlerFunc
	BindGetMySubShelvesAndItemsByPrevSubShelfId(controllerFunc controllers.Func[*apicontract.GetMySubShelvesAndItemsByPrevSubShelfIdRequestDto]) gin.HandlerFunc
	BindCreateSubShelfByRootShelfId(controllerFunc controllers.Func[*apicontract.CreateSubShelfByRootShelfIdRequestDto]) gin.HandlerFunc
	BindCreateSubShelvesByRootShelfIds(controllerFunc controllers.Func[*apicontract.CreateSubShelvesByRootShelfIdsRequestDto]) gin.HandlerFunc
	BindUpdateMySubShelfById(controllerFunc controllers.Func[*apicontract.UpdateMySubShelfByIdRequestDto]) gin.HandlerFunc
	BindUpdateMySubShelvesByIds(controllerFunc controllers.Func[*apicontract.UpdateMySubShelvesByIdsRequestDto]) gin.HandlerFunc
	BindMoveMySubShelfByRootShelfId(controllerFunc controllers.Func[*apicontract.MoveMySubShelfByRootShelfIdRequestDto]) gin.HandlerFunc
	BindMoveMySubShelvesByRootShelfId(controllerFunc controllers.Func[*apicontract.MoveMySubShelvesByRootShelfIdRequestDto]) gin.HandlerFunc
	BindMoveMySubShelvesByRootShelfIds(controllerFunc controllers.Func[*apicontract.MoveMySubShelvesByRootShelfIdsRequestDto]) gin.HandlerFunc
	BindRestoreMySubShelfById(controllerFunc controllers.Func[*apicontract.RestoreMySubShelfByIdRequestDto]) gin.HandlerFunc
	BindRestoreMySubShelvesByIds(controllerFunc controllers.Func[*apicontract.RestoreMySubShelvesByIdsRequestDto]) gin.HandlerFunc
	BindDeleteMySubShelfById(controllerFunc controllers.Func[*apicontract.DeleteMySubShelfByIdRequestDto]) gin.HandlerFunc
	BindDeleteMySubShelvesByIds(controllerFunc controllers.Func[*apicontract.DeleteMySubShelvesByIdsRequestDto]) gin.HandlerFunc
}

type SubShelfBinder struct{}

func NewSubShelfBinder() SubShelfBinderInterface {
	return &SubShelfBinder{}
}

func (b *SubShelfBinder) BindGetMySubShelfById(controllerFunc controllers.Func[*apicontract.GetMySubShelfByIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestDto apicontract.GetMySubShelfByIdRequestDto

		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		isDeletedString := ctx.Query("isDeleted")
		if isDeletedString != "" {
			isDeleted, err := strconv.ParseBool(isDeletedString)
			if err != nil {
				exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Shelf").WithOrigin(err), ctx)
				return
			}
			requestDto.Param.IsDeleted = &isDeleted
		}

		subShelfId, err := uuid.Parse(ctx.Param("sub-shelf-id"))
		if err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Shelf").WithOrigin(err), ctx)
			return
		}
		requestDto.Param.SubShelfId = subShelfId

		controllerFunc(ctx, &requestDto)
	}
}

func (b *SubShelfBinder) BindGetMySubShelvesByPrevSubShelfId(controllerFunc controllers.Func[*apicontract.GetMySubShelvesByPrevSubShelfIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestDto apicontract.GetMySubShelvesByPrevSubShelfIdRequestDto

		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		areDeletedString := ctx.Query("areDeleted")
		if areDeletedString != "" {
			areDeleted, err := strconv.ParseBool(areDeletedString)
			if err != nil {
				exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Shelf").WithOrigin(err), ctx)
				return
			}
			requestDto.Param.AreDeleted = &areDeleted
		}

		prevSubShelfId, err := uuid.Parse(ctx.Param("prev-sub-shelf-id"))
		if err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Shelf").WithOrigin(err), ctx)
			return
		}
		requestDto.Param.PrevSubShelfId = prevSubShelfId

		controllerFunc(ctx, &requestDto)
	}
}

func (b *SubShelfBinder) BindGetAllMySubShelvesByRootShelfId(controllerFunc controllers.Func[*apicontract.GetAllMySubShelvesByRootShelfIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestDto apicontract.GetAllMySubShelvesByRootShelfIdRequestDto

		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		areDeletedString := ctx.Query("areDeleted")
		if areDeletedString != "" {
			areDeleted, err := strconv.ParseBool(areDeletedString)
			if err != nil {
				exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Shelf").WithOrigin(err), ctx)
				return
			}
			requestDto.Param.AreDeleted = &areDeleted
		}

		rootShelfId, err := uuid.Parse(ctx.Param("root-shelf-id"))
		if err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Shelf").WithOrigin(err), ctx)
			return
		}
		requestDto.Param.RootShelfId = rootShelfId

		controllerFunc(ctx, &requestDto)
	}
}

func (b *SubShelfBinder) BindGetMySubShelvesAndItemsByPrevSubShelfId(controllerFunc controllers.Func[*apicontract.GetMySubShelvesAndItemsByPrevSubShelfIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestDto apicontract.GetMySubShelvesAndItemsByPrevSubShelfIdRequestDto

		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		areDeletedString := ctx.Query("areDeleted")
		if areDeletedString != "" {
			areDeleted, err := strconv.ParseBool(areDeletedString)
			if err != nil {
				exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Shelf").WithOrigin(err), ctx)
				return
			}
			requestDto.Param.AreDeleted = &areDeleted
		}

		prevSubShelfId, err := uuid.Parse(ctx.Param("prev-sub-shelf-id"))
		if err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Shelf").WithOrigin(err), ctx)
			return
		}
		requestDto.Param.PrevSubShelfId = prevSubShelfId

		controllerFunc(ctx, &requestDto)
	}
}

func (b *SubShelfBinder) BindCreateSubShelfByRootShelfId(controllerFunc controllers.Func[*apicontract.CreateSubShelfByRootShelfIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestDto apicontract.CreateSubShelfByRootShelfIdRequestDto

		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			exception := exceptions.InvalidDto("Shelf").WithOrigin(err)
			exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
			return
		}

		rootShelfId, err := uuid.Parse(ctx.Param("root-shelf-id"))
		if err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Shelf").WithOrigin(err), ctx)
			return
		}
		requestDto.Body.RootShelfId = rootShelfId

		controllerFunc(ctx, &requestDto)
	}
}

func (b *SubShelfBinder) BindCreateSubShelvesByRootShelfIds(controllerFunc controllers.Func[*apicontract.CreateSubShelvesByRootShelfIdsRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestDto apicontract.CreateSubShelvesByRootShelfIdsRequestDto

		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			exception := exceptions.InvalidDto("Shelf").WithOrigin(err)
			exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
			return
		}

		controllerFunc(ctx, &requestDto)
	}
}

func (b *SubShelfBinder) BindUpdateMySubShelfById(controllerFunc controllers.Func[*apicontract.UpdateMySubShelfByIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestDto apicontract.UpdateMySubShelfByIdRequestDto

		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			exception := exceptions.InvalidDto("Shelf").WithOrigin(err)
			exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
			return
		}

		subShelfId, err := uuid.Parse(ctx.Param("sub-shelf-id"))
		if err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Shelf").WithOrigin(err), ctx)
			return
		}
		requestDto.Param.SubShelfId = subShelfId

		controllerFunc(ctx, &requestDto)
	}
}

func (b *SubShelfBinder) BindUpdateMySubShelvesByIds(controllerFunc controllers.Func[*apicontract.UpdateMySubShelvesByIdsRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestDto apicontract.UpdateMySubShelvesByIdsRequestDto

		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			exception := exceptions.InvalidDto("Shelf").WithOrigin(err)
			exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
			return
		}

		controllerFunc(ctx, &requestDto)
	}
}

func (b *SubShelfBinder) BindMoveMySubShelfByRootShelfId(controllerFunc controllers.Func[*apicontract.MoveMySubShelfByRootShelfIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestDto apicontract.MoveMySubShelfByRootShelfIdRequestDto

		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		subShelfId, err := uuid.Parse(ctx.Param("sub-shelf-id"))
		if err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Shelf").WithOrigin(err), ctx)
			return
		}
		requestDto.Body.SourceSubShelfId = subShelfId

		controllerFunc(ctx, &requestDto)
	}
}

func (b *SubShelfBinder) BindMoveMySubShelvesByRootShelfId(controllerFunc controllers.Func[*apicontract.MoveMySubShelvesByRootShelfIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestDto apicontract.MoveMySubShelvesByRootShelfIdRequestDto

		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			exception := exceptions.InvalidDto("Shelf").WithOrigin(err)
			exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
			return
		}

		controllerFunc(ctx, &requestDto)
	}
}

func (b *SubShelfBinder) BindMoveMySubShelvesByRootShelfIds(controllerFunc controllers.Func[*apicontract.MoveMySubShelvesByRootShelfIdsRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestDto apicontract.MoveMySubShelvesByRootShelfIdsRequestDto

		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			exception := exceptions.InvalidDto("Shelf").WithOrigin(err)
			exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
			return
		}

		controllerFunc(ctx, &requestDto)
	}
}

func (b *SubShelfBinder) BindRestoreMySubShelfById(controllerFunc controllers.Func[*apicontract.RestoreMySubShelfByIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestDto apicontract.RestoreMySubShelfByIdRequestDto

		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		subShelfId, err := uuid.Parse(ctx.Param("sub-shelf-id"))
		if err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Shelf").WithOrigin(err), ctx)
			return
		}
		requestDto.Param.SubShelfId = subShelfId

		controllerFunc(ctx, &requestDto)
	}
}

func (b *SubShelfBinder) BindRestoreMySubShelvesByIds(controllerFunc controllers.Func[*apicontract.RestoreMySubShelvesByIdsRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestDto apicontract.RestoreMySubShelvesByIdsRequestDto

		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			exception := exceptions.InvalidDto("Shelf").WithOrigin(err)
			exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
			return
		}

		controllerFunc(ctx, &requestDto)
	}
}

func (b *SubShelfBinder) BindDeleteMySubShelfById(controllerFunc controllers.Func[*apicontract.DeleteMySubShelfByIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestDto apicontract.DeleteMySubShelfByIdRequestDto

		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		subShelfId, err := uuid.Parse(ctx.Param("sub-shelf-id"))
		if err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Shelf").WithOrigin(err), ctx)
			return
		}
		requestDto.Param.SubShelfId = subShelfId

		controllerFunc(ctx, &requestDto)
	}
}

func (b *SubShelfBinder) BindDeleteMySubShelvesByIds(controllerFunc controllers.Func[*apicontract.DeleteMySubShelvesByIdsRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestDto apicontract.DeleteMySubShelvesByIdsRequestDto

		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			exception := exceptions.InvalidDto("Shelf").WithOrigin(err)
			exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
			return
		}

		controllerFunc(ctx, &requestDto)
	}
}
