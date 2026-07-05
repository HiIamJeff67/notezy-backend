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

type RoutineHandler struct {
	db                   *gorm.DB
	patternResolver      resolvers.PatternResolverInterface
	templateBlockMatcher matchers.TemplateBlockMatcherInterface
	routineRepository    repositories.RoutineRepositoryInterface
}

func NewRoutineHandler(
	db *gorm.DB,
	patternResolver resolvers.PatternResolverInterface,
	templateBlockMatcher matchers.TemplateBlockMatcherInterface,
	routineRepository repositories.RoutineRepositoryInterface,
) RoutineHandler {
	if patternResolver == nil {
		patternResolver = resolvers.NewPatternResolver(db, nil, nil)
	}
	if templateBlockMatcher == nil {
		templateBlockMatcher = matchers.NewTemplateBlockMatcher()
	}
	return RoutineHandler{
		db:                   db,
		patternResolver:      patternResolver,
		templateBlockMatcher: templateBlockMatcher,
		routineRepository:    routineRepository,
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
		patternValues, exception := h.patternResolver.Resolve(ctx, task, ownerId, payload.Pattern)
		if exception != nil {
			continue
		}
		title := h.templateBlockMatcher.MatchString(payload.Title, patternValues)
		description := h.templateBlockMatcher.MatchString(payload.Description, patternValues)
		bulkInputs = append(bulkInputs, inputs.BulkCreateRoutineInput{
			UserId:           ownerId,
			Id:               payload.Id,
			StationId:        payload.StationId,
			Title:            title,
			Description:      description,
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
		patternValues, exception := h.patternResolver.Resolve(ctx, task, ownerId, payload.Pattern)
		if exception != nil {
			continue
		}
		title := payload.Title
		if title != nil {
			matchedTitle := h.templateBlockMatcher.MatchString(*title, patternValues)
			title = &matchedTitle
		}
		description := payload.Description
		if description != nil {
			matchedDescription := h.templateBlockMatcher.MatchString(*description, patternValues)
			description = &matchedDescription
		}
		bulkInputs = append(bulkInputs, inputs.BulkUpdateRoutineInput{
			UserId: ownerId,
			Id:     payload.RoutineId,
			PartialUpdateInput: inputs.PartialUpdateRoutineInput{
				Values: inputs.UpdateRoutineInput{
					Title:            title,
					Description:      description,
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
