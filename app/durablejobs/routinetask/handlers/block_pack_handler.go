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
	resolvers "github.com/HiIamJeff67/notezy-backend/app/durablejobs/routinetask/handlers/resolvers"
	exceptions "github.com/HiIamJeff67/notezy-backend/app/exceptions"
	inputs "github.com/HiIamJeff67/notezy-backend/app/models/inputs"
	repositories "github.com/HiIamJeff67/notezy-backend/app/models/repositories"
	schemas "github.com/HiIamJeff67/notezy-backend/app/models/schemas"
	enums "github.com/HiIamJeff67/notezy-backend/app/models/schemas/enums"
	options "github.com/HiIamJeff67/notezy-backend/app/options"
	constants "github.com/HiIamJeff67/notezy-backend/shared/constants"
	types "github.com/HiIamJeff67/notezy-backend/shared/types"
)

type BlockPackHandler struct {
	db                   *gorm.DB
	editableBlockAdapter adapters.EditableBlockAdapterInterface
	patternResolver      resolvers.PatternResolverInterface
	templateBlockMatcher matchers.TemplateBlockMatcherInterface
	blockPackRepository  repositories.BlockPackRepositoryInterface
	blockRepository      repositories.BlockRepositoryInterface
}

func NewBlockPackHandler(
	db *gorm.DB,
	editableBlockAdapter adapters.EditableBlockAdapterInterface,
	patternResolver resolvers.PatternResolverInterface,
	templateBlockMatcher matchers.TemplateBlockMatcherInterface,
	blockPackRepository repositories.BlockPackRepositoryInterface,
	blockRepository repositories.BlockRepositoryInterface,
) BlockPackHandler {
	if editableBlockAdapter == nil {
		editableBlockAdapter = adapters.NewEditableBlockAdapter()
	}
	if patternResolver == nil {
		patternResolver = resolvers.NewPatternResolver(db, blockRepository, blockPackRepository)
	}
	if templateBlockMatcher == nil {
		templateBlockMatcher = matchers.NewTemplateBlockMatcher()
	}
	return BlockPackHandler{
		db:                   db,
		editableBlockAdapter: editableBlockAdapter,
		patternResolver:      patternResolver,
		templateBlockMatcher: templateBlockMatcher,
		blockPackRepository:  blockPackRepository,
		blockRepository:      blockRepository,
	}
}

func (h BlockPackHandler) HandleCreateBlockPack(ctx context.Context, tasks []schemas.RoutineTask, taskIdToOwnerId map[uuid.UUID]uuid.UUID) ([]bool, *exceptions.Exception) {
	successes := make([]bool, len(tasks))
	candidateTaskIndexes := make([]int, 0, len(tasks))
	candidateTasks := make([]schemas.RoutineTask, 0, len(tasks))
	candidateOwnerIds := make([]uuid.UUID, 0, len(tasks))
	candidatePayloads := make([]dtos.CreateBlockPackRoutineTaskPayload, 0, len(tasks))
	candidatePatterns := make([]dtos.RoutineTaskPattern, 0, len(tasks))

	for taskIndex, task := range tasks {
		ownerId, exists := taskIdToOwnerId[task.Id]
		if !exists {
			continue
		}
		payload, exception := decodePayload[dtos.CreateBlockPackRoutineTaskPayload](task)
		if exception != nil {
			continue
		}
		candidateTaskIndexes = append(candidateTaskIndexes, taskIndex)
		candidateTasks = append(candidateTasks, task)
		candidateOwnerIds = append(candidateOwnerIds, ownerId)
		candidatePayloads = append(candidatePayloads, *payload)
		candidatePatterns = append(candidatePatterns, payload.Pattern)
	}
	if len(candidateTasks) == 0 {
		return successes, nil
	}

	patternValuesByCandidate, patternSuccesses, exception := h.patternResolver.ResolveMany(ctx, candidateTasks, candidateOwnerIds, candidatePatterns)
	if exception != nil {
		return successes, exception
	}

	blockPackInputs := make([]inputs.BulkCreateBlockPackInput, 0, len(candidateTasks))
	blockContentInputs := make([]inputs.BulkCreateBlockPackContentInput, 0, len(candidateTasks))
	preparedTaskIndexes := make([]int, 0, len(candidateTasks))

	for candidateIndex, payload := range candidatePayloads {
		if !patternSuccesses[candidateIndex] {
			continue
		}
		patternValues := patternValuesByCandidate[candidateIndex]
		blockPackId := uuid.New()
		name := h.templateBlockMatcher.MatchString(payload.Template.Name, patternValues)
		var prevRootId *uuid.UUID
		taskFailed := false
		taskBlocks := make([]inputs.CreateBlockInput, 0)
		prevRootInputIndex := -1
		for _, block := range payload.Template.Blocks {
			matchedBlock, exception := h.templateBlockMatcher.MatchArborizedEditableBlock(block.ArborizedEditableBlock, patternValues)
			if exception != nil {
				taskFailed = true
				break
			}
			blocks, _, _, exception := flattenArborizedBlock(h.editableBlockAdapter, blockPackId, &matchedBlock)
			if exception != nil || len(blocks) == 0 {
				taskFailed = true
				break
			}
			blocks[0].PrevBlockId = prevRootId
			if prevRootInputIndex >= 0 {
				nextBlockId := blocks[0].Id
				taskBlocks[prevRootInputIndex].NextBlockId = &nextBlockId
			}
			prevRootId = &blocks[0].Id
			prevRootInputIndex = len(taskBlocks)
			for _, block := range blocks {
				taskBlocks = append(taskBlocks, inputs.CreateBlockInput{
					Id:            block.Id,
					BlockPackId:   block.BlockPackId,
					ParentBlockId: block.ParentBlockId,
					PrevBlockId:   block.PrevBlockId,
					NextBlockId:   block.NextBlockId,
					Type:          block.Type,
					Props:         block.Props,
					Content:       block.Content,
				})
			}
		}
		if taskFailed || len(taskBlocks) == 0 {
			continue
		}
		blockPackInputs = append(blockPackInputs, inputs.BulkCreateBlockPackInput{
			UserId:              candidateOwnerIds[candidateIndex],
			Id:                  &blockPackId,
			ParentSubShelfId:    payload.TargetSubShelfId,
			Name:                name,
			Icon:                payload.Template.Icon,
			HeaderBackgroundURL: payload.Template.HeaderBackgroundURL,
		})
		blockContentInputs = append(blockContentInputs, inputs.BulkCreateBlockPackContentInput{
			UserId:      candidateOwnerIds[candidateIndex],
			BlockPackId: blockPackId,
			Blocks:      taskBlocks,
		})
		preparedTaskIndexes = append(preparedTaskIndexes, candidateTaskIndexes[candidateIndex])
	}
	if len(blockPackInputs) == 0 {
		return successes, nil
	}

	tx := h.db.WithContext(ctx).Begin()

	blockPackSuccesses, exception := h.blockPackRepository.BulkCreateMany(
		blockPackInputs,
		options.WithTransactionDB(tx),
		options.WithOnlyDeleted(types.Ternary_Negative),
	)
	if exception != nil {
		tx.Rollback()
		return successes, exception
	}

	successfulBlockContentInputs := make([]inputs.BulkCreateBlockPackContentInput, 0, len(blockContentInputs))
	successfulTaskIndexes := make([]int, 0, len(preparedTaskIndexes))
	for index, success := range blockPackSuccesses {
		if success {
			successfulBlockContentInputs = append(successfulBlockContentInputs, blockContentInputs[index])
			successfulTaskIndexes = append(successfulTaskIndexes, preparedTaskIndexes[index])
		}
	}
	if len(successfulBlockContentInputs) == 0 {
		tx.Rollback()
		return successes, nil
	}

	documents := make([]schemas.BlockPackYjsDocument, len(successfulBlockContentInputs))
	for index, successfulBlockContentInput := range successfulBlockContentInputs {
		documents[index] = schemas.BlockPackYjsDocument{BlockPackId: successfulBlockContentInput.BlockPackId}
	}
	if err := tx.CreateInBatches(&documents, constants.MaxBatchCreateBlockSize).Error; err != nil {
		tx.Rollback()
		return successes, exceptions.BlockPack.FailedToCreate().WithOrigin(err)
	}

	blockSuccesses, exception := h.blockRepository.BulkCreateMany(
		successfulBlockContentInputs,
		options.WithTransactionDB(tx),
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

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return successes, exceptions.BlockPack.FailedToCommitTransaction().WithOrigin(err)
	}

	for _, taskIndex := range successfulTaskIndexes {
		successes[taskIndex] = true
	}

	return successes, nil
}

func (h BlockPackHandler) HandleUpdateBlockPack(ctx context.Context, tasks []schemas.RoutineTask, taskIdToOwnerId map[uuid.UUID]uuid.UUID) ([]bool, *exceptions.Exception) {
	successes := make([]bool, len(tasks))
	candidateTaskIndexes := make([]int, 0, len(tasks))
	candidateTasks := make([]schemas.RoutineTask, 0, len(tasks))
	candidateOwnerIds := make([]uuid.UUID, 0, len(tasks))
	candidatePayloads := make([]dtos.UpdateBlockPackRoutineTaskPayload, 0, len(tasks))
	candidatePatterns := make([]dtos.RoutineTaskPattern, 0, len(tasks))

	for taskIndex, task := range tasks {
		ownerId, exists := taskIdToOwnerId[task.Id]
		if !exists {
			continue
		}
		payload, exception := decodePayload[dtos.UpdateBlockPackRoutineTaskPayload](task)
		if exception != nil {
			continue
		}
		candidateTaskIndexes = append(candidateTaskIndexes, taskIndex)
		candidateTasks = append(candidateTasks, task)
		candidateOwnerIds = append(candidateOwnerIds, ownerId)
		candidatePayloads = append(candidatePayloads, *payload)
		candidatePatterns = append(candidatePatterns, payload.Pattern)
	}
	if len(candidateTasks) == 0 {
		return successes, nil
	}

	patternValuesByCandidate, patternSuccesses, exception := h.patternResolver.ResolveMany(ctx, candidateTasks, candidateOwnerIds, candidatePatterns)
	if exception != nil {
		return successes, exception
	}

	preparedInputs := make([]inputs.BulkUpdateBlockInput, 0)
	taskIndexes := make([]int, 0)
	pairPlaceholders := make([]string, 0)
	pairArgs := make([]any, 0)

	for candidateIndex, payload := range candidatePayloads {
		if !patternSuccesses[candidateIndex] {
			continue
		}
		ownerId := candidateOwnerIds[candidateIndex]
		patternValues := patternValuesByCandidate[candidateIndex]
		for _, block := range payload.UpdatedBlocks {
			if block.ArborizedEditableBlock == nil {
				continue
			}
			matchedBlock, exception := h.templateBlockMatcher.MatchArborizedEditableBlock(*block.ArborizedEditableBlock, patternValues)
			if exception != nil {
				continue
			}
			flattenedBlocks, _, _, exception := flattenArborizedBlock(h.editableBlockAdapter, payload.BlockPackId, &matchedBlock)
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
			taskIndexes = append(taskIndexes, candidateTaskIndexes[candidateIndex])
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
	checkInputs := make([]inputs.BulkCheckBlockPackPermissionInput, 0, len(tasks))
	taskIndexes := make([]int, 0, len(tasks))
	blockPackIds := make([]uuid.UUID, 0, len(tasks))

	for taskIndex, task := range tasks {
		ownerId, exists := taskIdToOwnerId[task.Id]
		if !exists {
			continue
		}
		payload, exception := decodePayload[dtos.ResetBlockPackRoutineTaskPayload](task)
		if exception != nil {
			continue
		}
		checkInputs = append(checkInputs, inputs.BulkCheckBlockPackPermissionInput{
			UserId: ownerId,
			Id:     payload.BlockPackId,
		})
		taskIndexes = append(taskIndexes, taskIndex)
		blockPackIds = append(blockPackIds, payload.BlockPackId)
	}
	if len(checkInputs) == 0 {
		return successes, nil
	}

	allowedPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
	}

	tx := h.db.WithContext(ctx).Begin()

	checkSuccesses, _, exception := h.blockPackRepository.BulkCheckPermissionsAndGetManyByIds(
		checkInputs,
		nil,
		allowedPermissions,
		options.WithTransactionDB(tx),
		options.WithOnlyDeleted(types.Ternary_Negative),
		options.WithLockingStrength(options.LockingStrengthNoKeyUpdate),
	)
	if exception != nil {
		tx.Rollback()
		return successes, exception
	}

	validBlockPackIds := make([]uuid.UUID, 0, len(blockPackIds))
	for index, success := range checkSuccesses {
		if success {
			validBlockPackIds = append(validBlockPackIds, blockPackIds[index])
		}
	}
	if len(validBlockPackIds) == 0 {
		tx.Rollback()
		return successes, nil
	}

	if err := tx.Model(&schemas.Block{}).
		Where("block_pack_id IN ? AND deleted_at IS NULL", validBlockPackIds).
		Updates(map[string]any{"deleted_at": time.Now(), "prev_block_id": nil, "next_block_id": nil}).Error; err != nil {
		tx.Rollback()
		return successes, exceptions.Block.FailedToUpdate().WithOrigin(err)
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return successes, exceptions.BlockPack.FailedToCommitTransaction().WithOrigin(err)
	}

	for index, success := range checkSuccesses {
		successes[taskIndexes[index]] = success
	}

	return successes, nil
}
