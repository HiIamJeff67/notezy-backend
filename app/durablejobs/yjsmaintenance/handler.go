package yjsmaintenance

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	inputs "github.com/HiIamJeff67/notezy-backend/app/models/inputs"
	repositories "github.com/HiIamJeff67/notezy-backend/app/models/repositories"
	options "github.com/HiIamJeff67/notezy-backend/app/options"
	realtimetypes "github.com/HiIamJeff67/notezy-backend/app/realtime/types"
	constants "github.com/HiIamJeff67/notezy-backend/shared/constants"
)

type Handler struct {
	db                     *gorm.DB
	blockPackYjsRepository repositories.BlockPackYjsRepositoryInterface
}

func NewHandler(db *gorm.DB) Handler {
	return Handler{
		db:                     db,
		blockPackYjsRepository: repositories.NewBlockPackYjsRepository(),
	}
}

func (h Handler) Handle(
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
