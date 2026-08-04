package workers

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"

	eventscontract "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/events"
	config "github.com/HiIamJeff67/notezy-backend/internal/platform/config"
	platformkafka "github.com/HiIamJeff67/notezy-backend/internal/platform/kafka"
	logs "github.com/HiIamJeff67/notezy-backend/internal/platform/observability/logs"
	realtimelease "github.com/HiIamJeff67/notezy-backend/internal/realtimegateway/data/cache/realtimelease"
)

type LifecycleConsumer struct {
	leaseStore *realtimelease.RealtimeLeaseCacheClient
}

func NewLifecycleConsumer(
	leaseStore *realtimelease.RealtimeLeaseCacheClient,
) *LifecycleConsumer {
	return &LifecycleConsumer{
		leaseStore: leaseStore,
	}
}

func (c *LifecycleConsumer) Start(ctx context.Context) func() {
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

func (c *LifecycleConsumer) run(ctx context.Context) {
	for ctx.Err() == nil {
		consumer, err := platformkafka.NewConsumer(
			config.Kafka(),
			eventscontract.CoreLifecycleTopic.String(),
		)
		if err == nil {
			err = consumer.Run(ctx, c.handle)
			consumer.Close()
		}
		if ctx.Err() != nil {
			return
		}
		if logs.NotezyLogger != nil {
			logs.NotezyLogger.Error(ctx, err, "RealtimeGateway lifecycle Kafka consumer stopped")
		}

		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Second):
		}
	}
}

func (c *LifecycleConsumer) handle(
	_ context.Context,
	_ platformkafka.ConsumerRecord,
	envelope eventscontract.EventEnvelope[json.RawMessage],
) error {
	switch envelope.EventType {
	case eventscontract.EventType_BlockPackAccessRevoked:
		var data eventscontract.BlockPackAccessRevokedData
		if err := json.Unmarshal(envelope.Data, &data); err != nil {
			return &platformkafka.ConsumerError{
				Classification: platformkafka.ErrorClassification_SchemaIncompatible,
				Origin:         err,
			}
		}
		if data.Reason != eventscontract.BlockPackAccessRevocationReason_PermissionRevoked &&
			data.Reason != eventscontract.BlockPackAccessRevocationReason_ResourceUnavailable {
			return &platformkafka.ConsumerError{
				Classification: platformkafka.ErrorClassification_SchemaIncompatible,
				Origin:         errors.New("Kafka BlockPack access revocation has an unsupported reason"),
			}
		}

		return c.leaseStore.PublishBlockPackChannelRevocation(realtimelease.BlockPackChannelRevocation{
			EventId:            envelope.EventId,
			BlockPackId:        envelope.AggregateId,
			TargetUserPublicId: data.TargetUserPublicId,
			Reason:             data.Reason,
		})
	case eventscontract.EventType_UserSessionsRevoked:
		return c.leaseStore.PublishUserSessionRevocation(realtimelease.UserSessionRevocation{
			EventId:      envelope.EventId,
			UserPublicId: envelope.AggregateId,
		})
	case eventscontract.EventType_BlockPackRoomPolicyChanged,
		eventscontract.EventType_RootShelfPermissionRevoked:
		return nil
	default:
		return &platformkafka.ConsumerError{
			Classification: platformkafka.ErrorClassification_PoisonMessage,
			Origin:         errors.New("Kafka lifecycle event type is unsupported by RealtimeGateway"),
		}
	}
}
