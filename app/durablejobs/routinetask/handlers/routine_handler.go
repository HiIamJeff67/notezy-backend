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

type RoutineHandler struct {
	db                 *gorm.DB
	namePatternMatcher matchers.NamePatternMatcherInterface
	routineRepository  repositories.RoutineRepositoryInterface
}

func NewRoutineHandler(db *gorm.DB, routineRepository repositories.RoutineRepositoryInterface) RoutineHandler {
	return RoutineHandler{
		db:                 db,
		namePatternMatcher: matchers.NewNamePatternMatcher(),
		routineRepository:  routineRepository,
	}
}

func (h RoutineHandler) HandleCreateRoutine(
	ctx context.Context,
	tasks []schemas.RoutineTask,
	taskIdToOwnerId map[uuid.UUID]uuid.UUID,
) ([]bool, *exceptions.Exception) {
	successes := make([]bool, len(tasks))
	bulkInputs := make([]inputs.BulkCreateRoutineInput, 0, len(tasks))
	taskIndexes := make([]int, 0, len(tasks))

	for taskIndex, task := range tasks {
		ownerId, exists := taskIdToOwnerId[task.Id]
		if !exists {
			continue
		}

		payload, exception := decodePayload[dtos.CreateRoutineRoutineTaskPayload](task)
		if exception != nil {
			continue
		}
		title, exception := h.namePatternMatcher.Match(payload.Title, payload.TitlePattern, task)
		if exception != nil {
			continue
		}
		bulkInputs = append(bulkInputs, inputs.BulkCreateRoutineInput{
			UserId:           ownerId,
			Id:               payload.Id,
			StationId:        payload.StationId,
			Title:            title,
			Description:      payload.Description,
			Status:           payload.Status,
			IsPinned:         payload.IsPinned,
			ScheduledStartAt: payload.ScheduledStartAt,
			ScheduledEndAt:   payload.ScheduledEndAt,
			Period:           payload.Period,
			Timezone:         payload.Timezone,
		})
		taskIndexes = append(taskIndexes, taskIndex)
	}

	if len(bulkInputs) == 0 {
		return successes, nil
	}
	bulkSuccesses, exception := h.routineRepository.BulkCreateMany(
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

func (h RoutineHandler) HandleUpdateRoutine(
	ctx context.Context,
	tasks []schemas.RoutineTask,
	taskIdToOwnerId map[uuid.UUID]uuid.UUID,
) ([]bool, *exceptions.Exception) {
	successes := make([]bool, len(tasks))
	bulkInputs := make([]inputs.BulkUpdateRoutineInput, 0, len(tasks))
	taskIndexes := make([]int, 0, len(tasks))

	for taskIndex, task := range tasks {
		ownerId, exists := taskIdToOwnerId[task.Id]
		if !exists {
			continue
		}

		payload, exception := decodePayload[dtos.UpdateRoutineRoutineTaskPayload](task)
		if exception != nil {
			continue
		}
		title := payload.Title
		if title != nil && payload.TitlePattern != nil {
			matchedTitle, exception := h.namePatternMatcher.Match(*title, *payload.TitlePattern, task)
			if exception != nil {
				continue
			}
			title = &matchedTitle
		}
		bulkInputs = append(bulkInputs, inputs.BulkUpdateRoutineInput{
			UserId: ownerId,
			Id:     payload.RoutineId,
			PartialUpdateInput: inputs.PartialUpdateRoutineInput{
				Values: inputs.UpdateRoutineInput{
					Title:            title,
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
		taskIndexes = append(taskIndexes, taskIndex)
	}

	if len(bulkInputs) == 0 {
		return successes, nil
	}
	bulkSuccesses, exception := h.routineRepository.BulkUpdateMany(
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
