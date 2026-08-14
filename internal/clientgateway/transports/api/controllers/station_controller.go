package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	exceptionwriter "github.com/HiIamJeff67/notezy-backend/shared/util/exceptionwriter"

	apicontract "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api/stations"

	coreadapters "github.com/HiIamJeff67/notezy-backend/internal/clientgateway/transports/core/adapters"
)

type StationControllerInterface interface {
	GetMyStationById(ctx *gin.Context, request *apicontract.GetMyStationByIdRequestDto)
	GetAllMyStations(ctx *gin.Context, request *apicontract.GetAllMyStationsRequestDto)
	CreateStation(ctx *gin.Context, request *apicontract.CreateStationRequestDto)
	CreateStations(ctx *gin.Context, request *apicontract.CreateStationsRequestDto)
	UpdateMyStationById(ctx *gin.Context, request *apicontract.UpdateMyStationByIdRequestDto)
	UpdateMyStationsByIds(ctx *gin.Context, request *apicontract.UpdateMyStationsByIdsRequestDto)
	RestoreMyStationById(ctx *gin.Context, request *apicontract.RestoreMyStationByIdRequestDto)
	RestoreMyStationsByIds(ctx *gin.Context, request *apicontract.RestoreMyStationsByIdsRequestDto)
	DeleteMyStationById(ctx *gin.Context, request *apicontract.DeleteMyStationByIdRequestDto)
	DeleteMyStationsByIds(ctx *gin.Context, request *apicontract.DeleteMyStationsByIdsRequestDto)
	HardDeleteMyStationById(ctx *gin.Context, request *apicontract.HardDeleteMyStationByIdRequestDto)
	HardDeleteMyStationsByIds(ctx *gin.Context, request *apicontract.HardDeleteMyStationsByIdsRequestDto)

	/* ============================== Visualization Methods ============================== */
	VisualizeMyTotalCount(ctx *gin.Context, request *apicontract.VisualizeMyTotalCountRequestDto)

	/* ============================== Station Permission Methods ============================== */
	GetMyStationPermission(ctx *gin.Context, request *apicontract.GetMyStationPermissionRequestDto)
	CreateMyStationPermission(ctx *gin.Context, request *apicontract.CreateMyStationPermissionRequestDto)
	UpsertMyStationPermission(ctx *gin.Context, request *apicontract.UpsertMyStationPermissionRequestDto)
	UpsertMyStationPermissions(ctx *gin.Context, request *apicontract.UpsertMyStationPermissionsRequestDto)
	UpdateMyStationPermission(ctx *gin.Context, request *apicontract.UpdateMyStationPermissionRequestDto)
	TransferMyStationOwnership(ctx *gin.Context, request *apicontract.TransferMyStationOwnershipRequestDto)
	DeleteMyStationPermission(ctx *gin.Context, request *apicontract.DeleteMyStationPermissionRequestDto)
	DeleteMyStationPermissions(ctx *gin.Context, request *apicontract.DeleteMyStationPermissionsRequestDto)
	LeaveMyStation(ctx *gin.Context, request *apicontract.LeaveMyStationRequestDto)
	LeaveMyStations(ctx *gin.Context, request *apicontract.LeaveMyStationsRequestDto)
}

type StationController struct {
	coreClient *coreadapters.CoreAdapter
}

func NewStationController(coreClient *coreadapters.CoreAdapter) StationControllerInterface {
	return &StationController{
		coreClient: coreClient,
	}
}

func (c *StationController) GetMyStationById(ctx *gin.Context, request *apicontract.GetMyStationByIdRequestDto) {
	response, exception := coreadapters.CallSecurly[
		apicontract.GetMyStationByIdRequestDto,
		apicontract.GetMyStationByIdResponseDto,
	](
		ctx,
		c.coreClient,
		request,
		apicontract.GetMyStationByIdOperation,
		"/core/v1/stations/get-by-id",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *StationController) GetAllMyStations(ctx *gin.Context, request *apicontract.GetAllMyStationsRequestDto) {
	response, exception := coreadapters.CallSecurly[
		apicontract.GetAllMyStationsRequestDto,
		apicontract.GetAllMyStationsResponseDto,
	](
		ctx,
		c.coreClient,
		request,
		apicontract.GetAllMyStationsOperation,
		"/core/v1/stations/get-all",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *StationController) CreateStation(ctx *gin.Context, request *apicontract.CreateStationRequestDto) {
	response, exception := coreadapters.CallSecurly[
		apicontract.CreateStationRequestDto,
		apicontract.CreateStationResponseDto,
	](
		ctx,
		c.coreClient,
		request,
		apicontract.CreateStationOperation,
		"/core/v1/stations/create",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeCreatedClientResponse(ctx, response.Data)
}

func (c *StationController) CreateStations(ctx *gin.Context, request *apicontract.CreateStationsRequestDto) {
	response, exception := coreadapters.CallSecurly[
		apicontract.CreateStationsRequestDto,
		apicontract.CreateStationsResponseDto,
	](
		ctx,
		c.coreClient,
		request,
		apicontract.CreateStationsOperation,
		"/core/v1/stations/create-many",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeCreatedClientResponse(ctx, response.Data)
}

func (c *StationController) UpdateMyStationById(ctx *gin.Context, request *apicontract.UpdateMyStationByIdRequestDto) {
	response, exception := coreadapters.CallSecurly[
		apicontract.UpdateMyStationByIdRequestDto,
		apicontract.UpdateMyStationByIdResponseDto,
	](
		ctx,
		c.coreClient,
		request,
		apicontract.UpdateMyStationByIdOperation,
		"/core/v1/stations/update",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *StationController) UpdateMyStationsByIds(ctx *gin.Context, request *apicontract.UpdateMyStationsByIdsRequestDto) {
	response, exception := coreadapters.CallSecurly[
		apicontract.UpdateMyStationsByIdsRequestDto,
		apicontract.UpdateMyStationsByIdsResponseDto,
	](
		ctx,
		c.coreClient,
		request,
		apicontract.UpdateMyStationsByIdsOperation,
		"/core/v1/stations/update-many",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *StationController) RestoreMyStationById(ctx *gin.Context, request *apicontract.RestoreMyStationByIdRequestDto) {
	response, exception := coreadapters.CallSecurly[
		apicontract.RestoreMyStationByIdRequestDto,
		apicontract.RestoreMyStationByIdResponseDto,
	](
		ctx,
		c.coreClient,
		request,
		apicontract.RestoreMyStationByIdOperation,
		"/core/v1/stations/restore",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *StationController) RestoreMyStationsByIds(ctx *gin.Context, request *apicontract.RestoreMyStationsByIdsRequestDto) {
	response, exception := coreadapters.CallSecurly[
		apicontract.RestoreMyStationsByIdsRequestDto,
		apicontract.RestoreMyStationsByIdsResponseDto,
	](
		ctx,
		c.coreClient,
		request,
		apicontract.RestoreMyStationsByIdsOperation,
		"/core/v1/stations/restore-many",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *StationController) DeleteMyStationById(ctx *gin.Context, request *apicontract.DeleteMyStationByIdRequestDto) {
	response, exception := coreadapters.CallSecurly[
		apicontract.DeleteMyStationByIdRequestDto,
		apicontract.DeleteMyStationByIdResponseDto,
	](
		ctx,
		c.coreClient,
		request,
		apicontract.DeleteMyStationByIdOperation,
		"/core/v1/stations/delete",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *StationController) DeleteMyStationsByIds(ctx *gin.Context, request *apicontract.DeleteMyStationsByIdsRequestDto) {
	response, exception := coreadapters.CallSecurly[
		apicontract.DeleteMyStationsByIdsRequestDto,
		apicontract.DeleteMyStationsByIdsResponseDto,
	](
		ctx,
		c.coreClient,
		request,
		apicontract.DeleteMyStationsByIdsOperation,
		"/core/v1/stations/delete-many",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *StationController) HardDeleteMyStationById(ctx *gin.Context, request *apicontract.HardDeleteMyStationByIdRequestDto) {
	response, exception := coreadapters.CallSecurly[
		apicontract.HardDeleteMyStationByIdRequestDto,
		apicontract.HardDeleteMyStationByIdResponseDto,
	](
		ctx,
		c.coreClient,
		request,
		apicontract.HardDeleteMyStationByIdOperation,
		"/core/v1/stations/hard-delete",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *StationController) HardDeleteMyStationsByIds(ctx *gin.Context, request *apicontract.HardDeleteMyStationsByIdsRequestDto) {
	response, exception := coreadapters.CallSecurly[
		apicontract.HardDeleteMyStationsByIdsRequestDto,
		apicontract.HardDeleteMyStationsByIdsResponseDto,
	](
		ctx,
		c.coreClient,
		request,
		apicontract.HardDeleteMyStationsByIdsOperation,
		"/core/v1/stations/hard-delete-many",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

/* ============================== Controller Methods for Visualization ============================== */

func (c *StationController) VisualizeMyTotalCount(ctx *gin.Context, request *apicontract.VisualizeMyTotalCountRequestDto) {
	response, exception := coreadapters.CallSecurly[
		apicontract.VisualizeMyTotalCountRequestDto,
		apicontract.VisualizeMyTotalCountResponseDto,
	](
		ctx,
		c.coreClient,
		request,
		apicontract.VisualizeMyTotalCountOperation,
		"/core/v1/stations/visualizations/total-count",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

/* ============================== Controller Methods for Station Permissions ============================== */

func (c *StationController) GetMyStationPermission(ctx *gin.Context, request *apicontract.GetMyStationPermissionRequestDto) {
	response, exception := coreadapters.CallSecurly[
		apicontract.GetMyStationPermissionRequestDto,
		apicontract.GetMyStationPermissionResponseDto,
	](
		ctx,
		c.coreClient,
		request,
		apicontract.GetMyStationPermissionOperation,
		"/core/v1/stations/permissions/get",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *StationController) CreateMyStationPermission(ctx *gin.Context, request *apicontract.CreateMyStationPermissionRequestDto) {
	response, exception := coreadapters.CallSecurly[
		apicontract.CreateMyStationPermissionRequestDto,
		apicontract.CreateMyStationPermissionResponseDto,
	](
		ctx,
		c.coreClient,
		request,
		apicontract.CreateMyStationPermissionOperation,
		"/core/v1/stations/permissions/create",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeCreatedClientResponse(ctx, response.Data)
}

func (c *StationController) UpsertMyStationPermission(
	ctx *gin.Context, request *apicontract.UpsertMyStationPermissionRequestDto,
) {
	response, exception := coreadapters.CallSecurly[
		apicontract.UpsertMyStationPermissionRequestDto,
		apicontract.UpsertMyStationPermissionResponseDto,
	](
		ctx,
		c.coreClient,
		request,
		apicontract.UpsertMyStationPermissionOperation,
		"/core/v1/stations/permissions/upsert",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *StationController) UpsertMyStationPermissions(
	ctx *gin.Context, request *apicontract.UpsertMyStationPermissionsRequestDto,
) {
	response, exception := coreadapters.CallSecurly[
		apicontract.UpsertMyStationPermissionsRequestDto,
		apicontract.UpsertMyStationPermissionsResponseDto,
	](
		ctx,
		c.coreClient,
		request,
		apicontract.UpsertMyStationPermissionsOperation,
		"/core/v1/stations/permissions/upsert-many",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *StationController) UpdateMyStationPermission(ctx *gin.Context, request *apicontract.UpdateMyStationPermissionRequestDto) {
	response, exception := coreadapters.CallSecurly[
		apicontract.UpdateMyStationPermissionRequestDto,
		apicontract.UpdateMyStationPermissionResponseDto,
	](
		ctx,
		c.coreClient,
		request,
		apicontract.UpdateMyStationPermissionOperation,
		"/core/v1/stations/permissions/update",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *StationController) TransferMyStationOwnership(
	ctx *gin.Context, request *apicontract.TransferMyStationOwnershipRequestDto,
) {
	response, exception := coreadapters.CallSecurly[
		apicontract.TransferMyStationOwnershipRequestDto,
		apicontract.TransferMyStationOwnershipResponseDto,
	](
		ctx,
		c.coreClient,
		request,
		apicontract.TransferMyStationOwnershipOperation,
		"/core/v1/stations/ownership/transfer",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *StationController) DeleteMyStationPermission(
	ctx *gin.Context, request *apicontract.DeleteMyStationPermissionRequestDto,
) {
	_, exception := coreadapters.CallSecurly[
		apicontract.DeleteMyStationPermissionRequestDto,
		apicontract.DeleteMyStationPermissionResponseDto,
	](
		ctx,
		c.coreClient,
		request,
		apicontract.DeleteMyStationPermissionOperation,
		"/core/v1/stations/permissions/delete",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, apicontract.DeleteMyStationPermissionResponseDto{})
}

func (c *StationController) DeleteMyStationPermissions(
	ctx *gin.Context, request *apicontract.DeleteMyStationPermissionsRequestDto,
) {
	_, exception := coreadapters.CallSecurly[
		apicontract.DeleteMyStationPermissionsRequestDto,
		apicontract.DeleteMyStationPermissionsResponseDto,
	](
		ctx,
		c.coreClient,
		request,
		apicontract.DeleteMyStationPermissionsOperation,
		"/core/v1/stations/permissions/delete-many",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *StationController) LeaveMyStation(ctx *gin.Context, request *apicontract.LeaveMyStationRequestDto) {
	_, exception := coreadapters.CallSecurly[
		apicontract.LeaveMyStationRequestDto,
		apicontract.LeaveMyStationResponseDto,
	](
		ctx,
		c.coreClient,
		request,
		apicontract.LeaveMyStationOperation,
		"/core/v1/stations/memberships/leave",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *StationController) LeaveMyStations(ctx *gin.Context, request *apicontract.LeaveMyStationsRequestDto) {
	_, exception := coreadapters.CallSecurly[
		apicontract.LeaveMyStationsRequestDto,
		apicontract.LeaveMyStationsResponseDto,
	](
		ctx,
		c.coreClient,
		request,
		apicontract.LeaveMyStationsOperation,
		"/core/v1/stations/memberships/leave-many",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	ctx.Status(http.StatusNoContent)
}
