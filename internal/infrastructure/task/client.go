package task

import (
	"context"
	"encoding/json"

	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"

	"jcourse_go/internal/domain/task"
)

type TaskQueueClient struct {
	client *asynq.Client
}

func NewTaskQueueClient(rdb *redis.Client) *TaskQueueClient {
	client := asynq.NewClientFromRedisClient(rdb)
	return &TaskQueueClient{client: client}
}

func (c *TaskQueueClient) Enqueue(ctx context.Context, t task.AsyncTask) error {
	payload, err := json.Marshal(t)
	if err != nil {
		return err
	}
	asynqTask := asynq.NewTask(string(t.Type), payload)
	_, err = c.client.Enqueue(asynqTask)
	if err != nil {
		return err
	}
	return nil
}
