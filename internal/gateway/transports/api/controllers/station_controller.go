package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	stationsdto "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api/stations"
	coreadapters "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/core/adapters"
	exceptionwriter "github.com/HiIamJeff67/notezy-backend/shared/exceptionwriter"
)

type StationControllerInterface interface {
	GetMyStationById(ctx *gin.Context, request *stationsdto.GetMyStationByIdRequestDto)
	GetAllMyStations(ctx *gin.Context, request *stationsdto.GetAllMyStationsRequestDto)
	CreateStation(ctx *gin.Context, request *stationsdto.CreateStationRequestDto)
	CreateStations(ctx *gin.Context, request *stationsdto.CreateStationsRequestDto)
	UpdateMyStationById(ctx *gin.Context, request *stationsdto.UpdateMyStationByIdRequestDto)
	UpdateMyStationsByIds(ctx *gin.Context, request *stationsdto.UpdateMyStationsByIdsRequestDto)
	RestoreMyStationById(ctx *gin.Context, request *stationsdto.RestoreMyStationByIdRequestDto)
	RestoreMyStationsByIds(ctx *gin.Context, request *stationsdto.RestoreMyStationsByIdsRequestDto)
	DeleteMyStationById(ctx *gin.Context, request *stationsdto.DeleteMyStationByIdRequestDto)
	DeleteMyStationsByIds(ctx *gin.Context, request *stationsdto.DeleteMyStationsByIdsRequestDto)
	HardDeleteMyStationById(ctx *gin.Context, request *stationsdto.HardDeleteMyStationByIdRequestDto)
	HardDeleteMyStationsByIds(ctx *gin.Context, request *stationsdto.HardDeleteMyStationsByIdsRequestDto)

	/* ============================== Visualization Methods ============================== */
	VisualizeMyTotalCount(ctx *gin.Context, request *stationsdto.VisualizeMyTotalCountRequestDto)

	/* ============================== Station Permission Methods ============================== */
	GetMyStationPermission(ctx *gin.Context, request *stationsdto.GetMyStationPermissionRequestDto)
	CreateMyStationPermission(ctx *gin.Context, request *stationsdto.CreateMyStationPermissionRequestDto)
	UpsertMyStationPermission(ctx *gin.Context, request *stationsdto.UpsertMyStationPermissionRequestDto)
	UpsertMyStationPermissions(ctx *gin.Context, request *stationsdto.UpsertMyStationPermissionsRequestDto)
	UpdateMyStationPermission(ctx *gin.Context, request *stationsdto.UpdateMyStationPermissionRequestDto)
	TransferMyStationOwnership(ctx *gin.Context, request *stationsdto.TransferMyStationOwnershipRequestDto)
	DeleteMyStationPermission(ctx *gin.Context, request *stationsdto.DeleteMyStationPermissionRequestDto)
	DeleteMyStationPermissions(ctx *gin.Context, request *stationsdto.DeleteMyStationPermissionsRequestDto)
	LeaveMyStation(ctx *gin.Context, request *stationsdto.LeaveMyStationRequestDto)
	LeaveMyStations(ctx *gin.Context, request *stationsdto.LeaveMyStationsRequestDto)
}

type StationController struct {
	coreClient *coreadapters.CoreClient
}

func NewStationController(coreClient *coreadapters.CoreClient) StationControllerInterface {
	return &StationController{
		coreClient: coreClient,
	}
}

func (c *StationController) GetMyStationById(ctx *gin.Context, request *stationsdto.GetMyStationByIdRequestDto) {
	response, exception := coreadapters.CallSecurly[
		stationsdto.GetMyStationByIdRequestDto,
		stationsdto.GetMyStationByIdResponseDto,
	](
		ctx,
		c.coreClient,
		request,
		stationsdto.GetMyStationByIdOperation,
		"/core/v1/stations/get-by-id",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      response.Data,
		"exception": nil,
	})
}

func (c *StationController) GetAllMyStations(ctx *gin.Context, request *stationsdto.GetAllMyStationsRequestDto) {
	response, exception := coreadapters.CallSecurly[
		stationsdto.GetAllMyStationsRequestDto,
		stationsdto.GetAllMyStationsResponseDto,
	](
		ctx,
		c.coreClient,
		request,
		stationsdto.GetAllMyStationsOperation,
		"/core/v1/stations/get-all",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      response.Data,
		"exception": nil,
	})
}

func (c *StationController) CreateStation(ctx *gin.Context, request *stationsdto.CreateStationRequestDto) {
	response, exception := coreadapters.CallSecurly[
		stationsdto.CreateStationRequestDto,
		stationsdto.CreateStationResponseDto,
	](
		ctx,
		c.coreClient,
		request,
		stationsdto.CreateStationOperation,
		"/core/v1/stations/create",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      response.Data,
		"exception": nil,
	})
}

func (c *StationController) CreateStations(ctx *gin.Context, request *stationsdto.CreateStationsRequestDto) {
	response, exception := coreadapters.CallSecurly[
		stationsdto.CreateStationsRequestDto,
		stationsdto.CreateStationsResponseDto,
	](
		ctx,
		c.coreClient,
		request,
		stationsdto.CreateStationsOperation,
		"/core/v1/stations/create-many",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      response.Data,
		"exception": nil,
	})
}

func (c *StationController) UpdateMyStationById(ctx *gin.Context, request *stationsdto.UpdateMyStationByIdRequestDto) {
	response, exception := coreadapters.CallSecurly[
		stationsdto.UpdateMyStationByIdRequestDto,
		stationsdto.UpdateMyStationByIdResponseDto,
	](
		ctx,
		c.coreClient,
		request,
		stationsdto.UpdateMyStationByIdOperation,
		"/core/v1/stations/update",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      response.Data,
		"exception": nil,
	})
}

func (c *StationController) UpdateMyStationsByIds(ctx *gin.Context, request *stationsdto.UpdateMyStationsByIdsRequestDto) {
	response, exception := coreadapters.CallSecurly[
		stationsdto.UpdateMyStationsByIdsRequestDto,
		stationsdto.UpdateMyStationsByIdsResponseDto,
	](
		ctx,
		c.coreClient,
		request,
		stationsdto.UpdateMyStationsByIdsOperation,
		"/core/v1/stations/update-many",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      response.Data,
		"exception": nil,
	})
}

func (c *StationController) RestoreMyStationById(ctx *gin.Context, request *stationsdto.RestoreMyStationByIdRequestDto) {
	response, exception := coreadapters.CallSecurly[
		stationsdto.RestoreMyStationByIdRequestDto,
		stationsdto.RestoreMyStationByIdResponseDto,
	](
		ctx,
		c.coreClient,
		request,
		stationsdto.RestoreMyStationByIdOperation,
		"/core/v1/stations/restore",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      response.Data,
		"exception": nil,
	})
}

func (c *StationController) RestoreMyStationsByIds(ctx *gin.Context, request *stationsdto.RestoreMyStationsByIdsRequestDto) {
	response, exception := coreadapters.CallSecurly[
		stationsdto.RestoreMyStationsByIdsRequestDto,
		stationsdto.RestoreMyStationsByIdsResponseDto,
	](
		ctx,
		c.coreClient,
		request,
		stationsdto.RestoreMyStationsByIdsOperation,
		"/core/v1/stations/restore-many",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      response.Data,
		"exception": nil,
	})
}

func (c *StationController) DeleteMyStationById(ctx *gin.Context, request *stationsdto.DeleteMyStationByIdRequestDto) {
	response, exception := coreadapters.CallSecurly[
		stationsdto.DeleteMyStationByIdRequestDto,
		stationsdto.DeleteMyStationByIdResponseDto,
	](
		ctx,
		c.coreClient,
		request,
		stationsdto.DeleteMyStationByIdOperation,
		"/core/v1/stations/delete",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      response.Data,
		"exception": nil,
	})
}

func (c *StationController) DeleteMyStationsByIds(ctx *gin.Context, request *stationsdto.DeleteMyStationsByIdsRequestDto) {
	response, exception := coreadapters.CallSecurly[
		stationsdto.DeleteMyStationsByIdsRequestDto,
		stationsdto.DeleteMyStationsByIdsResponseDto,
	](
		ctx,
		c.coreClient,
		request,
		stationsdto.DeleteMyStationsByIdsOperation,
		"/core/v1/stations/delete-many",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      response.Data,
		"exception": nil,
	})
}

func (c *StationController) HardDeleteMyStationById(ctx *gin.Context, request *stationsdto.HardDeleteMyStationByIdRequestDto) {
	response, exception := coreadapters.CallSecurly[
		stationsdto.HardDeleteMyStationByIdRequestDto,
		stationsdto.HardDeleteMyStationByIdResponseDto,
	](
		ctx,
		c.coreClient,
		request,
		stationsdto.HardDeleteMyStationByIdOperation,
		"/core/v1/stations/hard-delete",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      response.Data,
		"exception": nil,
	})
}

func (c *StationController) HardDeleteMyStationsByIds(ctx *gin.Context, request *stationsdto.HardDeleteMyStationsByIdsRequestDto) {
	response, exception := coreadapters.CallSecurly[
		stationsdto.HardDeleteMyStationsByIdsRequestDto,
		stationsdto.HardDeleteMyStationsByIdsResponseDto,
	](
		ctx,
		c.coreClient,
		request,
		stationsdto.HardDeleteMyStationsByIdsOperation,
		"/core/v1/stations/hard-delete-many",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      response.Data,
		"exception": nil,
	})
}

/* ============================== Controller Methods for Visualization ============================== */

func (c *StationController) VisualizeMyTotalCount(ctx *gin.Context, request *stationsdto.VisualizeMyTotalCountRequestDto) {
	response, exception := coreadapters.CallSecurly[
		stationsdto.VisualizeMyTotalCountRequestDto,
		stationsdto.VisualizeMyTotalCountResponseDto,
	](
		ctx,
		c.coreClient,
		request,
		stationsdto.VisualizeMyTotalCountOperation,
		"/core/v1/stations/visualizations/total-count",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      response.Data,
		"exception": nil,
	})
}

/* ============================== Controller Methods for Station Permissions ============================== */

func (c *StationController) GetMyStationPermission(ctx *gin.Context, request *stationsdto.GetMyStationPermissionRequestDto) {
	response, exception := coreadapters.CallSecurly[
		stationsdto.GetMyStationPermissionRequestDto,
		stationsdto.GetMyStationPermissionResponseDto,
	](
		ctx,
		c.coreClient,
		request,
		stationsdto.GetMyStationPermissionOperation,
		"/core/v1/stations/permissions/get",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      response.Data,
		"exception": nil,
	})
}

func (c *StationController) CreateMyStationPermission(ctx *gin.Context, request *stationsdto.CreateMyStationPermissionRequestDto) {
	response, exception := coreadapters.CallSecurly[
		stationsdto.CreateMyStationPermissionRequestDto,
		stationsdto.CreateMyStationPermissionResponseDto,
	](
		ctx,
		c.coreClient,
		request,
		stationsdto.CreateMyStationPermissionOperation,
		"/core/v1/stations/permissions/create",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      response.Data,
		"exception": nil,
	})
}

func (c *StationController) UpsertMyStationPermission(
	ctx *gin.Context, request *stationsdto.UpsertMyStationPermissionRequestDto,
) {
	response, exception := coreadapters.CallSecurly[
		stationsdto.UpsertMyStationPermissionRequestDto,
		stationsdto.UpsertMyStationPermissionResponseDto,
	](
		ctx,
		c.coreClient,
		request,
		stationsdto.UpsertMyStationPermissionOperation,
		"/core/v1/stations/permissions/upsert",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      response.Data,
		"exception": nil,
	})
}

func (c *StationController) UpsertMyStationPermissions(
	ctx *gin.Context, request *stationsdto.UpsertMyStationPermissionsRequestDto,
) {
	response, exception := coreadapters.CallSecurly[
		stationsdto.UpsertMyStationPermissionsRequestDto,
		stationsdto.UpsertMyStationPermissionsResponseDto,
	](
		ctx,
		c.coreClient,
		request,
		stationsdto.UpsertMyStationPermissionsOperation,
		"/core/v1/stations/permissions/upsert-many",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      response.Data,
		"exception": nil,
	})
}

func (c *StationController) UpdateMyStationPermission(ctx *gin.Context, request *stationsdto.UpdateMyStationPermissionRequestDto) {
	response, exception := coreadapters.CallSecurly[
		stationsdto.UpdateMyStationPermissionRequestDto,
		stationsdto.UpdateMyStationPermissionResponseDto,
	](
		ctx,
		c.coreClient,
		request,
		stationsdto.UpdateMyStationPermissionOperation,
		"/core/v1/stations/permissions/update",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      response.Data,
		"exception": nil,
	})
}

func (c *StationController) TransferMyStationOwnership(
	ctx *gin.Context, request *stationsdto.TransferMyStationOwnershipRequestDto,
) {
	response, exception := coreadapters.CallSecurly[
		stationsdto.TransferMyStationOwnershipRequestDto,
		stationsdto.TransferMyStationOwnershipResponseDto,
	](
		ctx,
		c.coreClient,
		request,
		stationsdto.TransferMyStationOwnershipOperation,
		"/core/v1/stations/ownership/transfer",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      response.Data,
		"exception": nil,
	})
}

func (c *StationController) DeleteMyStationPermission(
	ctx *gin.Context, request *stationsdto.DeleteMyStationPermissionRequestDto,
) {
	_, exception := coreadapters.CallSecurly[
		stationsdto.DeleteMyStationPermissionRequestDto,
		stationsdto.DeleteMyStationPermissionResponseDto,
	](
		ctx,
		c.coreClient,
		request,
		stationsdto.DeleteMyStationPermissionOperation,
		"/core/v1/stations/permissions/delete",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      stationsdto.DeleteMyStationPermissionResponseDto{},
		"exception": nil,
	})
}

func (c *StationController) DeleteMyStationPermissions(
	ctx *gin.Context, request *stationsdto.DeleteMyStationPermissionsRequestDto,
) {
	_, exception := coreadapters.CallSecurly[
		stationsdto.DeleteMyStationPermissionsRequestDto,
		stationsdto.DeleteMyStationPermissionsResponseDto,
	](
		ctx,
		c.coreClient,
		request,
		stationsdto.DeleteMyStationPermissionsOperation,
		"/core/v1/stations/permissions/delete-many",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *StationController) LeaveMyStation(ctx *gin.Context, request *stationsdto.LeaveMyStationRequestDto) {
	_, exception := coreadapters.CallSecurly[
		stationsdto.LeaveMyStationRequestDto,
		stationsdto.LeaveMyStationResponseDto,
	](
		ctx,
		c.coreClient,
		request,
		stationsdto.LeaveMyStationOperation,
		"/core/v1/stations/memberships/leave",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *StationController) LeaveMyStations(ctx *gin.Context, request *stationsdto.LeaveMyStationsRequestDto) {
	_, exception := coreadapters.CallSecurly[
		stationsdto.LeaveMyStationsRequestDto,
		stationsdto.LeaveMyStationsResponseDto,
	](
		ctx,
		c.coreClient,
		request,
		stationsdto.LeaveMyStationsOperation,
		"/core/v1/stations/memberships/leave-many",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	ctx.Status(http.StatusNoContent)
}
