package handlers

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	dtos "github.com/HiIamJeff67/notezy-backend/app/dtos"
	matchers "github.com/HiIamJeff67/notezy-backend/app/durablejobs/routinetask/handlers/matchers"
	resolvers "github.com/HiIamJeff67/notezy-backend/app/durablejobs/routinetask/handlers/resolvers"
	exceptions "github.com/HiIamJeff67/notezy-backend/app/exceptions"
	inputs "github.com/HiIamJeff67/notezy-backend/app/models/inputs"
	repositories "github.com/HiIamJeff67/notezy-backend/app/models/repositories"
	schemas "github.com/HiIamJeff67/notezy-backend/app/models/schemas"
	options "github.com/HiIamJeff67/notezy-backend/app/options"
	types "github.com/HiIamJeff67/notezy-backend/shared/types"
)

type SubShelfHandler struct {
	db                   *gorm.DB
	patternResolver      resolvers.PatternResolverInterface
	templateBlockMatcher matchers.TemplateBlockMatcherInterface
	subShelfRepository   repositories.SubShelfRepositoryInterface
	materialRepository   repositories.MaterialRepositoryInterface
	blockPackRepository  repositories.BlockPackRepositoryInterface
}

func NewSubShelfHandler(
	db *gorm.DB,
	patternResolver resolvers.PatternResolverInterface,
	templateBlockMatcher matchers.TemplateBlockMatcherInterface,
	subShelfRepository repositories.SubShelfRepositoryInterface,
	materialRepository repositories.MaterialRepositoryInterface,
	blockPackRepository repositories.BlockPackRepositoryInterface,
) SubShelfHandler {
	if patternResolver == nil {
		patternResolver = resolvers.NewPatternResolver(db, nil, nil)
	}
	if templateBlockMatcher == nil {
		templateBlockMatcher = matchers.NewTemplateBlockMatcher()
	}
	return SubShelfHandler{
		db:                   db,
		patternResolver:      patternResolver,
		templateBlockMatcher: templateBlockMatcher,
		subShelfRepository:   subShelfRepository,
		materialRepository:   materialRepository,
		blockPackRepository:  blockPackRepository,
	}
}

func (h SubShelfHandler) HandleCreateSubShelf(
	ctx context.Context,
	tasks []schemas.RoutineTask,
	taskIdToOwnerId map[uuid.UUID]uuid.UUID,
) ([]bool, *exceptions.Exception) {
	successes := make([]bool, len(tasks))
	candidateTaskIndexes := make([]int, 0, len(tasks))
	candidateTasks := make([]schemas.RoutineTask, 0, len(tasks))
	candidateOwnerIds := make([]uuid.UUID, 0, len(tasks))
	candidatePayloads := make([]dtos.CreateSubShelfRoutineTaskPayload, 0, len(tasks))
	candidatePatterns := make([]dtos.RoutineTaskPattern, 0, len(tasks))

	for taskIndex, task := range tasks {
		ownerId, exists := taskIdToOwnerId[task.Id]
		if !exists {
			continue
		}

		payload, exception := decodePayload[dtos.CreateSubShelfRoutineTaskPayload](task)
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

	bulkInputs := make([]inputs.BulkCreateSubShelfInput, 0, len(candidateTasks))
	taskIndexes := make([]int, 0, len(candidateTasks))
	for candidateIndex, payload := range candidatePayloads {
		if !patternSuccesses[candidateIndex] {
			continue
		}
		patternValues := patternValuesByCandidate[candidateIndex]
		name := h.templateBlockMatcher.MatchString(payload.Name, patternValues)
		bulkInputs = append(bulkInputs, inputs.BulkCreateSubShelfInput{
			UserId:         candidateOwnerIds[candidateIndex],
			Id:             payload.Id,
			RootShelfId:    payload.RootShelfId,
			PrevSubShelfId: payload.PrevSubShelfId,
			Name:           name,
		})
		taskIndexes = append(taskIndexes, candidateTaskIndexes[candidateIndex])
	}
	if len(bulkInputs) == 0 {
		return successes, nil
	}

	bulkSuccesses, exception := h.subShelfRepository.BulkCreateMany(
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

func (h SubShelfHandler) HandleUpdateSubShelf(
	ctx context.Context,
	tasks []schemas.RoutineTask,
	taskIdToOwnerId map[uuid.UUID]uuid.UUID,
) ([]bool, *exceptions.Exception) {
	successes := make([]bool, len(tasks))
	candidateTaskIndexes := make([]int, 0, len(tasks))
	candidateTasks := make([]schemas.RoutineTask, 0, len(tasks))
	candidateOwnerIds := make([]uuid.UUID, 0, len(tasks))
	candidatePayloads := make([]dtos.UpdateSubShelfRoutineTaskPayload, 0, len(tasks))
	candidatePatterns := make([]dtos.RoutineTaskPattern, 0, len(tasks))

	for taskIndex, task := range tasks {
		ownerId, exists := taskIdToOwnerId[task.Id]
		if !exists {
			continue
		}

		payload, exception := decodePayload[dtos.UpdateSubShelfRoutineTaskPayload](task)
		if exception != nil {
			continue
		}
		if payload.Name == nil {
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

	bulkInputs := make([]inputs.BulkUpdateSubShelfInput, 0, len(candidateTasks))
	taskIndexes := make([]int, 0, len(candidateTasks))
	for candidateIndex, payload := range candidatePayloads {
		if !patternSuccesses[candidateIndex] || payload.Name == nil {
			continue
		}
		patternValues := patternValuesByCandidate[candidateIndex]
		name := h.templateBlockMatcher.MatchString(*payload.Name, patternValues)
		bulkInputs = append(bulkInputs, inputs.BulkUpdateSubShelfInput{
			UserId: candidateOwnerIds[candidateIndex],
			Id:     payload.SubShelfId,
			PartialUpdateInput: inputs.PartialUpdateSubShelfInput{
				Values: inputs.UpdateSubShelfInput{
					Name: &name,
				},
			},
		})
		taskIndexes = append(taskIndexes, candidateTaskIndexes[candidateIndex])
	}

	if len(bulkInputs) == 0 {
		return successes, nil
	}
	bulkSuccesses, exception := h.subShelfRepository.BulkUpdateMany(
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

func (h SubShelfHandler) HandleResetSubShelf(
	ctx context.Context,
	tasks []schemas.RoutineTask,
	taskIdToOwnerId map[uuid.UUID]uuid.UUID,
) ([]bool, *exceptions.Exception) {
	successes := make([]bool, len(tasks))
	subShelfIds := make([]uuid.UUID, 0, len(tasks))
	ownerIdBySubShelfId := make(map[uuid.UUID]uuid.UUID, len(tasks))
	taskIndexesBySubShelfId := make(map[uuid.UUID][]int, len(tasks))

	for taskIndex, task := range tasks {
		ownerId, exists := taskIdToOwnerId[task.Id]
		if !exists {
			continue
		}

		payload, exception := decodePayload[dtos.ResetSubShelfRoutineTaskPayload](task)
		if exception != nil {
			continue
		}
		subShelfIds = append(subShelfIds, payload.SubShelfId)
		ownerIdBySubShelfId[payload.SubShelfId] = ownerId
		taskIndexesBySubShelfId[payload.SubShelfId] = append(taskIndexesBySubShelfId[payload.SubShelfId], taskIndex)
	}

	if len(subShelfIds) == 0 {
		return successes, nil
	}

	tx := h.db.WithContext(ctx).Begin()

	var childSubShelves []struct {
		Id             uuid.UUID `gorm:"column:id"`
		PrevSubShelfId uuid.UUID `gorm:"column:prev_sub_shelf_id"`
	}
	if err := tx.Model(&schemas.SubShelf{}).
		Select("id, prev_sub_shelf_id").
		Where("prev_sub_shelf_id IN ? AND deleted_at IS NULL", subShelfIds).
		Find(&childSubShelves).Error; err != nil {
		return successes, exceptions.Shelf.NotFound().WithOrigin(err)
	}

	var blockPacks []struct {
		Id               uuid.UUID `gorm:"column:id"`
		ParentSubShelfId uuid.UUID `gorm:"column:parent_sub_shelf_id"`
	}
	if err := tx.Model(&schemas.BlockPack{}).
		Select("id, parent_sub_shelf_id").
		Where("parent_sub_shelf_id IN ? AND deleted_at IS NULL", subShelfIds).
		Find(&blockPacks).Error; err != nil {
		return successes, exceptions.BlockPack.NotFound().WithOrigin(err)
	}

	var materials []struct {
		Id               uuid.UUID `gorm:"column:id"`
		ParentSubShelfId uuid.UUID `gorm:"column:parent_sub_shelf_id"`
	}
	if err := tx.Model(&schemas.Material{}).
		Select("id, parent_sub_shelf_id").
		Where("parent_sub_shelf_id IN ? AND deleted_at IS NULL", subShelfIds).
		Find(&materials).Error; err != nil {
		return successes, exceptions.Material.NotFound().WithOrigin(err)
	}

	for _, taskIndexes := range taskIndexesBySubShelfId {
		for _, taskIndex := range taskIndexes {
			successes[taskIndex] = true
		}
	}

	if len(childSubShelves) > 0 {
		bulkInputs := make([]inputs.BulkDeleteSubShelfInput, 0, len(childSubShelves))
		taskIndexes := make([][]int, 0, len(childSubShelves))
		for _, childSubShelf := range childSubShelves {
			bulkInputs = append(bulkInputs, inputs.BulkDeleteSubShelfInput{
				UserId: ownerIdBySubShelfId[childSubShelf.PrevSubShelfId],
				Id:     childSubShelf.Id,
			})
			taskIndexes = append(taskIndexes, taskIndexesBySubShelfId[childSubShelf.PrevSubShelfId])
		}
		bulkSuccesses, exception := h.subShelfRepository.BulkDeleteMany(
			bulkInputs,
			options.WithTransactionDB(tx),
			options.WithLockingStrength(options.LockingStrengthNoKeyUpdate),
			options.WithOnlyDeleted(types.Ternary_Negative),
		)
		if exception != nil {
			tx.Rollback()
			return successes, exception
		}
		for index, success := range bulkSuccesses {
			if !success {
				for _, taskIndex := range taskIndexes[index] {
					successes[taskIndex] = false
				}
			}
		}
	}

	if len(blockPacks) > 0 {
		bulkInputs := make([]inputs.BulkDeleteBlockPackInput, 0, len(blockPacks))
		taskIndexes := make([][]int, 0, len(blockPacks))
		for _, blockPack := range blockPacks {
			bulkInputs = append(bulkInputs, inputs.BulkDeleteBlockPackInput{
				UserId: ownerIdBySubShelfId[blockPack.ParentSubShelfId],
				Id:     blockPack.Id,
			})
			taskIndexes = append(taskIndexes, taskIndexesBySubShelfId[blockPack.ParentSubShelfId])
		}
		bulkSuccesses, exception := h.blockPackRepository.BulkDeleteMany(
			bulkInputs,
			options.WithTransactionDB(tx),
			options.WithLockingStrength(options.LockingStrengthNoKeyUpdate),
			options.WithOnlyDeleted(types.Ternary_Negative),
		)
		if exception != nil {
			tx.Rollback()
			return successes, exception
		}
		for index, success := range bulkSuccesses {
			if !success {
				for _, taskIndex := range taskIndexes[index] {
					successes[taskIndex] = false
				}
			}
		}
	}

	if len(materials) > 0 {
		bulkInputs := make([]inputs.BulkDeleteMaterialInput, 0, len(materials))
		taskIndexes := make([][]int, 0, len(materials))
		for _, material := range materials {
			bulkInputs = append(bulkInputs, inputs.BulkDeleteMaterialInput{
				UserId: ownerIdBySubShelfId[material.ParentSubShelfId],
				Id:     material.Id,
			})
			taskIndexes = append(taskIndexes, taskIndexesBySubShelfId[material.ParentSubShelfId])
		}
		bulkSuccesses, exception := h.materialRepository.BulkDeleteMany(
			bulkInputs,
			options.WithTransactionDB(tx),
			options.WithLockingStrength(options.LockingStrengthNoKeyUpdate),
			options.WithOnlyDeleted(types.Ternary_Negative),
		)
		if exception != nil {
			tx.Rollback()
			return successes, exception
		}
		for index, success := range bulkSuccesses {
			if !success {
				for _, taskIndex := range taskIndexes[index] {
					successes[taskIndex] = false
				}
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return successes, exceptions.Shelf.FailedToCommitTransaction().WithOrigin(err)
	}

	return successes, nil
}
