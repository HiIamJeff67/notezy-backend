package binders

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	stationsdto "github.com/HiIamJeff67/notezy-backend/contracts/api/v1/stations"
	exceptions "github.com/HiIamJeff67/notezy-backend/internal/exceptions"
	responsewriter "github.com/HiIamJeff67/notezy-backend/internal/gateway/responsewriter"
	apitransport "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/api"
)

type StationBinderInterface interface {
	BindGetMyStationById(controllerFunc apitransport.ControllerFunc[*stationsdto.GetMyStationByIdRequestDto]) gin.HandlerFunc
	BindGetAllMyStations(controllerFunc apitransport.ControllerFunc[*stationsdto.GetAllMyStationsRequestDto]) gin.HandlerFunc
	BindCreateStation(controllerFunc apitransport.ControllerFunc[*stationsdto.CreateStationRequestDto]) gin.HandlerFunc
	BindCreateStations(controllerFunc apitransport.ControllerFunc[*stationsdto.CreateStationsRequestDto]) gin.HandlerFunc
	BindUpdateMyStationById(controllerFunc apitransport.ControllerFunc[*stationsdto.UpdateMyStationByIdRequestDto]) gin.HandlerFunc
	BindUpdateMyStationsByIds(controllerFunc apitransport.ControllerFunc[*stationsdto.UpdateMyStationsByIdsRequestDto]) gin.HandlerFunc
	BindRestoreMyStationById(controllerFunc apitransport.ControllerFunc[*stationsdto.RestoreMyStationByIdRequestDto]) gin.HandlerFunc
	BindRestoreMyStationsByIds(controllerFunc apitransport.ControllerFunc[*stationsdto.RestoreMyStationsByIdsRequestDto]) gin.HandlerFunc
	BindDeleteMyStationById(controllerFunc apitransport.ControllerFunc[*stationsdto.DeleteMyStationByIdRequestDto]) gin.HandlerFunc
	BindDeleteMyStationsByIds(controllerFunc apitransport.ControllerFunc[*stationsdto.DeleteMyStationsByIdsRequestDto]) gin.HandlerFunc
	BindHardDeleteMyStationById(controllerFunc apitransport.ControllerFunc[*stationsdto.HardDeleteMyStationByIdRequestDto]) gin.HandlerFunc
	BindHardDeleteMyStationsByIds(controllerFunc apitransport.ControllerFunc[*stationsdto.HardDeleteMyStationsByIdsRequestDto]) gin.HandlerFunc

	/* ============================== Visualization Methods ============================== */
	BindVisualizeMyTotalCount(controllerFunc apitransport.ControllerFunc[*stationsdto.VisualizeMyTotalCountRequestDto]) gin.HandlerFunc

	/* ============================== Station Permission Methods ============================== */
	BindGetMyStationPermission(controllerFunc apitransport.ControllerFunc[*stationsdto.GetMyStationPermissionRequestDto]) gin.HandlerFunc
	BindCreateMyStationPermission(controllerFunc apitransport.ControllerFunc[*stationsdto.CreateMyStationPermissionRequestDto]) gin.HandlerFunc
	BindUpsertMyStationPermission(controllerFunc apitransport.ControllerFunc[*stationsdto.UpsertMyStationPermissionRequestDto]) gin.HandlerFunc
	BindUpsertMyStationPermissions(controllerFunc apitransport.ControllerFunc[*stationsdto.UpsertMyStationPermissionsRequestDto]) gin.HandlerFunc
	BindUpdateMyStationPermission(controllerFunc apitransport.ControllerFunc[*stationsdto.UpdateMyStationPermissionRequestDto]) gin.HandlerFunc
	BindTransferMyStationOwnership(controllerFunc apitransport.ControllerFunc[*stationsdto.TransferMyStationOwnershipRequestDto]) gin.HandlerFunc
	BindDeleteMyStationPermission(controllerFunc apitransport.ControllerFunc[*stationsdto.DeleteMyStationPermissionRequestDto]) gin.HandlerFunc
	BindDeleteMyStationPermissions(controllerFunc apitransport.ControllerFunc[*stationsdto.DeleteMyStationPermissionsRequestDto]) gin.HandlerFunc
	BindLeaveMyStation(controllerFunc apitransport.ControllerFunc[*stationsdto.LeaveMyStationRequestDto]) gin.HandlerFunc
	BindLeaveMyStations(controllerFunc apitransport.ControllerFunc[*stationsdto.LeaveMyStationsRequestDto]) gin.HandlerFunc
}

type StationBinder struct{}

func NewStationBinder() StationBinderInterface {
	return &StationBinder{}
}

func (b *StationBinder) BindGetMyStationById(controllerFunc apitransport.ControllerFunc[*stationsdto.GetMyStationByIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &stationsdto.GetMyStationByIdRequestDto{}
		request.Header.UserAgent = ctx.GetHeader("User-Agent")

		isDeletedString := ctx.Query("isDeleted")
		if isDeletedString != "" {
			isDeleted, err := strconv.ParseBool(isDeletedString)
			if err != nil {
				responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Station").WithOrigin(err), ctx)
				return
			}
			request.Param.IsDeleted = &isDeleted
		}

		stationId, err := uuid.Parse(ctx.Param("stationId"))
		if err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Station").WithOrigin(err), ctx)
			return
		}
		request.Param.StationId = stationId

		controllerFunc(ctx, request)
	}
}

func (b *StationBinder) BindGetAllMyStations(controllerFunc apitransport.ControllerFunc[*stationsdto.GetAllMyStationsRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &stationsdto.GetAllMyStationsRequestDto{}
		request.Header.UserAgent = ctx.GetHeader("User-Agent")

		areDeletedString := ctx.Query("areDeleted")
		if areDeletedString != "" {
			areDeleted, err := strconv.ParseBool(areDeletedString)
			if err != nil {
				responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Station").WithOrigin(err), ctx)
				return
			}
			request.Query.AreDeleted = &areDeleted
		}

		controllerFunc(ctx, request)
	}
}

func (b *StationBinder) BindCreateStation(controllerFunc apitransport.ControllerFunc[*stationsdto.CreateStationRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &stationsdto.CreateStationRequestDto{}
		request.Header.UserAgent = ctx.GetHeader("User-Agent")
		if err := ctx.ShouldBindJSON(&request.Body); err != nil {
			exception := exceptions.InvalidDto("Station").WithOrigin(err)
			responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
			return
		}

		controllerFunc(ctx, request)
	}
}

func (b *StationBinder) BindCreateStations(controllerFunc apitransport.ControllerFunc[*stationsdto.CreateStationsRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &stationsdto.CreateStationsRequestDto{}
		request.Header.UserAgent = ctx.GetHeader("User-Agent")
		if err := ctx.ShouldBindJSON(&request.Body); err != nil {
			exception := exceptions.InvalidDto("Station").WithOrigin(err)
			responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
			return
		}

		controllerFunc(ctx, request)
	}
}

func (b *StationBinder) BindUpdateMyStationById(controllerFunc apitransport.ControllerFunc[*stationsdto.UpdateMyStationByIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &stationsdto.UpdateMyStationByIdRequestDto{}
		request.Header.UserAgent = ctx.GetHeader("User-Agent")
		if err := ctx.ShouldBindJSON(&request.Body); err != nil {
			exception := exceptions.InvalidDto("Station").WithOrigin(err)
			responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
			return
		}

		stationId, err := uuid.Parse(ctx.Param("stationId"))
		if err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Station").WithOrigin(err), ctx)
			return
		}
		request.Param.StationId = stationId

		controllerFunc(ctx, request)
	}
}

func (b *StationBinder) BindUpdateMyStationsByIds(controllerFunc apitransport.ControllerFunc[*stationsdto.UpdateMyStationsByIdsRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &stationsdto.UpdateMyStationsByIdsRequestDto{}
		request.Header.UserAgent = ctx.GetHeader("User-Agent")
		if err := ctx.ShouldBindJSON(&request.Body); err != nil {
			exception := exceptions.InvalidDto("Station").WithOrigin(err)
			responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
			return
		}

		controllerFunc(ctx, request)
	}
}

func (b *StationBinder) BindRestoreMyStationById(controllerFunc apitransport.ControllerFunc[*stationsdto.RestoreMyStationByIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &stationsdto.RestoreMyStationByIdRequestDto{}
		request.Header.UserAgent = ctx.GetHeader("User-Agent")
		if err := ctx.ShouldBindJSON(&request.Body); err != nil {
			exception := exceptions.InvalidDto("Station").WithOrigin(err)
			responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
			return
		}

		stationId, err := uuid.Parse(ctx.Param("stationId"))
		if err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Station").WithOrigin(err), ctx)
			return
		}
		request.Body.StationId = stationId

		controllerFunc(ctx, request)
	}
}

func (b *StationBinder) BindRestoreMyStationsByIds(controllerFunc apitransport.ControllerFunc[*stationsdto.RestoreMyStationsByIdsRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &stationsdto.RestoreMyStationsByIdsRequestDto{}
		request.Header.UserAgent = ctx.GetHeader("User-Agent")
		if err := ctx.ShouldBindJSON(&request.Body); err != nil {
			exception := exceptions.InvalidDto("Station").WithOrigin(err)
			responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
			return
		}

		controllerFunc(ctx, request)
	}
}

func (b *StationBinder) BindDeleteMyStationById(controllerFunc apitransport.ControllerFunc[*stationsdto.DeleteMyStationByIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &stationsdto.DeleteMyStationByIdRequestDto{}
		request.Header.UserAgent = ctx.GetHeader("User-Agent")
		if err := ctx.ShouldBindJSON(&request.Body); err != nil {
			exception := exceptions.InvalidDto("Station").WithOrigin(err)
			responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
			return
		}

		stationId, err := uuid.Parse(ctx.Param("stationId"))
		if err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Station").WithOrigin(err), ctx)
			return
		}
		request.Body.StationId = stationId

		controllerFunc(ctx, request)
	}
}

func (b *StationBinder) BindDeleteMyStationsByIds(controllerFunc apitransport.ControllerFunc[*stationsdto.DeleteMyStationsByIdsRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &stationsdto.DeleteMyStationsByIdsRequestDto{}
		request.Header.UserAgent = ctx.GetHeader("User-Agent")
		if err := ctx.ShouldBindJSON(&request.Body); err != nil {
			exception := exceptions.InvalidDto("Station").WithOrigin(err)
			responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
			return
		}

		controllerFunc(ctx, request)
	}
}

func (b *StationBinder) BindHardDeleteMyStationById(controllerFunc apitransport.ControllerFunc[*stationsdto.HardDeleteMyStationByIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &stationsdto.HardDeleteMyStationByIdRequestDto{}
		request.Header.UserAgent = ctx.GetHeader("User-Agent")
		if err := ctx.ShouldBindJSON(&request.Body); err != nil {
			exception := exceptions.InvalidDto("Station").WithOrigin(err)
			responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
			return
		}

		stationId, err := uuid.Parse(ctx.Param("stationId"))
		if err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Station").WithOrigin(err), ctx)
			return
		}
		request.Body.StationId = stationId

		controllerFunc(ctx, request)
	}
}

func (b *StationBinder) BindHardDeleteMyStationsByIds(controllerFunc apitransport.ControllerFunc[*stationsdto.HardDeleteMyStationsByIdsRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &stationsdto.HardDeleteMyStationsByIdsRequestDto{}
		request.Header.UserAgent = ctx.GetHeader("User-Agent")
		if err := ctx.ShouldBindJSON(&request.Body); err != nil {
			exception := exceptions.InvalidDto("Station").WithOrigin(err)
			responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
			return
		}

		controllerFunc(ctx, request)
	}
}

/* ============================== Visualization Methods ============================== */

func (b *StationBinder) BindVisualizeMyTotalCount(controllerFunc apitransport.ControllerFunc[*stationsdto.VisualizeMyTotalCountRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &stationsdto.VisualizeMyTotalCountRequestDto{}
		request.Header.UserAgent = ctx.GetHeader("User-Agent")
		permissionString := ctx.Query("permission")
		if permissionString == "" {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Station").WithOrigin(fmt.Errorf("permission is required")), ctx)
			return
		}
		request.Query.Permission = permissionString

		controllerFunc(ctx, request)
	}
}

/* ============================== Station Permission Methods ============================== */

func (b *StationBinder) BindGetMyStationPermission(controllerFunc apitransport.ControllerFunc[*stationsdto.GetMyStationPermissionRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &stationsdto.GetMyStationPermissionRequestDto{}
		request.Header.UserAgent = ctx.GetHeader("User-Agent")

		stationId, err := uuid.Parse(ctx.Param("stationId"))
		if err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Station").WithOrigin(err), ctx)
			return
		}
		request.Param.StationId = stationId

		userPublicId, err := uuid.Parse(ctx.Param("userPublicId"))
		if err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Station").WithOrigin(err), ctx)
			return
		}
		request.Param.UserPublicId = userPublicId

		controllerFunc(ctx, request)
	}
}

func (b *StationBinder) BindCreateMyStationPermission(controllerFunc apitransport.ControllerFunc[*stationsdto.CreateMyStationPermissionRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &stationsdto.CreateMyStationPermissionRequestDto{}
		request.Header.UserAgent = ctx.GetHeader("User-Agent")

		stationId, err := uuid.Parse(ctx.Param("stationId"))
		if err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Station").WithOrigin(err), ctx)
			return
		}
		request.Param.StationId = stationId

		userPublicId, err := uuid.Parse(ctx.Param("userPublicId"))
		if err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Station").WithOrigin(err), ctx)
			return
		}
		request.Param.UserPublicId = userPublicId

		if err := ctx.ShouldBindJSON(&request.Body); err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("Station").WithOrigin(err), ctx)
			return
		}

		controllerFunc(ctx, request)
	}
}

func (b *StationBinder) BindUpsertMyStationPermission(controllerFunc apitransport.ControllerFunc[*stationsdto.UpsertMyStationPermissionRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &stationsdto.UpsertMyStationPermissionRequestDto{}
		request.Header.UserAgent = ctx.GetHeader("User-Agent")

		stationId, err := uuid.Parse(ctx.Param("stationId"))
		if err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Station").WithOrigin(err), ctx)
			return
		}
		request.Param.StationId = stationId

		userPublicId, err := uuid.Parse(ctx.Param("userPublicId"))
		if err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Station").WithOrigin(err), ctx)
			return
		}
		request.Param.UserPublicId = userPublicId

		if err := ctx.ShouldBindJSON(&request.Body); err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("Station").WithOrigin(err), ctx)
			return
		}

		controllerFunc(ctx, request)
	}
}

func (b *StationBinder) BindUpsertMyStationPermissions(controllerFunc apitransport.ControllerFunc[*stationsdto.UpsertMyStationPermissionsRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &stationsdto.UpsertMyStationPermissionsRequestDto{}
		request.Header.UserAgent = ctx.GetHeader("User-Agent")

		stationId, err := uuid.Parse(ctx.Param("stationId"))
		if err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Station").WithOrigin(err), ctx)
			return
		}
		request.Param.StationId = stationId

		if err := ctx.ShouldBindJSON(&request.Body); err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("Station").WithOrigin(err), ctx)
			return
		}

		controllerFunc(ctx, request)
	}
}

func (b *StationBinder) BindUpdateMyStationPermission(controllerFunc apitransport.ControllerFunc[*stationsdto.UpdateMyStationPermissionRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &stationsdto.UpdateMyStationPermissionRequestDto{}
		request.Header.UserAgent = ctx.GetHeader("User-Agent")

		stationId, err := uuid.Parse(ctx.Param("stationId"))
		if err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Station").WithOrigin(err), ctx)
			return
		}
		request.Param.StationId = stationId

		userPublicId, err := uuid.Parse(ctx.Param("userPublicId"))
		if err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Station").WithOrigin(err), ctx)
			return
		}
		request.Param.UserPublicId = userPublicId

		if err := ctx.ShouldBindJSON(&request.Body); err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("Station").WithOrigin(err), ctx)
			return
		}

		controllerFunc(ctx, request)
	}
}

func (b *StationBinder) BindTransferMyStationOwnership(controllerFunc apitransport.ControllerFunc[*stationsdto.TransferMyStationOwnershipRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &stationsdto.TransferMyStationOwnershipRequestDto{}
		request.Header.UserAgent = ctx.GetHeader("User-Agent")

		stationId, err := uuid.Parse(ctx.Param("stationId"))
		if err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Station").WithOrigin(err), ctx)
			return
		}
		request.Param.StationId = stationId

		if err := ctx.ShouldBindJSON(&request.Body); err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("Station").WithOrigin(err), ctx)
			return
		}

		controllerFunc(ctx, request)
	}
}

func (b *StationBinder) BindDeleteMyStationPermission(controllerFunc apitransport.ControllerFunc[*stationsdto.DeleteMyStationPermissionRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &stationsdto.DeleteMyStationPermissionRequestDto{}
		request.Header.UserAgent = ctx.GetHeader("User-Agent")

		stationId, err := uuid.Parse(ctx.Param("stationId"))
		if err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Station").WithOrigin(err), ctx)
			return
		}
		request.Param.StationId = stationId

		userPublicId, err := uuid.Parse(ctx.Param("userPublicId"))
		if err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Station").WithOrigin(err), ctx)
			return
		}
		request.Param.UserPublicId = userPublicId

		controllerFunc(ctx, request)
	}
}

func (b *StationBinder) BindDeleteMyStationPermissions(controllerFunc apitransport.ControllerFunc[*stationsdto.DeleteMyStationPermissionsRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &stationsdto.DeleteMyStationPermissionsRequestDto{}
		request.Header.UserAgent = ctx.GetHeader("User-Agent")

		stationId, err := uuid.Parse(ctx.Param("stationId"))
		if err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Station").WithOrigin(err), ctx)
			return
		}
		request.Param.StationId = stationId

		if err := ctx.ShouldBindJSON(&request.Body); err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("Station").WithOrigin(err), ctx)
			return
		}

		controllerFunc(ctx, request)
	}
}

func (b *StationBinder) BindLeaveMyStation(controllerFunc apitransport.ControllerFunc[*stationsdto.LeaveMyStationRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &stationsdto.LeaveMyStationRequestDto{}
		request.Header.UserAgent = ctx.GetHeader("User-Agent")

		stationId, err := uuid.Parse(ctx.Param("stationId"))
		if err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Station").WithOrigin(err), ctx)
			return
		}
		request.Param.StationId = stationId

		controllerFunc(ctx, request)
	}
}

func (b *StationBinder) BindLeaveMyStations(controllerFunc apitransport.ControllerFunc[*stationsdto.LeaveMyStationsRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &stationsdto.LeaveMyStationsRequestDto{}
		request.Header.UserAgent = ctx.GetHeader("User-Agent")
		if err := ctx.ShouldBindJSON(&request.Body); err != nil {
			responsewriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("Station").WithOrigin(err), ctx)
			return
		}

		controllerFunc(ctx, request)
	}
}
