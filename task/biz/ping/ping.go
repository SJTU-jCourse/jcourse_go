package ping

import (
	"context"
	"jcourse_go/task"
	"log"
)

const (
	TypePing = "test:ping"
)

func Init() {
	task.TaskManager.RegisterTaskHandler(TypePing, TaskPingHandler)
}

func TaskPingHandler(ctx context.Context, t task.Task) error {
	log.Printf("pong for task: %s", t.Type())
	_, err := t.ResultWriter().Write([]byte("pong"))
	if err != nil {
		return err
	}
	return nil
}
