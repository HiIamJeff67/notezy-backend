package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	validator "github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	blockpacksdto "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api/block-packs"
	typescontract "github.com/HiIamJeff67/notezy-backend/contracts/types"
	exceptions "github.com/HiIamJeff67/notezy-backend/internal/exceptions"
	inputs "github.com/HiIamJeff67/notezy-backend/internal/services/durablejob/data/inputs"
	options "github.com/HiIamJeff67/notezy-backend/internal/services/durablejob/data/options"
	repositories "github.com/HiIamJeff67/notezy-backend/internal/services/durablejob/data/repositories"
	schemas "github.com/HiIamJeff67/notezy-backend/internal/services/durablejob/data/schemas"
	enums "github.com/HiIamJeff67/notezy-backend/internal/services/durablejob/data/schemas/enums"
	matchers "github.com/HiIamJeff67/notezy-backend/internal/services/durablejob/routinetask/handlers/matchers"
	resolvers "github.com/HiIamJeff67/notezy-backend/internal/services/durablejob/routinetask/handlers/resolvers"
	yjsmaintenance "github.com/HiIamJeff67/notezy-backend/internal/services/durablejob/yjsmaintenance"
	constants "github.com/HiIamJeff67/notezy-backend/shared/constants"
	types "github.com/HiIamJeff67/notezy-backend/shared/types"
)

type BlockPackHandler struct {
	validator            *validator.Validate
	db                   *gorm.DB
	patternResolver      resolvers.PatternResolverInterface
	templateBlockMatcher matchers.TemplateBlockMatcherInterface
	blockPackRepository  repositories.BlockPackRepositoryInterface
	blockRepository      repositories.BlockRepositoryInterface
	yjsWorkerClient      yjsmaintenance.WorkerClient
}

func NewBlockPackHandler(
	validator *validator.Validate,
	db *gorm.DB,
	patternResolver resolvers.PatternResolverInterface,
	templateBlockMatcher matchers.TemplateBlockMatcherInterface,
	blockPackRepository repositories.BlockPackRepositoryInterface,
	blockRepository repositories.BlockRepositoryInterface,
	yjsWorkerClient yjsmaintenance.WorkerClient,
) BlockPackHandler {
	if patternResolver == nil {
		patternResolver = resolvers.NewPatternResolver(db, blockRepository, blockPackRepository)
	}
	if templateBlockMatcher == nil {
		templateBlockMatcher = matchers.NewTemplateBlockMatcher()
	}
	return BlockPackHandler{
		validator:            validator,
		db:                   db,
		patternResolver:      patternResolver,
		templateBlockMatcher: templateBlockMatcher,
		blockPackRepository:  blockPackRepository,
		blockRepository:      blockRepository,
		yjsWorkerClient:      yjsWorkerClient,
	}
}

func (h BlockPackHandler) HandleCreateBlockPack(
	ctx context.Context,
	tasks []schemas.RoutineTask,
	taskIdToActorUserId map[uuid.UUID]uuid.UUID,
	allowedPermissions []enums.AccessControlPermission,
) ([]bool, *exceptions.Exception) {
	successes := make([]bool, len(tasks))
	candidateTaskIndexes := make([]int, 0, len(tasks))
	candidateTasks := make([]schemas.RoutineTask, 0, len(tasks))
	candidateActorUserIds := make([]uuid.UUID, 0, len(tasks))
	candidatePayloads := make([]typescontract.CreateBlockPackRoutineTaskPayload, 0, len(tasks))
	candidatePatterns := make([]typescontract.RoutineTaskPattern, 0, len(tasks))

	for taskIndex, task := range tasks {
		actorUserId, exists := taskIdToActorUserId[task.Id]
		if !exists {
			continue
		}
		payload, exception := decodePayload[typescontract.CreateBlockPackRoutineTaskPayload](h.validator, task)
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

	blockPackInputs := make([]inputs.BulkCreateBlockPackInput, 0, len(candidateTasks))
	blockContentInputs := make([]inputs.BulkCreateBlockPackContentInput, 0, len(candidateTasks))
	initializationReqDtos := make([]blockpacksdto.InitializeBlockPackYjsDocumentReqDto, 0, len(candidateTasks))
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
		matchedRootBlocks := make([]typescontract.ArborizedEditableBlock, 0, len(payload.Template.Blocks))
		prevRootInputIndex := -1
		for _, block := range payload.Template.Blocks {
			matchedBlock, exception := h.templateBlockMatcher.MatchArborizedEditableBlock(block.ArborizedEditableBlock, patternValues)
			if exception != nil {
				taskFailed = true
				break
			}
			matchedRootBlocks = append(matchedRootBlocks, matchedBlock)
			blocks, _, _, exception := flattenArborizedBlock(blockPackId, &matchedBlock)
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
			UserId:              candidateActorUserIds[candidateIndex],
			Id:                  &blockPackId,
			ParentSubShelfId:    payload.TargetSubShelfId,
			Name:                name,
			Icon:                (*enums.SupportedIcon)(payload.Template.Icon),
			HeaderBackgroundURL: payload.Template.HeaderBackgroundURL,
		})
		blockContentInputs = append(blockContentInputs, inputs.BulkCreateBlockPackContentInput{
			UserId:      candidateActorUserIds[candidateIndex],
			BlockPackId: blockPackId,
			Blocks:      taskBlocks,
		})
		initializationReqDtos = append(initializationReqDtos, blockpacksdto.InitializeBlockPackYjsDocumentReqDto{
			Blocks: matchedRootBlocks,
		})
		preparedTaskIndexes = append(preparedTaskIndexes, candidateTaskIndexes[candidateIndex])
	}
	if len(blockPackInputs) == 0 {
		return successes, nil
	}
	initializationResDtos, err := h.yjsWorkerClient.InitializeDocuments(ctx, initializationReqDtos)
	if err != nil {
		return successes, exceptions.New(
			"FailedToCreate",
			"BlockPack",
			"Create",
			"Failed to initialize block pack documents",
			http.StatusInternalServerError,
			true,
		).WithOrigin(err)
	}

	tx := h.db.WithContext(ctx).Begin()

	blockPackSuccesses, exception := h.blockPackRepository.BulkCreateMany(
		blockPackInputs,
		options.WithTransactionDB(tx),
		options.WithAllowedPermissions(allowedPermissions),
		options.WithOnlyDeleted(types.Ternary_Negative),
	)
	if exception != nil {
		tx.Rollback()
		return successes, exception
	}

	successfulBlockContentInputs := make([]inputs.BulkCreateBlockPackContentInput, 0, len(blockContentInputs))
	successfulInitializationResDtos := make([]blockpacksdto.InitializeBlockPackYjsDocumentResDto, 0, len(initializationResDtos))
	successfulTaskIndexes := make([]int, 0, len(preparedTaskIndexes))
	for index, success := range blockPackSuccesses {
		if success {
			successfulBlockContentInputs = append(successfulBlockContentInputs, blockContentInputs[index])
			successfulInitializationResDtos = append(successfulInitializationResDtos, initializationResDtos[index])
			successfulTaskIndexes = append(successfulTaskIndexes, preparedTaskIndexes[index])
		}
	}
	if len(successfulBlockContentInputs) == 0 {
		tx.Rollback()
		return successes, nil
	}

	documents := make([]schemas.BlockPackYjsDocument, len(successfulBlockContentInputs))
	for index, successfulBlockContentInput := range successfulBlockContentInputs {
		documents[index] = schemas.BlockPackYjsDocument{
			BlockPackId:            successfulBlockContentInput.BlockPackId,
			Snapshot:               successfulInitializationResDtos[index].Snapshot,
			StateVector:            successfulInitializationResDtos[index].StateVector,
			ProjectedUntilSequence: 0,
		}
	}
	if err := tx.CreateInBatches(&documents, constants.MaxBatchCreateBlockSize).Error; err != nil {
		tx.Rollback()
		return successes, exceptions.New(
			"FailedToCreate",
			"BlockPack",
			"Create",
			"Failed to create block pack documents",
			http.StatusInternalServerError,
			true,
		).WithOrigin(err)
	}

	blockSuccesses, exception := h.blockRepository.BulkCreateMany(
		successfulBlockContentInputs,
		options.WithTransactionDB(tx),
		options.WithAllowedPermissions(allowedPermissions),
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
		return successes, exceptions.New(
			"TransactionCommitFailed",
			"BlockPack",
			"Create",
			"Failed to commit the block pack creation transaction",
			http.StatusInternalServerError,
			true,
		).WithOrigin(err)
	}

	for _, taskIndex := range successfulTaskIndexes {
		successes[taskIndex] = true
	}

	return successes, nil
}

func (h BlockPackHandler) HandleUpdateBlockPack(
	ctx context.Context,
	tasks []schemas.RoutineTask,
	taskIdToActorUserId map[uuid.UUID]uuid.UUID,
	allowedPermissions []enums.AccessControlPermission,
) ([]bool, *exceptions.Exception) {
	successes := make([]bool, len(tasks))
	candidateTaskIndexes := make([]int, 0, len(tasks))
	candidateTasks := make([]schemas.RoutineTask, 0, len(tasks))
	candidateActorUserIds := make([]uuid.UUID, 0, len(tasks))
	candidatePayloads := make([]typescontract.UpdateBlockPackRoutineTaskPayload, 0, len(tasks))
	candidatePatterns := make([]typescontract.RoutineTaskPattern, 0, len(tasks))

	for taskIndex, task := range tasks {
		actorUserId, exists := taskIdToActorUserId[task.Id]
		if !exists {
			continue
		}
		payload, exception := decodePayload[typescontract.UpdateBlockPackRoutineTaskPayload](h.validator, task)
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

	preparedInputs := make([]inputs.BulkUpdateBlockInput, 0)
	taskIndexes := make([]int, 0)
	pairPlaceholders := make([]string, 0)
	pairArgs := make([]any, 0)

	for candidateIndex, payload := range candidatePayloads {
		if !patternSuccesses[candidateIndex] {
			continue
		}
		actorUserId := candidateActorUserIds[candidateIndex]
		patternValues := patternValuesByCandidate[candidateIndex]
		for _, block := range payload.UpdatedBlocks {
			if block.ArborizedEditableBlock == nil {
				continue
			}
			matchedBlock, exception := h.templateBlockMatcher.MatchArborizedEditableBlock(*block.ArborizedEditableBlock, patternValues)
			if exception != nil {
				continue
			}
			flattenedBlocks, _, _, exception := flattenArborizedBlock(payload.BlockPackId, &matchedBlock)
			if exception != nil || len(flattenedBlocks) != 1 {
				continue
			}
			blockType := flattenedBlocks[0].Type
			props := datatypes.JSON(flattenedBlocks[0].Props)
			content := datatypes.JSON(flattenedBlocks[0].Content)
			pairPlaceholders = append(pairPlaceholders, "(?::uuid, ?::uuid)")
			pairArgs = append(pairArgs, block.BlockId, payload.BlockPackId)
			preparedInputs = append(preparedInputs, inputs.BulkUpdateBlockInput{
				UserId: actorUserId,
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
		INNER JOIN "BlockTable" b ON b.id = p.block_id::uuid AND b.block_pack_id = p.block_pack_id::uuid
	`, strings.Join(pairPlaceholders, ","))
	if err := h.db.WithContext(ctx).Raw(sql, pairArgs...).Scan(&validRows).Error; err != nil {
		return successes, exceptions.New(
			"QueryFailed",
			"Block",
			"Update",
			"Failed to validate block pack blocks",
			http.StatusInternalServerError,
			true,
		).WithOrigin(err)
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

	bulkSuccesses, exception := h.blockRepository.BulkUpdateMany(
		filteredInputs,
		options.WithDB(h.db.WithContext(ctx)),
		options.WithAllowedPermissions(allowedPermissions),
		options.WithOnlyDeleted(types.Ternary_Negative),
	)
	if exception != nil {
		return successes, exception
	}
	for index, success := range bulkSuccesses {
		successes[filteredTaskIndexes[index]] = success
	}

	return successes, nil
}

func (h BlockPackHandler) HandleResetBlockPack(
	ctx context.Context,
	tasks []schemas.RoutineTask,
	taskIdToActorUserId map[uuid.UUID]uuid.UUID,
	allowedPermissions []enums.AccessControlPermission,
) ([]bool, *exceptions.Exception) {
	successes := make([]bool, len(tasks))
	checkInputs := make([]inputs.BulkCheckBlockPackPermissionInput, 0, len(tasks))
	taskIndexes := make([]int, 0, len(tasks))
	blockPackIds := make([]uuid.UUID, 0, len(tasks))

	for taskIndex, task := range tasks {
		actorUserId, exists := taskIdToActorUserId[task.Id]
		if !exists {
			continue
		}
		payload, exception := decodePayload[typescontract.ResetBlockPackRoutineTaskPayload](h.validator, task)
		if exception != nil {
			continue
		}
		checkInputs = append(checkInputs, inputs.BulkCheckBlockPackPermissionInput{
			UserId: actorUserId,
			Id:     payload.BlockPackId,
		})
		taskIndexes = append(taskIndexes, taskIndex)
		blockPackIds = append(blockPackIds, payload.BlockPackId)
	}
	if len(checkInputs) == 0 {
		return successes, nil
	}

	tx := h.db.WithContext(ctx).Begin()

	checkSuccesses, _, exception := h.blockPackRepository.BulkCheckPermissionsAndGetManyByIds(
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
		return successes, exceptions.New(
			"FailedToUpdate",
			"Block",
			"Reset",
			"Failed to reset block pack blocks",
			http.StatusInternalServerError,
			true,
		).WithOrigin(err)
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return successes, exceptions.New(
			"TransactionCommitFailed",
			"BlockPack",
			"Reset",
			"Failed to commit the block pack reset transaction",
			http.StatusInternalServerError,
			true,
		).WithOrigin(err)
	}

	for index, success := range checkSuccesses {
		successes[taskIndexes[index]] = success
	}

	return successes, nil
}
