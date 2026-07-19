package yjsmaintenance

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	repositories "github.com/HiIamJeff67/notezy-backend/app/models/repositories"
	schemas "github.com/HiIamJeff67/notezy-backend/app/models/schemas"
	metrics "github.com/HiIamJeff67/notezy-backend/app/monitor/metrics"
	options "github.com/HiIamJeff67/notezy-backend/app/options"
	realtimetypes "github.com/HiIamJeff67/notezy-backend/app/realtime/types"
	constants "github.com/HiIamJeff67/notezy-backend/shared/constants"
)

type Claimer struct {
	db                     *gorm.DB
	blockPackYjsRepository repositories.BlockPackYjsRepositoryInterface
}

func NewClaimer(db *gorm.DB) Claimer {
	return Claimer{
		db:                     db,
		blockPackYjsRepository: repositories.NewBlockPackYjsRepository(),
	}
}

func (c Claimer) ClaimCompactions(ctx context.Context) ([]realtimetypes.YjsCompactionBatchInput, error) {
	tx := c.db.WithContext(ctx).Begin()

	type claimableDocument struct {
		BlockPackId     uuid.UUID  `gorm:"column:block_pack_id"`
		LastCompactedAt *time.Time `gorm:"column:last_compacted_at"`
	}

	var claimableDocuments []claimableDocument
	if err := tx.Model(&schemas.BlockPackYjsDocument{}).
		Select(`"BlockPackYjsDocumentTable".block_pack_id, "BlockPackYjsDocumentTable".last_compacted_at`).
		Joins(`INNER JOIN "BlockPackTable" bp ON bp.id = "BlockPackYjsDocumentTable".block_pack_id AND bp.deleted_at IS NULL`).
		Where(`"BlockPackYjsDocumentTable".deleted_at IS NULL`).
		Where(`"BlockPackYjsDocumentTable".last_update_sequence > "BlockPackYjsDocumentTable".compacted_until_sequence`).
		Order(`"BlockPackYjsDocumentTable".last_update_sequence - "BlockPackYjsDocumentTable".compacted_until_sequence DESC`).
		Order(`"BlockPackYjsDocumentTable".last_compacted_at ASC NULLS FIRST`).
		Order(`"BlockPackYjsDocumentTable".updated_at ASC`).
		Clauses(clause.Locking{
			Strength: "UPDATE",
			Options:  "SKIP LOCKED",
			Table:    clause.Table{Name: schemas.BlockPackYjsDocument{}.TableName()},
		}).
		Limit(constants.YjsMaintenanceMaxDocumentsPerRun).
		Find(&claimableDocuments).Error; err != nil {
		tx.Rollback()

		return nil, err
	}
	if len(claimableDocuments) == 0 {
		if err := tx.Commit().Error; err != nil {
			return nil, err
		}

		return nil, nil
	}

	blockPackIds := make([]uuid.UUID, len(claimableDocuments))
	for index, document := range claimableDocuments {
		blockPackIds[index] = document.BlockPackId
		if document.LastCompactedAt != nil {
			metrics.NotezyMeter.Duration(ctx, "yjs.compaction.age", time.Since(*document.LastCompactedAt))
		}
	}

	_, pairs, err := c.blockPackYjsRepository.BulkCheckPermissionsAndGetManyByIds(
		blockPackIds,
		options.WithTransactionDB(tx),
		options.WithLockingStrength(options.LockingStrengthUpdate),
	)
	if err != nil {
		tx.Rollback()

		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	inputs := make([]realtimetypes.YjsCompactionBatchInput, 0, len(pairs))
	payloadSize := 4
	for _, pair := range pairs {
		updates := make([]realtimetypes.YjsDocumentUpdate, len(pair.Updates))
		for updateIndex, update := range pair.Updates {
			updates[updateIndex] = realtimetypes.YjsDocumentUpdate{
				UpdateSequence: update.UpdateSequence,
				Payload:        update.Payload,
			}
		}

		input := realtimetypes.YjsCompactionBatchInput{
			BlockPackId: pair.Document.BlockPackId,
			Input: realtimetypes.YjsCompactionInput{
				Snapshot:                   pair.Document.Snapshot,
				StateVector:                pair.Document.StateVector,
				BaseCompactedUntilSequence: pair.Document.CompactedUntilSequence,
				CutoffSequence:             pair.Document.LastUpdateSequence,
				Updates:                    updates,
			},
		}

		inputPayload, err := input.MarshalBytes()
		if err != nil {
			return nil, err
		}
		if len(inputPayload)+4 > constants.YjsMaintenanceWorkerMaxPayloadBytes {
			return nil, fmt.Errorf("yjs document %s exceeds the maintenance worker payload limit", input.BlockPackId)
		}
		if payloadSize+4+len(inputPayload) > constants.YjsMaintenanceWorkerMaxPayloadBytes {
			break
		}

		inputs = append(inputs, input)
		payloadSize += 4 + len(inputPayload)
	}
	if len(inputs) == 0 {
		return nil, nil
	}

	return inputs, nil
}

func (c Claimer) ClaimProjections(ctx context.Context) ([]realtimetypes.YjsProjectionBatchInput, error) {
	tx := c.db.WithContext(ctx).Begin()

	type claimableDocument struct {
		BlockPackId uuid.UUID `gorm:"column:block_pack_id"`
	}

	var claimableDocuments []claimableDocument
	if err := tx.Model(&schemas.BlockPackYjsDocument{}).
		Select(`"BlockPackYjsDocumentTable".block_pack_id`).
		Joins(`INNER JOIN "BlockPackTable" bp ON bp.id = "BlockPackYjsDocumentTable".block_pack_id AND bp.deleted_at IS NULL`).
		Where(`"BlockPackYjsDocumentTable".deleted_at IS NULL`).
		Where(`"BlockPackYjsDocumentTable".last_update_sequence > "BlockPackYjsDocumentTable".projected_until_sequence`).
		Order(`"BlockPackYjsDocumentTable".last_update_sequence - "BlockPackYjsDocumentTable".projected_until_sequence DESC`).
		Order(`"BlockPackYjsDocumentTable".last_compacted_at ASC NULLS FIRST`).
		Order(`"BlockPackYjsDocumentTable".updated_at ASC`).
		Clauses(clause.Locking{
			Strength: "UPDATE",
			Options:  "SKIP LOCKED",
			Table:    clause.Table{Name: schemas.BlockPackYjsDocument{}.TableName()},
		}).
		Limit(constants.YjsMaintenanceMaxDocumentsPerRun).
		Find(&claimableDocuments).Error; err != nil {
		tx.Rollback()

		return nil, err
	}
	if len(claimableDocuments) == 0 {
		if err := tx.Commit().Error; err != nil {
			return nil, err
		}

		return nil, nil
	}

	blockPackIds := make([]uuid.UUID, len(claimableDocuments))
	for index, document := range claimableDocuments {
		blockPackIds[index] = document.BlockPackId
	}

	pairs, err := c.blockPackYjsRepository.BulkCheckProjectionDocumentsByIds(
		blockPackIds,
		options.WithTransactionDB(tx),
		options.WithLockingStrength(options.LockingStrengthUpdate),
	)
	if err != nil {
		tx.Rollback()

		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	inputs := make([]realtimetypes.YjsProjectionBatchInput, 0, len(pairs))
	payloadSize := 4
	for _, pair := range pairs {
		updates := make([]realtimetypes.YjsDocumentUpdate, len(pair.Updates))
		for updateIndex, update := range pair.Updates {
			updates[updateIndex] = realtimetypes.YjsDocumentUpdate{
				UpdateSequence: update.UpdateSequence,
				Payload:        update.Payload,
			}
		}

		input := realtimetypes.YjsProjectionBatchInput{
			BlockPackId: pair.Document.BlockPackId,
			State: realtimetypes.YjsDocumentState{
				Snapshot:               pair.Document.Snapshot,
				StateVector:            pair.Document.StateVector,
				LastUpdateSequence:     pair.Document.LastUpdateSequence,
				CompactedUntilSequence: pair.Document.CompactedUntilSequence,
				ProjectedUntilSequence: pair.Document.ProjectedUntilSequence,
				Updates:                updates,
			},
		}

		inputPayload, err := input.MarshalBytes()
		if err != nil {
			return nil, err
		}
		if len(inputPayload)+4 > constants.YjsMaintenanceWorkerMaxPayloadBytes {
			continue
		}
		if payloadSize+4+len(inputPayload) > constants.YjsMaintenanceWorkerMaxPayloadBytes {
			break
		}

		inputs = append(inputs, input)
		payloadSize += 4 + len(inputPayload)
	}
	if len(inputs) == 0 {
		return nil, nil
	}

	return inputs, nil
}
