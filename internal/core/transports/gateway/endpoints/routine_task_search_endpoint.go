package endpoints

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	routinetasksdto "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api/routine-tasks"
	gatewaycontract "github.com/HiIamJeff67/notezy-backend/contracts/gateway/v1"

	contexts "github.com/HiIamJeff67/notezy-backend/internal/core/contexts"
)

func (t *RoutineTaskEndpoint) SearchRoutineTasks(ctx *gin.Context) {
	request := &gatewaycontract.Request[routinetasksdto.SearchRoutineTasksRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	userId, exception := contexts.GetActorUserId(ctx.Request.Context())
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

	responseDto, exception := t.routineTaskService.SearchPrivateRoutineTasks(ctx.Request.Context(), userId, request.Dto)
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

	ctx.JSON(http.StatusOK, gatewaycontract.Response[routinetasksdto.SearchRoutineTasksResponseDto]{
		Version: gatewaycontract.Version,
		Metadata: gatewaycontract.ResponseMetadata{
			RequestId:   request.Metadata.RequestId,
			RespondedAt: time.Now(),
		},
		Data: *responseDto,
	})
}
