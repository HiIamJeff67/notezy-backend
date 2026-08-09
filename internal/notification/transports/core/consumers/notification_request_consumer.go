package consumers

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"

	coreeventscontract "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/events"
	eventcontract "github.com/HiIamJeff67/notezy-backend/contracts/types/events"

	platformkafka "github.com/HiIamJeff67/notezy-backend/shared/platform/kafka"
	logs "github.com/HiIamJeff67/notezy-backend/shared/platform/observability/logs"

	services "github.com/HiIamJeff67/notezy-backend/internal/notification/services"
)

type NotificationRequestConsumer struct {
	service     *services.NotificationService
	kafkaConfig platformkafka.ConsumerConfig
}

func NewNotificationRequestConsumer(
	service *services.NotificationService,
	kafkaConfig platformkafka.ConsumerConfig,
) *NotificationRequestConsumer {
	return &NotificationRequestConsumer{
		service:     service,
		kafkaConfig: kafkaConfig,
	}
}

func (c *NotificationRequestConsumer) Start(ctx context.Context) func() {
	workerCtx, cancel := context.WithCancel(ctx)
	var waitGroup sync.WaitGroup
	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		c.run(workerCtx)
	}()

	return func() {
		cancel()
		waitGroup.Wait()
	}
}

func (c *NotificationRequestConsumer) run(ctx context.Context) {
	for ctx.Err() == nil {
		consumer, err := platformkafka.NewConsumer(
			c.kafkaConfig,
			coreeventscontract.CoreNotificationTopic.String(),
			coreeventscontract.CoreLifecycleTopic.String(),
		)
		if err == nil {
			err = consumer.Run(ctx, c.consume)
			consumer.Close()
		}
		if ctx.Err() != nil {
			return
		}
		if logs.NotezyLogger != nil {
			logs.NotezyLogger.Error(ctx, err, "Notification request consumer stopped")
		}
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Second):
		}
	}
}

func (c *NotificationRequestConsumer) consume(
	ctx context.Context,
	_ platformkafka.ConsumerRecord,
	event eventcontract.EventEnvelope[json.RawMessage],
) error {
	if event.EventType == coreeventscontract.EventType_UserDeleted {
		if event.AggregateType != coreeventscontract.AggregateType_User || event.AggregateId == uuid.Nil {
			return &platformkafka.ConsumerError{
				Classification: platformkafka.ErrorClassification_SchemaIncompatible,
				Origin:         errors.New("user deletion event aggregate is invalid"),
			}
		}
		var data coreeventscontract.UserDeletedData
		if err := json.Unmarshal(event.Data, &data); err != nil {
			return &platformkafka.ConsumerError{
				Classification: platformkafka.ErrorClassification_SchemaIncompatible,
				Origin:         err,
			}
		}
		if data.DeletedAt.IsZero() {
			return &platformkafka.ConsumerError{
				Classification: platformkafka.ErrorClassification_SchemaIncompatible,
				Origin:         errors.New("user deletion event timestamp is missing"),
			}
		}
		if err := c.service.DeleteForUser(ctx, event.AggregateId); err != nil {
			return &platformkafka.ConsumerError{
				Classification: platformkafka.ErrorClassification_Transient,
				Origin:         err,
			}
		}
		return nil
	}
	if event.EventType != coreeventscontract.EventType_NotificationRequested {
		return nil
	}
	var data coreeventscontract.NotificationRequestedData
	if err := json.Unmarshal(event.Data, &data); err != nil {
		return &platformkafka.ConsumerError{
			Classification: platformkafka.ErrorClassification_SchemaIncompatible,
			Origin:         err,
		}
	}
	eventWithData := eventcontract.EventEnvelope[coreeventscontract.NotificationRequestedData]{
		SchemaVersion: event.SchemaVersion,
		EventId:       event.EventId,
		EventType:     event.EventType,
		AggregateType: event.AggregateType,
		AggregateId:   event.AggregateId,
		KafkaKey:      event.KafkaKey,
		OccurredAt:    event.OccurredAt,
		CorrelationId: event.CorrelationId,
		CausationId:   event.CausationId,
		Trace:         event.Trace,
		Data:          data,
	}
	if err := c.service.ConsumeRequested(ctx, eventWithData); err != nil {
		return &platformkafka.ConsumerError{
			Classification: platformkafka.ErrorClassification_Transient,
			Origin:         err,
		}
	}

	return nil
}
