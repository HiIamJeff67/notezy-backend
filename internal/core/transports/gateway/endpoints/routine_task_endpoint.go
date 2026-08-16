package endpoints

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	apicontract "github.com/HiIamJeff67/notegic-backend/contracts/core/v1/api/routine-tasks"
	gatewaycontract "github.com/HiIamJeff67/notegic-backend/contracts/gateway/v1"

	routineservices "github.com/HiIamJeff67/notegic-backend/internal/core/services/routines"
)

type RoutineTaskEndpointInterface interface {
	GetMyRoutineTaskById(ctx *gin.Context)
	GetAllMyRoutineTasksByRoutineIds(ctx *gin.Context)
	GetAllMyRoutineTasks(ctx *gin.Context)
	CreateRoutineTaskByRoutineId(ctx *gin.Context)
	UpdateMyRoutineTaskById(ctx *gin.Context)
	PauseMyRoutineTaskById(ctx *gin.Context)
	ResumeMyRoutineTaskById(ctx *gin.Context)
	HardDeleteMyRoutineTaskById(ctx *gin.Context)
	HardDeleteMyRoutineTasksByIds(ctx *gin.Context)

	/* ============================== Visualization Methods ============================== */
	VisualizeMyRoutineTaskStatusCount(ctx *gin.Context)
	VisualizeMyRoutineTaskPurposeCount(ctx *gin.Context)
	VisualizeMyRoutineTaskScheduledAtCount(ctx *gin.Context)
	VisualizeMyRoutineTaskActualStartedAtCount(ctx *gin.Context)
	VisualizeMyRoutineTaskActualEndedAtCount(ctx *gin.Context)

	/* ============================== GraphQL Methods ============================== */
	SearchRoutineTasks(ctx *gin.Context)
}

type RoutineTaskEndpoint struct {
	routineTaskService routineservices.RoutineTaskServiceInterface
}

func NewRoutineTaskEndpoint(routineTaskService routineservices.RoutineTaskServiceInterface) RoutineTaskEndpointInterface {
	return &RoutineTaskEndpoint{routineTaskService: routineTaskService}
}

func (t *RoutineTaskEndpoint) GetMyRoutineTaskById(ctx *gin.Context) {
	request := &gatewaycontract.Request[apicontract.GetMyRoutineTaskByIdRequestDto]{}
	if err := ctx.ShouldBindBodyWithJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineTaskService.GetMyRoutineTaskById(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), gatewaycontract.Response[struct{}]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, gatewaycontract.Response[apicontract.GetMyRoutineTaskByIdResponseDto]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RoutineTaskEndpoint) GetAllMyRoutineTasksByRoutineIds(ctx *gin.Context) {
	request := &gatewaycontract.Request[apicontract.GetAllMyRoutineTasksByRoutineIdsRequestDto]{}
	if err := ctx.ShouldBindBodyWithJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineTaskService.GetAllMyRoutineTasksByRoutineIds(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), gatewaycontract.Response[struct{}]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, gatewaycontract.Response[apicontract.GetAllMyRoutineTasksByRoutineIdsResponseDto]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RoutineTaskEndpoint) GetAllMyRoutineTasks(ctx *gin.Context) {
	request := &gatewaycontract.Request[apicontract.GetAllMyRoutineTasksRequestDto]{}
	if err := ctx.ShouldBindBodyWithJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineTaskService.GetAllMyRoutineTasks(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), gatewaycontract.Response[struct{}]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, gatewaycontract.Response[apicontract.GetAllMyRoutineTasksResponseDto]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RoutineTaskEndpoint) CreateRoutineTaskByRoutineId(ctx *gin.Context) {
	request := &gatewaycontract.Request[apicontract.CreateRoutineTaskByRoutineIdRequestDto]{}
	if err := ctx.ShouldBindBodyWithJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineTaskService.CreateRoutineTaskByRoutineId(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), gatewaycontract.Response[struct{}]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, gatewaycontract.Response[apicontract.CreateRoutineTaskByRoutineIdResponseDto]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RoutineTaskEndpoint) UpdateMyRoutineTaskById(ctx *gin.Context) {
	request := &gatewaycontract.Request[apicontract.UpdateMyRoutineTaskByIdRequestDto]{}
	if err := ctx.ShouldBindBodyWithJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineTaskService.UpdateMyRoutineTaskById(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), gatewaycontract.Response[struct{}]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, gatewaycontract.Response[apicontract.UpdateMyRoutineTaskByIdResponseDto]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RoutineTaskEndpoint) PauseMyRoutineTaskById(ctx *gin.Context) {
	request := &gatewaycontract.Request[apicontract.PauseMyRoutineTaskByIdRequestDto]{}
	if err := ctx.ShouldBindBodyWithJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineTaskService.PauseMyRoutineTaskById(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), gatewaycontract.Response[struct{}]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, gatewaycontract.Response[apicontract.PauseMyRoutineTaskByIdResponseDto]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RoutineTaskEndpoint) ResumeMyRoutineTaskById(ctx *gin.Context) {
	request := &gatewaycontract.Request[apicontract.ResumeMyRoutineTaskByIdRequestDto]{}
	if err := ctx.ShouldBindBodyWithJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineTaskService.ResumeMyRoutineTaskById(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), gatewaycontract.Response[struct{}]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, gatewaycontract.Response[apicontract.ResumeMyRoutineTaskByIdResponseDto]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RoutineTaskEndpoint) HardDeleteMyRoutineTaskById(ctx *gin.Context) {
	request := &gatewaycontract.Request[apicontract.HardDeleteMyRoutineTaskByIdRequestDto]{}
	if err := ctx.ShouldBindBodyWithJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineTaskService.HardDeleteMyRoutineTaskById(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), gatewaycontract.Response[struct{}]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, gatewaycontract.Response[apicontract.HardDeleteMyRoutineTaskByIdResponseDto]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RoutineTaskEndpoint) HardDeleteMyRoutineTasksByIds(ctx *gin.Context) {
	request := &gatewaycontract.Request[apicontract.HardDeleteMyRoutineTasksByIdsRequestDto]{}
	if err := ctx.ShouldBindBodyWithJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineTaskService.HardDeleteMyRoutineTasksByIds(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), gatewaycontract.Response[struct{}]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, gatewaycontract.Response[apicontract.HardDeleteMyRoutineTasksByIdsResponseDto]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RoutineTaskEndpoint) VisualizeMyRoutineTaskStatusCount(ctx *gin.Context) {
	request := &gatewaycontract.Request[apicontract.VisualizeMyRoutineTaskStatusCountRequestDto]{}
	if err := ctx.ShouldBindBodyWithJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineTaskService.VisualizeMyRoutineTaskStatusCount(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), gatewaycontract.Response[struct{}]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, gatewaycontract.Response[apicontract.VisualizeMyRoutineTaskStatusCountResponseDto]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RoutineTaskEndpoint) VisualizeMyRoutineTaskPurposeCount(ctx *gin.Context) {
	request := &gatewaycontract.Request[apicontract.VisualizeMyRoutineTaskPurposeCountRequestDto]{}
	if err := ctx.ShouldBindBodyWithJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineTaskService.VisualizeMyRoutineTaskPurposeCount(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), gatewaycontract.Response[struct{}]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, gatewaycontract.Response[apicontract.VisualizeMyRoutineTaskPurposeCountResponseDto]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RoutineTaskEndpoint) VisualizeMyRoutineTaskScheduledAtCount(ctx *gin.Context) {
	request := &gatewaycontract.Request[apicontract.VisualizeMyRoutineTaskScheduledAtCountRequestDto]{}
	if err := ctx.ShouldBindBodyWithJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineTaskService.VisualizeMyRoutineTaskScheduledAtCount(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), gatewaycontract.Response[struct{}]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, gatewaycontract.Response[apicontract.VisualizeMyRoutineTaskScheduledAtCountResponseDto]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RoutineTaskEndpoint) VisualizeMyRoutineTaskActualStartedAtCount(ctx *gin.Context) {
	request := &gatewaycontract.Request[apicontract.VisualizeMyRoutineTaskActualStartedAtCountRequestDto]{}
	if err := ctx.ShouldBindBodyWithJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineTaskService.VisualizeMyRoutineTaskActualStartedAtCount(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), gatewaycontract.Response[struct{}]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, gatewaycontract.Response[apicontract.VisualizeMyRoutineTaskActualStartedAtCountResponseDto]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RoutineTaskEndpoint) VisualizeMyRoutineTaskActualEndedAtCount(ctx *gin.Context) {
	request := &gatewaycontract.Request[apicontract.VisualizeMyRoutineTaskActualEndedAtCountRequestDto]{}
	if err := ctx.ShouldBindBodyWithJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineTaskService.VisualizeMyRoutineTaskActualEndedAtCount(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), gatewaycontract.Response[struct{}]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, gatewaycontract.Response[apicontract.VisualizeMyRoutineTaskActualEndedAtCountResponseDto]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}
