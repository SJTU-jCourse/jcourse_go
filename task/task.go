package task

import (
	"time"
)

type ScheduleType = string

const (
	ScheduleTypePeriodic = "periodic"
	ScheduleTypeOneTime  = "one-time"
)

// IAsyncTaskManager defines an interface for asynchronous task management without exposing implementation details.
type IAsyncTaskManager interface {
	// Task runner
	IOneTimeTaskRunner
	// IScheduleTaskRunner TODO @huangjunqing
	// StartServer registers task types and handlers, then starts the server to process tasks.
	StartServer() error
}

type IOneTimeTaskRunner interface {
	// Enqueue enqueues a task based on the scheduleType.
	// scheduleType: "one-time" only
	Enqueue(taskType string, payload interface{}, options ...interface{}) error
}

type IScheduleTaskRunner interface {
	// Submit schedules a new task.
	Submit(taskType string, payload interface{}, scheduleTime time.Time) (string, error)

	// Kill removes a scheduled task.
	Kill(taskID string) error
}
