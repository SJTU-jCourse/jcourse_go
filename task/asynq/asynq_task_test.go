package asynq

// import (
// 	"context"
// 	"log"
// 	"sync"
// 	"testing"
// 	"time"

// 	"jcourse_go/task/base"

// 	"github.com/stretchr/testify/assert"
// )

// TestAsynqTaskManager_Enqueue tests the Enqueue method of AsynqTaskManager using local Redis
// func TestAsynqTaskManager_Enqueue(t *testing.T) {
// 	// Initialize AsynqTaskManager with local Redis
// 	redisConfig := base.RedisConfig{
// 		DSN:      "localhost:6379",
// 		Password: "", // Update if your Redis requires a password
// 	}
// 	taskManager := NewAsynqTaskManager(redisConfig)

// 	// Channel to signal handler execution
// 	handlerCalled := make(chan bool, 1)

// 	// Register a test handler
// 	taskType := "test:task"
// 	err := taskManager.RegisterTaskHandler(taskType, func(ctx context.Context, t base.Task) error {
// 		handlerCalled <- true
// 		return nil
// 	})
// 	if err != nil {
// 		t.Fatalf("Failed to register task handler: %v", err)
// 	}

// 	// Run the server in a separate goroutine
// 	var wg sync.WaitGroup
// 	wg.Add(1)
// 	go func() {
// 		defer wg.Done()
// 		err := taskManager.RunServer()
// 		if err != nil {
// 			log.Fatalf("RunServer failed: %v", err)
// 		}
// 	}()

// 	wg.Wait()

// 	// Create and enqueue a task
// 	payload := []byte(`{"key":"value"}`)
// 	task := taskManager.CreateTask(taskType, payload)
// 	err = taskManager.Enqueue(task)
// 	assert.Nil(t, err, "Enqueue should not return an error")

// 	// Wait for the handler to be called
// 	select {
// 	case <-handlerCalled:
// 		t.Log("Handler was called, \n Test Success!!! \n")
// 	case <-time.After(1 * time.Second):
// 		t.Fatal("Handler was not called within the expected time")
// 	}

// 	// Shutdown the server gracefully
// 	taskManager.server.Shutdown()
// 	taskManager.scheduler.Shutdown()
// }

// // TestAsynqTaskManager_Submit tests the Submit method of AsynqTaskManager using local Redis
// func TestAsynqTaskManager_Submit(t *testing.T) {
// 	// Initialize AsynqTaskManager with local Redis
// 	redisConfig := base.RedisConfig{
// 		DSN:      "localhost:6379",
// 		Password: "", // Update if your Redis requires a password
// 	}
// 	taskManager := NewAsynqTaskManager(redisConfig)

// 	// WaitGroup to manage goroutines
// 	var wg sync.WaitGroup

// 	// Mutex and counter to track handler executions
// 	var mu sync.Mutex
// 	executionCount := 0
// 	expectedExecutions := 3

// 	// Register a test handler
// 	taskType := "test:periodic_task"
// 	err := taskManager.RegisterTaskHandler(taskType, func(ctx context.Context, t base.Task) error {
// 		mu.Lock()
// 		executionCount++
// 		mu.Unlock()

// 		// Optional: Log each execution
// 		log.Printf("Handler executed %d time(s)", executionCount)

// 		return nil
// 	})
// 	if err != nil {
// 		t.Fatalf("Failed to register task handler: %v", err)
// 	}

// 	// Run the server in a separate goroutine
// 	wg.Add(1)
// 	go func() {
// 		defer wg.Done()
// 		err := taskManager.RunServer()
// 		if err != nil {
// 			log.Fatalf("RunServer failed: %v", err)
// 		}
// 	}()

// 	wg.Wait()

// 	// Submit a periodic task with a short interval
// 	interval := time.Second * 1
// 	task := taskManager.CreateTask(taskType, nil)
// 	_, err = taskManager.Submit(interval, task)
// 	assert.Nil(t, err, "Submit should not return an error")

// 	// Function to wait for expected executions
// 	waitForExecutions := func(t *testing.T) {
// 		timeout := time.After(time.Duration(expectedExecutions+2) * time.Second)
// 		ticker := time.NewTicker(100 * time.Millisecond)
// 		defer ticker.Stop()

// 		for {
// 			select {
// 			case <-timeout:
// 				mu.Lock()
// 				currentCount := executionCount
// 				mu.Unlock()
// 				t.Fatalf("Expected handler to be called %d times, but was called %d times", expectedExecutions, currentCount)
// 			case <-ticker.C:
// 				mu.Lock()
// 				if executionCount >= expectedExecutions {
// 					mu.Unlock()
// 					return
// 				}
// 				mu.Unlock()
// 			}
// 		}
// 	}

// 	// Wait for the expected number of handler executions
// 	waitForExecutions(t)

// 	// Shutdown the scheduler and server gracefully
// 	taskManager.scheduler.Shutdown()
// 	taskManager.server.Shutdown()

// 	// Final assertion
// 	assert.Equal(t, expectedExecutions, executionCount, "Handler should be called the expected number of times")
// }

// // TestAsynqTaskManager_RegisterTaskHandler tests the RegisterTaskHandler method using local Redis
// func TestAsynqTaskManager_RegisterTaskHandler(t *testing.T) {
// 	// Initialize AsynqTaskManager with local Redis
// 	redisConfig := base.RedisConfig{
// 		DSN:      "localhost:6379",
// 		Password: "", // Update if your Redis requires a password
// 	}
// 	taskManager := NewAsynqTaskManager(redisConfig)

// 	// Register a test handler
// 	taskType := "test:register_handler"
// 	err := taskManager.RegisterTaskHandler(taskType, func(ctx context.Context, t base.Task) error {
// 		return nil
// 	})
// 	assert.Nil(t, err, "RegisterTaskHandler should not return an error")

// 	// Verify that the handler is registered
// 	taskManager.mu.Lock()
// 	defer taskManager.mu.Unlock()
// 	_, exists := taskManager.muxRegistry[taskType]
// 	assert.True(t, exists, "Handler should be registered in muxRegistry")
// }

// // TestAsynqTaskManager_CreateTask tests the CreateTask method without requiring Redis
// func TestAsynqTaskManager_CreateTask(t *testing.T) {
// 	// Initialize AsynqTaskManager without connecting to Redis
// 	redisConfig := base.RedisConfig{
// 		DSN:      "localhost:6379", // This won't be used in this test
// 		Password: "",
// 	}
// 	taskManager := NewAsynqTaskManager(redisConfig)

// 	// Create a task
// 	taskType := "test:create_task"
// 	payload := []byte(`{"foo":"bar"}`)
// 	task := taskManager.CreateTask(taskType, payload)

// 	// Verify the task
// 	assert.Equal(t, taskType, task.Type(), "Task type should match")
// 	assert.Equal(t, payload, task.Payload(), "Task payload should match")
// }
