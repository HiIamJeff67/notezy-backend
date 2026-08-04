package routinetask

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	coreeventscontract "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/events"
	durablejobcontract "github.com/HiIamJeff67/notezy-backend/contracts/durablejob/v1"
	durablejobroutinetasktypes "github.com/HiIamJeff67/notezy-backend/contracts/durablejob/v1/types/routine-tasks"
	enumcontract "github.com/HiIamJeff67/notezy-backend/contracts/types/enums"
	exceptions "github.com/HiIamJeff67/notezy-backend/internal/exceptions"
	platformkafka "github.com/HiIamJeff67/notezy-backend/internal/platform/kafka"
	inputs "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/database/inputs"
	options "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/database/options"
	repositories "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/database/repositories"
	schemas "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/database/schemas"
	enums "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/database/schemas/enums"
	scopes "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/database/scopes"
	durablejobconfig "github.com/HiIamJeff67/notezy-backend/internal/services/durablejob/config"
	handlers "github.com/HiIamJeff67/notezy-backend/internal/services/durablejob/routinetask/handlers"
	matchers "github.com/HiIamJeff67/notezy-backend/internal/services/durablejob/routinetask/handlers/matchers"
	resolvers "github.com/HiIamJeff67/notezy-backend/internal/services/durablejob/routinetask/handlers/resolvers"
	validation "github.com/HiIamJeff67/notezy-backend/internal/services/durablejob/validation"
	yjsmaintenance "github.com/HiIamJeff67/notezy-backend/internal/services/durablejob/yjsmaintenance"
)

type HandlerManager struct {
	maxWorkers        int
	activeWorkers     atomic.Int32
	workerPool        sync.WaitGroup
	sem               chan struct{}
	workerId          uuid.UUID
	failed            []routineTaskWithRecord
	failedMutex       sync.Mutex
	success           []routineTaskWithRecord
	successMutex      sync.Mutex
	db                *gorm.DB
	routineRepository repositories.RoutineRepositoryInterface
	registries        map[enums.RoutineTaskPurpose]handlers.PurposeHandler
}

type routineTaskWithRecord struct {
	task   schemas.RoutineTask
	record schemas.RoutineTaskRecord
}

type purposeTaskGroup struct {
	handlerFunc        handlers.PurposeHandlerFunc
	allowedPermissions []enums.AccessControlPermission
	tasks              []schemas.RoutineTask
}

func NewHandlerManager(
	maxWorkers int,
	db *gorm.DB,
	config durablejobconfig.Config,
	workerIds ...uuid.UUID,
) HandlerManager {
	validator := validation.New()

	if maxWorkers <= 0 {
		maxWorkers = 1
	}
	rootShelfRepository := repositories.NewRootShelfRepository(scopes.NewRootShelfScope())
	subShelfRepository := repositories.NewSubShelfRepository(scopes.NewSubShelfScope())
	materialRepository := repositories.NewMaterialRepository(scopes.NewMaterialScope())
	blockPackRepository := repositories.NewBlockPackRepository(scopes.NewBlockPackScope())
	blockRepository := repositories.NewBlockRepository(scopes.NewBlockScope())
	routineRepository := repositories.NewRoutineRepository(scopes.NewRoutineScope())
	patternResolver := resolvers.NewPatternResolver(db, blockRepository, blockPackRepository)
	templateBlockMatcher := matchers.NewTemplateBlockMatcher()
	yjsWorkerClient := yjsmaintenance.NewWorkerClient(config)

	blockPackHandler := handlers.NewBlockPackHandler(
		validator,
		db,
		patternResolver,
		templateBlockMatcher,
		blockPackRepository,
		blockRepository,
		yjsWorkerClient,
	)
	blockHandler := handlers.NewBlockHandler(
		validator,
		db,
		patternResolver,
		templateBlockMatcher,
		blockPackRepository,
		blockRepository,
	)
	rootShelfHandler := handlers.NewRootShelfHandler(
		validator,
		db,
		patternResolver,
		templateBlockMatcher,
		rootShelfRepository,
		subShelfRepository,
	)
	subShelfHandler := handlers.NewSubShelfHandler(
		validator,
		db,
		patternResolver,
		templateBlockMatcher,
		subShelfRepository,
		materialRepository,
		blockPackRepository,
	)
	routineHandler := handlers.NewRoutineHandler(
		validator,
		db,
		patternResolver,
		templateBlockMatcher,
		routineRepository,
	)
	readPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
		enums.AccessControlPermission_Read,
	}
	adminPermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
	}
	writePermissions := []enums.AccessControlPermission{
		enums.AccessControlPermission_Owner,
		enums.AccessControlPermission_Admin,
		enums.AccessControlPermission_Write,
	}

	workerId := uuid.New()
	if len(workerIds) > 0 && workerIds[0] != uuid.Nil {
		workerId = workerIds[0]
	}

	return HandlerManager{
		maxWorkers:        maxWorkers,
		activeWorkers:     atomic.Int32{},
		sem:               make(chan struct{}, maxWorkers),
		workerId:          workerId,
		db:                db,
		routineRepository: routineRepository,
		registries: map[enums.RoutineTaskPurpose]handlers.PurposeHandler{
			enums.RoutineTaskPurpose_CreateRootShelf: {
				HandlerFunc:        rootShelfHandler.HandleCreateRootShelf,
				AllowedPermissions: readPermissions,
			},
			enums.RoutineTaskPurpose_UpdateRootShelf: {
				HandlerFunc:        rootShelfHandler.HandleUpdateRootShelf,
				AllowedPermissions: adminPermissions,
			},
			enums.RoutineTaskPurpose_ResetRootShelf: {
				HandlerFunc:        rootShelfHandler.HandleResetRootShelf,
				AllowedPermissions: adminPermissions,
			},
			enums.RoutineTaskPurpose_CreateSubShelf: {
				HandlerFunc:        subShelfHandler.HandleCreateSubShelf,
				AllowedPermissions: adminPermissions,
			},
			enums.RoutineTaskPurpose_UpdateSubShelf: {
				HandlerFunc:        subShelfHandler.HandleUpdateSubShelf,
				AllowedPermissions: adminPermissions,
			},
			enums.RoutineTaskPurpose_ResetSubShelf: {
				HandlerFunc:        subShelfHandler.HandleResetSubShelf,
				AllowedPermissions: adminPermissions,
			},
			enums.RoutineTaskPurpose_CreateBlockPack: {
				HandlerFunc:        blockPackHandler.HandleCreateBlockPack,
				AllowedPermissions: writePermissions,
			},
			enums.RoutineTaskPurpose_UpdateBlockPack: {
				HandlerFunc:        blockPackHandler.HandleUpdateBlockPack,
				AllowedPermissions: writePermissions,
			},
			enums.RoutineTaskPurpose_ResetBlockPack: {
				HandlerFunc:        blockPackHandler.HandleResetBlockPack,
				AllowedPermissions: writePermissions,
			},
			enums.RoutineTaskPurpose_AppendBlock: {
				HandlerFunc:        blockHandler.HandleAppendBlock,
				AllowedPermissions: writePermissions,
			},
			enums.RoutineTaskPurpose_UpdateBlock: {
				HandlerFunc:        blockHandler.HandleUpdateBlock,
				AllowedPermissions: writePermissions,
			},
			enums.RoutineTaskPurpose_ResetBlock: {
				HandlerFunc:        blockHandler.HandleResetBlock,
				AllowedPermissions: writePermissions,
			},
			enums.RoutineTaskPurpose_CreateRoutine: {
				HandlerFunc:        routineHandler.HandleCreateRoutine,
				AllowedPermissions: writePermissions,
			},
			enums.RoutineTaskPurpose_UpdateRoutine: {
				HandlerFunc:        routineHandler.HandleUpdateRoutine,
				AllowedPermissions: writePermissions,
			},
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

	if errors.Is(exception.Origin(), context.Canceled) {
		return enums.RoutineTaskRecordErrorCode_Canceled, reason
	}
	if errors.Is(exception.Origin(), context.DeadlineExceeded) {
		return enums.RoutineTaskRecordErrorCode_Timeout, reason
	}

	switch exception.HTTPStatusCode() {
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

	correlationId := uuid.New().String()
	if len(success) > 0 {
		completed := durablejobcontract.MarkCompletedRoutineTasksRequestDto{
			WorkerId: hm.workerId,
			Tasks:    make([]durablejobroutinetasktypes.CompletedRoutineTask, len(success)),
		}
		for index, item := range success {
			completedAt := time.Now().UTC()
			if item.record.ActualEndedAt != nil {
				completedAt = *item.record.ActualEndedAt
			}
			completed.Tasks[index] = durablejobroutinetasktypes.CompletedRoutineTask{
				RoutineTaskId:       item.task.Id,
				RoutineTaskRecordId: item.record.Id,
				CompletedAt:         completedAt,
			}
		}
		if exception := hm.publishResult(ctx, coreeventscontract.CoreDurableJobRoutineTaskTopic, coreeventscontract.EventType_RoutineTasksCompleted, correlationId, completed); exception != nil {
			return exception
		}
	}
	if len(failed) > 0 {
		failedBatch := durablejobcontract.MarkFailedRoutineTasksRequestDto{
			WorkerId: hm.workerId,
			Tasks:    make([]durablejobroutinetasktypes.FailedRoutineTask, len(failed)),
		}
		for index, item := range failed {
			errorCode := enumcontract.RoutineTaskRecordErrorCode_HandlerFailed
			if item.record.ErrorCode != nil {
				if contractErrorCode := item.record.ErrorCode.ToContractable(); contractErrorCode != nil {
					errorCode = *contractErrorCode
				}
			}
			errorReason := "Routine task handler returned unsuccessful result"
			if item.record.ErrorReason != nil {
				errorReason = *item.record.ErrorReason
			}
			failedAt := time.Now().UTC()
			if item.record.ActualEndedAt != nil {
				failedAt = *item.record.ActualEndedAt
			}
			failedBatch.Tasks[index] = durablejobroutinetasktypes.FailedRoutineTask{
				RoutineTaskId:       item.task.Id,
				RoutineTaskRecordId: item.record.Id,
				FailedAt:            failedAt,
				ErrorCode:           errorCode,
				ErrorReason:         errorReason,
			}
		}
		if exception := hm.publishResult(ctx, coreeventscontract.CoreDurableJobRoutineTaskTopic, coreeventscontract.EventType_RoutineTasksFailed, correlationId, failedBatch); exception != nil {
			return exception
		}
	}

	return nil
}

func (hm *HandlerManager) publishResult(
	ctx context.Context,
	topic coreeventscontract.Topic,
	eventType coreeventscontract.EventType,
	correlationId string,
	data any,
) *exceptions.Exception {
	payload, err := json.Marshal(coreeventscontract.EventEnvelope[any]{
		SchemaVersion: coreeventscontract.Version,
		EventId:       uuid.New(),
		EventType:     eventType,
		AggregateType: coreeventscontract.AggregateType_DurableJobWorker,
		AggregateId:   hm.workerId,
		KafkaKey:      hm.workerId.String(),
		OccurredAt:    time.Now().UTC(),
		CorrelationId: correlationId,
		Data:          data,
	})
	if err != nil {
		return exceptions.New("FailedToPublishResult", "RoutineTask", "Finalize", "Failed to serialize the routine task result", http.StatusInternalServerError, true).WithOrigin(err)
	}
	if err := platformkafka.ProduceWithDefaultProducer(ctx, topic.String(), hm.workerId.String(), payload); err != nil {
		return exceptions.New("FailedToPublishResult", "RoutineTask", "Finalize", "Failed to publish the routine task result", http.StatusInternalServerError, true).WithOrigin(err)
	}
	return nil
}

/* ============================== Core logic ============================== */

func (hm *HandlerManager) Manage(
	ctx context.Context,
	claimedTasks []schemas.RoutineTask,
) *exceptions.Exception {
	if len(claimedTasks) == 0 {
		return nil
	}

	hm.resetRoutineTasksWithRecords(len(claimedTasks))

	taskIdToActorUserId := make(map[uuid.UUID]uuid.UUID, len(claimedTasks))
	groupsByPurpose := make(map[enums.RoutineTaskPurpose]purposeTaskGroup)
	for _, task := range claimedTasks {
		if task.ActorUserId == uuid.Nil {
			endedAt := time.Now()
			tempErrorCode := enums.RoutineTaskRecordErrorCode_PermissionDenied
			tempErrorReason := "Routine task actor was not found"
			hm.appendFailedRoutineTaskWithRecord(task, hm.newRecord(
				task,
				enums.RoutineTaskRecordStatus_Failed,
				endedAt,
				&tempErrorCode,
				&tempErrorReason,
			))
			continue
		}
		taskIdToActorUserId[task.Id] = task.ActorUserId

		registry, exists := hm.registries[task.Purpose]
		if !exists || registry.HandlerFunc == nil || len(registry.AllowedPermissions) == 0 {
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
		group.handlerFunc = registry.HandlerFunc
		group.allowedPermissions = registry.AllowedPermissions
		group.tasks = append(group.tasks, task)
		groupsByPurpose[task.Purpose] = group
	}
	if len(groupsByPurpose) == 0 {
		return hm.finalize(ctx)
	}

	for purpose, taskGroup := range groupsByPurpose {
		checkInputs := make([]inputs.BulkCheckRoutinePermissionInput, len(taskGroup.tasks))
		for index, task := range taskGroup.tasks {
			checkInputs[index] = inputs.BulkCheckRoutinePermissionInput{
				UserId: taskIdToActorUserId[task.Id],
				Id:     task.RoutineId,
			}
		}

		permissionSuccesses, _, exception := hm.routineRepository.BulkCheckPermissionsAndGetManyByIds(
			checkInputs,
			nil,
			taskGroup.allowedPermissions,
			options.WithDB(hm.db.WithContext(ctx)),
			options.WithAllowedPermissions(taskGroup.allowedPermissions),
		)
		if exception != nil {
			return exception
		}

		permittedTasks := make([]schemas.RoutineTask, 0, len(taskGroup.tasks))
		for index, task := range taskGroup.tasks {
			if permissionSuccesses[index] {
				permittedTasks = append(permittedTasks, task)
				continue
			}

			endedAt := time.Now()
			errorCode := enums.RoutineTaskRecordErrorCode_PermissionDenied
			errorReason := "Routine task actor no longer has the required routine permission"
			hm.appendFailedRoutineTaskWithRecord(task, hm.newRecord(
				task,
				enums.RoutineTaskRecordStatus_Failed,
				endedAt,
				&errorCode,
				&errorReason,
			))
		}
		if len(permittedTasks) == 0 {
			delete(groupsByPurpose, purpose)
			continue
		}

		taskGroup.tasks = permittedTasks
		groupsByPurpose[purpose] = taskGroup
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

			handlerResults, exception := group.handlerFunc(
				ctx,
				group.tasks,
				taskIdToActorUserId,
				group.allowedPermissions,
			)
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
