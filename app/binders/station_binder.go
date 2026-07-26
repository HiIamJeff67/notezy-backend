package binders

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	contexts "github.com/HiIamJeff67/notezy-backend/app/contexts"
	dtos "github.com/HiIamJeff67/notezy-backend/app/dtos"
	exceptions "github.com/HiIamJeff67/notezy-backend/app/exceptions"
	enums "github.com/HiIamJeff67/notezy-backend/app/models/schemas/enums"
	types "github.com/HiIamJeff67/notezy-backend/shared/types"
)

type StationBinderInterface interface {
	BindGetMyStationById(controllerFunc types.ControllerFunc[*dtos.GetMyStationByIdReqDto]) gin.HandlerFunc
	BindGetAllMyStations(controllerFunc types.ControllerFunc[*dtos.GetAllMyStationsReqDto]) gin.HandlerFunc
	BindCreateStation(controllerFunc types.ControllerFunc[*dtos.CreateStationReqDto]) gin.HandlerFunc
	BindCreateStations(controllerFunc types.ControllerFunc[*dtos.CreateStationsReqDto]) gin.HandlerFunc
	BindUpdateMyStationById(controllerFunc types.ControllerFunc[*dtos.UpdateMyStationByIdReqDto]) gin.HandlerFunc
	BindUpdateMyStationsByIds(controllerFunc types.ControllerFunc[*dtos.UpdateMyStationsByIdsReqDto]) gin.HandlerFunc
	BindRestoreMyStationById(controllerFunc types.ControllerFunc[*dtos.RestoreMyStationByIdReqDto]) gin.HandlerFunc
	BindRestoreMyStationsByIds(controllerFunc types.ControllerFunc[*dtos.RestoreMyStationsByIdsReqDto]) gin.HandlerFunc
	BindDeleteMyStationById(controllerFunc types.ControllerFunc[*dtos.DeleteMyStationByIdReqDto]) gin.HandlerFunc
	BindDeleteMyStationsByIds(controllerFunc types.ControllerFunc[*dtos.DeleteMyStationsByIdsReqDto]) gin.HandlerFunc
	BindHardDeleteMyStationById(controllerFunc types.ControllerFunc[*dtos.HardDeleteMyStationByIdReqDto]) gin.HandlerFunc
	BindHardDeleteMyStationsByIds(controllerFunc types.ControllerFunc[*dtos.HardDeleteMyStationsByIdsReqDto]) gin.HandlerFunc

	BindVisualizeMyTotalCount(controllerFunc types.ControllerFunc[*dtos.VisualizeMyTotalCountReqDto]) gin.HandlerFunc

	BindGetMyStationPermission(controllerFunc types.ControllerFunc[*dtos.GetMyStationPermissionReqDto]) gin.HandlerFunc
	BindCreateMyStationPermission(controllerFunc types.ControllerFunc[*dtos.CreateMyStationPermissionReqDto]) gin.HandlerFunc
	BindUpsertMyStationPermission(controllerFunc types.ControllerFunc[*dtos.UpsertMyStationPermissionReqDto]) gin.HandlerFunc
	BindUpsertMyStationPermissions(controllerFunc types.ControllerFunc[*dtos.UpsertMyStationPermissionsReqDto]) gin.HandlerFunc
	BindUpdateMyStationPermission(controllerFunc types.ControllerFunc[*dtos.UpdateMyStationPermissionReqDto]) gin.HandlerFunc
	BindTransferMyStationOwnership(controllerFunc types.ControllerFunc[*dtos.TransferMyStationOwnershipReqDto]) gin.HandlerFunc
	BindDeleteMyStationPermission(controllerFunc types.ControllerFunc[*dtos.DeleteMyStationPermissionReqDto]) gin.HandlerFunc
	BindDeleteMyStationPermissions(controllerFunc types.ControllerFunc[*dtos.DeleteMyStationPermissionsReqDto]) gin.HandlerFunc
	BindLeaveMyStation(controllerFunc types.ControllerFunc[*dtos.LeaveMyStationReqDto]) gin.HandlerFunc
	BindLeaveMyStations(controllerFunc types.ControllerFunc[*dtos.LeaveMyStationsReqDto]) gin.HandlerFunc
}

type StationBinder struct{}

func NewStationBinder() StationBinderInterface {
	return &StationBinder{}
}

func (b *StationBinder) BindGetMyStationById(controllerFunc types.ControllerFunc[*dtos.GetMyStationByIdReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.GetMyStationByIdReqDto

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
				exceptions.Station.InvalidInput().WithOrigin(err).SafelyAbortAndResponseWithJSON(ctx)
				return
			}
			reqDto.Param.IsDeleted = &isDeleted
		}

		stationIdString := ctx.Query("stationId")
		if stationIdString == "" {
			exceptions.Station.InvalidInput().WithOrigin(fmt.Errorf("stationId is required")).SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		stationId, err := uuid.Parse(stationIdString)
		if err != nil {
			exceptions.Station.InvalidInput().WithOrigin(err).SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.Param.StationId = stationId

		controllerFunc(ctx, &reqDto)
	}
}

func (b *StationBinder) BindGetAllMyStations(controllerFunc types.ControllerFunc[*dtos.GetAllMyStationsReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.GetAllMyStationsReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		areDeletedString := ctx.Query("areDeleted")
		if areDeletedString != "" {
			areDeleted, err := strconv.ParseBool(areDeletedString)
			if err != nil {
				exceptions.Station.InvalidInput().WithOrigin(err).SafelyAbortAndResponseWithJSON(ctx)
				return
			}
			reqDto.Param.AreDeleted = &areDeleted
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *StationBinder) BindCreateStation(controllerFunc types.ControllerFunc[*dtos.CreateStationReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.CreateStationReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.Station.InvalidDto().WithOrigin(err)
			exception.SafelyAbortAndResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *StationBinder) BindCreateStations(controllerFunc types.ControllerFunc[*dtos.CreateStationsReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.CreateStationsReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.Station.InvalidDto().WithOrigin(err)
			exception.SafelyAbortAndResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *StationBinder) BindUpdateMyStationById(controllerFunc types.ControllerFunc[*dtos.UpdateMyStationByIdReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.UpdateMyStationByIdReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.Station.InvalidDto().WithOrigin(err)
			exception.SafelyAbortAndResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *StationBinder) BindUpdateMyStationsByIds(controllerFunc types.ControllerFunc[*dtos.UpdateMyStationsByIdsReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.UpdateMyStationsByIdsReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.Station.InvalidDto().WithOrigin(err)
			exception.SafelyAbortAndResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *StationBinder) BindRestoreMyStationById(controllerFunc types.ControllerFunc[*dtos.RestoreMyStationByIdReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.RestoreMyStationByIdReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.Station.InvalidDto().WithOrigin(err)
			exception.SafelyAbortAndResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *StationBinder) BindRestoreMyStationsByIds(controllerFunc types.ControllerFunc[*dtos.RestoreMyStationsByIdsReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.RestoreMyStationsByIdsReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.Station.InvalidDto().WithOrigin(err)
			exception.SafelyAbortAndResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *StationBinder) BindDeleteMyStationById(controllerFunc types.ControllerFunc[*dtos.DeleteMyStationByIdReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.DeleteMyStationByIdReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.Station.InvalidDto().WithOrigin(err)
			exception.SafelyAbortAndResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *StationBinder) BindDeleteMyStationsByIds(controllerFunc types.ControllerFunc[*dtos.DeleteMyStationsByIdsReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.DeleteMyStationsByIdsReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.Station.InvalidDto().WithOrigin(err)
			exception.SafelyAbortAndResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *StationBinder) BindHardDeleteMyStationById(controllerFunc types.ControllerFunc[*dtos.HardDeleteMyStationByIdReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.HardDeleteMyStationByIdReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.Station.InvalidDto().WithOrigin(err)
			exception.SafelyAbortAndResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *StationBinder) BindHardDeleteMyStationsByIds(controllerFunc types.ControllerFunc[*dtos.HardDeleteMyStationsByIdsReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.HardDeleteMyStationsByIdsReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exception := exceptions.Station.InvalidDto().WithOrigin(err)
			exception.SafelyAbortAndResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

/* ============================== Binder Methods for Visualization ============================== */

func (b *StationBinder) BindVisualizeMyTotalCount(controllerFunc types.ControllerFunc[*dtos.VisualizeMyTotalCountReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.VisualizeMyTotalCountReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		permissionString := ctx.Query("permission")
		if permissionString == "" {
			exceptions.Station.InvalidInput().WithOrigin(fmt.Errorf("permission is required")).SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		permission, err := enums.ConvertStringToAccessControlPermission(permissionString)
		if err != nil {
			exceptions.Station.InvalidInput().WithOrigin(err).SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.Param.Permission = *permission

		controllerFunc(ctx, &reqDto)
	}
}

/* ============================== Binder Methods for Station Permissions ============================== */

func (b *StationBinder) BindGetMyStationPermission(controllerFunc types.ControllerFunc[*dtos.GetMyStationPermissionReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.GetMyStationPermissionReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		stationId, err := uuid.Parse(ctx.Param("stationId"))
		if err != nil {
			exceptions.Station.InvalidInput().WithOrigin(err).SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.Param.StationId = stationId

		userPublicId, err := uuid.Parse(ctx.Param("userPublicId"))
		if err != nil {
			exceptions.Station.InvalidInput().WithOrigin(err).SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.Param.UserPublicId = userPublicId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exceptions.Station.InvalidDto().WithOrigin(err).SafelyAbortAndResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *StationBinder) BindCreateMyStationPermission(controllerFunc types.ControllerFunc[*dtos.CreateMyStationPermissionReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.CreateMyStationPermissionReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		stationId, err := uuid.Parse(ctx.Param("stationId"))
		if err != nil {
			exceptions.Station.InvalidInput().WithOrigin(err).SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.Param.StationId = stationId

		userPublicId, err := uuid.Parse(ctx.Param("userPublicId"))
		if err != nil {
			exceptions.Station.InvalidInput().WithOrigin(err).SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.Param.UserPublicId = userPublicId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exceptions.Station.InvalidDto().WithOrigin(err).SafelyAbortAndResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *StationBinder) BindUpsertMyStationPermission(controllerFunc types.ControllerFunc[*dtos.UpsertMyStationPermissionReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.UpsertMyStationPermissionReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		stationId, err := uuid.Parse(ctx.Param("stationId"))
		if err != nil {
			exceptions.Station.InvalidInput().WithOrigin(err).SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.Param.StationId = stationId

		userPublicId, err := uuid.Parse(ctx.Param("userPublicId"))
		if err != nil {
			exceptions.Station.InvalidInput().WithOrigin(err).SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.Param.UserPublicId = userPublicId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exceptions.Station.InvalidDto().WithOrigin(err).SafelyAbortAndResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *StationBinder) BindUpsertMyStationPermissions(
	controllerFunc types.ControllerFunc[*dtos.UpsertMyStationPermissionsReqDto],
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.UpsertMyStationPermissionsReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		stationId, err := uuid.Parse(ctx.Param("stationId"))
		if err != nil {
			exceptions.Station.InvalidInput().WithOrigin(err).SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.Param.StationId = stationId

		controllerFunc(ctx, &reqDto)
	}
}

func (b *StationBinder) BindUpdateMyStationPermission(controllerFunc types.ControllerFunc[*dtos.UpdateMyStationPermissionReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.UpdateMyStationPermissionReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		stationId, err := uuid.Parse(ctx.Param("stationId"))
		if err != nil {
			exceptions.Station.InvalidInput().WithOrigin(err).SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.Param.StationId = stationId

		userPublicId, err := uuid.Parse(ctx.Param("userPublicId"))
		if err != nil {
			exceptions.Station.InvalidInput().WithOrigin(err).SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.Param.UserPublicId = userPublicId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exceptions.Station.InvalidDto().WithOrigin(err).SafelyAbortAndResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *StationBinder) BindTransferMyStationOwnership(
	controllerFunc types.ControllerFunc[*dtos.TransferMyStationOwnershipReqDto],
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.TransferMyStationOwnershipReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		stationId, err := uuid.Parse(ctx.Param("stationId"))
		if err != nil {
			exceptions.Station.InvalidInput().WithOrigin(err).SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.Param.StationId = stationId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exceptions.Station.InvalidDto().WithOrigin(err).SafelyAbortAndResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *StationBinder) BindDeleteMyStationPermission(controllerFunc types.ControllerFunc[*dtos.DeleteMyStationPermissionReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.DeleteMyStationPermissionReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		stationId, err := uuid.Parse(ctx.Param("stationId"))
		if err != nil {
			exceptions.Station.InvalidInput().WithOrigin(err).SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.Param.StationId = stationId

		userPublicId, err := uuid.Parse(ctx.Param("userPublicId"))
		if err != nil {
			exceptions.Station.InvalidInput().WithOrigin(err).SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.Param.UserPublicId = userPublicId

		controllerFunc(ctx, &reqDto)
	}
}

func (b *StationBinder) BindDeleteMyStationPermissions(
	controllerFunc types.ControllerFunc[*dtos.DeleteMyStationPermissionsReqDto],
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.DeleteMyStationPermissionsReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		stationId, err := uuid.Parse(ctx.Param("stationId"))
		if err != nil {
			exceptions.Station.InvalidInput().WithOrigin(err).SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.Param.StationId = stationId

		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exceptions.Station.InvalidDto().WithOrigin(err).SafelyAbortAndResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}

func (b *StationBinder) BindLeaveMyStation(controllerFunc types.ControllerFunc[*dtos.LeaveMyStationReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.LeaveMyStationReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId

		stationId, err := uuid.Parse(ctx.Param("stationId"))
		if err != nil {
			exceptions.Station.InvalidInput().WithOrigin(err).SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.Param.StationId = stationId

		controllerFunc(ctx, &reqDto)
	}
}

func (b *StationBinder) BindLeaveMyStations(controllerFunc types.ControllerFunc[*dtos.LeaveMyStationsReqDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqDto dtos.LeaveMyStationsReqDto

		reqDto.Header.UserAgent = ctx.GetHeader("User-Agent")

		userId, exception := contexts.GetAndConvertContextFieldToUUID(ctx, types.ContextFieldName_User_Id)
		if exception != nil {
			exception.Log().SafelyAbortAndResponseWithJSON(ctx)
			return
		}
		reqDto.ContextFields.UserId = *userId
		if err := ctx.ShouldBindJSON(&reqDto.Body); err != nil {
			exceptions.Station.InvalidDto().WithOrigin(err).SafelyAbortAndResponseWithJSON(ctx)
			return
		}

		controllerFunc(ctx, &reqDto)
	}
}
