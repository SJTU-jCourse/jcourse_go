package asynq

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/hibiken/asynq"

	"jcourse_go/internal/interface/task/base"
)

const (
	DefaultQueue             = "default"
	StandalonePeriodicQueue  = "standalone_periodic"
	DistributedPeriodicQueue = "distributed_periodic"
)

var (
	usedQueues = map[string]int{
		DefaultQueue:             1,
		StandalonePeriodicQueue:  1,
		DistributedPeriodicQueue: 1,
	}
)

type AsynqTaskManager struct {
	client    *asynq.Client
	server    *asynq.Server
	scheduler *asynq.Scheduler
	mux       *asynq.ServeMux
	redisOpt  asynq.RedisClientOpt

	muxRegistry map[string]asynq.HandlerFunc
	mu          sync.Mutex
}

type rawTask = *asynq.Task
type adapterTask struct {
	rawTask
}

func (a *adapterTask) ResultWriter() base.ResultWriter {
	return a.rawTask.ResultWriter()
}

// ----------------------------------------------------------------------------------
// Creation & Registration
// ----------------------------------------------------------------------------------

func newTask(taskType string, payload []byte, opts ...base.TaskOption) base.Task {
	return &adapterTask{
		asynq.NewTask(taskType, payload, convertOptions(opts...)...),
	}
}

func (m *AsynqTaskManager) CreateTask(taskType string, payload []byte, options ...base.TaskOption) base.Task {
	return newTask(taskType, payload, options...)
}

func (m *AsynqTaskManager) RegisterTaskHandler(taskType string, handler base.TaskHandler) error {
	asynqHandler := func(ctx context.Context, t *asynq.Task) error {
		return handler(ctx, &adapterTask{t})
	}

	m.mu.Lock()
	m.muxRegistry[taskType] = asynqHandler
	m.mu.Unlock()

	return nil
}

// ----------------------------------------------------------------------------------
// Queueing
// ----------------------------------------------------------------------------------

func (m *AsynqTaskManager) Enqueue(task base.Task, opts ...base.TaskOption) error {
	asynqTask, ok := task.(*adapterTask)
	if !ok {
		return errors.New("task is not an instance of asynq.Task")
	}
	_, err := m.client.Enqueue(asynqTask.rawTask, convertOptions(opts...)...)
	return err
}

// ----------------------------------------------------------------------------------
// Initialization
// ----------------------------------------------------------------------------------

// NewAsynqTaskManager creates and returns a new AsynqTaskManager, but does not start anything yet.
func NewAsynqTaskManager(redisConfig base.RedisConfig) *AsynqTaskManager {
	redisOpt := asynq.RedisClientOpt{
		Addr:     redisConfig.DSN,
		Password: redisConfig.Password,
	}

	client := asynq.NewClient(redisOpt)
	mux := asynq.NewServeMux()

	return &AsynqTaskManager{
		client:      client,
		redisOpt:    redisOpt,
		mux:         mux,
		muxRegistry: make(map[string]asynq.HandlerFunc),
	}
}

// RunServer starts the main Asynq server and scheduler concurrently.
// If either fails, the error is returned immediately.
func (m *AsynqTaskManager) RunServer() error {
	m.server = asynq.NewServer(
		m.redisOpt,
		asynq.Config{
			Concurrency: 10,
			Queues:      usedQueues,
		},
	)

	m.mu.Lock()
	for pattern, handler := range m.muxRegistry {
		m.mux.HandleFunc(pattern, handler)
	}
	m.mu.Unlock()

	m.scheduler = asynq.NewScheduler(m.redisOpt, nil)

	var wg sync.WaitGroup
	wg.Add(2)
	fail := make(chan error, 2)

	// Start server
	go func() {
		wg.Done()
		log.Println("[AsynqTaskManager] Server is starting...")
		if err := m.server.Run(m.mux); err != nil {
			log.Printf("[AsynqTaskManager] Could not run server: %v", err)
			fail <- err
		}
	}()

	// Start scheduler
	go func() {
		wg.Done()
		log.Println("[AsynqTaskManager] Scheduler is starting...")
		if err := m.scheduler.Run(); err != nil {
			log.Printf("[AsynqTaskManager] Could not run scheduler: %v", err)
			fail <- err
		}
	}()

	wg.Wait()

	select {
	case err := <-fail:
		return err
	case <-time.After(2 * time.Second):
		// If no error after a little wait, initialization is considered successful
		return nil
	}
}

// ----------------------------------------------------------------------------------
// Periodic tasks
// Only For Standalone Env, i.e. only one instance
// ----------------------------------------------------------------------------------

func (m *AsynqTaskManager) Submit(interval base.TaskInterval, task base.Task) (base.PeriodicTaskId, error) {
	if m.scheduler == nil {
		return "", errors.New("scheduler is not initialized")
	}

	asynqTask, ok := task.(*adapterTask)
	if !ok {
		return "", errors.New("task is not an instance of asynq.Task")
	}

	intervalStr := convertInterval(interval)
	entryID, err := m.scheduler.Register(intervalStr, asynqTask.rawTask, asynq.Queue(StandalonePeriodicQueue))
	if err != nil {
		return "", err
	}
	return entryID, nil
}

func (m *AsynqTaskManager) Kill(id base.PeriodicTaskId) error {
	if m.scheduler == nil {
		return errors.New("scheduler is not initialized")
	}
	return m.scheduler.Unregister(id)
}

// ----------------------------------------------------------------------------------
// Graceful Shutdown
// ----------------------------------------------------------------------------------

func (m *AsynqTaskManager) Shutdown() error {
	log.Println("[AsynqTaskManager] Shutting down AsynqTaskManager...")

	if m.scheduler != nil {
		m.scheduler.Shutdown()
		log.Println("[AsynqTaskManager] Scheduler shutdown complete.")
	}

	if m.server != nil {
		m.server.Shutdown()
		log.Println("[AsynqTaskManager] Server shutdown complete.")
	}

	log.Println("[AsynqTaskManager] Shutdown sequence complete.")
	return nil
}

// ----------------------------------------------------------------------------------
// Utilities
// ----------------------------------------------------------------------------------

func (m *AsynqTaskManager) HealthCheck() error {
	// TODO: Use asynq.Inspector for additional health checks
	return nil
}

func convertOptions(options ...base.TaskOption) []asynq.Option {
	var asynqOptions []asynq.Option
	for _, opt := range options {
		if opt.WithQueue != nil {
			asynqOptions = append(asynqOptions, asynq.Queue(*opt.WithQueue))
		}
	}
	return asynqOptions
}

func convertInterval(interval base.TaskInterval) string {
	return fmt.Sprintf("@every %s", interval.String())
}
