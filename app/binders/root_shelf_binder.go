package binders

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	contexts "github.com/HiIamJeff67/notezy-backend/app/contexts"
	dtos "github.com/HiIamJeff67/notezy-backend/app/dtos"
	exceptions "github.com/HiIamJeff67/notezy-backend/app/exceptions"
	types "github.com/HiIamJeff67/notezy-backend/shared/types"
)

type RootShelfBinderInterface interface {
	BindGetMyRootShelfById(controllerFunc types.ControllerFunc[*dtos.GetMyRootShelfByIdReqDto]) gin.HandlerFunc
	BindCreateRootShelf(controllerFunc types.ControllerFunc[*dtos.CreateRootShelfReqDto]) gin.HandlerFunc
	BindCreateRootShelves(controllerFunc types.ControllerFunc[*dtos.CreateRootShelvesReqDto]) gin.HandlerFunc
	BindUpdateMyRootShelfById(controllerFunc types.ControllerFunc[*dtos.UpdateMyRootShelfByIdReqDto]) gin.HandlerFunc
	BindUpdateMyRootShelvesByIds(controllerFunc types.ControllerFunc[*dtos.UpdateMyRootShelvesByIdsReqDto]) gin.HandlerFunc
	BindRestoreMyRootShelfById(controllerFunc types.ControllerFunc[*dtos.RestoreMyRootShelfByIdReqDto]) gin.HandlerFunc
	BindRestoreMyRootShelvesByIds(controllerFunc types.ControllerFunc[*dtos.RestoreMyRootShelvesByIdsReqDto]) gin.HandlerFunc
	BindDeleteMyRootShelfById(controllerFunc types.ControllerFunc[*dtos.DeleteMyRootShelfByIdReqDto]) gin.HandlerFunc
	BindDeleteMyRootShelvesByIds(controllerFunc types.ControllerFunc[*dtos.DeleteMyRootShelvesByIdsReqDto]) gin.HandlerFunc

	BindGetMyRootShelfPermission(controllerFunc types.ControllerFunc[*dtos.GetMyRootShelfPermissionReqDto]) gin.HandlerFunc
	BindCreateMyRootShelfPermission(controllerFunc types.ControllerFunc[*dtos.CreateMyRootShelfPermissionReqDto]) gin.HandlerFunc
	BindUpsertMyRootShelfPermission(controllerFunc types.ControllerFunc[*dtos.UpsertMyRootShelfPermissionReqDto]) gin.HandlerFunc
	BindUpsertMyRootShelfPermissions(controllerFunc types.ControllerFunc[*dtos.UpsertMyRootShelfPermissionsReqDto]) gin.HandlerFunc
	BindUpdateMyRootShelfPermission(controllerFunc types.ControllerFunc[*dtos.UpdateMyRootShelfPermissionReqDto]) gin.HandlerFunc
	BindTransferMyRootShelfOwnership(controllerFunc types.ControllerFunc[*dtos.TransferMyRootShelfOwnershipReqDto]) gin.HandlerFunc
	BindDeleteMyRootShelfPermission(controllerFunc types.ControllerFunc[*dtos.DeleteMyRootShelfPermissionReqDto]) gin.HandlerFunc
	BindDeleteMyRootShelfPermissions(controllerFunc types.ControllerFunc[*dtos.DeleteMyRootShelfPermissionsReqDto]) gin.HandlerFunc
	BindLeaveMyRootShelf(controllerFunc types.ControllerFunc[*dtos.LeaveMyRootShelfReqDto]) gin.HandlerFunc
	BindLeaveMyRootShelves(controllerFunc types.ControllerFunc[*dtos.LeaveMyRootShelvesReqDto]) gin.HandlerFunc
}

type RootShelfBinder struct{}

func NewRootShelfBinder() RootShelfBinderInterface {
	return &RootShelfBinder{}
}

func (b *RootShelfBinder) BindGetMyRootShelfById(controllerFunc types.ControllerFunc[*dtos.GetMyRootShelfByIdReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.GetMyRootShelfByIdReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		isDeletedString := ctx.Query("isDeleted")
		if isDeletedString != "" {
			isDeleted, err := strconv.ParseBool(isDeletedString)
			if err != nil {
				exceptions.Shelf.InvalidInput().WithOrigin(err).SafelyAbortAndResponseWithJSON(ctx)
				return
			}
			reqDto.Param.IsDeleted = &isDeleted
		}

		rootShelfId, err := uuid.Parse(ctx.Param("rootShelfId"))
		if err != nil {
			exceptions.Shelf.InvalidInput().WithOrigin(err).SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.Param.RootShelfId = rootShelfId

		controllerFunc(ctx, &reqDto)
	}
}

func (b *RootShelfBinder) BindCreateRootShelf(controllerFunc types.ControllerFunc[*dtos.CreateRootShelfReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.CreateRootShelfReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.Shelf.InvalidDto().WithOrigin(err)
			exception.SafelyAbortAndResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *RootShelfBinder) BindCreateRootShelves(controllerFunc types.ControllerFunc[*dtos.CreateRootShelvesReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.CreateRootShelvesReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.Shelf.InvalidDto().WithOrigin(err)
			exception.SafelyAbortAndResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *RootShelfBinder) BindUpdateMyRootShelfById(controllerFunc types.ControllerFunc[*dtos.UpdateMyRootShelfByIdReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.UpdateMyRootShelfByIdReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.Shelf.InvalidDto().WithOrigin(err)
			exception.SafelyAbortAndResponseWithJSON(ctx)
			return
		}

		rootShelfId, err := uuid.Parse(ctx.Param("rootShelfId"))
		if err != nil {
			exceptions.Shelf.InvalidInput().WithOrigin(err).SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.Body.RootShelfId = rootShelfId

		controllerFunc(ctx, &reqDto)
	}
}

func (b *RootShelfBinder) BindUpdateMyRootShelvesByIds(controllerFunc types.ControllerFunc[*dtos.UpdateMyRootShelvesByIdsReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.UpdateMyRootShelvesByIdsReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.Shelf.InvalidDto().WithOrigin(err)
			exception.SafelyAbortAndResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *RootShelfBinder) BindRestoreMyRootShelfById(controllerFunc types.ControllerFunc[*dtos.RestoreMyRootShelfByIdReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.RestoreMyRootShelfByIdReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.Shelf.InvalidDto().WithOrigin(err)
			exception.SafelyAbortAndResponseWithJSON(ctx)
			return
		}

		rootShelfId, err := uuid.Parse(ctx.Param("rootShelfId"))
		if err != nil {
			exceptions.Shelf.InvalidInput().WithOrigin(err).SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.Body.RootShelfId = rootShelfId

		controllerFunc(ctx, &reqDto)
	}
}

func (b *RootShelfBinder) BindRestoreMyRootShelvesByIds(controllerFunc types.ControllerFunc[*dtos.RestoreMyRootShelvesByIdsReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.RestoreMyRootShelvesByIdsReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.Shelf.InvalidDto().WithOrigin(err)
			exception.SafelyAbortAndResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *RootShelfBinder) BindDeleteMyRootShelfById(controllerFunc types.ControllerFunc[*dtos.DeleteMyRootShelfByIdReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.DeleteMyRootShelfByIdReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.Shelf.InvalidDto().WithOrigin(err)
			exception.SafelyAbortAndResponseWithJSON(ctx)
			return
		}

		rootShelfId, err := uuid.Parse(ctx.Param("rootShelfId"))
		if err != nil {
			exceptions.Shelf.InvalidInput().WithOrigin(err).SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.Body.RootShelfId = rootShelfId

		controllerFunc(ctx, &reqDto)
	}
}

func (b *RootShelfBinder) BindDeleteMyRootShelvesByIds(controllerFunc types.ControllerFunc[*dtos.DeleteMyRootShelvesByIdsReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.DeleteMyRootShelvesByIdsReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.Shelf.InvalidDto().WithOrigin(err)
			exception.SafelyAbortAndResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *RootShelfBinder) BindGetMyRootShelfPermission(controllerFunc types.ControllerFunc[*dtos.GetMyRootShelfPermissionReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.GetMyRootShelfPermissionReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		rootShelfId, err := uuid.Parse(ctx.Param("rootShelfId"))
		if err != nil {
			exceptions.Shelf.InvalidInput().WithOrigin(err).Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.Param.RootShelfId = rootShelfId

		userPublicId, err := uuid.Parse(ctx.Param("userPublicId"))
		if err != nil {
			exceptions.Shelf.InvalidInput().WithOrigin(err).Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.Param.UserPublicId = userPublicId

		controllerFunc(ctx, &reqDto)
	}
}

func (b *RootShelfBinder) BindCreateMyRootShelfPermission(controllerFunc types.ControllerFunc[*dtos.CreateMyRootShelfPermissionReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.CreateMyRootShelfPermissionReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		rootShelfId, err := uuid.Parse(ctx.Param("rootShelfId"))
		if err != nil {
			exceptions.Shelf.InvalidInput().WithOrigin(err).Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.Param.RootShelfId = rootShelfId

		userPublicId, err := uuid.Parse(ctx.Param("userPublicId"))
		if err != nil {
			exceptions.Shelf.InvalidInput().WithOrigin(err).Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.Param.UserPublicId = userPublicId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exceptions.Shelf.InvalidDto().WithOrigin(err).Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *RootShelfBinder) BindUpsertMyRootShelfPermission(
	controllerFunc types.ControllerFunc[*dtos.UpsertMyRootShelfPermissionReqDto],
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.UpsertMyRootShelfPermissionReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		rootShelfId, err := uuid.Parse(ctx.Param("rootShelfId"))
		if err != nil {
			exceptions.Shelf.InvalidInput().WithOrigin(err).Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.Param.RootShelfId = rootShelfId

		userPublicId, err := uuid.Parse(ctx.Param("userPublicId"))
		if err != nil {
			exceptions.Shelf.InvalidInput().WithOrigin(err).Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.Param.UserPublicId = userPublicId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exceptions.Shelf.InvalidDto().WithOrigin(err).Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *RootShelfBinder) BindUpsertMyRootShelfPermissions(
	controllerFunc types.ControllerFunc[*dtos.UpsertMyRootShelfPermissionsReqDto],
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.UpsertMyRootShelfPermissionsReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		rootShelfId, err := uuid.Parse(ctx.Param("rootShelfId"))
		if err != nil {
			exceptions.Shelf.InvalidInput().WithOrigin(err).Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.Param.RootShelfId = rootShelfId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exceptions.Shelf.InvalidDto().WithOrigin(err).Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *RootShelfBinder) BindUpdateMyRootShelfPermission(controllerFunc types.ControllerFunc[*dtos.UpdateMyRootShelfPermissionReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.UpdateMyRootShelfPermissionReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		rootShelfId, err := uuid.Parse(ctx.Param("rootShelfId"))
		if err != nil {
			exceptions.Shelf.InvalidInput().WithOrigin(err).Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.Param.RootShelfId = rootShelfId

		userPublicId, err := uuid.Parse(ctx.Param("userPublicId"))
		if err != nil {
			exceptions.Shelf.InvalidInput().WithOrigin(err).Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.Param.UserPublicId = userPublicId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exceptions.Shelf.InvalidDto().WithOrigin(err).Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *RootShelfBinder) BindTransferMyRootShelfOwnership(
	controllerFunc types.ControllerFunc[*dtos.TransferMyRootShelfOwnershipReqDto],
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.TransferMyRootShelfOwnershipReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		rootShelfId, err := uuid.Parse(ctx.Param("rootShelfId"))
		if err != nil {
			exceptions.Shelf.InvalidInput().WithOrigin(err).Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.Param.RootShelfId = rootShelfId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exceptions.Shelf.InvalidDto().WithOrigin(err).Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *RootShelfBinder) BindDeleteMyRootShelfPermission(
	controllerFunc types.ControllerFunc[*dtos.DeleteMyRootShelfPermissionReqDto],
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.DeleteMyRootShelfPermissionReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		rootShelfId, err := uuid.Parse(ctx.Param("rootShelfId"))
		if err != nil {
			exceptions.Shelf.InvalidInput().WithOrigin(err).Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.Param.RootShelfId = rootShelfId

		userPublicId, err := uuid.Parse(ctx.Param("userPublicId"))
		if err != nil {
			exceptions.Shelf.InvalidInput().WithOrigin(err).Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.Param.UserPublicId = userPublicId

		controllerFunc(ctx, &reqDto)
	}
}

func (b *RootShelfBinder) BindDeleteMyRootShelfPermissions(
	controllerFunc types.ControllerFunc[*dtos.DeleteMyRootShelfPermissionsReqDto],
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.DeleteMyRootShelfPermissionsReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		rootShelfId, err := uuid.Parse(ctx.Param("rootShelfId"))
		if err != nil {
			exceptions.Shelf.InvalidInput().WithOrigin(err).Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.Param.RootShelfId = rootShelfId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exceptions.Shelf.InvalidDto().WithOrigin(err).Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *RootShelfBinder) BindLeaveMyRootShelf(controllerFunc types.ControllerFunc[*dtos.LeaveMyRootShelfReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.LeaveMyRootShelfReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")
		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		rootShelfId, err := uuid.Parse(ctx.Param("rootShelfId"))
		if err != nil {
			exceptions.Shelf.InvalidInput().WithOrigin(err).Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.Param.RootShelfId = rootShelfId

		controllerFunc(ctx, &reqDto)
	}
}

func (b *RootShelfBinder) BindLeaveMyRootShelves(controllerFunc types.ControllerFunc[*dtos.LeaveMyRootShelvesReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.LeaveMyRootShelvesReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")
		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exceptions.Shelf.InvalidDto().WithOrigin(err).Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		controllerFunc(ctx, &reqDto)
	}
}
