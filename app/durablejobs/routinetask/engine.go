package routinetask

import (
	"context"
	"sync/atomic"
	"time"

	"gorm.io/gorm"
)

type Engine struct {
	ticker         *time.Ticker
	isHealthy      int32
	claimer        Claimer
	handlerManager HandlerManager
}

func NewEngine(maxWorkers int, db *gorm.DB) Engine {
	return Engine{
		claimer:        NewClaimer(db),
		handlerManager: NewHandlerManager(maxWorkers, db),
		ticker:         time.NewTicker(time.Minute),
		isHealthy:      1,
	}
}

func (e *Engine) Start(ctx context.Context) {
	e.runOnce(ctx)

	for {
		select {
		case <-ctx.Done():
			e.Stop()
			return
		case <-e.ticker.C:
			e.runOnce(ctx)
		}
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

func (e *Engine) runOnce(ctx context.Context) {
	routineTasks, exception := e.claimer.Claim(ctx)
	if exception != nil {
		atomic.StoreInt32(&e.isHealthy, 0)
		return
	}
	if exception = e.handlerManager.Manage(ctx, routineTasks); exception != nil {
		atomic.StoreInt32(&e.isHealthy, 0)
		return
	}
	atomic.StoreInt32(&e.isHealthy, 1)
}
