package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"gorm.io/gorm"

	inputs "github.com/HiIamJeff67/notezy-backend/app/models/inputs"
	repositories "github.com/HiIamJeff67/notezy-backend/app/models/repositories"
	logs "github.com/HiIamJeff67/notezy-backend/app/monitor/logs"
	metrics "github.com/HiIamJeff67/notezy-backend/app/monitor/metrics"
	traces "github.com/HiIamJeff67/notezy-backend/app/monitor/traces"
	options "github.com/HiIamJeff67/notezy-backend/app/options"
	realtimetypes "github.com/HiIamJeff67/notezy-backend/app/realtime/types"
)

type YjsPersistenceServiceInterface interface {
	LoadDocument(ctx context.Context, blockPackId uuid.UUID) (*realtimetypes.YjsDocumentState, error)
	AppendUpdate(ctx context.Context, blockPackId uuid.UUID, persistenceBatchId uuid.UUID, originConnectionId *uuid.UUID, payload []byte) (int64, error)
	GetCompactableYjsDocumentWithUpdates(ctx context.Context, blockPackId uuid.UUID) (*realtimetypes.YjsCompactionInput, error)
	ApplyCompactedYjsDocument(ctx context.Context, blockPackId uuid.UUID, result realtimetypes.YjsCompactionResult) (bool, error)
}

type YjsPersistenceService struct {
	db                     *gorm.DB
	blockPackYjsRepository repositories.BlockPackYjsRepositoryInterface
}

func NewYjsPersistenceService(db *gorm.DB) YjsPersistenceServiceInterface {
	return &YjsPersistenceService{
		db:                     db,
		blockPackYjsRepository: repositories.NewBlockPackYjsRepository(),
	}
}

func (s *YjsPersistenceService) LoadDocument(
	ctx context.Context, blockPackId uuid.UUID,
) (state *realtimetypes.YjsDocumentState, err error) {
	start := time.Now()
	ctx, span := traces.NotezyTracer.Start(ctx, "yjs.document.load")
	defer func() { traces.NotezyTracer.End(span, err) }()

	db := s.db.WithContext(ctx)

	document, updates, err := s.blockPackYjsRepository.LoadDocumentAndUpdates(
		blockPackId,
		options.WithDB(db),
	)
	if err != nil {
		metrics.NotezyMeter.Count(ctx, "yjs.operation.count", 1,
			attribute.String("operation", "document.load"),
			attribute.String("outcome", "error"),
		)
		metrics.NotezyMeter.Duration(ctx, "yjs.operation.duration", time.Since(start),
			attribute.String("operation", "document.load"),
			attribute.String("outcome", "error"),
		)
		logs.NotezyLogger.Error(ctx, err, "failed to load Yjs document", attribute.String("operation", "document.load"))

		return nil, err
	}

	state = &realtimetypes.YjsDocumentState{
		Snapshot:               document.Snapshot,
		StateVector:            document.StateVector,
		LastUpdateSequence:     document.LastUpdateSequence,
		CompactedUntilSequence: document.CompactedUntilSequence,
		ProjectedUntilSequence: document.ProjectedUntilSequence,
		Updates:                make([]realtimetypes.YjsDocumentUpdate, len(updates)),
	}
	for index, update := range updates {
		state.Updates[index] = realtimetypes.YjsDocumentUpdate{
			UpdateSequence: update.UpdateSequence,
			Payload:        update.Payload,
		}
	}

	span.SetAttributes(attribute.Int("yjs.update_count", len(updates)))
	metrics.NotezyMeter.Count(ctx, "yjs.operation.count", 1,
		attribute.String("operation", "document.load"),
		attribute.String("outcome", "success"),
	)
	metrics.NotezyMeter.Duration(ctx, "yjs.operation.duration", time.Since(start),
		attribute.String("operation", "document.load"),
		attribute.String("outcome", "success"),
	)

	return state, nil
}

func (s *YjsPersistenceService) AppendUpdate(
	ctx context.Context,
	blockPackId uuid.UUID,
	persistenceBatchId uuid.UUID,
	originConnectionId *uuid.UUID,
	payload []byte,
) (updateSequence int64, err error) {
	start := time.Now()
	ctx, span := traces.NotezyTracer.Start(ctx, "yjs.update.append")
	defer func() { traces.NotezyTracer.End(span, err) }()

	db := s.db.WithContext(ctx)

	updateSequence, err = s.blockPackYjsRepository.AppendUpdate(
		blockPackId,
		inputs.AppendBlockPackYjsUpdateInput{
			PersistenceBatchId: persistenceBatchId,
			OriginConnectionId: originConnectionId,
			Payload:            payload,
		},
		options.WithDB(db),
	)
	if err != nil {
		metrics.NotezyMeter.Count(ctx, "yjs.operation.count", 1,
			attribute.String("operation", "update.append"),
			attribute.String("outcome", "error"),
		)
		metrics.NotezyMeter.Duration(ctx, "yjs.operation.duration", time.Since(start),
			attribute.String("operation", "update.append"),
			attribute.String("outcome", "error"),
		)
		metrics.NotezyMeter.Bytes(ctx, "yjs.payload.bytes", int64(len(payload)), attribute.String("operation", "update.append"))
		logs.NotezyLogger.Error(ctx, err, "failed to append Yjs update", attribute.String("operation", "update.append"))

		return 0, err
	}

	span.SetAttributes(attribute.Int("yjs.payload_bytes", len(payload)))
	metrics.NotezyMeter.Count(ctx, "yjs.operation.count", 1,
		attribute.String("operation", "update.append"),
		attribute.String("outcome", "success"),
	)
	metrics.NotezyMeter.Duration(ctx, "yjs.operation.duration", time.Since(start),
		attribute.String("operation", "update.append"),
		attribute.String("outcome", "success"),
	)
	metrics.NotezyMeter.Bytes(ctx, "yjs.payload.bytes", int64(len(payload)), attribute.String("operation", "update.append"))

	return updateSequence, nil
}

func (s *YjsPersistenceService) GetCompactableYjsDocumentWithUpdates(
	ctx context.Context, blockPackId uuid.UUID,
) (input *realtimetypes.YjsCompactionInput, err error) {
	start := time.Now()
	ctx, span := traces.NotezyTracer.Start(ctx, "yjs.compaction.load")
	defer func() { traces.NotezyTracer.End(span, err) }()

	db := s.db.WithContext(ctx)

	document, updates, err := s.blockPackYjsRepository.GetCompactableYjsDocumentWithUpdates(
		blockPackId,
		options.WithDB(db),
	)
	if err != nil || document == nil {
		if err != nil {
			logs.NotezyLogger.Error(ctx, err, "failed to load compactable Yjs document", attribute.String("operation", "compaction.load"))
		}
		metrics.NotezyMeter.Count(ctx, "yjs.operation.count", 1,
			attribute.String("operation", "compaction.load"),
			attribute.String("outcome", "error"),
		)
		metrics.NotezyMeter.Duration(ctx, "yjs.operation.duration", time.Since(start),
			attribute.String("operation", "compaction.load"),
			attribute.String("outcome", "error"),
		)

		return nil, err
	}

	input = &realtimetypes.YjsCompactionInput{
		Snapshot:                   document.Snapshot,
		StateVector:                document.StateVector,
		BaseCompactedUntilSequence: document.CompactedUntilSequence,
		CutoffSequence:             document.LastUpdateSequence,
		Updates:                    make([]realtimetypes.YjsDocumentUpdate, len(updates)),
	}
	for index, update := range updates {
		input.Updates[index] = realtimetypes.YjsDocumentUpdate{
			UpdateSequence: update.UpdateSequence,
			Payload:        update.Payload,
		}
	}

	span.SetAttributes(attribute.Int("yjs.update_count", len(updates)))
	metrics.NotezyMeter.Count(ctx, "yjs.operation.count", 1,
		attribute.String("operation", "compaction.load"),
		attribute.String("outcome", "success"),
	)
	metrics.NotezyMeter.Duration(ctx, "yjs.operation.duration", time.Since(start),
		attribute.String("operation", "compaction.load"),
		attribute.String("outcome", "success"),
	)

	return input, nil
}

func (s *YjsPersistenceService) ApplyCompactedYjsDocument(
	ctx context.Context, blockPackId uuid.UUID, result realtimetypes.YjsCompactionResult,
) (applied bool, err error) {
	start := time.Now()
	ctx, span := traces.NotezyTracer.Start(ctx, "yjs.compaction.apply")
	defer func() { traces.NotezyTracer.End(span, err) }()

	db := s.db.WithContext(ctx)

	applied, err = s.blockPackYjsRepository.ApplyCompactedYjsDocument(
		blockPackId,
		inputs.ApplyCompactedBlockPackYjsDocumentInput{
			BaseCompactedUntilSequence: result.BaseCompactedUntilSequence,
			CutoffSequence:             result.CutoffSequence,
			Snapshot:                   result.Snapshot,
			StateVector:                result.StateVector,
		},
		options.WithDB(db),
	)
	if err != nil {
		metrics.NotezyMeter.Count(ctx, "yjs.operation.count", 1,
			attribute.String("operation", "compaction.apply"),
			attribute.String("outcome", "error"),
		)
		metrics.NotezyMeter.Duration(ctx, "yjs.operation.duration", time.Since(start),
			attribute.String("operation", "compaction.apply"),
			attribute.String("outcome", "error"),
		)
		metrics.NotezyMeter.Bytes(ctx, "yjs.payload.bytes", int64(len(result.Snapshot)), attribute.String("operation", "compaction.apply"))
		logs.NotezyLogger.Error(ctx, err, "failed to apply compacted Yjs document", attribute.String("operation", "compaction.apply"))

		return false, err
	}

	span.SetAttributes(attribute.Bool("yjs.applied", applied))
	metrics.NotezyMeter.Count(ctx, "yjs.operation.count", 1,
		attribute.String("operation", "compaction.apply"),
		attribute.String("outcome", "success"),
	)
	metrics.NotezyMeter.Duration(ctx, "yjs.operation.duration", time.Since(start),
		attribute.String("operation", "compaction.apply"),
		attribute.String("outcome", "success"),
	)
	metrics.NotezyMeter.Bytes(ctx, "yjs.payload.bytes", int64(len(result.Snapshot)), attribute.String("operation", "compaction.apply"))

	return applied, nil
}
