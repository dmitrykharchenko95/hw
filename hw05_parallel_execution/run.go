package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	tasksCh := make(chan struct{}, n)
	maxErrCh := make(chan struct{}, 1)

	var wg sync.WaitGroup
	var mu sync.Mutex

	if m == 0 {
		return ErrErrorsLimitExceeded
	}

	for _, task := range tasks {
		wg.Add(1)
		tasksCh <- struct{}{}

		select {
		case <-maxErrCh:
			wg.Done()
			wg.Wait()
			return ErrErrorsLimitExceeded
		default:
		}

		go func(task Task, m *int, tasksCh chan struct{}) {
			defer wg.Done()
			err := task()
			if err != nil {
				mu.Lock()
				if *m == 0 {
					maxErrCh <- struct{}{}
				}
				*m--
				mu.Unlock()
			}
			<-tasksCh
		}(task, &m, tasksCh)
	}

	wg.Wait()
	return nil
}
