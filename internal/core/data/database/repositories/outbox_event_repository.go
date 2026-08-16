package repositories

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	exceptions "github.com/HiIamJeff67/notegic-backend/contracts/types/exceptions"

	coreeventscontract "github.com/HiIamJeff67/notegic-backend/contracts/core/v1/events"
	eventcontract "github.com/HiIamJeff67/notegic-backend/contracts/types/events"

	inputs "github.com/HiIamJeff67/notegic-backend/internal/core/data/database/inputs"
	options "github.com/HiIamJeff67/notegic-backend/internal/core/data/database/options"
	schemas "github.com/HiIamJeff67/notegic-backend/internal/core/data/database/schemas"
	durablejobeventbuilders "github.com/HiIamJeff67/notegic-backend/internal/core/transports/durablejob/eventbuilders"
)

type OutboxEventRepositoryInterface interface {
	CreateMany(createInputs []inputs.CreateOutboxEventInput, opts ...options.RepositoryOptions) *exceptions.Exception
	EnqueueBlockPackAccessRevocations(tx *gorm.DB, correlationId string, blockPackIds []uuid.UUID, targetUserPublicIds []uuid.UUID, reason coreeventscontract.BlockPackAccessRevocationReason) error
	EnqueueRootShelfPermissionChanged(tx *gorm.DB, correlationId string, rootShelfId uuid.UUID, targetUserPublicId uuid.UUID, permission string) error
	EnqueueManyRootShelfPermissionChanges(tx *gorm.DB, correlationId string, rootShelfId uuid.UUID, permissions []schemas.UsersToShelves, userPublicIdByUserId map[uuid.UUID]uuid.UUID) error
	EnqueueRootShelfPermissionRevoked(tx *gorm.DB, correlationId string, rootShelfId uuid.UUID, targetUserPublicId uuid.UUID) error
	EnqueueManyRootShelfPermissionRevocations(tx *gorm.DB, correlationId string, rootShelfIds []uuid.UUID, targetUserPublicIds []uuid.UUID) error
	EnqueueRootShelfDeleted(tx *gorm.DB, correlationId string, rootShelfId uuid.UUID, targetUserPublicIds []uuid.UUID) error
	EnqueueManyRootShelfDeleted(tx *gorm.DB, correlationId string, rootShelfIds []uuid.UUID, targetUserPublicIdsByRootShelfId map[uuid.UUID][]uuid.UUID) error
	EnqueueBlockPackChanged(tx *gorm.DB, correlationId string, blockPackIds []uuid.UUID) error
	EnqueueBlockPackDeleted(tx *gorm.DB, correlationId string, blockPackIds []uuid.UUID) error
	EnqueueUserSessionsRevoked(tx *gorm.DB, correlationId string, userPublicId uuid.UUID) error
	EnqueueUserDeleted(tx *gorm.DB, correlationId string, userPublicId uuid.UUID, deletedAt time.Time) error
	EnqueueNotificationRequested(tx *gorm.DB, correlationId string, data coreeventscontract.NotificationRequestedData) error
	EnqueueYjsMaintenanceHint(tx *gorm.DB, correlationId string, blockPackId uuid.UUID, reason string) error
	EnqueueManyYjsMaintenanceHints(tx *gorm.DB, correlationId string, blockPackIds []uuid.UUID, reason string) error
	ClaimAvailable(ctx context.Context, workerId string, batchSize int, claimTimeout time.Duration, opts ...options.RepositoryOptions) ([]schemas.OutboxEvent, *exceptions.Exception)
	MarkPublishedMany(ctx context.Context, eventIds []uuid.UUID, workerId string, opts ...options.RepositoryOptions) *exceptions.Exception
	MarkFailedMany(ctx context.Context, failureInputs []inputs.FailedOutboxEventInput, workerId string, opts ...options.RepositoryOptions) *exceptions.Exception
	DeletePublishedBefore(ctx context.Context, publishedBefore time.Time, opts ...options.RepositoryOptions) (int64, *exceptions.Exception)
}

type OutboxEventRepository struct{}

type outboxEventMetadata struct {
	SchemaVersion string                      `json:"schemaVersion"`
	CorrelationId string                      `json:"correlationId"`
	CausationId   *uuid.UUID                  `json:"causationId,omitempty"`
	OccurredAt    time.Time                   `json:"occurredAt"`
	Trace         eventcontract.TraceMetadata `json:"trace"`
}

func NewOutboxEventRepository() OutboxEventRepositoryInterface {
	return &OutboxEventRepository{}
}

func ConvertEnvelopeToCreateOutboxEventInput[D any](
	topic eventcontract.Topic,
	envelope eventcontract.EventEnvelope[D],
) (inputs.CreateOutboxEventInput, error) {
	if topic == "" || envelope.EventId == uuid.Nil || envelope.AggregateId == uuid.Nil ||
		envelope.AggregateType == "" || envelope.EventType == "" || envelope.KafkaKey == "" {
		return inputs.CreateOutboxEventInput{}, errors.New("outbox event envelope is incomplete")
	}
	if envelope.KafkaKey != envelope.AggregateId.String() {
		return inputs.CreateOutboxEventInput{}, errors.New("outbox event Kafka key must equal the aggregate ID")
	}

	payload, err := json.Marshal(envelope.Data)
	if err != nil {
		return inputs.CreateOutboxEventInput{}, err
	}
	metadata, err := json.Marshal(outboxEventMetadata{
		SchemaVersion: envelope.SchemaVersion,
		CorrelationId: envelope.CorrelationId,
		CausationId:   envelope.CausationId,
		OccurredAt:    envelope.OccurredAt,
		Trace:         envelope.Trace,
	})
	if err != nil {
		return inputs.CreateOutboxEventInput{}, err
	}

	return inputs.CreateOutboxEventInput{
		Id:            envelope.EventId,
		AggregateType: envelope.AggregateType,
		AggregateId:   envelope.AggregateId,
		EventType:     envelope.EventType,
		Topic:         topic,
		KafkaKey:      envelope.KafkaKey,
		Payload:       payload,
		Metadata:      metadata,
		AvailableAt:   time.Now(),
	}, nil
}

func EnqueueOutboxEvents[D any](
	tx *gorm.DB,
	topic eventcontract.Topic,
	envelopes []eventcontract.EventEnvelope[D],
) error {
	if len(envelopes) == 0 {
		return nil
	}

	createInputs := make([]inputs.CreateOutboxEventInput, len(envelopes))
	for index, envelope := range envelopes {
		createInput, err := ConvertEnvelopeToCreateOutboxEventInput(topic, envelope)
		if err != nil {
			return err
		}
		createInputs[index] = createInput
	}

	exception := NewOutboxEventRepository().CreateMany(
		createInputs,
		options.WithTransactionDB(tx),
	)
	if exception != nil {
		return exception
	}

	return nil
}

func (r *OutboxEventRepository) CreateMany(
	createInputs []inputs.CreateOutboxEventInput,
	opts ...options.RepositoryOptions,
) *exceptions.Exception {
	if len(createInputs) == 0 {
		return nil
	}

	parsedOptions := options.ParseRepositoryOptions(opts...)
	if !parsedOptions.IsTransactionStarted {
		return exceptions.New(
			"TransactionRequired",
			"Outbox",
			"Create",
			"Outbox events must be created in the domain transaction",
			http.StatusInternalServerError,
		)
	}

	events := make([]schemas.OutboxEvent, len(createInputs))
	for index, createInput := range createInputs {
		events[index] = schemas.OutboxEvent{
			Id:            createInput.Id,
			AggregateType: createInput.AggregateType,
			AggregateId:   createInput.AggregateId,
			EventType:     createInput.EventType,
			Topic:         createInput.Topic,
			KafkaKey:      createInput.KafkaKey,
			Payload:       datatypes.JSON(createInput.Payload),
			Metadata:      datatypes.JSON(createInput.Metadata),
			AvailableAt:   createInput.AvailableAt,
		}
	}

	result := parsedOptions.DB.CreateInBatches(&events, parsedOptions.BatchSize)
	if result.Error != nil {
		return exceptions.New(
			"FailedToCreate",
			"Outbox",
			"Create",
			"Failed to create outbox events",
			http.StatusInternalServerError,
			true,
		).WithOrigin(result.Error)
	}

	return nil
}

func (r *OutboxEventRepository) EnqueueBlockPackAccessRevocations(
	tx *gorm.DB,
	correlationId string,
	blockPackIds []uuid.UUID,
	targetUserPublicIds []uuid.UUID,
	reason coreeventscontract.BlockPackAccessRevocationReason,
) error {
	if len(blockPackIds) == 0 {
		return nil
	}

	targetCount := len(targetUserPublicIds)
	if targetCount == 0 {
		targetCount = 1
	}
	events := make(
		[]eventcontract.EventEnvelope[coreeventscontract.BlockPackAccessRevokedData],
		0,
		len(blockPackIds)*targetCount,
	)
	occurredAt := time.Now().UTC()
	for _, blockPackId := range blockPackIds {
		if len(targetUserPublicIds) == 0 {
			events = append(events, eventcontract.EventEnvelope[coreeventscontract.BlockPackAccessRevokedData]{
				SchemaVersion: eventcontract.Version,
				EventId:       uuid.New(),
				EventType:     coreeventscontract.EventType_BlockPackAccessRevoked,
				AggregateType: coreeventscontract.AggregateType_BlockPack,
				AggregateId:   blockPackId,
				KafkaKey:      blockPackId.String(),
				OccurredAt:    occurredAt,
				CorrelationId: correlationId,
				Data: coreeventscontract.BlockPackAccessRevokedData{
					Reason: reason,
				},
			})
			continue
		}

		for _, targetUserPublicId := range targetUserPublicIds {
			targetUserPublicId := targetUserPublicId
			events = append(events, eventcontract.EventEnvelope[coreeventscontract.BlockPackAccessRevokedData]{
				SchemaVersion: eventcontract.Version,
				EventId:       uuid.New(),
				EventType:     coreeventscontract.EventType_BlockPackAccessRevoked,
				AggregateType: coreeventscontract.AggregateType_BlockPack,
				AggregateId:   blockPackId,
				KafkaKey:      blockPackId.String(),
				OccurredAt:    occurredAt,
				CorrelationId: correlationId,
				Data: coreeventscontract.BlockPackAccessRevokedData{
					TargetUserPublicId: &targetUserPublicId,
					Reason:             reason,
				},
			})
		}
	}

	return EnqueueOutboxEvents(tx, coreeventscontract.CoreLifecycleTopic, events)
}

func (r *OutboxEventRepository) EnqueueRootShelfPermissionChanged(
	tx *gorm.DB,
	correlationId string,
	rootShelfId uuid.UUID,
	targetUserPublicId uuid.UUID,
	permission string,
) error {
	return EnqueueOutboxEvents(
		tx,
		coreeventscontract.CoreLifecycleTopic,
		[]eventcontract.EventEnvelope[coreeventscontract.ResourceChangedData]{
			{
				SchemaVersion: eventcontract.Version,
				EventId:       uuid.New(),
				EventType:     coreeventscontract.EventType_RootShelfPermissionChanged,
				AggregateType: coreeventscontract.AggregateType_RootShelf,
				AggregateId:   rootShelfId,
				KafkaKey:      rootShelfId.String(),
				OccurredAt:    time.Now().UTC(),
				CorrelationId: correlationId,
				Data: coreeventscontract.ResourceChangedData{
					ResourceId:         rootShelfId,
					TargetUserPublicId: &targetUserPublicId,
					Change:             coreeventscontract.ResourceEventChange_PermissionUpdated,
					Permission:         permission,
				},
			},
		},
	)
}

func (r *OutboxEventRepository) EnqueueManyRootShelfPermissionChanges(
	tx *gorm.DB,
	correlationId string,
	rootShelfId uuid.UUID,
	permissions []schemas.UsersToShelves,
	userPublicIdByUserId map[uuid.UUID]uuid.UUID,
) error {
	events := make([]eventcontract.EventEnvelope[coreeventscontract.ResourceChangedData], 0, len(permissions))
	occurredAt := time.Now().UTC()
	for _, permission := range permissions {
		targetUserPublicId, exists := userPublicIdByUserId[permission.UserId]
		if !exists {
			return errors.New("root shelf permission event target user is unavailable")
		}

		events = append(events, eventcontract.EventEnvelope[coreeventscontract.ResourceChangedData]{
			SchemaVersion: eventcontract.Version,
			EventId:       uuid.New(),
			EventType:     coreeventscontract.EventType_RootShelfPermissionChanged,
			AggregateType: coreeventscontract.AggregateType_RootShelf,
			AggregateId:   rootShelfId,
			KafkaKey:      rootShelfId.String(),
			OccurredAt:    occurredAt,
			CorrelationId: correlationId,
			Data: coreeventscontract.ResourceChangedData{
				ResourceId:         rootShelfId,
				TargetUserPublicId: &targetUserPublicId,
				Change:             coreeventscontract.ResourceEventChange_PermissionUpdated,
				Permission:         permission.Permission.String(),
			},
		})
	}

	return EnqueueOutboxEvents(tx, coreeventscontract.CoreLifecycleTopic, events)
}

func (r *OutboxEventRepository) EnqueueRootShelfPermissionRevoked(
	tx *gorm.DB,
	correlationId string,
	rootShelfId uuid.UUID,
	targetUserPublicId uuid.UUID,
) error {
	return EnqueueOutboxEvents(
		tx,
		coreeventscontract.CoreLifecycleTopic,
		[]eventcontract.EventEnvelope[coreeventscontract.ResourceChangedData]{
			{
				SchemaVersion: eventcontract.Version,
				EventId:       uuid.New(),
				EventType:     coreeventscontract.EventType_RootShelfPermissionRevoked,
				AggregateType: coreeventscontract.AggregateType_RootShelf,
				AggregateId:   rootShelfId,
				KafkaKey:      rootShelfId.String(),
				OccurredAt:    time.Now().UTC(),
				CorrelationId: correlationId,
				Data: coreeventscontract.ResourceChangedData{
					ResourceId:         rootShelfId,
					TargetUserPublicId: &targetUserPublicId,
					Change:             coreeventscontract.ResourceEventChange_PermissionRevoked,
				},
			},
		},
	)
}

func (r *OutboxEventRepository) EnqueueManyRootShelfPermissionRevocations(
	tx *gorm.DB,
	correlationId string,
	rootShelfIds []uuid.UUID,
	targetUserPublicIds []uuid.UUID,
) error {
	events := make([]eventcontract.EventEnvelope[coreeventscontract.ResourceChangedData], 0, len(rootShelfIds)*len(targetUserPublicIds))
	occurredAt := time.Now().UTC()
	for _, rootShelfId := range rootShelfIds {
		for _, targetUserPublicId := range targetUserPublicIds {
			targetUserPublicId := targetUserPublicId
			events = append(events, eventcontract.EventEnvelope[coreeventscontract.ResourceChangedData]{
				SchemaVersion: eventcontract.Version,
				EventId:       uuid.New(),
				EventType:     coreeventscontract.EventType_RootShelfPermissionRevoked,
				AggregateType: coreeventscontract.AggregateType_RootShelf,
				AggregateId:   rootShelfId,
				KafkaKey:      rootShelfId.String(),
				OccurredAt:    occurredAt,
				CorrelationId: correlationId,
				Data: coreeventscontract.ResourceChangedData{
					ResourceId:         rootShelfId,
					TargetUserPublicId: &targetUserPublicId,
					Change:             coreeventscontract.ResourceEventChange_PermissionRevoked,
				},
			})
		}
	}

	return EnqueueOutboxEvents(tx, coreeventscontract.CoreLifecycleTopic, events)
}

func (r *OutboxEventRepository) EnqueueRootShelfDeleted(
	tx *gorm.DB,
	correlationId string,
	rootShelfId uuid.UUID,
	targetUserPublicIds []uuid.UUID,
) error {
	events := make([]eventcontract.EventEnvelope[coreeventscontract.ResourceChangedData], 0, len(targetUserPublicIds))
	for _, targetUserPublicId := range targetUserPublicIds {
		targetUserPublicId := targetUserPublicId
		events = append(events, eventcontract.EventEnvelope[coreeventscontract.ResourceChangedData]{
			SchemaVersion: eventcontract.Version,
			EventId:       uuid.New(),
			EventType:     coreeventscontract.EventType_RootShelfDeleted,
			AggregateType: coreeventscontract.AggregateType_RootShelf,
			AggregateId:   rootShelfId,
			KafkaKey:      rootShelfId.String(),
			OccurredAt:    time.Now().UTC(),
			CorrelationId: correlationId,
			Data: coreeventscontract.ResourceChangedData{
				ResourceId:         rootShelfId,
				TargetUserPublicId: &targetUserPublicId,
				Change:             coreeventscontract.ResourceEventChange_Deleted,
			},
		})
	}

	return EnqueueOutboxEvents(tx, coreeventscontract.CoreLifecycleTopic, events)
}

func (r *OutboxEventRepository) EnqueueManyRootShelfDeleted(
	tx *gorm.DB,
	correlationId string,
	rootShelfIds []uuid.UUID,
	targetUserPublicIdsByRootShelfId map[uuid.UUID][]uuid.UUID,
) error {
	events := make([]eventcontract.EventEnvelope[coreeventscontract.ResourceChangedData], 0)
	occurredAt := time.Now().UTC()
	for _, rootShelfId := range rootShelfIds {
		for _, targetUserPublicId := range targetUserPublicIdsByRootShelfId[rootShelfId] {
			targetUserPublicId := targetUserPublicId
			events = append(events, eventcontract.EventEnvelope[coreeventscontract.ResourceChangedData]{
				SchemaVersion: eventcontract.Version,
				EventId:       uuid.New(),
				EventType:     coreeventscontract.EventType_RootShelfDeleted,
				AggregateType: coreeventscontract.AggregateType_RootShelf,
				AggregateId:   rootShelfId,
				KafkaKey:      rootShelfId.String(),
				OccurredAt:    occurredAt,
				CorrelationId: correlationId,
				Data: coreeventscontract.ResourceChangedData{
					ResourceId:         rootShelfId,
					TargetUserPublicId: &targetUserPublicId,
					Change:             coreeventscontract.ResourceEventChange_Deleted,
				},
			})
		}
	}

	return EnqueueOutboxEvents(tx, coreeventscontract.CoreLifecycleTopic, events)
}

func (r *OutboxEventRepository) EnqueueBlockPackChanged(
	tx *gorm.DB,
	correlationId string,
	blockPackIds []uuid.UUID,
) error {
	events := make([]eventcontract.EventEnvelope[coreeventscontract.ResourceChangedData], 0, len(blockPackIds))
	for _, blockPackId := range blockPackIds {
		events = append(events, eventcontract.EventEnvelope[coreeventscontract.ResourceChangedData]{
			SchemaVersion: eventcontract.Version,
			EventId:       uuid.New(),
			EventType:     coreeventscontract.EventType_BlockPackChanged,
			AggregateType: coreeventscontract.AggregateType_BlockPack,
			AggregateId:   blockPackId,
			KafkaKey:      blockPackId.String(),
			OccurredAt:    time.Now().UTC(),
			CorrelationId: correlationId,
			Data: coreeventscontract.ResourceChangedData{
				ResourceId: blockPackId,
				Change:     coreeventscontract.ResourceEventChange_Updated,
			},
		})
	}

	return EnqueueOutboxEvents(tx, coreeventscontract.CoreLifecycleTopic, events)
}

func (r *OutboxEventRepository) EnqueueBlockPackDeleted(
	tx *gorm.DB,
	correlationId string,
	blockPackIds []uuid.UUID,
) error {
	events := make([]eventcontract.EventEnvelope[coreeventscontract.ResourceChangedData], 0, len(blockPackIds))
	for _, blockPackId := range blockPackIds {
		events = append(events, eventcontract.EventEnvelope[coreeventscontract.ResourceChangedData]{
			SchemaVersion: eventcontract.Version,
			EventId:       uuid.New(),
			EventType:     coreeventscontract.EventType_BlockPackDeleted,
			AggregateType: coreeventscontract.AggregateType_BlockPack,
			AggregateId:   blockPackId,
			KafkaKey:      blockPackId.String(),
			OccurredAt:    time.Now().UTC(),
			CorrelationId: correlationId,
			Data: coreeventscontract.ResourceChangedData{
				ResourceId: blockPackId,
				Change:     coreeventscontract.ResourceEventChange_Deleted,
			},
		})
	}

	return EnqueueOutboxEvents(tx, coreeventscontract.CoreLifecycleTopic, events)
}

func (r *OutboxEventRepository) EnqueueUserSessionsRevoked(
	tx *gorm.DB,
	correlationId string,
	userPublicId uuid.UUID,
) error {
	return EnqueueOutboxEvents(
		tx,
		coreeventscontract.CoreLifecycleTopic,
		[]eventcontract.EventEnvelope[coreeventscontract.UserSessionsRevokedData]{
			{
				SchemaVersion: eventcontract.Version,
				EventId:       uuid.New(),
				EventType:     coreeventscontract.EventType_UserSessionsRevoked,
				AggregateType: coreeventscontract.AggregateType_User,
				AggregateId:   userPublicId,
				KafkaKey:      userPublicId.String(),
				OccurredAt:    time.Now().UTC(),
				CorrelationId: correlationId,
				Data:          coreeventscontract.UserSessionsRevokedData{},
			},
		},
	)
}

func (r *OutboxEventRepository) EnqueueUserDeleted(
	tx *gorm.DB,
	correlationId string,
	userPublicId uuid.UUID,
	deletedAt time.Time,
) error {
	if tx == nil || userPublicId == uuid.Nil || deletedAt.IsZero() {
		return errors.New("user deletion event requires a transaction, user public ID, and deletion time")
	}

	return EnqueueOutboxEvents(
		tx,
		coreeventscontract.CoreLifecycleTopic,
		[]eventcontract.EventEnvelope[coreeventscontract.UserDeletedData]{
			{
				SchemaVersion: eventcontract.Version,
				EventId:       uuid.New(),
				EventType:     coreeventscontract.EventType_UserDeleted,
				AggregateType: coreeventscontract.AggregateType_User,
				AggregateId:   userPublicId,
				KafkaKey:      userPublicId.String(),
				OccurredAt:    deletedAt.UTC(),
				CorrelationId: correlationId,
				Data: coreeventscontract.UserDeletedData{
					DeletedAt: deletedAt.UTC(),
				},
			},
		},
	)
}

func (r *OutboxEventRepository) EnqueueNotificationRequested(
	tx *gorm.DB,
	correlationId string,
	data coreeventscontract.NotificationRequestedData,
) error {
	if tx == nil || data.RecipientUserPublicId == uuid.Nil || data.Type == "" ||
		data.TemplateKey == "" || data.TemplateVersion <= 0 || data.DedupeKey == "" {
		return errors.New("notification request is incomplete")
	}

	envelope := eventcontract.EventEnvelope[coreeventscontract.NotificationRequestedData]{
		SchemaVersion: eventcontract.Version,
		EventId:       uuid.New(),
		EventType:     coreeventscontract.EventType_NotificationRequested,
		AggregateType: coreeventscontract.AggregateType_Notification,
		AggregateId:   data.RecipientUserPublicId,
		KafkaKey:      data.RecipientUserPublicId.String(),
		OccurredAt:    time.Now().UTC(),
		CorrelationId: correlationId,
		Data:          data,
	}

	return EnqueueOutboxEvents(
		tx,
		coreeventscontract.CoreNotificationTopic,
		[]eventcontract.EventEnvelope[coreeventscontract.NotificationRequestedData]{envelope},
	)
}

func (r *OutboxEventRepository) EnqueueYjsMaintenanceHint(
	tx *gorm.DB,
	correlationId string,
	blockPackId uuid.UUID,
	reason string,
) error {
	if tx == nil || blockPackId == uuid.Nil {
		return errors.New("Yjs maintenance hint requires a transaction and BlockPack ID")
	}

	return r.EnqueueManyYjsMaintenanceHints(tx, correlationId, []uuid.UUID{blockPackId}, reason)
}

func (r *OutboxEventRepository) EnqueueManyYjsMaintenanceHints(
	tx *gorm.DB,
	correlationId string,
	blockPackIds []uuid.UUID,
	reason string,
) error {
	if tx == nil {
		return errors.New("Yjs maintenance hints require a transaction")
	}
	if len(blockPackIds) == 0 {
		return nil
	}

	var documents []schemas.BlockPackYjsDocument
	if err := tx.Model(&schemas.BlockPackYjsDocument{}).
		Where("block_pack_id IN ? AND deleted_at IS NULL", blockPackIds).
		Find(&documents).Error; err != nil {
		return err
	}
	if len(documents) != len(blockPackIds) {
		return errors.New("Yjs maintenance hints require documents for every BlockPack ID")
	}

	eventBuilder := durablejobeventbuilders.NewYjsMaintenanceHintEventBuilder()
	occurredAt := time.Now().UTC()
	events := make([]eventcontract.EventEnvelope[coreeventscontract.YjsMaintenanceHintData], 0, len(documents))
	for _, document := range documents {
		events = append(events, eventBuilder.Build(coreeventscontract.YjsMaintenanceHintData{
			BlockPackId:            document.BlockPackId,
			DocumentId:             document.Id,
			LatestUpdateSequence:   document.LastUpdateSequence,
			CompactedUntilSequence: document.CompactedUntilSequence,
			ProjectedUntilSequence: document.ProjectedUntilSequence,
			LastCompactedAt:        document.LastCompactedAt,
			UncompactedUpdateCount: document.LastUpdateSequence - document.CompactedUntilSequence,
			SnapshotBytes:          len(document.Snapshot),
			StateVectorBytes:       len(document.StateVector),
			Reason:                 reason,
		}, correlationId, occurredAt))
	}

	return EnqueueOutboxEvents(tx, coreeventscontract.CoreDurableJobYjsMaintenanceHintTopic, events)
}

func SerializeOutboxEvent(event schemas.OutboxEvent) ([]byte, error) {
	var metadata outboxEventMetadata
	if err := json.Unmarshal(event.Metadata, &metadata); err != nil {
		return nil, err
	}

	var payload json.RawMessage
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		return nil, err
	}

	return json.Marshal(eventcontract.EventEnvelope[json.RawMessage]{
		SchemaVersion: metadata.SchemaVersion,
		EventId:       event.Id,
		EventType:     event.EventType,
		AggregateType: event.AggregateType,
		AggregateId:   event.AggregateId,
		KafkaKey:      event.KafkaKey,
		OccurredAt:    metadata.OccurredAt,
		CorrelationId: metadata.CorrelationId,
		CausationId:   metadata.CausationId,
		Trace:         metadata.Trace,
		Data:          payload,
	})
}

func (r *OutboxEventRepository) ClaimAvailable(
	ctx context.Context,
	workerId string,
	batchSize int,
	claimTimeout time.Duration,
	opts ...options.RepositoryOptions,
) ([]schemas.OutboxEvent, *exceptions.Exception) {
	parsedOptions := options.ParseRepositoryOptions(opts...)
	now := time.Now()
	expiredAt := now.Add(-claimTimeout)
	tx := parsedOptions.DB.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, exceptions.New(
			"TransactionBeginFailed",
			"Outbox",
			"Claim",
			"Failed to begin the outbox claim transaction",
			http.StatusInternalServerError,
			true,
		).WithOrigin(tx.Error)
	}

	var events []schemas.OutboxEvent
	result := tx.
		Model(&schemas.OutboxEvent{}).
		Where("published_at IS NULL").
		Where("available_at <= ?", now).
		Where("claimed_at IS NULL OR claimed_at <= ?", expiredAt).
		Order("created_at ASC").
		Limit(batchSize).
		Clauses(clause.Locking{
			Strength: "UPDATE",
			Options:  "SKIP LOCKED",
		}).
		Find(&events)
	if result.Error != nil {
		tx.Rollback()
		return nil, exceptions.New(
			"FailedToGet",
			"Outbox",
			"Claim",
			"Failed to claim available outbox events",
			http.StatusInternalServerError,
			true,
		).WithOrigin(result.Error)
	}
	if len(events) == 0 {
		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			return nil, exceptions.New(
				"TransactionCommitFailed",
				"Outbox",
				"Claim",
				"Failed to commit the empty outbox claim transaction",
				http.StatusInternalServerError,
				true,
			).WithOrigin(err)
		}

		return events, nil
	}

	eventIds := make([]uuid.UUID, len(events))
	for index, event := range events {
		eventIds[index] = event.Id
	}
	result = tx.
		Model(&schemas.OutboxEvent{}).
		Where("id IN ?", eventIds).
		Updates(map[string]any{
			"claimed_by": workerId,
			"claimed_at": now,
		})
	if result.Error != nil {
		tx.Rollback()
		return nil, exceptions.New(
			"FailedToUpdate",
			"Outbox",
			"Claim",
			"Failed to claim available outbox events",
			http.StatusInternalServerError,
			true,
		).WithOrigin(result.Error)
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, exceptions.New(
			"TransactionCommitFailed",
			"Outbox",
			"Claim",
			"Failed to commit the outbox claim transaction",
			http.StatusInternalServerError,
			true,
		).WithOrigin(err)
	}

	for index := range events {
		events[index].ClaimedBy = &workerId
		events[index].ClaimedAt = &now
	}

	return events, nil
}

func (r *OutboxEventRepository) MarkPublishedMany(
	ctx context.Context,
	eventIds []uuid.UUID,
	workerId string,
	opts ...options.RepositoryOptions,
) *exceptions.Exception {
	if len(eventIds) == 0 {
		return nil
	}

	parsedOptions := options.ParseRepositoryOptions(opts...)
	now := time.Now()
	result := parsedOptions.DB.WithContext(ctx).
		Model(&schemas.OutboxEvent{}).
		Where("id IN ? AND claimed_by = ? AND published_at IS NULL", eventIds, workerId).
		Updates(map[string]any{
			"published_at": now,
			"last_error":   nil,
			"claimed_by":   nil,
			"claimed_at":   nil,
		})
	if result.Error != nil {
		return exceptions.New(
			"FailedToUpdate",
			"Outbox",
			"MarkPublished",
			"Failed to mark outbox events as published",
			http.StatusInternalServerError,
			true,
		).WithOrigin(result.Error)
	}

	return nil
}

func (r *OutboxEventRepository) MarkFailedMany(
	ctx context.Context,
	failureInputs []inputs.FailedOutboxEventInput,
	workerId string,
	opts ...options.RepositoryOptions,
) *exceptions.Exception {
	if len(failureInputs) == 0 {
		return nil
	}

	parsedOptions := options.ParseRepositoryOptions(opts...)
	valuePlaceholders := make([]string, 0, len(failureInputs))
	valueArguments := make([]any, 0, len(failureInputs)*3+1)
	for _, failureInput := range failureInputs {
		valuePlaceholders = append(valuePlaceholders, "(?::uuid, ?::text, ?::timestamptz)")
		valueArguments = append(
			valueArguments,
			failureInput.Id,
			failureInput.LastError,
			failureInput.AvailableAt,
		)
	}
	valueArguments = append(valueArguments, workerId)

	result := parsedOptions.DB.WithContext(ctx).Exec(
		fmt.Sprintf(`
			UPDATE "OutboxEventTable" AS outbox_event
			SET
				available_at = value.available_at,
				publish_count = outbox_event.publish_count + 1,
				last_error = value.last_error,
				claimed_by = NULL,
				claimed_at = NULL
			FROM (VALUES %s) AS value(id, last_error, available_at)
			WHERE outbox_event.id = value.id
				AND outbox_event.claimed_by = ?
				AND outbox_event.published_at IS NULL
		`, strings.Join(valuePlaceholders, ",")),
		valueArguments...,
	)
	if result.Error != nil {
		return exceptions.New(
			"FailedToUpdate",
			"Outbox",
			"MarkFailed",
			"Failed to schedule outbox event retries",
			http.StatusInternalServerError,
			true,
		).WithOrigin(result.Error)
	}

	return nil
}

func (r *OutboxEventRepository) DeletePublishedBefore(
	ctx context.Context,
	publishedBefore time.Time,
	opts ...options.RepositoryOptions,
) (int64, *exceptions.Exception) {
	parsedOptions := options.ParseRepositoryOptions(opts...)
	result := parsedOptions.DB.WithContext(ctx).
		Where("published_at IS NOT NULL AND published_at < ?", publishedBefore).
		Delete(&schemas.OutboxEvent{})
	if result.Error != nil {
		return 0, exceptions.New(
			"FailedToDelete",
			"Outbox",
			"Cleanup",
			"Failed to delete published outbox events",
			http.StatusInternalServerError,
			true,
		).WithOrigin(result.Error)
	}

	return result.RowsAffected, nil
}
