package base

import (
	"context"
	"time"
)

type RedisConfig struct {
	DSN      string
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

type TaskOption struct {
	WithQueue *string
}

type TaskHandler = func(ctx context.Context, t Task) error

type TaskInterval = time.Duration

type PeriodicTaskId = string
