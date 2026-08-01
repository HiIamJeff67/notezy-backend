package yjsmaintenance

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	blocksdto "github.com/HiIamJeff67/notezy-backend/contracts/api/v1/blocks"
	inputs "github.com/HiIamJeff67/notezy-backend/internal/services/durablejob/data/inputs"
	options "github.com/HiIamJeff67/notezy-backend/internal/services/durablejob/data/options"
	repositories "github.com/HiIamJeff67/notezy-backend/internal/services/durablejob/data/repositories"
	coretransport "github.com/HiIamJeff67/notezy-backend/internal/services/durablejob/transports/core"
	constants "github.com/HiIamJeff67/notezy-backend/internal/shared/constants"
	sharedtypes "github.com/HiIamJeff67/notezy-backend/internal/shared/types"
)

type Handler struct {
	db                     *gorm.DB
	blockPackYjsRepository repositories.BlockPackYjsRepositoryInterface
	projectionClient       coretransport.Client
}

func NewHandler(db *gorm.DB, projectionClient coretransport.Client) Handler {
	return Handler{
		db:                     db,
		blockPackYjsRepository: repositories.NewBlockPackYjsRepository(),
		projectionClient:       projectionClient,
	}
}

func (h Handler) HandleCompactions(
	ctx context.Context,
	batchInputs []sharedtypes.YjsCompactionBatchInput,
	results []sharedtypes.YjsCompactionBatchResult,
) ([]uuid.UUID, error) {
	if len(batchInputs) == 0 || len(results) == 0 {
		return []uuid.UUID{}, nil
	}
	if len(batchInputs) != len(results) {
		return nil, fmt.Errorf("incomplete yjs compaction batch result")
	}

	inputByBlockPackId := make(map[uuid.UUID]sharedtypes.YjsCompactionBatchInput, len(batchInputs))
	for _, input := range batchInputs {
		inputByBlockPackId[input.BlockPackId] = input
	}

	bulkInputs := make([]inputs.BulkApplyCompactedBlockPackYjsDocumentInput, 0, len(results))
	for _, result := range results {
		input, exists := inputByBlockPackId[result.BlockPackId]
		if !exists ||
			result.Result.BaseCompactedUntilSequence != input.Input.BaseCompactedUntilSequence ||
			result.Result.CutoffSequence != input.Input.CutoffSequence {
			return nil, fmt.Errorf("invalid yjs compaction batch result")
		}

		bulkInputs = append(bulkInputs, inputs.BulkApplyCompactedBlockPackYjsDocumentInput{
			BlockPackId: result.BlockPackId,
			ApplyCompactedBlockPackYjsDocumentInput: inputs.ApplyCompactedBlockPackYjsDocumentInput{
				BaseCompactedUntilSequence: result.Result.BaseCompactedUntilSequence,
				CutoffSequence:             result.Result.CutoffSequence,
				Snapshot:                   result.Result.Snapshot,
				StateVector:                result.Result.StateVector,
			},
		})
	}

	tx := h.db.WithContext(ctx).Begin()

	appliedBlockPackIds, err := h.blockPackYjsRepository.BulkApplyCompactedYjsDocuments(
		bulkInputs,
		options.WithTransactionDB(tx),
	)
	if err != nil {
		tx.Rollback()

		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return appliedBlockPackIds, nil
}

func (h Handler) HandleProjections(
	ctx context.Context,
	inputs []sharedtypes.YjsProjectionBatchInput,
	results []sharedtypes.YjsProjectionBatchResult,
) (blocksdto.ApplyBlockProjectionDocumentResponseDto, error) {
	if len(inputs) == 0 && len(results) == 0 {
		return blocksdto.ApplyBlockProjectionDocumentResponseDto{}, nil
	}
	if len(inputs) != len(results) {
		return nil, fmt.Errorf("incomplete yjs projection batch result")
	}

	inputByBlockPackId := make(map[uuid.UUID]sharedtypes.YjsProjectionBatchInput, len(inputs))
	for _, input := range inputs {
		if _, exists := inputByBlockPackId[input.BlockPackId]; exists {
			return nil, fmt.Errorf("duplicate yjs projection input")
		}
		inputByBlockPackId[input.BlockPackId] = input
	}

	projectionRequestDtos := make([]blocksdto.ApplyBlockProjectionDocumentRequestDto, 0, len(results))
	resultBlockPackIdSet := make(map[uuid.UUID]bool, len(results))
	for _, result := range results {
		input, exists := inputByBlockPackId[result.BlockPackId]
		if !exists || resultBlockPackIdSet[result.BlockPackId] {
			return nil, fmt.Errorf("invalid yjs projection batch result")
		}

		var projectionRequestDto blocksdto.ApplyBlockProjectionRequestDto
		if err := json.Unmarshal(result.Payload, &projectionRequestDto); err != nil {
			return nil, fmt.Errorf("failed to decode yjs projection result: %w", err)
		}
		if projectionRequestDto.ProjectedSequence != input.State.LastUpdateSequence {
			return nil, fmt.Errorf("yjs projection result sequence does not match the claimed document")
		}

		resultBlockPackIdSet[result.BlockPackId] = true
		projectionRequestDtos = append(projectionRequestDtos, blocksdto.ApplyBlockProjectionDocumentRequestDto{
			BlockPackId: result.BlockPackId,
			Projection:  projectionRequestDto,
		})
	}
	if len(resultBlockPackIdSet) != len(inputByBlockPackId) {
		return nil, fmt.Errorf("incomplete yjs projection batch result")
	}

	appliedBlockPackIds, err := h.projectionClient.ApplyBlockProjections(ctx, projectionRequestDtos)
	if err != nil {
		return nil, err
	}

	return blocksdto.ApplyBlockProjectionDocumentResponseDto(appliedBlockPackIds), nil
}

func (h Handler) Cleanup(ctx context.Context) error {
	db := h.db.WithContext(ctx)

	_, err := h.blockPackYjsRepository.DeleteCompactedUpdates(
		inputs.DeleteCompactedBlockPackYjsUpdatesInput{
			Before: time.Now().Add(-constants.YjsCompactedUpdateRetention),
			Limit:  constants.YjsCleanupMaxUpdatesPerRun,
		},
		options.WithDB(db),
	)

	return err
}
