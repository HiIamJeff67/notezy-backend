package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	adapters "github.com/HiIamJeff67/notezy-backend/internal/adapters"
	exceptions "github.com/HiIamJeff67/notezy-backend/internal/exceptions"
	inputs "github.com/HiIamJeff67/notezy-backend/internal/services/durablejob/data/inputs"
	options "github.com/HiIamJeff67/notezy-backend/internal/services/durablejob/data/options"
	repositories "github.com/HiIamJeff67/notezy-backend/internal/services/durablejob/data/repositories"
	schemas "github.com/HiIamJeff67/notezy-backend/internal/services/durablejob/data/schemas"
	enums "github.com/HiIamJeff67/notezy-backend/internal/services/durablejob/data/schemas/enums"
	matchers "github.com/HiIamJeff67/notezy-backend/internal/services/durablejob/routinetask/handlers/matchers"
	resolvers "github.com/HiIamJeff67/notezy-backend/internal/services/durablejob/routinetask/handlers/resolvers"
	payloads "github.com/HiIamJeff67/notezy-backend/internal/services/durablejob/routinetask/payloads"
	types "github.com/HiIamJeff67/notezy-backend/internal/shared/types"
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

func (h BlockHandler) HandleAppendBlock(
	ctx context.Context,
	tasks []schemas.RoutineTask,
	taskIdToActorUserId map[uuid.UUID]uuid.UUID,
	allowedPermissions []enums.AccessControlPermission,
) ([]bool, *exceptions.Exception) {
	successes := make([]bool, len(tasks))
	candidateTaskIndexes := make([]int, 0, len(tasks))
	candidateTasks := make([]schemas.RoutineTask, 0, len(tasks))
	candidateActorUserIds := make([]uuid.UUID, 0, len(tasks))
	candidatePayloads := make([]payloads.AppendBlockRoutineTaskPayload, 0, len(tasks))
	candidatePatterns := make([]payloads.RoutineTaskPattern, 0, len(tasks))

	for taskIndex, task := range tasks {
		actorUserId, exists := taskIdToActorUserId[task.Id]
		if !exists {
			continue
		}
		payload, exception := decodePayload[payloads.AppendBlockRoutineTaskPayload](task)
		if exception != nil {
			continue
		}
		candidateTaskIndexes = append(candidateTaskIndexes, taskIndex)
		candidateTasks = append(candidateTasks, task)
		candidateActorUserIds = append(candidateActorUserIds, actorUserId)
		candidatePayloads = append(candidatePayloads, *payload)
		candidatePatterns = append(candidatePatterns, payload.Pattern)
	}
	if len(candidateTasks) == 0 {
		return successes, nil
	}

	patternValuesByCandidate, patternSuccesses, exception := h.patternResolver.ResolveMany(
		ctx,
		candidateTasks,
		candidateActorUserIds,
		candidatePatterns,
		allowedPermissions,
	)
	if exception != nil {
		return successes, exception
	}

	preparedTaskIndexes := make([]int, 0, len(candidatePayloads))
	checkInputs := make([]inputs.BulkCheckBlockPackPermissionInput, 0, len(candidatePayloads))
	blocksByPreparedTask := make([][]schemas.Block, 0, len(candidatePayloads))
	for candidateIndex, payload := range candidatePayloads {
		if !patternSuccesses[candidateIndex] {
			continue
		}
		actorUserId := candidateActorUserIds[candidateIndex]
		patternValues := patternValuesByCandidate[candidateIndex]
		matchedBlock, exception := h.templateBlockMatcher.MatchArborizedEditableBlock(payload.ArborizedEditableBlock, patternValues)
		if exception != nil {
			continue
		}
		blocks, _, _, exception := flattenArborizedBlock(h.editableBlockAdapter, payload.BlockPackId, &matchedBlock)
		if exception != nil {
			continue
		}
		if len(blocks) == 0 {
			continue
		}

		preparedTaskIndexes = append(preparedTaskIndexes, candidateTaskIndexes[candidateIndex])
		checkInputs = append(checkInputs, inputs.BulkCheckBlockPackPermissionInput{
			UserId: actorUserId,
			Id:     payload.BlockPackId,
		})
		blocksByPreparedTask = append(blocksByPreparedTask, blocks)
	}
	if len(checkInputs) == 0 {
		return successes, nil
	}

	tx := h.db.WithContext(ctx).Begin()
	if err := tx.Error; err != nil {
		return successes, exceptions.New(
			"TransactionBeginFailed",
			"Block",
			"Append",
			"Failed to start the block append transaction",
			http.StatusInternalServerError,
			true,
		).WithOrigin(err)
	}
	permissionSuccesses, _, exception := h.blockPackRepository.BulkCheckPermissionsAndGetManyByIds(
		checkInputs,
		nil,
		allowedPermissions,
		options.WithTransactionDB(tx),
		options.WithAllowedPermissions(allowedPermissions),
		options.WithOnlyDeleted(types.Ternary_Negative),
		options.WithLockingStrength(options.LockingStrengthNoKeyUpdate),
	)
	if exception != nil {
		tx.Rollback()
		return successes, exception
	}

	blockPackIds := make([]uuid.UUID, 0, len(checkInputs))
	for index, checkInput := range checkInputs {
		if permissionSuccesses[index] {
			blockPackIds = append(blockPackIds, checkInput.Id)
		}
	}
	if len(blockPackIds) == 0 {
		tx.Rollback()
		return successes, nil
	}

	var tails []struct {
		Id          uuid.UUID `gorm:"column:id"`
		BlockPackId uuid.UUID `gorm:"column:block_pack_id"`
	}
	if err := tx.Model(&schemas.Block{}).
		Select("id, block_pack_id").
		Where("block_pack_id IN ? AND parent_block_id IS NULL AND next_block_id IS NULL AND deleted_at IS NULL", blockPackIds).
		Find(&tails).Error; err != nil {
		tx.Rollback()
		return successes, exceptions.New(
			"QueryFailed",
			"Block",
			"Append",
			"Failed to retrieve existing blocks",
			http.StatusInternalServerError,
			true,
		).WithOrigin(err)
	}

	tailIdByBlockPackId := make(map[uuid.UUID]uuid.UUID, len(tails))
	for _, tail := range tails {
		tailIdByBlockPackId[tail.BlockPackId] = tail.Id
	}

	blocksToCreate := make([]schemas.Block, 0, len(blocksByPreparedTask))
	linkPlaceholders := make([]string, 0, len(blocksByPreparedTask))
	linkArgs := make([]any, 0, len(blocksByPreparedTask)*2)
	for index, blocks := range blocksByPreparedTask {
		if !permissionSuccesses[index] {
			continue
		}

		blockPackId := checkInputs[index].Id
		if tailId, exists := tailIdByBlockPackId[blockPackId]; exists {
			blocks[0].PrevBlockId = &tailId
			linkPlaceholders = append(linkPlaceholders, "(?::uuid, ?::uuid)")
			linkArgs = append(linkArgs, tailId, blocks[0].Id)
		}

		for _, block := range blocks {
			blocksToCreate = append(blocksToCreate, block)
		}
		for _, block := range blocks {
			if block.ParentBlockId == nil && block.NextBlockId == nil {
				tailIdByBlockPackId[blockPackId] = block.Id
			}
		}
		successes[preparedTaskIndexes[index]] = true
	}
	if len(blocksToCreate) == 0 {
		tx.Rollback()
		return successes, nil
	}

	if err := tx.CreateInBatches(&blocksToCreate, 100).Error; err != nil {
		tx.Rollback()
		return successes, exceptions.New(
			"FailedToCreate",
			"Block",
			"Append",
			"Failed to create blocks",
			http.StatusInternalServerError,
			true,
		).WithOrigin(err)
	}
	if len(linkPlaceholders) > 0 {
		result := tx.Exec(fmt.Sprintf(`
			UPDATE "BlockTable" AS block
			SET next_block_id = value.next_block_id
			FROM (VALUES %s) AS value(id, next_block_id)
			WHERE block.id = value.id::uuid
		`, strings.Join(linkPlaceholders, ",")), linkArgs...)
		if result.Error != nil {
			tx.Rollback()
			return successes, exceptions.New(
				"FailedToUpdate",
				"Block",
				"Append",
				"Failed to link appended blocks",
				http.StatusInternalServerError,
				true,
			).WithOrigin(result.Error)
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return successes, exceptions.New(
			"TransactionCommitFailed",
			"Block",
			"Append",
			"Failed to commit the block append transaction",
			http.StatusInternalServerError,
			true,
		).WithOrigin(err)
	}

	return successes, nil
}

func (h BlockHandler) HandleUpdateBlock(
	ctx context.Context,
	tasks []schemas.RoutineTask,
	taskIdToActorUserId map[uuid.UUID]uuid.UUID,
	allowedPermissions []enums.AccessControlPermission,
) ([]bool, *exceptions.Exception) {
	successes := make([]bool, len(tasks))
	candidateTaskIndexes := make([]int, 0, len(tasks))
	candidateTasks := make([]schemas.RoutineTask, 0, len(tasks))
	candidateActorUserIds := make([]uuid.UUID, 0, len(tasks))
	candidatePayloads := make([]payloads.UpdateBlockRoutineTaskPayload, 0, len(tasks))
	candidatePatterns := make([]payloads.RoutineTaskPattern, 0, len(tasks))

	for taskIndex, task := range tasks {
		actorUserId, exists := taskIdToActorUserId[task.Id]
		if !exists {
			continue
		}
		payload, exception := decodePayload[payloads.UpdateBlockRoutineTaskPayload](task)
		if exception != nil {
			continue
		}
		if payload.ArborizedEditableBlock == nil {
			continue
		}
		candidateTaskIndexes = append(candidateTaskIndexes, taskIndex)
		candidateTasks = append(candidateTasks, task)
		candidateActorUserIds = append(candidateActorUserIds, actorUserId)
		candidatePayloads = append(candidatePayloads, *payload)
		candidatePatterns = append(candidatePatterns, payload.Pattern)
	}
	if len(candidateTasks) == 0 {
		return successes, nil
	}

	patternValuesByCandidate, patternSuccesses, exception := h.patternResolver.ResolveMany(
		ctx,
		candidateTasks,
		candidateActorUserIds,
		candidatePatterns,
		allowedPermissions,
	)
	if exception != nil {
		return successes, exception
	}

	bulkInputs := make([]inputs.BulkUpdateBlockInput, 0, len(candidateTasks))
	taskIndexes := make([]int, 0, len(candidateTasks))
	for candidateIndex, payload := range candidatePayloads {
		if !patternSuccesses[candidateIndex] || payload.ArborizedEditableBlock == nil {
			continue
		}
		patternValues := patternValuesByCandidate[candidateIndex]
		matchedBlock, exception := h.templateBlockMatcher.MatchArborizedEditableBlock(*payload.ArborizedEditableBlock, patternValues)
		if exception != nil {
			continue
		}
		rawBlocks, _, exception := h.editableBlockAdapter.FlattenToRaw(&matchedBlock)
		if exception != nil || len(rawBlocks) != 1 {
			continue
		}
		blockType := enums.BlockType(rawBlocks[0].Type)
		props := datatypes.JSON(rawBlocks[0].Props)
		content := datatypes.JSON(rawBlocks[0].Content)
		bulkInputs = append(bulkInputs, inputs.BulkUpdateBlockInput{
			UserId: candidateActorUserIds[candidateIndex],
			Id:     payload.BlockId,
			PartialUpdateInput: inputs.PartialUpdateBlockInput{Values: inputs.UpdateBlockInput{
				Type:    &blockType,
				Props:   &props,
				Content: &content,
			}},
		})
		taskIndexes = append(taskIndexes, candidateTaskIndexes[candidateIndex])
	}
	if len(bulkInputs) == 0 {
		return successes, nil
	}

	bulkSuccesses, exception := h.blockRepository.BulkUpdateMany(
		bulkInputs,
		options.WithDB(h.db.WithContext(ctx)),
		options.WithAllowedPermissions(allowedPermissions),
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
	taskIdToActorUserId map[uuid.UUID]uuid.UUID,
	allowedPermissions []enums.AccessControlPermission,
) ([]bool, *exceptions.Exception) {
	successes := make([]bool, len(tasks))
	blockType := enums.BlockType_Paragraph
	props := datatypes.JSON([]byte("{}"))
	content := datatypes.JSON([]byte("[]"))
	bulkInputs := make([]inputs.BulkUpdateBlockInput, 0, len(tasks))
	taskIndexes := make([]int, 0, len(tasks))

	for taskIndex, task := range tasks {
		actorUserId, exists := taskIdToActorUserId[task.Id]
		if !exists {
			continue
		}
		payload, exception := decodePayload[payloads.ResetBlockRoutineTaskPayload](task)
		if exception != nil {
			continue
		}
		bulkInputs = append(bulkInputs, inputs.BulkUpdateBlockInput{
			UserId: actorUserId,
			Id:     payload.BlockId,
			PartialUpdateInput: inputs.PartialUpdateBlockInput{
				Values: inputs.UpdateBlockInput{Type: &blockType, Props: &props, Content: &content},
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
		options.WithAllowedPermissions(allowedPermissions),
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
