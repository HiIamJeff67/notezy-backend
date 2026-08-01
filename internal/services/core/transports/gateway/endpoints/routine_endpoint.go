package endpoints

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	routinesdto "github.com/HiIamJeff67/notezy-backend/contracts/api/v1/routines"
	core "github.com/HiIamJeff67/notezy-backend/contracts/core/v1"
	services "github.com/HiIamJeff67/notezy-backend/internal/services/core/services"
)

type RoutineEndpointInterface interface {
	GetMyRoutineById(ctx *gin.Context)
	GetMyRoutinesByStationId(ctx *gin.Context)
	GetAllMyRoutinesByTimeRange(ctx *gin.Context)
	CreateRoutineByStationId(ctx *gin.Context)
	CreateRoutinesByStationIds(ctx *gin.Context)
	UpdateMyRoutineById(ctx *gin.Context)
	UpdateMyRoutinesByIds(ctx *gin.Context)
	LinkRoutineTagById(ctx *gin.Context)
	LinkRoutineTagsByIds(ctx *gin.Context)
	LinkRoutineItemById(ctx *gin.Context)
	LinkRoutineItemsByIds(ctx *gin.Context)
	RestoreMyRoutineById(ctx *gin.Context)
	RestoreMyRoutinesByIds(ctx *gin.Context)
	DeleteMyRoutineById(ctx *gin.Context)
	DeleteMyRoutinesByIds(ctx *gin.Context)
	HardDeleteMyRoutineById(ctx *gin.Context)
	HardDeleteMyRoutinesByIds(ctx *gin.Context)

	/* ============================== Visualization Methods ============================== */
	VisualizeMyRoutineStatusCount(ctx *gin.Context)
	VisualizeMyRoutinePeriodCount(ctx *gin.Context)
	VisualizeMyRoutineScheduledStartAtCount(ctx *gin.Context)
	VisualizeMyRoutineScheduledEndAtCount(ctx *gin.Context)

	/* ============================== GraphQL Methods ============================== */
	SearchRoutines(ctx *gin.Context)
}

type RoutineEndpoint struct {
	routineService services.RoutineServiceInterface
}

func NewRoutineEndpoint(routineService services.RoutineServiceInterface) RoutineEndpointInterface {
	return &RoutineEndpoint{routineService: routineService}
}

func (t *RoutineEndpoint) GetMyRoutineById(ctx *gin.Context) {
	request := &core.Request[routinesdto.GetMyRoutineByIdRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineService.GetMyRoutineById(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), core.Response[struct{}]{Version: core.Version, Metadata: core.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, core.Response[routinesdto.GetMyRoutineByIdResponseDto]{Version: core.Version, Metadata: core.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RoutineEndpoint) GetMyRoutinesByStationId(ctx *gin.Context) {
	request := &core.Request[routinesdto.GetMyRoutinesByStationIdRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineService.GetMyRoutinesByStationId(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), core.Response[struct{}]{Version: core.Version, Metadata: core.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, core.Response[routinesdto.GetMyRoutinesByStationIdResponseDto]{Version: core.Version, Metadata: core.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RoutineEndpoint) GetAllMyRoutinesByTimeRange(ctx *gin.Context) {
	request := &core.Request[routinesdto.GetAllMyRoutinesByTimeRangeRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineService.GetAllMyRoutinesByTimeRange(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), core.Response[struct{}]{Version: core.Version, Metadata: core.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, core.Response[routinesdto.GetAllMyRoutinesByTimeRangeResponseDto]{Version: core.Version, Metadata: core.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RoutineEndpoint) CreateRoutineByStationId(ctx *gin.Context) {
	request := &core.Request[routinesdto.CreateRoutineByStationIdRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineService.CreateRoutineByStationId(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), core.Response[struct{}]{Version: core.Version, Metadata: core.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, core.Response[routinesdto.CreateRoutineByStationIdResponseDto]{Version: core.Version, Metadata: core.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RoutineEndpoint) CreateRoutinesByStationIds(ctx *gin.Context) {
	request := &core.Request[routinesdto.CreateRoutinesByStationIdsRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineService.CreateRoutinesByStationIds(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), core.Response[struct{}]{Version: core.Version, Metadata: core.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, core.Response[routinesdto.CreateRoutinesByStationIdsResponseDto]{Version: core.Version, Metadata: core.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RoutineEndpoint) UpdateMyRoutineById(ctx *gin.Context) {
	request := &core.Request[routinesdto.UpdateMyRoutineByIdRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineService.UpdateMyRoutineById(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), core.Response[struct{}]{Version: core.Version, Metadata: core.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, core.Response[routinesdto.UpdateMyRoutineByIdResponseDto]{Version: core.Version, Metadata: core.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RoutineEndpoint) UpdateMyRoutinesByIds(ctx *gin.Context) {
	request := &core.Request[routinesdto.UpdateMyRoutinesByIdsRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineService.UpdateMyRoutinesByIds(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), core.Response[struct{}]{Version: core.Version, Metadata: core.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, core.Response[routinesdto.UpdateMyRoutinesByIdsResponseDto]{Version: core.Version, Metadata: core.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RoutineEndpoint) LinkRoutineTagById(ctx *gin.Context) {
	request := &core.Request[routinesdto.LinkRoutineTagByIdRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineService.LinkRoutineTagById(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), core.Response[struct{}]{Version: core.Version, Metadata: core.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, core.Response[routinesdto.LinkRoutineTagByIdResponseDto]{Version: core.Version, Metadata: core.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RoutineEndpoint) LinkRoutineTagsByIds(ctx *gin.Context) {
	request := &core.Request[routinesdto.LinkRoutineTagsByIdsRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineService.LinkRoutineTagsByIds(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), core.Response[struct{}]{Version: core.Version, Metadata: core.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, core.Response[routinesdto.LinkRoutineTagsByIdsResponseDto]{Version: core.Version, Metadata: core.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RoutineEndpoint) LinkRoutineItemById(ctx *gin.Context) {
	request := &core.Request[routinesdto.LinkRoutineItemByIdRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineService.LinkRoutineItemById(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), core.Response[struct{}]{Version: core.Version, Metadata: core.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, core.Response[routinesdto.LinkRoutineItemByIdResponseDto]{Version: core.Version, Metadata: core.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RoutineEndpoint) LinkRoutineItemsByIds(ctx *gin.Context) {
	request := &core.Request[routinesdto.LinkRoutineItemsByIdsRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineService.LinkRoutineItemsByIds(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), core.Response[struct{}]{Version: core.Version, Metadata: core.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, core.Response[routinesdto.LinkRoutineItemsByIdsResponseDto]{Version: core.Version, Metadata: core.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RoutineEndpoint) RestoreMyRoutineById(ctx *gin.Context) {
	request := &core.Request[routinesdto.RestoreMyRoutineByIdRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineService.RestoreMyRoutineById(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), core.Response[struct{}]{Version: core.Version, Metadata: core.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, core.Response[routinesdto.RestoreMyRoutineByIdResponseDto]{Version: core.Version, Metadata: core.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RoutineEndpoint) RestoreMyRoutinesByIds(ctx *gin.Context) {
	request := &core.Request[routinesdto.RestoreMyRoutinesByIdsRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineService.RestoreMyRoutinesByIds(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), core.Response[struct{}]{Version: core.Version, Metadata: core.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, core.Response[routinesdto.RestoreMyRoutinesByIdsResponseDto]{Version: core.Version, Metadata: core.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RoutineEndpoint) DeleteMyRoutineById(ctx *gin.Context) {
	request := &core.Request[routinesdto.DeleteMyRoutineByIdRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineService.DeleteMyRoutineById(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), core.Response[struct{}]{Version: core.Version, Metadata: core.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, core.Response[routinesdto.DeleteMyRoutineByIdResponseDto]{Version: core.Version, Metadata: core.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RoutineEndpoint) DeleteMyRoutinesByIds(ctx *gin.Context) {
	request := &core.Request[routinesdto.DeleteMyRoutinesByIdsRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineService.DeleteMyRoutinesByIds(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), core.Response[struct{}]{Version: core.Version, Metadata: core.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, core.Response[routinesdto.DeleteMyRoutinesByIdsResponseDto]{Version: core.Version, Metadata: core.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RoutineEndpoint) HardDeleteMyRoutineById(ctx *gin.Context) {
	request := &core.Request[routinesdto.HardDeleteMyRoutineByIdRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineService.HardDeleteMyRoutineById(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), core.Response[struct{}]{Version: core.Version, Metadata: core.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, core.Response[routinesdto.HardDeleteMyRoutineByIdResponseDto]{Version: core.Version, Metadata: core.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RoutineEndpoint) HardDeleteMyRoutinesByIds(ctx *gin.Context) {
	request := &core.Request[routinesdto.HardDeleteMyRoutinesByIdsRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineService.HardDeleteMyRoutinesByIds(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), core.Response[struct{}]{Version: core.Version, Metadata: core.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, core.Response[routinesdto.HardDeleteMyRoutinesByIdsResponseDto]{Version: core.Version, Metadata: core.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RoutineEndpoint) VisualizeMyRoutineStatusCount(ctx *gin.Context) {
	request := &core.Request[routinesdto.VisualizeMyRoutineStatusCountRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineService.VisualizeMyRoutineStatusCount(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), core.Response[struct{}]{Version: core.Version, Metadata: core.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, core.Response[routinesdto.VisualizeMyRoutineStatusCountResponseDto]{Version: core.Version, Metadata: core.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RoutineEndpoint) VisualizeMyRoutinePeriodCount(ctx *gin.Context) {
	request := &core.Request[routinesdto.VisualizeMyRoutinePeriodCountRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineService.VisualizeMyRoutinePeriodCount(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), core.Response[struct{}]{Version: core.Version, Metadata: core.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, core.Response[routinesdto.VisualizeMyRoutinePeriodCountResponseDto]{Version: core.Version, Metadata: core.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RoutineEndpoint) VisualizeMyRoutineScheduledStartAtCount(ctx *gin.Context) {
	request := &core.Request[routinesdto.VisualizeMyRoutineScheduledStartAtCountRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineService.VisualizeMyRoutineScheduledStartAtCount(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), core.Response[struct{}]{Version: core.Version, Metadata: core.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, core.Response[routinesdto.VisualizeMyRoutineScheduledStartAtCountResponseDto]{Version: core.Version, Metadata: core.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RoutineEndpoint) VisualizeMyRoutineScheduledEndAtCount(ctx *gin.Context) {
	request := &core.Request[routinesdto.VisualizeMyRoutineScheduledEndAtCountRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineService.VisualizeMyRoutineScheduledEndAtCount(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), core.Response[struct{}]{Version: core.Version, Metadata: core.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, core.Response[routinesdto.VisualizeMyRoutineScheduledEndAtCountResponseDto]{Version: core.Version, Metadata: core.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}
