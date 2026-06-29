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
) map[uuid.UUID]*exceptions.Exception {
	results := make(map[uuid.UUID]*exceptions.Exception)
	ownerIdToTasks := make(map[uuid.UUID][]schemas.RoutineTask)
	for _, task := range tasks {
		ownerId, exists := taskIdToOwnerId[task.Id]
		if !exists {
			results[task.Id] = exceptions.Station.NoPermission("run this routine task")
			continue
		}
		ownerIdToTasks[ownerId] = append(ownerIdToTasks[ownerId], task)
	}

	for ownerId, createTasks := range ownerIdToTasks {
		blockPackInputs := make([]inputs.BulkCreateBlockPackInput, 0, len(createTasks))
		blockGroupInputs := make([]inputs.BulkCreateBlockGroupInput, 0)
		blockGroupSizeUpdates := make([]inputs.BulkUpdateBlockGroupsInput, 0)
		blockContents := make([]inputs.CreateBlockGroupContentInput, 0)

		for _, task := range createTasks {
			payload, exception := decodePayload[dtos.CreateBlockPackRoutineTaskPayload](task)
			if exception != nil {
				results[task.Id] = exception
				continue
			}

			blockPackId := uuid.New()
			taskBlockGroupInputs := make([]inputs.BulkCreateBlockGroupInput, 0, len(payload.Template.BlockGroups))
			taskBlockGroupSizeUpdates := make([]inputs.BulkUpdateBlockGroupsInput, 0, len(payload.Template.BlockGroups))
			taskBlockContents := make([]inputs.CreateBlockGroupContentInput, 0, len(payload.Template.BlockGroups))
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
					results[task.Id] = exception
					isTaskValid = false
					continue
				}

				taskBlockGroupInputs = append(taskBlockGroupInputs, inputs.BulkCreateBlockGroupInput{
					BlockPackId:      blockPackId,
					BlockGroupId:     &blockGroupId,
					PrevBlockGroupId: prevBlockGroupId,
				})
				taskBlockGroupSizeUpdates = append(taskBlockGroupSizeUpdates, inputs.BulkUpdateBlockGroupsInput{
					Id: blockGroupId,
					PartialUpdateInput: inputs.PartialUpdateBlockGroupInput{
						Values: inputs.UpdateBlockGroupInput{Size: &totalSize},
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
				taskBlockContents = append(taskBlockContents, inputs.CreateBlockGroupContentInput{
					BlockGroupId: blockGroupId,
					Blocks:       blockInputs,
				})
			}
			if !isTaskValid {
				continue
			}

			blockPackInputs = append(blockPackInputs, inputs.BulkCreateBlockPackInput{
				Id:                  &blockPackId,
				ParentSubShelfId:    payload.TargetSubShelfId,
				Name:                payload.Template.Name,
				Icon:                payload.Template.Icon,
				HeaderBackgroundURL: payload.Template.HeaderBackgroundURL,
			})
			blockGroupInputs = append(blockGroupInputs, taskBlockGroupInputs...)
			blockGroupSizeUpdates = append(blockGroupSizeUpdates, taskBlockGroupSizeUpdates...)
			blockContents = append(blockContents, taskBlockContents...)
		}

		if len(blockPackInputs) == 0 {
			continue
		}

		tx := h.db.WithContext(ctx).Begin()
		if _, exception := h.blockPackRepository.BulkCreateManyBySubShelfIds(
			ownerId,
			blockPackInputs,
			options.WithTransactionDB(tx),
			options.WithLockingStrength(options.LockingStrengthNoKeyUpdate),
			options.WithOnlyDeleted(types.Ternary_Negative),
		); exception != nil {
			tx.Rollback()
			for _, task := range createTasks {
				if results[task.Id] == nil {
					results[task.Id] = exception
				}
			}
			continue
		}
		if len(blockGroupInputs) > 0 {
			if _, exception := h.blockGroupRepository.InsertManyByBlockPackIds(
				ownerId,
				blockGroupInputs,
				options.WithTransactionDB(tx),
				options.WithLockingStrength(options.LockingStrengthNoKeyUpdate),
				options.WithOnlyDeleted(types.Ternary_Negative),
			); exception != nil {
				tx.Rollback()
				for _, task := range createTasks {
					if results[task.Id] == nil {
						results[task.Id] = exception
					}
				}
				continue
			}
		}
		if len(blockGroupSizeUpdates) > 0 {
			if exception := h.blockGroupRepository.BulkUpdateManyByIds(
				ownerId,
				blockGroupSizeUpdates,
				options.WithTransactionDB(tx),
				options.WithLockingStrength(options.LockingStrengthNoKeyUpdate),
				options.WithOnlyDeleted(types.Ternary_Negative),
			); exception != nil {
				tx.Rollback()
				for _, task := range createTasks {
					if results[task.Id] == nil {
						results[task.Id] = exception
					}
				}
				continue
			}
		}
		if len(blockContents) > 0 {
			if _, exception := h.blockRepository.CreateManyByBlockGroupIds(
				ownerId,
				blockContents,
				options.WithTransactionDB(tx),
				options.WithLockingStrength(options.LockingStrengthNoKeyUpdate),
				options.WithBatchSize(constants.MaxBatchCreateBlockSize),
				options.WithOnlyDeleted(types.Ternary_Negative),
			); exception != nil {
				tx.Rollback()
				for _, task := range createTasks {
					if results[task.Id] == nil {
						results[task.Id] = exception
					}
				}
				continue
			}
		}
		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			for _, task := range createTasks {
				if results[task.Id] == nil {
					results[task.Id] = exceptions.BlockPack.FailedToCommitTransaction().WithOrigin(err)
				}
			}
		}
	}

	return results
}

func (h BlockPackHandler) HandleUpdateBlockPack(
	ctx context.Context,
	tasks []schemas.RoutineTask,
	taskIdToOwnerId map[uuid.UUID]uuid.UUID,
) map[uuid.UUID]*exceptions.Exception {
	results := make(map[uuid.UUID]*exceptions.Exception)
	ownerIdToTasks := make(map[uuid.UUID][]schemas.RoutineTask)
	for _, task := range tasks {
		ownerId, exists := taskIdToOwnerId[task.Id]
		if !exists {
			results[task.Id] = exceptions.Station.NoPermission("run this routine task")
			continue
		}
		ownerIdToTasks[ownerId] = append(ownerIdToTasks[ownerId], task)
	}

	for ownerId, updateTasks := range ownerIdToTasks {
		blockIds := make([]uuid.UUID, 0)
		updateInputs := make([]inputs.BulkUpdateBlocksInput, 0)
		pairPlaceholders := make([]string, 0)
		pairArgs := make([]any, 0)

		for _, task := range updateTasks {
			payload, exception := decodePayload[dtos.UpdateBlockPackRoutineTaskPayload](task)
			if exception != nil {
				results[task.Id] = exception
				continue
			}

			taskBlockIds := make([]uuid.UUID, 0, len(payload.UpdatedBlocks))
			taskUpdateInputs := make([]inputs.BulkUpdateBlocksInput, 0, len(payload.UpdatedBlocks))
			taskPairPlaceholders := make([]string, 0, len(payload.UpdatedBlocks))
			taskPairArgs := make([]any, 0, len(payload.UpdatedBlocks)*2)
			isTaskValid := true

			for _, block := range payload.UpdatedBlocks {
				blockGroupId := uuid.New()
				flattenedBlocks, _, _, exception := flattenArborizedBlock(
					h.editableBlockAdapter,
					blockGroupId,
					block.ArborizedEditableBlock,
				)
				if exception != nil {
					results[task.Id] = exception
					isTaskValid = false
					continue
				}
				if len(flattenedBlocks) != 1 {
					results[task.Id] = exceptions.RoutineTask.InvalidDto().
						WithDetails("UpdateBlockPack updatedBlocks must not contain children")
					isTaskValid = false
					continue
				}

				blockType := flattenedBlocks[0].Type
				props := datatypes.JSON(flattenedBlocks[0].Props)
				content := datatypes.JSON(flattenedBlocks[0].Content)
				taskBlockIds = append(taskBlockIds, block.BlockId)
				taskPairPlaceholders = append(taskPairPlaceholders, "(?::uuid, ?::uuid)")
				taskPairArgs = append(taskPairArgs, block.BlockId, payload.BlockPackId)
				taskUpdateInputs = append(taskUpdateInputs, inputs.BulkUpdateBlocksInput{
					Id: block.BlockId,
					PartialUpdateInput: inputs.PartialUpdateBlockInput{
						Values: inputs.UpdateBlockInput{
							Type:    &blockType,
							Props:   &props,
							Content: &content,
						},
					},
				})
			}
			if !isTaskValid {
				continue
			}

			blockIds = append(blockIds, taskBlockIds...)
			pairPlaceholders = append(pairPlaceholders, taskPairPlaceholders...)
			pairArgs = append(pairArgs, taskPairArgs...)
			updateInputs = append(updateInputs, taskUpdateInputs...)
		}

		if len(updateInputs) == 0 {
			continue
		}

		tx := h.db.WithContext(ctx).Begin()
		var validCount int64
		sql := fmt.Sprintf(`
			WITH pairs(block_id, block_pack_id) AS (VALUES %s)
			SELECT COUNT(*)
			FROM pairs p
			INNER JOIN "BlockTable" b ON b.id = p.block_id::uuid AND b.deleted_at IS NULL
			INNER JOIN "BlockGroupTable" bg ON bg.id = b.block_group_id AND bg.deleted_at IS NULL
			WHERE bg.block_pack_id = p.block_pack_id::uuid
		`, strings.Join(pairPlaceholders, ","))
		if err := tx.Raw(sql, pairArgs...).Scan(&validCount).Error; err != nil {
			tx.Rollback()
			for _, task := range updateTasks {
				if results[task.Id] == nil {
					results[task.Id] = exceptions.Block.NotFound().WithOrigin(err)
				}
			}
			continue
		}
		if validCount != int64(len(blockIds)) {
			tx.Rollback()
			exception := exceptions.Block.NoPermission("update blocks outside this block pack")
			for _, task := range updateTasks {
				if results[task.Id] == nil {
					results[task.Id] = exception
				}
			}
			continue
		}

		if exception := h.blockRepository.BulkUpdateManyByIds(
			ownerId,
			updateInputs,
			options.WithTransactionDB(tx),
			options.WithLockingStrength(options.LockingStrengthNoKeyUpdate),
			options.WithOnlyDeleted(types.Ternary_Negative),
		); exception != nil {
			tx.Rollback()
			for _, task := range updateTasks {
				if results[task.Id] == nil {
					results[task.Id] = exception
				}
			}
			continue
		}
		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			for _, task := range updateTasks {
				if results[task.Id] == nil {
					results[task.Id] = exceptions.BlockPack.FailedToCommitTransaction().WithOrigin(err)
				}
			}
		}
	}

	return results
}

func (h BlockPackHandler) HandleResetBlockPack(
	ctx context.Context,
	tasks []schemas.RoutineTask,
	taskIdToOwnerId map[uuid.UUID]uuid.UUID,
) map[uuid.UUID]*exceptions.Exception {
	results := make(map[uuid.UUID]*exceptions.Exception)
	blockPackIdsByOwnerId := make(map[uuid.UUID][]uuid.UUID)
	ownerIdToTasks := make(map[uuid.UUID][]schemas.RoutineTask)

	for _, task := range tasks {
		ownerId, exists := taskIdToOwnerId[task.Id]
		if !exists {
			results[task.Id] = exceptions.Station.NoPermission("run this routine task")
			continue
		}

		payload, exception := decodePayload[dtos.ResetBlockPackRoutineTaskPayload](task)
		if exception != nil {
			results[task.Id] = exception
			continue
		}
		blockPackIdsByOwnerId[ownerId] = append(blockPackIdsByOwnerId[ownerId], payload.BlockPackId)
		ownerIdToTasks[ownerId] = append(ownerIdToTasks[ownerId], task)
	}

	for ownerId, blockPackIds := range blockPackIdsByOwnerId {
		var blockGroupIds []uuid.UUID
		if err := h.db.WithContext(ctx).
			Model(&schemas.BlockGroup{}).
			Where("block_pack_id IN ? AND deleted_at IS NULL", blockPackIds).
			Pluck("id", &blockGroupIds).Error; err != nil {
			for _, task := range ownerIdToTasks[ownerId] {
				results[task.Id] = exceptions.BlockGroup.NotFound().WithOrigin(err)
			}
			continue
		}
		if len(blockGroupIds) == 0 {
			continue
		}
		if exception := h.blockGroupRepository.SoftDeleteManyByIds(
			blockGroupIds,
			ownerId,
			options.WithDB(h.db.WithContext(ctx)),
			options.WithLockingStrength(options.LockingStrengthNoKeyUpdate),
			options.WithOnlyDeleted(types.Ternary_Negative),
		); exception != nil {
			for _, task := range ownerIdToTasks[ownerId] {
				results[task.Id] = exception
			}
		}
	}

	return results
}
