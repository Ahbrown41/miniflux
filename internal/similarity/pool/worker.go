package pool

import (
	"sync"
)

// Task defines the task to be processed by the thread pool.
type Task[T any, R any] struct {
	ID       int
	Data     T
	Function func(T) (R, error)
}

// Result defines the structure for the result returned by the thread pool.
type Result[R any] struct {
	ID    int
	Value R
}

// Worker represents a single worker in the thread pool.
type Worker[T any, R any] struct {
	ID         int
	TaskQueue  chan Task[T, R]
	ResultChan chan Result[R]
	ErrChan    chan error
	WorkerPool chan chan Task[T, R]
	QuitChan   chan bool
}

// ThreadPool represents the thread pool.
type ThreadPool[T any, R any] struct {
	TaskQueue  chan Task[T, R]
	ResultChan chan Result[R]
	ErrChan    chan error
	WorkerPool chan chan Task[T, R]
	Workers    []Worker[T, R]
	QuitChan   chan bool
	wg         sync.WaitGroup
}

// NewWorker creates a new worker.
func NewWorker[T any, R any](id int, workerPool chan chan Task[T, R], resultChan chan Result[R], errChan chan error) Worker[T, R] {
	return Worker[T, R]{
		ID:         id,
		TaskQueue:  make(chan Task[T, R]),
		ResultChan: resultChan,
		ErrChan:    errChan,
		WorkerPool: workerPool,
		QuitChan:   make(chan bool),
	}
}

// Start starts the worker.
func (w Worker[T, R]) Start(wg *sync.WaitGroup) {
	go func() {
		defer wg.Done()
		for {
			// Add this worker's task queue to the worker pool.
			w.WorkerPool <- w.TaskQueue

			select {
			case task := <-w.TaskQueue:
				result, err := task.Function(task.Data)
				if err != nil {
					w.ErrChan <- err
				} else {
					w.ResultChan <- Result[R]{ID: task.ID, Value: result}
				}
			case <-w.QuitChan:
				return
			}
		}
	}()
}

// Stop stops the worker.
func (w Worker[T, R]) Stop() {
	go func() {
		w.QuitChan <- true
	}()
}

// NewThreadPool creates a new thread pool.
func NewThreadPool[T any, R any](numWorkers int, buffer int) *ThreadPool[T, R] {
	taskQueue := make(chan Task[T, R], buffer) // Buffered channel to hold a large queue of tasks
	resultChan := make(chan Result[R], buffer)
	errChan := make(chan error, buffer)
	workerPool := make(chan chan Task[T, R], numWorkers)
	quitChan := make(chan bool)

	workers := make([]Worker[T, R], numWorkers)
	for i := 0; i < numWorkers; i++ {
		workers[i] = NewWorker(i, workerPool, resultChan, errChan)
	}

	return &ThreadPool[T, R]{
		TaskQueue:  taskQueue,
		ResultChan: resultChan,
		ErrChan:    errChan,
		WorkerPool: workerPool,
		Workers:    workers,
		QuitChan:   quitChan,
	}
}

// Start starts the thread pool.
func (tp *ThreadPool[T, R]) Start() {
	for _, worker := range tp.Workers {
		tp.wg.Add(1)
		worker.Start(&tp.wg)
	}

	go func() {
		for {
			select {
			case task := <-tp.TaskQueue:
				// Get a worker's task queue from the worker pool.
				taskQueue := <-tp.WorkerPool
				taskQueue <- task
			case <-tp.QuitChan:
				for _, worker := range tp.Workers {
					worker.Stop()
				}
				return
			}
		}
	}()
}

// Stop stops the thread pool.
func (tp *ThreadPool[T, R]) Stop() {
	tp.QuitChan <- true
	tp.wg.Wait()
	close(tp.ResultChan)
	close(tp.ErrChan)
}
