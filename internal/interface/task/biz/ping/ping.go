package ping

import (
	"context"
	"log"
	"time"

	"jcourse_go/internal/interface/task/base"
)

const (
	TypePing     = "test:ping"
	IntervalPing = 2 * time.Second // "@every 2s"
)

func TaskPingHandler(ctx context.Context, t base.Task) error {
	log.Printf("pong for task: %s", t.Type())
	_, err := t.ResultWriter().Write([]byte("pong"))
	if err != nil {
		return err
	}
	return nil
}
