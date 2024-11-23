package impl

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/hibiken/asynq"
)

type redisConfig struct {
	Host     string
	Port     string
	Password string
}

// AsynqTaskManager implements the IAsyncTaskManager interface using Asynq.
type AsynqTaskManager struct {
	client    *asynq.Client
	server    *asynq.Server
	scheduler *asynq.Scheduler
	// prepared fields
	mux      *asynq.ServeMux
	redisOpt asynq.RedisClientOpt
}

// NewAsynqTaskManager creates a new instance of AsynqTaskManager with the provided Redis configuration.
func NewAsynqTaskManager(redisConfig redisConfig) *AsynqTaskManager {
	redisAddr := fmt.Sprintf("%s:%s", redisConfig.Host, redisConfig.Port)

	client := asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})

	redisOpt := asynq.RedisClientOpt{
		Addr:     redisAddr,
		Password: redisConfig.Password,
	}

	mux := asynq.NewServeMux()
	return &AsynqTaskManager{
		client:   client,
		mux:      mux,
		redisOpt: redisOpt,
	}
}

// Enqueue enqueues a task based on the scheduleType.
// scheduleType: "periodic" or "one-time"
func (m *AsynqTaskManager) Enqueue(taskType string, payload interface{}, options ...interface{}) error {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	task := asynq.NewTask(taskType, payloadBytes)
	_, err = m.client.Enqueue(task, convertOptions(options...)...)
	return err
}

// StartServer initializes and starts the Asynq server and scheduler.
func (m *AsynqTaskManager) StartServer() error {
	// Initialize Server
	m.server = asynq.NewServer(
		m.redisOpt,
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
		},
	)

	// Register Handlers
	for pattern, handler := range muxRegistry {
		m.mux.HandleFunc(pattern, handler)
	}

	// Initialize Scheduler
	m.scheduler = asynq.NewScheduler(m.redisOpt, nil)

	var wg sync.WaitGroup
	wg.Add(1)

	// Start Server
	go func() {
		defer wg.Done()
		if err := m.server.Run(m.mux); err != nil {
			log.Fatalf("could not run server: %v", err)
		}
	}()

	// wait for server to start to register periodic tasks
	wg.Wait()

	wg.Add(1)
	// Start Scheduler
	go func() {
		defer wg.Done()
		for _, task := range deamonTaskRegistry {
			_, err := m.scheduler.Register(task.cronspec, task.task, task.opts...) // TODO @huangjunqing collect taskId for cancellation
			if err != nil {
				return
			}
		}
		if err := m.scheduler.Run(); err != nil {
			log.Fatal(err)
		}
	}()

	// Wait for both server and scheduler to finish
	wg.Wait()

	return nil
}

// convertOptions converts generic interface options to Asynq options.
// This is a helper function to handle variadic options.
func convertOptions(options ...interface{}) []asynq.Option {
	var asynqOptions []asynq.Option
	for _, opt := range options {
		if ao, ok := opt.(asynq.Option); ok {
			asynqOptions = append(asynqOptions, ao)
		}
	}
	return asynqOptions
}
