package endpoints

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	apicontract "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api/routines"
	gatewaycontract "github.com/HiIamJeff67/notezy-backend/contracts/gateway/v1"

	routineservices "github.com/HiIamJeff67/notezy-backend/internal/core/services/routines"
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
	routineService routineservices.RoutineServiceInterface
}

func NewRoutineEndpoint(routineService routineservices.RoutineServiceInterface) RoutineEndpointInterface {
	return &RoutineEndpoint{routineService: routineService}
}

func (t *RoutineEndpoint) GetMyRoutineById(ctx *gin.Context) {
	request := &gatewaycontract.Request[apicontract.GetMyRoutineByIdRequestDto]{}
	if err := ctx.ShouldBindBodyWithJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineService.GetMyRoutineById(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), gatewaycontract.Response[struct{}]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, gatewaycontract.Response[apicontract.GetMyRoutineByIdResponseDto]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RoutineEndpoint) GetMyRoutinesByStationId(ctx *gin.Context) {
	request := &gatewaycontract.Request[apicontract.GetMyRoutinesByStationIdRequestDto]{}
	if err := ctx.ShouldBindBodyWithJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineService.GetMyRoutinesByStationId(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), gatewaycontract.Response[struct{}]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, gatewaycontract.Response[apicontract.GetMyRoutinesByStationIdResponseDto]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RoutineEndpoint) GetAllMyRoutinesByTimeRange(ctx *gin.Context) {
	request := &gatewaycontract.Request[apicontract.GetAllMyRoutinesByTimeRangeRequestDto]{}
	if err := ctx.ShouldBindBodyWithJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineService.GetAllMyRoutinesByTimeRange(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), gatewaycontract.Response[struct{}]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, gatewaycontract.Response[apicontract.GetAllMyRoutinesByTimeRangeResponseDto]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RoutineEndpoint) CreateRoutineByStationId(ctx *gin.Context) {
	request := &gatewaycontract.Request[apicontract.CreateRoutineByStationIdRequestDto]{}
	if err := ctx.ShouldBindBodyWithJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineService.CreateRoutineByStationId(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), gatewaycontract.Response[struct{}]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, gatewaycontract.Response[apicontract.CreateRoutineByStationIdResponseDto]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RoutineEndpoint) CreateRoutinesByStationIds(ctx *gin.Context) {
	request := &gatewaycontract.Request[apicontract.CreateRoutinesByStationIdsRequestDto]{}
	if err := ctx.ShouldBindBodyWithJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineService.CreateRoutinesByStationIds(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), gatewaycontract.Response[struct{}]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, gatewaycontract.Response[apicontract.CreateRoutinesByStationIdsResponseDto]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RoutineEndpoint) UpdateMyRoutineById(ctx *gin.Context) {
	request := &gatewaycontract.Request[apicontract.UpdateMyRoutineByIdRequestDto]{}
	if err := ctx.ShouldBindBodyWithJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineService.UpdateMyRoutineById(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), gatewaycontract.Response[struct{}]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, gatewaycontract.Response[apicontract.UpdateMyRoutineByIdResponseDto]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RoutineEndpoint) UpdateMyRoutinesByIds(ctx *gin.Context) {
	request := &gatewaycontract.Request[apicontract.UpdateMyRoutinesByIdsRequestDto]{}
	if err := ctx.ShouldBindBodyWithJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineService.UpdateMyRoutinesByIds(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), gatewaycontract.Response[struct{}]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, gatewaycontract.Response[apicontract.UpdateMyRoutinesByIdsResponseDto]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RoutineEndpoint) LinkRoutineTagById(ctx *gin.Context) {
	request := &gatewaycontract.Request[apicontract.LinkRoutineTagByIdRequestDto]{}
	if err := ctx.ShouldBindBodyWithJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineService.LinkRoutineTagById(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), gatewaycontract.Response[struct{}]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, gatewaycontract.Response[apicontract.LinkRoutineTagByIdResponseDto]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RoutineEndpoint) LinkRoutineTagsByIds(ctx *gin.Context) {
	request := &gatewaycontract.Request[apicontract.LinkRoutineTagsByIdsRequestDto]{}
	if err := ctx.ShouldBindBodyWithJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineService.LinkRoutineTagsByIds(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), gatewaycontract.Response[struct{}]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, gatewaycontract.Response[apicontract.LinkRoutineTagsByIdsResponseDto]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RoutineEndpoint) LinkRoutineItemById(ctx *gin.Context) {
	request := &gatewaycontract.Request[apicontract.LinkRoutineItemByIdRequestDto]{}
	if err := ctx.ShouldBindBodyWithJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineService.LinkRoutineItemById(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), gatewaycontract.Response[struct{}]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, gatewaycontract.Response[apicontract.LinkRoutineItemByIdResponseDto]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RoutineEndpoint) LinkRoutineItemsByIds(ctx *gin.Context) {
	request := &gatewaycontract.Request[apicontract.LinkRoutineItemsByIdsRequestDto]{}
	if err := ctx.ShouldBindBodyWithJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineService.LinkRoutineItemsByIds(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), gatewaycontract.Response[struct{}]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, gatewaycontract.Response[apicontract.LinkRoutineItemsByIdsResponseDto]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RoutineEndpoint) RestoreMyRoutineById(ctx *gin.Context) {
	request := &gatewaycontract.Request[apicontract.RestoreMyRoutineByIdRequestDto]{}
	if err := ctx.ShouldBindBodyWithJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineService.RestoreMyRoutineById(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), gatewaycontract.Response[struct{}]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, gatewaycontract.Response[apicontract.RestoreMyRoutineByIdResponseDto]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RoutineEndpoint) RestoreMyRoutinesByIds(ctx *gin.Context) {
	request := &gatewaycontract.Request[apicontract.RestoreMyRoutinesByIdsRequestDto]{}
	if err := ctx.ShouldBindBodyWithJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineService.RestoreMyRoutinesByIds(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), gatewaycontract.Response[struct{}]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, gatewaycontract.Response[apicontract.RestoreMyRoutinesByIdsResponseDto]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RoutineEndpoint) DeleteMyRoutineById(ctx *gin.Context) {
	request := &gatewaycontract.Request[apicontract.DeleteMyRoutineByIdRequestDto]{}
	if err := ctx.ShouldBindBodyWithJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineService.DeleteMyRoutineById(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), gatewaycontract.Response[struct{}]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, gatewaycontract.Response[apicontract.DeleteMyRoutineByIdResponseDto]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RoutineEndpoint) DeleteMyRoutinesByIds(ctx *gin.Context) {
	request := &gatewaycontract.Request[apicontract.DeleteMyRoutinesByIdsRequestDto]{}
	if err := ctx.ShouldBindBodyWithJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineService.DeleteMyRoutinesByIds(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), gatewaycontract.Response[struct{}]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, gatewaycontract.Response[apicontract.DeleteMyRoutinesByIdsResponseDto]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RoutineEndpoint) HardDeleteMyRoutineById(ctx *gin.Context) {
	request := &gatewaycontract.Request[apicontract.HardDeleteMyRoutineByIdRequestDto]{}
	if err := ctx.ShouldBindBodyWithJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineService.HardDeleteMyRoutineById(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), gatewaycontract.Response[struct{}]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, gatewaycontract.Response[apicontract.HardDeleteMyRoutineByIdResponseDto]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RoutineEndpoint) HardDeleteMyRoutinesByIds(ctx *gin.Context) {
	request := &gatewaycontract.Request[apicontract.HardDeleteMyRoutinesByIdsRequestDto]{}
	if err := ctx.ShouldBindBodyWithJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineService.HardDeleteMyRoutinesByIds(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), gatewaycontract.Response[struct{}]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, gatewaycontract.Response[apicontract.HardDeleteMyRoutinesByIdsResponseDto]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RoutineEndpoint) VisualizeMyRoutineStatusCount(ctx *gin.Context) {
	request := &gatewaycontract.Request[apicontract.VisualizeMyRoutineStatusCountRequestDto]{}
	if err := ctx.ShouldBindBodyWithJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineService.VisualizeMyRoutineStatusCount(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), gatewaycontract.Response[struct{}]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, gatewaycontract.Response[apicontract.VisualizeMyRoutineStatusCountResponseDto]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RoutineEndpoint) VisualizeMyRoutinePeriodCount(ctx *gin.Context) {
	request := &gatewaycontract.Request[apicontract.VisualizeMyRoutinePeriodCountRequestDto]{}
	if err := ctx.ShouldBindBodyWithJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineService.VisualizeMyRoutinePeriodCount(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), gatewaycontract.Response[struct{}]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, gatewaycontract.Response[apicontract.VisualizeMyRoutinePeriodCountResponseDto]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RoutineEndpoint) VisualizeMyRoutineScheduledStartAtCount(ctx *gin.Context) {
	request := &gatewaycontract.Request[apicontract.VisualizeMyRoutineScheduledStartAtCountRequestDto]{}
	if err := ctx.ShouldBindBodyWithJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineService.VisualizeMyRoutineScheduledStartAtCount(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), gatewaycontract.Response[struct{}]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, gatewaycontract.Response[apicontract.VisualizeMyRoutineScheduledStartAtCountResponseDto]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RoutineEndpoint) VisualizeMyRoutineScheduledEndAtCount(ctx *gin.Context) {
	request := &gatewaycontract.Request[apicontract.VisualizeMyRoutineScheduledEndAtCountRequestDto]{}
	if err := ctx.ShouldBindBodyWithJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineService.VisualizeMyRoutineScheduledEndAtCount(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), gatewaycontract.Response[struct{}]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, gatewaycontract.Response[apicontract.VisualizeMyRoutineScheduledEndAtCountResponseDto]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}
