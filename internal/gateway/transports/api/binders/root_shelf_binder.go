package binders

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	rootshelvesdto "github.com/HiIamJeff67/notezy-backend/contracts/gateway/v1/api/root-shelves"
	exceptions "github.com/HiIamJeff67/notezy-backend/internal/exceptions"
	responsewriter "github.com/HiIamJeff67/notezy-backend/internal/shared/responsewriter"
	controllers "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/api/controllers"
)

type RootShelfBinderInterface interface {
	BindGetMyRootShelfById(controllerFunc controllers.Func[*rootshelvesdto.GetMyRootShelfByIdRequestDto]) gin.HandlerFunc
	BindCreateRootShelf(controllerFunc controllers.Func[*rootshelvesdto.CreateRootShelfRequestDto]) gin.HandlerFunc
	BindCreateRootShelves(controllerFunc controllers.Func[*rootshelvesdto.CreateRootShelvesRequestDto]) gin.HandlerFunc
	BindUpdateMyRootShelfById(controllerFunc controllers.Func[*rootshelvesdto.UpdateMyRootShelfByIdRequestDto]) gin.HandlerFunc
	BindUpdateMyRootShelvesByIds(controllerFunc controllers.Func[*rootshelvesdto.UpdateMyRootShelvesByIdsRequestDto]) gin.HandlerFunc
	BindRestoreMyRootShelfById(controllerFunc controllers.Func[*rootshelvesdto.RestoreMyRootShelfByIdRequestDto]) gin.HandlerFunc
	BindRestoreMyRootShelvesByIds(controllerFunc controllers.Func[*rootshelvesdto.RestoreMyRootShelvesByIdsRequestDto]) gin.HandlerFunc
	BindDeleteMyRootShelfById(controllerFunc controllers.Func[*rootshelvesdto.DeleteMyRootShelfByIdRequestDto]) gin.HandlerFunc
	BindDeleteMyRootShelvesByIds(controllerFunc controllers.Func[*rootshelvesdto.DeleteMyRootShelvesByIdsRequestDto]) gin.HandlerFunc
	BindGetMyRootShelfPermission(controllerFunc controllers.Func[*rootshelvesdto.GetMyRootShelfPermissionRequestDto]) gin.HandlerFunc
	BindCreateMyRootShelfPermission(controllerFunc controllers.Func[*rootshelvesdto.CreateMyRootShelfPermissionRequestDto]) gin.HandlerFunc
	BindUpsertMyRootShelfPermission(controllerFunc controllers.Func[*rootshelvesdto.UpsertMyRootShelfPermissionRequestDto]) gin.HandlerFunc
	BindUpsertMyRootShelfPermissions(controllerFunc controllers.Func[*rootshelvesdto.UpsertMyRootShelfPermissionsRequestDto]) gin.HandlerFunc
	BindUpdateMyRootShelfPermission(controllerFunc controllers.Func[*rootshelvesdto.UpdateMyRootShelfPermissionRequestDto]) gin.HandlerFunc
	BindTransferMyRootShelfOwnership(controllerFunc controllers.Func[*rootshelvesdto.TransferMyRootShelfOwnershipRequestDto]) gin.HandlerFunc
	BindDeleteMyRootShelfPermission(controllerFunc controllers.Func[*rootshelvesdto.DeleteMyRootShelfPermissionRequestDto]) gin.HandlerFunc
	BindDeleteMyRootShelfPermissions(controllerFunc controllers.Func[*rootshelvesdto.DeleteMyRootShelfPermissionsRequestDto]) gin.HandlerFunc
	BindLeaveMyRootShelf(controllerFunc controllers.Func[*rootshelvesdto.LeaveMyRootShelfRequestDto]) gin.HandlerFunc
	BindLeaveMyRootShelves(controllerFunc controllers.Func[*rootshelvesdto.LeaveMyRootShelvesRequestDto]) gin.HandlerFunc
}

type RootShelfBinder struct{}

func NewRootShelfBinder() RootShelfBinderInterface {
	return &RootShelfBinder{}
}

func (b *RootShelfBinder) BindGetMyRootShelfById(controllerFunc controllers.Func[*rootshelvesdto.GetMyRootShelfByIdRequestDto]) gin.HandlerFunc {
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

func (b *RootShelfBinder) BindCreateRootShelf(controllerFunc controllers.Func[*rootshelvesdto.CreateRootShelfRequestDto]) gin.HandlerFunc {
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

func (b *RootShelfBinder) BindCreateRootShelves(controllerFunc controllers.Func[*rootshelvesdto.CreateRootShelvesRequestDto]) gin.HandlerFunc {
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

func (b *RootShelfBinder) BindUpdateMyRootShelfById(controllerFunc controllers.Func[*rootshelvesdto.UpdateMyRootShelfByIdRequestDto]) gin.HandlerFunc {
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

func (b *RootShelfBinder) BindUpdateMyRootShelvesByIds(controllerFunc controllers.Func[*rootshelvesdto.UpdateMyRootShelvesByIdsRequestDto]) gin.HandlerFunc {
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

func (b *RootShelfBinder) BindRestoreMyRootShelfById(controllerFunc controllers.Func[*rootshelvesdto.RestoreMyRootShelfByIdRequestDto]) gin.HandlerFunc {
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

func (b *RootShelfBinder) BindRestoreMyRootShelvesByIds(controllerFunc controllers.Func[*rootshelvesdto.RestoreMyRootShelvesByIdsRequestDto]) gin.HandlerFunc {
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

func (b *RootShelfBinder) BindDeleteMyRootShelfById(controllerFunc controllers.Func[*rootshelvesdto.DeleteMyRootShelfByIdRequestDto]) gin.HandlerFunc {
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

func (b *RootShelfBinder) BindDeleteMyRootShelvesByIds(controllerFunc controllers.Func[*rootshelvesdto.DeleteMyRootShelvesByIdsRequestDto]) gin.HandlerFunc {
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

func (b *RootShelfBinder) BindGetMyRootShelfPermission(controllerFunc controllers.Func[*rootshelvesdto.GetMyRootShelfPermissionRequestDto]) gin.HandlerFunc {
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

func (b *RootShelfBinder) BindCreateMyRootShelfPermission(controllerFunc controllers.Func[*rootshelvesdto.CreateMyRootShelfPermissionRequestDto]) gin.HandlerFunc {
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

func (b *RootShelfBinder) BindUpsertMyRootShelfPermission(controllerFunc controllers.Func[*rootshelvesdto.UpsertMyRootShelfPermissionRequestDto]) gin.HandlerFunc {
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

func (b *RootShelfBinder) BindUpsertMyRootShelfPermissions(controllerFunc controllers.Func[*rootshelvesdto.UpsertMyRootShelfPermissionsRequestDto]) gin.HandlerFunc {
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

func (b *RootShelfBinder) BindUpdateMyRootShelfPermission(controllerFunc controllers.Func[*rootshelvesdto.UpdateMyRootShelfPermissionRequestDto]) gin.HandlerFunc {
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

func (b *RootShelfBinder) BindTransferMyRootShelfOwnership(controllerFunc controllers.Func[*rootshelvesdto.TransferMyRootShelfOwnershipRequestDto]) gin.HandlerFunc {
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

func (b *RootShelfBinder) BindDeleteMyRootShelfPermission(controllerFunc controllers.Func[*rootshelvesdto.DeleteMyRootShelfPermissionRequestDto]) gin.HandlerFunc {
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

func (b *RootShelfBinder) BindDeleteMyRootShelfPermissions(controllerFunc controllers.Func[*rootshelvesdto.DeleteMyRootShelfPermissionsRequestDto]) gin.HandlerFunc {
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

func (b *RootShelfBinder) BindLeaveMyRootShelf(controllerFunc controllers.Func[*rootshelvesdto.LeaveMyRootShelfRequestDto]) gin.HandlerFunc {
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

func (b *RootShelfBinder) BindLeaveMyRootShelves(controllerFunc controllers.Func[*rootshelvesdto.LeaveMyRootShelvesRequestDto]) gin.HandlerFunc {
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
