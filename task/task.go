package task

import (
	"context"
	"io"
	"time"
)

type TaskId = string
type TaskOption struct {
	WithQueue string
}

type Task interface {
	Type() string
	Payload() interface{}
	Options() []TaskOption
	ResultWriter() io.Writer
}

type TaskPeriodic interface {
	Task
	ScheduleInterval() time.Duration
}

type TaskHandler func(ctx context.Context, t Task) error

var TaskManager IAsyncTaskManager

// IAsyncTaskManager defines an interface for asynchronous task management without exposing implementation details.
type IAsyncTaskManager interface {
	ITaskFactory

	// RegisterTaskHandler registers a task handler for a task type.
	ITaskHandlerRegister

	IOneTimeTaskRunner
	IScheduleTaskRunner
	// StartServer registers task types and handlers, then starts the server to process tasks.
	StartServer() error
}

type ITaskFactory interface {
	CreateTask(taskType string, payload interface{}, options ...TaskOption) Task
	CreatePeriodicTask(taskType string, payload interface{}, duration time.Duration, options ...TaskOption) TaskPeriodic
}

type ITaskHandlerRegister interface {
	RegisterTaskHandler(taskType string, handler TaskHandler) error
}

type IOneTimeTaskRunner interface {
	Enqueue(task Task) error
}

type IScheduleTaskRunner interface {
	// Submit schedules a new task.
	Submit(task TaskPeriodic) (TaskId, error)

	// Kill removes a scheduled task.
	Kill(taskID TaskId) error
}
