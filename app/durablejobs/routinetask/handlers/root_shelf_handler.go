package handlers

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	dtos "github.com/HiIamJeff67/notezy-backend/app/dtos"
	matchers "github.com/HiIamJeff67/notezy-backend/app/durablejobs/routinetask/handlers/matchers"
	exceptions "github.com/HiIamJeff67/notezy-backend/app/exceptions"
	inputs "github.com/HiIamJeff67/notezy-backend/app/models/inputs"
	repositories "github.com/HiIamJeff67/notezy-backend/app/models/repositories"
	schemas "github.com/HiIamJeff67/notezy-backend/app/models/schemas"
	options "github.com/HiIamJeff67/notezy-backend/app/options"
	types "github.com/HiIamJeff67/notezy-backend/shared/types"
)

type RootShelfHandler struct {
	db                  *gorm.DB
	namePatternMatcher  matchers.NamePatternMatcherInterface
	rootShelfRepository repositories.RootShelfRepositoryInterface
	subShelfRepository  repositories.SubShelfRepositoryInterface
}

func NewRootShelfHandler(
	db *gorm.DB,
	rootShelfRepository repositories.RootShelfRepositoryInterface,
	subShelfRepository repositories.SubShelfRepositoryInterface,
) RootShelfHandler {
	return RootShelfHandler{
		db:                  db,
		namePatternMatcher:  matchers.NewNamePatternMatcher(),
		rootShelfRepository: rootShelfRepository,
		subShelfRepository:  subShelfRepository,
	}
}

func (h RootShelfHandler) HandleCreateRootShelf(
	ctx context.Context,
	tasks []schemas.RoutineTask,
	taskIdToOwnerId map[uuid.UUID]uuid.UUID,
) ([]bool, *exceptions.Exception) {
	successes := make([]bool, len(tasks))
	bulkInputs := make([]inputs.BulkCreateRootShelfInput, 0, len(tasks))
	taskIndexes := make([]int, 0, len(tasks))

	for taskIndex, task := range tasks {
		ownerId, exists := taskIdToOwnerId[task.Id]
		if !exists {
			continue
		}

		payload, exception := decodePayload[dtos.CreateRootShelfRoutineTaskPayload](task)
		if exception != nil {
			continue
		}
		name, exception := h.namePatternMatcher.Match(payload.Name, payload.NamePattern, task)
		if exception != nil {
			continue
		}
		bulkInputs = append(bulkInputs, inputs.BulkCreateRootShelfInput{
			UserId: ownerId,
			Id:     payload.Id,
			Name:   name,
		})
		taskIndexes = append(taskIndexes, taskIndex)
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
	bulkInputs := make([]inputs.BulkUpdateRootShelfInput, 0, len(tasks))
	taskIndexes := make([]int, 0, len(tasks))

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
		name := *payload.Name
		if payload.NamePattern != nil {
			matchedName, exception := h.namePatternMatcher.Match(name, *payload.NamePattern, task)
			if exception != nil {
				continue
			}
			name = matchedName
		}
		bulkInputs = append(bulkInputs, inputs.BulkUpdateRootShelfInput{
			UserId: ownerId,
			Id:     payload.RootShelfId,
			PartialUpdateInput: inputs.PartialUpdateRootShelfInput{
				Values: inputs.UpdateRootShelfInput{
					Name: &name,
				},
			},
		})
		taskIndexes = append(taskIndexes, taskIndex)
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
