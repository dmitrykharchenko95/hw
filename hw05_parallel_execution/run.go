package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	var max int
	tasksCh := make(chan struct{}, n)

	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, task := range tasks {
		tasksCh <- struct{}{}
		wg.Add(1)

		mu.Lock()
		max = m
		mu.Unlock()

		if max == 0 {
			wg.Done()
			wg.Wait()
			return ErrErrorsLimitExceeded
		}
		go func(task Task, m *int, tasksCh chan struct{}) {
			defer wg.Done()
			err := task()
			<-tasksCh

			if err != nil {
				mu.Lock()
				*m--
				mu.Unlock()
			}
		}(task, &m, tasksCh)
	}

	wg.Wait()
	return nil
}
