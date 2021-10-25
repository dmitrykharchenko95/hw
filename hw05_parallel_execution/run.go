package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	M := int32(m)
	var err error

	maxErrCh := make(chan struct{})
	tasksCh := make(chan Task)

	var wg sync.WaitGroup

	if m == 0 {
		return ErrErrorsLimitExceeded
	}

	go func(err *error) {
		defer close(tasksCh)
		for _, task := range tasks {
			select {
			case <-maxErrCh:
				*err = ErrErrorsLimitExceeded
				return
			default:
				select {
				case <-maxErrCh:
					*err = ErrErrorsLimitExceeded
					return
				case tasksCh <- task:
				}
			}
		}
	}(&err)

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(M *int32) {
			defer wg.Done()
			for task := range tasksCh {
				if task() != nil {
					if atomic.CompareAndSwapInt32(M, 0, 0) {
						close(maxErrCh)
					}
					atomic.AddInt32(M, -1)
				}
			}
		}(&M)
	}

	wg.Wait()
	return err
}
