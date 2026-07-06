package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	dtos "github.com/HiIamJeff67/notezy-backend/app/dtos"
	services "github.com/HiIamJeff67/notezy-backend/app/services"
)

type RoutineControllerInterface interface {
	GetMyRoutineById(ctx *gin.Context, reqDto *dtos.GetMyRoutineByIdReqDto)
	GetMyRoutinesByStationId(ctx *gin.Context, reqDto *dtos.GetMyRoutinesByStationIdReqDto)
	GetAllMyRoutinesByTimeRange(ctx *gin.Context, reqDto *dtos.GetAllMyRoutinesByTimeRangeReqDto)
	CreateRoutineByStationId(ctx *gin.Context, reqDto *dtos.CreateRoutineByStationIdReqDto)
	CreateRoutinesByStationIds(ctx *gin.Context, reqDto *dtos.CreateRoutinesByStationIdsReqDto)
	UpdateMyRoutineById(ctx *gin.Context, reqDto *dtos.UpdateMyRoutineByIdReqDto)
	UpdateMyRoutinesByIds(ctx *gin.Context, reqDto *dtos.UpdateMyRoutinesByIdsReqDto)
	LinkRoutineTagById(ctx *gin.Context, reqDto *dtos.LinkRoutineTagByIdReqDto)
	LinkRoutineTagsByIds(ctx *gin.Context, reqDto *dtos.LinkRoutineTagsByIdsReqDto)
	LinkRoutineItemById(ctx *gin.Context, reqDto *dtos.LinkRoutineItemByIdReqDto)
	LinkRoutineItemsByIds(ctx *gin.Context, reqDto *dtos.LinkRoutineItemsByIdsReqDto)
	RestoreMyRoutineById(ctx *gin.Context, reqDto *dtos.RestoreMyRoutineByIdReqDto)
	RestoreMyRoutinesByIds(ctx *gin.Context, reqDto *dtos.RestoreMyRoutinesByIdsReqDto)
	DeleteMyRoutineById(ctx *gin.Context, reqDto *dtos.DeleteMyRoutineByIdReqDto)
	DeleteMyRoutinesByIds(ctx *gin.Context, reqDto *dtos.DeleteMyRoutinesByIdsReqDto)
	HardDeleteMyRoutineById(ctx *gin.Context, reqDto *dtos.HardDeleteMyRoutineByIdReqDto)
	HardDeleteMyRoutinesByIds(ctx *gin.Context, reqDto *dtos.HardDeleteMyRoutinesByIdsReqDto)
	VisualizeMyRoutineStatusCount(ctx *gin.Context, reqDto *dtos.VisualizeMyRoutineStatusCountReqDto)
	VisualizeMyRoutinePeriodCount(ctx *gin.Context, reqDto *dtos.VisualizeMyRoutinePeriodCountReqDto)
	VisualizeMyRoutineScheduledStartAtCount(ctx *gin.Context, reqDto *dtos.VisualizeMyRoutineScheduledStartAtCountReqDto)
	VisualizeMyRoutineScheduledEndAtCount(ctx *gin.Context, reqDto *dtos.VisualizeMyRoutineScheduledEndAtCountReqDto)
}

type RoutineController struct {
	routineService services.RoutineServiceInterface
}

func NewRoutineController(service services.RoutineServiceInterface) RoutineControllerInterface {
	return &RoutineController{
		routineService: service,
	}
}

func (c *RoutineController) GetMyRoutineById(ctx *gin.Context, reqDto *dtos.GetMyRoutineByIdReqDto) {
	resDto, exception := c.routineService.GetMyRoutineById(ctx.Request.Context(), reqDto)
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

func (c *RoutineController) GetMyRoutinesByStationId(ctx *gin.Context, reqDto *dtos.GetMyRoutinesByStationIdReqDto) {
	resDto, exception := c.routineService.GetMyRoutinesByStationId(ctx.Request.Context(), reqDto)
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

func (c *RoutineController) GetAllMyRoutinesByTimeRange(
	ctx *gin.Context,
	reqDto *dtos.GetAllMyRoutinesByTimeRangeReqDto,
) {
	resDto, exception := c.routineService.GetAllMyRoutinesByTimeRange(ctx.Request.Context(), reqDto)
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

func (c *RoutineController) CreateRoutineByStationId(ctx *gin.Context, reqDto *dtos.CreateRoutineByStationIdReqDto) {
	resDto, exception := c.routineService.CreateRoutineByStationId(ctx.Request.Context(), reqDto)
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

func (c *RoutineController) CreateRoutinesByStationIds(ctx *gin.Context, reqDto *dtos.CreateRoutinesByStationIdsReqDto) {
	resDto, exception := c.routineService.CreateRoutinesByStationIds(ctx.Request.Context(), reqDto)
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

func (c *RoutineController) UpdateMyRoutineById(ctx *gin.Context, reqDto *dtos.UpdateMyRoutineByIdReqDto) {
	resDto, exception := c.routineService.UpdateMyRoutineById(ctx.Request.Context(), reqDto)
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

func (c *RoutineController) UpdateMyRoutinesByIds(ctx *gin.Context, reqDto *dtos.UpdateMyRoutinesByIdsReqDto) {
	resDto, exception := c.routineService.UpdateMyRoutinesByIds(ctx.Request.Context(), reqDto)
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

func (c *RoutineController) LinkRoutineTagById(ctx *gin.Context, reqDto *dtos.LinkRoutineTagByIdReqDto) {
	resDto, exception := c.routineService.LinkRoutineTagById(ctx.Request.Context(), reqDto)
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

func (c *RoutineController) LinkRoutineTagsByIds(ctx *gin.Context, reqDto *dtos.LinkRoutineTagsByIdsReqDto) {
	resDto, exception := c.routineService.LinkRoutineTagsByIds(ctx.Request.Context(), reqDto)
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

func (c *RoutineController) LinkRoutineItemById(ctx *gin.Context, reqDto *dtos.LinkRoutineItemByIdReqDto) {
	resDto, exception := c.routineService.LinkRoutineItemById(ctx.Request.Context(), reqDto)
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

func (c *RoutineController) LinkRoutineItemsByIds(ctx *gin.Context, reqDto *dtos.LinkRoutineItemsByIdsReqDto) {
	resDto, exception := c.routineService.LinkRoutineItemsByIds(ctx.Request.Context(), reqDto)
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

func (c *RoutineController) RestoreMyRoutineById(ctx *gin.Context, reqDto *dtos.RestoreMyRoutineByIdReqDto) {
	resDto, exception := c.routineService.RestoreMyRoutineById(ctx.Request.Context(), reqDto)
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

func (c *RoutineController) RestoreMyRoutinesByIds(ctx *gin.Context, reqDto *dtos.RestoreMyRoutinesByIdsReqDto) {
	resDto, exception := c.routineService.RestoreMyRoutinesByIds(ctx.Request.Context(), reqDto)
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

func (c *RoutineController) DeleteMyRoutineById(ctx *gin.Context, reqDto *dtos.DeleteMyRoutineByIdReqDto) {
	resDto, exception := c.routineService.DeleteMyRoutineById(ctx.Request.Context(), reqDto)
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

func (c *RoutineController) DeleteMyRoutinesByIds(ctx *gin.Context, reqDto *dtos.DeleteMyRoutinesByIdsReqDto) {
	resDto, exception := c.routineService.DeleteMyRoutinesByIds(ctx.Request.Context(), reqDto)
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

func (c *RoutineController) HardDeleteMyRoutineById(ctx *gin.Context, reqDto *dtos.HardDeleteMyRoutineByIdReqDto) {
	resDto, exception := c.routineService.HardDeleteMyRoutineById(ctx.Request.Context(), reqDto)
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

func (c *RoutineController) HardDeleteMyRoutinesByIds(ctx *gin.Context, reqDto *dtos.HardDeleteMyRoutinesByIdsReqDto) {
	resDto, exception := c.routineService.HardDeleteMyRoutinesByIds(ctx.Request.Context(), reqDto)
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

func (c *RoutineController) VisualizeMyRoutineStatusCount(ctx *gin.Context, reqDto *dtos.VisualizeMyRoutineStatusCountReqDto) {
	resDto, exception := c.routineService.VisualizeMyRoutineStatusCount(ctx.Request.Context(), reqDto)
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

func (c *RoutineController) VisualizeMyRoutinePeriodCount(ctx *gin.Context, reqDto *dtos.VisualizeMyRoutinePeriodCountReqDto) {
	resDto, exception := c.routineService.VisualizeMyRoutinePeriodCount(ctx.Request.Context(), reqDto)
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

func (c *RoutineController) VisualizeMyRoutineScheduledStartAtCount(ctx *gin.Context, reqDto *dtos.VisualizeMyRoutineScheduledStartAtCountReqDto) {
	resDto, exception := c.routineService.VisualizeMyRoutineScheduledStartAtCount(ctx.Request.Context(), reqDto)
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

func (c *RoutineController) VisualizeMyRoutineScheduledEndAtCount(ctx *gin.Context, reqDto *dtos.VisualizeMyRoutineScheduledEndAtCountReqDto) {
	resDto, exception := c.routineService.VisualizeMyRoutineScheduledEndAtCount(ctx.Request.Context(), reqDto)
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
