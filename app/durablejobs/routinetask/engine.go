package routinetask

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"gorm.io/gorm"

	constants "github.com/HiIamJeff67/notezy-backend/shared/constants"
)

type Engine struct {
	ticker         *time.Ticker
	isHealthy      int32
	claimer        Claimer
	handlerManager HandlerManager
}

func NewEngine(db *gorm.DB, maxWorkers ...int) Engine {
	initialMaxWorkers := constants.RoutineTaskEngineMaxWorkers
	if len(maxWorkers) > 0 {
		initialMaxWorkers = min(initialMaxWorkers, maxWorkers[0])
	}
	return Engine{
		claimer:        NewClaimer(db),
		handlerManager: NewHandlerManager(initialMaxWorkers, db),
		ticker:         time.NewTicker(time.Minute),
		isHealthy:      1,
	}
}

func (e *Engine) runOnce(ctx context.Context) {
	routineTasks, taskIdToOwnerId, exception := e.claimer.Claim(ctx)
	if exception != nil {
		atomic.StoreInt32(&e.isHealthy, 0)
		return
	}
	if exception = e.handlerManager.Manage(ctx, routineTasks, taskIdToOwnerId); exception != nil {
		atomic.StoreInt32(&e.isHealthy, 0)
		return
	}
	atomic.StoreInt32(&e.isHealthy, 1)
}

func (e *Engine) Start(ctx context.Context) func() {
	ctx, cancel := context.WithCancel(ctx)
	done := make(chan struct{})
	var shutdownOnce sync.Once

	go func() {
		// note that defer is added using stack (LIFO data structure)
		// hence `e.Stop()` will be executed before the `close(done)` below
		defer close(done) // executed last
		defer e.Stop()    // executed first

		e.runOnce(ctx) // run once right after started
		for {
			select {
			case <-ctx.Done():
				return
			case <-e.ticker.C:
				e.runOnce(ctx)
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
	atomic.StoreInt32(&e.isHealthy, 0)
}

func (e *Engine) IsHealthy() bool {
	return atomic.LoadInt32(&e.isHealthy) == 1
}
