package endpoints

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	gatewaycontract "github.com/HiIamJeff67/notezy-backend/contracts/gateway/v1"
	routinesdto "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api/routines"
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
	request := &gatewaycontract.Request[routinesdto.GetMyRoutineByIdRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineService.GetMyRoutineById(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), gatewaycontract.Response[struct{}]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, gatewaycontract.Response[routinesdto.GetMyRoutineByIdResponseDto]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RoutineEndpoint) GetMyRoutinesByStationId(ctx *gin.Context) {
	request := &gatewaycontract.Request[routinesdto.GetMyRoutinesByStationIdRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineService.GetMyRoutinesByStationId(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), gatewaycontract.Response[struct{}]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, gatewaycontract.Response[routinesdto.GetMyRoutinesByStationIdResponseDto]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RoutineEndpoint) GetAllMyRoutinesByTimeRange(ctx *gin.Context) {
	request := &gatewaycontract.Request[routinesdto.GetAllMyRoutinesByTimeRangeRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineService.GetAllMyRoutinesByTimeRange(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), gatewaycontract.Response[struct{}]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, gatewaycontract.Response[routinesdto.GetAllMyRoutinesByTimeRangeResponseDto]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RoutineEndpoint) CreateRoutineByStationId(ctx *gin.Context) {
	request := &gatewaycontract.Request[routinesdto.CreateRoutineByStationIdRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineService.CreateRoutineByStationId(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), gatewaycontract.Response[struct{}]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, gatewaycontract.Response[routinesdto.CreateRoutineByStationIdResponseDto]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RoutineEndpoint) CreateRoutinesByStationIds(ctx *gin.Context) {
	request := &gatewaycontract.Request[routinesdto.CreateRoutinesByStationIdsRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineService.CreateRoutinesByStationIds(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), gatewaycontract.Response[struct{}]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, gatewaycontract.Response[routinesdto.CreateRoutinesByStationIdsResponseDto]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RoutineEndpoint) UpdateMyRoutineById(ctx *gin.Context) {
	request := &gatewaycontract.Request[routinesdto.UpdateMyRoutineByIdRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineService.UpdateMyRoutineById(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), gatewaycontract.Response[struct{}]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, gatewaycontract.Response[routinesdto.UpdateMyRoutineByIdResponseDto]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RoutineEndpoint) UpdateMyRoutinesByIds(ctx *gin.Context) {
	request := &gatewaycontract.Request[routinesdto.UpdateMyRoutinesByIdsRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineService.UpdateMyRoutinesByIds(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), gatewaycontract.Response[struct{}]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, gatewaycontract.Response[routinesdto.UpdateMyRoutinesByIdsResponseDto]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RoutineEndpoint) LinkRoutineTagById(ctx *gin.Context) {
	request := &gatewaycontract.Request[routinesdto.LinkRoutineTagByIdRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineService.LinkRoutineTagById(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), gatewaycontract.Response[struct{}]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, gatewaycontract.Response[routinesdto.LinkRoutineTagByIdResponseDto]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RoutineEndpoint) LinkRoutineTagsByIds(ctx *gin.Context) {
	request := &gatewaycontract.Request[routinesdto.LinkRoutineTagsByIdsRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineService.LinkRoutineTagsByIds(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), gatewaycontract.Response[struct{}]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, gatewaycontract.Response[routinesdto.LinkRoutineTagsByIdsResponseDto]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RoutineEndpoint) LinkRoutineItemById(ctx *gin.Context) {
	request := &gatewaycontract.Request[routinesdto.LinkRoutineItemByIdRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineService.LinkRoutineItemById(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), gatewaycontract.Response[struct{}]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, gatewaycontract.Response[routinesdto.LinkRoutineItemByIdResponseDto]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RoutineEndpoint) LinkRoutineItemsByIds(ctx *gin.Context) {
	request := &gatewaycontract.Request[routinesdto.LinkRoutineItemsByIdsRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineService.LinkRoutineItemsByIds(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), gatewaycontract.Response[struct{}]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, gatewaycontract.Response[routinesdto.LinkRoutineItemsByIdsResponseDto]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RoutineEndpoint) RestoreMyRoutineById(ctx *gin.Context) {
	request := &gatewaycontract.Request[routinesdto.RestoreMyRoutineByIdRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineService.RestoreMyRoutineById(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), gatewaycontract.Response[struct{}]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, gatewaycontract.Response[routinesdto.RestoreMyRoutineByIdResponseDto]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RoutineEndpoint) RestoreMyRoutinesByIds(ctx *gin.Context) {
	request := &gatewaycontract.Request[routinesdto.RestoreMyRoutinesByIdsRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineService.RestoreMyRoutinesByIds(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), gatewaycontract.Response[struct{}]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, gatewaycontract.Response[routinesdto.RestoreMyRoutinesByIdsResponseDto]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RoutineEndpoint) DeleteMyRoutineById(ctx *gin.Context) {
	request := &gatewaycontract.Request[routinesdto.DeleteMyRoutineByIdRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineService.DeleteMyRoutineById(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), gatewaycontract.Response[struct{}]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, gatewaycontract.Response[routinesdto.DeleteMyRoutineByIdResponseDto]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RoutineEndpoint) DeleteMyRoutinesByIds(ctx *gin.Context) {
	request := &gatewaycontract.Request[routinesdto.DeleteMyRoutinesByIdsRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineService.DeleteMyRoutinesByIds(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), gatewaycontract.Response[struct{}]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, gatewaycontract.Response[routinesdto.DeleteMyRoutinesByIdsResponseDto]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RoutineEndpoint) HardDeleteMyRoutineById(ctx *gin.Context) {
	request := &gatewaycontract.Request[routinesdto.HardDeleteMyRoutineByIdRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineService.HardDeleteMyRoutineById(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), gatewaycontract.Response[struct{}]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, gatewaycontract.Response[routinesdto.HardDeleteMyRoutineByIdResponseDto]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RoutineEndpoint) HardDeleteMyRoutinesByIds(ctx *gin.Context) {
	request := &gatewaycontract.Request[routinesdto.HardDeleteMyRoutinesByIdsRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineService.HardDeleteMyRoutinesByIds(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), gatewaycontract.Response[struct{}]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, gatewaycontract.Response[routinesdto.HardDeleteMyRoutinesByIdsResponseDto]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RoutineEndpoint) VisualizeMyRoutineStatusCount(ctx *gin.Context) {
	request := &gatewaycontract.Request[routinesdto.VisualizeMyRoutineStatusCountRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineService.VisualizeMyRoutineStatusCount(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), gatewaycontract.Response[struct{}]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, gatewaycontract.Response[routinesdto.VisualizeMyRoutineStatusCountResponseDto]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RoutineEndpoint) VisualizeMyRoutinePeriodCount(ctx *gin.Context) {
	request := &gatewaycontract.Request[routinesdto.VisualizeMyRoutinePeriodCountRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineService.VisualizeMyRoutinePeriodCount(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), gatewaycontract.Response[struct{}]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, gatewaycontract.Response[routinesdto.VisualizeMyRoutinePeriodCountResponseDto]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RoutineEndpoint) VisualizeMyRoutineScheduledStartAtCount(ctx *gin.Context) {
	request := &gatewaycontract.Request[routinesdto.VisualizeMyRoutineScheduledStartAtCountRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineService.VisualizeMyRoutineScheduledStartAtCount(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), gatewaycontract.Response[struct{}]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, gatewaycontract.Response[routinesdto.VisualizeMyRoutineScheduledStartAtCountResponseDto]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RoutineEndpoint) VisualizeMyRoutineScheduledEndAtCount(ctx *gin.Context) {
	request := &gatewaycontract.Request[routinesdto.VisualizeMyRoutineScheduledEndAtCountRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineService.VisualizeMyRoutineScheduledEndAtCount(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), gatewaycontract.Response[struct{}]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, gatewaycontract.Response[routinesdto.VisualizeMyRoutineScheduledEndAtCountResponseDto]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}
