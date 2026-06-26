package routinetask

import (
	"context"
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
		handlerManager: NewHandlerManager(maxWorkers),
		ticker:         time.NewTicker(time.Minute),
		isHealthy:      1,
	}
}

func (e *Engine) Start(ctx context.Context) {
	// routineTasks, exception := e.claimer.Claim(ctx)
	// if exception != nil {

	// }
}
