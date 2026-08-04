package routinetask

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	coreeventscontract "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/events"
	durablejobcontract "github.com/HiIamJeff67/notezy-backend/contracts/durablejob/v1"
	platformkafka "github.com/HiIamJeff67/notezy-backend/internal/platform/kafka"
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
	kafkaConfig    platformkafka.ConsumerConfig
	config         durablejobconfig.Config
}

func NewEngine(
	db *gorm.DB,
	kafkaConfig platformkafka.ConsumerConfig,
	config durablejobconfig.Config,
	maxWorkers ...int,
) *Engine {
	initialMaxWorkers := constants.RoutineTaskEngineMaxWorkers
	if len(maxWorkers) > 0 {
		initialMaxWorkers = min(initialMaxWorkers, maxWorkers[0])
	}

	engine := &Engine{
		ticker:      time.NewTicker(constants.RoutineTaskEngineTickerDuration),
		workerId:    uuid.New(),
		batchSize:   initialMaxWorkers,
		kafkaConfig: kafkaConfig,
		config:      config,
	}
	engine.handlerManager = NewHandlerManager(initialMaxWorkers, db, config, engine.workerId)
	engine.isHealthy.Store(true)

	return engine
}

func (e *Engine) requestRoutineTasks(ctx context.Context) {
	if e.isManagingWork.Load() {
		return
	}

	request := durablejobcontract.ClaimRoutineTasksRequestDto{
		RequestId: uuid.New(),
		WorkerId:  e.workerId,
		BatchSize: e.batchSize,
	}
	payload, err := json.Marshal(coreeventscontract.EventEnvelope[durablejobcontract.ClaimRoutineTasksRequestDto]{
		SchemaVersion: coreeventscontract.Version,
		EventId:       uuid.New(),
		EventType:     coreeventscontract.EventType_RoutineTaskClaimRequested,
		AggregateType: coreeventscontract.AggregateType_DurableJobWorker,
		AggregateId:   e.workerId,
		KafkaKey:      e.workerId.String(),
		OccurredAt:    time.Now().UTC(),
		CorrelationId: request.RequestId.String(),
		Data:          request,
	})
	if err != nil {
		e.isHealthy.Store(false)
		if logs.NotezyLogger != nil {
			logs.NotezyLogger.Error(ctx, err, "Failed to serialize routine task claim request")
		}

		return
	}
	if err := platformkafka.ProduceWithDefaultProducer(
		ctx,
		coreeventscontract.CoreDurableJobRoutineTaskTopic.String(),
		e.workerId.String(),
		payload,
	); err != nil {
		e.isHealthy.Store(false)
		if logs.NotezyLogger != nil {
			logs.NotezyLogger.Error(ctx, err, "Failed to publish routine task claim request")
		}

		return
	}

	e.isHealthy.Store(true)
}

func (e *Engine) consumeAssignment(
	ctx context.Context,
	_ platformkafka.ConsumerRecord,
	event coreeventscontract.EventEnvelope[json.RawMessage],
) error {
	if event.EventType != coreeventscontract.EventType_RoutineTasksAssigned {
		return nil
	}
	if event.AggregateType != coreeventscontract.AggregateType_DurableJobWorker {
		return &platformkafka.ConsumerError{
			Classification: platformkafka.ErrorClassification_SchemaIncompatible,
			Origin:         fmt.Errorf("DurableJob routine task assignment event is invalid"),
		}
	}

	var response durablejobcontract.ClaimRoutineTasksResponseDto
	if err := json.Unmarshal(event.Data, &response); err != nil {
		return &platformkafka.ConsumerError{
			Classification: platformkafka.ErrorClassification_SchemaIncompatible,
			Origin:         fmt.Errorf("decode DurableJob routine task assignments: %w", err),
		}
	}
	if response.RequestId == uuid.Nil || response.WorkerId != event.AggregateId {
		return &platformkafka.ConsumerError{
			Classification: platformkafka.ErrorClassification_SchemaIncompatible,
			Origin:         fmt.Errorf("DurableJob routine task assignment response is invalid"),
		}
	}
	if len(response.Assignments) == 0 {
		return nil
	}
	if !e.isManagingWork.CompareAndSwap(false, true) {
		return &platformkafka.ConsumerError{
			Classification: platformkafka.ErrorClassification_Transient,
			Origin:         fmt.Errorf("DurableJob routine task worker is at capacity"),
		}
	}
	defer e.isManagingWork.Store(false)

	routineTasks := make([]schemas.RoutineTask, len(response.Assignments))
	for index, assignment := range response.Assignments {
		if assignment.RoutineTaskId == uuid.Nil || assignment.RoutineTaskRecordId == uuid.Nil ||
			assignment.ActorUserId == uuid.Nil || assignment.RoutineId == uuid.Nil ||
			assignment.Purpose == "" || len(assignment.Payload) == 0 || assignment.StartedAt.IsZero() {
			return &platformkafka.ConsumerError{
				Classification: platformkafka.ErrorClassification_SchemaIncompatible,
				Origin:         fmt.Errorf("DurableJob routine task assignment at index %d is invalid", index),
			}
		}

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
		return &platformkafka.ConsumerError{
			Classification: platformkafka.ErrorClassification_Transient,
			Origin:         exception,
		}
	}

	e.isHealthy.Store(true)

	return nil
}

func (e *Engine) Start(ctx context.Context) func() {
	workerCtx, cancel := context.WithCancel(ctx)
	done := make(chan struct{})
	var shutdownOnce sync.Once

	consumer, err := platformkafka.NewConsumer(
		e.kafkaConfig,
		coreeventscontract.CoreDurableJobRoutineTaskTopic.String(),
	)
	if err != nil {
		e.isHealthy.Store(false)
		if logs.NotezyLogger != nil {
			logs.NotezyLogger.Error(workerCtx, err, "Failed to create routine task assignment consumer")
		}

		return func() {
			shutdownOnce.Do(func() {
				cancel()
				e.Stop()
			})
		}
	}
	go func() {
		if err := consumer.Run(workerCtx, e.consumeAssignment); err != nil && workerCtx.Err() == nil && logs.NotezyLogger != nil {
			logs.NotezyLogger.Error(workerCtx, err, "Routine task assignment consumer stopped")
		}
	}()

	go func() {
		defer close(done)
		defer e.Stop()

		e.requestRoutineTasks(workerCtx)
		for {
			select {
			case <-workerCtx.Done():
				return
			case <-e.ticker.C:
				e.requestRoutineTasks(workerCtx)
			}
		}
	}()

	return func() {
		shutdownOnce.Do(func() {
			cancel()
			if consumer != nil {
				consumer.Close()
			}
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
