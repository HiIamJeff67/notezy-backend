package routinetask

import (
	"context"

	"gorm.io/gorm"

	exceptions "github.com/HiIamJeff67/notezy-backend/app/exceptions"
	schemas "github.com/HiIamJeff67/notezy-backend/app/models/schemas"
	enums "github.com/HiIamJeff67/notezy-backend/app/models/schemas/enums"
	constants "github.com/HiIamJeff67/notezy-backend/shared/constants"
)

type Claimer struct {
	db *gorm.DB
}

func NewClaimer(db *gorm.DB) Claimer {
	return Claimer{
		db: db,
	}
}

func (c *Claimer) Claim(
	ctx context.Context,
) ([]schemas.RoutineTask, *exceptions.Exception) {
	var claimedTasks []schemas.RoutineTask
	result := c.db.
		WithContext(ctx).
		Raw(`WITH claimable_tasks AS (
				SELECT id
				FROM "RoutineTaskTable"
				WHERE status = ?
				AND scheduled_at <= NOW()
				AND attempts < max_attempts
				ORDER BY priority DESC, scheduled_at ASC, id ASC
				FOR UPDATE SKIP LOCKED
				LIMIT ?
			)
			UPDATE "RoutineTaskTable" AS routine_task
			SET
				status = ?,
				attempts = routine_task.attempts + 1,
				actual_started_at = NOW(),
				actual_ended_at = NULL,
				updated_at = NOW()
			FROM claimable_tasks
			WHERE routine_task.id = claimable_tasks.id
			RETURNING routine_task.*;
			`,
			enums.RoutineTaskStatus_Waiting,
			constants.RoutineTaskClaimerMaxClaimableTasks,
			enums.RoutineTaskStatus_Running,
		).
		Scan(&claimedTasks)
	if result.Error != nil {
		return nil, exceptions.RoutineTask.
			FailedToClaim("routine tasks").
			WithOrigin(result.Error)
	}

	return claimedTasks, nil
}
