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
