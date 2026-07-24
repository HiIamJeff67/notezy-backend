package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	dtos "github.com/HiIamJeff67/notezy-backend/app/dtos"
	services "github.com/HiIamJeff67/notezy-backend/app/services"
)

type StationControllerInterface interface {
	GetMyStationById(ctx *gin.Context, reqDto *dtos.GetMyStationByIdReqDto)
	GetAllMyStations(ctx *gin.Context, reqDto *dtos.GetAllMyStationsReqDto)
	CreateStation(ctx *gin.Context, reqDto *dtos.CreateStationReqDto)
	CreateStations(ctx *gin.Context, reqDto *dtos.CreateStationsReqDto)
	UpdateMyStationById(ctx *gin.Context, reqDto *dtos.UpdateMyStationByIdReqDto)
	UpdateMyStationsByIds(ctx *gin.Context, reqDto *dtos.UpdateMyStationsByIdsReqDto)
	RestoreMyStationById(ctx *gin.Context, reqDto *dtos.RestoreMyStationByIdReqDto)
	RestoreMyStationsByIds(ctx *gin.Context, reqDto *dtos.RestoreMyStationsByIdsReqDto)
	DeleteMyStationById(ctx *gin.Context, reqDto *dtos.DeleteMyStationByIdReqDto)
	DeleteMyStationsByIds(ctx *gin.Context, reqDto *dtos.DeleteMyStationsByIdsReqDto)
	HardDeleteMyStationById(ctx *gin.Context, reqDto *dtos.HardDeleteMyStationByIdReqDto)
	HardDeleteMyStationsByIds(ctx *gin.Context, reqDto *dtos.HardDeleteMyStationsByIdsReqDto)
	VisualizeMyTotalCount(ctx *gin.Context, reqDto *dtos.VisualizeMyTotalCountReqDto)

	GetMyStationPermission(ctx *gin.Context, reqDto *dtos.GetMyStationPermissionReqDto)
	CreateMyStationPermission(ctx *gin.Context, reqDto *dtos.CreateMyStationPermissionReqDto)
	UpsertMyStationPermission(ctx *gin.Context, reqDto *dtos.UpsertMyStationPermissionReqDto)
	UpsertMyStationPermissions(ctx *gin.Context, reqDto *dtos.UpsertMyStationPermissionsReqDto)
	UpdateMyStationPermission(ctx *gin.Context, reqDto *dtos.UpdateMyStationPermissionReqDto)
	TransferMyStationOwnership(ctx *gin.Context, reqDto *dtos.TransferMyStationOwnershipReqDto)
	DeleteMyStationPermission(ctx *gin.Context, reqDto *dtos.DeleteMyStationPermissionReqDto)
	DeleteMyStationPermissions(ctx *gin.Context, reqDto *dtos.DeleteMyStationPermissionsReqDto)
	LeaveMyStation(ctx *gin.Context, reqDto *dtos.LeaveMyStationReqDto)
	LeaveMyStations(ctx *gin.Context, reqDto *dtos.LeaveMyStationsReqDto)
}

type StationController struct {
	stationService services.StationServiceInterface
}

func NewStationController(service services.StationServiceInterface) StationControllerInterface {
	return &StationController{
		stationService: service,
	}
}

func (c *StationController) GetMyStationById(ctx *gin.Context, reqDto *dtos.GetMyStationByIdReqDto) {
	resDto, exception := c.stationService.GetMyStationById(ctx.Request.Context(), reqDto)
	if exception != nil {
		exception.Log().SafelyAbortAndResponseWithJSON(ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      resDto,
		"exception": nil,
	})
}

func (c *StationController) GetAllMyStations(ctx *gin.Context, reqDto *dtos.GetAllMyStationsReqDto) {
	resDto, exception := c.stationService.GetAllMyStations(ctx.Request.Context(), reqDto)
	if exception != nil {
		exception.Log().SafelyAbortAndResponseWithJSON(ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      resDto,
		"exception": nil,
	})
}

func (c *StationController) CreateStation(ctx *gin.Context, reqDto *dtos.CreateStationReqDto) {
	resDto, exception := c.stationService.CreateStation(ctx.Request.Context(), reqDto)
	if exception != nil {
		exception.Log().SafelyAbortAndResponseWithJSON(ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      resDto,
		"exception": nil,
	})
}

func (c *StationController) CreateStations(ctx *gin.Context, reqDto *dtos.CreateStationsReqDto) {
	resDto, exception := c.stationService.CreateStations(ctx.Request.Context(), reqDto)
	if exception != nil {
		exception.Log().SafelyAbortAndResponseWithJSON(ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      resDto,
		"exception": nil,
	})
}

func (c *StationController) UpdateMyStationById(ctx *gin.Context, reqDto *dtos.UpdateMyStationByIdReqDto) {
	resDto, exception := c.stationService.UpdateMyStationById(ctx.Request.Context(), reqDto)
	if exception != nil {
		exception.Log().SafelyAbortAndResponseWithJSON(ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      resDto,
		"exception": nil,
	})
}

func (c *StationController) UpdateMyStationsByIds(ctx *gin.Context, reqDto *dtos.UpdateMyStationsByIdsReqDto) {
	resDto, exception := c.stationService.UpdateMyStationsByIds(ctx.Request.Context(), reqDto)
	if exception != nil {
		exception.Log().SafelyAbortAndResponseWithJSON(ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      resDto,
		"exception": nil,
	})
}

func (c *StationController) RestoreMyStationById(ctx *gin.Context, reqDto *dtos.RestoreMyStationByIdReqDto) {
	resDto, exception := c.stationService.RestoreMyStationById(ctx.Request.Context(), reqDto)
	if exception != nil {
		exception.Log().SafelyAbortAndResponseWithJSON(ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      resDto,
		"exception": nil,
	})
}

func (c *StationController) RestoreMyStationsByIds(ctx *gin.Context, reqDto *dtos.RestoreMyStationsByIdsReqDto) {
	resDto, exception := c.stationService.RestoreMyStationsByIds(ctx.Request.Context(), reqDto)
	if exception != nil {
		exception.Log().SafelyAbortAndResponseWithJSON(ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      resDto,
		"exception": nil,
	})
}

func (c *StationController) DeleteMyStationById(ctx *gin.Context, reqDto *dtos.DeleteMyStationByIdReqDto) {
	resDto, exception := c.stationService.DeleteMyStationById(ctx.Request.Context(), reqDto)
	if exception != nil {
		exception.Log().SafelyAbortAndResponseWithJSON(ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      resDto,
		"exception": nil,
	})
}

func (c *StationController) DeleteMyStationsByIds(ctx *gin.Context, reqDto *dtos.DeleteMyStationsByIdsReqDto) {
	resDto, exception := c.stationService.DeleteMyStationsByIds(ctx.Request.Context(), reqDto)
	if exception != nil {
		exception.Log().SafelyAbortAndResponseWithJSON(ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      resDto,
		"exception": nil,
	})
}

func (c *StationController) HardDeleteMyStationById(ctx *gin.Context, reqDto *dtos.HardDeleteMyStationByIdReqDto) {
	resDto, exception := c.stationService.HardDeleteMyStationById(ctx.Request.Context(), reqDto)
	if exception != nil {
		exception.Log().SafelyAbortAndResponseWithJSON(ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      resDto,
		"exception": nil,
	})
}

func (c *StationController) HardDeleteMyStationsByIds(ctx *gin.Context, reqDto *dtos.HardDeleteMyStationsByIdsReqDto) {
	resDto, exception := c.stationService.HardDeleteMyStationsByIds(ctx.Request.Context(), reqDto)
	if exception != nil {
		exception.Log().SafelyAbortAndResponseWithJSON(ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      resDto,
		"exception": nil,
	})
}

/* ============================== Controller Methods for Visualization ============================== */

func (c *StationController) VisualizeMyTotalCount(ctx *gin.Context, reqDto *dtos.VisualizeMyTotalCountReqDto) {
	resDto, exception := c.stationService.VisualizeMyTotalCount(ctx.Request.Context(), reqDto)
	if exception != nil {
		exception.Log().SafelyAbortAndResponseWithJSON(ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      resDto,
		"exception": nil,
	})
}

/* ============================== Controller Methods for Station Permissions ============================== */

func (c *StationController) GetMyStationPermission(ctx *gin.Context, reqDto *dtos.GetMyStationPermissionReqDto) {
	resDto, exception := c.stationService.GetMyStationPermission(ctx.Request.Context(), reqDto)
	if exception != nil {
		exception.Log().SafelyAbortAndResponseWithJSON(ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      resDto,
		"exception": nil,
	})
}

func (c *StationController) CreateMyStationPermission(ctx *gin.Context, reqDto *dtos.CreateMyStationPermissionReqDto) {
	resDto, exception := c.stationService.CreateMyStationPermission(ctx.Request.Context(), reqDto)
	if exception != nil {
		exception.Log().SafelyAbortAndResponseWithJSON(ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      resDto,
		"exception": nil,
	})
}

func (c *StationController) UpsertMyStationPermission(
	ctx *gin.Context, reqDto *dtos.UpsertMyStationPermissionReqDto,
) {
	resDto, exception := c.stationService.UpsertMyStationPermission(ctx.Request.Context(), reqDto)
	if exception != nil {
		exception.Log().SafelyAbortAndResponseWithJSON(ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      resDto,
		"exception": nil,
	})
}

func (c *StationController) UpsertMyStationPermissions(
	ctx *gin.Context, reqDto *dtos.UpsertMyStationPermissionsReqDto,
) {
	resDto, exception := c.stationService.UpsertMyStationPermissions(ctx.Request.Context(), reqDto)
	if exception != nil {
		exception.Log().SafelyAbortAndResponseWithJSON(ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      resDto,
		"exception": nil,
	})
}

func (c *StationController) UpdateMyStationPermission(ctx *gin.Context, reqDto *dtos.UpdateMyStationPermissionReqDto) {
	resDto, exception := c.stationService.UpdateMyStationPermission(ctx.Request.Context(), reqDto)
	if exception != nil {
		exception.Log().SafelyAbortAndResponseWithJSON(ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      resDto,
		"exception": nil,
	})
}

func (c *StationController) TransferMyStationOwnership(
	ctx *gin.Context, reqDto *dtos.TransferMyStationOwnershipReqDto,
) {
	resDto, exception := c.stationService.TransferMyStationOwnership(ctx.Request.Context(), reqDto)
	if exception != nil {
		exception.Log().SafelyAbortAndResponseWithJSON(ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      resDto,
		"exception": nil,
	})
}

func (c *StationController) DeleteMyStationPermission(ctx *gin.Context, reqDto *dtos.DeleteMyStationPermissionReqDto) {
	if exception := c.stationService.DeleteMyStationPermission(ctx.Request.Context(), reqDto); exception != nil {
		exception.Log().SafelyAbortAndResponseWithJSON(ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      nil,
		"exception": nil,
	})
}

func (c *StationController) DeleteMyStationPermissions(
	ctx *gin.Context, reqDto *dtos.DeleteMyStationPermissionsReqDto,
) {
	if exception := c.stationService.DeleteMyStationPermissions(ctx.Request.Context(), reqDto); exception != nil {
		exception.Log().SafelyAbortAndResponseWithJSON(ctx)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *StationController) LeaveMyStation(ctx *gin.Context, reqDto *dtos.LeaveMyStationReqDto) {
	if exception := c.stationService.LeaveMyStation(ctx.Request.Context(), reqDto); exception != nil {
		exception.Log().SafelyAbortAndResponseWithJSON(ctx)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *StationController) LeaveMyStations(ctx *gin.Context, reqDto *dtos.LeaveMyStationsReqDto) {
	if exception := c.stationService.LeaveMyStations(ctx.Request.Context(), reqDto); exception != nil {
		exception.Log().SafelyAbortAndResponseWithJSON(ctx)
		return
	}

	ctx.Status(http.StatusNoContent)
}
