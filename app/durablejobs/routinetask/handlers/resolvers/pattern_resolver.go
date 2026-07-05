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
	values := make(map[string]string, len(pattern))
	if len(pattern) == 0 {
		return values, nil
	}

	blockValues, exception := r.blockPatternResolver.Resolve(ctx, ownerId, pattern)
	if exception != nil {
		return nil, exception
	}
	for key, value := range blockValues {
		values[key] = value
	}

	blockPackValues, exception := r.blockPackPatternResolver.Resolve(ctx, ownerId, pattern)
	if exception != nil {
		return nil, exception
	}
	for key, value := range blockPackValues {
		values[key] = value
	}

	for key, binding := range pattern {
		switch binding.Source {
		case PatternSourceScheduledAt:
			scheduledAt := task.RecordScheduledAt
			if scheduledAt.IsZero() {
				scheduledAt = task.ScheduledAt
			}
			if binding.Timezone != nil && *binding.Timezone != "" {
				location, err := time.LoadLocation(*binding.Timezone)
				if err != nil {
					return nil, exceptions.RoutineTask.InvalidDto().WithOrigin(err)
				}
				scheduledAt = scheduledAt.In(location)
			}
			format := time.RFC3339
			if binding.Format != nil && *binding.Format != "" {
				format = *binding.Format
			}
			values[key] = scheduledAt.Format(format)

		case PatternSourceRecordId:
			values[key] = task.RecordId.String()

		case PatternSourceShortRecordId:
			recordId := task.RecordId.String()
			if len(recordId) > 8 {
				recordId = recordId[:8]
			}
			values[key] = recordId

		case PatternSourceRoutineTaskId:
			values[key] = task.Id.String()

		case PatternSourceBlockText, PatternSourceBlockCheckboxCount:
			continue

		default:
			return nil, exceptions.RoutineTask.InvalidDto().
				WithOrigin(fmt.Errorf("unsupported pattern source: %s", binding.Source))
		}
	}

	return values, nil
}
