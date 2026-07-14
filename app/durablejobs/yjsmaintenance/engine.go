package yjsmaintenance

import (
	"context"
	"sync"
	"time"

	"gorm.io/gorm"

	constants "github.com/HiIamJeff67/notezy-backend/shared/constants"
)

type Engine struct {
	ticker       *time.Ticker
	claimer      Claimer
	handler      Handler
	workerClient WorkerClient
}

func NewEngine(db *gorm.DB) *Engine {
	return &Engine{
		ticker:       time.NewTicker(constants.YjsMaintenanceScanInterval),
		claimer:      NewClaimer(db),
		handler:      NewHandler(db),
		workerClient: NewWorkerClient(),
	}
}

func (e *Engine) Start(ctx context.Context) func() {
	ctx, cancel := context.WithCancel(ctx)
	done := make(chan struct{})
	var shutdownOnce sync.Once

	go func() {
		defer close(done)
		defer e.ticker.Stop()

		e.runOnce(ctx)
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

func (e *Engine) runOnce(ctx context.Context) {
	inputs, err := e.claimer.Claim(ctx)
	if err == nil && len(inputs) > 0 {
		results, err := e.workerClient.Compact(ctx, inputs)
		if err == nil {
			_, _ = e.handler.Handle(ctx, inputs, results)
		}
	}

	_ = e.handler.Cleanup(ctx)
}
