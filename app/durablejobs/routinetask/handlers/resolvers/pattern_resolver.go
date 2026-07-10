package resolvers

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	dtos "github.com/HiIamJeff67/notezy-backend/app/dtos"
	exceptions "github.com/HiIamJeff67/notezy-backend/app/exceptions"
	repositories "github.com/HiIamJeff67/notezy-backend/app/models/repositories"
	schemas "github.com/HiIamJeff67/notezy-backend/app/models/schemas"
)

const (
	PatternSourceScheduledAt        = "scheduledAt"
	PatternSourceRecordId           = "recordId"
	PatternSourceShortRecordId      = "shortRecordId"
	PatternSourceRoutineTaskId      = "routineTaskId"
	PatternSourceBlockText          = "blockText"
	PatternSourceBlockCheckboxCount = "blockCheckboxCount"
)

type PatternResolverInterface interface {
	Resolve(ctx context.Context, task schemas.RoutineTask, ownerId uuid.UUID, pattern dtos.RoutineTaskPattern) (map[string]string, *exceptions.Exception)
	ResolveMany(ctx context.Context, tasks []schemas.RoutineTask, ownerIds []uuid.UUID, patterns []dtos.RoutineTaskPattern) ([]map[string]string, []bool, *exceptions.Exception)
}

type PatternResolver struct {
	blockPatternResolver     BlockPatternResolverInterface
	blockPackPatternResolver BlockPackPatternResolverInterface
}

func NewPatternResolver(
	db *gorm.DB,
	blockRepository repositories.BlockRepositoryInterface,
	blockPackRepository repositories.BlockPackRepositoryInterface,
) PatternResolverInterface {
	return PatternResolver{
		blockPatternResolver:     NewBlockPatternResolver(db, blockRepository),
		blockPackPatternResolver: NewBlockPackPatternResolver(db, blockPackRepository),
	}
}

func (r PatternResolver) Resolve(
	ctx context.Context,
	task schemas.RoutineTask,
	ownerId uuid.UUID,
	pattern dtos.RoutineTaskPattern,
) (map[string]string, *exceptions.Exception) {
	values, successes, exception := r.ResolveMany(ctx, []schemas.RoutineTask{task}, []uuid.UUID{ownerId}, []dtos.RoutineTaskPattern{pattern})
	if exception != nil {
		return nil, exception
	}
	if len(successes) == 0 || !successes[0] {
		return nil, exceptions.RoutineTask.InvalidDto()
	}
	return values[0], nil
}

func (r PatternResolver) ResolveMany(
	ctx context.Context,
	tasks []schemas.RoutineTask,
	ownerIds []uuid.UUID,
	patterns []dtos.RoutineTaskPattern,
) ([]map[string]string, []bool, *exceptions.Exception) {
	values := make([]map[string]string, len(patterns))
	successes := make([]bool, len(patterns))
	for index := range patterns {
		values[index] = make(map[string]string, len(patterns[index]))
		successes[index] = true
	}
	if len(tasks) != len(patterns) || len(ownerIds) != len(patterns) {
		return nil, nil, exceptions.RoutineTask.InvalidDto().
			WithOrigin(fmt.Errorf("tasks, ownerIds and patterns length mismatch"))
	}

	hasBlockPatternSource := false
	hasBlockPackPatternSource := false
	for _, pattern := range patterns {
		for _, binding := range pattern {
			switch binding.Source {
			case PatternSourceBlockText:
				hasBlockPatternSource = true
			case PatternSourceBlockCheckboxCount:
				hasBlockPackPatternSource = true
			}
		}
	}

	if hasBlockPatternSource {
		blockValues, blockSuccesses, exception := r.blockPatternResolver.ResolveMany(ctx, ownerIds, patterns)
		if exception != nil {
			return nil, nil, exception
		}
		for index, success := range blockSuccesses {
			if !success {
				successes[index] = false
			}
			for key, value := range blockValues[index] {
				values[index][key] = value
			}
		}
	}

	if hasBlockPackPatternSource {
		blockPackValues, blockPackSuccesses, exception := r.blockPackPatternResolver.ResolveMany(ctx, ownerIds, patterns)
		if exception != nil {
			return nil, nil, exception
		}
		for index, success := range blockPackSuccesses {
			if !success {
				successes[index] = false
			}
			for key, value := range blockPackValues[index] {
				values[index][key] = value
			}
		}
	}

	for patternIndex, pattern := range patterns {
		for key, binding := range pattern {
			switch binding.Source {
			case PatternSourceScheduledAt:
				scheduledAt := tasks[patternIndex].RecordScheduledAt
				if scheduledAt.IsZero() {
					scheduledAt = tasks[patternIndex].ScheduledAt
				}
				if binding.Timezone != nil && *binding.Timezone != "" {
					location, err := time.LoadLocation(*binding.Timezone)
					if err != nil {
						successes[patternIndex] = false
						continue
					}
					scheduledAt = scheduledAt.In(location)
				}
				format := time.RFC3339
				if binding.Format != nil && *binding.Format != "" {
					format = *binding.Format
				}
				values[patternIndex][key] = scheduledAt.Format(format)

			case PatternSourceRecordId:
				values[patternIndex][key] = tasks[patternIndex].RecordId.String()

			case PatternSourceShortRecordId:
				recordId := tasks[patternIndex].RecordId.String()
				if len(recordId) > 8 {
					recordId = recordId[:8]
				}
				values[patternIndex][key] = recordId

			case PatternSourceRoutineTaskId:
				values[patternIndex][key] = tasks[patternIndex].Id.String()

			case PatternSourceBlockText, PatternSourceBlockCheckboxCount:
				continue

			default:
				successes[patternIndex] = false
			}
		}
	}

	return values, successes, nil
}
