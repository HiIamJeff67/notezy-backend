package binders

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	exceptions "github.com/HiIamJeff67/notegic-backend/contracts/types/exceptions"

	exceptionwriter "github.com/HiIamJeff67/notegic-backend/shared/util/exceptionwriter"

	apicontract "github.com/HiIamJeff67/notegic-backend/contracts/core/v1/api/root-shelves"

	controllers "github.com/HiIamJeff67/notegic-backend/internal/apigateway/transports/api/controllers"
)

type RootShelfBinderInterface interface {
	BindGetMyRootShelfById(controllerFunc controllers.Func[*apicontract.GetMyRootShelfByIdRequestDto]) gin.HandlerFunc
	BindCreateRootShelf(controllerFunc controllers.Func[*apicontract.CreateRootShelfRequestDto]) gin.HandlerFunc
	BindCreateRootShelves(controllerFunc controllers.Func[*apicontract.CreateRootShelvesRequestDto]) gin.HandlerFunc
	BindUpdateMyRootShelfById(controllerFunc controllers.Func[*apicontract.UpdateMyRootShelfByIdRequestDto]) gin.HandlerFunc
	BindUpdateMyRootShelvesByIds(controllerFunc controllers.Func[*apicontract.UpdateMyRootShelvesByIdsRequestDto]) gin.HandlerFunc
	BindRestoreMyRootShelfById(controllerFunc controllers.Func[*apicontract.RestoreMyRootShelfByIdRequestDto]) gin.HandlerFunc
	BindRestoreMyRootShelvesByIds(controllerFunc controllers.Func[*apicontract.RestoreMyRootShelvesByIdsRequestDto]) gin.HandlerFunc
	BindDeleteMyRootShelfById(controllerFunc controllers.Func[*apicontract.DeleteMyRootShelfByIdRequestDto]) gin.HandlerFunc
	BindDeleteMyRootShelvesByIds(controllerFunc controllers.Func[*apicontract.DeleteMyRootShelvesByIdsRequestDto]) gin.HandlerFunc
	BindGetMyRootShelfPermission(controllerFunc controllers.Func[*apicontract.GetMyRootShelfPermissionRequestDto]) gin.HandlerFunc
	BindCreateMyRootShelfPermission(controllerFunc controllers.Func[*apicontract.CreateMyRootShelfPermissionRequestDto]) gin.HandlerFunc
	BindUpsertMyRootShelfPermission(controllerFunc controllers.Func[*apicontract.UpsertMyRootShelfPermissionRequestDto]) gin.HandlerFunc
	BindUpsertMyRootShelfPermissions(controllerFunc controllers.Func[*apicontract.UpsertMyRootShelfPermissionsRequestDto]) gin.HandlerFunc
	BindUpdateMyRootShelfPermission(controllerFunc controllers.Func[*apicontract.UpdateMyRootShelfPermissionRequestDto]) gin.HandlerFunc
	BindTransferMyRootShelfOwnership(controllerFunc controllers.Func[*apicontract.TransferMyRootShelfOwnershipRequestDto]) gin.HandlerFunc
	BindDeleteMyRootShelfPermission(controllerFunc controllers.Func[*apicontract.DeleteMyRootShelfPermissionRequestDto]) gin.HandlerFunc
	BindDeleteMyRootShelfPermissions(controllerFunc controllers.Func[*apicontract.DeleteMyRootShelfPermissionsRequestDto]) gin.HandlerFunc
	BindLeaveMyRootShelf(controllerFunc controllers.Func[*apicontract.LeaveMyRootShelfRequestDto]) gin.HandlerFunc
	BindLeaveMyRootShelves(controllerFunc controllers.Func[*apicontract.LeaveMyRootShelvesRequestDto]) gin.HandlerFunc
}

type RootShelfBinder struct{}

func NewRootShelfBinder() RootShelfBinderInterface {
	return &RootShelfBinder{}
}

func (b *RootShelfBinder) BindGetMyRootShelfById(controllerFunc controllers.Func[*apicontract.GetMyRootShelfByIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &apicontract.GetMyRootShelfByIdRequestDto{}
		request.Header.UserAgent = ctx.GetHeader("User-Agent")

		isDeletedString := ctx.Query("isDeleted")
		if isDeletedString != "" {
			isDeleted, err := strconv.ParseBool(isDeletedString)
			if err != nil {
				exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Shelf").WithOrigin(err), ctx)
				return
			}
			request.Param.IsDeleted = &isDeleted
		}

		rootShelfId, err := uuid.Parse(ctx.Param("root-shelf-id"))
		if err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Shelf").WithOrigin(err), ctx)
			return
		}
		request.Param.RootShelfId = rootShelfId

		controllerFunc(ctx, request)
	}
}

func (b *RootShelfBinder) BindCreateRootShelf(controllerFunc controllers.Func[*apicontract.CreateRootShelfRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &apicontract.CreateRootShelfRequestDto{}
		request.Header.UserAgent = ctx.GetHeader("User-Agent")
		if err := ctx.ShouldBindJSON(&request.Body); err != nil {
			exception := exceptions.InvalidDto("Shelf").WithOrigin(err)
			exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
			return
		}

		controllerFunc(ctx, request)
	}
}

func (b *RootShelfBinder) BindCreateRootShelves(controllerFunc controllers.Func[*apicontract.CreateRootShelvesRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &apicontract.CreateRootShelvesRequestDto{}
		request.Header.UserAgent = ctx.GetHeader("User-Agent")
		if err := ctx.ShouldBindJSON(&request.Body); err != nil {
			exception := exceptions.InvalidDto("Shelf").WithOrigin(err)
			exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
			return
		}

		controllerFunc(ctx, request)
	}
}

func (b *RootShelfBinder) BindUpdateMyRootShelfById(controllerFunc controllers.Func[*apicontract.UpdateMyRootShelfByIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &apicontract.UpdateMyRootShelfByIdRequestDto{}
		request.Header.UserAgent = ctx.GetHeader("User-Agent")
		if err := ctx.ShouldBindJSON(&request.Body); err != nil {
			exception := exceptions.InvalidDto("Shelf").WithOrigin(err)
			exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
			return
		}

		rootShelfId, err := uuid.Parse(ctx.Param("root-shelf-id"))
		if err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Shelf").WithOrigin(err), ctx)
			return
		}
		request.Param.RootShelfId = rootShelfId

		controllerFunc(ctx, request)
	}
}

func (b *RootShelfBinder) BindUpdateMyRootShelvesByIds(controllerFunc controllers.Func[*apicontract.UpdateMyRootShelvesByIdsRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto apicontract.UpdateMyRootShelvesByIdsRequestDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.InvalidDto("Shelf").WithOrigin(err)
			exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *RootShelfBinder) BindRestoreMyRootShelfById(controllerFunc controllers.Func[*apicontract.RestoreMyRootShelfByIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto apicontract.RestoreMyRootShelfByIdRequestDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.InvalidDto("Shelf").WithOrigin(err)
			exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
			return
		}

		rootShelfId, err := uuid.Parse(ctx.Param("root-shelf-id"))
		if err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Shelf").WithOrigin(err), ctx)
			return
		}
		reqDto.Body.RootShelfId = rootShelfId

		controllerFunc(ctx, &reqDto)
	}
}

func (b *RootShelfBinder) BindRestoreMyRootShelvesByIds(controllerFunc controllers.Func[*apicontract.RestoreMyRootShelvesByIdsRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto apicontract.RestoreMyRootShelvesByIdsRequestDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.InvalidDto("Shelf").WithOrigin(err)
			exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *RootShelfBinder) BindDeleteMyRootShelfById(controllerFunc controllers.Func[*apicontract.DeleteMyRootShelfByIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto apicontract.DeleteMyRootShelfByIdRequestDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.InvalidDto("Shelf").WithOrigin(err)
			exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
			return
		}

		rootShelfId, err := uuid.Parse(ctx.Param("root-shelf-id"))
		if err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Shelf").WithOrigin(err), ctx)
			return
		}
		reqDto.Body.RootShelfId = rootShelfId

		controllerFunc(ctx, &reqDto)
	}
}

func (b *RootShelfBinder) BindDeleteMyRootShelvesByIds(controllerFunc controllers.Func[*apicontract.DeleteMyRootShelvesByIdsRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto apicontract.DeleteMyRootShelvesByIdsRequestDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.InvalidDto("Shelf").WithOrigin(err)
			exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *RootShelfBinder) BindGetMyRootShelfPermission(controllerFunc controllers.Func[*apicontract.GetMyRootShelfPermissionRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &apicontract.GetMyRootShelfPermissionRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		rootShelfId, err := uuid.Parse(ctx.Param("root-shelf-id"))
		if err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Shelf").WithOrigin(err), ctx)
			return
		}
		requestDto.Param.RootShelfId = rootShelfId

		userPublicId, err := uuid.Parse(ctx.Param("user-public-id"))
		if err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Shelf").WithOrigin(err), ctx)
			return
		}
		requestDto.Param.UserPublicId = userPublicId

		controllerFunc(ctx, requestDto)
	}
}

func (b *RootShelfBinder) BindCreateMyRootShelfPermission(controllerFunc controllers.Func[*apicontract.CreateMyRootShelfPermissionRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &apicontract.CreateMyRootShelfPermissionRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		rootShelfId, err := uuid.Parse(ctx.Param("root-shelf-id"))
		if err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Shelf").WithOrigin(err), ctx)
			return
		}
		requestDto.Param.RootShelfId = rootShelfId

		userPublicId, err := uuid.Parse(ctx.Param("user-public-id"))
		if err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Shelf").WithOrigin(err), ctx)
			return
		}
		requestDto.Param.UserPublicId = userPublicId

		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("Shelf").WithOrigin(err), ctx)
			return
		}

		controllerFunc(ctx, requestDto)
	}
}

func (b *RootShelfBinder) BindUpsertMyRootShelfPermission(controllerFunc controllers.Func[*apicontract.UpsertMyRootShelfPermissionRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &apicontract.UpsertMyRootShelfPermissionRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		rootShelfId, err := uuid.Parse(ctx.Param("root-shelf-id"))
		if err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Shelf").WithOrigin(err), ctx)
			return
		}
		requestDto.Param.RootShelfId = rootShelfId

		userPublicId, err := uuid.Parse(ctx.Param("user-public-id"))
		if err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Shelf").WithOrigin(err), ctx)
			return
		}
		requestDto.Param.UserPublicId = userPublicId

		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("Shelf").WithOrigin(err), ctx)
			return
		}

		controllerFunc(ctx, requestDto)
	}
}

func (b *RootShelfBinder) BindUpsertMyRootShelfPermissions(controllerFunc controllers.Func[*apicontract.UpsertMyRootShelfPermissionsRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &apicontract.UpsertMyRootShelfPermissionsRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		rootShelfId, err := uuid.Parse(ctx.Param("root-shelf-id"))
		if err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Shelf").WithOrigin(err), ctx)
			return
		}
		requestDto.Param.RootShelfId = rootShelfId

		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("Shelf").WithOrigin(err), ctx)
			return
		}

		controllerFunc(ctx, requestDto)
	}
}

func (b *RootShelfBinder) BindUpdateMyRootShelfPermission(controllerFunc controllers.Func[*apicontract.UpdateMyRootShelfPermissionRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &apicontract.UpdateMyRootShelfPermissionRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		rootShelfId, err := uuid.Parse(ctx.Param("root-shelf-id"))
		if err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Shelf").WithOrigin(err), ctx)
			return
		}
		requestDto.Param.RootShelfId = rootShelfId

		userPublicId, err := uuid.Parse(ctx.Param("user-public-id"))
		if err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Shelf").WithOrigin(err), ctx)
			return
		}
		requestDto.Param.UserPublicId = userPublicId

		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("Shelf").WithOrigin(err), ctx)
			return
		}

		controllerFunc(ctx, requestDto)
	}
}

func (b *RootShelfBinder) BindTransferMyRootShelfOwnership(controllerFunc controllers.Func[*apicontract.TransferMyRootShelfOwnershipRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &apicontract.TransferMyRootShelfOwnershipRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		rootShelfId, err := uuid.Parse(ctx.Param("root-shelf-id"))
		if err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Shelf").WithOrigin(err), ctx)
			return
		}
		requestDto.Param.RootShelfId = rootShelfId

		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("Shelf").WithOrigin(err), ctx)
			return
		}

		controllerFunc(ctx, requestDto)
	}
}

func (b *RootShelfBinder) BindDeleteMyRootShelfPermission(controllerFunc controllers.Func[*apicontract.DeleteMyRootShelfPermissionRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &apicontract.DeleteMyRootShelfPermissionRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		rootShelfId, err := uuid.Parse(ctx.Param("root-shelf-id"))
		if err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Shelf").WithOrigin(err), ctx)
			return
		}
		requestDto.Param.RootShelfId = rootShelfId

		userPublicId, err := uuid.Parse(ctx.Param("user-public-id"))
		if err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Shelf").WithOrigin(err), ctx)
			return
		}
		requestDto.Param.UserPublicId = userPublicId

		controllerFunc(ctx, requestDto)
	}
}

func (b *RootShelfBinder) BindDeleteMyRootShelfPermissions(controllerFunc controllers.Func[*apicontract.DeleteMyRootShelfPermissionsRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &apicontract.DeleteMyRootShelfPermissionsRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		rootShelfId, err := uuid.Parse(ctx.Param("root-shelf-id"))
		if err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Shelf").WithOrigin(err), ctx)
			return
		}
		requestDto.Param.RootShelfId = rootShelfId

		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("Shelf").WithOrigin(err), ctx)
			return
		}

		controllerFunc(ctx, requestDto)
	}
}

func (b *RootShelfBinder) BindLeaveMyRootShelf(controllerFunc controllers.Func[*apicontract.LeaveMyRootShelfRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &apicontract.LeaveMyRootShelfRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		rootShelfId, err := uuid.Parse(ctx.Param("root-shelf-id"))
		if err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Shelf").WithOrigin(err), ctx)
			return
		}
		requestDto.Param.RootShelfId = rootShelfId

		controllerFunc(ctx, requestDto)
	}
}

func (b *RootShelfBinder) BindLeaveMyRootShelves(controllerFunc controllers.Func[*apicontract.LeaveMyRootShelvesRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestDto := &apicontract.LeaveMyRootShelvesRequestDto{}
		requestDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		if err := ctx.ShouldBindJSON(&requestDto.Body); err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("Shelf").WithOrigin(err), ctx)
			return
		}
		controllerFunc(ctx, requestDto)
	}
}
