package emails

import (
	"container/heap"
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	exceptions "notezy-backend/app/exceptions"
	logs "notezy-backend/app/logs"
	constants "notezy-backend/shared/constants"
)

/* ============================== Initialization & Instance ============================== */

type EmailWorkerManager struct {
	maxWorkers    int
	activeWorkers int32
	workerPool    sync.WaitGroup
	ctx           context.Context
	cancel        context.CancelFunc
	emailSender   EmailSender

	buffer      *EmailBuffer
	bufferMutex sync.RWMutex

	monitorTicker *time.Ticker
	isMonitoring  int32
}

func NewEmailWorkerManager(maxWorkers int, sender EmailSender) *EmailWorkerManager {
	ctx, cancel := context.WithCancel(context.Background())
	return &EmailWorkerManager{
		maxWorkers:   maxWorkers,
		ctx:          ctx,
		cancel:       cancel,
		emailSender:  sender,
		buffer:       NewEmailBuffer(),
		isMonitoring: 0,
	}
}

var (
	NotezyEmailWorkerManager = NewEmailWorkerManager(16, *NotezyEmailSender)
)

/* ============================== Aulixary Functions ============================== */

func (ewm *EmailWorkerManager) generateTaskId() string {
	return fmt.Sprintf("email_task_%d", time.Now().UnixNano())
}

/* ============================== Private Methods ============================== */

func (ewm *EmailWorkerManager) processTask(task *EmailTask, workerId int) {
	exception := ewm.emailSender.SyncSend(task.Object.To, task.Object.Subject, task.Object.Body, task.Object.EmailContentType)
	if exception != nil {
		exceptions.Email.FailedToSendEmailByWorkers(workerId, task.Retries+1, task.MaxRetries).Log()

		task.Retries++
		if task.Retries < task.MaxRetries {
			task.Priority = max(0, task.Priority-1)

			go func() {
				retryDelay := time.Duration(task.Retries) * time.Second * 30
				time.Sleep(retryDelay)
				ewm.enqueueTask(task)
			}()
		} else {
			logs.FInfo("Worker %d gave up on task %s after %d of retries", workerId, task.ID, task.MaxRetries)
		}
	} else {
		logs.FInfo("Worker %d successfully sent email task Id is %s", workerId, task.ID)
	}
}

func (ewm *EmailWorkerManager) createWorker(task *EmailTask) {
	current := atomic.AddInt32(&ewm.activeWorkers, 1)
	workerId := int(current)

	ewm.workerPool.Add(1)
	go func() {
		defer func() {
			atomic.AddInt32(&ewm.activeWorkers, -1)
			ewm.workerPool.Done()
			logs.FInfo("Worker %d completed and stopped", workerId)
		}()

		logs.FInfo("Worker %d started for task: %s", workerId, task.ID)
		ewm.processTask(task, workerId)
	}()
}

func (ewm *EmailWorkerManager) dispatchTasks() {
	ewm.bufferMutex.Lock()
	defer ewm.bufferMutex.Unlock()

	activeWorkers := ewm.GetActiveWorkerCount()
	numOfTasks := ewm.buffer.Len()
	if numOfTasks == 0 {
		return
	}

	workersNeeded := min(numOfTasks, ewm.maxWorkers-activeWorkers)
	// the worker manager will process the job once the number workers is enough
	// else it will wait, and the tasks are still storing in the buffer
	for i := 0; i < workersNeeded; i++ {
		if ewm.buffer.Len() == 0 {
			break
		}

		task := heap.Pop(ewm.buffer).(*EmailTask)
		ewm.createWorker(task)
	}
}

// this function is just for monitoring the constaints,
// so that the entire producer & consumer model of handling the email will work
// it won't logging enoggh informations for actual monitoring
// it's only used to maintain the necessary conditions between workers and tasks(as the producer & consumer model)
func (ewm *EmailWorkerManager) tryStartMonitoring() {
	if !atomic.CompareAndSwapInt32(&ewm.isMonitoring, 0, 1) {
		// if the current ewm.isMonitoring == 0 (compare "addr" and "old"(0 in this case))
		// then ewm.isMonitoring = "new"(1 in this case) and return true
		// else return false

		// the task of the current go routine is unnecessary to process anymore,
		// since there's some other (1 go routine) which has ensured the email worker manager is working
		// on dispatching the email worker to send the emails including the enqueued emails of current go routine
		return
	}

	ewm.monitorTicker = time.NewTicker(constants.EmailWorkerManagerTickerDuration)

	go func() {
		defer func() {
			atomic.StoreInt32(&ewm.isMonitoring, 0)
			logs.FInfo("Stopped queue monitoring")
		}()

		logs.FInfo("Started queue monitoring")

		for {
			select {
			case <-ewm.ctx.Done():
				return
			case <-ewm.monitorTicker.C:
				ewm.dispatchTasks() // call the method to dispatch the task

				ewm.bufferMutex.RLock()
				if ewm.buffer.IsEmpty() && ewm.GetActiveWorkerCount() == 0 {
					ewm.monitorTicker.Stop()
					ewm.bufferMutex.RUnlock()
					return
				}
				ewm.bufferMutex.RUnlock()
			}
		}
	}()
}

func (ewm *EmailWorkerManager) enqueueTask(task *EmailTask) error {
	ewm.bufferMutex.Lock()
	ewm.buffer.EnqueueTask(task)
	bufferSize := ewm.buffer.Len()
	ewm.bufferMutex.Unlock()

	logs.FInfo("Enqueued email task: ID=%s, Type=%s, Priority=%d, Queue size: %d",
		task.ID, task.Type, task.Priority, bufferSize)

	ewm.tryStartMonitoring()

	return nil
}

/* ============================== Public Methods ============================== */

func (ewm *EmailWorkerManager) GetActiveWorkerCount() int {
	return int(atomic.LoadInt32(&ewm.activeWorkers))
}

func (ewm *EmailWorkerManager) Shutdown() {
	logs.FInfo("Shutting down email worker manager...")

	ewm.cancel()
	if ewm.monitorTicker != nil {
		ewm.monitorTicker.Stop()
	}
	ewm.workerPool.Wait()

	logs.FInfo("Email worker manager stopped")
}

func (ewm *EmailWorkerManager) GetStats() map[string]interface{} {
	ewm.bufferMutex.RLock()
	bufferSize := ewm.buffer.Len()
	ewm.bufferMutex.RUnlock()

	return map[string]interface{}{
		"bufferSize":    bufferSize,
		"activeWorkers": ewm.GetActiveWorkerCount(),
		"maxWorkers":    ewm.maxWorkers,
		"isMonitoring":  atomic.LoadInt32(&ewm.isMonitoring) == 1,
	}
}

func (ewm *EmailWorkerManager) Enqueue(
	emailObject EmailObject,
	emailTaskType EmailTaskType,
	maxRetries int,
	priority int,
) *exceptions.Exception {
	task := &EmailTask{
		ID:         ewm.generateTaskId(),
		Type:       emailTaskType,
		Object:     emailObject,
		CreatedAt:  time.Now(),
		Retries:    0,
		MaxRetries: maxRetries,
		Priority:   priority,
	}
	err := ewm.enqueueTask(task)
	if err != nil {
		return exceptions.Email.FailedToEnqueueTaskToEmailWorkerManager().WithError(err)
	}

	return nil
}
