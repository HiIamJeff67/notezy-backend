package routinetask

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	handlers "github.com/HiIamJeff67/notezy-backend/app/durablejobs/routinetask/handlers"
	exceptions "github.com/HiIamJeff67/notezy-backend/app/exceptions"
	repositories "github.com/HiIamJeff67/notezy-backend/app/models/repositories"
	schemas "github.com/HiIamJeff67/notezy-backend/app/models/schemas"
	enums "github.com/HiIamJeff67/notezy-backend/app/models/schemas/enums"
	scopes "github.com/HiIamJeff67/notezy-backend/app/models/scopes"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type HandlerManager struct {
	maxWorkers    int
	activeWorkers atomic.Int32
	workerPool    sync.WaitGroup
	db            *gorm.DB
	registries    map[enums.RoutineTaskPurpose]handlers.PurposeHandlerFunc
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

func (hm *HandlerManager) Manage(ctx context.Context, claimedTasks []schemas.RoutineTask) *exceptions.Exception {
	if len(claimedTasks) == 0 {
		return nil
	}

	taskIdToOwnerId, exception := hm.resolveTaskOwners(ctx, claimedTasks)
	if exception != nil {
		return exception
	}

	type purposeTaskGroup struct {
		handler handlers.PurposeHandlerFunc
		tasks   []schemas.RoutineTask
	}
	groupsByPurpose := make(map[enums.RoutineTaskPurpose]purposeTaskGroup)
	for _, task := range claimedTasks {
		if _, exists := taskIdToOwnerId[task.Id]; !exists {
			_ = hm.markFailed(ctx, task)
			continue
		}

		registry, exists := hm.registries[task.Purpose]
		if !exists {
			_ = hm.markFailed(ctx, task)
			continue
		}

		group := groupsByPurpose[task.Purpose]
		group.handler = registry
		group.tasks = append(group.tasks, task)
		groupsByPurpose[task.Purpose] = group
	}
	if len(groupsByPurpose) == 0 {
		return nil
	}

	sem := make(chan struct{}, hm.maxWorkers)
	for _, taskGroup := range groupsByPurpose {
		group := taskGroup
		sem <- struct{}{}
		hm.workerPool.Add(1)
		hm.activeWorkers.Add(1)
		go func() {
			defer func() {
				<-sem
				hm.activeWorkers.Add(-1)
				hm.workerPool.Done()
			}()

			handlerResults := group.handler(ctx, group.tasks, taskIdToOwnerId)
			for _, task := range group.tasks {
				if handlerResults[task.Id] != nil {
					_ = hm.markFailed(ctx, task)
					continue
				}
				_ = hm.markSucceeded(ctx, task)
			}
		}()
	}

	hm.workerPool.Wait()
	return nil
}

func (hm *HandlerManager) resolveTaskOwners(
	ctx context.Context,
	tasks []schemas.RoutineTask,
) (map[uuid.UUID]uuid.UUID, *exceptions.Exception) {
	stationIdSet := make(map[uuid.UUID]bool)
	stationIds := make([]uuid.UUID, 0, len(tasks))
	for _, task := range tasks {
		if stationIdSet[task.StationId] {
			continue
		}
		stationIdSet[task.StationId] = true
		stationIds = append(stationIds, task.StationId)
	}

	var rows []struct {
		Id      uuid.UUID `gorm:"column:id;"`
		OwnerId uuid.UUID `gorm:"column:owner_id;"`
	}
	if err := hm.db.WithContext(ctx).
		Model(&schemas.Station{}).
		Select("id, owner_id").
		Where("id IN ? AND deleted_at IS NULL", stationIds).
		Find(&rows).Error; err != nil {
		return nil, exceptions.Station.NotFound().WithOrigin(err)
	}

	ownerIdByStationId := make(map[uuid.UUID]uuid.UUID, len(rows))
	for _, row := range rows {
		ownerIdByStationId[row.Id] = row.OwnerId
	}

	taskIdToOwnerId := make(map[uuid.UUID]uuid.UUID, len(tasks))
	for _, task := range tasks {
		ownerId, exists := ownerIdByStationId[task.StationId]
		if !exists {
			continue
		}
		taskIdToOwnerId[task.Id] = ownerId
	}
	return taskIdToOwnerId, nil
}

func (hm *HandlerManager) markSucceeded(ctx context.Context, task schemas.RoutineTask) *exceptions.Exception {
	now := time.Now()
	updates := map[string]any{
		"actual_ended_at": now,
		"updated_at":      now,
	}
	if task.Period == nil {
		updates["status"] = enums.RoutineTaskStatus_Success
	} else {
		nextScheduledAt := task.ScheduledAt
		for !nextScheduledAt.After(now) {
			switch *task.Period {
			case enums.RoutinePeriod_Daily:
				nextScheduledAt = nextScheduledAt.AddDate(0, 0, 1)
			case enums.RoutinePeriod_Weekly:
				nextScheduledAt = nextScheduledAt.AddDate(0, 0, 7)
			case enums.RoutinePeriod_Monthly:
				nextScheduledAt = nextScheduledAt.AddDate(0, 1, 0)
			default:
				nextScheduledAt = now.Add(time.Minute)
			}
		}

		updates["status"] = enums.RoutineTaskStatus_Waiting
		updates["attempts"] = 0
		updates["scheduled_at"] = nextScheduledAt
	}

	result := hm.db.WithContext(ctx).
		Model(&schemas.RoutineTask{}).
		Where("id = ? AND status = ?", task.Id, enums.RoutineTaskStatus_Running).
		Updates(updates)
	if result.Error != nil {
		return exceptions.RoutineTask.FailedToUpdate().WithOrigin(result.Error)
	}
	return nil
}

func (hm *HandlerManager) markFailed(ctx context.Context, task schemas.RoutineTask) *exceptions.Exception {
	now := time.Now()
	status := enums.RoutineTaskStatus_Waiting
	if task.Attempts >= task.MaxAttempts {
		status = enums.RoutineTaskStatus_Fail
	}

	result := hm.db.WithContext(ctx).
		Model(&schemas.RoutineTask{}).
		Where("id = ? AND status = ?", task.Id, enums.RoutineTaskStatus_Running).
		Updates(map[string]any{
			"status":          status,
			"actual_ended_at": now,
			"updated_at":      now,
		})
	if result.Error != nil {
		return exceptions.RoutineTask.FailedToUpdate().WithOrigin(result.Error)
	}
	return nil
}
