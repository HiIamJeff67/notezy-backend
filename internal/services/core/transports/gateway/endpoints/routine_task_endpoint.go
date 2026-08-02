package endpoints

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	core "github.com/HiIamJeff67/notezy-backend/contracts/core/v1"
	routinetasksdto "github.com/HiIamJeff67/notezy-backend/contracts/gateway/v1/api/routine-tasks"
	services "github.com/HiIamJeff67/notezy-backend/internal/services/core/services"
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
	routineTaskService services.RoutineTaskServiceInterface
}

func NewRoutineTaskEndpoint(routineTaskService services.RoutineTaskServiceInterface) RoutineTaskEndpointInterface {
	return &RoutineTaskEndpoint{routineTaskService: routineTaskService}
}

func (t *RoutineTaskEndpoint) GetMyRoutineTaskById(ctx *gin.Context) {
	request := &core.Request[routinetasksdto.GetMyRoutineTaskByIdRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineTaskService.GetMyRoutineTaskById(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), core.Response[struct{}]{Version: core.Version, Metadata: core.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, core.Response[routinetasksdto.GetMyRoutineTaskByIdResponseDto]{Version: core.Version, Metadata: core.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RoutineTaskEndpoint) GetAllMyRoutineTasksByRoutineIds(ctx *gin.Context) {
	request := &core.Request[routinetasksdto.GetAllMyRoutineTasksByRoutineIdsRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineTaskService.GetAllMyRoutineTasksByRoutineIds(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), core.Response[struct{}]{Version: core.Version, Metadata: core.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, core.Response[routinetasksdto.GetAllMyRoutineTasksByRoutineIdsResponseDto]{Version: core.Version, Metadata: core.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RoutineTaskEndpoint) GetAllMyRoutineTasks(ctx *gin.Context) {
	request := &core.Request[routinetasksdto.GetAllMyRoutineTasksRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineTaskService.GetAllMyRoutineTasks(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), core.Response[struct{}]{Version: core.Version, Metadata: core.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, core.Response[routinetasksdto.GetAllMyRoutineTasksResponseDto]{Version: core.Version, Metadata: core.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RoutineTaskEndpoint) CreateRoutineTaskByRoutineId(ctx *gin.Context) {
	request := &core.Request[routinetasksdto.CreateRoutineTaskByRoutineIdRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineTaskService.CreateRoutineTaskByRoutineId(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), core.Response[struct{}]{Version: core.Version, Metadata: core.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, core.Response[routinetasksdto.CreateRoutineTaskByRoutineIdResponseDto]{Version: core.Version, Metadata: core.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RoutineTaskEndpoint) UpdateMyRoutineTaskById(ctx *gin.Context) {
	request := &core.Request[routinetasksdto.UpdateMyRoutineTaskByIdRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineTaskService.UpdateMyRoutineTaskById(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), core.Response[struct{}]{Version: core.Version, Metadata: core.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, core.Response[routinetasksdto.UpdateMyRoutineTaskByIdResponseDto]{Version: core.Version, Metadata: core.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RoutineTaskEndpoint) PauseMyRoutineTaskById(ctx *gin.Context) {
	request := &core.Request[routinetasksdto.PauseMyRoutineTaskByIdRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineTaskService.PauseMyRoutineTaskById(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), core.Response[struct{}]{Version: core.Version, Metadata: core.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, core.Response[routinetasksdto.PauseMyRoutineTaskByIdResponseDto]{Version: core.Version, Metadata: core.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RoutineTaskEndpoint) ResumeMyRoutineTaskById(ctx *gin.Context) {
	request := &core.Request[routinetasksdto.ResumeMyRoutineTaskByIdRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineTaskService.ResumeMyRoutineTaskById(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), core.Response[struct{}]{Version: core.Version, Metadata: core.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, core.Response[routinetasksdto.ResumeMyRoutineTaskByIdResponseDto]{Version: core.Version, Metadata: core.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RoutineTaskEndpoint) HardDeleteMyRoutineTaskById(ctx *gin.Context) {
	request := &core.Request[routinetasksdto.HardDeleteMyRoutineTaskByIdRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineTaskService.HardDeleteMyRoutineTaskById(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), core.Response[struct{}]{Version: core.Version, Metadata: core.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, core.Response[routinetasksdto.HardDeleteMyRoutineTaskByIdResponseDto]{Version: core.Version, Metadata: core.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RoutineTaskEndpoint) HardDeleteMyRoutineTasksByIds(ctx *gin.Context) {
	request := &core.Request[routinetasksdto.HardDeleteMyRoutineTasksByIdsRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineTaskService.HardDeleteMyRoutineTasksByIds(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), core.Response[struct{}]{Version: core.Version, Metadata: core.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, core.Response[routinetasksdto.HardDeleteMyRoutineTasksByIdsResponseDto]{Version: core.Version, Metadata: core.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RoutineTaskEndpoint) VisualizeMyRoutineTaskStatusCount(ctx *gin.Context) {
	request := &core.Request[routinetasksdto.VisualizeMyRoutineTaskStatusCountRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineTaskService.VisualizeMyRoutineTaskStatusCount(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), core.Response[struct{}]{Version: core.Version, Metadata: core.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, core.Response[routinetasksdto.VisualizeMyRoutineTaskStatusCountResponseDto]{Version: core.Version, Metadata: core.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RoutineTaskEndpoint) VisualizeMyRoutineTaskPurposeCount(ctx *gin.Context) {
	request := &core.Request[routinetasksdto.VisualizeMyRoutineTaskPurposeCountRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineTaskService.VisualizeMyRoutineTaskPurposeCount(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), core.Response[struct{}]{Version: core.Version, Metadata: core.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, core.Response[routinetasksdto.VisualizeMyRoutineTaskPurposeCountResponseDto]{Version: core.Version, Metadata: core.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RoutineTaskEndpoint) VisualizeMyRoutineTaskScheduledAtCount(ctx *gin.Context) {
	request := &core.Request[routinetasksdto.VisualizeMyRoutineTaskScheduledAtCountRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineTaskService.VisualizeMyRoutineTaskScheduledAtCount(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), core.Response[struct{}]{Version: core.Version, Metadata: core.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, core.Response[routinetasksdto.VisualizeMyRoutineTaskScheduledAtCountResponseDto]{Version: core.Version, Metadata: core.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RoutineTaskEndpoint) VisualizeMyRoutineTaskActualStartedAtCount(ctx *gin.Context) {
	request := &core.Request[routinetasksdto.VisualizeMyRoutineTaskActualStartedAtCountRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineTaskService.VisualizeMyRoutineTaskActualStartedAtCount(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), core.Response[struct{}]{Version: core.Version, Metadata: core.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, core.Response[routinetasksdto.VisualizeMyRoutineTaskActualStartedAtCountResponseDto]{Version: core.Version, Metadata: core.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RoutineTaskEndpoint) VisualizeMyRoutineTaskActualEndedAtCount(ctx *gin.Context) {
	request := &core.Request[routinetasksdto.VisualizeMyRoutineTaskActualEndedAtCountRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineTaskService.VisualizeMyRoutineTaskActualEndedAtCount(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), core.Response[struct{}]{Version: core.Version, Metadata: core.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, core.Response[routinetasksdto.VisualizeMyRoutineTaskActualEndedAtCountResponseDto]{Version: core.Version, Metadata: core.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}
