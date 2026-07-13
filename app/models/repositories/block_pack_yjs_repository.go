package repositories

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	schemas "github.com/HiIamJeff67/notezy-backend/app/models/schemas"
	scopes "github.com/HiIamJeff67/notezy-backend/app/models/scopes"
)

type BlockPackYjsRepositoryInterface interface {
	LoadDocumentAndUpdates(ctx context.Context, blockPackId uuid.UUID) (*schemas.BlockPackYjsDocument, []schemas.BlockPackYjsUpdate, error)
	AppendUpdate(
		ctx context.Context,
		blockPackId uuid.UUID,
		persistenceBatchId uuid.UUID,
		originConnectionId *uuid.UUID,
		payload []byte,
	) (int64, error)
}

type BlockPackYjsRepository struct {
	db *gorm.DB
}

func NewBlockPackYjsRepository(db *gorm.DB) BlockPackYjsRepositoryInterface {
	return &BlockPackYjsRepository{db: db}
}

func (r *BlockPackYjsRepository) LoadDocumentAndUpdates(
	ctx context.Context, blockPackId uuid.UUID,
) (*schemas.BlockPackYjsDocument, []schemas.BlockPackYjsUpdate, error) {
	db := r.db.WithContext(ctx)

	var blockPack schemas.BlockPack
	if err := db.Model(&schemas.BlockPack{}).
		Select("id").
		Where("id = ? AND deleted_at IS NULL", blockPackId).
		First(&blockPack).Error; err != nil {
		return nil, nil, err
	}

	var document schemas.BlockPackYjsDocument
	if err := db.Model(&schemas.BlockPackYjsDocument{}).
		Where("block_pack_id = ? AND deleted_at IS NULL", blockPackId).
		First(&document).Error; err != nil {
		return nil, nil, err
	}

	var updates []schemas.BlockPackYjsUpdate
	if err := db.Model(&schemas.BlockPackYjsUpdate{}).
		Where("block_pack_id = ?", blockPackId).
		Where("update_sequence > ? AND update_sequence <= ?", document.CompactedUntilSequence, document.LastUpdateSequence).
		Order("update_sequence ASC").
		Find(&updates).Error; err != nil {
		return nil, nil, err
	}

	return &document, updates, nil
}

func (r *BlockPackYjsRepository) AppendUpdate(
	ctx context.Context,
	blockPackId uuid.UUID,
	persistenceBatchId uuid.UUID,
	originConnectionId *uuid.UUID,
	payload []byte,
) (int64, error) {
	tx := r.db.WithContext(ctx).Begin()

	lockingStrength := "UPDATE"
	var blockPack schemas.BlockPack
	if err := tx.Model(&schemas.BlockPack{}).
		Select("id").
		Scopes(scopes.Locking(&lockingStrength)).
		Where("id = ? AND deleted_at IS NULL", blockPackId).
		First(&blockPack).Error; err != nil {
		tx.Rollback()

		return 0, err
	}

	var document schemas.BlockPackYjsDocument
	if err := tx.Model(&schemas.BlockPackYjsDocument{}).
		Scopes(scopes.Locking(&lockingStrength)).
		Where("block_pack_id = ? AND deleted_at IS NULL", blockPackId).
		First(&document).Error; err != nil {
		tx.Rollback()

		return 0, err
	}

	var existingUpdate schemas.BlockPackYjsUpdate
	existingUpdateResult := tx.Model(&schemas.BlockPackYjsUpdate{}).
		Select("update_sequence").
		Where("block_pack_id = ? AND persistence_batch_id = ?", blockPackId, persistenceBatchId).
		Limit(1).
		Find(&existingUpdate)
	if existingUpdateResult.Error != nil {
		tx.Rollback()

		return 0, existingUpdateResult.Error
	}
	if existingUpdateResult.RowsAffected > 0 {
		if err := tx.Commit().Error; err != nil {
			return 0, err
		}

		return existingUpdate.UpdateSequence, nil
	}

	nextUpdateSequence := document.LastUpdateSequence + 1
	update := schemas.BlockPackYjsUpdate{
		BlockPackId:        blockPackId,
		UpdateSequence:     nextUpdateSequence,
		PersistenceBatchId: persistenceBatchId,
		Payload:            payload,
		OriginConnectionId: originConnectionId,
	}
	if err := tx.Create(&update).Error; err != nil {
		tx.Rollback()

		return 0, err
	}

	if err := tx.Model(&schemas.BlockPackYjsDocument{}).
		Where("id = ?", document.Id).
		Update("last_update_sequence", nextUpdateSequence).Error; err != nil {
		tx.Rollback()

		return 0, err
	}

	if err := tx.Commit().Error; err != nil {
		return 0, err
	}

	return nextUpdateSequence, nil
}
