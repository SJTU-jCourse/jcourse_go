package base

import (
	"context"
	"time"
)

type RedisConfig struct {
	Host     string
	Port     string
	Password string
}

type TaskId = string

type ResultWriter interface {
	Write(data []byte) (n int, err error)
}

type Task interface {
	Type() string
	Payload() []byte
	ResultWriter() ResultWriter
}

// type TaskPeriodic interface {
// 	Task
// 	ScheduleInterval() time.Duration
// }

type TaskOption struct {
	WithQueue *string
}

type TaskHandler = func(ctx context.Context, t Task) error

type TaskInterval = time.Duration

type PeriodicTaskId = string

// ============================ Implement ============================ //

// type taskPeriodic struct {
// 	task
// 	scheduleInterval time.Duration
// }

// func (t *taskPeriodic) ScheduleInterval() time.Duration {
// 	return t.scheduleInterval
// }

// func NewTaskPeriodic(taskType string, payload interface{}, duration time.Duration, options ...TaskOption) TaskPeriodic {
// 	return &taskPeriodic{
// 		task:             task{taskType: taskType, payload: payload, options: options},
// 		scheduleInterval: duration,
// 	}
// }

// ============================ End ============================ //
