package workers

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	logs "github.com/HiIamJeff67/notezy-backend/shared/platform/observability/logs"
	metrics "github.com/HiIamJeff67/notezy-backend/shared/platform/observability/metrics"

	repositories "github.com/HiIamJeff67/notezy-backend/internal/core/data/database/repositories"
	schemas "github.com/HiIamJeff67/notezy-backend/internal/core/data/database/schemas"
)

type YjsMaintenanceReconciliationWorkerInterface interface {
	Start(ctx context.Context) func()
	Reconcile(ctx context.Context) error
}

type YjsMaintenanceReconciliationWorker struct {
	db                    *gorm.DB
	outboxEventRepository repositories.OutboxEventRepositoryInterface
}

func NewYjsMaintenanceReconciliationWorker(
	db *gorm.DB,
	outboxEventRepository repositories.OutboxEventRepositoryInterface,
) YjsMaintenanceReconciliationWorkerInterface {
	return &YjsMaintenanceReconciliationWorker{
		db:                    db,
		outboxEventRepository: outboxEventRepository,
	}
}

/* ============================== Constants ============================== */

const (
	yjsMaintenanceReconciliationBatchSize = 256
	yjsMaintenanceReconciliationInterval  = time.Hour
)

/* ============================== Auxiliary Functions ============================== */

func (w *YjsMaintenanceReconciliationWorker) reconcile(ctx context.Context) {
	if err := w.Reconcile(ctx); err != nil && ctx.Err() == nil {
		// Reconciliation is a safety net. The normal outbox path remains available
		// when a scan fails, so the next interval can retry without stopping Core.
		if logs.NotezyLogger != nil {
			logs.NotezyLogger.Error(ctx, err, "Yjs maintenance reconciliation failed")
		}
		return
	}
}

/* ============================== Worker Methods ============================== */

func (w *YjsMaintenanceReconciliationWorker) Start(ctx context.Context) func() {
	workerCtx, cancel := context.WithCancel(ctx)
	done := make(chan struct{})

	go func() {
		defer close(done)
		w.reconcile(workerCtx)

		ticker := time.NewTicker(yjsMaintenanceReconciliationInterval)
		defer ticker.Stop()
		for {
			select {
			case <-workerCtx.Done():
				return
			case <-ticker.C:
				w.reconcile(workerCtx)
			}
		}
	}()

	return func() {
		cancel()
		<-done
	}
}

func (w *YjsMaintenanceReconciliationWorker) Reconcile(ctx context.Context) error {
	if w == nil || w.db == nil || w.outboxEventRepository == nil {
		return errors.New("Yjs maintenance reconciliation dependencies are required")
	}

	var documents []schemas.BlockPackYjsDocument
	result := w.db.WithContext(ctx).
		Select("id, block_pack_id, last_update_sequence, compacted_until_sequence, projected_until_sequence, last_compacted_at, snapshot, state_vector").
		Where("deleted_at IS NULL").
		Where("last_update_sequence > 0").
		Where("last_update_sequence > compacted_until_sequence OR last_update_sequence > projected_until_sequence").
		Order("updated_at ASC").
		Limit(yjsMaintenanceReconciliationBatchSize).
		Find(&documents)
	if result.Error != nil {
		return fmt.Errorf("load stale Yjs documents: %w", result.Error)
	}
	if len(documents) == 0 {
		if metrics.NotezyMeter != nil {
			metrics.NotezyMeter.Value(ctx, "yjs.maintenance.reconciliation.stale_documents", 0)
		}
		return nil
	}
	if metrics.NotezyMeter != nil {
		metrics.NotezyMeter.Value(ctx, "yjs.maintenance.reconciliation.stale_documents", int64(len(documents)))
	}

	tx := w.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return fmt.Errorf("begin Yjs maintenance reconciliation transaction: %w", tx.Error)
	}
	for _, document := range documents {
		if err := w.outboxEventRepository.EnqueueYjsMaintenanceHint(
			tx,
			uuid.NewString(),
			document.BlockPackId,
			"reconciliation",
		); err != nil {
			tx.Rollback()
			return fmt.Errorf("enqueue Yjs maintenance hint for %s: %w", document.BlockPackId, err)
		}
	}
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("commit Yjs maintenance reconciliation transaction: %w", err)
	}
	if metrics.NotezyMeter != nil {
		metrics.NotezyMeter.Count(ctx, "yjs.maintenance.reconciliation.hints_enqueued", int64(len(documents)))
	}

	return nil
}
