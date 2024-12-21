package asynq

import (
	"context"
	"errors"
	"fmt"
	"jcourse_go/task/base"
	"log"
	"sync"

	"github.com/hibiken/asynq"
)

const (
	DefaultQueue  = "default"
	PeriodicQueue = "periodic"
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

func (m *AsynqTaskManager) Enqueue(task base.Task) error {
	asynqTask, ok := task.(*adapterTask)
	if !ok {
		return errors.New("task is not an instance of asynq.Task")
	}
	_, err := m.client.Enqueue(asynqTask.rawTask)
	return err
}

func NewAsynqTaskManager(redisConfig base.RedisConfig) *AsynqTaskManager {
	redisOpt := asynq.RedisClientOpt{
		Addr:     redisConfig.DSN,
		Password: redisConfig.Password,
	}

	client := asynq.NewClient(redisOpt)

	mux := asynq.NewServeMux()

	return &AsynqTaskManager{
		// server:      nil,
		// scheduler:   nil,
		client:      client,
		redisOpt:    redisOpt,
		mux:         mux,
		muxRegistry: make(map[string]asynq.HandlerFunc),
	}
}

// must call RunServer first
func (m *AsynqTaskManager) Submit(interval base.TaskInterval, task base.Task) (base.PeriodicTaskId, error) {
	if m.scheduler == nil {
		return "", errors.New("scheduler is not initialized")
	}

	asynqTask, ok := task.(*adapterTask)
	if !ok {
		return "", errors.New("task is not an instance of asynq.Task")
	}

	intervalStr := convertInterval(interval)
	entryID, err := m.scheduler.Register(intervalStr, asynqTask.rawTask, asynq.Queue(PeriodicQueue))
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

func (m *AsynqTaskManager) RunServer() error {
	m.server = asynq.NewServer(
		m.redisOpt,
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				// TODO @huangjunqing configable queue
				DefaultQueue:  1,
				PeriodicQueue: 1,
			},
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

	go func() {
		wg.Done()
		log.Println("server is starting")
		if err := m.server.Run(m.mux); err != nil {
			log.Fatalf("could not run server: %v", err)
		}
	}()

	go func() {
		wg.Done()
		log.Println("scheduler is starting")
		if err := m.scheduler.Run(); err != nil {
			log.Fatalf("could not run scheduler: %v", err)
		}
	}()

	wg.Wait()

	return nil
}

func (m *AsynqTaskManager) HealthCheck() error {
	// TODO @huangjunqing using asynq's Inspector to implement health check
	return nil
}

func convertOptions(options ...base.TaskOption) []asynq.Option {
	var asynqOptions []asynq.Option
	for _, opt := range options {
		if opt.WithQueue != nil {
			// TODO @huangjunqing add more converter
			asynqOptions = append(asynqOptions, asynq.Queue(*opt.WithQueue))
		}
	}
	return asynqOptions
}

// Start of Selection
/*
24 * time.Hour //"@every 24h"
1 * time.Hour  //"@every 1h"
*/
func convertInterval(interval base.TaskInterval) string {
	return fmt.Sprintf("@every %s", interval.String())
}
