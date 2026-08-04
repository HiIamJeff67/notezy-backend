package binders

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	stationsdto "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api/stations"
	exceptions "github.com/HiIamJeff67/notezy-backend/internal/exceptions"
	controllers "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/api/controllers"
	exceptionwriter "github.com/HiIamJeff67/notezy-backend/shared/exceptionwriter"
)

type StationBinderInterface interface {
	BindGetMyStationById(controllerFunc controllers.Func[*stationsdto.GetMyStationByIdRequestDto]) gin.HandlerFunc
	BindGetAllMyStations(controllerFunc controllers.Func[*stationsdto.GetAllMyStationsRequestDto]) gin.HandlerFunc
	BindCreateStation(controllerFunc controllers.Func[*stationsdto.CreateStationRequestDto]) gin.HandlerFunc
	BindCreateStations(controllerFunc controllers.Func[*stationsdto.CreateStationsRequestDto]) gin.HandlerFunc
	BindUpdateMyStationById(controllerFunc controllers.Func[*stationsdto.UpdateMyStationByIdRequestDto]) gin.HandlerFunc
	BindUpdateMyStationsByIds(controllerFunc controllers.Func[*stationsdto.UpdateMyStationsByIdsRequestDto]) gin.HandlerFunc
	BindRestoreMyStationById(controllerFunc controllers.Func[*stationsdto.RestoreMyStationByIdRequestDto]) gin.HandlerFunc
	BindRestoreMyStationsByIds(controllerFunc controllers.Func[*stationsdto.RestoreMyStationsByIdsRequestDto]) gin.HandlerFunc
	BindDeleteMyStationById(controllerFunc controllers.Func[*stationsdto.DeleteMyStationByIdRequestDto]) gin.HandlerFunc
	BindDeleteMyStationsByIds(controllerFunc controllers.Func[*stationsdto.DeleteMyStationsByIdsRequestDto]) gin.HandlerFunc
	BindHardDeleteMyStationById(controllerFunc controllers.Func[*stationsdto.HardDeleteMyStationByIdRequestDto]) gin.HandlerFunc
	BindHardDeleteMyStationsByIds(controllerFunc controllers.Func[*stationsdto.HardDeleteMyStationsByIdsRequestDto]) gin.HandlerFunc

	/* ============================== Visualization Methods ============================== */
	BindVisualizeMyTotalCount(controllerFunc controllers.Func[*stationsdto.VisualizeMyTotalCountRequestDto]) gin.HandlerFunc

	/* ============================== Station Permission Methods ============================== */
	BindGetMyStationPermission(controllerFunc controllers.Func[*stationsdto.GetMyStationPermissionRequestDto]) gin.HandlerFunc
	BindCreateMyStationPermission(controllerFunc controllers.Func[*stationsdto.CreateMyStationPermissionRequestDto]) gin.HandlerFunc
	BindUpsertMyStationPermission(controllerFunc controllers.Func[*stationsdto.UpsertMyStationPermissionRequestDto]) gin.HandlerFunc
	BindUpsertMyStationPermissions(controllerFunc controllers.Func[*stationsdto.UpsertMyStationPermissionsRequestDto]) gin.HandlerFunc
	BindUpdateMyStationPermission(controllerFunc controllers.Func[*stationsdto.UpdateMyStationPermissionRequestDto]) gin.HandlerFunc
	BindTransferMyStationOwnership(controllerFunc controllers.Func[*stationsdto.TransferMyStationOwnershipRequestDto]) gin.HandlerFunc
	BindDeleteMyStationPermission(controllerFunc controllers.Func[*stationsdto.DeleteMyStationPermissionRequestDto]) gin.HandlerFunc
	BindDeleteMyStationPermissions(controllerFunc controllers.Func[*stationsdto.DeleteMyStationPermissionsRequestDto]) gin.HandlerFunc
	BindLeaveMyStation(controllerFunc controllers.Func[*stationsdto.LeaveMyStationRequestDto]) gin.HandlerFunc
	BindLeaveMyStations(controllerFunc controllers.Func[*stationsdto.LeaveMyStationsRequestDto]) gin.HandlerFunc
}

type StationBinder struct{}

func NewStationBinder() StationBinderInterface {
	return &StationBinder{}
}

func (b *StationBinder) BindGetMyStationById(controllerFunc controllers.Func[*stationsdto.GetMyStationByIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &stationsdto.GetMyStationByIdRequestDto{}
		request.Header.UserAgent = ctx.GetHeader("User-Agent")

		isDeletedString := ctx.Query("isDeleted")
		if isDeletedString != "" {
			isDeleted, err := strconv.ParseBool(isDeletedString)
			if err != nil {
				exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Station").WithOrigin(err), ctx)
				return
			}
			request.Param.IsDeleted = &isDeleted
		}

		stationId, err := uuid.Parse(ctx.Param("stationId"))
		if err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Station").WithOrigin(err), ctx)
			return
		}
		request.Param.StationId = stationId

		controllerFunc(ctx, request)
	}
}

func (b *StationBinder) BindGetAllMyStations(controllerFunc controllers.Func[*stationsdto.GetAllMyStationsRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &stationsdto.GetAllMyStationsRequestDto{}
		request.Header.UserAgent = ctx.GetHeader("User-Agent")

		areDeletedString := ctx.Query("areDeleted")
		if areDeletedString != "" {
			areDeleted, err := strconv.ParseBool(areDeletedString)
			if err != nil {
				exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Station").WithOrigin(err), ctx)
				return
			}
			request.Query.AreDeleted = &areDeleted
		}

		controllerFunc(ctx, request)
	}
}

func (b *StationBinder) BindCreateStation(controllerFunc controllers.Func[*stationsdto.CreateStationRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &stationsdto.CreateStationRequestDto{}
		request.Header.UserAgent = ctx.GetHeader("User-Agent")
		if err := ctx.ShouldBindJSON(&request.Body); err != nil {
			exception := exceptions.InvalidDto("Station").WithOrigin(err)
			exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
			return
		}

		controllerFunc(ctx, request)
	}
}

func (b *StationBinder) BindCreateStations(controllerFunc controllers.Func[*stationsdto.CreateStationsRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &stationsdto.CreateStationsRequestDto{}
		request.Header.UserAgent = ctx.GetHeader("User-Agent")
		if err := ctx.ShouldBindJSON(&request.Body); err != nil {
			exception := exceptions.InvalidDto("Station").WithOrigin(err)
			exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
			return
		}

		controllerFunc(ctx, request)
	}
}

func (b *StationBinder) BindUpdateMyStationById(controllerFunc controllers.Func[*stationsdto.UpdateMyStationByIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &stationsdto.UpdateMyStationByIdRequestDto{}
		request.Header.UserAgent = ctx.GetHeader("User-Agent")
		if err := ctx.ShouldBindJSON(&request.Body); err != nil {
			exception := exceptions.InvalidDto("Station").WithOrigin(err)
			exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
			return
		}

		stationId, err := uuid.Parse(ctx.Param("stationId"))
		if err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Station").WithOrigin(err), ctx)
			return
		}
		request.Param.StationId = stationId

		controllerFunc(ctx, request)
	}
}

func (b *StationBinder) BindUpdateMyStationsByIds(controllerFunc controllers.Func[*stationsdto.UpdateMyStationsByIdsRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &stationsdto.UpdateMyStationsByIdsRequestDto{}
		request.Header.UserAgent = ctx.GetHeader("User-Agent")
		if err := ctx.ShouldBindJSON(&request.Body); err != nil {
			exception := exceptions.InvalidDto("Station").WithOrigin(err)
			exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
			return
		}

		controllerFunc(ctx, request)
	}
}

func (b *StationBinder) BindRestoreMyStationById(controllerFunc controllers.Func[*stationsdto.RestoreMyStationByIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &stationsdto.RestoreMyStationByIdRequestDto{}
		request.Header.UserAgent = ctx.GetHeader("User-Agent")
		if err := ctx.ShouldBindJSON(&request.Body); err != nil {
			exception := exceptions.InvalidDto("Station").WithOrigin(err)
			exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
			return
		}

		stationId, err := uuid.Parse(ctx.Param("stationId"))
		if err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Station").WithOrigin(err), ctx)
			return
		}
		request.Body.StationId = stationId

		controllerFunc(ctx, request)
	}
}

func (b *StationBinder) BindRestoreMyStationsByIds(controllerFunc controllers.Func[*stationsdto.RestoreMyStationsByIdsRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &stationsdto.RestoreMyStationsByIdsRequestDto{}
		request.Header.UserAgent = ctx.GetHeader("User-Agent")
		if err := ctx.ShouldBindJSON(&request.Body); err != nil {
			exception := exceptions.InvalidDto("Station").WithOrigin(err)
			exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
			return
		}

		controllerFunc(ctx, request)
	}
}

func (b *StationBinder) BindDeleteMyStationById(controllerFunc controllers.Func[*stationsdto.DeleteMyStationByIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &stationsdto.DeleteMyStationByIdRequestDto{}
		request.Header.UserAgent = ctx.GetHeader("User-Agent")
		if err := ctx.ShouldBindJSON(&request.Body); err != nil {
			exception := exceptions.InvalidDto("Station").WithOrigin(err)
			exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
			return
		}

		stationId, err := uuid.Parse(ctx.Param("stationId"))
		if err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Station").WithOrigin(err), ctx)
			return
		}
		request.Body.StationId = stationId

		controllerFunc(ctx, request)
	}
}

func (b *StationBinder) BindDeleteMyStationsByIds(controllerFunc controllers.Func[*stationsdto.DeleteMyStationsByIdsRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &stationsdto.DeleteMyStationsByIdsRequestDto{}
		request.Header.UserAgent = ctx.GetHeader("User-Agent")
		if err := ctx.ShouldBindJSON(&request.Body); err != nil {
			exception := exceptions.InvalidDto("Station").WithOrigin(err)
			exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
			return
		}

		controllerFunc(ctx, request)
	}
}

func (b *StationBinder) BindHardDeleteMyStationById(controllerFunc controllers.Func[*stationsdto.HardDeleteMyStationByIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &stationsdto.HardDeleteMyStationByIdRequestDto{}
		request.Header.UserAgent = ctx.GetHeader("User-Agent")
		if err := ctx.ShouldBindJSON(&request.Body); err != nil {
			exception := exceptions.InvalidDto("Station").WithOrigin(err)
			exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
			return
		}

		stationId, err := uuid.Parse(ctx.Param("stationId"))
		if err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Station").WithOrigin(err), ctx)
			return
		}
		request.Body.StationId = stationId

		controllerFunc(ctx, request)
	}
}

func (b *StationBinder) BindHardDeleteMyStationsByIds(controllerFunc controllers.Func[*stationsdto.HardDeleteMyStationsByIdsRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &stationsdto.HardDeleteMyStationsByIdsRequestDto{}
		request.Header.UserAgent = ctx.GetHeader("User-Agent")
		if err := ctx.ShouldBindJSON(&request.Body); err != nil {
			exception := exceptions.InvalidDto("Station").WithOrigin(err)
			exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
			return
		}

		controllerFunc(ctx, request)
	}
}

/* ============================== Visualization Methods ============================== */

func (b *StationBinder) BindVisualizeMyTotalCount(controllerFunc controllers.Func[*stationsdto.VisualizeMyTotalCountRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &stationsdto.VisualizeMyTotalCountRequestDto{}
		request.Header.UserAgent = ctx.GetHeader("User-Agent")
		permissionString := ctx.Query("permission")
		if permissionString == "" {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Station").WithOrigin(fmt.Errorf("permission is required")), ctx)
			return
		}
		request.Query.Permission = permissionString

		controllerFunc(ctx, request)
	}
}

/* ============================== Station Permission Methods ============================== */

func (b *StationBinder) BindGetMyStationPermission(controllerFunc controllers.Func[*stationsdto.GetMyStationPermissionRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &stationsdto.GetMyStationPermissionRequestDto{}
		request.Header.UserAgent = ctx.GetHeader("User-Agent")

		stationId, err := uuid.Parse(ctx.Param("stationId"))
		if err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Station").WithOrigin(err), ctx)
			return
		}
		request.Param.StationId = stationId

		userPublicId, err := uuid.Parse(ctx.Param("userPublicId"))
		if err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Station").WithOrigin(err), ctx)
			return
		}
		request.Param.UserPublicId = userPublicId

		controllerFunc(ctx, request)
	}
}

func (b *StationBinder) BindCreateMyStationPermission(controllerFunc controllers.Func[*stationsdto.CreateMyStationPermissionRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &stationsdto.CreateMyStationPermissionRequestDto{}
		request.Header.UserAgent = ctx.GetHeader("User-Agent")

		stationId, err := uuid.Parse(ctx.Param("stationId"))
		if err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Station").WithOrigin(err), ctx)
			return
		}
		request.Param.StationId = stationId

		userPublicId, err := uuid.Parse(ctx.Param("userPublicId"))
		if err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Station").WithOrigin(err), ctx)
			return
		}
		request.Param.UserPublicId = userPublicId

		if err := ctx.ShouldBindJSON(&request.Body); err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("Station").WithOrigin(err), ctx)
			return
		}

		controllerFunc(ctx, request)
	}
}

func (b *StationBinder) BindUpsertMyStationPermission(controllerFunc controllers.Func[*stationsdto.UpsertMyStationPermissionRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &stationsdto.UpsertMyStationPermissionRequestDto{}
		request.Header.UserAgent = ctx.GetHeader("User-Agent")

		stationId, err := uuid.Parse(ctx.Param("stationId"))
		if err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Station").WithOrigin(err), ctx)
			return
		}
		request.Param.StationId = stationId

		userPublicId, err := uuid.Parse(ctx.Param("userPublicId"))
		if err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Station").WithOrigin(err), ctx)
			return
		}
		request.Param.UserPublicId = userPublicId

		if err := ctx.ShouldBindJSON(&request.Body); err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("Station").WithOrigin(err), ctx)
			return
		}

		controllerFunc(ctx, request)
	}
}

func (b *StationBinder) BindUpsertMyStationPermissions(controllerFunc controllers.Func[*stationsdto.UpsertMyStationPermissionsRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &stationsdto.UpsertMyStationPermissionsRequestDto{}
		request.Header.UserAgent = ctx.GetHeader("User-Agent")

		stationId, err := uuid.Parse(ctx.Param("stationId"))
		if err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Station").WithOrigin(err), ctx)
			return
		}
		request.Param.StationId = stationId

		if err := ctx.ShouldBindJSON(&request.Body); err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("Station").WithOrigin(err), ctx)
			return
		}

		controllerFunc(ctx, request)
	}
}

func (b *StationBinder) BindUpdateMyStationPermission(controllerFunc controllers.Func[*stationsdto.UpdateMyStationPermissionRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &stationsdto.UpdateMyStationPermissionRequestDto{}
		request.Header.UserAgent = ctx.GetHeader("User-Agent")

		stationId, err := uuid.Parse(ctx.Param("stationId"))
		if err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Station").WithOrigin(err), ctx)
			return
		}
		request.Param.StationId = stationId

		userPublicId, err := uuid.Parse(ctx.Param("userPublicId"))
		if err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Station").WithOrigin(err), ctx)
			return
		}
		request.Param.UserPublicId = userPublicId

		if err := ctx.ShouldBindJSON(&request.Body); err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("Station").WithOrigin(err), ctx)
			return
		}

		controllerFunc(ctx, request)
	}
}

func (b *StationBinder) BindTransferMyStationOwnership(controllerFunc controllers.Func[*stationsdto.TransferMyStationOwnershipRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &stationsdto.TransferMyStationOwnershipRequestDto{}
		request.Header.UserAgent = ctx.GetHeader("User-Agent")

		stationId, err := uuid.Parse(ctx.Param("stationId"))
		if err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Station").WithOrigin(err), ctx)
			return
		}
		request.Param.StationId = stationId

		if err := ctx.ShouldBindJSON(&request.Body); err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("Station").WithOrigin(err), ctx)
			return
		}

		controllerFunc(ctx, request)
	}
}

func (b *StationBinder) BindDeleteMyStationPermission(controllerFunc controllers.Func[*stationsdto.DeleteMyStationPermissionRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &stationsdto.DeleteMyStationPermissionRequestDto{}
		request.Header.UserAgent = ctx.GetHeader("User-Agent")

		stationId, err := uuid.Parse(ctx.Param("stationId"))
		if err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Station").WithOrigin(err), ctx)
			return
		}
		request.Param.StationId = stationId

		userPublicId, err := uuid.Parse(ctx.Param("userPublicId"))
		if err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Station").WithOrigin(err), ctx)
			return
		}
		request.Param.UserPublicId = userPublicId

		controllerFunc(ctx, request)
	}
}

func (b *StationBinder) BindDeleteMyStationPermissions(controllerFunc controllers.Func[*stationsdto.DeleteMyStationPermissionsRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &stationsdto.DeleteMyStationPermissionsRequestDto{}
		request.Header.UserAgent = ctx.GetHeader("User-Agent")

		stationId, err := uuid.Parse(ctx.Param("stationId"))
		if err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Station").WithOrigin(err), ctx)
			return
		}
		request.Param.StationId = stationId

		if err := ctx.ShouldBindJSON(&request.Body); err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("Station").WithOrigin(err), ctx)
			return
		}

		controllerFunc(ctx, request)
	}
}

func (b *StationBinder) BindLeaveMyStation(controllerFunc controllers.Func[*stationsdto.LeaveMyStationRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &stationsdto.LeaveMyStationRequestDto{}
		request.Header.UserAgent = ctx.GetHeader("User-Agent")

		stationId, err := uuid.Parse(ctx.Param("stationId"))
		if err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Station").WithOrigin(err), ctx)
			return
		}
		request.Param.StationId = stationId

		controllerFunc(ctx, request)
	}
}

func (b *StationBinder) BindLeaveMyStations(controllerFunc controllers.Func[*stationsdto.LeaveMyStationsRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &stationsdto.LeaveMyStationsRequestDto{}
		request.Header.UserAgent = ctx.GetHeader("User-Agent")
		if err := ctx.ShouldBindJSON(&request.Body); err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("Station").WithOrigin(err), ctx)
			return
		}

		controllerFunc(ctx, request)
	}
}
