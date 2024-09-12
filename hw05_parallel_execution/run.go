package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	tasksChan := make(chan Task, n)
	errorsChan := make(chan error, len(tasks))
	stopChan := make(chan struct{}) // flag that all tasks are finished or m errors occurred

	countErrors := 0
	wg := sync.WaitGroup{}

	// starting n workers
	for i := 1; i <= n; i++ {
		wg.Add(1)
		go func() { // each function call is one worker
			defer wg.Done()
			for task := range tasksChan {
				select { // если уберу этот селект, то код успевает делать n+m+2 тасков пока не закроется tasksChan
				case <-stopChan: // before executing task, check if we have to stop
					return
				default:
					err := task()
					if err != nil {
						errorsChan <- err
					}
				}
			}
		}()
	}

	go func() {
		defer close(tasksChan)
		for _, task := range tasks {
			select {
			case <-stopChan: // stop sending tasks to channel if we have m errors
				return
			case tasksChan <- task:
			}
		}
	}()

	go func() {
		defer close(stopChan)
		for range errorsChan {
			countErrors++
			if m > 0 && countErrors >= m {
				return
			}
		}
	}()

	wg.Wait()
	close(errorsChan) // close only when all workers are done with their work

	if m > 0 && countErrors >= m {
		return ErrErrorsLimitExceeded
	}

	return nil
}
