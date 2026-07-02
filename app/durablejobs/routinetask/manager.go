package routinetask

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	handlers "github.com/HiIamJeff67/notezy-backend/app/durablejobs/routinetask/handlers"
	exceptions "github.com/HiIamJeff67/notezy-backend/app/exceptions"
	repositories "github.com/HiIamJeff67/notezy-backend/app/models/repositories"
	schemas "github.com/HiIamJeff67/notezy-backend/app/models/schemas"
	enums "github.com/HiIamJeff67/notezy-backend/app/models/schemas/enums"
	scopes "github.com/HiIamJeff67/notezy-backend/app/models/scopes"
)

type HandlerManager struct {
	maxWorkers    int
	activeWorkers atomic.Int32
	workerPool    sync.WaitGroup
	sem           chan struct{}
	failed        []routineTaskWithRecord
	failedMutex   sync.Mutex
	success       []routineTaskWithRecord
	successMutex  sync.Mutex
	db            *gorm.DB
	registries    map[enums.RoutineTaskPurpose]handlers.PurposeHandlerFunc
}

type routineTaskWithRecord struct {
	task   schemas.RoutineTask
	record schemas.RoutineTaskRecord
}

type purposeTaskGroup struct {
	handlerFunc handlers.PurposeHandlerFunc
	tasks       []schemas.RoutineTask
}

func NewHandlerManager(maxWorkers int, db *gorm.DB) HandlerManager {
	if maxWorkers <= 0 {
		maxWorkers = 1
	}
	rootShelfRepository := repositories.NewRootShelfRepository(scopes.NewRootShelfScope())
	subShelfRepository := repositories.NewSubShelfRepository(scopes.NewSubShelfScope())
	materialRepository := repositories.NewMaterialRepository(scopes.NewMaterialScope())
	blockPackRepository := repositories.NewBlockPackRepository(scopes.NewBlockPackScope())
	blockGroupRepository := repositories.NewBlockGroupRepository(scopes.NewBlockGroupScope())
	blockRepository := repositories.NewBlockRepository(scopes.NewBlockScope())
	routineRepository := repositories.NewRoutineRepository(scopes.NewRoutineScope())

	blockPackHandler := handlers.NewBlockPackHandler(
		db,
		nil,
		blockPackRepository,
		blockGroupRepository,
		blockRepository,
	)
	blockHandler := handlers.NewBlockHandler(
		db,
		nil,
		blockPackRepository,
		blockGroupRepository,
		blockRepository,
	)
	rootShelfHandler := handlers.NewRootShelfHandler(db, rootShelfRepository, subShelfRepository)
	subShelfHandler := handlers.NewSubShelfHandler(db, subShelfRepository, materialRepository, blockPackRepository)
	routineHandler := handlers.NewRoutineHandler(db, routineRepository)

	return HandlerManager{
		maxWorkers:    maxWorkers,
		activeWorkers: atomic.Int32{},
		sem:           make(chan struct{}, maxWorkers),
		db:            db,
		registries: map[enums.RoutineTaskPurpose]handlers.PurposeHandlerFunc{
			enums.RoutineTaskPurpose_CreateRootShelf: rootShelfHandler.HandleCreateRootShelf,
			enums.RoutineTaskPurpose_UpdateRootShelf: rootShelfHandler.HandleUpdateRootShelf,
			enums.RoutineTaskPurpose_ResetRootShelf:  rootShelfHandler.HandleResetRootShelf,
			enums.RoutineTaskPurpose_CreateSubShelf:  subShelfHandler.HandleCreateSubShelf,
			enums.RoutineTaskPurpose_UpdateSubShelf:  subShelfHandler.HandleUpdateSubShelf,
			enums.RoutineTaskPurpose_ResetSubShelf:   subShelfHandler.HandleResetSubShelf,
			enums.RoutineTaskPurpose_CreateBlockPack: blockPackHandler.HandleCreateBlockPack,
			enums.RoutineTaskPurpose_UpdateBlockPack: blockPackHandler.HandleUpdateBlockPack,
			enums.RoutineTaskPurpose_ResetBlockPack:  blockPackHandler.HandleResetBlockPack,
			enums.RoutineTaskPurpose_AppendBlock:     blockHandler.HandleAppendBlock,
			enums.RoutineTaskPurpose_UpdateBlock:     blockHandler.HandleUpdateBlock,
			enums.RoutineTaskPurpose_ResetBlock:      blockHandler.HandleResetBlock,
			enums.RoutineTaskPurpose_CreateRoutine:   routineHandler.HandleCreateRoutine,
			enums.RoutineTaskPurpose_UpdateRoutine:   routineHandler.HandleUpdateRoutine,
		},
	}
}

/* ============================== Util & Helpers for Routine Tasks and Routine Task Records ============================== */

func (hm *HandlerManager) getErrorDetails(exception *exceptions.Exception) (enums.RoutineTaskRecordErrorCode, string) {
	if exception == nil {
		return enums.RoutineTaskRecordErrorCode_HandlerFailed, "Routine task handler returned unsuccessful result"
	}
	reason := exception.Reason
	if len(reason) > 256 {
		reason = reason[:256]
	}

	if errors.Is(exception.Origin, context.Canceled) {
		return enums.RoutineTaskRecordErrorCode_Canceled, reason
	}
	if errors.Is(exception.Origin, context.DeadlineExceeded) {
		return enums.RoutineTaskRecordErrorCode_Timeout, reason
	}

	switch exception.HTTPStatusCode {
	case http.StatusUnauthorized, http.StatusForbidden:
		return enums.RoutineTaskRecordErrorCode_PermissionDenied, reason
	case http.StatusNotFound:
		return enums.RoutineTaskRecordErrorCode_TargetNotFound, reason
	case http.StatusBadRequest:
		return enums.RoutineTaskRecordErrorCode_PayloadInvalid, reason
	case http.StatusInternalServerError:
		return enums.RoutineTaskRecordErrorCode_DatabaseError, reason
	default:
		return enums.RoutineTaskRecordErrorCode_HandlerFailed, reason
	}
}

func (hm *HandlerManager) newRecord(
	task schemas.RoutineTask,
	status enums.RoutineTaskRecordStatus,
	endedAt time.Time,
	errorCode *enums.RoutineTaskRecordErrorCode,
	errorReason *string,
) schemas.RoutineTaskRecord {
	scheduledAt := task.RecordScheduledAt
	if scheduledAt.IsZero() {
		scheduledAt = task.ScheduledAt
	}

	return schemas.RoutineTaskRecord{
		Id:              task.RecordId,
		RoutineTaskId:   task.Id,
		Purpose:         task.Purpose,
		Status:          status,
		ErrorCode:       errorCode,
		ErrorReason:     errorReason,
		CostUnit:        task.CostUnit,
		TotalAttempts:   int64(task.Attempts),
		ScheduledAt:     scheduledAt,
		ActualStartedAt: task.ActualStartedAt,
		ActualEndedAt:   &endedAt,
	}
}

func (hm *HandlerManager) resetRoutineTasksWithRecords(capacity int) {
	hm.failedMutex.Lock()
	hm.failed = make([]routineTaskWithRecord, 0, capacity)
	hm.failedMutex.Unlock()

	hm.successMutex.Lock()
	hm.success = make([]routineTaskWithRecord, 0, capacity)
	hm.successMutex.Unlock()
}

func (hm *HandlerManager) appendFailedRoutineTaskWithRecord(task schemas.RoutineTask, record schemas.RoutineTaskRecord) {
	hm.failedMutex.Lock()
	hm.failed = append(hm.failed, routineTaskWithRecord{task: task, record: record})
	hm.failedMutex.Unlock()
}

func (hm *HandlerManager) appendSuccessRoutineTaskWithRecord(task schemas.RoutineTask, record schemas.RoutineTaskRecord) {
	hm.successMutex.Lock()
	hm.success = append(hm.success, routineTaskWithRecord{task: task, record: record})
	hm.successMutex.Unlock()
}

func (hm *HandlerManager) finalize(ctx context.Context) *exceptions.Exception {
	hm.failedMutex.Lock()
	failed := append([]routineTaskWithRecord(nil), hm.failed...)
	hm.failedMutex.Unlock()

	hm.successMutex.Lock()
	success := append([]routineTaskWithRecord(nil), hm.success...)
	hm.successMutex.Unlock()

	if len(failed)+len(success) == 0 {
		return nil
	}

	tx := hm.db.WithContext(ctx).Begin()

	now := time.Now()
	successTaskIds := make([]uuid.UUID, len(success))
	successRecordIds := make([]uuid.UUID, len(success))
	for index, item := range success {
		successTaskIds[index] = item.task.Id
		successRecordIds[index] = item.record.Id
	}

	if len(successTaskIds) > 0 {
		result := tx.
			Model(&schemas.RoutineTask{}).
			Where("id IN ? AND status = ?", successTaskIds, enums.RoutineTaskStatus_Running).
			Updates(map[string]any{
				"status":          enums.RoutineTaskStatus_Idle,
				"attempts":        0,
				"actual_ended_at": now,
				"updated_at":      now,
			})
		if result.Error != nil {
			tx.Rollback()
			return exceptions.RoutineTask.FailedToUpdate().WithOrigin(result.Error)
		}
		if result.RowsAffected != int64(len(successTaskIds)) {
			tx.Rollback()
			return exceptions.RoutineTask.FailedToUpdate().WithOrigin(errors.New("routine task success finalize count mismatch"))
		}

		result = tx.
			Model(&schemas.RoutineTaskRecord{}).
			Where("id IN ? AND status = ?", successRecordIds, enums.RoutineTaskRecordStatus_Running).
			Updates(map[string]any{
				"status":          enums.RoutineTaskRecordStatus_Success,
				"actual_ended_at": now,
				"error_code":      nil,
				"error_reason":    nil,
				"updated_at":      now,
			})
		if result.Error != nil {
			tx.Rollback()
			return exceptions.RoutineTask.FailedToUpdate().WithOrigin(result.Error)
		}
		if result.RowsAffected != int64(len(successRecordIds)) {
			tx.Rollback()
			return exceptions.RoutineTask.FailedToUpdate().WithOrigin(errors.New("routine task record success finalize count mismatch"))
		}
	}

	failedTaskIds := make([]uuid.UUID, 0, len(failed))
	type failedRecordGroupKey struct {
		errorCode   enums.RoutineTaskRecordErrorCode
		errorReason string
	}
	failedRecordIdsByGroup := make(map[failedRecordGroupKey][]uuid.UUID)
	for _, item := range failed {
		failedTaskIds = append(failedTaskIds, item.task.Id)
		errorCode := enums.RoutineTaskRecordErrorCode_HandlerFailed
		if item.record.ErrorCode != nil {
			errorCode = *item.record.ErrorCode
		}
		errorReason := "Routine task handler returned unsuccessful result"
		if item.record.ErrorReason != nil {
			errorReason = *item.record.ErrorReason
		}
		key := failedRecordGroupKey{errorCode: errorCode, errorReason: errorReason}
		failedRecordIdsByGroup[key] = append(failedRecordIdsByGroup[key], item.record.Id)
	}

	if len(failedTaskIds) > 0 {
		result := tx.
			Model(&schemas.RoutineTask{}).
			Where("id IN ? AND status = ?", failedTaskIds, enums.RoutineTaskStatus_Running).
			Updates(map[string]any{
				"status":          enums.RoutineTaskStatus_Idle,
				"actual_ended_at": now,
				"updated_at":      now,
			})
		if result.Error != nil {
			tx.Rollback()
			return exceptions.RoutineTask.FailedToUpdate().WithOrigin(result.Error)
		}
		if result.RowsAffected != int64(len(failedTaskIds)) {
			tx.Rollback()
			return exceptions.RoutineTask.FailedToUpdate().WithOrigin(errors.New("routine task failed finalize count mismatch"))
		}

		updatedFailedRecordCount := int64(0)
		for key, failedRecordIds := range failedRecordIdsByGroup {
			errorCode := key.errorCode
			errorReason := key.errorReason
			result = tx.
				Model(&schemas.RoutineTaskRecord{}).
				Where("id IN ? AND status = ?", failedRecordIds, enums.RoutineTaskRecordStatus_Running).
				Updates(map[string]any{
					"status":          enums.RoutineTaskRecordStatus_Failed,
					"actual_ended_at": now,
					"error_code":      errorCode,
					"error_reason":    errorReason,
					"updated_at":      now,
				})
			if result.Error != nil {
				tx.Rollback()
				return exceptions.RoutineTask.FailedToUpdate().WithOrigin(result.Error)
			}
			updatedFailedRecordCount += result.RowsAffected
		}
		if updatedFailedRecordCount != int64(len(failed)) {
			tx.Rollback()
			return exceptions.RoutineTask.FailedToUpdate().WithOrigin(errors.New("routine task record failed finalize count mismatch"))
		}
	}

	if err := tx.Commit().Error; err != nil {
		return exceptions.RoutineTask.FailedToCommitTransaction().WithOrigin(err)
	}

	return nil
}

/* ============================== Core logic ============================== */

func (hm *HandlerManager) Manage(
	ctx context.Context,
	claimedTasks []schemas.RoutineTask,
	taskIdToOwnerId map[uuid.UUID]uuid.UUID,
) *exceptions.Exception {
	if len(claimedTasks) == 0 {
		return nil
	}

	hm.resetRoutineTasksWithRecords(len(claimedTasks))

	groupsByPurpose := make(map[enums.RoutineTaskPurpose]purposeTaskGroup)
	for _, task := range claimedTasks {
		if _, exists := taskIdToOwnerId[task.Id]; !exists {
			endedAt := time.Now()
			tempErrorCode := enums.RoutineTaskRecordErrorCode_TargetNotFound
			tempErrorReason := "Routine task station owner was not found"
			hm.appendFailedRoutineTaskWithRecord(task, hm.newRecord(
				task,
				enums.RoutineTaskRecordStatus_Failed,
				endedAt,
				&tempErrorCode,
				&tempErrorReason,
			))
			continue
		}

		registry, exists := hm.registries[task.Purpose]
		if !exists {
			endedAt := time.Now()
			tempErrorCode := enums.RoutineTaskRecordErrorCode_HandlerFailed
			tempErrorReason := "Routine task purpose handler was not found"
			hm.appendFailedRoutineTaskWithRecord(task, hm.newRecord(
				task,
				enums.RoutineTaskRecordStatus_Failed,
				endedAt,
				&tempErrorCode,
				&tempErrorReason,
			))
			continue
		}

		group := groupsByPurpose[task.Purpose]
		group.handlerFunc = registry
		group.tasks = append(group.tasks, task)
		groupsByPurpose[task.Purpose] = group
	}
	if len(groupsByPurpose) == 0 {
		return hm.finalize(ctx)
	}

	for _, taskGroup := range groupsByPurpose {
		group := taskGroup
		hm.sem <- struct{}{}
		hm.workerPool.Add(1)
		hm.activeWorkers.Add(1)
		go func() {
			defer func() {
				<-hm.sem
				hm.activeWorkers.Add(-1)
				hm.workerPool.Done()
			}()

			handlerResults, exception := group.handlerFunc(ctx, group.tasks, taskIdToOwnerId)
			for index, task := range group.tasks {
				endedAt := time.Now()
				if exception != nil || index >= len(handlerResults) || !handlerResults[index] { // if the task was failed
					errorCode, errorReason := hm.getErrorDetails(exception)
					hm.appendFailedRoutineTaskWithRecord(task, hm.newRecord(task, enums.RoutineTaskRecordStatus_Failed, endedAt, &errorCode, &errorReason))
				} else { // if the task was success
					hm.appendSuccessRoutineTaskWithRecord(task, hm.newRecord(task, enums.RoutineTaskRecordStatus_Success, endedAt, nil, nil))
				}
			}
		}()
	}

	hm.workerPool.Wait()
	if exception := hm.finalize(ctx); exception != nil {
		return exception
	}
	return nil
}
