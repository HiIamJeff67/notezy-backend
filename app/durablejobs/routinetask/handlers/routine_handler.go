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

type RoutineHandler struct {
	db                *gorm.DB
	routineRepository repositories.RoutineRepositoryInterface
}

func NewRoutineHandler(db *gorm.DB, routineRepository repositories.RoutineRepositoryInterface) RoutineHandler {
	return RoutineHandler{
		db:                db,
		routineRepository: routineRepository,
	}
}

func (h RoutineHandler) HandleCreateRoutine(
	ctx context.Context,
	tasks []schemas.RoutineTask,
	taskIdToOwnerId map[uuid.UUID]uuid.UUID,
) map[uuid.UUID]*exceptions.Exception {
	results := make(map[uuid.UUID]*exceptions.Exception)
	inputsByOwnerId := make(map[uuid.UUID][]inputs.BulkCreateRoutineInput)

	for _, task := range tasks {
		ownerId, exists := taskIdToOwnerId[task.Id]
		if !exists {
			results[task.Id] = exceptions.Station.NoPermission("run this routine task")
			continue
		}

		payload, exception := decodePayload[dtos.CreateRoutineRoutineTaskPayload](task)
		if exception != nil {
			results[task.Id] = exception
			continue
		}
		inputsByOwnerId[ownerId] = append(inputsByOwnerId[ownerId], inputs.BulkCreateRoutineInput{
			Id:               payload.Id,
			StationId:        payload.StationId,
			Title:            payload.Title,
			Description:      payload.Description,
			Status:           payload.Status,
			IsPinned:         payload.IsPinned,
			ScheduledStartAt: payload.ScheduledStartAt,
			ScheduledEndAt:   payload.ScheduledEndAt,
			Period:           payload.Period,
			Timezone:         payload.Timezone,
		})
	}

	for ownerId, createInputs := range inputsByOwnerId {
		if _, exception := h.routineRepository.BulkCreateManyByStationIds(
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

func (h RoutineHandler) HandleUpdateRoutine(
	ctx context.Context,
	tasks []schemas.RoutineTask,
	taskIdToOwnerId map[uuid.UUID]uuid.UUID,
) map[uuid.UUID]*exceptions.Exception {
	results := make(map[uuid.UUID]*exceptions.Exception)
	inputsByOwnerId := make(map[uuid.UUID][]inputs.BulkUpdateRoutineInput)

	for _, task := range tasks {
		ownerId, exists := taskIdToOwnerId[task.Id]
		if !exists {
			results[task.Id] = exceptions.Station.NoPermission("run this routine task")
			continue
		}

		payload, exception := decodePayload[dtos.UpdateRoutineRoutineTaskPayload](task)
		if exception != nil {
			results[task.Id] = exception
			continue
		}
		inputsByOwnerId[ownerId] = append(inputsByOwnerId[ownerId], inputs.BulkUpdateRoutineInput{
			Id: payload.RoutineId,
			PartialUpdateInput: inputs.PartialUpdateRoutineInput{
				Values: inputs.UpdateRoutineInput{
					Title:            payload.Title,
					Description:      payload.Description,
					Status:           payload.Status,
					IsPinned:         payload.IsPinned,
					ScheduledStartAt: payload.ScheduledStartAt,
					ScheduledEndAt:   payload.ScheduledEndAt,
					Period:           payload.Period,
					Timezone:         payload.Timezone,
				},
			},
		})
	}

	for ownerId, updateInputs := range inputsByOwnerId {
		if exception := h.routineRepository.BulkUpdateManyByIds(
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
