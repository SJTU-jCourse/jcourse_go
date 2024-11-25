package asynq

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"sync"

	"jcourse_go/task"

	"github.com/hibiken/asynq"
)

// AsynqTaskManager implements the IAsyncTaskManager interface using Asynq.
type AsynqTaskManager struct {
	client    *asynq.Client
	server    *asynq.Server
	scheduler *asynq.Scheduler
	mux       *asynq.ServeMux
	redisOpt  asynq.RedisClientOpt

	// Internal registries for handlers and daemon tasks
	muxRegistry        map[string]asynq.HandlerFunc
	daemonTaskRegistry []daemonTask
	mu                 sync.Mutex
}

// daemonTask represents a scheduled daemon task.
type daemonTask struct {
	cronspec string
	task     *asynq.Task
	opts     []asynq.Option
}

// NewAsynqTaskManager creates a new instance of AsynqTaskManager with the provided Redis configuration.
func NewAsynqTaskManager(redisConfig task.RedisConfig) *AsynqTaskManager {
	redisAddr := fmt.Sprintf("%s:%s", redisConfig.Host, redisConfig.Port)

	client := asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})

	redisOpt := asynq.RedisClientOpt{
		Addr:     redisAddr,
		Password: redisConfig.Password,
	}

	mux := asynq.NewServeMux()

	return &AsynqTaskManager{
		client:             client,
		mux:                mux,
		redisOpt:           redisOpt,
		muxRegistry:        make(map[string]asynq.HandlerFunc),
		daemonTaskRegistry: []daemonTask{},
	}
}

// RegisterTaskHandler registers a task handler for a given task type.
func (m *AsynqTaskManager) RegisterTaskHandler(taskType string, handler task.TaskHandler) error {
	asynqHandler := func(ctx context.Context, t *asynq.Task) error {
		// Wrap the Asynq Task interface to match the task.Task interface
		wrappedTask := &asynqTask{t}
		return handler(ctx, wrappedTask)
	}

	m.muxRegistry[taskType] = asynqHandler
	m.mux.HandleFunc(taskType, asynqHandler)
	return nil
}

// Enqueue enqueues a one-time task.
func (m *AsynqTaskManager) Enqueue(task task.Task) error {
	payloadBytes, err := json.Marshal(task.Payload())
	if err != nil {
		return err
	}
	asynqTask := asynq.NewTask(task.Type(), payloadBytes)
	_, err = m.client.Enqueue(asynqTask, convertOptions(task.Options()...)...)
	return err
}

// Submit schedules a periodic task.
func (m *AsynqTaskManager) Submit(task task.TaskPeriodic) (task.TaskId, error) {
	payloadBytes, err := json.Marshal(task.Payload())
	if err != nil {
		return "", err
	}
	asynqTask := asynq.NewTask(task.Type(), payloadBytes)
	taskID, err := m.scheduler.Register(task.ScheduleInterval(), asynqTask, convertOptions(task.Options()...)...)
	if err != nil {
		return "", err
	}
	return task.TaskId(taskID), nil
}

// Kill removes a scheduled task.
func (m *AsynqTaskManager) Kill(taskID task.TaskId) error {
	err := m.scheduler.Cancel(string(taskID))
	if err != nil {
		return err
	}
	return nil
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
	for pattern, handler := range m.muxRegistry {
		m.mux.HandleFunc(pattern, handler)
	}

	// Initialize Scheduler
	m.scheduler = asynq.NewScheduler(m.redisOpt, nil)

	var wg sync.WaitGroup
	wg.Add(2)

	// Start Server
	go func() {
		defer wg.Done()
		if err := m.server.Run(m.mux); err != nil {
			log.Fatalf("could not run server: %v", err)
		}
	}()

	// Start Scheduler
	go func() {
		defer wg.Done()
		if err := m.scheduler.Run(); err != nil {
			log.Fatalf("could not run scheduler: %v", err)
		}
	}()

	// Wait for both server and scheduler to finish
	wg.Wait()

	return nil
}

// RegisterDaemonTask registers a daemon task to be scheduled.
func (m *AsynqTaskManager) RegisterDaemonTask(cronspec string, t *asynq.Task, opts ...asynq.Option) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.daemonTaskRegistry = append(m.daemonTaskRegistry, daemonTask{
		cronspec: cronspec,
		task:     t,
		opts:     opts,
	})
}

// InitializeDaemonTasks registers all daemon tasks with the scheduler.
func (m *AsynqTaskManager) InitializeDaemonTasks() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, dt := range m.daemonTaskRegistry {
		_, err := m.scheduler.Register(dt.cronspec, dt.task, dt.opts...)
		if err != nil {
			return fmt.Errorf("failed to register daemon task %s: %w", dt.task.Type(), err)
		}
	}
	return nil
}

// convertOptions converts generic TaskOption interfaces to Asynq options.
func convertOptions(options ...task.TaskOption) []asynq.Option {
	var asynqOptions []asynq.Option
	for _, opt := range options {
		if ao, ok := opt.(asynq.Option); ok {
			asynqOptions = append(asynqOptions, ao)
		}
	}
	return asynqOptions
}

// asynqTask is a wrapper to adapt asynq.Task to the task.Task interface.
type asynqTask struct {
	*asynq.Task
}

// Type returns the task type.
func (t *asynqTask) Type() string {
	return t.Task.Type()
}

// Payload returns the task payload.
func (t *asynqTask) Payload() interface{} {
	var payload interface{}
	if err := json.Unmarshal(t.Task.Payload(), &payload); err != nil {
		return nil
	}
	return payload
}

// Options returns the task options.
func (t *asynqTask) Options() []task.TaskOption {
	// Asynq does not expose options after task creation
	return nil
}

// ResultWriter returns the result writer.
func (t *asynqTask) ResultWriter() io.Writer {
	// Implement if needed
	return nil
}
