package handlers

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	adapters "github.com/HiIamJeff67/notezy-backend/app/adapters"
	dtos "github.com/HiIamJeff67/notezy-backend/app/dtos"
	matchers "github.com/HiIamJeff67/notezy-backend/app/durablejobs/routinetask/handlers/matchers"
	exceptions "github.com/HiIamJeff67/notezy-backend/app/exceptions"
	inputs "github.com/HiIamJeff67/notezy-backend/app/models/inputs"
	repositories "github.com/HiIamJeff67/notezy-backend/app/models/repositories"
	schemas "github.com/HiIamJeff67/notezy-backend/app/models/schemas"
	options "github.com/HiIamJeff67/notezy-backend/app/options"
	types "github.com/HiIamJeff67/notezy-backend/shared/types"
)

type BlockPackHandler struct {
	db                   *gorm.DB
	editableBlockAdapter adapters.EditableBlockAdapterInterface
	namePatternMatcher   matchers.NamePatternMatcherInterface
	blockPackRepository  repositories.BlockPackRepositoryInterface
	blockRepository      repositories.BlockRepositoryInterface
}

func NewBlockPackHandler(
	db *gorm.DB,
	editableBlockAdapter adapters.EditableBlockAdapterInterface,
	blockPackRepository repositories.BlockPackRepositoryInterface,
	blockRepository repositories.BlockRepositoryInterface,
) BlockPackHandler {
	if editableBlockAdapter == nil {
		editableBlockAdapter = adapters.NewEditableBlockAdapter()
	}
	return BlockPackHandler{
		db:                   db,
		editableBlockAdapter: editableBlockAdapter,
		namePatternMatcher:   matchers.NewNamePatternMatcher(),
		blockPackRepository:  blockPackRepository,
		blockRepository:      blockRepository,
	}
}

func (h BlockPackHandler) HandleCreateBlockPack(ctx context.Context, tasks []schemas.RoutineTask, taskIdToOwnerId map[uuid.UUID]uuid.UUID) ([]bool, *exceptions.Exception) {
	successes := make([]bool, len(tasks))
	for taskIndex, task := range tasks {
		ownerId, exists := taskIdToOwnerId[task.Id]
		if !exists {
			continue
		}
		payload, exception := decodePayload[dtos.CreateBlockPackRoutineTaskPayload](task)
		if exception != nil {
			continue
		}
		name, exception := h.namePatternMatcher.Match(payload.Template.Name, payload.Template.NamePattern, task)
		if exception != nil {
			continue
		}
		blockPackId := uuid.New()
		tx := h.db.WithContext(ctx).Begin()
		blockPackSuccesses, exception := h.blockPackRepository.BulkCreateMany([]inputs.BulkCreateBlockPackInput{{
			UserId:              ownerId,
			Id:                  &blockPackId,
			ParentSubShelfId:    payload.TargetSubShelfId,
			Name:                name,
			Icon:                payload.Template.Icon,
			HeaderBackgroundURL: payload.Template.HeaderBackgroundURL,
		}}, options.WithTransactionDB(tx), options.WithOnlyDeleted(types.Ternary_Negative))
		if exception != nil || len(blockPackSuccesses) == 0 || !blockPackSuccesses[0] {
			tx.Rollback()
			continue
		}
		var prevRootId *uuid.UUID
		taskFailed := false
		for _, block := range payload.Template.Blocks {
			blocks, _, _, exception := flattenArborizedBlock(h.editableBlockAdapter, blockPackId, &block.ArborizedEditableBlock)
			if exception != nil || len(blocks) == 0 {
				tx.Rollback()
				taskFailed = true
				break
			}
			blocks[0].PrevBlockId = prevRootId
			if err := tx.CreateInBatches(&blocks, 100).Error; err != nil {
				tx.Rollback()
				taskFailed = true
				break
			}
			if prevRootId != nil {
				if err := tx.Model(&schemas.Block{}).Where("id = ?", *prevRootId).Update("next_block_id", blocks[0].Id).Error; err != nil {
					tx.Rollback()
					taskFailed = true
					break
				}
			}
			prevRootId = &blocks[0].Id
		}
		if taskFailed {
			continue
		}
		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			continue
		}
		successes[taskIndex] = true
	}
	return successes, nil
}

func (h BlockPackHandler) HandleUpdateBlockPack(ctx context.Context, tasks []schemas.RoutineTask, taskIdToOwnerId map[uuid.UUID]uuid.UUID) ([]bool, *exceptions.Exception) {
	successes := make([]bool, len(tasks))
	preparedInputs := make([]inputs.BulkUpdateBlockInput, 0)
	taskIndexes := make([]int, 0)
	pairPlaceholders := make([]string, 0)
	pairArgs := make([]any, 0)
	for taskIndex, task := range tasks {
		ownerId, exists := taskIdToOwnerId[task.Id]
		if !exists {
			continue
		}
		payload, exception := decodePayload[dtos.UpdateBlockPackRoutineTaskPayload](task)
		if exception != nil {
			continue
		}
		for _, block := range payload.UpdatedBlocks {
			flattenedBlocks, _, _, exception := flattenArborizedBlock(h.editableBlockAdapter, payload.BlockPackId, block.ArborizedEditableBlock)
			if exception != nil || len(flattenedBlocks) != 1 {
				continue
			}
			blockType := flattenedBlocks[0].Type
			props := datatypes.JSON(flattenedBlocks[0].Props)
			content := datatypes.JSON(flattenedBlocks[0].Content)
			pairPlaceholders = append(pairPlaceholders, "(?::uuid, ?::uuid)")
			pairArgs = append(pairArgs, block.BlockId, payload.BlockPackId)
			preparedInputs = append(preparedInputs, inputs.BulkUpdateBlockInput{
				UserId: ownerId,
				Id:     block.BlockId,
				PartialUpdateInput: inputs.PartialUpdateBlockInput{Values: inputs.UpdateBlockInput{
					Type:    &blockType,
					Props:   &props,
					Content: &content,
				}},
			})
			taskIndexes = append(taskIndexes, taskIndex)
		}
	}
	if len(preparedInputs) == 0 {
		return successes, nil
	}
	var validRows []struct {
		BlockId     uuid.UUID `gorm:"column:block_id"`
		BlockPackId uuid.UUID `gorm:"column:block_pack_id"`
	}
	sql := fmt.Sprintf(`
		WITH pairs(block_id, block_pack_id) AS (VALUES %s)
		SELECT p.block_id::uuid, p.block_pack_id::uuid
		FROM pairs p
		INNER JOIN "BlockTable" b ON b.id = p.block_id::uuid AND b.block_pack_id = p.block_pack_id::uuid AND b.deleted_at IS NULL
	`, strings.Join(pairPlaceholders, ","))
	if err := h.db.WithContext(ctx).Raw(sql, pairArgs...).Scan(&validRows).Error; err != nil {
		return successes, exceptions.Block.NotFound().WithOrigin(err)
	}
	valid := make(map[[2]uuid.UUID]bool, len(validRows))
	for _, row := range validRows {
		valid[[2]uuid.UUID{row.BlockId, row.BlockPackId}] = true
	}
	filteredInputs := make([]inputs.BulkUpdateBlockInput, 0, len(preparedInputs))
	filteredTaskIndexes := make([]int, 0, len(taskIndexes))
	for index, input := range preparedInputs {
		blockPackId := pairArgs[index*2+1].(uuid.UUID)
		if valid[[2]uuid.UUID{input.Id, blockPackId}] {
			filteredInputs = append(filteredInputs, input)
			filteredTaskIndexes = append(filteredTaskIndexes, taskIndexes[index])
		}
	}
	if len(filteredInputs) == 0 {
		return successes, nil
	}
	bulkSuccesses, exception := h.blockRepository.BulkUpdateMany(filteredInputs, options.WithDB(h.db.WithContext(ctx)), options.WithOnlyDeleted(types.Ternary_Negative))
	if exception != nil {
		return successes, exception
	}
	for index, success := range bulkSuccesses {
		successes[filteredTaskIndexes[index]] = success
	}
	return successes, nil
}

func (h BlockPackHandler) HandleResetBlockPack(ctx context.Context, tasks []schemas.RoutineTask, taskIdToOwnerId map[uuid.UUID]uuid.UUID) ([]bool, *exceptions.Exception) {
	successes := make([]bool, len(tasks))
	for taskIndex, task := range tasks {
		ownerId, exists := taskIdToOwnerId[task.Id]
		if !exists {
			continue
		}
		payload, exception := decodePayload[dtos.ResetBlockPackRoutineTaskPayload](task)
		if exception != nil {
			continue
		}
		if !h.blockPackRepository.HasPermission(payload.BlockPackId, ownerId, writePermissionsForJobs(), options.WithDB(h.db.WithContext(ctx)), options.WithOnlyDeleted(types.Ternary_Negative)) {
			continue
		}
		if err := h.db.WithContext(ctx).Model(&schemas.Block{}).
			Where("block_pack_id = ? AND deleted_at IS NULL", payload.BlockPackId).
			Updates(map[string]any{"deleted_at": time.Now(), "prev_block_id": nil, "next_block_id": nil}).Error; err != nil {
			continue
		}
		successes[taskIndex] = true
	}
	return successes, nil
}
