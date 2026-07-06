package routinetask

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

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
) ([]schemas.RoutineTask, map[uuid.UUID]uuid.UUID, *exceptions.Exception) {
	var claimedTasks []schemas.RoutineTask
	taskIdToOwnerId := make(map[uuid.UUID]uuid.UUID)

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
		return nil, nil, exceptions.RoutineTask.FailedToClaim("routine tasks").WithOrigin(result.Error)
	}

	if len(claimableTasks) == 0 {
		if err := tx.Commit().Error; err != nil {
			return nil, nil, exceptions.RoutineTask.FailedToCommitTransaction().WithOrigin(err)
		}
		return claimedTasks, taskIdToOwnerId, nil
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
		return nil, nil, exceptions.RoutineTask.FailedToClaim("routine tasks").WithOrigin(result.Error)
	}

	result = tx.
		Model(&schemas.RoutineTask{}).
		Where("id IN ?", claimedTaskIds).
		Find(&claimedTasks)
	if result.Error != nil {
		tx.Rollback()
		return nil, nil, exceptions.RoutineTask.FailedToClaim("routine tasks").WithOrigin(result.Error)
	}

	routineTaskRecords := make([]schemas.RoutineTaskRecord, len(claimedTasks))
	routineIds := make([]uuid.UUID, 0, len(claimedTasks))
	routineIdSet := make(map[uuid.UUID]bool, len(claimedTasks))
	for index, claimedTask := range claimedTasks {
		recordScheduledAt := recordScheduledAtByTaskId[claimedTask.Id]
		claimedTasks[index].RecordScheduledAt = recordScheduledAt
		claimedTasks[index].RecordId = uuid.New()
		if !routineIdSet[claimedTask.RoutineId] {
			routineIdSet[claimedTask.RoutineId] = true
			routineIds = append(routineIds, claimedTask.RoutineId)
		}
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
		return nil, nil, exceptions.RoutineTask.FailedToClaim("routine task records").WithOrigin(result.Error)
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
		return nil, nil, exceptions.RoutineTask.FailedToClaim("routine task records").WithOrigin(result.Error)
	}

	if len(fetchedRoutineTaskRecords) != len(routineTaskRecords) {
		tx.Rollback()
		return nil, nil, exceptions.RoutineTask.FailedToClaim("routine task records")
	}

	var routines []struct {
		RoutineId uuid.UUID `gorm:"column:routine_id;"`
		OwnerId   uuid.UUID `gorm:"column:owner_id;"`
	}
	result = tx.
		Model(&schemas.Routine{}).
		Select(`"RoutineTable".id AS routine_id, station.owner_id AS owner_id`).
		Joins(`INNER JOIN "StationTable" station ON station.id = "RoutineTable".station_id AND station.deleted_at IS NULL`).
		Where(`"RoutineTable".id IN ? AND "RoutineTable".deleted_at IS NULL`, routineIds).
		Find(&routines)
	if result.Error != nil {
		tx.Rollback()
		return nil, nil, exceptions.RoutineTask.FailedToClaim("routine task owners").WithOrigin(result.Error)
	}

	ownerIdByRoutineId := make(map[uuid.UUID]uuid.UUID, len(routines))
	for _, routine := range routines {
		ownerIdByRoutineId[routine.RoutineId] = routine.OwnerId
	}
	for _, claimedTask := range claimedTasks {
		ownerId, exists := ownerIdByRoutineId[claimedTask.RoutineId]
		if exists {
			taskIdToOwnerId[claimedTask.Id] = ownerId
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, nil, exceptions.RoutineTask.FailedToCommitTransaction().WithOrigin(err)
	}

	return claimedTasks, taskIdToOwnerId, nil
}
