package handlers

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	adapters "github.com/HiIamJeff67/notezy-backend/app/adapters"
	dtos "github.com/HiIamJeff67/notezy-backend/app/dtos"
	exceptions "github.com/HiIamJeff67/notezy-backend/app/exceptions"
	inputs "github.com/HiIamJeff67/notezy-backend/app/models/inputs"
	repositories "github.com/HiIamJeff67/notezy-backend/app/models/repositories"
	schemas "github.com/HiIamJeff67/notezy-backend/app/models/schemas"
	options "github.com/HiIamJeff67/notezy-backend/app/options"
	constants "github.com/HiIamJeff67/notezy-backend/shared/constants"
	types "github.com/HiIamJeff67/notezy-backend/shared/types"
)

type BlockPackHandler struct {
	db                   *gorm.DB
	editableBlockAdapter adapters.EditableBlockAdapterInterface
	blockPackRepository  repositories.BlockPackRepositoryInterface
	blockGroupRepository repositories.BlockGroupRepositoryInterface
	blockRepository      repositories.BlockRepositoryInterface
}

func NewBlockPackHandler(
	db *gorm.DB,
	editableBlockAdapter adapters.EditableBlockAdapterInterface,
	blockPackRepository repositories.BlockPackRepositoryInterface,
	blockGroupRepository repositories.BlockGroupRepositoryInterface,
	blockRepository repositories.BlockRepositoryInterface,
) BlockPackHandler {
	if editableBlockAdapter == nil {
		editableBlockAdapter = adapters.NewEditableBlockAdapter()
	}
	return BlockPackHandler{
		db:                   db,
		editableBlockAdapter: editableBlockAdapter,
		blockPackRepository:  blockPackRepository,
		blockGroupRepository: blockGroupRepository,
		blockRepository:      blockRepository,
	}
}

func (h BlockPackHandler) HandleCreateBlockPack(
	ctx context.Context,
	tasks []schemas.RoutineTask,
	taskIdToOwnerId map[uuid.UUID]uuid.UUID,
) ([]bool, *exceptions.Exception) {
	successes := make([]bool, len(tasks))
	blockPackInputs := make([]inputs.BulkCreateBlockPackInput, 0, len(tasks))
	blockGroupInputsByTask := make([][]inputs.BulkCreateBlockGroupInput, 0, len(tasks))
	blockGroupSizeUpdatesByTask := make([][]inputs.BulkUpdateBlockGroupInput, 0, len(tasks))
	blockContentsByTask := make([][]inputs.BulkCreateBlockGroupContentInput, 0, len(tasks))
	taskIndexes := make([]int, 0, len(tasks))

	for taskIndex, task := range tasks {
		ownerId, exists := taskIdToOwnerId[task.Id]
		if !exists {
			continue
		}

		payload, exception := decodePayload[dtos.CreateBlockPackRoutineTaskPayload](task)
		if exception != nil {
			continue
		}

		blockPackId := uuid.New()
		taskBlockGroupInputs := make([]inputs.BulkCreateBlockGroupInput, 0, len(payload.Template.BlockGroups))
		taskBlockGroupSizeUpdates := make([]inputs.BulkUpdateBlockGroupInput, 0, len(payload.Template.BlockGroups))
		taskBlockContents := make([]inputs.BulkCreateBlockGroupContentInput, 0, len(payload.Template.BlockGroups))
		isTaskValid := true

		blockGroupIdsByClientId := make(map[string]uuid.UUID, len(payload.Template.BlockGroups))
		for _, blockGroup := range payload.Template.BlockGroups {
			blockGroupIdsByClientId[blockGroup.ClientId] = uuid.New()
		}

		for _, blockGroup := range payload.Template.BlockGroups {
			blockGroupId := blockGroupIdsByClientId[blockGroup.ClientId]
			var prevBlockGroupId *uuid.UUID
			if blockGroup.PrevClientId != nil {
				id := blockGroupIdsByClientId[*blockGroup.PrevClientId]
				prevBlockGroupId = &id
			}

			flattenedBlocks, _, totalSize, exception := flattenArborizedBlock(
				h.editableBlockAdapter,
				blockGroupId,
				&blockGroup.ArborizedEditableBlock,
			)
			if exception != nil {
				isTaskValid = false
				continue
			}

			taskBlockGroupInputs = append(taskBlockGroupInputs, inputs.BulkCreateBlockGroupInput{
				UserId:           ownerId,
				BlockPackId:      blockPackId,
				BlockGroupId:     &blockGroupId,
				PrevBlockGroupId: prevBlockGroupId,
			})
			size := totalSize
			taskBlockGroupSizeUpdates = append(taskBlockGroupSizeUpdates, inputs.BulkUpdateBlockGroupInput{
				UserId: ownerId,
				Id:     blockGroupId,
				PartialUpdateInput: inputs.PartialUpdateBlockGroupInput{
					Values: inputs.UpdateBlockGroupInput{Size: &size},
				},
			})

			blockInputs := make([]inputs.CreateBlockInput, len(flattenedBlocks))
			for index, block := range flattenedBlocks {
				blockInputs[index] = inputs.CreateBlockInput{
					Id:            block.Id,
					ParentBlockId: block.ParentBlockId,
					Type:          block.Type,
					Props:         block.Props,
					Content:       block.Content,
				}
			}
			taskBlockContents = append(taskBlockContents, inputs.BulkCreateBlockGroupContentInput{
				UserId:       ownerId,
				BlockGroupId: blockGroupId,
				Blocks:       blockInputs,
			})
		}
		if !isTaskValid {
			continue
		}

		blockPackInputs = append(blockPackInputs, inputs.BulkCreateBlockPackInput{
			UserId:              ownerId,
			Id:                  &blockPackId,
			ParentSubShelfId:    payload.TargetSubShelfId,
			Name:                payload.Template.Name,
			Icon:                payload.Template.Icon,
			HeaderBackgroundURL: payload.Template.HeaderBackgroundURL,
		})
		blockGroupInputsByTask = append(blockGroupInputsByTask, taskBlockGroupInputs)
		blockGroupSizeUpdatesByTask = append(blockGroupSizeUpdatesByTask, taskBlockGroupSizeUpdates)
		blockContentsByTask = append(blockContentsByTask, taskBlockContents)
		taskIndexes = append(taskIndexes, taskIndex)
	}
	if len(blockPackInputs) == 0 {
		return successes, nil
	}

	tx := h.db.WithContext(ctx).Begin()

	blockPackSuccesses, exception := h.blockPackRepository.BulkCreateMany(
		blockPackInputs,
		options.WithTransactionDB(tx),
		options.WithLockingStrength(options.LockingStrengthNoKeyUpdate),
		options.WithOnlyDeleted(types.Ternary_Negative),
	)
	if exception != nil {
		tx.Rollback()
		return successes, exception
	}

	blockGroupInputs := make([]inputs.BulkCreateBlockGroupInput, 0)
	blockGroupSizeUpdates := make([]inputs.BulkUpdateBlockGroupInput, 0)
	blockContents := make([]inputs.BulkCreateBlockGroupContentInput, 0)
	successfulTaskIndexes := make([]int, 0, len(taskIndexes))
	for index, success := range blockPackSuccesses {
		if success {
			blockGroupInputs = append(blockGroupInputs, blockGroupInputsByTask[index]...)
			blockGroupSizeUpdates = append(blockGroupSizeUpdates, blockGroupSizeUpdatesByTask[index]...)
			blockContents = append(blockContents, blockContentsByTask[index]...)
			successfulTaskIndexes = append(successfulTaskIndexes, taskIndexes[index])
		}
	}

	if len(blockGroupInputs) > 0 {
		blockGroupSuccesses, exception := h.blockGroupRepository.BulkCreateMany(
			blockGroupInputs,
			options.WithTransactionDB(tx),
			options.WithLockingStrength(options.LockingStrengthNoKeyUpdate),
			options.WithOnlyDeleted(types.Ternary_Negative),
		)
		if exception != nil {
			tx.Rollback()
			return successes, exception
		}
		for _, success := range blockGroupSuccesses {
			if !success {
				tx.Rollback()
				return successes, nil
			}
		}
	}

	if len(blockGroupSizeUpdates) > 0 {
		sizeSuccesses, exception := h.blockGroupRepository.BulkUpdateMany(
			blockGroupSizeUpdates,
			options.WithTransactionDB(tx),
			options.WithLockingStrength(options.LockingStrengthNoKeyUpdate),
			options.WithOnlyDeleted(types.Ternary_Negative),
		)
		if exception != nil {
			tx.Rollback()
			return successes, exception
		}
		for _, success := range sizeSuccesses {
			if !success {
				tx.Rollback()
				return successes, nil
			}
		}
	}

	if len(blockContents) > 0 {
		blockSuccesses, exception := h.blockRepository.BulkCreateMany(
			blockContents,
			options.WithTransactionDB(tx),
			options.WithLockingStrength(options.LockingStrengthNoKeyUpdate),
			options.WithBatchSize(constants.MaxBatchCreateBlockSize),
			options.WithOnlyDeleted(types.Ternary_Negative),
		)
		if exception != nil {
			tx.Rollback()
			return successes, exception
		}
		for _, success := range blockSuccesses {
			if !success {
				tx.Rollback()
				return successes, nil
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return successes, exceptions.BlockPack.FailedToCommitTransaction().WithOrigin(err)
	}

	for _, taskIndex := range successfulTaskIndexes {
		successes[taskIndex] = true
	}

	return successes, nil
}

func (h BlockPackHandler) HandleUpdateBlockPack(
	ctx context.Context,
	tasks []schemas.RoutineTask,
	taskIdToOwnerId map[uuid.UUID]uuid.UUID,
) ([]bool, *exceptions.Exception) {
	successes := make([]bool, len(tasks))
	preparedInputs := make([]struct {
		TaskIndex   int
		UserId      uuid.UUID
		BlockId     uuid.UUID
		BlockPackId uuid.UUID
		Input       inputs.BulkUpdateBlockInput
	}, 0)
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

		isTaskValid := true
		taskPreparedInputs := make([]struct {
			TaskIndex   int
			UserId      uuid.UUID
			BlockId     uuid.UUID
			BlockPackId uuid.UUID
			Input       inputs.BulkUpdateBlockInput
		}, 0, len(payload.UpdatedBlocks))
		taskPairPlaceholders := make([]string, 0, len(payload.UpdatedBlocks))
		taskPairArgs := make([]any, 0, len(payload.UpdatedBlocks)*2)
		for _, block := range payload.UpdatedBlocks {
			blockGroupId := uuid.New()
			flattenedBlocks, _, _, exception := flattenArborizedBlock(
				h.editableBlockAdapter,
				blockGroupId,
				block.ArborizedEditableBlock,
			)
			if exception != nil {
				isTaskValid = false
				continue
			}
			if len(flattenedBlocks) != 1 {
				isTaskValid = false
				continue
			}

			blockType := flattenedBlocks[0].Type
			props := datatypes.JSON(flattenedBlocks[0].Props)
			content := datatypes.JSON(flattenedBlocks[0].Content)
			taskPairPlaceholders = append(taskPairPlaceholders, "(?::uuid, ?::uuid)")
			taskPairArgs = append(taskPairArgs, block.BlockId, payload.BlockPackId)
			taskPreparedInputs = append(taskPreparedInputs, struct {
				TaskIndex   int
				UserId      uuid.UUID
				BlockId     uuid.UUID
				BlockPackId uuid.UUID
				Input       inputs.BulkUpdateBlockInput
			}{
				TaskIndex:   taskIndex,
				UserId:      ownerId,
				BlockId:     block.BlockId,
				BlockPackId: payload.BlockPackId,
				Input: inputs.BulkUpdateBlockInput{
					UserId: ownerId,
					Id:     block.BlockId,
					PartialUpdateInput: inputs.PartialUpdateBlockInput{
						Values: inputs.UpdateBlockInput{
							Type:    &blockType,
							Props:   &props,
							Content: &content,
						},
					},
				},
			})
		}
		if !isTaskValid {
			continue
		}
		pairPlaceholders = append(pairPlaceholders, taskPairPlaceholders...)
		pairArgs = append(pairArgs, taskPairArgs...)
		preparedInputs = append(preparedInputs, taskPreparedInputs...)
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
		INNER JOIN "BlockTable" b ON b.id = p.block_id::uuid AND b.deleted_at IS NULL
		INNER JOIN "BlockGroupTable" bg ON bg.id = b.block_group_id AND bg.deleted_at IS NULL
		WHERE bg.block_pack_id = p.block_pack_id::uuid
	`, strings.Join(pairPlaceholders, ","))
	if err := h.db.WithContext(ctx).Raw(sql, pairArgs...).Scan(&validRows).Error; err != nil {
		return successes, exceptions.Block.NotFound().WithOrigin(err)
	}
	validBlockByBlockPack := make(map[[2]uuid.UUID]bool, len(validRows))
	for _, validRow := range validRows {
		validBlockByBlockPack[[2]uuid.UUID{validRow.BlockId, validRow.BlockPackId}] = true
	}

	bulkInputs := make([]inputs.BulkUpdateBlockInput, 0, len(preparedInputs))
	taskIndexes := make([]int, 0, len(preparedInputs))
	for _, preparedInput := range preparedInputs {
		if !validBlockByBlockPack[[2]uuid.UUID{preparedInput.BlockId, preparedInput.BlockPackId}] {
			continue
		}
		bulkInputs = append(bulkInputs, preparedInput.Input)
		taskIndexes = append(taskIndexes, preparedInput.TaskIndex)
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

func (h BlockPackHandler) HandleResetBlockPack(
	ctx context.Context,
	tasks []schemas.RoutineTask,
	taskIdToOwnerId map[uuid.UUID]uuid.UUID,
) ([]bool, *exceptions.Exception) {
	successes := make([]bool, len(tasks))
	blockPackIds := make([]uuid.UUID, 0, len(tasks))
	ownerIdByBlockPackId := make(map[uuid.UUID]uuid.UUID, len(tasks))
	taskIndexesByBlockPackId := make(map[uuid.UUID][]int, len(tasks))

	for taskIndex, task := range tasks {
		ownerId, exists := taskIdToOwnerId[task.Id]
		if !exists {
			continue
		}

		payload, exception := decodePayload[dtos.ResetBlockPackRoutineTaskPayload](task)
		if exception != nil {
			continue
		}
		blockPackIds = append(blockPackIds, payload.BlockPackId)
		ownerIdByBlockPackId[payload.BlockPackId] = ownerId
		taskIndexesByBlockPackId[payload.BlockPackId] = append(taskIndexesByBlockPackId[payload.BlockPackId], taskIndex)
	}

	if len(blockPackIds) == 0 {
		return successes, nil
	}

	var blockGroups []struct {
		Id          uuid.UUID `gorm:"column:id"`
		BlockPackId uuid.UUID `gorm:"column:block_pack_id"`
	}
	if err := h.db.WithContext(ctx).
		Model(&schemas.BlockGroup{}).
		Select("id, block_pack_id").
		Where("block_pack_id IN ? AND deleted_at IS NULL", blockPackIds).
		Find(&blockGroups).Error; err != nil {
		return successes, exceptions.BlockGroup.NotFound().WithOrigin(err)
	}
	if len(blockGroups) == 0 {
		for _, taskIndexes := range taskIndexesByBlockPackId {
			for _, taskIndex := range taskIndexes {
				successes[taskIndex] = true
			}
		}
		return successes, nil
	}

	bulkInputs := make([]inputs.BulkDeleteBlockGroupInput, 0, len(blockGroups))
	taskIndexes := make([][]int, 0, len(blockGroups))
	for _, blockGroup := range blockGroups {
		bulkInputs = append(bulkInputs, inputs.BulkDeleteBlockGroupInput{
			UserId: ownerIdByBlockPackId[blockGroup.BlockPackId],
			Id:     blockGroup.Id,
		})
		taskIndexes = append(taskIndexes, taskIndexesByBlockPackId[blockGroup.BlockPackId])
	}

	bulkSuccesses, exception := h.blockGroupRepository.BulkDeleteMany(
		bulkInputs,
		options.WithDB(h.db.WithContext(ctx)),
		options.WithLockingStrength(options.LockingStrengthNoKeyUpdate),
		options.WithOnlyDeleted(types.Ternary_Negative),
	)
	if exception != nil {
		return successes, exception
	}

	for _, indexes := range taskIndexesByBlockPackId {
		for _, taskIndex := range indexes {
			successes[taskIndex] = true
		}
	}
	for index, success := range bulkSuccesses {
		if !success {
			for _, taskIndex := range taskIndexes[index] {
				successes[taskIndex] = false
			}
		}
	}

	return successes, nil
}
