package async_task

import "context"

type Type string

const (
	TypeUserCreated   Type = "user.created"
	TypeReviewCreated Type = "review.created"
)

type AsyncTask struct {
	ID      string
	Type    Type
	Payload any
}

type TaskQueue interface {
	Enqueue(ctx context.Context, t AsyncTask) error
}
