package yjsmaintenance

import (
	"context"
	"sync"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"gorm.io/gorm"

	logs "github.com/HiIamJeff67/notezy-backend/app/monitor/logs"
	metrics "github.com/HiIamJeff67/notezy-backend/app/monitor/metrics"
	traces "github.com/HiIamJeff67/notezy-backend/app/monitor/traces"
	services "github.com/HiIamJeff67/notezy-backend/app/services"
	constants "github.com/HiIamJeff67/notezy-backend/shared/constants"
)

type Engine struct {
	ticker       *time.Ticker
	claimer      Claimer
	handler      Handler
	workerClient WorkerClient
}

func NewEngine(db *gorm.DB) *Engine {
	blockProjectionService := services.NewBlockProjectionService(db)

	return &Engine{
		ticker:       time.NewTicker(constants.YjsMaintenanceScanInterval),
		claimer:      NewClaimer(db),
		handler:      NewHandler(db, blockProjectionService),
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
	start := time.Now()
	ctx, span := traces.NotezyTracer.Start(ctx, "yjs.maintenance.run")
	var finalErr error
	defer func() { traces.NotezyTracer.End(span, finalErr) }()

	outcome := "success"
	inputs, err := e.claimer.ClaimCompactions(ctx)
	if err != nil {
		outcome = "error"
		finalErr = err
		logs.NotezyLogger.Error(ctx, err, "failed to claim Yjs maintenance inputs", attribute.String("operation", "maintenance.claim"))
	} else if len(inputs) > 0 {
		results, err := e.workerClient.Compact(ctx, inputs)
		if err != nil {
			outcome = "error"
			finalErr = err
			logs.NotezyLogger.Error(ctx, err, "failed to compact Yjs maintenance inputs", attribute.String("operation", "maintenance.compact"))
		} else if _, err := e.handler.HandleCompactions(ctx, inputs, results); err != nil {
			outcome = "error"
			finalErr = err
			logs.NotezyLogger.Error(ctx, err, "failed to persist compacted Yjs documents", attribute.String("operation", "maintenance.apply"))
		}
	}

	if err := e.handler.Cleanup(ctx); err != nil {
		outcome = "error"
		finalErr = err
		logs.NotezyLogger.Error(ctx, err, "failed to clean compacted Yjs updates", attribute.String("operation", "maintenance.cleanup"))
	}

	projectionInputs, err := e.claimer.ClaimProjections(ctx)
	if err != nil {
		outcome = "error"
		finalErr = err
		logs.NotezyLogger.Error(ctx, err, "failed to claim Yjs projection inputs", attribute.String("operation", "maintenance.projection.claim"))
	} else if len(projectionInputs) > 0 {
		projectionResults, err := e.workerClient.Project(ctx, projectionInputs)
		if err != nil {
			outcome = "error"
			finalErr = err
			logs.NotezyLogger.Error(ctx, err, "failed to project Yjs documents", attribute.String("operation", "maintenance.projection.project"))
		} else if _, err := e.handler.HandleProjections(ctx, projectionInputs, projectionResults); err != nil {
			outcome = "error"
			finalErr = err
			logs.NotezyLogger.Error(ctx, err, "failed to persist projected blocks", attribute.String("operation", "maintenance.projection.apply"))
		}
	}

	metrics.NotezyMeter.Count(ctx, "yjs.operation.count", 1,
		attribute.String("operation", "maintenance.run"),
		attribute.String("outcome", outcome),
	)
	metrics.NotezyMeter.Duration(ctx, "yjs.operation.duration", time.Since(start),
		attribute.String("operation", "maintenance.run"),
		attribute.String("outcome", outcome),
	)
}
