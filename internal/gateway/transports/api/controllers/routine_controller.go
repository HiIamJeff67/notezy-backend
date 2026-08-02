package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	routinesdto "github.com/HiIamJeff67/notezy-backend/contracts/gateway/v1/api/routines"
	responsewriter "github.com/HiIamJeff67/notezy-backend/internal/shared/responsewriter"
	coreadapters "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/core/adapters"
)

type RoutineControllerInterface interface {
	GetMyRoutineById(ctx *gin.Context, requestDto *routinesdto.GetMyRoutineByIdRequestDto)
	GetMyRoutinesByStationId(ctx *gin.Context, requestDto *routinesdto.GetMyRoutinesByStationIdRequestDto)
	GetAllMyRoutinesByTimeRange(ctx *gin.Context, requestDto *routinesdto.GetAllMyRoutinesByTimeRangeRequestDto)
	CreateRoutineByStationId(ctx *gin.Context, requestDto *routinesdto.CreateRoutineByStationIdRequestDto)
	CreateRoutinesByStationIds(ctx *gin.Context, requestDto *routinesdto.CreateRoutinesByStationIdsRequestDto)
	UpdateMyRoutineById(ctx *gin.Context, requestDto *routinesdto.UpdateMyRoutineByIdRequestDto)
	UpdateMyRoutinesByIds(ctx *gin.Context, requestDto *routinesdto.UpdateMyRoutinesByIdsRequestDto)
	LinkRoutineTagById(ctx *gin.Context, requestDto *routinesdto.LinkRoutineTagByIdRequestDto)
	LinkRoutineTagsByIds(ctx *gin.Context, requestDto *routinesdto.LinkRoutineTagsByIdsRequestDto)
	LinkRoutineItemById(ctx *gin.Context, requestDto *routinesdto.LinkRoutineItemByIdRequestDto)
	LinkRoutineItemsByIds(ctx *gin.Context, requestDto *routinesdto.LinkRoutineItemsByIdsRequestDto)
	RestoreMyRoutineById(ctx *gin.Context, requestDto *routinesdto.RestoreMyRoutineByIdRequestDto)
	RestoreMyRoutinesByIds(ctx *gin.Context, requestDto *routinesdto.RestoreMyRoutinesByIdsRequestDto)
	DeleteMyRoutineById(ctx *gin.Context, requestDto *routinesdto.DeleteMyRoutineByIdRequestDto)
	DeleteMyRoutinesByIds(ctx *gin.Context, requestDto *routinesdto.DeleteMyRoutinesByIdsRequestDto)
	HardDeleteMyRoutineById(ctx *gin.Context, requestDto *routinesdto.HardDeleteMyRoutineByIdRequestDto)
	HardDeleteMyRoutinesByIds(ctx *gin.Context, requestDto *routinesdto.HardDeleteMyRoutinesByIdsRequestDto)

	/* ============================== Visualization Methods ============================== */
	VisualizeMyRoutineStatusCount(ctx *gin.Context, requestDto *routinesdto.VisualizeMyRoutineStatusCountRequestDto)
	VisualizeMyRoutinePeriodCount(ctx *gin.Context, requestDto *routinesdto.VisualizeMyRoutinePeriodCountRequestDto)
	VisualizeMyRoutineScheduledStartAtCount(ctx *gin.Context, requestDto *routinesdto.VisualizeMyRoutineScheduledStartAtCountRequestDto)
	VisualizeMyRoutineScheduledEndAtCount(ctx *gin.Context, requestDto *routinesdto.VisualizeMyRoutineScheduledEndAtCountRequestDto)
}

type RoutineController struct {
	coreClient *coreadapters.CoreClient
}

func NewRoutineController(coreClient *coreadapters.CoreClient) RoutineControllerInterface {
	return &RoutineController{coreClient: coreClient}
}

func (c *RoutineController) GetMyRoutineById(ctx *gin.Context, requestDto *routinesdto.GetMyRoutineByIdRequestDto) {
	response, exception := coreadapters.CallSecurly[routinesdto.GetMyRoutineByIdRequestDto, routinesdto.GetMyRoutineByIdResponseDto](
		ctx,
		c.coreClient,
		requestDto,
		routinesdto.GetMyRoutineByIdOperation,
		"/core/v1/routines/get-by-id",
	)
	if exception != nil {
		responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      response.Data,
		"exception": nil,
	})
}

func (c *RoutineController) GetMyRoutinesByStationId(ctx *gin.Context, requestDto *routinesdto.GetMyRoutinesByStationIdRequestDto) {
	response, exception := coreadapters.CallSecurly[routinesdto.GetMyRoutinesByStationIdRequestDto, routinesdto.GetMyRoutinesByStationIdResponseDto](
		ctx,
		c.coreClient,
		requestDto,
		routinesdto.GetMyRoutinesByStationIdOperation,
		"/core/v1/routines/get-by-station-id",
	)
	if exception != nil {
		responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      response.Data,
		"exception": nil,
	})
}

func (c *RoutineController) GetAllMyRoutinesByTimeRange(ctx *gin.Context, requestDto *routinesdto.GetAllMyRoutinesByTimeRangeRequestDto) {
	response, exception := coreadapters.CallSecurly[routinesdto.GetAllMyRoutinesByTimeRangeRequestDto, routinesdto.GetAllMyRoutinesByTimeRangeResponseDto](
		ctx,
		c.coreClient,
		requestDto,
		routinesdto.GetAllMyRoutinesByTimeRangeOperation,
		"/core/v1/routines/get-all-by-time-range",
	)
	if exception != nil {
		responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      response.Data,
		"exception": nil,
	})
}

func (c *RoutineController) CreateRoutineByStationId(ctx *gin.Context, requestDto *routinesdto.CreateRoutineByStationIdRequestDto) {
	response, exception := coreadapters.CallSecurly[routinesdto.CreateRoutineByStationIdRequestDto, routinesdto.CreateRoutineByStationIdResponseDto](
		ctx,
		c.coreClient,
		requestDto,
		routinesdto.CreateRoutineByStationIdOperation,
		"/core/v1/routines/create-by-station-id",
	)
	if exception != nil {
		responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      response.Data,
		"exception": nil,
	})
}

func (c *RoutineController) CreateRoutinesByStationIds(ctx *gin.Context, requestDto *routinesdto.CreateRoutinesByStationIdsRequestDto) {
	response, exception := coreadapters.CallSecurly[routinesdto.CreateRoutinesByStationIdsRequestDto, routinesdto.CreateRoutinesByStationIdsResponseDto](
		ctx,
		c.coreClient,
		requestDto,
		routinesdto.CreateRoutinesByStationIdsOperation,
		"/core/v1/routines/create-many-by-station-ids",
	)
	if exception != nil {
		responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      response.Data,
		"exception": nil,
	})
}

func (c *RoutineController) UpdateMyRoutineById(ctx *gin.Context, requestDto *routinesdto.UpdateMyRoutineByIdRequestDto) {
	response, exception := coreadapters.CallSecurly[routinesdto.UpdateMyRoutineByIdRequestDto, routinesdto.UpdateMyRoutineByIdResponseDto](
		ctx,
		c.coreClient,
		requestDto,
		routinesdto.UpdateMyRoutineByIdOperation,
		"/core/v1/routines/update",
	)
	if exception != nil {
		responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      response.Data,
		"exception": nil,
	})
}

func (c *RoutineController) UpdateMyRoutinesByIds(ctx *gin.Context, requestDto *routinesdto.UpdateMyRoutinesByIdsRequestDto) {
	response, exception := coreadapters.CallSecurly[routinesdto.UpdateMyRoutinesByIdsRequestDto, routinesdto.UpdateMyRoutinesByIdsResponseDto](
		ctx,
		c.coreClient,
		requestDto,
		routinesdto.UpdateMyRoutinesByIdsOperation,
		"/core/v1/routines/update-many",
	)
	if exception != nil {
		responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      response.Data,
		"exception": nil,
	})
}

func (c *RoutineController) LinkRoutineTagById(ctx *gin.Context, requestDto *routinesdto.LinkRoutineTagByIdRequestDto) {
	response, exception := coreadapters.CallSecurly[routinesdto.LinkRoutineTagByIdRequestDto, routinesdto.LinkRoutineTagByIdResponseDto](
		ctx,
		c.coreClient,
		requestDto,
		routinesdto.LinkRoutineTagByIdOperation,
		"/core/v1/routines/link-tag",
	)
	if exception != nil {
		responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      response.Data,
		"exception": nil,
	})
}

func (c *RoutineController) LinkRoutineTagsByIds(ctx *gin.Context, requestDto *routinesdto.LinkRoutineTagsByIdsRequestDto) {
	response, exception := coreadapters.CallSecurly[routinesdto.LinkRoutineTagsByIdsRequestDto, routinesdto.LinkRoutineTagsByIdsResponseDto](
		ctx,
		c.coreClient,
		requestDto,
		routinesdto.LinkRoutineTagsByIdsOperation,
		"/core/v1/routines/link-tags",
	)
	if exception != nil {
		responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      response.Data,
		"exception": nil,
	})
}

func (c *RoutineController) LinkRoutineItemById(ctx *gin.Context, requestDto *routinesdto.LinkRoutineItemByIdRequestDto) {
	response, exception := coreadapters.CallSecurly[routinesdto.LinkRoutineItemByIdRequestDto, routinesdto.LinkRoutineItemByIdResponseDto](
		ctx,
		c.coreClient,
		requestDto,
		routinesdto.LinkRoutineItemByIdOperation,
		"/core/v1/routines/link-item",
	)
	if exception != nil {
		responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      response.Data,
		"exception": nil,
	})
}

func (c *RoutineController) LinkRoutineItemsByIds(ctx *gin.Context, requestDto *routinesdto.LinkRoutineItemsByIdsRequestDto) {
	response, exception := coreadapters.CallSecurly[routinesdto.LinkRoutineItemsByIdsRequestDto, routinesdto.LinkRoutineItemsByIdsResponseDto](
		ctx,
		c.coreClient,
		requestDto,
		routinesdto.LinkRoutineItemsByIdsOperation,
		"/core/v1/routines/link-items",
	)
	if exception != nil {
		responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      response.Data,
		"exception": nil,
	})
}

func (c *RoutineController) RestoreMyRoutineById(ctx *gin.Context, requestDto *routinesdto.RestoreMyRoutineByIdRequestDto) {
	response, exception := coreadapters.CallSecurly[routinesdto.RestoreMyRoutineByIdRequestDto, routinesdto.RestoreMyRoutineByIdResponseDto](
		ctx,
		c.coreClient,
		requestDto,
		routinesdto.RestoreMyRoutineByIdOperation,
		"/core/v1/routines/restore",
	)
	if exception != nil {
		responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      response.Data,
		"exception": nil,
	})
}

func (c *RoutineController) RestoreMyRoutinesByIds(ctx *gin.Context, requestDto *routinesdto.RestoreMyRoutinesByIdsRequestDto) {
	response, exception := coreadapters.CallSecurly[routinesdto.RestoreMyRoutinesByIdsRequestDto, routinesdto.RestoreMyRoutinesByIdsResponseDto](
		ctx,
		c.coreClient,
		requestDto,
		routinesdto.RestoreMyRoutinesByIdsOperation,
		"/core/v1/routines/restore-many",
	)
	if exception != nil {
		responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      response.Data,
		"exception": nil,
	})
}

func (c *RoutineController) DeleteMyRoutineById(ctx *gin.Context, requestDto *routinesdto.DeleteMyRoutineByIdRequestDto) {
	response, exception := coreadapters.CallSecurly[routinesdto.DeleteMyRoutineByIdRequestDto, routinesdto.DeleteMyRoutineByIdResponseDto](
		ctx,
		c.coreClient,
		requestDto,
		routinesdto.DeleteMyRoutineByIdOperation,
		"/core/v1/routines/delete",
	)
	if exception != nil {
		responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      response.Data,
		"exception": nil,
	})
}

func (c *RoutineController) DeleteMyRoutinesByIds(ctx *gin.Context, requestDto *routinesdto.DeleteMyRoutinesByIdsRequestDto) {
	response, exception := coreadapters.CallSecurly[routinesdto.DeleteMyRoutinesByIdsRequestDto, routinesdto.DeleteMyRoutinesByIdsResponseDto](
		ctx,
		c.coreClient,
		requestDto,
		routinesdto.DeleteMyRoutinesByIdsOperation,
		"/core/v1/routines/delete-many",
	)
	if exception != nil {
		responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      response.Data,
		"exception": nil,
	})
}

func (c *RoutineController) HardDeleteMyRoutineById(ctx *gin.Context, requestDto *routinesdto.HardDeleteMyRoutineByIdRequestDto) {
	response, exception := coreadapters.CallSecurly[routinesdto.HardDeleteMyRoutineByIdRequestDto, routinesdto.HardDeleteMyRoutineByIdResponseDto](
		ctx,
		c.coreClient,
		requestDto,
		routinesdto.HardDeleteMyRoutineByIdOperation,
		"/core/v1/routines/hard-delete",
	)
	if exception != nil {
		responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      response.Data,
		"exception": nil,
	})
}

func (c *RoutineController) HardDeleteMyRoutinesByIds(ctx *gin.Context, requestDto *routinesdto.HardDeleteMyRoutinesByIdsRequestDto) {
	response, exception := coreadapters.CallSecurly[routinesdto.HardDeleteMyRoutinesByIdsRequestDto, routinesdto.HardDeleteMyRoutinesByIdsResponseDto](
		ctx,
		c.coreClient,
		requestDto,
		routinesdto.HardDeleteMyRoutinesByIdsOperation,
		"/core/v1/routines/hard-delete-many",
	)
	if exception != nil {
		responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      response.Data,
		"exception": nil,
	})
}

func (c *RoutineController) VisualizeMyRoutineStatusCount(ctx *gin.Context, requestDto *routinesdto.VisualizeMyRoutineStatusCountRequestDto) {
	response, exception := coreadapters.CallSecurly[routinesdto.VisualizeMyRoutineStatusCountRequestDto, routinesdto.VisualizeMyRoutineStatusCountResponseDto](
		ctx,
		c.coreClient,
		requestDto,
		routinesdto.VisualizeMyRoutineStatusCountOperation,
		"/core/v1/routines/visualize-status-count",
	)
	if exception != nil {
		responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      response.Data,
		"exception": nil,
	})
}

func (c *RoutineController) VisualizeMyRoutinePeriodCount(ctx *gin.Context, requestDto *routinesdto.VisualizeMyRoutinePeriodCountRequestDto) {
	response, exception := coreadapters.CallSecurly[routinesdto.VisualizeMyRoutinePeriodCountRequestDto, routinesdto.VisualizeMyRoutinePeriodCountResponseDto](
		ctx,
		c.coreClient,
		requestDto,
		routinesdto.VisualizeMyRoutinePeriodCountOperation,
		"/core/v1/routines/visualize-period-count",
	)
	if exception != nil {
		responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      response.Data,
		"exception": nil,
	})
}

func (c *RoutineController) VisualizeMyRoutineScheduledStartAtCount(ctx *gin.Context, requestDto *routinesdto.VisualizeMyRoutineScheduledStartAtCountRequestDto) {
	response, exception := coreadapters.CallSecurly[routinesdto.VisualizeMyRoutineScheduledStartAtCountRequestDto, routinesdto.VisualizeMyRoutineScheduledStartAtCountResponseDto](
		ctx,
		c.coreClient,
		requestDto,
		routinesdto.VisualizeMyRoutineScheduledStartAtCountOperation,
		"/core/v1/routines/visualize-scheduled-start-at-count",
	)
	if exception != nil {
		responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      response.Data,
		"exception": nil,
	})
}

func (c *RoutineController) VisualizeMyRoutineScheduledEndAtCount(ctx *gin.Context, requestDto *routinesdto.VisualizeMyRoutineScheduledEndAtCountRequestDto) {
	response, exception := coreadapters.CallSecurly[routinesdto.VisualizeMyRoutineScheduledEndAtCountRequestDto, routinesdto.VisualizeMyRoutineScheduledEndAtCountResponseDto](
		ctx,
		c.coreClient,
		requestDto,
		routinesdto.VisualizeMyRoutineScheduledEndAtCountOperation,
		"/core/v1/routines/visualize-scheduled-end-at-count",
	)
	if exception != nil {
		responsewriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      response.Data,
		"exception": nil,
	})
}
