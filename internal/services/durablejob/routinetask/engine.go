package routinetask

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	durablejobcontract "github.com/HiIamJeff67/notezy-backend/contracts/durablejob/v1"
	durablejobroutetasktypes "github.com/HiIamJeff67/notezy-backend/contracts/durablejob/v1/types/routine-tasks"
	logs "github.com/HiIamJeff67/notezy-backend/internal/platform/observability/logs"
	schemas "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/database/schemas"
	enums "github.com/HiIamJeff67/notezy-backend/internal/services/core/data/database/schemas/enums"
	durablejobconfig "github.com/HiIamJeff67/notezy-backend/internal/services/durablejob/config"
	constants "github.com/HiIamJeff67/notezy-backend/shared/constants"
)

type Engine struct {
	ticker         *time.Ticker
	workerId       uuid.UUID
	batchSize      int
	isHealthy      atomic.Bool
	isManagingWork atomic.Bool
	handlerManager HandlerManager
}

func NewEngine(
	db *gorm.DB,
	config durablejobconfig.Config,
	maxWorkers ...int,
) *Engine {
	initialMaxWorkers := constants.RoutineTaskEngineMaxWorkers
	if len(maxWorkers) > 0 {
		initialMaxWorkers = min(initialMaxWorkers, maxWorkers[0])
	}

	engine := &Engine{
		ticker:    time.NewTicker(constants.RoutineTaskEngineTickerDuration),
		workerId:  uuid.New(),
		batchSize: initialMaxWorkers,
	}
	engine.handlerManager = NewHandlerManager(initialMaxWorkers, db, config, engine.workerId)
	engine.isHealthy.Store(true)

	return engine
}

func (e *Engine) SetResultPublisher(publisher ResultPublisher) {
	e.handlerManager.SetResultPublisher(publisher)
}

func (e *Engine) NewClaimRoutineTasksRequest() (durablejobcontract.ClaimRoutineTasksRequestDto, bool) {
	if e.isManagingWork.Load() {
		return durablejobcontract.ClaimRoutineTasksRequestDto{}, false
	}

	return durablejobcontract.ClaimRoutineTasksRequestDto{
		RequestId: uuid.New(),
		WorkerId:  e.workerId,
		BatchSize: e.batchSize,
	}, true
}

func (e *Engine) HandleRoutineTaskAssignments(
	ctx context.Context,
	assignments []durablejobroutetasktypes.RoutineTaskAssignment,
) error {
	if len(assignments) == 0 {
		return nil
	}
	if !e.isManagingWork.CompareAndSwap(false, true) {
		return errors.New("DurableJob routine task worker is at capacity")
	}
	defer e.isManagingWork.Store(false)

	routineTasks := make([]schemas.RoutineTask, len(assignments))
	for index, assignment := range assignments {
		startedAt := assignment.StartedAt
		routineTasks[index] = schemas.RoutineTask{
			Id:                assignment.RoutineTaskId,
			RoutineId:         assignment.RoutineId,
			ActorUserId:       assignment.ActorUserId,
			Title:             assignment.Title,
			Purpose:           enums.RoutineTaskPurpose(assignment.Purpose),
			Payload:           datatypes.JSON(assignment.Payload),
			CostUnit:          assignment.CostUnit,
			Priority:          assignment.Priority,
			Status:            enums.RoutineTaskStatus_Running,
			Attempts:          assignment.Attempt,
			ActualStartedAt:   &startedAt,
			RecordScheduledAt: assignment.ScheduledAt,
			RecordId:          assignment.RoutineTaskRecordId,
		}
	}

	if exception := e.handlerManager.Manage(ctx, routineTasks); exception != nil {
		return exception
	}

	e.isHealthy.Store(true)

	return nil
}

func (e *Engine) Start(
	ctx context.Context,
	requestRoutineTasks func(context.Context, durablejobcontract.ClaimRoutineTasksRequestDto) error,
) func() {
	workerCtx, cancel := context.WithCancel(ctx)
	done := make(chan struct{})
	var shutdownOnce sync.Once

	go func() {
		defer close(done)
		defer e.Stop()

		request, shouldRequest := e.NewClaimRoutineTasksRequest()
		if shouldRequest {
			if err := requestRoutineTasks(workerCtx, request); err != nil {
				e.isHealthy.Store(false)
				if logs.NotezyLogger != nil {
					logs.NotezyLogger.Error(workerCtx, err, "Failed to publish routine task claim request")
				}
			}
		}
		for {
			select {
			case <-workerCtx.Done():
				return
			case <-e.ticker.C:
				request, shouldRequest := e.NewClaimRoutineTasksRequest()
				if !shouldRequest {
					continue
				}
				if err := requestRoutineTasks(workerCtx, request); err != nil {
					e.isHealthy.Store(false)
					if logs.NotezyLogger != nil {
						logs.NotezyLogger.Error(workerCtx, err, "Failed to publish routine task claim request")
					}
					continue
				}
				e.isHealthy.Store(true)
			}
		}
	}()

	return func() {
		shutdownOnce.Do(func() {
			cancel()
			<-done
		})
	}
}

func (e *Engine) Stop() {
	if e.ticker != nil {
		e.ticker.Stop()
	}
	e.isHealthy.Store(false)
}

func (e *Engine) IsHealthy() bool {
	return e.isHealthy.Load()
}
