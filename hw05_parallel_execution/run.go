package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func Run(tasks []Task, goroutinesCount, maxErrors int) error {
	if maxErrors <= 0 {
		return ErrErrorsLimitExceeded
	}
	taskChan := make(chan Task, len(tasks))

	wg := &sync.WaitGroup{}
	wg.Add(goroutinesCount)

	for _, t := range tasks {
		taskChan <- t
	}
	close(taskChan)
	var errorsCount int32

	for i := 0; i < goroutinesCount; i++ {
		go worker(&errorsCount, int32(maxErrors), wg, taskChan)
	}
	wg.Wait()
	if errorsCount > int32(maxErrors) {
		return ErrErrorsLimitExceeded
	}
	return nil
}

func worker(errorsCount *int32, maxErrors int32, wg *sync.WaitGroup, taskChan chan Task) {
	defer wg.Done()
	for atomic.LoadInt32(errorsCount) < maxErrors {
		task, ok := <-taskChan
		if !ok {
			return
		}
		err := task()
		if err != nil {
			atomic.AddInt32(errorsCount, 1)
		}
	}
}
