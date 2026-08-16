package endpoints

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	apicontract "github.com/HiIamJeff67/notegic-backend/contracts/core/v1/api/routine-task-records"
	gatewaycontract "github.com/HiIamJeff67/notegic-backend/contracts/gateway/v1"

	routineservices "github.com/HiIamJeff67/notegic-backend/internal/core/services/routines"
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
	routineTaskRecordService routineservices.RoutineTaskRecordServiceInterface
}

func NewRoutineTaskRecordEndpoint(
	routineTaskRecordService routineservices.RoutineTaskRecordServiceInterface,
) RoutineTaskRecordEndpointInterface {
	return &RoutineTaskRecordEndpoint{
		routineTaskRecordService: routineTaskRecordService,
	}
}

func (t *RoutineTaskRecordEndpoint) GetAllMyRoutineTaskRecordsByRoutineTaskId(ctx *gin.Context) {
	request := &gatewaycontract.Request[apicontract.GetAllMyRoutineTaskRecordsByRoutineTaskIdRequestDto]{}
	if err := ctx.ShouldBindBodyWithJSON(request); err != nil {
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

	ctx.JSON(http.StatusOK, gatewaycontract.Response[apicontract.GetAllMyRoutineTaskRecordsByRoutineTaskIdResponseDto]{
		Version:  gatewaycontract.Version,
		Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()},
		Data:     *responseDto,
	})
}

/* ============================== Visualization Methods ============================== */

func (t *RoutineTaskRecordEndpoint) VisualizeMyRoutineTaskRecordStatusCount(ctx *gin.Context) {
	request := &gatewaycontract.Request[apicontract.VisualizeMyRoutineTaskRecordStatusCountRequestDto]{}
	if err := ctx.ShouldBindBodyWithJSON(request); err != nil {
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

	ctx.JSON(http.StatusOK, gatewaycontract.Response[apicontract.VisualizeMyRoutineTaskRecordStatusCountResponseDto]{
		Version:  gatewaycontract.Version,
		Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()},
		Data:     *responseDto,
	})
}

func (t *RoutineTaskRecordEndpoint) VisualizeMyRoutineTaskRecordPurposeCount(ctx *gin.Context) {
	request := &gatewaycontract.Request[apicontract.VisualizeMyRoutineTaskRecordPurposeCountRequestDto]{}
	if err := ctx.ShouldBindBodyWithJSON(request); err != nil {
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

	ctx.JSON(http.StatusOK, gatewaycontract.Response[apicontract.VisualizeMyRoutineTaskRecordPurposeCountResponseDto]{
		Version:  gatewaycontract.Version,
		Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()},
		Data:     *responseDto,
	})
}

func (t *RoutineTaskRecordEndpoint) VisualizeMyRoutineTaskRecordScheduledAtCount(ctx *gin.Context) {
	request := &gatewaycontract.Request[apicontract.VisualizeMyRoutineTaskRecordScheduledAtCountRequestDto]{}
	if err := ctx.ShouldBindBodyWithJSON(request); err != nil {
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

	ctx.JSON(http.StatusOK, gatewaycontract.Response[apicontract.VisualizeMyRoutineTaskRecordScheduledAtCountResponseDto]{
		Version:  gatewaycontract.Version,
		Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()},
		Data:     *responseDto,
	})
}

func (t *RoutineTaskRecordEndpoint) VisualizeMyRoutineTaskRecordActualStartedAtCount(ctx *gin.Context) {
	request := &gatewaycontract.Request[apicontract.VisualizeMyRoutineTaskRecordActualStartedAtCountRequestDto]{}
	if err := ctx.ShouldBindBodyWithJSON(request); err != nil {
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

	ctx.JSON(http.StatusOK, gatewaycontract.Response[apicontract.VisualizeMyRoutineTaskRecordActualStartedAtCountResponseDto]{
		Version:  gatewaycontract.Version,
		Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()},
		Data:     *responseDto,
	})
}

func (t *RoutineTaskRecordEndpoint) VisualizeMyRoutineTaskRecordActualEndedAtCount(ctx *gin.Context) {
	request := &gatewaycontract.Request[apicontract.VisualizeMyRoutineTaskRecordActualEndedAtCountRequestDto]{}
	if err := ctx.ShouldBindBodyWithJSON(request); err != nil {
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

	ctx.JSON(http.StatusOK, gatewaycontract.Response[apicontract.VisualizeMyRoutineTaskRecordActualEndedAtCountResponseDto]{
		Version:  gatewaycontract.Version,
		Metadata: gatewaycontract.ResponseMetadata{RequestId: request.Metadata.RequestId, RespondedAt: time.Now()},
		Data:     *responseDto,
	})
}
