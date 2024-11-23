package ping

import (
	"context"
	"log"

	"github.com/hibiken/asynq"
)

const (
	TypePing = "test:ping"
)

func TaskPingHandler(ctx context.Context, t *asynq.Task) error {
	log.Printf("pong for task: %s", t.Type())
	_, err := t.ResultWriter().Write([]byte("pong"))
	if err != nil {
		return err
	}
	return nil
}
