package routinetask

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	exceptions "github.com/HiIamJeff67/notezy-backend/internal/exceptions"
	schemas "github.com/HiIamJeff67/notezy-backend/internal/services/durablejob/data/schemas"
	enums "github.com/HiIamJeff67/notezy-backend/internal/services/durablejob/data/schemas/enums"
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

	tx := c.db.WithContext(ctx).Begin()

	type claimableTask struct {
		Id          uuid.UUID `gorm:"column:id;"`
		ScheduledAt time.Time `gorm:"column:scheduled_at;"`
	}

	now := time.Now()
	var claimableTasks []claimableTask
	result := tx.
		Model(&schemas.RoutineTask{}).
		Select("id, scheduled_at").
		Where("status IN ?", []enums.RoutineTaskStatus{enums.RoutineTaskStatus_Idle}).
		Where("scheduled_at <= ?", now).
		Where("attempts < max_attempts").
		Order("priority DESC, scheduled_at ASC, id ASC").
		Clauses(clause.Locking{Strength: "UPDATE", Options: "SKIP LOCKED"}).
		Limit(constants.RoutineTaskClaimerMaxClaimableTasks).
		Find(&claimableTasks)
	if result.Error != nil {
		tx.Rollback()
		return nil, exceptions.New(
			"ClaimFailed",
			"RoutineTask",
			"Claim",
			"Failed to claim routine tasks",
			http.StatusInternalServerError,
			true,
		).WithOrigin(result.Error)
	}

	if len(claimableTasks) == 0 {
		if err := tx.Commit().Error; err != nil {
			return nil, exceptions.New(
				"TransactionCommitFailed",
				"RoutineTask",
				"Claim",
				"Failed to commit the routine task claim transaction",
				http.StatusInternalServerError,
				true,
			).WithOrigin(err)
		}
		return claimedTasks, nil
	}

	claimedTaskIds := make([]uuid.UUID, len(claimableTasks))
	recordScheduledAtByTaskId := make(map[uuid.UUID]time.Time, len(claimableTasks))
	for index, claimableTask := range claimableTasks {
		claimedTaskIds[index] = claimableTask.Id
		recordScheduledAtByTaskId[claimableTask.Id] = claimableTask.ScheduledAt
	}

	result = tx.
		Model(&schemas.RoutineTask{}).
		Where("id IN ?", claimedTaskIds).
		Updates(map[string]any{
			"status":   enums.RoutineTaskStatus_Running,
			"attempts": gorm.Expr("attempts + 1"),
			"scheduled_at": gorm.Expr(
				`CASE period
					WHEN ? THEN GREATEST(scheduled_at, next_scheduled_at) + INTERVAL '1 day'
					WHEN ? THEN GREATEST(scheduled_at, next_scheduled_at) + INTERVAL '7 days'
					WHEN ? THEN GREATEST(scheduled_at, next_scheduled_at) + INTERVAL '30 days'
					ELSE GREATEST(scheduled_at, next_scheduled_at)
				END`,
				enums.RoutinePeriod_Daily,
				enums.RoutinePeriod_Weekly,
				enums.RoutinePeriod_Monthly,
			),
			"next_scheduled_at": gorm.Expr(
				`CASE period
					WHEN ? THEN GREATEST(scheduled_at, next_scheduled_at) + INTERVAL '1 day'
					WHEN ? THEN GREATEST(scheduled_at, next_scheduled_at) + INTERVAL '7 days'
					WHEN ? THEN GREATEST(scheduled_at, next_scheduled_at) + INTERVAL '30 days'
					ELSE GREATEST(scheduled_at, next_scheduled_at)
				END`,
				enums.RoutinePeriod_Daily,
				enums.RoutinePeriod_Weekly,
				enums.RoutinePeriod_Monthly,
			),
			"actual_started_at": now,
			"actual_ended_at":   nil,
			"updated_at":        now,
		})
	if result.Error != nil {
		tx.Rollback()
		return nil, exceptions.New(
			"ClaimFailed",
			"RoutineTask",
			"Claim",
			"Failed to claim routine tasks",
			http.StatusInternalServerError,
			true,
		).WithOrigin(result.Error)
	}

	result = tx.
		Model(&schemas.RoutineTask{}).
		Where("id IN ?", claimedTaskIds).
		Find(&claimedTasks)
	if result.Error != nil {
		tx.Rollback()
		return nil, exceptions.New(
			"ClaimFailed",
			"RoutineTask",
			"Claim",
			"Failed to claim routine tasks",
			http.StatusInternalServerError,
			true,
		).WithOrigin(result.Error)
	}

	routineTaskRecords := make([]schemas.RoutineTaskRecord, len(claimedTasks))
	for index, claimedTask := range claimedTasks {
		recordScheduledAt := recordScheduledAtByTaskId[claimedTask.Id]
		claimedTasks[index].RecordScheduledAt = recordScheduledAt
		claimedTasks[index].RecordId = uuid.New()
		routineTaskRecords[index] = schemas.RoutineTaskRecord{
			Id:              claimedTasks[index].RecordId,
			RoutineTaskId:   claimedTask.Id,
			Purpose:         claimedTask.Purpose,
			Status:          enums.RoutineTaskRecordStatus_Running,
			CostUnit:        claimedTask.CostUnit,
			TotalAttempts:   int64(claimedTask.Attempts),
			ScheduledAt:     recordScheduledAt,
			ActualStartedAt: claimedTask.ActualStartedAt,
		}
	}

	result = tx.CreateInBatches(&routineTaskRecords, constants.RoutineTaskClaimerMaxClaimableTasks)
	if result.Error != nil {
		tx.Rollback()
		return nil, exceptions.New(
			"ClaimFailed",
			"RoutineTaskRecord",
			"Claim",
			"Failed to create routine task records for claimed tasks",
			http.StatusInternalServerError,
			true,
		).WithOrigin(result.Error)
	}

	recordIds := make([]uuid.UUID, len(routineTaskRecords))
	for index, routineTaskRecord := range routineTaskRecords {
		recordIds[index] = routineTaskRecord.Id
	}

	var fetchedRoutineTaskRecords []schemas.RoutineTaskRecord
	result = tx.
		Model(&schemas.RoutineTaskRecord{}).
		Where("id IN ?", recordIds).
		Find(&fetchedRoutineTaskRecords)
	if result.Error != nil {
		tx.Rollback()
		return nil, exceptions.New(
			"ClaimFailed",
			"RoutineTaskRecord",
			"Claim",
			"Failed to retrieve routine task records for claimed tasks",
			http.StatusInternalServerError,
			true,
		).WithOrigin(result.Error)
	}

	if len(fetchedRoutineTaskRecords) != len(routineTaskRecords) {
		tx.Rollback()
		return nil, exceptions.New(
			"ClaimFailed",
			"RoutineTaskRecord",
			"Claim",
			"Routine task record claim count did not match",
			http.StatusInternalServerError,
			true,
		)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, exceptions.New(
			"TransactionCommitFailed",
			"RoutineTask",
			"Claim",
			"Failed to commit the routine task claim transaction",
			http.StatusInternalServerError,
			true,
		).WithOrigin(err)
	}

	return claimedTasks, nil
}
