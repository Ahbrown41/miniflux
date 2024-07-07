package pool

import (
	"sync"
	"testing"
	"time"
)

// Example task function for testing
func testTaskFunction(data int) (int, error) {
	time.Sleep(10 * time.Millisecond) // Simulate some work
	result := data * data             // Example task: square the number
	return result, nil
}

func TestThreadPool(t *testing.T) {
	numWorkers := 10
	numTasks := 300
	buffer := 200

	threadPool := NewThreadPool[int, int](numWorkers, buffer)
	threadPool.Start()

	var taskWG sync.WaitGroup
	taskWG.Add(numTasks)

	results := make(map[int]int)
	resultsMutex := sync.Mutex{}
	errors := make([]error, 0)

	// Create tasks and add them to the task queue.
	for i := 1; i <= numTasks; i++ {
		task := Task[int, int]{
			ID:       i,
			Data:     i,
			Function: testTaskFunction,
		}
		threadPool.TaskQueue <- task
	}

	// Collect results.
	go func() {
		for result := range threadPool.ResultChan {
			resultsMutex.Lock()
			results[result.ID] = result.Value
			resultsMutex.Unlock()
			taskWG.Done()
		}
	}()

	// Collect errors.
	go func() {
		for err := range threadPool.ErrChan {
			resultsMutex.Lock()
			errors = append(errors, err)
			resultsMutex.Unlock()
			taskWG.Done()
		}
	}()

	// Wait for all tasks to complete.
	taskWG.Wait()

	// Stop the thread pool.
	threadPool.Stop()

	// Check results
	if len(errors) != 0 {
		t.Fatalf("Expected no errors, but got %d errors: %v", len(errors), errors)
	}

	if len(results) != numTasks {
		t.Fatalf("Expected %d results, but got %d", numTasks, len(results))
	}

	for i := 1; i <= numTasks; i++ {
		expected := i * i
		if results[i] != expected {
			t.Errorf("Task %d: expected result %d, but got %d", i, expected, results[i])
		}
	}
}
