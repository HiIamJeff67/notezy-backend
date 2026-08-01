package binders

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	rootshelvesdto "github.com/HiIamJeff67/notezy-backend/contracts/api/v1/root-shelves"
	exceptions "github.com/HiIamJeff67/notezy-backend/internal/exceptions"
	responsewriter "github.com/HiIamJeff67/notezy-backend/internal/gateway/responsewriter"
	apitransport "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/api"
)

type RootShelfBinderInterface interface {
	BindGetMyRootShelfById(controllerFunc apitransport.ControllerFunc[*rootshelvesdto.GetMyRootShelfByIdRequestDto]) gin.HandlerFunc
	BindCreateRootShelf(controllerFunc apitransport.ControllerFunc[*rootshelvesdto.CreateRootShelfRequestDto]) gin.HandlerFunc
	BindCreateRootShelves(controllerFunc apitransport.ControllerFunc[*rootshelvesdto.CreateRootShelvesRequestDto]) gin.HandlerFunc
	BindUpdateMyRootShelfById(controllerFunc apitransport.ControllerFunc[*rootshelvesdto.UpdateMyRootShelfByIdRequestDto]) gin.HandlerFunc
	BindUpdateMyRootShelvesByIds(controllerFunc apitransport.ControllerFunc[*rootshelvesdto.UpdateMyRootShelvesByIdsRequestDto]) gin.HandlerFunc
	BindRestoreMyRootShelfById(controllerFunc apitransport.ControllerFunc[*rootshelvesdto.RestoreMyRootShelfByIdRequestDto]) gin.HandlerFunc
	BindRestoreMyRootShelvesByIds(controllerFunc apitransport.ControllerFunc[*rootshelvesdto.RestoreMyRootShelvesByIdsRequestDto]) gin.HandlerFunc
	BindDeleteMyRootShelfById(controllerFunc apitransport.ControllerFunc[*rootshelvesdto.DeleteMyRootShelfByIdRequestDto]) gin.HandlerFunc
	BindDeleteMyRootShelvesByIds(controllerFunc apitransport.ControllerFunc[*rootshelvesdto.DeleteMyRootShelvesByIdsRequestDto]) gin.HandlerFunc
	BindGetMyRootShelfPermission(controllerFunc apitransport.ControllerFunc[*rootshelvesdto.GetMyRootShelfPermissionRequestDto]) gin.HandlerFunc
	BindCreateMyRootShelfPermission(controllerFunc apitransport.ControllerFunc[*rootshelvesdto.CreateMyRootShelfPermissionRequestDto]) gin.HandlerFunc
	BindUpsertMyRootShelfPermission(controllerFunc apitransport.ControllerFunc[*rootshelvesdto.UpsertMyRootShelfPermissionRequestDto]) gin.HandlerFunc
	BindUpsertMyRootShelfPermissions(controllerFunc apitransport.ControllerFunc[*rootshelvesdto.UpsertMyRootShelfPermissionsRequestDto]) gin.HandlerFunc
	BindUpdateMyRootShelfPermission(controllerFunc apitransport.ControllerFunc[*rootshelvesdto.UpdateMyRootShelfPermissionRequestDto]) gin.HandlerFunc
	BindTransferMyRootShelfOwnership(controllerFunc apitransport.ControllerFunc[*rootshelvesdto.TransferMyRootShelfOwnershipRequestDto]) gin.HandlerFunc
	BindDeleteMyRootShelfPermission(controllerFunc apitransport.ControllerFunc[*rootshelvesdto.DeleteMyRootShelfPermissionRequestDto]) gin.HandlerFunc
	BindDeleteMyRootShelfPermissions(controllerFunc apitransport.ControllerFunc[*rootshelvesdto.DeleteMyRootShelfPermissionsRequestDto]) gin.HandlerFunc
	BindLeaveMyRootShelf(controllerFunc apitransport.ControllerFunc[*rootshelvesdto.LeaveMyRootShelfRequestDto]) gin.HandlerFunc
	BindLeaveMyRootShelves(controllerFunc apitransport.ControllerFunc[*rootshelvesdto.LeaveMyRootShelvesRequestDto]) gin.HandlerFunc
}

type RootShelfBinder struct{}

func NewRootShelfBinder() RootShelfBinderInterface {
	return &RootShelfBinder{}
}

func (b *RootShelfBinder) BindGetMyRootShelfById(controllerFunc apitransport.ControllerFunc[*rootshelvesdto.GetMyRootShelfByIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &rootshelvesdto.GetMyRootShelfByIdRequestDto{}
		request.Header.UserAgent = ctx.GetHeader("User-Agent")

		isDeletedString := ctx.Query("isDeleted")
		if isDeletedString != "" {
			isDeleted, err := strconv.ParseBool(isDeletedString)
			if err != nil {
				responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Shelf").WithOrigin(err), ctx)
				return
			}
			request.Param.IsDeleted = &isDeleted
		}

		rootShelfId, err := uuid.Parse(ctx.Param("rootShelfId"))
		if err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Shelf").WithOrigin(err), ctx)
			return
		}
		request.Param.RootShelfId = rootShelfId

		controllerFunc(ctx, request)
	}
}

func (b *RootShelfBinder) BindCreateRootShelf(controllerFunc apitransport.ControllerFunc[*rootshelvesdto.CreateRootShelfRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &rootshelvesdto.CreateRootShelfRequestDto{}
		request.Header.UserAgent = ctx.GetHeader("User-Agent")
		if err := ctx.ShouldBindJSON(&request.Body); err != nil {
			exception := exceptions.InvalidDto("Shelf").WithOrigin(err)
			responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
			return
		}

		controllerFunc(ctx, request)
	}
}

func (b *RootShelfBinder) BindCreateRootShelves(controllerFunc apitransport.ControllerFunc[*rootshelvesdto.CreateRootShelvesRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &rootshelvesdto.CreateRootShelvesRequestDto{}
		request.Header.UserAgent = ctx.GetHeader("User-Agent")
		if err := ctx.ShouldBindJSON(&request.Body); err != nil {
			exception := exceptions.InvalidDto("Shelf").WithOrigin(err)
			responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
			return
		}

		controllerFunc(ctx, request)
	}
}

func (b *RootShelfBinder) BindUpdateMyRootShelfById(controllerFunc apitransport.ControllerFunc[*rootshelvesdto.UpdateMyRootShelfByIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &rootshelvesdto.UpdateMyRootShelfByIdRequestDto{}
		request.Header.UserAgent = ctx.GetHeader("User-Agent")
		if err := ctx.ShouldBindJSON(&request.Body); err != nil {
			exception := exceptions.InvalidDto("Shelf").WithOrigin(err)
			responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
			return
		}

		rootShelfId, err := uuid.Parse(ctx.Param("rootShelfId"))
		if err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Shelf").WithOrigin(err), ctx)
			return
		}
		request.Param.RootShelfId = rootShelfId

		controllerFunc(ctx, request)
	}
}

func (b *RootShelfBinder) BindUpdateMyRootShelvesByIds(controllerFunc apitransport.ControllerFunc[*rootshelvesdto.UpdateMyRootShelvesByIdsRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto rootshelvesdto.UpdateMyRootShelvesByIdsRequestDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.InvalidDto("Shelf").WithOrigin(err)
			responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *RootShelfBinder) BindRestoreMyRootShelfById(controllerFunc apitransport.ControllerFunc[*rootshelvesdto.RestoreMyRootShelfByIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto rootshelvesdto.RestoreMyRootShelfByIdRequestDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.InvalidDto("Shelf").WithOrigin(err)
			responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
			return
		}

		rootShelfId, err := uuid.Parse(ctx.Param("rootShelfId"))
		if err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Shelf").WithOrigin(err), ctx)
			return
		}
		reqDto.Body.RootShelfId = rootShelfId

		controllerFunc(ctx, &reqDto)
	}
}

func (b *RootShelfBinder) BindRestoreMyRootShelvesByIds(controllerFunc apitransport.ControllerFunc[*rootshelvesdto.RestoreMyRootShelvesByIdsRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto rootshelvesdto.RestoreMyRootShelvesByIdsRequestDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.InvalidDto("Shelf").WithOrigin(err)
			responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *RootShelfBinder) BindDeleteMyRootShelfById(controllerFunc apitransport.ControllerFunc[*rootshelvesdto.DeleteMyRootShelfByIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto rootshelvesdto.DeleteMyRootShelfByIdRequestDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.InvalidDto("Shelf").WithOrigin(err)
			responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
			return
		}

		rootShelfId, err := uuid.Parse(ctx.Param("rootShelfId"))
		if err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Shelf").WithOrigin(err), ctx)
			return
		}
		reqDto.Body.RootShelfId = rootShelfId

		controllerFunc(ctx, &reqDto)
	}
}

func (b *RootShelfBinder) BindDeleteMyRootShelvesByIds(controllerFunc apitransport.ControllerFunc[*rootshelvesdto.DeleteMyRootShelvesByIdsRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto rootshelvesdto.DeleteMyRootShelvesByIdsRequestDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.InvalidDto("Shelf").WithOrigin(err)
			responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *RootShelfBinder) BindGetMyRootShelfPermission(controllerFunc apitransport.ControllerFunc[*rootshelvesdto.GetMyRootShelfPermissionRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &rootshelvesdto.GetMyRootShelfPermissionRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		rootShelfId, err := uuid.Parse(ctx.Param("rootShelfId"))
		if err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Shelf").WithOrigin(err), ctx)
			return
		}
		requestDto.Param.RootShelfId = rootShelfId

		userPublicId, err := uuid.Parse(ctx.Param("userPublicId"))
		if err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Shelf").WithOrigin(err), ctx)
			return
		}
		requestDto.Param.UserPublicId = userPublicId

		controllerFunc(ctx, requestDto)
	}
}

func (b *RootShelfBinder) BindCreateMyRootShelfPermission(controllerFunc apitransport.ControllerFunc[*rootshelvesdto.CreateMyRootShelfPermissionRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &rootshelvesdto.CreateMyRootShelfPermissionRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		rootShelfId, err := uuid.Parse(ctx.Param("rootShelfId"))
		if err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Shelf").WithOrigin(err), ctx)
			return
		}
		requestDto.Param.RootShelfId = rootShelfId

		userPublicId, err := uuid.Parse(ctx.Param("userPublicId"))
		if err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Shelf").WithOrigin(err), ctx)
			return
		}
		requestDto.Param.UserPublicId = userPublicId

		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("Shelf").WithOrigin(err), ctx)
			return
		}

		controllerFunc(ctx, requestDto)
	}
}

func (b *RootShelfBinder) BindUpsertMyRootShelfPermission(controllerFunc apitransport.ControllerFunc[*rootshelvesdto.UpsertMyRootShelfPermissionRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &rootshelvesdto.UpsertMyRootShelfPermissionRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		rootShelfId, err := uuid.Parse(ctx.Param("rootShelfId"))
		if err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Shelf").WithOrigin(err), ctx)
			return
		}
		requestDto.Param.RootShelfId = rootShelfId

		userPublicId, err := uuid.Parse(ctx.Param("userPublicId"))
		if err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Shelf").WithOrigin(err), ctx)
			return
		}
		requestDto.Param.UserPublicId = userPublicId

		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("Shelf").WithOrigin(err), ctx)
			return
		}

		controllerFunc(ctx, requestDto)
	}
}

func (b *RootShelfBinder) BindUpsertMyRootShelfPermissions(controllerFunc apitransport.ControllerFunc[*rootshelvesdto.UpsertMyRootShelfPermissionsRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &rootshelvesdto.UpsertMyRootShelfPermissionsRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		rootShelfId, err := uuid.Parse(ctx.Param("rootShelfId"))
		if err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Shelf").WithOrigin(err), ctx)
			return
		}
		requestDto.Param.RootShelfId = rootShelfId

		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("Shelf").WithOrigin(err), ctx)
			return
		}

		controllerFunc(ctx, requestDto)
	}
}

func (b *RootShelfBinder) BindUpdateMyRootShelfPermission(controllerFunc apitransport.ControllerFunc[*rootshelvesdto.UpdateMyRootShelfPermissionRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &rootshelvesdto.UpdateMyRootShelfPermissionRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		rootShelfId, err := uuid.Parse(ctx.Param("rootShelfId"))
		if err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Shelf").WithOrigin(err), ctx)
			return
		}
		requestDto.Param.RootShelfId = rootShelfId

		userPublicId, err := uuid.Parse(ctx.Param("userPublicId"))
		if err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Shelf").WithOrigin(err), ctx)
			return
		}
		requestDto.Param.UserPublicId = userPublicId

		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("Shelf").WithOrigin(err), ctx)
			return
		}

		controllerFunc(ctx, requestDto)
	}
}

func (b *RootShelfBinder) BindTransferMyRootShelfOwnership(controllerFunc apitransport.ControllerFunc[*rootshelvesdto.TransferMyRootShelfOwnershipRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &rootshelvesdto.TransferMyRootShelfOwnershipRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		rootShelfId, err := uuid.Parse(ctx.Param("rootShelfId"))
		if err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Shelf").WithOrigin(err), ctx)
			return
		}
		requestDto.Param.RootShelfId = rootShelfId

		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("Shelf").WithOrigin(err), ctx)
			return
		}

		controllerFunc(ctx, requestDto)
	}
}

func (b *RootShelfBinder) BindDeleteMyRootShelfPermission(controllerFunc apitransport.ControllerFunc[*rootshelvesdto.DeleteMyRootShelfPermissionRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &rootshelvesdto.DeleteMyRootShelfPermissionRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		rootShelfId, err := uuid.Parse(ctx.Param("rootShelfId"))
		if err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Shelf").WithOrigin(err), ctx)
			return
		}
		requestDto.Param.RootShelfId = rootShelfId

		userPublicId, err := uuid.Parse(ctx.Param("userPublicId"))
		if err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Shelf").WithOrigin(err), ctx)
			return
		}
		requestDto.Param.UserPublicId = userPublicId

		controllerFunc(ctx, requestDto)
	}
}

func (b *RootShelfBinder) BindDeleteMyRootShelfPermissions(controllerFunc apitransport.ControllerFunc[*rootshelvesdto.DeleteMyRootShelfPermissionsRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &rootshelvesdto.DeleteMyRootShelfPermissionsRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		rootShelfId, err := uuid.Parse(ctx.Param("rootShelfId"))
		if err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Shelf").WithOrigin(err), ctx)
			return
		}
		requestDto.Param.RootShelfId = rootShelfId

		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("Shelf").WithOrigin(err), ctx)
			return
		}

		controllerFunc(ctx, requestDto)
	}
}

func (b *RootShelfBinder) BindLeaveMyRootShelf(controllerFunc apitransport.ControllerFunc[*rootshelvesdto.LeaveMyRootShelfRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &rootshelvesdto.LeaveMyRootShelfRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		rootShelfId, err := uuid.Parse(ctx.Param("rootShelfId"))
		if err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Shelf").WithOrigin(err), ctx)
			return
		}
		requestDto.Param.RootShelfId = rootShelfId

		controllerFunc(ctx, requestDto)
	}
}

func (b *RootShelfBinder) BindLeaveMyRootShelves(controllerFunc apitransport.ControllerFunc[*rootshelvesdto.LeaveMyRootShelvesRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &rootshelvesdto.LeaveMyRootShelvesRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("Shelf").WithOrigin(err), ctx)
			return
		}
		controllerFunc(ctx, requestDto)
	}
}
