package controllers

import (
	"github.com/gin-gonic/gin"

	exceptionwriter "github.com/HiIamJeff67/notegic-backend/shared/util/exceptionwriter"

	apicontract "github.com/HiIamJeff67/notegic-backend/contracts/core/v1/api/routines"

	coreadapters "github.com/HiIamJeff67/notegic-backend/internal/apigateway/transports/core/adapters"
)

type RoutineControllerInterface interface {
	GetMyRoutineById(ctx *gin.Context, requestDto *apicontract.GetMyRoutineByIdRequestDto)
	GetMyRoutinesByStationId(ctx *gin.Context, requestDto *apicontract.GetMyRoutinesByStationIdRequestDto)
	GetAllMyRoutinesByTimeRange(ctx *gin.Context, requestDto *apicontract.GetAllMyRoutinesByTimeRangeRequestDto)
	CreateRoutineByStationId(ctx *gin.Context, requestDto *apicontract.CreateRoutineByStationIdRequestDto)
	CreateRoutinesByStationIds(ctx *gin.Context, requestDto *apicontract.CreateRoutinesByStationIdsRequestDto)
	UpdateMyRoutineById(ctx *gin.Context, requestDto *apicontract.UpdateMyRoutineByIdRequestDto)
	UpdateMyRoutinesByIds(ctx *gin.Context, requestDto *apicontract.UpdateMyRoutinesByIdsRequestDto)
	LinkRoutineTagById(ctx *gin.Context, requestDto *apicontract.LinkRoutineTagByIdRequestDto)
	LinkRoutineTagsByIds(ctx *gin.Context, requestDto *apicontract.LinkRoutineTagsByIdsRequestDto)
	LinkRoutineItemById(ctx *gin.Context, requestDto *apicontract.LinkRoutineItemByIdRequestDto)
	LinkRoutineItemsByIds(ctx *gin.Context, requestDto *apicontract.LinkRoutineItemsByIdsRequestDto)
	RestoreMyRoutineById(ctx *gin.Context, requestDto *apicontract.RestoreMyRoutineByIdRequestDto)
	RestoreMyRoutinesByIds(ctx *gin.Context, requestDto *apicontract.RestoreMyRoutinesByIdsRequestDto)
	DeleteMyRoutineById(ctx *gin.Context, requestDto *apicontract.DeleteMyRoutineByIdRequestDto)
	DeleteMyRoutinesByIds(ctx *gin.Context, requestDto *apicontract.DeleteMyRoutinesByIdsRequestDto)
	HardDeleteMyRoutineById(ctx *gin.Context, requestDto *apicontract.HardDeleteMyRoutineByIdRequestDto)
	HardDeleteMyRoutinesByIds(ctx *gin.Context, requestDto *apicontract.HardDeleteMyRoutinesByIdsRequestDto)

	/* ============================== Visualization Methods ============================== */
	VisualizeMyRoutineStatusCount(ctx *gin.Context, requestDto *apicontract.VisualizeMyRoutineStatusCountRequestDto)
	VisualizeMyRoutinePeriodCount(ctx *gin.Context, requestDto *apicontract.VisualizeMyRoutinePeriodCountRequestDto)
	VisualizeMyRoutineScheduledStartAtCount(ctx *gin.Context, requestDto *apicontract.VisualizeMyRoutineScheduledStartAtCountRequestDto)
	VisualizeMyRoutineScheduledEndAtCount(ctx *gin.Context, requestDto *apicontract.VisualizeMyRoutineScheduledEndAtCountRequestDto)
}

type RoutineController struct {
	coreAdapter *coreadapters.CoreAdapter
}

func NewRoutineController(coreAdapter *coreadapters.CoreAdapter) RoutineControllerInterface {
	return &RoutineController{coreAdapter: coreAdapter}
}

func (c *RoutineController) GetMyRoutineById(ctx *gin.Context, requestDto *apicontract.GetMyRoutineByIdRequestDto) {
	response, exception := coreadapters.CallSecurly[apicontract.GetMyRoutineByIdRequestDto, apicontract.GetMyRoutineByIdResponseDto](
		ctx,
		c.coreAdapter,
		requestDto,
		apicontract.GetMyRoutineByIdOperation,
		"/core/v1/routines/get-by-id",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *RoutineController) GetMyRoutinesByStationId(ctx *gin.Context, requestDto *apicontract.GetMyRoutinesByStationIdRequestDto) {
	response, exception := coreadapters.CallSecurly[apicontract.GetMyRoutinesByStationIdRequestDto, apicontract.GetMyRoutinesByStationIdResponseDto](
		ctx,
		c.coreAdapter,
		requestDto,
		apicontract.GetMyRoutinesByStationIdOperation,
		"/core/v1/routines/get-by-station-id",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *RoutineController) GetAllMyRoutinesByTimeRange(ctx *gin.Context, requestDto *apicontract.GetAllMyRoutinesByTimeRangeRequestDto) {
	response, exception := coreadapters.CallSecurly[apicontract.GetAllMyRoutinesByTimeRangeRequestDto, apicontract.GetAllMyRoutinesByTimeRangeResponseDto](
		ctx,
		c.coreAdapter,
		requestDto,
		apicontract.GetAllMyRoutinesByTimeRangeOperation,
		"/core/v1/routines/get-all-by-time-range",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *RoutineController) CreateRoutineByStationId(ctx *gin.Context, requestDto *apicontract.CreateRoutineByStationIdRequestDto) {
	response, exception := coreadapters.CallSecurly[apicontract.CreateRoutineByStationIdRequestDto, apicontract.CreateRoutineByStationIdResponseDto](
		ctx,
		c.coreAdapter,
		requestDto,
		apicontract.CreateRoutineByStationIdOperation,
		"/core/v1/routines/create-by-station-id",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeCreatedClientResponse(ctx, response.Data)
}

func (c *RoutineController) CreateRoutinesByStationIds(ctx *gin.Context, requestDto *apicontract.CreateRoutinesByStationIdsRequestDto) {
	response, exception := coreadapters.CallSecurly[apicontract.CreateRoutinesByStationIdsRequestDto, apicontract.CreateRoutinesByStationIdsResponseDto](
		ctx,
		c.coreAdapter,
		requestDto,
		apicontract.CreateRoutinesByStationIdsOperation,
		"/core/v1/routines/create-many-by-station-ids",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeCreatedClientResponse(ctx, response.Data)
}

func (c *RoutineController) UpdateMyRoutineById(ctx *gin.Context, requestDto *apicontract.UpdateMyRoutineByIdRequestDto) {
	response, exception := coreadapters.CallSecurly[apicontract.UpdateMyRoutineByIdRequestDto, apicontract.UpdateMyRoutineByIdResponseDto](
		ctx,
		c.coreAdapter,
		requestDto,
		apicontract.UpdateMyRoutineByIdOperation,
		"/core/v1/routines/update",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *RoutineController) UpdateMyRoutinesByIds(ctx *gin.Context, requestDto *apicontract.UpdateMyRoutinesByIdsRequestDto) {
	response, exception := coreadapters.CallSecurly[apicontract.UpdateMyRoutinesByIdsRequestDto, apicontract.UpdateMyRoutinesByIdsResponseDto](
		ctx,
		c.coreAdapter,
		requestDto,
		apicontract.UpdateMyRoutinesByIdsOperation,
		"/core/v1/routines/update-many",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *RoutineController) LinkRoutineTagById(ctx *gin.Context, requestDto *apicontract.LinkRoutineTagByIdRequestDto) {
	response, exception := coreadapters.CallSecurly[apicontract.LinkRoutineTagByIdRequestDto, apicontract.LinkRoutineTagByIdResponseDto](
		ctx,
		c.coreAdapter,
		requestDto,
		apicontract.LinkRoutineTagByIdOperation,
		"/core/v1/routines/link-tag",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *RoutineController) LinkRoutineTagsByIds(ctx *gin.Context, requestDto *apicontract.LinkRoutineTagsByIdsRequestDto) {
	response, exception := coreadapters.CallSecurly[apicontract.LinkRoutineTagsByIdsRequestDto, apicontract.LinkRoutineTagsByIdsResponseDto](
		ctx,
		c.coreAdapter,
		requestDto,
		apicontract.LinkRoutineTagsByIdsOperation,
		"/core/v1/routines/link-tags",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *RoutineController) LinkRoutineItemById(ctx *gin.Context, requestDto *apicontract.LinkRoutineItemByIdRequestDto) {
	response, exception := coreadapters.CallSecurly[apicontract.LinkRoutineItemByIdRequestDto, apicontract.LinkRoutineItemByIdResponseDto](
		ctx,
		c.coreAdapter,
		requestDto,
		apicontract.LinkRoutineItemByIdOperation,
		"/core/v1/routines/link-item",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *RoutineController) LinkRoutineItemsByIds(ctx *gin.Context, requestDto *apicontract.LinkRoutineItemsByIdsRequestDto) {
	response, exception := coreadapters.CallSecurly[apicontract.LinkRoutineItemsByIdsRequestDto, apicontract.LinkRoutineItemsByIdsResponseDto](
		ctx,
		c.coreAdapter,
		requestDto,
		apicontract.LinkRoutineItemsByIdsOperation,
		"/core/v1/routines/link-items",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *RoutineController) RestoreMyRoutineById(ctx *gin.Context, requestDto *apicontract.RestoreMyRoutineByIdRequestDto) {
	response, exception := coreadapters.CallSecurly[apicontract.RestoreMyRoutineByIdRequestDto, apicontract.RestoreMyRoutineByIdResponseDto](
		ctx,
		c.coreAdapter,
		requestDto,
		apicontract.RestoreMyRoutineByIdOperation,
		"/core/v1/routines/restore",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *RoutineController) RestoreMyRoutinesByIds(ctx *gin.Context, requestDto *apicontract.RestoreMyRoutinesByIdsRequestDto) {
	response, exception := coreadapters.CallSecurly[apicontract.RestoreMyRoutinesByIdsRequestDto, apicontract.RestoreMyRoutinesByIdsResponseDto](
		ctx,
		c.coreAdapter,
		requestDto,
		apicontract.RestoreMyRoutinesByIdsOperation,
		"/core/v1/routines/restore-many",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *RoutineController) DeleteMyRoutineById(ctx *gin.Context, requestDto *apicontract.DeleteMyRoutineByIdRequestDto) {
	response, exception := coreadapters.CallSecurly[apicontract.DeleteMyRoutineByIdRequestDto, apicontract.DeleteMyRoutineByIdResponseDto](
		ctx,
		c.coreAdapter,
		requestDto,
		apicontract.DeleteMyRoutineByIdOperation,
		"/core/v1/routines/delete",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *RoutineController) DeleteMyRoutinesByIds(ctx *gin.Context, requestDto *apicontract.DeleteMyRoutinesByIdsRequestDto) {
	response, exception := coreadapters.CallSecurly[apicontract.DeleteMyRoutinesByIdsRequestDto, apicontract.DeleteMyRoutinesByIdsResponseDto](
		ctx,
		c.coreAdapter,
		requestDto,
		apicontract.DeleteMyRoutinesByIdsOperation,
		"/core/v1/routines/delete-many",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *RoutineController) HardDeleteMyRoutineById(ctx *gin.Context, requestDto *apicontract.HardDeleteMyRoutineByIdRequestDto) {
	response, exception := coreadapters.CallSecurly[apicontract.HardDeleteMyRoutineByIdRequestDto, apicontract.HardDeleteMyRoutineByIdResponseDto](
		ctx,
		c.coreAdapter,
		requestDto,
		apicontract.HardDeleteMyRoutineByIdOperation,
		"/core/v1/routines/hard-delete",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *RoutineController) HardDeleteMyRoutinesByIds(ctx *gin.Context, requestDto *apicontract.HardDeleteMyRoutinesByIdsRequestDto) {
	response, exception := coreadapters.CallSecurly[apicontract.HardDeleteMyRoutinesByIdsRequestDto, apicontract.HardDeleteMyRoutinesByIdsResponseDto](
		ctx,
		c.coreAdapter,
		requestDto,
		apicontract.HardDeleteMyRoutinesByIdsOperation,
		"/core/v1/routines/hard-delete-many",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *RoutineController) VisualizeMyRoutineStatusCount(ctx *gin.Context, requestDto *apicontract.VisualizeMyRoutineStatusCountRequestDto) {
	response, exception := coreadapters.CallSecurly[apicontract.VisualizeMyRoutineStatusCountRequestDto, apicontract.VisualizeMyRoutineStatusCountResponseDto](
		ctx,
		c.coreAdapter,
		requestDto,
		apicontract.VisualizeMyRoutineStatusCountOperation,
		"/core/v1/routines/visualize-status-count",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *RoutineController) VisualizeMyRoutinePeriodCount(ctx *gin.Context, requestDto *apicontract.VisualizeMyRoutinePeriodCountRequestDto) {
	response, exception := coreadapters.CallSecurly[apicontract.VisualizeMyRoutinePeriodCountRequestDto, apicontract.VisualizeMyRoutinePeriodCountResponseDto](
		ctx,
		c.coreAdapter,
		requestDto,
		apicontract.VisualizeMyRoutinePeriodCountOperation,
		"/core/v1/routines/visualize-period-count",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *RoutineController) VisualizeMyRoutineScheduledStartAtCount(ctx *gin.Context, requestDto *apicontract.VisualizeMyRoutineScheduledStartAtCountRequestDto) {
	response, exception := coreadapters.CallSecurly[apicontract.VisualizeMyRoutineScheduledStartAtCountRequestDto, apicontract.VisualizeMyRoutineScheduledStartAtCountResponseDto](
		ctx,
		c.coreAdapter,
		requestDto,
		apicontract.VisualizeMyRoutineScheduledStartAtCountOperation,
		"/core/v1/routines/visualize-scheduled-start-at-count",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}

func (c *RoutineController) VisualizeMyRoutineScheduledEndAtCount(ctx *gin.Context, requestDto *apicontract.VisualizeMyRoutineScheduledEndAtCountRequestDto) {
	response, exception := coreadapters.CallSecurly[apicontract.VisualizeMyRoutineScheduledEndAtCountRequestDto, apicontract.VisualizeMyRoutineScheduledEndAtCountResponseDto](
		ctx,
		c.coreAdapter,
		requestDto,
		apicontract.VisualizeMyRoutineScheduledEndAtCountOperation,
		"/core/v1/routines/visualize-scheduled-end-at-count",
	)
	if exception != nil {
		exceptionwriter.SafelyAbortAndResponseWithJSON(exception, ctx)
		return
	}

	writeClientResponse(ctx, response.Data)
}
