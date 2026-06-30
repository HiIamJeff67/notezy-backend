package handlers

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	adapters "github.com/HiIamJeff67/notezy-backend/app/adapters"
	dtos "github.com/HiIamJeff67/notezy-backend/app/dtos"
	exceptions "github.com/HiIamJeff67/notezy-backend/app/exceptions"
	inputs "github.com/HiIamJeff67/notezy-backend/app/models/inputs"
	repositories "github.com/HiIamJeff67/notezy-backend/app/models/repositories"
	schemas "github.com/HiIamJeff67/notezy-backend/app/models/schemas"
	enums "github.com/HiIamJeff67/notezy-backend/app/models/schemas/enums"
	options "github.com/HiIamJeff67/notezy-backend/app/options"
	constants "github.com/HiIamJeff67/notezy-backend/shared/constants"
	types "github.com/HiIamJeff67/notezy-backend/shared/types"
)

type BlockHandler struct {
	db                   *gorm.DB
	editableBlockAdapter adapters.EditableBlockAdapterInterface
	blockPackRepository  repositories.BlockPackRepositoryInterface
	blockGroupRepository repositories.BlockGroupRepositoryInterface
	blockRepository      repositories.BlockRepositoryInterface
}

func NewBlockHandler(
	db *gorm.DB,
	editableBlockAdapter adapters.EditableBlockAdapterInterface,
	blockPackRepository repositories.BlockPackRepositoryInterface,
	blockGroupRepository repositories.BlockGroupRepositoryInterface,
	blockRepository repositories.BlockRepositoryInterface,
) BlockHandler {
	if editableBlockAdapter == nil {
		editableBlockAdapter = adapters.NewEditableBlockAdapter()
	}
	return BlockHandler{
		db:                   db,
		editableBlockAdapter: editableBlockAdapter,
		blockPackRepository:  blockPackRepository,
		blockGroupRepository: blockGroupRepository,
		blockRepository:      blockRepository,
	}
}

func (h BlockHandler) HandleAppendBlock(
	ctx context.Context,
	tasks []schemas.RoutineTask,
	taskIdToOwnerId map[uuid.UUID]uuid.UUID,
) ([]bool, *exceptions.Exception) {
	successes := make([]bool, len(tasks))
	preparedInputs := make([]struct {
		TaskIndex    int
		UserId       uuid.UUID
		BlockPackId  uuid.UUID
		BlockGroupId uuid.UUID
		Size         int64
		Blocks       []inputs.CreateBlockInput
	}, 0, len(tasks))
	blockPackIds := make([]uuid.UUID, 0, len(tasks))

	for taskIndex, task := range tasks {
		ownerId, exists := taskIdToOwnerId[task.Id]
		if !exists {
			continue
		}

		payload, exception := decodePayload[dtos.AppendBlockRoutineTaskPayload](task)
		if exception != nil {
			continue
		}

		blockGroupId := uuid.New()
		blocks, _, totalSize, exception := flattenArborizedBlock(
			h.editableBlockAdapter,
			blockGroupId,
			&payload.ArborizedEditableBlock,
		)
		if exception != nil {
			continue
		}

		blockInputs := make([]inputs.CreateBlockInput, len(blocks))
		for index, block := range blocks {
			blockInputs[index] = inputs.CreateBlockInput{
				Id:            block.Id,
				ParentBlockId: block.ParentBlockId,
				Type:          block.Type,
				Props:         block.Props,
				Content:       block.Content,
			}
		}

		preparedInputs = append(preparedInputs, struct {
			TaskIndex    int
			UserId       uuid.UUID
			BlockPackId  uuid.UUID
			BlockGroupId uuid.UUID
			Size         int64
			Blocks       []inputs.CreateBlockInput
		}{
			TaskIndex:    taskIndex,
			UserId:       ownerId,
			BlockPackId:  payload.BlockPackId,
			BlockGroupId: blockGroupId,
			Size:         totalSize,
			Blocks:       blockInputs,
		})
		blockPackIds = append(blockPackIds, payload.BlockPackId)
	}
	if len(preparedInputs) == 0 {
		return successes, nil
	}

	var blockPacks []struct {
		Id                uuid.UUID  `gorm:"column:id"`
		FinalBlockGroupId *uuid.UUID `gorm:"column:final_block_group_id"`
	}
	if err := h.db.WithContext(ctx).
		Model(&schemas.BlockPack{}).
		Select("id, final_block_group_id").
		Where("id IN ? AND deleted_at IS NULL", blockPackIds).
		Find(&blockPacks).Error; err != nil {
		return successes, exceptions.BlockPack.NotFound().WithOrigin(err)
	}
	finalBlockGroupIdByBlockPackId := make(map[uuid.UUID]*uuid.UUID, len(blockPacks))
	for _, blockPack := range blockPacks {
		finalBlockGroupIdByBlockPackId[blockPack.Id] = blockPack.FinalBlockGroupId
	}

	lastBlockGroupIdByBlockPackId := make(map[uuid.UUID]uuid.UUID)
	blockGroupInputs := make([]inputs.BulkCreateBlockGroupInput, 0, len(preparedInputs))
	blockGroupSizeUpdates := make([]inputs.BulkUpdateBlockGroupInput, 0, len(preparedInputs))
	blockContents := make([]inputs.BulkCreateBlockGroupContentInput, 0, len(preparedInputs))
	taskIndexes := make([]int, 0, len(preparedInputs))
	for _, preparedInput := range preparedInputs {
		finalBlockGroupId, exists := finalBlockGroupIdByBlockPackId[preparedInput.BlockPackId]
		if !exists {
			continue
		}

		var prevBlockGroupId *uuid.UUID
		if lastBlockGroupId, exists := lastBlockGroupIdByBlockPackId[preparedInput.BlockPackId]; exists {
			prev := lastBlockGroupId
			prevBlockGroupId = &prev
		} else if finalBlockGroupId != nil {
			prev := *finalBlockGroupId
			prevBlockGroupId = &prev
		}

		blockGroupId := preparedInput.BlockGroupId
		blockGroupInputs = append(blockGroupInputs, inputs.BulkCreateBlockGroupInput{
			UserId:           preparedInput.UserId,
			BlockPackId:      preparedInput.BlockPackId,
			BlockGroupId:     &blockGroupId,
			PrevBlockGroupId: prevBlockGroupId,
		})
		size := preparedInput.Size
		blockGroupSizeUpdates = append(blockGroupSizeUpdates, inputs.BulkUpdateBlockGroupInput{
			UserId: preparedInput.UserId,
			Id:     preparedInput.BlockGroupId,
			PartialUpdateInput: inputs.PartialUpdateBlockGroupInput{
				Values: inputs.UpdateBlockGroupInput{
					Size: &size,
				},
			},
		})
		blockContents = append(blockContents, inputs.BulkCreateBlockGroupContentInput{
			UserId:       preparedInput.UserId,
			BlockGroupId: preparedInput.BlockGroupId,
			Blocks:       preparedInput.Blocks,
		})
		taskIndexes = append(taskIndexes, preparedInput.TaskIndex)
		lastBlockGroupIdByBlockPackId[preparedInput.BlockPackId] = preparedInput.BlockGroupId
	}
	if len(blockGroupInputs) == 0 {
		return successes, nil
	}

	tx := h.db.WithContext(ctx).Begin()

	groupSuccesses, exception := h.blockGroupRepository.BulkCreateMany(
		blockGroupInputs,
		options.WithTransactionDB(tx),
		options.WithLockingStrength(options.LockingStrengthNoKeyUpdate),
		options.WithOnlyDeleted(types.Ternary_Negative),
	)
	if exception != nil {
		tx.Rollback()
		return successes, exception
	}

	successfulSizeUpdates := make([]inputs.BulkUpdateBlockGroupInput, 0, len(blockGroupSizeUpdates))
	successfulBlockContents := make([]inputs.BulkCreateBlockGroupContentInput, 0, len(blockContents))
	successfulTaskIndexes := make([]int, 0, len(taskIndexes))
	for index, success := range groupSuccesses {
		if success {
			successfulSizeUpdates = append(successfulSizeUpdates, blockGroupSizeUpdates[index])
			successfulBlockContents = append(successfulBlockContents, blockContents[index])
			successfulTaskIndexes = append(successfulTaskIndexes, taskIndexes[index])
		}
	}

	if len(successfulSizeUpdates) > 0 {
		sizeSuccesses, exception := h.blockGroupRepository.BulkUpdateMany(
			successfulSizeUpdates,
			options.WithTransactionDB(tx),
			options.WithLockingStrength(options.LockingStrengthNoKeyUpdate),
			options.WithOnlyDeleted(types.Ternary_Negative),
		)
		if exception != nil {
			tx.Rollback()
			return successes, exception
		}
		nextBlockContents := make([]inputs.BulkCreateBlockGroupContentInput, 0, len(successfulBlockContents))
		nextTaskIndexes := make([]int, 0, len(successfulTaskIndexes))
		for index, success := range sizeSuccesses {
			if success {
				nextBlockContents = append(nextBlockContents, successfulBlockContents[index])
				nextTaskIndexes = append(nextTaskIndexes, successfulTaskIndexes[index])
			}
		}
		successfulBlockContents = nextBlockContents
		successfulTaskIndexes = nextTaskIndexes
	}

	if len(successfulBlockContents) > 0 {
		blockSuccesses, exception := h.blockRepository.BulkCreateMany(
			successfulBlockContents,
			options.WithTransactionDB(tx),
			options.WithLockingStrength(options.LockingStrengthNoKeyUpdate),
			options.WithBatchSize(constants.MaxBatchCreateBlockSize),
			options.WithOnlyDeleted(types.Ternary_Negative),
		)
		if exception != nil {
			tx.Rollback()
			return successes, exception
		}
		nextTaskIndexes := make([]int, 0, len(successfulTaskIndexes))
		for index, success := range blockSuccesses {
			if success {
				nextTaskIndexes = append(nextTaskIndexes, successfulTaskIndexes[index])
			}
		}
		successfulTaskIndexes = nextTaskIndexes
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return successes, exceptions.Block.FailedToCommitTransaction().WithOrigin(err)
	}

	for _, taskIndex := range successfulTaskIndexes {
		successes[taskIndex] = true
	}

	return successes, nil
}

func (h BlockHandler) HandleUpdateBlock(
	ctx context.Context,
	tasks []schemas.RoutineTask,
	taskIdToOwnerId map[uuid.UUID]uuid.UUID,
) ([]bool, *exceptions.Exception) {
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

		blockGroupId := uuid.New()
		blocks, _, _, exception := flattenArborizedBlock(
			h.editableBlockAdapter,
			blockGroupId,
			payload.ArborizedEditableBlock,
		)
		if exception != nil {
			continue
		}
		if len(blocks) != 1 {
			continue
		}

		blockType := blocks[0].Type
		props := datatypes.JSON(blocks[0].Props)
		content := datatypes.JSON(blocks[0].Content)
		bulkInputs = append(bulkInputs, inputs.BulkUpdateBlockInput{
			UserId: ownerId,
			Id:     payload.BlockId,
			PartialUpdateInput: inputs.PartialUpdateBlockInput{
				Values: inputs.UpdateBlockInput{
					Type:    &blockType,
					Props:   &props,
					Content: &content,
				},
			},
		})
		taskIndexes = append(taskIndexes, taskIndex)
	}

	if len(bulkInputs) == 0 {
		return successes, nil
	}
	bulkSuccesses, exception := h.blockRepository.BulkUpdateMany(
		bulkInputs,
		options.WithDB(h.db.WithContext(ctx)),
		options.WithLockingStrength(options.LockingStrengthNoKeyUpdate),
		options.WithOnlyDeleted(types.Ternary_Negative),
	)
	if exception != nil {
		return successes, exception
	}
	for index, success := range bulkSuccesses {
		successes[taskIndexes[index]] = success
	}

	return successes, nil
}

func (h BlockHandler) HandleResetBlock(
	ctx context.Context,
	tasks []schemas.RoutineTask,
	taskIdToOwnerId map[uuid.UUID]uuid.UUID,
) ([]bool, *exceptions.Exception) {
	successes := make([]bool, len(tasks))
	bulkInputs := make([]inputs.BulkUpdateBlockInput, 0, len(tasks))
	taskIndexes := make([]int, 0, len(tasks))

	for taskIndex, task := range tasks {
		ownerId, exists := taskIdToOwnerId[task.Id]
		if !exists {
			continue
		}

		payload, exception := decodePayload[dtos.ResetBlockRoutineTaskPayload](task)
		if exception != nil {
			continue
		}
		blockType := enums.BlockType_Paragraph
		props := datatypes.JSON([]byte("{}"))
		content := datatypes.JSON([]byte("[]"))
		bulkInputs = append(bulkInputs, inputs.BulkUpdateBlockInput{
			UserId: ownerId,
			Id:     payload.BlockId,
			PartialUpdateInput: inputs.PartialUpdateBlockInput{
				Values: inputs.UpdateBlockInput{
					Type:    &blockType,
					Props:   &props,
					Content: &content,
				},
			},
		})
		taskIndexes = append(taskIndexes, taskIndex)
	}

	if len(bulkInputs) == 0 {
		return successes, nil
	}
	bulkSuccesses, exception := h.blockRepository.BulkUpdateMany(
		bulkInputs,
		options.WithDB(h.db.WithContext(ctx)),
		options.WithLockingStrength(options.LockingStrengthNoKeyUpdate),
		options.WithOnlyDeleted(types.Ternary_Negative),
	)
	if exception != nil {
		return successes, exception
	}
	for index, success := range bulkSuccesses {
		successes[taskIndexes[index]] = success
	}

	return successes, nil
}
