package asynq

// import (
// 	"context"
// 	"log"
// 	"testing"
// 	"time"

// 	"jcourse_go/task/base"

// 	"github.com/redis/go-redis/v9"
// 	"github.com/stretchr/testify/assert"
// )

// func mockTaskHandler(ctx context.Context, t base.Task, executionChan chan<- string) error {
// 	executionChan <- t.Type()
// 	return nil
// }

// func setupRedisClient() *redis.Client {
// 	return redis.NewClient(&redis.Options{
// 		Addr:     "localhost:6379", // Ensure Redis is running locally
// 		Password: "",               // Update if Redis requires a password
// 	})
// }

// func cleanupRedis(ctx context.Context, rdb *redis.Client, keys ...string) {
// 	for _, key := range keys {
// 		err := rdb.Del(ctx, key).Err()
// 		if err != nil {
// 			log.Printf("Failed to delete key %s: %v", key, err)
// 		}
// 	}
// }

// func TestDistributedAsynqTaskManager_Submit(t *testing.T) {
// 	ctx := context.Background()
// 	rdb1 := setupRedisClient()
// 	rdb2 := setupRedisClient()

// 	taskType := "test:distributed_task"

// 	executionChan := make(chan string, 100)

// 	asynqManager1 := NewAsynqTaskManager(base.RedisConfig{
// 		DSN:      "localhost:6379",
// 		Password: "",
// 	})

// 	err := asynqManager1.RegisterTaskHandler(taskType, func(ctx context.Context, t base.Task) error {
// 		return mockTaskHandler(ctx, t, executionChan)
// 	})
// 	assert.Nil(t, err, "RegisterTaskHandler should not return an error")

// 	go func() {
// 		err := asynqManager1.RunServer()
// 		if err != nil {
// 			log.Fatalf("AsynqTaskManager1 RunServer failed: %v", err)
// 		}
// 	}()

// 	distManager1 := NewDistributedAsynqTaskManager(asynqManager1, rdb1)
// 	defer distManager1.Shutdown()

// 	asynqManager2 := NewAsynqTaskManager(base.RedisConfig{
// 		DSN:      "localhost:6379",
// 		Password: "",
// 	})

// 	err = asynqManager2.RegisterTaskHandler(taskType, func(ctx context.Context, t base.Task) error {
// 		return mockTaskHandler(ctx, t, executionChan)
// 	})
// 	assert.Nil(t, err, "RegisterTaskHandler should not return an error")

// 	go func() {
// 		err := asynqManager2.RunServer()
// 		if err != nil {
// 			log.Fatalf("AsynqTaskManager2 RunServer failed: %v", err)
// 		}
// 	}()

// 	distManager2 := NewDistributedAsynqTaskManager(asynqManager2, rdb2)
// 	defer distManager2.Shutdown()

// 	defer cleanupRedis(ctx, rdb1, "dist_lock:"+taskType, "task_kill_channel")
// 	defer cleanupRedis(ctx, rdb2, "dist_lock:"+taskType, "task_kill_channel")

// 	jobID1, err := distManager1.Submit(2*time.Second, asynqManager1.CreateTask(taskType, []byte(`{"data":"test1"}`)))
// 	assert.Nil(t, err, "Submit should not return an error")
// 	assert.NotEmpty(t, jobID1, "Job ID should not be empty")

// 	jobID2, err := distManager2.Submit(2*time.Second, asynqManager2.CreateTask(taskType, []byte(`{"data":"test2"}`)))
// 	assert.Nil(t, err, "Submit should not return an error")
// 	assert.NotEmpty(t, jobID2, "Job ID should not be empty")

// 	time.Sleep(10 * time.Second)

// 	assert.Equal(t, 5, len(executionChan), "Only one task should have been executed")
// }

// func TestDistributedAsynqTaskManager_Kill(t *testing.T) {
// 	ctx := context.Background()
// 	rdb := setupRedisClient()
// 	defer cleanupRedis(ctx, rdb, "dist_lock:test:kill_task", "task_kill_channel")

// 	executionChan := make(chan string, 100)
// 	asynqManager := NewAsynqTaskManager(base.RedisConfig{
// 		DSN:      "localhost:6379",
// 		Password: "",
// 	})

// 	// Register a handler for the test task
// 	taskType := "test:kill_task"
// 	err := asynqManager.RegisterTaskHandler(taskType, func(ctx context.Context, t base.Task) error {
// 		executionChan <- t.Type()
// 		return nil
// 	})
// 	assert.Nil(t, err, "RegisterTaskHandler should not return an error")

// 	// Run the server
// 	go func() {
// 		err := asynqManager.RunServer()
// 		if err != nil {
// 			log.Fatalf("AsynqTaskManager RunServer failed: %v", err)
// 		}
// 	}()
// 	// defer asynqManager.Shutdown()

// 	// Distributed manager
// 	distManager := NewDistributedAsynqTaskManager(asynqManager, rdb)
// 	defer distManager.Shutdown()

// 	// Submit a periodic task with an interval of 2s
// 	jobID, err := distManager.Submit(2*time.Second, asynqManager.CreateTask(taskType, []byte(`{"kill_test":"1"}`)))
// 	assert.Nil(t, err, "Submit should not return an error")
// 	assert.NotEmpty(t, jobID, "Job ID should not be empty")

// 	// Wait ~5s to allow at least a couple of executions
// 	time.Sleep(5 * time.Second)
// 	firstCount := len(executionChan)
// 	assert.Greater(t, firstCount, 0, "Task should have executed at least once before kill")

// 	// Kill the job
// 	err = distManager.Kill(jobID)
// 	assert.Nil(t, err, "Kill should not return an error")
// 	// Clear current executions seen so far
// 	// (Not strictly necessary, but helps check fresh executions below.)
// 	for len(executionChan) > 0 {
// 		<-executionChan
// 	}

// 	// Wait again for ~5s to see if new executions come in
// 	time.Sleep(5 * time.Second)
// 	secondCount := len(executionChan)
// 	assert.Equal(t, 0, secondCount, "No new executions should appear after Kill")
// }
