package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

type ErrorsCounter struct {
	limit int
	sync.Mutex
	count int
}

func (ec *ErrorsCounter) Increase() {
	ec.Lock()
	defer ec.Unlock()

	ec.count++
}

func (ec *ErrorsCounter) IsErrorLimitExceed() bool {
	ec.Lock()
	defer ec.Unlock()

	return ec.count >= ec.limit
}

func produce(tasks []Task, workersCount int, ec *ErrorsCounter) <-chan Task {
	queue := make(chan Task, workersCount)

	go func() {
		defer close(queue)

		for _, task := range tasks {
			if ec.IsErrorLimitExceed() {
				return
			}

			queue <- task
		}
	}()

	return queue
}

func consume(queue <-chan Task, wg *sync.WaitGroup, ec *ErrorsCounter) {
	defer wg.Done()

	for task := range queue {
		if ec.IsErrorLimitExceed() {
			return
		}

		err := task()
		if err != nil {
			ec.Increase()
		}
	}
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	errorsCounter := ErrorsCounter{
		limit: m,
	}

	wg := sync.WaitGroup{}
	wg.Add(n)

	queue := produce(tasks, n, &errorsCounter)

	for i := 0; i < n; i++ {
		go consume(queue, &wg, &errorsCounter)
	}
	wg.Wait()

	if errorsCounter.IsErrorLimitExceed() {
		return ErrErrorsLimitExceeded
	}

	return nil
}
