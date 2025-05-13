package task

import (
	"jcourse_go/pkg/task/base"
)

type IAsyncTaskManager interface {
	ITaskFactory

	ITaskHandlerRegister

	IOneTimeTaskRunner
	IScheduleTaskRunner

	HealthChecker

	RunServer() error

	Shutdown() error
}

type ITaskFactory interface {
	CreateTask(taskType string, payload []byte, options ...base.TaskOption) base.Task
}

type ITaskHandlerRegister interface {
	RegisterTaskHandler(taskType string, handler base.TaskHandler) error
}

type IOneTimeTaskRunner interface {
	Enqueue(task base.Task, opts ...base.TaskOption) error
}

type IScheduleTaskRunner interface {
	Submit(interval base.TaskInterval, task base.Task) (base.PeriodicTaskId, error)

	Kill(taskID base.PeriodicTaskId) error
}

type HealthChecker interface {
	HealthCheck() error
}
