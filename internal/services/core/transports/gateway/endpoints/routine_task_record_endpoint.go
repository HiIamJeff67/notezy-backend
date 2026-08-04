package endpoints

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	gatewaycontract "github.com/HiIamJeff67/notezy-backend/contracts/gateway/v1"
	routinetaskrecordsdto "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api/routine-task-records"
	services "github.com/HiIamJeff67/notezy-backend/internal/services/core/services"
)

type RoutineTaskRecordEndpointInterface interface {
	GetAllMyRoutineTaskRecordsByRoutineTaskId(ctx *gin.Context)

	/* ============================== Visualization Methods ============================== */
	VisualizeMyRoutineTaskRecordStatusCount(ctx *gin.Context)
	VisualizeMyRoutineTaskRecordPurposeCount(ctx *gin.Context)
	VisualizeMyRoutineTaskRecordScheduledAtCount(ctx *gin.Context)
	VisualizeMyRoutineTaskRecordActualStartedAtCount(ctx *gin.Context)
	VisualizeMyRoutineTaskRecordActualEndedAtCount(ctx *gin.Context)

	/* ============================== GraphQL Methods ============================== */
	SearchRoutineTaskRecords(ctx *gin.Context)
}

type RoutineTaskRecordEndpoint struct {
	routineTaskRecordService services.RoutineTaskRecordServiceInterface
}

func NewRoutineTaskRecordEndpoint(
	routineTaskRecordService services.RoutineTaskRecordServiceInterface,
) RoutineTaskRecordEndpointInterface {
	return &RoutineTaskRecordEndpoint{
		routineTaskRecordService: routineTaskRecordService,
	}
}

func (t *RoutineTaskRecordEndpoint) GetAllMyRoutineTaskRecordsByRoutineTaskId(ctx *gin.Context) {
	request := &gatewaycontract.Request[routinetaskrecordsdto.GetAllMyRoutineTaskRecordsByRoutineTaskIdRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineTaskRecordService.GetAllMyRoutineTaskRecordsByRoutineTaskId(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), gatewaycontract.Response[struct{}]{
			Version: gatewaycontract.Version,
			Metadata: gatewaycontract.ResponseMetadata{
				RequestId:   request.Metadata.RequestId,
				RespondedAt: time.Now(),
			},
			Data:      struct{}{},
			Exception: publicException,
		})
		return
	}

	ctx.JSON(http.StatusOK, gatewaycontract.Response[routinetaskrecordsdto.GetAllMyRoutineTaskRecordsByRoutineTaskIdResponseDto]{
		Version:  gatewaycontract.Version,
		Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()},
		Data:     *responseDto,
	})
}

/* ============================== Visualization Methods ============================== */

func (t *RoutineTaskRecordEndpoint) VisualizeMyRoutineTaskRecordStatusCount(ctx *gin.Context) {
	request := &gatewaycontract.Request[routinetaskrecordsdto.VisualizeMyRoutineTaskRecordStatusCountRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineTaskRecordService.VisualizeMyRoutineTaskRecordStatusCount(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), gatewaycontract.Response[struct{}]{
			Version: gatewaycontract.Version,
			Metadata: gatewaycontract.ResponseMetadata{
				RequestId:   request.Metadata.RequestId,
				RespondedAt: time.Now(),
			},
			Data:      struct{}{},
			Exception: publicException,
		})
		return
	}

	ctx.JSON(http.StatusOK, gatewaycontract.Response[routinetaskrecordsdto.VisualizeMyRoutineTaskRecordStatusCountResponseDto]{
		Version:  gatewaycontract.Version,
		Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()},
		Data:     *responseDto,
	})
}

func (t *RoutineTaskRecordEndpoint) VisualizeMyRoutineTaskRecordPurposeCount(ctx *gin.Context) {
	request := &gatewaycontract.Request[routinetaskrecordsdto.VisualizeMyRoutineTaskRecordPurposeCountRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineTaskRecordService.VisualizeMyRoutineTaskRecordPurposeCount(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), gatewaycontract.Response[struct{}]{
			Version: gatewaycontract.Version,
			Metadata: gatewaycontract.ResponseMetadata{
				RequestId:   request.Metadata.RequestId,
				RespondedAt: time.Now(),
			},
			Data:      struct{}{},
			Exception: publicException,
		})
		return
	}

	ctx.JSON(http.StatusOK, gatewaycontract.Response[routinetaskrecordsdto.VisualizeMyRoutineTaskRecordPurposeCountResponseDto]{
		Version:  gatewaycontract.Version,
		Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()},
		Data:     *responseDto,
	})
}

func (t *RoutineTaskRecordEndpoint) VisualizeMyRoutineTaskRecordScheduledAtCount(ctx *gin.Context) {
	request := &gatewaycontract.Request[routinetaskrecordsdto.VisualizeMyRoutineTaskRecordScheduledAtCountRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineTaskRecordService.VisualizeMyRoutineTaskRecordScheduledAtCount(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), gatewaycontract.Response[struct{}]{
			Version: gatewaycontract.Version,
			Metadata: gatewaycontract.ResponseMetadata{
				RequestId:   request.Metadata.RequestId,
				RespondedAt: time.Now(),
			},
			Data:      struct{}{},
			Exception: publicException,
		})
		return
	}

	ctx.JSON(http.StatusOK, gatewaycontract.Response[routinetaskrecordsdto.VisualizeMyRoutineTaskRecordScheduledAtCountResponseDto]{
		Version:  gatewaycontract.Version,
		Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()},
		Data:     *responseDto,
	})
}

func (t *RoutineTaskRecordEndpoint) VisualizeMyRoutineTaskRecordActualStartedAtCount(ctx *gin.Context) {
	request := &gatewaycontract.Request[routinetaskrecordsdto.VisualizeMyRoutineTaskRecordActualStartedAtCountRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineTaskRecordService.VisualizeMyRoutineTaskRecordActualStartedAtCount(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), gatewaycontract.Response[struct{}]{
			Version: gatewaycontract.Version,
			Metadata: gatewaycontract.ResponseMetadata{
				RequestId:   request.Metadata.RequestId,
				RespondedAt: time.Now(),
			},
			Data:      struct{}{},
			Exception: publicException,
		})
		return
	}

	ctx.JSON(http.StatusOK, gatewaycontract.Response[routinetaskrecordsdto.VisualizeMyRoutineTaskRecordActualStartedAtCountResponseDto]{
		Version:  gatewaycontract.Version,
		Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()},
		Data:     *responseDto,
	})
}

func (t *RoutineTaskRecordEndpoint) VisualizeMyRoutineTaskRecordActualEndedAtCount(ctx *gin.Context) {
	request := &gatewaycontract.Request[routinetaskrecordsdto.VisualizeMyRoutineTaskRecordActualEndedAtCountRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	responseDto, exception := t.routineTaskRecordService.VisualizeMyRoutineTaskRecordActualEndedAtCount(ctx.Request.Context(), &request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), gatewaycontract.Response[struct{}]{
			Version: gatewaycontract.Version,
			Metadata: gatewaycontract.ResponseMetadata{
				RequestId:   request.Metadata.RequestId,
				RespondedAt: time.Now(),
			},
			Data:      struct{}{},
			Exception: publicException,
		})
		return
	}

	ctx.JSON(http.StatusOK, gatewaycontract.Response[routinetaskrecordsdto.VisualizeMyRoutineTaskRecordActualEndedAtCountResponseDto]{
		Version:  gatewaycontract.Version,
		Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()},
		Data:     *responseDto,
	})
}
