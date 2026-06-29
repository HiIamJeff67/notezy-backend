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

type SubShelfHandler struct {
	db                  *gorm.DB
	subShelfRepository  repositories.SubShelfRepositoryInterface
	materialRepository  repositories.MaterialRepositoryInterface
	blockPackRepository repositories.BlockPackRepositoryInterface
}

func NewSubShelfHandler(
	db *gorm.DB,
	subShelfRepository repositories.SubShelfRepositoryInterface,
	materialRepository repositories.MaterialRepositoryInterface,
	blockPackRepository repositories.BlockPackRepositoryInterface,
) SubShelfHandler {
	return SubShelfHandler{
		db:                  db,
		subShelfRepository:  subShelfRepository,
		materialRepository:  materialRepository,
		blockPackRepository: blockPackRepository,
	}
}

func (h SubShelfHandler) HandleCreateSubShelf(
	ctx context.Context,
	tasks []schemas.RoutineTask,
	taskIdToOwnerId map[uuid.UUID]uuid.UUID,
) map[uuid.UUID]*exceptions.Exception {
	results := make(map[uuid.UUID]*exceptions.Exception)
	inputsByOwnerId := make(map[uuid.UUID][]inputs.BulkCreateSubShelfInput)

	for _, task := range tasks {
		ownerId, exists := taskIdToOwnerId[task.Id]
		if !exists {
			results[task.Id] = exceptions.Station.NoPermission("run this routine task")
			continue
		}

		payload, exception := decodePayload[dtos.CreateSubShelfRoutineTaskPayload](task)
		if exception != nil {
			results[task.Id] = exception
			continue
		}
		inputsByOwnerId[ownerId] = append(inputsByOwnerId[ownerId], inputs.BulkCreateSubShelfInput{
			Id:             payload.Id,
			RootShelfId:    payload.RootShelfId,
			PrevSubShelfId: payload.PrevSubShelfId,
			Name:           payload.Name,
		})
	}

	for ownerId, createInputs := range inputsByOwnerId {
		if _, exception := h.subShelfRepository.BulkCreateManyByRootShelfIds(
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

func (h SubShelfHandler) HandleUpdateSubShelf(
	ctx context.Context,
	tasks []schemas.RoutineTask,
	taskIdToOwnerId map[uuid.UUID]uuid.UUID,
) map[uuid.UUID]*exceptions.Exception {
	results := make(map[uuid.UUID]*exceptions.Exception)
	inputsByOwnerId := make(map[uuid.UUID][]inputs.BulkUpdateSubShelfInput)

	for _, task := range tasks {
		ownerId, exists := taskIdToOwnerId[task.Id]
		if !exists {
			results[task.Id] = exceptions.Station.NoPermission("run this routine task")
			continue
		}

		payload, exception := decodePayload[dtos.UpdateSubShelfRoutineTaskPayload](task)
		if exception != nil {
			results[task.Id] = exception
			continue
		}
		if payload.Name == nil {
			continue
		}
		inputsByOwnerId[ownerId] = append(inputsByOwnerId[ownerId], inputs.BulkUpdateSubShelfInput{
			Id: payload.SubShelfId,
			PartialUpdateInput: inputs.PartialUpdateSubShelfInput{
				Values: inputs.UpdateSubShelfInput{
					Name: payload.Name,
				},
			},
		})
	}

	for ownerId, updateInputs := range inputsByOwnerId {
		if exception := h.subShelfRepository.BulkUpdateManyByIds(
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

func (h SubShelfHandler) HandleResetSubShelf(
	ctx context.Context,
	tasks []schemas.RoutineTask,
	taskIdToOwnerId map[uuid.UUID]uuid.UUID,
) map[uuid.UUID]*exceptions.Exception {
	results := make(map[uuid.UUID]*exceptions.Exception)
	subShelfIdsByOwnerId := make(map[uuid.UUID][]uuid.UUID)
	ownerIdToTasks := make(map[uuid.UUID][]schemas.RoutineTask)

	for _, task := range tasks {
		ownerId, exists := taskIdToOwnerId[task.Id]
		if !exists {
			results[task.Id] = exceptions.Station.NoPermission("run this routine task")
			continue
		}

		payload, exception := decodePayload[dtos.ResetSubShelfRoutineTaskPayload](task)
		if exception != nil {
			results[task.Id] = exception
			continue
		}
		subShelfIdsByOwnerId[ownerId] = append(subShelfIdsByOwnerId[ownerId], payload.SubShelfId)
		ownerIdToTasks[ownerId] = append(ownerIdToTasks[ownerId], task)
	}

	for ownerId, subShelfIds := range subShelfIdsByOwnerId {
		db := h.db.WithContext(ctx)
		var childSubShelfIds []uuid.UUID
		if err := db.Model(&schemas.SubShelf{}).
			Where("prev_sub_shelf_id IN ? AND deleted_at IS NULL", subShelfIds).
			Pluck("id", &childSubShelfIds).Error; err != nil {
			for _, task := range ownerIdToTasks[ownerId] {
				results[task.Id] = exceptions.Shelf.NotFound().WithOrigin(err)
			}
			continue
		}

		var blockPackIds []uuid.UUID
		if err := db.Model(&schemas.BlockPack{}).
			Where("parent_sub_shelf_id IN ? AND deleted_at IS NULL", subShelfIds).
			Pluck("id", &blockPackIds).Error; err != nil {
			for _, task := range ownerIdToTasks[ownerId] {
				results[task.Id] = exceptions.BlockPack.NotFound().WithOrigin(err)
			}
			continue
		}

		var materialIds []uuid.UUID
		if err := db.Model(&schemas.Material{}).
			Where("parent_sub_shelf_id IN ? AND deleted_at IS NULL", subShelfIds).
			Pluck("id", &materialIds).Error; err != nil {
			for _, task := range ownerIdToTasks[ownerId] {
				results[task.Id] = exceptions.Material.NotFound().WithOrigin(err)
			}
			continue
		}

		tx := db.Begin()
		if len(childSubShelfIds) > 0 {
			if exception := h.subShelfRepository.SoftDeleteManyByIds(
				childSubShelfIds,
				ownerId,
				options.WithTransactionDB(tx),
				options.WithLockingStrength(options.LockingStrengthNoKeyUpdate),
				options.WithOnlyDeleted(types.Ternary_Negative),
			); exception != nil {
				tx.Rollback()
				for _, task := range ownerIdToTasks[ownerId] {
					results[task.Id] = exception
				}
				continue
			}
		}
		if len(blockPackIds) > 0 {
			if exception := h.blockPackRepository.SoftDeleteManyByIds(
				blockPackIds,
				ownerId,
				options.WithTransactionDB(tx),
				options.WithLockingStrength(options.LockingStrengthNoKeyUpdate),
				options.WithOnlyDeleted(types.Ternary_Negative),
			); exception != nil {
				tx.Rollback()
				for _, task := range ownerIdToTasks[ownerId] {
					results[task.Id] = exception
				}
				continue
			}
		}
		if len(materialIds) > 0 {
			if exception := h.materialRepository.SoftDeleteManyByIds(
				materialIds,
				ownerId,
				options.WithTransactionDB(tx),
				options.WithLockingStrength(options.LockingStrengthNoKeyUpdate),
				options.WithOnlyDeleted(types.Ternary_Negative),
			); exception != nil {
				tx.Rollback()
				for _, task := range ownerIdToTasks[ownerId] {
					results[task.Id] = exception
				}
				continue
			}
		}
		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			for _, task := range ownerIdToTasks[ownerId] {
				results[task.Id] = exceptions.Shelf.FailedToCommitTransaction().WithOrigin(err)
			}
		}
	}

	return results
}
