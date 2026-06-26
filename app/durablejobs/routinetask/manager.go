package routinetask

import (
	"sync"
	"sync/atomic"

	schemas "github.com/HiIamJeff67/notezy-backend/app/models/schemas"
)

type HandlerManager struct {
	maxWorkers    int
	activeWorkers atomic.Int32
	workerPool    sync.WaitGroup
}

func NewHandlerManager(maxWorkers int) HandlerManager {
	return HandlerManager{
		maxWorkers:    maxWorkers,
		activeWorkers: atomic.Int32{},
	}
}

func (hm *HandlerManager) Manage(claimedTasks []schemas.RoutineTask) {

}
