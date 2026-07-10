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

type RootShelfHandler struct {
	db                   *gorm.DB
	patternResolver      resolvers.PatternResolverInterface
	templateBlockMatcher matchers.TemplateBlockMatcherInterface
	rootShelfRepository  repositories.RootShelfRepositoryInterface
	subShelfRepository   repositories.SubShelfRepositoryInterface
}

func NewRootShelfHandler(
	db *gorm.DB,
	patternResolver resolvers.PatternResolverInterface,
	templateBlockMatcher matchers.TemplateBlockMatcherInterface,
	rootShelfRepository repositories.RootShelfRepositoryInterface,
	subShelfRepository repositories.SubShelfRepositoryInterface,
) RootShelfHandler {
	if patternResolver == nil {
		patternResolver = resolvers.NewPatternResolver(db, nil, nil)
	}
	if templateBlockMatcher == nil {
		templateBlockMatcher = matchers.NewTemplateBlockMatcher()
	}
	return RootShelfHandler{
		db:                   db,
		patternResolver:      patternResolver,
		templateBlockMatcher: templateBlockMatcher,
		rootShelfRepository:  rootShelfRepository,
		subShelfRepository:   subShelfRepository,
	}
}

func (h RootShelfHandler) HandleCreateRootShelf(
	ctx context.Context,
	tasks []schemas.RoutineTask,
	taskIdToOwnerId map[uuid.UUID]uuid.UUID,
) ([]bool, *exceptions.Exception) {
	successes := make([]bool, len(tasks))
	candidateTaskIndexes := make([]int, 0, len(tasks))
	candidateTasks := make([]schemas.RoutineTask, 0, len(tasks))
	candidateOwnerIds := make([]uuid.UUID, 0, len(tasks))
	candidatePayloads := make([]dtos.CreateRootShelfRoutineTaskPayload, 0, len(tasks))
	candidatePatterns := make([]dtos.RoutineTaskPattern, 0, len(tasks))

	for taskIndex, task := range tasks {
		ownerId, exists := taskIdToOwnerId[task.Id]
		if !exists {
			continue
		}

		payload, exception := decodePayload[dtos.CreateRootShelfRoutineTaskPayload](task)
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

	bulkInputs := make([]inputs.BulkCreateRootShelfInput, 0, len(candidateTasks))
	taskIndexes := make([]int, 0, len(candidateTasks))
	for candidateIndex, payload := range candidatePayloads {
		if !patternSuccesses[candidateIndex] {
			continue
		}
		patternValues := patternValuesByCandidate[candidateIndex]
		name := h.templateBlockMatcher.MatchString(payload.Name, patternValues)
		bulkInputs = append(bulkInputs, inputs.BulkCreateRootShelfInput{
			UserId: candidateOwnerIds[candidateIndex],
			Id:     payload.Id,
			Name:   name,
		})
		taskIndexes = append(taskIndexes, candidateTaskIndexes[candidateIndex])
	}

	if len(bulkInputs) == 0 {
		return successes, nil
	}
	bulkSuccesses, exception := h.rootShelfRepository.BulkCreateMany(
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

func (h RootShelfHandler) HandleUpdateRootShelf(
	ctx context.Context,
	tasks []schemas.RoutineTask,
	taskIdToOwnerId map[uuid.UUID]uuid.UUID,
) ([]bool, *exceptions.Exception) {
	successes := make([]bool, len(tasks))
	candidateTaskIndexes := make([]int, 0, len(tasks))
	candidateTasks := make([]schemas.RoutineTask, 0, len(tasks))
	candidateOwnerIds := make([]uuid.UUID, 0, len(tasks))
	candidatePayloads := make([]dtos.UpdateRootShelfRoutineTaskPayload, 0, len(tasks))
	candidatePatterns := make([]dtos.RoutineTaskPattern, 0, len(tasks))

	for taskIndex, task := range tasks {
		ownerId, exists := taskIdToOwnerId[task.Id]
		if !exists {
			continue
		}

		payload, exception := decodePayload[dtos.UpdateRootShelfRoutineTaskPayload](task)
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

	bulkInputs := make([]inputs.BulkUpdateRootShelfInput, 0, len(candidateTasks))
	taskIndexes := make([]int, 0, len(candidateTasks))
	for candidateIndex, payload := range candidatePayloads {
		if !patternSuccesses[candidateIndex] || payload.Name == nil {
			continue
		}
		patternValues := patternValuesByCandidate[candidateIndex]
		name := h.templateBlockMatcher.MatchString(*payload.Name, patternValues)
		bulkInputs = append(bulkInputs, inputs.BulkUpdateRootShelfInput{
			UserId: candidateOwnerIds[candidateIndex],
			Id:     payload.RootShelfId,
			PartialUpdateInput: inputs.PartialUpdateRootShelfInput{
				Values: inputs.UpdateRootShelfInput{
					Name: &name,
				},
			},
		})
		taskIndexes = append(taskIndexes, candidateTaskIndexes[candidateIndex])
	}

	if len(bulkInputs) == 0 {
		return successes, nil
	}

	bulkSuccesses, exception := h.rootShelfRepository.BulkUpdateMany(
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

func (h RootShelfHandler) HandleResetRootShelf(
	ctx context.Context,
	tasks []schemas.RoutineTask,
	taskIdToOwnerId map[uuid.UUID]uuid.UUID,
) ([]bool, *exceptions.Exception) {
	successes := make([]bool, len(tasks))
	rootShelfIds := make([]uuid.UUID, 0, len(tasks))
	taskIndexesByRootShelfId := make(map[uuid.UUID][]int, len(tasks))
	ownerIdByRootShelfId := make(map[uuid.UUID]uuid.UUID, len(tasks))

	for taskIndex, task := range tasks {
		ownerId, exists := taskIdToOwnerId[task.Id]
		if !exists {
			continue
		}

		payload, exception := decodePayload[dtos.ResetRootShelfRoutineTaskPayload](task)
		if exception != nil {
			continue
		}
		rootShelfIds = append(rootShelfIds, payload.RootShelfId)
		taskIndexesByRootShelfId[payload.RootShelfId] = append(taskIndexesByRootShelfId[payload.RootShelfId], taskIndex)
		ownerIdByRootShelfId[payload.RootShelfId] = ownerId
	}

	if len(rootShelfIds) == 0 {
		return successes, nil
	}

	var rows []struct {
		Id          uuid.UUID `gorm:"column:id"`
		RootShelfId uuid.UUID `gorm:"column:root_shelf_id"`
	}
	if err := h.db.WithContext(ctx).
		Model(&schemas.SubShelf{}).
		Select("id, root_shelf_id").
		Where("root_shelf_id IN ? AND deleted_at IS NULL", rootShelfIds).
		Find(&rows).Error; err != nil {
		return successes, exceptions.Shelf.NotFound().WithOrigin(err)
	}

	bulkInputs := make([]inputs.BulkDeleteSubShelfInput, 0, len(rows))
	taskIndexes := make([][]int, 0, len(rows))
	for _, row := range rows {
		bulkInputs = append(bulkInputs, inputs.BulkDeleteSubShelfInput{
			UserId: ownerIdByRootShelfId[row.RootShelfId],
			Id:     row.Id,
		})
		taskIndexes = append(taskIndexes, taskIndexesByRootShelfId[row.RootShelfId])
	}
	if len(bulkInputs) == 0 {
		for _, indexes := range taskIndexesByRootShelfId {
			for _, taskIndex := range indexes {
				successes[taskIndex] = true
			}
		}
		return successes, nil
	}

	bulkSuccesses, exception := h.subShelfRepository.BulkDeleteMany(
		bulkInputs,
		options.WithDB(h.db.WithContext(ctx)),
		options.WithLockingStrength(options.LockingStrengthNoKeyUpdate),
		options.WithOnlyDeleted(types.Ternary_Negative),
	)
	if exception != nil {
		return successes, exception
	}

	for _, indexes := range taskIndexesByRootShelfId {
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
