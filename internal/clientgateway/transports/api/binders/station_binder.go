package binders

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	exceptions "github.com/HiIamJeff67/notezy-backend/contracts/types/exceptions"

	exceptionwriter "github.com/HiIamJeff67/notezy-backend/shared/util/exceptionwriter"

	apicontract "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api/stations"

	controllers "github.com/HiIamJeff67/notezy-backend/internal/clientgateway/transports/api/controllers"
)

type StationBinderInterface interface {
	BindGetMyStationById(controllerFunc controllers.Func[*apicontract.GetMyStationByIdRequestDto]) gin.HandlerFunc
	BindGetAllMyStations(controllerFunc controllers.Func[*apicontract.GetAllMyStationsRequestDto]) gin.HandlerFunc
	BindCreateStation(controllerFunc controllers.Func[*apicontract.CreateStationRequestDto]) gin.HandlerFunc
	BindCreateStations(controllerFunc controllers.Func[*apicontract.CreateStationsRequestDto]) gin.HandlerFunc
	BindUpdateMyStationById(controllerFunc controllers.Func[*apicontract.UpdateMyStationByIdRequestDto]) gin.HandlerFunc
	BindUpdateMyStationsByIds(controllerFunc controllers.Func[*apicontract.UpdateMyStationsByIdsRequestDto]) gin.HandlerFunc
	BindRestoreMyStationById(controllerFunc controllers.Func[*apicontract.RestoreMyStationByIdRequestDto]) gin.HandlerFunc
	BindRestoreMyStationsByIds(controllerFunc controllers.Func[*apicontract.RestoreMyStationsByIdsRequestDto]) gin.HandlerFunc
	BindDeleteMyStationById(controllerFunc controllers.Func[*apicontract.DeleteMyStationByIdRequestDto]) gin.HandlerFunc
	BindDeleteMyStationsByIds(controllerFunc controllers.Func[*apicontract.DeleteMyStationsByIdsRequestDto]) gin.HandlerFunc
	BindHardDeleteMyStationById(controllerFunc controllers.Func[*apicontract.HardDeleteMyStationByIdRequestDto]) gin.HandlerFunc
	BindHardDeleteMyStationsByIds(controllerFunc controllers.Func[*apicontract.HardDeleteMyStationsByIdsRequestDto]) gin.HandlerFunc

	/* ============================== Visualization Methods ============================== */
	BindVisualizeMyTotalCount(controllerFunc controllers.Func[*apicontract.VisualizeMyTotalCountRequestDto]) gin.HandlerFunc

	/* ============================== Station Permission Methods ============================== */
	BindGetMyStationPermission(controllerFunc controllers.Func[*apicontract.GetMyStationPermissionRequestDto]) gin.HandlerFunc
	BindCreateMyStationPermission(controllerFunc controllers.Func[*apicontract.CreateMyStationPermissionRequestDto]) gin.HandlerFunc
	BindUpsertMyStationPermission(controllerFunc controllers.Func[*apicontract.UpsertMyStationPermissionRequestDto]) gin.HandlerFunc
	BindUpsertMyStationPermissions(controllerFunc controllers.Func[*apicontract.UpsertMyStationPermissionsRequestDto]) gin.HandlerFunc
	BindUpdateMyStationPermission(controllerFunc controllers.Func[*apicontract.UpdateMyStationPermissionRequestDto]) gin.HandlerFunc
	BindTransferMyStationOwnership(controllerFunc controllers.Func[*apicontract.TransferMyStationOwnershipRequestDto]) gin.HandlerFunc
	BindDeleteMyStationPermission(controllerFunc controllers.Func[*apicontract.DeleteMyStationPermissionRequestDto]) gin.HandlerFunc
	BindDeleteMyStationPermissions(controllerFunc controllers.Func[*apicontract.DeleteMyStationPermissionsRequestDto]) gin.HandlerFunc
	BindLeaveMyStation(controllerFunc controllers.Func[*apicontract.LeaveMyStationRequestDto]) gin.HandlerFunc
	BindLeaveMyStations(controllerFunc controllers.Func[*apicontract.LeaveMyStationsRequestDto]) gin.HandlerFunc
}

type StationBinder struct{}

func NewStationBinder() StationBinderInterface {
	return &StationBinder{}
}

func (b *StationBinder) BindGetMyStationById(controllerFunc controllers.Func[*apicontract.GetMyStationByIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &apicontract.GetMyStationByIdRequestDto{}
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

		stationId, err := uuid.Parse(ctx.Param("station-id"))
		if err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Station").WithOrigin(err), ctx)
			return
		}
		request.Param.StationId = stationId

		controllerFunc(ctx, request)
	}
}

func (b *StationBinder) BindGetAllMyStations(controllerFunc controllers.Func[*apicontract.GetAllMyStationsRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &apicontract.GetAllMyStationsRequestDto{}
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

func (b *StationBinder) BindCreateStation(controllerFunc controllers.Func[*apicontract.CreateStationRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &apicontract.CreateStationRequestDto{}
		request.Header.UserAgent = ctx.GetHeader("User-Agent")
		if err := ctx.ShouldBindJSON(&request.Body); err != nil {
			exception := exceptions.InvalidDto("Station").WithOrigin(err)
			exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
			return
		}

		controllerFunc(ctx, request)
	}
}

func (b *StationBinder) BindCreateStations(controllerFunc controllers.Func[*apicontract.CreateStationsRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &apicontract.CreateStationsRequestDto{}
		request.Header.UserAgent = ctx.GetHeader("User-Agent")
		if err := ctx.ShouldBindJSON(&request.Body); err != nil {
			exception := exceptions.InvalidDto("Station").WithOrigin(err)
			exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
			return
		}

		controllerFunc(ctx, request)
	}
}

func (b *StationBinder) BindUpdateMyStationById(controllerFunc controllers.Func[*apicontract.UpdateMyStationByIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &apicontract.UpdateMyStationByIdRequestDto{}
		request.Header.UserAgent = ctx.GetHeader("User-Agent")
		if err := ctx.ShouldBindJSON(&request.Body); err != nil {
			exception := exceptions.InvalidDto("Station").WithOrigin(err)
			exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
			return
		}

		stationId, err := uuid.Parse(ctx.Param("station-id"))
		if err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Station").WithOrigin(err), ctx)
			return
		}
		request.Param.StationId = stationId

		controllerFunc(ctx, request)
	}
}

func (b *StationBinder) BindUpdateMyStationsByIds(controllerFunc controllers.Func[*apicontract.UpdateMyStationsByIdsRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &apicontract.UpdateMyStationsByIdsRequestDto{}
		request.Header.UserAgent = ctx.GetHeader("User-Agent")
		if err := ctx.ShouldBindJSON(&request.Body); err != nil {
			exception := exceptions.InvalidDto("Station").WithOrigin(err)
			exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
			return
		}

		controllerFunc(ctx, request)
	}
}

func (b *StationBinder) BindRestoreMyStationById(controllerFunc controllers.Func[*apicontract.RestoreMyStationByIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &apicontract.RestoreMyStationByIdRequestDto{}
		request.Header.UserAgent = ctx.GetHeader("User-Agent")
		if err := ctx.ShouldBindJSON(&request.Body); err != nil {
			exception := exceptions.InvalidDto("Station").WithOrigin(err)
			exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
			return
		}

		stationId, err := uuid.Parse(ctx.Param("station-id"))
		if err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Station").WithOrigin(err), ctx)
			return
		}
		request.Body.StationId = stationId

		controllerFunc(ctx, request)
	}
}

func (b *StationBinder) BindRestoreMyStationsByIds(controllerFunc controllers.Func[*apicontract.RestoreMyStationsByIdsRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &apicontract.RestoreMyStationsByIdsRequestDto{}
		request.Header.UserAgent = ctx.GetHeader("User-Agent")
		if err := ctx.ShouldBindJSON(&request.Body); err != nil {
			exception := exceptions.InvalidDto("Station").WithOrigin(err)
			exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
			return
		}

		controllerFunc(ctx, request)
	}
}

func (b *StationBinder) BindDeleteMyStationById(controllerFunc controllers.Func[*apicontract.DeleteMyStationByIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &apicontract.DeleteMyStationByIdRequestDto{}
		request.Header.UserAgent = ctx.GetHeader("User-Agent")
		if err := ctx.ShouldBindJSON(&request.Body); err != nil {
			exception := exceptions.InvalidDto("Station").WithOrigin(err)
			exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
			return
		}

		stationId, err := uuid.Parse(ctx.Param("station-id"))
		if err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Station").WithOrigin(err), ctx)
			return
		}
		request.Body.StationId = stationId

		controllerFunc(ctx, request)
	}
}

func (b *StationBinder) BindDeleteMyStationsByIds(controllerFunc controllers.Func[*apicontract.DeleteMyStationsByIdsRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &apicontract.DeleteMyStationsByIdsRequestDto{}
		request.Header.UserAgent = ctx.GetHeader("User-Agent")
		if err := ctx.ShouldBindJSON(&request.Body); err != nil {
			exception := exceptions.InvalidDto("Station").WithOrigin(err)
			exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
			return
		}

		controllerFunc(ctx, request)
	}
}

func (b *StationBinder) BindHardDeleteMyStationById(controllerFunc controllers.Func[*apicontract.HardDeleteMyStationByIdRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &apicontract.HardDeleteMyStationByIdRequestDto{}
		request.Header.UserAgent = ctx.GetHeader("User-Agent")
		if err := ctx.ShouldBindJSON(&request.Body); err != nil {
			exception := exceptions.InvalidDto("Station").WithOrigin(err)
			exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
			return
		}

		stationId, err := uuid.Parse(ctx.Param("station-id"))
		if err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Station").WithOrigin(err), ctx)
			return
		}
		request.Body.StationId = stationId

		controllerFunc(ctx, request)
	}
}

func (b *StationBinder) BindHardDeleteMyStationsByIds(controllerFunc controllers.Func[*apicontract.HardDeleteMyStationsByIdsRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &apicontract.HardDeleteMyStationsByIdsRequestDto{}
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

func (b *StationBinder) BindVisualizeMyTotalCount(controllerFunc controllers.Func[*apicontract.VisualizeMyTotalCountRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &apicontract.VisualizeMyTotalCountRequestDto{}
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

func (b *StationBinder) BindGetMyStationPermission(controllerFunc controllers.Func[*apicontract.GetMyStationPermissionRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &apicontract.GetMyStationPermissionRequestDto{}
		request.Header.UserAgent = ctx.GetHeader("User-Agent")

		stationId, err := uuid.Parse(ctx.Param("station-id"))
		if err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Station").WithOrigin(err), ctx)
			return
		}
		request.Param.StationId = stationId

		userPublicId, err := uuid.Parse(ctx.Param("user-public-id"))
		if err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Station").WithOrigin(err), ctx)
			return
		}
		request.Param.UserPublicId = userPublicId

		controllerFunc(ctx, request)
	}
}

func (b *StationBinder) BindCreateMyStationPermission(controllerFunc controllers.Func[*apicontract.CreateMyStationPermissionRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &apicontract.CreateMyStationPermissionRequestDto{}
		request.Header.UserAgent = ctx.GetHeader("User-Agent")

		stationId, err := uuid.Parse(ctx.Param("station-id"))
		if err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Station").WithOrigin(err), ctx)
			return
		}
		request.Param.StationId = stationId

		userPublicId, err := uuid.Parse(ctx.Param("user-public-id"))
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

func (b *StationBinder) BindUpsertMyStationPermission(controllerFunc controllers.Func[*apicontract.UpsertMyStationPermissionRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &apicontract.UpsertMyStationPermissionRequestDto{}
		request.Header.UserAgent = ctx.GetHeader("User-Agent")

		stationId, err := uuid.Parse(ctx.Param("station-id"))
		if err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Station").WithOrigin(err), ctx)
			return
		}
		request.Param.StationId = stationId

		userPublicId, err := uuid.Parse(ctx.Param("user-public-id"))
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

func (b *StationBinder) BindUpsertMyStationPermissions(controllerFunc controllers.Func[*apicontract.UpsertMyStationPermissionsRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &apicontract.UpsertMyStationPermissionsRequestDto{}
		request.Header.UserAgent = ctx.GetHeader("User-Agent")

		stationId, err := uuid.Parse(ctx.Param("station-id"))
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

func (b *StationBinder) BindUpdateMyStationPermission(controllerFunc controllers.Func[*apicontract.UpdateMyStationPermissionRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &apicontract.UpdateMyStationPermissionRequestDto{}
		request.Header.UserAgent = ctx.GetHeader("User-Agent")

		stationId, err := uuid.Parse(ctx.Param("station-id"))
		if err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Station").WithOrigin(err), ctx)
			return
		}
		request.Param.StationId = stationId

		userPublicId, err := uuid.Parse(ctx.Param("user-public-id"))
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

func (b *StationBinder) BindTransferMyStationOwnership(controllerFunc controllers.Func[*apicontract.TransferMyStationOwnershipRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &apicontract.TransferMyStationOwnershipRequestDto{}
		request.Header.UserAgent = ctx.GetHeader("User-Agent")

		stationId, err := uuid.Parse(ctx.Param("station-id"))
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

func (b *StationBinder) BindDeleteMyStationPermission(controllerFunc controllers.Func[*apicontract.DeleteMyStationPermissionRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &apicontract.DeleteMyStationPermissionRequestDto{}
		request.Header.UserAgent = ctx.GetHeader("User-Agent")

		stationId, err := uuid.Parse(ctx.Param("station-id"))
		if err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Station").WithOrigin(err), ctx)
			return
		}
		request.Param.StationId = stationId

		userPublicId, err := uuid.Parse(ctx.Param("user-public-id"))
		if err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Station").WithOrigin(err), ctx)
			return
		}
		request.Param.UserPublicId = userPublicId

		controllerFunc(ctx, request)
	}
}

func (b *StationBinder) BindDeleteMyStationPermissions(controllerFunc controllers.Func[*apicontract.DeleteMyStationPermissionsRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &apicontract.DeleteMyStationPermissionsRequestDto{}
		request.Header.UserAgent = ctx.GetHeader("User-Agent")

		stationId, err := uuid.Parse(ctx.Param("station-id"))
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

func (b *StationBinder) BindLeaveMyStation(controllerFunc controllers.Func[*apicontract.LeaveMyStationRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &apicontract.LeaveMyStationRequestDto{}
		request.Header.UserAgent = ctx.GetHeader("User-Agent")

		stationId, err := uuid.Parse(ctx.Param("station-id"))
		if err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidInput("Station").WithOrigin(err), ctx)
			return
		}
		request.Param.StationId = stationId

		controllerFunc(ctx, request)
	}
}

func (b *StationBinder) BindLeaveMyStations(controllerFunc controllers.Func[*apicontract.LeaveMyStationsRequestDto]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := &apicontract.LeaveMyStationsRequestDto{}
		request.Header.UserAgent = ctx.GetHeader("User-Agent")
		if err := ctx.ShouldBindJSON(&request.Body); err != nil {
			exceptionwriter.SafelyAbortAndResponseWithJSON(exceptions.InvalidDto("Station").WithOrigin(err), ctx)
			return
		}

		controllerFunc(ctx, request)
	}
}
