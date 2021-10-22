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

func produce(done <-chan struct{}, tasks []Task, workersCount int) <-chan Task {
	queue := make(chan Task, workersCount)

	go func() {
		defer close(queue)

		for _, task := range tasks {
			select {
			case queue <- task:
			case <-done:
				return
			}
		}
	}()

	return queue
}

func consume(done chan<- struct{}, queue <-chan Task, wg *sync.WaitGroup, ec *ErrorsCounter, once *sync.Once) {
	defer wg.Done()

	for task := range queue {
		if ec.IsErrorLimitExceed() {
			once.Do(func() {
				close(done)
			})
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
	done := make(chan struct{})
	errorsCounter := ErrorsCounter{
		limit: m,
	}

	wg := sync.WaitGroup{}
	wg.Add(n)

	queue := produce(done, tasks, n)
	once := sync.Once{}

	for i := 0; i < n; i++ {
		go consume(done, queue, &wg, &errorsCounter, &once)
	}
	wg.Wait()

	if errorsCounter.IsErrorLimitExceed() {
		return ErrErrorsLimitExceeded
	}

	return nil
}
