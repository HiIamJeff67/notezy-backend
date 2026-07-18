package yjsmaintenance

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	dtos "github.com/HiIamJeff67/notezy-backend/app/dtos"
	inputs "github.com/HiIamJeff67/notezy-backend/app/models/inputs"
	repositories "github.com/HiIamJeff67/notezy-backend/app/models/repositories"
	options "github.com/HiIamJeff67/notezy-backend/app/options"
	realtimetypes "github.com/HiIamJeff67/notezy-backend/app/realtime/types"
	constants "github.com/HiIamJeff67/notezy-backend/shared/constants"
)

type BlockProjectionService interface {
	ApplyMany(ctx context.Context, inputs []dtos.ApplyBlockProjectionDocumentInput) (dtos.ApplyBlockProjectionDocumentResult, error)
}

type Handler struct {
	db                     *gorm.DB
	blockPackYjsRepository repositories.BlockPackYjsRepositoryInterface
	blockService           BlockProjectionService
}

func NewHandler(db *gorm.DB, blockService BlockProjectionService) Handler {
	return Handler{
		db:                     db,
		blockPackYjsRepository: repositories.NewBlockPackYjsRepository(),
		blockService:           blockService,
	}
}

func (h Handler) HandleCompactions(
	ctx context.Context,
	batchInputs []realtimetypes.YjsCompactionBatchInput,
	results []realtimetypes.YjsCompactionBatchResult,
) ([]uuid.UUID, error) {
	if len(batchInputs) == 0 || len(results) == 0 {
		return []uuid.UUID{}, nil
	}
	if len(batchInputs) != len(results) {
		return nil, fmt.Errorf("incomplete yjs compaction batch result")
	}

	inputByBlockPackId := make(map[uuid.UUID]realtimetypes.YjsCompactionBatchInput, len(batchInputs))
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
	inputs []realtimetypes.YjsProjectionBatchInput,
	results []realtimetypes.YjsProjectionBatchResult,
) (dtos.ApplyBlockProjectionDocumentResult, error) {
	if len(inputs) == 0 && len(results) == 0 {
		return dtos.ApplyBlockProjectionDocumentResult{}, nil
	}
	if len(inputs) != len(results) {
		return nil, fmt.Errorf("incomplete yjs projection batch result")
	}

	inputByBlockPackId := make(map[uuid.UUID]realtimetypes.YjsProjectionBatchInput, len(inputs))
	for _, input := range inputs {
		if _, exists := inputByBlockPackId[input.BlockPackId]; exists {
			return nil, fmt.Errorf("duplicate yjs projection input")
		}
		inputByBlockPackId[input.BlockPackId] = input
	}

	projectionInputs := make([]dtos.ApplyBlockProjectionDocumentInput, 0, len(results))
	resultBlockPackIdSet := make(map[uuid.UUID]bool, len(results))
	for _, result := range results {
		input, exists := inputByBlockPackId[result.BlockPackId]
		if !exists || resultBlockPackIdSet[result.BlockPackId] {
			return nil, fmt.Errorf("invalid yjs projection batch result")
		}

		var projection dtos.ApplyBlockProjectionInput
		if err := json.Unmarshal(result.Payload, &projection); err != nil {
			return nil, fmt.Errorf("failed to decode yjs projection result: %w", err)
		}
		if projection.ProjectedSequence != input.State.LastUpdateSequence {
			return nil, fmt.Errorf("yjs projection result sequence does not match the claimed document")
		}

		resultBlockPackIdSet[result.BlockPackId] = true
		projectionInputs = append(projectionInputs, dtos.ApplyBlockProjectionDocumentInput{
			BlockPackId: result.BlockPackId,
			Projection:  projection,
		})
	}
	if len(resultBlockPackIdSet) != len(inputByBlockPackId) {
		return nil, fmt.Errorf("incomplete yjs projection batch result")
	}

	return h.blockService.ApplyMany(ctx, projectionInputs)
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
