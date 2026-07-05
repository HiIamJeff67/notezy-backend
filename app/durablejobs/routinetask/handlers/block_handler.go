package handlers

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	adapters "github.com/HiIamJeff67/notezy-backend/app/adapters"
	dtos "github.com/HiIamJeff67/notezy-backend/app/dtos"
	matchers "github.com/HiIamJeff67/notezy-backend/app/durablejobs/routinetask/handlers/matchers"
	resolvers "github.com/HiIamJeff67/notezy-backend/app/durablejobs/routinetask/handlers/resolvers"
	exceptions "github.com/HiIamJeff67/notezy-backend/app/exceptions"
	inputs "github.com/HiIamJeff67/notezy-backend/app/models/inputs"
	repositories "github.com/HiIamJeff67/notezy-backend/app/models/repositories"
	schemas "github.com/HiIamJeff67/notezy-backend/app/models/schemas"
	enums "github.com/HiIamJeff67/notezy-backend/app/models/schemas/enums"
	options "github.com/HiIamJeff67/notezy-backend/app/options"
	types "github.com/HiIamJeff67/notezy-backend/shared/types"
)

type BlockHandler struct {
	db                   *gorm.DB
	editableBlockAdapter adapters.EditableBlockAdapterInterface
	patternResolver      resolvers.PatternResolverInterface
	templateBlockMatcher matchers.TemplateBlockMatcherInterface
	blockPackRepository  repositories.BlockPackRepositoryInterface
	blockRepository      repositories.BlockRepositoryInterface
}

func NewBlockHandler(
	db *gorm.DB,
	editableBlockAdapter adapters.EditableBlockAdapterInterface,
	patternResolver resolvers.PatternResolverInterface,
	templateBlockMatcher matchers.TemplateBlockMatcherInterface,
	blockPackRepository repositories.BlockPackRepositoryInterface,
	blockRepository repositories.BlockRepositoryInterface,
) BlockHandler {
	if editableBlockAdapter == nil {
		editableBlockAdapter = adapters.NewEditableBlockAdapter()
	}
	if patternResolver == nil {
		patternResolver = resolvers.NewPatternResolver(db, blockRepository, blockPackRepository)
	}
	if templateBlockMatcher == nil {
		templateBlockMatcher = matchers.NewTemplateBlockMatcher()
	}
	return BlockHandler{
		db:                   db,
		editableBlockAdapter: editableBlockAdapter,
		patternResolver:      patternResolver,
		templateBlockMatcher: templateBlockMatcher,
		blockPackRepository:  blockPackRepository,
		blockRepository:      blockRepository,
	}
}

func (h BlockHandler) HandleAppendBlock(ctx context.Context, tasks []schemas.RoutineTask, taskIdToOwnerId map[uuid.UUID]uuid.UUID) ([]bool, *exceptions.Exception) {
	successes := make([]bool, len(tasks))

	for taskIndex, task := range tasks {
		ownerId, exists := taskIdToOwnerId[task.Id]
		if !exists {
			continue
		}
		payload, exception := decodePayload[dtos.AppendBlockRoutineTaskPayload](task)
		if exception != nil {
			continue
		}
		patternValues, exception := h.patternResolver.Resolve(ctx, task, ownerId, payload.Pattern)
		if exception != nil {
			continue
		}
		matchedBlock, exception := h.templateBlockMatcher.MatchArborizedEditableBlock(payload.ArborizedEditableBlock, patternValues)
		if exception != nil {
			continue
		}
		blocks, _, _, exception := flattenArborizedBlock(h.editableBlockAdapter, payload.BlockPackId, &matchedBlock)
		if exception != nil {
			continue
		}

		allowedPermissions := []enums.AccessControlPermission{
			enums.AccessControlPermission_Owner,
			enums.AccessControlPermission_Admin,
			enums.AccessControlPermission_Write,
		}

		tx := h.db.WithContext(ctx).Begin()
		if !h.blockPackRepository.HasPermission(payload.BlockPackId, ownerId, allowedPermissions, options.WithTransactionDB(tx), options.WithOnlyDeleted(types.Ternary_Negative)) {
			tx.Rollback()
			continue
		}
		var tail schemas.Block
		if err := tx.Model(&schemas.Block{}).
			Where("block_pack_id = ? AND parent_block_id IS NULL AND next_block_id IS NULL AND deleted_at IS NULL", payload.BlockPackId).
			First(&tail).Error; err == nil {
			blocks[0].PrevBlockId = &tail.Id
		}
		if err := tx.CreateInBatches(&blocks, 100).Error; err != nil {
			tx.Rollback()
			continue
		}
		if blocks[0].PrevBlockId != nil {
			if err := tx.Model(&schemas.Block{}).Where("id = ?", *blocks[0].PrevBlockId).Update("next_block_id", blocks[0].Id).Error; err != nil {
				tx.Rollback()
				continue
			}
		}
		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			continue
		}
		successes[taskIndex] = true
	}

	return successes, nil
}

func (h BlockHandler) HandleUpdateBlock(ctx context.Context, tasks []schemas.RoutineTask, taskIdToOwnerId map[uuid.UUID]uuid.UUID) ([]bool, *exceptions.Exception) {
	successes := make([]bool, len(tasks))
	bulkInputs := make([]inputs.BulkUpdateBlockInput, 0, len(tasks))
	taskIndexes := make([]int, 0, len(tasks))

	for taskIndex, task := range tasks {
		ownerId, exists := taskIdToOwnerId[task.Id]
		if !exists {
			continue
		}
		payload, exception := decodePayload[dtos.UpdateBlockRoutineTaskPayload](task)
		if exception != nil {
			continue
		}
		if payload.ArborizedEditableBlock == nil {
			continue
		}
		patternValues, exception := h.patternResolver.Resolve(ctx, task, ownerId, payload.Pattern)
		if exception != nil {
			continue
		}
		matchedBlock, exception := h.templateBlockMatcher.MatchArborizedEditableBlock(*payload.ArborizedEditableBlock, patternValues)
		if exception != nil {
			continue
		}
		rawBlocks, _, exception := h.editableBlockAdapter.FlattenToRaw(&matchedBlock)
		if exception != nil || len(rawBlocks) != 1 {
			continue
		}
		blockType := rawBlocks[0].Type
		props := datatypes.JSON(rawBlocks[0].Props)
		content := datatypes.JSON(rawBlocks[0].Content)
		bulkInputs = append(bulkInputs, inputs.BulkUpdateBlockInput{
			UserId: ownerId,
			Id:     payload.BlockId,
			PartialUpdateInput: inputs.PartialUpdateBlockInput{Values: inputs.UpdateBlockInput{
				Type:    &blockType,
				Props:   &props,
				Content: &content,
			}},
		})
		taskIndexes = append(taskIndexes, taskIndex)
	}
	if len(bulkInputs) == 0 {
		return successes, nil
	}

	bulkSuccesses, exception := h.blockRepository.BulkUpdateMany(bulkInputs, options.WithDB(h.db.WithContext(ctx)), options.WithOnlyDeleted(types.Ternary_Negative))
	if exception != nil {
		return successes, exception
	}
	for index, success := range bulkSuccesses {
		successes[taskIndexes[index]] = success
	}

	return successes, nil
}

func (h BlockHandler) HandleResetBlock(ctx context.Context, tasks []schemas.RoutineTask, taskIdToOwnerId map[uuid.UUID]uuid.UUID) ([]bool, *exceptions.Exception) {
	successes := make([]bool, len(tasks))
	blockType := enums.BlockType_Paragraph
	props := datatypes.JSON([]byte("{}"))
	content := datatypes.JSON([]byte("[]"))

	for taskIndex, task := range tasks {
		ownerId, exists := taskIdToOwnerId[task.Id]
		if !exists {
			continue
		}
		payload, exception := decodePayload[dtos.ResetBlockRoutineTaskPayload](task)
		if exception != nil {
			continue
		}
		_, exception = h.blockRepository.UpdateOneById(payload.BlockId, ownerId, inputs.PartialUpdateBlockInput{
			Values: inputs.UpdateBlockInput{Type: &blockType, Props: &props, Content: &content},
		}, options.WithDB(h.db.WithContext(ctx)), options.WithOnlyDeleted(types.Ternary_Negative))
		successes[taskIndex] = exception == nil
	}

	return successes, nil
}
