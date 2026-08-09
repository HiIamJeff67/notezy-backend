package workers

import (
	"context"
	"sync"
	"time"

	logs "github.com/HiIamJeff67/notezy-backend/shared/platform/observability/logs"

	services "github.com/HiIamJeff67/notezy-backend/internal/notification/services"
)

type CleanupWorker struct {
	service   *services.NotificationService
	interval  time.Duration
	retention time.Duration
}

func NewCleanupWorker(
	service *services.NotificationService,
	interval time.Duration,
	retention time.Duration,
) *CleanupWorker {
	return &CleanupWorker{
		service:   service,
		interval:  interval,
		retention: retention,
	}
}

func (w *CleanupWorker) Start(ctx context.Context) func() {
	workerCtx, cancel := context.WithCancel(ctx)
	var waitGroup sync.WaitGroup
	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		ticker := time.NewTicker(w.interval)
		defer ticker.Stop()
		for {
			if _, err := w.service.HardDelete(workerCtx, time.Now().UTC(), w.retention); err != nil && logs.NotezyLogger != nil {
				logs.NotezyLogger.Error(workerCtx, err, "Failed to clean Notification records")
			}
			select {
			case <-workerCtx.Done():
				return
			case <-ticker.C:
			}
		}
	}()

	return func() {
		cancel()
		waitGroup.Wait()
	}
}
