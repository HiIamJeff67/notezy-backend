package concurrency

import "sync"

func Execute[T any, R any](
	inputs []T,
	workerCount int,
	processorFunc ProcessorFunc[T, R],
) []Result[R] {
	if len(inputs) == 0 {
		return []Result[R]{}
	}

	inputCount := len(inputs)
	workerCount = min(workerCount, inputCount, MaxNumOfWorker)
	jobs := make(chan Job[T], inputCount)
	results := make(chan Result[R], inputCount)
	var wg sync.WaitGroup

	worker := func() {
		defer wg.Done()
		for job := range jobs {
			resultData, err := processorFunc(job.Data)
			results <- Result[R]{
				Index: job.Index,
				Data:  resultData,
				Err:   err,
			}
		}
	}

	wg.Add(workerCount)
	for i := 0; i < workerCount; i++ {
		go worker()
	}

	for i, input := range inputs {
		jobs <- Job[T]{
			Index: i,
			Data:  input,
		}
	}
	close(jobs) // remember to close the jobs channel while it is not used anymore here

	wg.Wait()
	close(results)

	orderedResults := make([]Result[R], inputCount)
	for result := range results {
		orderedResults[result.Index] = result
	}

	return orderedResults
}

// TODO: execute in sequential, it can have multiple stages,
// and once a job is done in stage_i, it can be passed to the next stage of stage_i+1
// for each job, it can only be execute in stage_i if all the stage_j where 0 < j < i are already finished
func SequentialExecute() {}
