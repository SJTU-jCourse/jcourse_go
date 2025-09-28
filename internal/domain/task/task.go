package task

import "context"

type AsyncTask struct {
	Type    Type
	Payload any
}

type TaskQueue interface {
	Enqueue(ctx context.Context, t AsyncTask) error
}
