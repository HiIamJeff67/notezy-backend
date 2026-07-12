package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	adapters "github.com/HiIamJeff67/notezy-backend/app/adapters"
	dtos "github.com/HiIamJeff67/notezy-backend/app/dtos"
	models "github.com/HiIamJeff67/notezy-backend/app/models"
	schemas "github.com/HiIamJeff67/notezy-backend/app/models/schemas"
	scopes "github.com/HiIamJeff67/notezy-backend/app/models/scopes"
	constants "github.com/HiIamJeff67/notezy-backend/shared/constants"
)

type BlockProjectionServiceInterface interface {
	Apply(ctx context.Context, blockPackId uuid.UUID, input dtos.ApplyBlockProjectionInput) (*dtos.ApplyBlockProjectionResult, error)
}

type BlockProjectionService struct {
	db                   *gorm.DB
	editableBlockAdapter adapters.EditableBlockAdapterInterface
}

func NewBlockProjectionService(
	db *gorm.DB,
	editableBlockAdapter adapters.EditableBlockAdapterInterface,
) BlockProjectionServiceInterface {
	if db == nil {
		db = models.NotezyDB
	}
	if editableBlockAdapter == nil {
		editableBlockAdapter = adapters.NewEditableBlockAdapter()
	}

	return &BlockProjectionService{
		db:                   db,
		editableBlockAdapter: editableBlockAdapter,
	}
}

func (s *BlockProjectionService) Apply(
	ctx context.Context,
	blockPackId uuid.UUID,
	input dtos.ApplyBlockProjectionInput,
) (*dtos.ApplyBlockProjectionResult, error) {
	if blockPackId == uuid.Nil {
		return nil, fmt.Errorf("block projection requires a block pack id")
	}
	if input.SchemaId != constants.YjsBlockPackSchemaId ||
		input.SchemaVersion != constants.YjsBlockPackSchemaVersion {
		return nil, fmt.Errorf("block projection source schema is not supported")
	}
	if input.ProjectedSequence < 0 {
		return nil, fmt.Errorf("block projection target update sequence must not be negative")
	}

	flattenedBlocks, _, exception := s.editableBlockAdapter.FlattenManyToRaw(input.Blocks)
	if exception != nil {
		return nil, fmt.Errorf("failed to flatten block projection: %w", exception)
	}

	blockIds := make([]uuid.UUID, len(flattenedBlocks))
	projectedBlocks := make([]schemas.Block, len(flattenedBlocks))
	for index, flattenedBlock := range flattenedBlocks {
		blockIds[index] = flattenedBlock.Id
		projectedBlocks[index] = schemas.Block{
			Id:            flattenedBlock.Id,
			BlockPackId:   blockPackId,
			ParentBlockId: flattenedBlock.ParentBlockId,
			PrevBlockId:   flattenedBlock.PrevBlockId,
			NextBlockId:   flattenedBlock.NextBlockId,
			Type:          flattenedBlock.Type,
			Props:         flattenedBlock.Props,
			Content:       flattenedBlock.Content,
		}
	}

	tx := s.db.WithContext(ctx).Begin()

	lockingStrength := "UPDATE"
	var blockPack schemas.BlockPack
	if err := tx.Model(&schemas.BlockPack{}).
		Select("id").
		Scopes(scopes.Locking(&lockingStrength)).
		Where("id = ? AND deleted_at IS NULL", blockPackId).
		First(&blockPack).Error; err != nil {
		tx.Rollback()

		return nil, fmt.Errorf("failed to lock block pack for projection: %w", err)
	}

	var document schemas.BlockPackYjsDocument
	if err := tx.Model(&schemas.BlockPackYjsDocument{}).
		Scopes(scopes.Locking(&lockingStrength)).
		Where("block_pack_id = ? AND deleted_at IS NULL", blockPackId).
		First(&document).Error; err != nil {
		tx.Rollback()

		return nil, fmt.Errorf("failed to lock yjs document for projection: %w", err)
	}

	if input.ProjectedSequence <= document.ProjectedUntilSequence {
		if err := tx.Commit().Error; err != nil {
			return nil, fmt.Errorf("failed to commit stale block projection: %w", err)
		}

		return &dtos.ApplyBlockProjectionResult{
			Applied:                false,
			ProjectedUntilSequence: document.ProjectedUntilSequence,
		}, nil
	}

	if input.ProjectedSequence > document.LastUpdateSequence {
		tx.Rollback()

		return nil, fmt.Errorf("block projection target update sequence exceeds durable yjs state")
	}

	type existingBlock struct {
		Id          uuid.UUID `gorm:"column:id"`
		BlockPackId uuid.UUID `gorm:"column:block_pack_id"`
	}

	existingBlocks := []existingBlock{}
	if len(blockIds) > 0 {
		if err := tx.Model(&schemas.Block{}).
			Select("id, block_pack_id").
			Scopes(scopes.Locking(&lockingStrength)).
			Where("id IN ?", blockIds).
			Find(&existingBlocks).Error; err != nil {
			tx.Rollback()

			return nil, fmt.Errorf("failed to lock projected blocks: %w", err)
		}
	}

	for _, existingBlock := range existingBlocks {
		if existingBlock.BlockPackId != blockPackId {
			tx.Rollback()

			return nil, fmt.Errorf("block projection contains an id owned by another block pack")
		}
	}

	now := time.Now()

	if len(projectedBlocks) > 0 {
		if err := tx.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "id"}},
			DoUpdates: clause.Assignments(map[string]any{
				"block_pack_id":   blockPackId,
				"parent_block_id": gorm.Expr("EXCLUDED.parent_block_id"),
				"prev_block_id":   gorm.Expr("EXCLUDED.prev_block_id"),
				"next_block_id":   gorm.Expr("EXCLUDED.next_block_id"),
				"type":            gorm.Expr("EXCLUDED.type"),
				"props":           gorm.Expr("EXCLUDED.props"),
				"content":         gorm.Expr("EXCLUDED.content"),
				"updated_at":      now,
			}),
		}).CreateInBatches(&projectedBlocks, constants.MaxBatchCreateBlockSize).Error; err != nil {
			tx.Rollback()

			return nil, fmt.Errorf("failed to bulk upsert block projection: %w", err)
		}
	}

	deleteQuery := tx.Where("block_pack_id = ?", blockPackId)
	if len(blockIds) > 0 {
		deleteQuery = deleteQuery.Where("id NOT IN ?", blockIds)
	}
	if err := deleteQuery.Delete(&schemas.Block{}).Error; err != nil {
		tx.Rollback()

		return nil, fmt.Errorf("failed to delete removed projected blocks: %w", err)
	}

	if err := tx.Model(&schemas.BlockPackYjsDocument{}).
		Where("id = ?", document.Id).
		Updates(map[string]any{
			"projected_until_sequence": input.ProjectedSequence,
			"updated_at":               now,
		}).Error; err != nil {
		tx.Rollback()

		return nil, fmt.Errorf("failed to update block projection checkpoint: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit block projection: %w", err)
	}

	return &dtos.ApplyBlockProjectionResult{
		Applied:                true,
		ProjectedUntilSequence: input.ProjectedSequence,
	}, nil
}
