package core

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	blocksdto "github.com/HiIamJeff67/notezy-backend/contracts/gateway/v1/api/blocks"
	websocketdto "github.com/HiIamJeff67/notezy-backend/contracts/gateway/v1/websocket"
	coreadapters "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/core/adapters"
	sharedtypes "github.com/HiIamJeff67/notezy-backend/internal/shared/types"
)

const (
	validateBlockPackChannelPermissionPath = "/websocket/block-pack-channel-permission/validate"
	loadCompactableYjsDocumentPath         = "/websocket/yjs-document-compaction/load"
	applyCompactedYjsDocumentPath          = "/websocket/yjs-document-compaction/apply"
	loadYjsDocumentPath                    = "/websocket/yjs-document/load"
	appendYjsUpdatePath                    = "/websocket/yjs-update/append"
	applyBlockProjectionPath               = "/websocket/block-projection/apply"
)

type CoreClient struct {
	client *coreadapters.CoreClient
}

func NewCoreClient(client *coreadapters.CoreClient) *CoreClient {
	return &CoreClient{
		client: client,
	}
}

func (c *CoreClient) ValidateBlockPackChannelPermission(
	ctx context.Context,
	userPublicId uuid.UUID,
	blockPackId uuid.UUID,
	permission sharedtypes.ChannelPermission,
) (sharedtypes.ErrorCode, error) {
	response, exception := coreadapters.CallAsComponent[
		websocketdto.ValidateBlockPackChannelPermissionRequestDto,
		websocketdto.ValidateBlockPackChannelPermissionResponseDto,
	](
		ctx,
		c.client,
		"websocket",
		&websocketdto.ValidateBlockPackChannelPermissionRequestDto{
			UserPublicId: userPublicId,
			BlockPackId:  blockPackId,
			Permission:   string(permission),
		},
		websocketdto.ValidateBlockPackChannelPermissionOperation,
		validateBlockPackChannelPermissionPath,
	)
	if exception != nil {
		return "", exception
	}
	if !response.Data.Permitted {
		return sharedtypes.ErrorCode(response.Data.ErrorCode), errors.New("block pack channel permission was rejected")
	}

	return "", nil
}

func (c *CoreClient) GetCompactableYjsDocumentWithUpdates(
	ctx context.Context,
	blockPackId uuid.UUID,
) (*sharedtypes.YjsCompactionInput, error) {
	response, exception := coreadapters.CallAsComponent[
		websocketdto.LoadCompactableYjsDocumentRequestDto,
		websocketdto.LoadCompactableYjsDocumentResponseDto,
	](
		ctx,
		c.client,
		"websocket",
		&websocketdto.LoadCompactableYjsDocumentRequestDto{
			BlockPackId: blockPackId,
		},
		websocketdto.LoadCompactableYjsDocumentOperation,
		loadCompactableYjsDocumentPath,
	)
	if exception != nil {
		return nil, exception
	}
	if !response.Data.Found {
		return nil, nil
	}

	input := &sharedtypes.YjsCompactionInput{}
	if err := input.UnmarshalBytes(response.Data.Payload); err != nil {
		return nil, err
	}

	return input, nil
}

func (c *CoreClient) ApplyCompactedYjsDocument(
	ctx context.Context,
	blockPackId uuid.UUID,
	result sharedtypes.YjsCompactionResult,
) (bool, error) {
	payload, err := result.MarshalBytes()
	if err != nil {
		return false, err
	}

	response, exception := coreadapters.CallAsComponent[
		websocketdto.ApplyCompactedYjsDocumentRequestDto,
		websocketdto.ApplyCompactedYjsDocumentResponseDto,
	](
		ctx,
		c.client,
		"websocket",
		&websocketdto.ApplyCompactedYjsDocumentRequestDto{
			BlockPackId: blockPackId,
			Payload:     payload,
		},
		websocketdto.ApplyCompactedYjsDocumentOperation,
		applyCompactedYjsDocumentPath,
	)
	if exception != nil {
		return false, exception
	}

	return response.Data.Applied, nil
}

func (c *CoreClient) LoadDocument(
	ctx context.Context,
	blockPackId uuid.UUID,
) (*sharedtypes.YjsDocumentState, error) {
	response, exception := coreadapters.CallAsComponent[
		websocketdto.LoadYjsDocumentRequestDto,
		websocketdto.LoadYjsDocumentResponseDto,
	](
		ctx,
		c.client,
		"websocket",
		&websocketdto.LoadYjsDocumentRequestDto{
			BlockPackId: blockPackId,
		},
		websocketdto.LoadYjsDocumentOperation,
		loadYjsDocumentPath,
	)
	if exception != nil {
		if exception.HTTPStatusCode() == 404 {
			return nil, gorm.ErrRecordNotFound
		}

		return nil, exception
	}

	state := &sharedtypes.YjsDocumentState{}
	if err := state.UnmarshalBytes(response.Data.Payload); err != nil {
		return nil, err
	}

	return state, nil
}

func (c *CoreClient) AppendUpdate(
	ctx context.Context,
	blockPackId uuid.UUID,
	persistenceBatchId uuid.UUID,
	originConnectionId *uuid.UUID,
	payload []byte,
) (int64, error) {
	response, exception := coreadapters.CallAsComponent[
		websocketdto.AppendYjsUpdateRequestDto,
		websocketdto.AppendYjsUpdateResponseDto,
	](
		ctx,
		c.client,
		"websocket",
		&websocketdto.AppendYjsUpdateRequestDto{
			BlockPackId:        blockPackId,
			PersistenceBatchId: persistenceBatchId,
			OriginConnectionId: originConnectionId,
			Payload:            payload,
		},
		websocketdto.AppendYjsUpdateOperation,
		appendYjsUpdatePath,
	)
	if exception != nil {
		if exception.HTTPStatusCode() == 404 {
			return 0, gorm.ErrRecordNotFound
		}

		return 0, exception
	}

	return response.Data.UpdateSequence, nil
}

func (c *CoreClient) ApplyBlockProjection(
	ctx context.Context,
	blockPackId uuid.UUID,
	requestDto blocksdto.ApplyBlockProjectionRequestDto,
) (*blocksdto.ApplyBlockProjectionResponseDto, error) {
	response, exception := coreadapters.CallAsComponent[
		websocketdto.ApplyBlockProjectionRequestDto,
		websocketdto.ApplyBlockProjectionResponseDto,
	](
		ctx,
		c.client,
		"websocket",
		&websocketdto.ApplyBlockProjectionRequestDto{
			BlockPackId: blockPackId,
			Projection:  requestDto,
		},
		websocketdto.ApplyBlockProjectionOperation,
		applyBlockProjectionPath,
	)
	if exception != nil {
		return nil, exception
	}

	return &blocksdto.ApplyBlockProjectionResponseDto{
		Applied:                response.Data.Applied,
		ProjectedUntilSequence: response.Data.ProjectedUntilSequence,
	}, nil
}
