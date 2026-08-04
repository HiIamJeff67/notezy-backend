package workers

import (
	"context"
	"hash/fnv"
	"os"
	"strings"
	"sync"

	realtimetypes "github.com/HiIamJeff67/notezy-backend/internal/realtimegateway/types"
	constants "github.com/HiIamJeff67/notezy-backend/shared/constants"
)

type WorkerManagerInterface interface {
	Attach(frame realtimetypes.InternalFrame) bool
	Detach(frame realtimetypes.InternalFrame)
	Forward(frame realtimetypes.InternalFrame) bool
	SetFrameHandler(handler func(realtimetypes.InternalFrame))
}

type WorkerManager struct {
	workers []*realtimeWorker
	cancel  context.CancelFunc
	waiter  sync.WaitGroup
}

func NewWorkerManager() *WorkerManager {
	ctx, cancel := context.WithCancel(context.Background())
	endpoints := strings.Split(os.Getenv("YJS_WORKER_URLS"), ",")
	workers := make([]*realtimeWorker, 0, len(endpoints))

	for _, endpoint := range endpoints {
		endpoint = strings.TrimSpace(endpoint)
		if endpoint == "" {
			continue
		}

		worker := &realtimeWorker{
			endpoint:       endpoint,
			activeChannels: make(map[string]realtimetypes.InternalFrame),
			outbound:       make(chan realtimetypes.InternalFrame, constants.RealtimeWorkerQueueSize),
		}
		workers = append(workers, worker)
	}

	manager := &WorkerManager{
		workers: workers,
		cancel:  cancel,
	}
	for _, worker := range manager.workers {
		manager.waiter.Add(1)
		go func() {
			defer manager.waiter.Done()
			worker.run(ctx)
		}()
	}

	return manager
}

func (m *WorkerManager) Shutdown() {
	m.cancel()
	for _, worker := range m.workers {
		worker.close()
	}
	m.waiter.Wait()
}

func (m *WorkerManager) Attach(frame realtimetypes.InternalFrame) bool {
	worker := m.getWorker(frame.ChannelId)
	if worker == nil || !worker.ready.Load() {
		return false
	}

	worker.activeChannelsMutex.Lock()
	worker.activeChannels[worker.channelKey(frame)] = frame
	worker.activeChannelsMutex.Unlock()

	return worker.enqueue(frame)
}

func (m *WorkerManager) Detach(frame realtimetypes.InternalFrame) {
	worker := m.getWorker(frame.ChannelId)
	if worker == nil {
		return
	}

	worker.activeChannelsMutex.Lock()
	delete(worker.activeChannels, worker.channelKey(frame))
	worker.activeChannelsMutex.Unlock()

	if worker.ready.Load() {
		worker.enqueue(frame)
	}
}

func (m *WorkerManager) Forward(frame realtimetypes.InternalFrame) bool {
	worker := m.getWorker(frame.ChannelId)
	if worker == nil || !worker.ready.Load() {
		return false
	}

	return worker.enqueue(frame)
}

func (m *WorkerManager) SetFrameHandler(handler func(realtimetypes.InternalFrame)) {
	for _, worker := range m.workers {
		worker.frameHandlerMutex.Lock()
		worker.frameHandler = handler
		worker.frameHandlerMutex.Unlock()
	}
}

func (m *WorkerManager) getWorker(channelId [16]byte) *realtimeWorker {
	if len(m.workers) == 0 {
		return nil
	}

	hasher := fnv.New32a()
	_, _ = hasher.Write(channelId[:])

	return m.workers[int(hasher.Sum32())%len(m.workers)]
}
