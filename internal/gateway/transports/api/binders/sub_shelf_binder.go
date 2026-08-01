package binders

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	subshelvesdto "github.com/HiIamJeff67/notezy-backend/contracts/api/v1/sub-shelves"
	exceptions "github.com/HiIamJeff67/notezy-backend/internal/exceptions"
	responsewriter "github.com/HiIamJeff67/notezy-backend/internal/gateway/responsewriter"
	apitransport "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/api"
)

type SubShelfBinderInterface interface {
	BindGetMySubShelfById(controllerFunc apitransport.ControllerFunc[*subshelvesdto.GetMySubShelfByIdRequestDto]) gin.HandlerFunc
	BindGetMySubShelvesByPrevSubShelfId(controllerFunc apitransport.ControllerFunc[*subshelvesdto.GetMySubShelvesByPrevSubShelfIdRequestDto]) gin.HandlerFunc
	BindGetAllMySubShelvesByRootShelfId(controllerFunc apitransport.ControllerFunc[*subshelvesdto.GetAllMySubShelvesByRootShelfIdRequestDto]) gin.HandlerFunc
	BindGetMySubShelvesAndItemsByPrevSubShelfId(controllerFunc apitransport.ControllerFunc[*subshelvesdto.GetMySubShelvesAndItemsByPrevSubShelfIdRequestDto]) gin.HandlerFunc
	BindCreateSubShelfByRootShelfId(controllerFunc apitransport.ControllerFunc[*subshelvesdto.CreateSubShelfByRootShelfIdRequestDto]) gin.HandlerFunc
	BindCreateSubShelvesByRootShelfIds(controllerFunc apitransport.ControllerFunc[*subshelvesdto.CreateSubShelvesByRootShelfIdsRequestDto]) gin.HandlerFunc
	BindUpdateMySubShelfById(controllerFunc apitransport.ControllerFunc[*subshelvesdto.UpdateMySubShelfByIdRequestDto]) gin.HandlerFunc
	BindUpdateMySubShelvesByIds(controllerFunc apitransport.ControllerFunc[*subshelvesdto.UpdateMySubShelvesByIdsRequestDto]) gin.HandlerFunc
	BindMoveMySubShelfByRootShelfId(controllerFunc apitransport.ControllerFunc[*subshelvesdto.MoveMySubShelfByRootShelfIdRequestDto]) gin.HandlerFunc
	BindMoveMySubShelvesByRootShelfId(controllerFunc apitransport.ControllerFunc[*subshelvesdto.MoveMySubShelvesByRootShelfIdRequestDto]) gin.HandlerFunc
	BindMoveMySubShelvesByRootShelfIds(controllerFunc apitransport.ControllerFunc[*subshelvesdto.MoveMySubShelvesByRootShelfIdsRequestDto]) gin.HandlerFunc
	BindRestoreMySubShelfById(controllerFunc apitransport.ControllerFunc[*subshelvesdto.RestoreMySubShelfByIdRequestDto]) gin.HandlerFunc
	BindRestoreMySubShelvesByIds(controllerFunc apitransport.ControllerFunc[*subshelvesdto.RestoreMySubShelvesByIdsRequestDto]) gin.HandlerFunc
	BindDeleteMySubShelfById(controllerFunc apitransport.ControllerFunc[*subshelvesdto.DeleteMySubShelfByIdRequestDto]) gin.HandlerFunc
	BindDeleteMySubShelvesByIds(controllerFunc apitransport.ControllerFunc[*subshelvesdto.DeleteMySubShelvesByIdsRequestDto]) gin.HandlerFunc
}

type SubShelfBinder struct{}

func NewSubShelfBinder() SubShelfBinderInterface {
	return &SubShelfBinder{}
}

func (b *SubShelfBinder) BindGetMySubShelfById(controllerFunc apitransport.ControllerFunc[*subshelvesdto.GetMySubShelfByIdRequestDto]) gin.HandlerFunc {
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

func (b *SubShelfBinder) BindGetMySubShelvesByPrevSubShelfId(controllerFunc apitransport.ControllerFunc[*subshelvesdto.GetMySubShelvesByPrevSubShelfIdRequestDto]) gin.HandlerFunc {
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

func (b *SubShelfBinder) BindGetAllMySubShelvesByRootShelfId(controllerFunc apitransport.ControllerFunc[*subshelvesdto.GetAllMySubShelvesByRootShelfIdRequestDto]) gin.HandlerFunc {
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

func (b *SubShelfBinder) BindGetMySubShelvesAndItemsByPrevSubShelfId(controllerFunc apitransport.ControllerFunc[*subshelvesdto.GetMySubShelvesAndItemsByPrevSubShelfIdRequestDto]) gin.HandlerFunc {
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

func (b *SubShelfBinder) BindCreateSubShelfByRootShelfId(controllerFunc apitransport.ControllerFunc[*subshelvesdto.CreateSubShelfByRootShelfIdRequestDto]) gin.HandlerFunc {
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

func (b *SubShelfBinder) BindCreateSubShelvesByRootShelfIds(controllerFunc apitransport.ControllerFunc[*subshelvesdto.CreateSubShelvesByRootShelfIdsRequestDto]) gin.HandlerFunc {
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

func (b *SubShelfBinder) BindUpdateMySubShelfById(controllerFunc apitransport.ControllerFunc[*subshelvesdto.UpdateMySubShelfByIdRequestDto]) gin.HandlerFunc {
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

func (b *SubShelfBinder) BindUpdateMySubShelvesByIds(controllerFunc apitransport.ControllerFunc[*subshelvesdto.UpdateMySubShelvesByIdsRequestDto]) gin.HandlerFunc {
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

func (b *SubShelfBinder) BindMoveMySubShelfByRootShelfId(controllerFunc apitransport.ControllerFunc[*subshelvesdto.MoveMySubShelfByRootShelfIdRequestDto]) gin.HandlerFunc {
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

func (b *SubShelfBinder) BindMoveMySubShelvesByRootShelfId(controllerFunc apitransport.ControllerFunc[*subshelvesdto.MoveMySubShelvesByRootShelfIdRequestDto]) gin.HandlerFunc {
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

func (b *SubShelfBinder) BindMoveMySubShelvesByRootShelfIds(controllerFunc apitransport.ControllerFunc[*subshelvesdto.MoveMySubShelvesByRootShelfIdsRequestDto]) gin.HandlerFunc {
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

func (b *SubShelfBinder) BindRestoreMySubShelfById(controllerFunc apitransport.ControllerFunc[*subshelvesdto.RestoreMySubShelfByIdRequestDto]) gin.HandlerFunc {
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

func (b *SubShelfBinder) BindRestoreMySubShelvesByIds(controllerFunc apitransport.ControllerFunc[*subshelvesdto.RestoreMySubShelvesByIdsRequestDto]) gin.HandlerFunc {
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

func (b *SubShelfBinder) BindDeleteMySubShelfById(controllerFunc apitransport.ControllerFunc[*subshelvesdto.DeleteMySubShelfByIdRequestDto]) gin.HandlerFunc {
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

func (b *SubShelfBinder) BindDeleteMySubShelvesByIds(controllerFunc apitransport.ControllerFunc[*subshelvesdto.DeleteMySubShelvesByIdsRequestDto]) gin.HandlerFunc {
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
