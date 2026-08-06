package endpoints

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	routinetasksdto "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api/routine-tasks"
	gatewaycontract "github.com/HiIamJeff67/notezy-backend/contracts/gateway/v1"

	routineservices "github.com/HiIamJeff67/notezy-backend/internal/core/services/routines"
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
	request := &gatewaycontract.Request[routinetasksdto.GetMyRoutineTaskByIdRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineTaskService.GetMyRoutineTaskById(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), gatewaycontract.Response[struct{}]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, gatewaycontract.Response[routinetasksdto.GetMyRoutineTaskByIdResponseDto]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RoutineTaskEndpoint) GetAllMyRoutineTasksByRoutineIds(ctx *gin.Context) {
	request := &gatewaycontract.Request[routinetasksdto.GetAllMyRoutineTasksByRoutineIdsRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineTaskService.GetAllMyRoutineTasksByRoutineIds(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), gatewaycontract.Response[struct{}]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, gatewaycontract.Response[routinetasksdto.GetAllMyRoutineTasksByRoutineIdsResponseDto]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RoutineTaskEndpoint) GetAllMyRoutineTasks(ctx *gin.Context) {
	request := &gatewaycontract.Request[routinetasksdto.GetAllMyRoutineTasksRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineTaskService.GetAllMyRoutineTasks(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), gatewaycontract.Response[struct{}]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, gatewaycontract.Response[routinetasksdto.GetAllMyRoutineTasksResponseDto]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RoutineTaskEndpoint) CreateRoutineTaskByRoutineId(ctx *gin.Context) {
	request := &gatewaycontract.Request[routinetasksdto.CreateRoutineTaskByRoutineIdRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineTaskService.CreateRoutineTaskByRoutineId(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), gatewaycontract.Response[struct{}]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, gatewaycontract.Response[routinetasksdto.CreateRoutineTaskByRoutineIdResponseDto]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RoutineTaskEndpoint) UpdateMyRoutineTaskById(ctx *gin.Context) {
	request := &gatewaycontract.Request[routinetasksdto.UpdateMyRoutineTaskByIdRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineTaskService.UpdateMyRoutineTaskById(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), gatewaycontract.Response[struct{}]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, gatewaycontract.Response[routinetasksdto.UpdateMyRoutineTaskByIdResponseDto]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RoutineTaskEndpoint) PauseMyRoutineTaskById(ctx *gin.Context) {
	request := &gatewaycontract.Request[routinetasksdto.PauseMyRoutineTaskByIdRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineTaskService.PauseMyRoutineTaskById(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), gatewaycontract.Response[struct{}]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, gatewaycontract.Response[routinetasksdto.PauseMyRoutineTaskByIdResponseDto]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RoutineTaskEndpoint) ResumeMyRoutineTaskById(ctx *gin.Context) {
	request := &gatewaycontract.Request[routinetasksdto.ResumeMyRoutineTaskByIdRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineTaskService.ResumeMyRoutineTaskById(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), gatewaycontract.Response[struct{}]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, gatewaycontract.Response[routinetasksdto.ResumeMyRoutineTaskByIdResponseDto]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RoutineTaskEndpoint) HardDeleteMyRoutineTaskById(ctx *gin.Context) {
	request := &gatewaycontract.Request[routinetasksdto.HardDeleteMyRoutineTaskByIdRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineTaskService.HardDeleteMyRoutineTaskById(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), gatewaycontract.Response[struct{}]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, gatewaycontract.Response[routinetasksdto.HardDeleteMyRoutineTaskByIdResponseDto]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RoutineTaskEndpoint) HardDeleteMyRoutineTasksByIds(ctx *gin.Context) {
	request := &gatewaycontract.Request[routinetasksdto.HardDeleteMyRoutineTasksByIdsRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineTaskService.HardDeleteMyRoutineTasksByIds(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), gatewaycontract.Response[struct{}]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, gatewaycontract.Response[routinetasksdto.HardDeleteMyRoutineTasksByIdsResponseDto]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RoutineTaskEndpoint) VisualizeMyRoutineTaskStatusCount(ctx *gin.Context) {
	request := &gatewaycontract.Request[routinetasksdto.VisualizeMyRoutineTaskStatusCountRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineTaskService.VisualizeMyRoutineTaskStatusCount(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), gatewaycontract.Response[struct{}]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, gatewaycontract.Response[routinetasksdto.VisualizeMyRoutineTaskStatusCountResponseDto]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RoutineTaskEndpoint) VisualizeMyRoutineTaskPurposeCount(ctx *gin.Context) {
	request := &gatewaycontract.Request[routinetasksdto.VisualizeMyRoutineTaskPurposeCountRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineTaskService.VisualizeMyRoutineTaskPurposeCount(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), gatewaycontract.Response[struct{}]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, gatewaycontract.Response[routinetasksdto.VisualizeMyRoutineTaskPurposeCountResponseDto]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RoutineTaskEndpoint) VisualizeMyRoutineTaskScheduledAtCount(ctx *gin.Context) {
	request := &gatewaycontract.Request[routinetasksdto.VisualizeMyRoutineTaskScheduledAtCountRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineTaskService.VisualizeMyRoutineTaskScheduledAtCount(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), gatewaycontract.Response[struct{}]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, gatewaycontract.Response[routinetasksdto.VisualizeMyRoutineTaskScheduledAtCountResponseDto]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RoutineTaskEndpoint) VisualizeMyRoutineTaskActualStartedAtCount(ctx *gin.Context) {
	request := &gatewaycontract.Request[routinetasksdto.VisualizeMyRoutineTaskActualStartedAtCountRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineTaskService.VisualizeMyRoutineTaskActualStartedAtCount(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), gatewaycontract.Response[struct{}]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, gatewaycontract.Response[routinetasksdto.VisualizeMyRoutineTaskActualStartedAtCountResponseDto]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}

func (t *RoutineTaskEndpoint) VisualizeMyRoutineTaskActualEndedAtCount(ctx *gin.Context) {
	request := &gatewaycontract.Request[routinetasksdto.VisualizeMyRoutineTaskActualEndedAtCountRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineTaskService.VisualizeMyRoutineTaskActualEndedAtCount(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), gatewaycontract.Response[struct{}]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: struct{}{}, Exception: publicException})
		return
	}
	ctx.JSON(http.StatusOK, gatewaycontract.Response[routinetasksdto.VisualizeMyRoutineTaskActualEndedAtCountResponseDto]{Version: gatewaycontract.Version, Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()}, Data: *responseDto})
}
