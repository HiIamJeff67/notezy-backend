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

	for ownerId, appendTasks := range ownerIdToTasks {
		appendInputsByBlockPackId := make(map[uuid.UUID][]inputs.CreateBlockGroupInput)
		blockGroupSizeUpdates := make([]inputs.BulkUpdateBlockGroupsInput, 0, len(appendTasks))
		blockContents := make([]inputs.CreateBlockGroupContentInput, 0, len(appendTasks))
		taskIdsByBlockPackId := make(map[uuid.UUID][]uuid.UUID)

		for _, task := range appendTasks {
			payload, exception := decodePayload[dtos.AppendBlockRoutineTaskPayload](task)
			if exception != nil {
				results[task.Id] = exception
				continue
			}

			blockGroupId := uuid.New()
			blocks, _, totalSize, exception := flattenArborizedBlock(
				h.editableBlockAdapter,
				blockGroupId,
				&payload.ArborizedEditableBlock,
			)
			if exception != nil {
				results[task.Id] = exception
				continue
			}

			appendInputs := appendInputsByBlockPackId[payload.BlockPackId]
			var prevBlockGroupId *uuid.UUID
			if len(appendInputs) > 0 {
				prev := *appendInputs[len(appendInputs)-1].BlockGroupId
				prevBlockGroupId = &prev
			}
			appendInputsByBlockPackId[payload.BlockPackId] = append(appendInputs, inputs.CreateBlockGroupInput{
				BlockGroupId:     &blockGroupId,
				PrevBlockGroupId: prevBlockGroupId,
			})
			taskIdsByBlockPackId[payload.BlockPackId] = append(taskIdsByBlockPackId[payload.BlockPackId], task.Id)
			blockGroupSizeUpdates = append(blockGroupSizeUpdates, inputs.BulkUpdateBlockGroupsInput{
				Id: blockGroupId,
				PartialUpdateInput: inputs.PartialUpdateBlockGroupInput{
					Values: inputs.UpdateBlockGroupInput{
						Size: &totalSize,
					},
				},
			})

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
			blockContents = append(blockContents, inputs.CreateBlockGroupContentInput{
				BlockGroupId: blockGroupId,
				Blocks:       blockInputs,
			})
		}

		tx := h.db.WithContext(ctx).Begin()
		hasAppendFailure := false
		for blockPackId, appendInputs := range appendInputsByBlockPackId {
			if len(appendInputs) == 0 {
				continue
			}
			if _, exception := h.blockGroupRepository.AppendManyByBlockPackId(
				blockPackId,
				ownerId,
				appendInputs,
				options.WithTransactionDB(tx),
				options.WithLockingStrength(options.LockingStrengthNoKeyUpdate),
				options.WithOnlyDeleted(types.Ternary_Negative),
			); exception != nil {
				hasAppendFailure = true
				for _, taskId := range taskIdsByBlockPackId[blockPackId] {
					results[taskId] = exception
				}
			}
		}
		if hasAppendFailure {
			tx.Rollback()
			continue
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
				for _, task := range appendTasks {
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
				for _, task := range appendTasks {
					if results[task.Id] == nil {
						results[task.Id] = exception
					}
				}
				continue
			}
		}
		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			for _, task := range appendTasks {
				if results[task.Id] == nil {
					results[task.Id] = exceptions.Block.FailedToCommitTransaction().WithOrigin(err)
				}
			}
		}
	}

	return results
}

func (h BlockHandler) HandleUpdateBlock(
	ctx context.Context,
	tasks []schemas.RoutineTask,
	taskIdToOwnerId map[uuid.UUID]uuid.UUID,
) map[uuid.UUID]*exceptions.Exception {
	results := make(map[uuid.UUID]*exceptions.Exception)
	inputsByOwnerId := make(map[uuid.UUID][]inputs.BulkUpdateBlocksInput)

	for _, task := range tasks {
		ownerId, exists := taskIdToOwnerId[task.Id]
		if !exists {
			results[task.Id] = exceptions.Station.NoPermission("run this routine task")
			continue
		}

		payload, exception := decodePayload[dtos.UpdateBlockRoutineTaskPayload](task)
		if exception != nil {
			results[task.Id] = exception
			continue
		}

		blockGroupId := uuid.New()
		blocks, _, _, exception := flattenArborizedBlock(
			h.editableBlockAdapter,
			blockGroupId,
			payload.ArborizedEditableBlock,
		)
		if exception != nil {
			results[task.Id] = exception
			continue
		}
		if len(blocks) != 1 {
			results[task.Id] = exceptions.RoutineTask.InvalidDto().WithDetails("UpdateBlock must not contain children")
			continue
		}

		blockType := blocks[0].Type
		props := datatypes.JSON(blocks[0].Props)
		content := datatypes.JSON(blocks[0].Content)
		inputsByOwnerId[ownerId] = append(inputsByOwnerId[ownerId], inputs.BulkUpdateBlocksInput{
			Id: payload.BlockId,
			PartialUpdateInput: inputs.PartialUpdateBlockInput{
				Values: inputs.UpdateBlockInput{
					Type:    &blockType,
					Props:   &props,
					Content: &content,
				},
			},
		})
	}

	for ownerId, updateInputs := range inputsByOwnerId {
		if exception := h.blockRepository.BulkUpdateManyByIds(
			ownerId,
			updateInputs,
			options.WithDB(h.db.WithContext(ctx)),
			options.WithLockingStrength(options.LockingStrengthNoKeyUpdate),
			options.WithOnlyDeleted(types.Ternary_Negative),
		); exception != nil {
			for _, task := range tasks {
				if taskIdToOwnerId[task.Id] == ownerId {
					results[task.Id] = exception
				}
			}
		}
	}

	return results
}

func (h BlockHandler) HandleResetBlock(
	ctx context.Context,
	tasks []schemas.RoutineTask,
	taskIdToOwnerId map[uuid.UUID]uuid.UUID,
) map[uuid.UUID]*exceptions.Exception {
	results := make(map[uuid.UUID]*exceptions.Exception)
	inputsByOwnerId := make(map[uuid.UUID][]inputs.BulkUpdateBlocksInput)

	for _, task := range tasks {
		ownerId, exists := taskIdToOwnerId[task.Id]
		if !exists {
			results[task.Id] = exceptions.Station.NoPermission("run this routine task")
			continue
		}

		payload, exception := decodePayload[dtos.ResetBlockRoutineTaskPayload](task)
		if exception != nil {
			results[task.Id] = exception
			continue
		}
		blockType := enums.BlockType_Paragraph
		props := datatypes.JSON([]byte("{}"))
		content := datatypes.JSON([]byte("[]"))
		inputsByOwnerId[ownerId] = append(inputsByOwnerId[ownerId], inputs.BulkUpdateBlocksInput{
			Id: payload.BlockId,
			PartialUpdateInput: inputs.PartialUpdateBlockInput{
				Values: inputs.UpdateBlockInput{
					Type:    &blockType,
					Props:   &props,
					Content: &content,
				},
			},
		})
	}

	for ownerId, resetInputs := range inputsByOwnerId {
		if exception := h.blockRepository.BulkUpdateManyByIds(
			ownerId,
			resetInputs,
			options.WithDB(h.db.WithContext(ctx)),
			options.WithLockingStrength(options.LockingStrengthNoKeyUpdate),
			options.WithOnlyDeleted(types.Ternary_Negative),
		); exception != nil {
			for _, task := range tasks {
				if taskIdToOwnerId[task.Id] == ownerId {
					results[task.Id] = exception
				}
			}
		}
	}

	return results
}
