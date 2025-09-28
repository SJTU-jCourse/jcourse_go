package handler

import (
	"context"

	"github.com/hibiken/asynq"
)

type ReviewAntiSpamHandler struct {
}

func NewReviewAntiSpamHandler() asynq.Handler {
	return &ReviewAntiSpamHandler{}
}

func (h *ReviewAntiSpamHandler) ProcessTask(ctx context.Context, t *asynq.Task) error {
	return nil
}
