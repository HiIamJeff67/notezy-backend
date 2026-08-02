package endpoints

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	core "github.com/HiIamJeff67/notezy-backend/contracts/core/v1"
	websocketdto "github.com/HiIamJeff67/notezy-backend/contracts/gateway/v1/websocket"
	exceptions "github.com/HiIamJeff67/notezy-backend/internal/exceptions"
	services "github.com/HiIamJeff67/notezy-backend/internal/services/core/services"
	validation "github.com/HiIamJeff67/notezy-backend/internal/services/core/validation"
	sharedtypes "github.com/HiIamJeff67/notezy-backend/internal/shared/types"
)

type YjsEndpoint struct {
	yjsPersistenceService services.YjsPersistenceServiceInterface
}

func NewYjsEndpoint(yjsPersistenceService services.YjsPersistenceServiceInterface) YjsEndpoint {
	return YjsEndpoint{
		yjsPersistenceService: yjsPersistenceService,
	}
}

func (e YjsEndpoint) LoadCompactableYjsDocument(ctx *gin.Context) {
	request := &core.Request[websocketdto.LoadCompactableYjsDocumentRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		exception := exceptions.New(
			"InvalidRequest",
			"WebSocket",
			websocketdto.LoadCompactableYjsDocumentOperation,
			"The compactable Yjs document request is invalid",
			http.StatusBadRequest,
		).WithOrigin(err)
		ctx.JSON(exception.HTTPStatusCode(), core.Response[websocketdto.LoadCompactableYjsDocumentResponseDto]{
			Version: core.Version,
			Metadata: core.ResponseMetadata{
				RequestId:   request.Metadata.RequestId,
				RespondedAt: time.Now(),
			},
			Exception: exception,
		})
		return
	}
	if err := validation.Validator.Struct(request.Dto); err != nil {
		exception := exceptions.New(
			"InvalidRequest",
			"WebSocket",
			websocketdto.LoadCompactableYjsDocumentOperation,
			"The compactable Yjs document request is invalid",
			http.StatusBadRequest,
		).WithOrigin(err)
		ctx.JSON(exception.HTTPStatusCode(), core.Response[websocketdto.LoadCompactableYjsDocumentResponseDto]{
			Version: core.Version,
			Metadata: core.ResponseMetadata{
				RequestId:   request.Metadata.RequestId,
				RespondedAt: time.Now(),
			},
			Exception: exception,
		})
		return
	}

	input, err := e.yjsPersistenceService.GetCompactableYjsDocumentWithUpdates(
		ctx.Request.Context(),
		request.Dto.BlockPackId,
	)
	if err != nil {
		exception := exceptions.New(
			"FailedToLoad",
			"BlockPackYjsDocument",
			websocketdto.LoadCompactableYjsDocumentOperation,
			"Failed to load the compactable Yjs document",
			http.StatusInternalServerError,
			true,
		).WithOrigin(err)
		ctx.JSON(exception.HTTPStatusCode(), core.Response[websocketdto.LoadCompactableYjsDocumentResponseDto]{
			Version: core.Version,
			Metadata: core.ResponseMetadata{
				RequestId:   request.Metadata.RequestId,
				RespondedAt: time.Now(),
			},
			Exception: exception,
		})
		return
	}
	if input == nil {
		ctx.JSON(http.StatusOK, core.Response[websocketdto.LoadCompactableYjsDocumentResponseDto]{
			Version: core.Version,
			Metadata: core.ResponseMetadata{
				RequestId:   request.Metadata.RequestId,
				RespondedAt: time.Now(),
			},
			Data: websocketdto.LoadCompactableYjsDocumentResponseDto{
				Found: false,
			},
		})
		return
	}

	payload, err := input.MarshalBytes()
	if err != nil {
		exception := exceptions.New(
			"FailedToEncodeResponse",
			"BlockPackYjsDocument",
			websocketdto.LoadCompactableYjsDocumentOperation,
			"Failed to encode the compactable Yjs document",
			http.StatusInternalServerError,
			true,
		).WithOrigin(err)
		ctx.JSON(exception.HTTPStatusCode(), core.Response[websocketdto.LoadCompactableYjsDocumentResponseDto]{
			Version: core.Version,
			Metadata: core.ResponseMetadata{
				RequestId:   request.Metadata.RequestId,
				RespondedAt: time.Now(),
			},
			Exception: exception,
		})
		return
	}

	ctx.JSON(http.StatusOK, core.Response[websocketdto.LoadCompactableYjsDocumentResponseDto]{
		Version: core.Version,
		Metadata: core.ResponseMetadata{
			RequestId:   request.Metadata.RequestId,
			RespondedAt: time.Now(),
		},
		Data: websocketdto.LoadCompactableYjsDocumentResponseDto{
			Found:   true,
			Payload: payload,
		},
	})
}

func (e YjsEndpoint) ApplyCompactedYjsDocument(ctx *gin.Context) {
	request := &core.Request[websocketdto.ApplyCompactedYjsDocumentRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		exception := exceptions.New(
			"InvalidRequest",
			"WebSocket",
			websocketdto.ApplyCompactedYjsDocumentOperation,
			"The compacted Yjs document request is invalid",
			http.StatusBadRequest,
		).WithOrigin(err)
		ctx.JSON(exception.HTTPStatusCode(), core.Response[websocketdto.ApplyCompactedYjsDocumentResponseDto]{
			Version: core.Version,
			Metadata: core.ResponseMetadata{
				RequestId:   request.Metadata.RequestId,
				RespondedAt: time.Now(),
			},
			Exception: exception,
		})
		return
	}
	if err := validation.Validator.Struct(request.Dto); err != nil {
		exception := exceptions.New(
			"InvalidRequest",
			"WebSocket",
			websocketdto.ApplyCompactedYjsDocumentOperation,
			"The compacted Yjs document request is invalid",
			http.StatusBadRequest,
		).WithOrigin(err)
		ctx.JSON(exception.HTTPStatusCode(), core.Response[websocketdto.ApplyCompactedYjsDocumentResponseDto]{
			Version: core.Version,
			Metadata: core.ResponseMetadata{
				RequestId:   request.Metadata.RequestId,
				RespondedAt: time.Now(),
			},
			Exception: exception,
		})
		return
	}

	var result sharedtypes.YjsCompactionResult
	if err := result.UnmarshalBytes(request.Dto.Payload); err != nil {
		exception := exceptions.New(
			"InvalidRequest",
			"WebSocket",
			websocketdto.ApplyCompactedYjsDocumentOperation,
			"The compacted Yjs document payload is invalid",
			http.StatusBadRequest,
		).WithOrigin(err)
		ctx.JSON(exception.HTTPStatusCode(), core.Response[websocketdto.ApplyCompactedYjsDocumentResponseDto]{
			Version: core.Version,
			Metadata: core.ResponseMetadata{
				RequestId:   request.Metadata.RequestId,
				RespondedAt: time.Now(),
			},
			Exception: exception,
		})
		return
	}

	applied, err := e.yjsPersistenceService.ApplyCompactedYjsDocument(
		ctx.Request.Context(),
		request.Dto.BlockPackId,
		result,
	)
	if err != nil {
		exception := exceptions.New(
			"FailedToApply",
			"BlockPackYjsDocument",
			websocketdto.ApplyCompactedYjsDocumentOperation,
			"Failed to apply the compacted Yjs document",
			http.StatusInternalServerError,
			true,
		).WithOrigin(err)
		ctx.JSON(exception.HTTPStatusCode(), core.Response[websocketdto.ApplyCompactedYjsDocumentResponseDto]{
			Version: core.Version,
			Metadata: core.ResponseMetadata{
				RequestId:   request.Metadata.RequestId,
				RespondedAt: time.Now(),
			},
			Exception: exception,
		})
		return
	}

	ctx.JSON(http.StatusOK, core.Response[websocketdto.ApplyCompactedYjsDocumentResponseDto]{
		Version: core.Version,
		Metadata: core.ResponseMetadata{
			RequestId:   request.Metadata.RequestId,
			RespondedAt: time.Now(),
		},
		Data: websocketdto.ApplyCompactedYjsDocumentResponseDto{
			Applied: applied,
		},
	})
}

func (e YjsEndpoint) LoadYjsDocument(ctx *gin.Context) {
	request := &core.Request[websocketdto.LoadYjsDocumentRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		exception := exceptions.New(
			"InvalidRequest",
			"WebSocket",
			websocketdto.LoadYjsDocumentOperation,
			"The Yjs document request is invalid",
			http.StatusBadRequest,
		).WithOrigin(err)
		ctx.JSON(exception.HTTPStatusCode(), core.Response[websocketdto.LoadYjsDocumentResponseDto]{
			Version: core.Version,
			Metadata: core.ResponseMetadata{
				RequestId:   request.Metadata.RequestId,
				RespondedAt: time.Now(),
			},
			Exception: exception,
		})
		return
	}
	if err := validation.Validator.Struct(request.Dto); err != nil {
		exception := exceptions.New(
			"InvalidRequest",
			"WebSocket",
			websocketdto.LoadYjsDocumentOperation,
			"The Yjs document request is invalid",
			http.StatusBadRequest,
		).WithOrigin(err)
		ctx.JSON(exception.HTTPStatusCode(), core.Response[websocketdto.LoadYjsDocumentResponseDto]{
			Version: core.Version,
			Metadata: core.ResponseMetadata{
				RequestId:   request.Metadata.RequestId,
				RespondedAt: time.Now(),
			},
			Exception: exception,
		})
		return
	}

	state, err := e.yjsPersistenceService.LoadDocument(
		ctx.Request.Context(),
		request.Dto.BlockPackId,
	)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, gorm.ErrRecordNotFound) {
			statusCode = http.StatusNotFound
		}
		exception := exceptions.New(
			"FailedToLoad",
			"BlockPackYjsDocument",
			websocketdto.LoadYjsDocumentOperation,
			"Failed to load the Yjs document",
			statusCode,
			statusCode >= http.StatusInternalServerError,
		).WithOrigin(err)
		ctx.JSON(exception.HTTPStatusCode(), core.Response[websocketdto.LoadYjsDocumentResponseDto]{
			Version: core.Version,
			Metadata: core.ResponseMetadata{
				RequestId:   request.Metadata.RequestId,
				RespondedAt: time.Now(),
			},
			Exception: exception,
		})
		return
	}

	payload, err := state.MarshalBytes()
	if err != nil {
		exception := exceptions.New(
			"FailedToEncodeResponse",
			"BlockPackYjsDocument",
			websocketdto.LoadYjsDocumentOperation,
			"Failed to encode the Yjs document",
			http.StatusInternalServerError,
			true,
		).WithOrigin(err)
		ctx.JSON(exception.HTTPStatusCode(), core.Response[websocketdto.LoadYjsDocumentResponseDto]{
			Version: core.Version,
			Metadata: core.ResponseMetadata{
				RequestId:   request.Metadata.RequestId,
				RespondedAt: time.Now(),
			},
			Exception: exception,
		})
		return
	}

	ctx.JSON(http.StatusOK, core.Response[websocketdto.LoadYjsDocumentResponseDto]{
		Version: core.Version,
		Metadata: core.ResponseMetadata{
			RequestId:   request.Metadata.RequestId,
			RespondedAt: time.Now(),
		},
		Data: websocketdto.LoadYjsDocumentResponseDto{
			Payload: payload,
		},
	})
}

func (e YjsEndpoint) AppendYjsUpdate(ctx *gin.Context) {
	request := &core.Request[websocketdto.AppendYjsUpdateRequestDto]{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		exception := exceptions.New(
			"InvalidRequest",
			"WebSocket",
			websocketdto.AppendYjsUpdateOperation,
			"The Yjs update request is invalid",
			http.StatusBadRequest,
		).WithOrigin(err)
		ctx.JSON(exception.HTTPStatusCode(), core.Response[websocketdto.AppendYjsUpdateResponseDto]{
			Version: core.Version,
			Metadata: core.ResponseMetadata{
				RequestId:   request.Metadata.RequestId,
				RespondedAt: time.Now(),
			},
			Exception: exception,
		})
		return
	}
	if err := validation.Validator.Struct(request.Dto); err != nil {
		exception := exceptions.New(
			"InvalidRequest",
			"WebSocket",
			websocketdto.AppendYjsUpdateOperation,
			"The Yjs update request is invalid",
			http.StatusBadRequest,
		).WithOrigin(err)
		ctx.JSON(exception.HTTPStatusCode(), core.Response[websocketdto.AppendYjsUpdateResponseDto]{
			Version: core.Version,
			Metadata: core.ResponseMetadata{
				RequestId:   request.Metadata.RequestId,
				RespondedAt: time.Now(),
			},
			Exception: exception,
		})
		return
	}

	updateSequence, err := e.yjsPersistenceService.AppendUpdate(
		ctx.Request.Context(),
		request.Dto.BlockPackId,
		request.Dto.PersistenceBatchId,
		request.Dto.OriginConnectionId,
		request.Dto.Payload,
	)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, gorm.ErrRecordNotFound) {
			statusCode = http.StatusNotFound
		}
		exception := exceptions.New(
			"FailedToAppend",
			"BlockPackYjsDocument",
			websocketdto.AppendYjsUpdateOperation,
			"Failed to append the Yjs update",
			statusCode,
			statusCode >= http.StatusInternalServerError,
		).WithOrigin(err)
		ctx.JSON(exception.HTTPStatusCode(), core.Response[websocketdto.AppendYjsUpdateResponseDto]{
			Version: core.Version,
			Metadata: core.ResponseMetadata{
				RequestId:   request.Metadata.RequestId,
				RespondedAt: time.Now(),
			},
			Exception: exception,
		})
		return
	}

	ctx.JSON(http.StatusOK, core.Response[websocketdto.AppendYjsUpdateResponseDto]{
		Version: core.Version,
		Metadata: core.ResponseMetadata{
			RequestId:   request.Metadata.RequestId,
			RespondedAt: time.Now(),
		},
		Data: websocketdto.AppendYjsUpdateResponseDto{
			UpdateSequence: updateSequence,
		},
	})
}
