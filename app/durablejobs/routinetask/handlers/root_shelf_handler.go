package handlers

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	dtos "github.com/HiIamJeff67/notezy-backend/app/dtos"
	exceptions "github.com/HiIamJeff67/notezy-backend/app/exceptions"
	inputs "github.com/HiIamJeff67/notezy-backend/app/models/inputs"
	repositories "github.com/HiIamJeff67/notezy-backend/app/models/repositories"
	schemas "github.com/HiIamJeff67/notezy-backend/app/models/schemas"
	options "github.com/HiIamJeff67/notezy-backend/app/options"
	types "github.com/HiIamJeff67/notezy-backend/shared/types"
)

type RootShelfHandler struct {
	db                  *gorm.DB
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
		rootShelfRepository: rootShelfRepository,
		subShelfRepository:  subShelfRepository,
	}
}

func (h RootShelfHandler) HandleCreateRootShelf(
	ctx context.Context,
	tasks []schemas.RoutineTask,
	taskIdToOwnerId map[uuid.UUID]uuid.UUID,
) map[uuid.UUID]*exceptions.Exception {
	results := make(map[uuid.UUID]*exceptions.Exception)
	inputsByOwnerId := make(map[uuid.UUID][]inputs.CreateRootShelfInput)

	for _, task := range tasks {
		ownerId, exists := taskIdToOwnerId[task.Id]
		if !exists {
			results[task.Id] = exceptions.Station.NoPermission("run this routine task")
			continue
		}

		payload, exception := decodePayload[dtos.CreateRootShelfRoutineTaskPayload](task)
		if exception != nil {
			results[task.Id] = exception
			continue
		}
		inputsByOwnerId[ownerId] = append(inputsByOwnerId[ownerId], inputs.CreateRootShelfInput{
			Id:   payload.Id,
			Name: payload.Name,
		})
	}

	for ownerId, createInputs := range inputsByOwnerId {
		if _, exception := h.rootShelfRepository.CreateManyByOwnerId(
			ownerId,
			createInputs,
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

func (h RootShelfHandler) HandleUpdateRootShelf(
	ctx context.Context,
	tasks []schemas.RoutineTask,
	taskIdToOwnerId map[uuid.UUID]uuid.UUID,
) map[uuid.UUID]*exceptions.Exception {
	results := make(map[uuid.UUID]*exceptions.Exception)
	inputsByOwnerId := make(map[uuid.UUID][]inputs.BulkUpdateRootShelfInput)

	for _, task := range tasks {
		ownerId, exists := taskIdToOwnerId[task.Id]
		if !exists {
			results[task.Id] = exceptions.Station.NoPermission("run this routine task")
			continue
		}

		payload, exception := decodePayload[dtos.UpdateRootShelfRoutineTaskPayload](task)
		if exception != nil {
			results[task.Id] = exception
			continue
		}
		if payload.Name == nil {
			continue
		}
		inputsByOwnerId[ownerId] = append(inputsByOwnerId[ownerId], inputs.BulkUpdateRootShelfInput{
			Id: payload.RootShelfId,
			PartialUpdateInput: inputs.PartialUpdateRootShelfInput{
				Values: inputs.UpdateRootShelfInput{
					Name: payload.Name,
				},
			},
		})
	}

	for ownerId, updateInputs := range inputsByOwnerId {
		if exception := h.rootShelfRepository.BulkUpdateManyByIds(
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

func (h RootShelfHandler) HandleResetRootShelf(
	ctx context.Context,
	tasks []schemas.RoutineTask,
	taskIdToOwnerId map[uuid.UUID]uuid.UUID,
) map[uuid.UUID]*exceptions.Exception {
	results := make(map[uuid.UUID]*exceptions.Exception)
	rootShelfIdsByOwnerId := make(map[uuid.UUID][]uuid.UUID)
	ownerIdToTasks := make(map[uuid.UUID][]schemas.RoutineTask)

	for _, task := range tasks {
		ownerId, exists := taskIdToOwnerId[task.Id]
		if !exists {
			results[task.Id] = exceptions.Station.NoPermission("run this routine task")
			continue
		}

		payload, exception := decodePayload[dtos.ResetRootShelfRoutineTaskPayload](task)
		if exception != nil {
			results[task.Id] = exception
			continue
		}
		rootShelfIdsByOwnerId[ownerId] = append(rootShelfIdsByOwnerId[ownerId], payload.RootShelfId)
		ownerIdToTasks[ownerId] = append(ownerIdToTasks[ownerId], task)
	}

	for ownerId, rootShelfIds := range rootShelfIdsByOwnerId {
		var subShelfIds []uuid.UUID
		if err := h.db.WithContext(ctx).
			Model(&schemas.SubShelf{}).
			Where("root_shelf_id IN ? AND deleted_at IS NULL", rootShelfIds).
			Pluck("id", &subShelfIds).Error; err != nil {
			for _, task := range ownerIdToTasks[ownerId] {
				results[task.Id] = exceptions.Shelf.NotFound().WithOrigin(err)
			}
			continue
		}
		if len(subShelfIds) == 0 {
			continue
		}
		if exception := h.subShelfRepository.SoftDeleteManyByIds(
			subShelfIds,
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
