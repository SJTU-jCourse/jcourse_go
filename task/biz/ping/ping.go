package ping

import (
	"context"
	"jcourse_go/task/base"
	"log"
)

const (
	TypePing = "test:ping"
)

func TaskPingHandler(ctx context.Context, t base.Task) error {
	log.Printf("pong for task: %s", t.Type())
	_, err := t.ResultWriter().Write([]byte("pong"))
	if err != nil {
		return err
	}
	return nil
}
