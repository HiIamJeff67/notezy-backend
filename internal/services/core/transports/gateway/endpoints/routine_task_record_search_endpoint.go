package endpoints

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	routinetaskrecordsdto "github.com/HiIamJeff67/notezy-backend/contracts/api/v1/routine-task-records"
	core "github.com/HiIamJeff67/notezy-backend/contracts/core/v1"
	contexts "github.com/HiIamJeff67/notezy-backend/internal/services/core/contexts"
)

func (t *RoutineTaskRecordEndpoint) SearchRoutineTaskRecords(ctx *gin.Context) {
	request := &core.Request[routinetaskrecordsdto.SearchRoutineTaskRecordsRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	userId, exception := contexts.GetActorUserId(ctx.Request.Context())
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), core.Response[struct{}]{
			Version: core.Version,
			Metadata: core.ResponseMetadata{
				RequestId:   request.Metadata.RequestId,
				RespondedAt: time.Now(),
			},
			Data:      struct{}{},
			Exception: publicException,
		})
		return
	}

	responseDto, exception := t.routineTaskRecordService.SearchPrivateRoutineTaskRecords(ctx.Request.Context(), userId, request.Dto)
	if exception != nil {
		publicException := exception.ToPublic()
		ctx.JSON(publicException.HTTPStatusCode(), core.Response[struct{}]{
			Version: core.Version,
			Metadata: core.ResponseMetadata{
				RequestId:   request.Metadata.RequestId,
				RespondedAt: time.Now(),
			},
			Data:      struct{}{},
			Exception: publicException,
		})
		return
	}

	ctx.JSON(http.StatusOK, core.Response[routinetaskrecordsdto.SearchRoutineTaskRecordsResponseDto]{
		Version: core.Version,
		Metadata: core.ResponseMetadata{
			RequestId:   request.Metadata.RequestId,
			RespondedAt: time.Now(),
		},
		Data: *responseDto,
	})
}
