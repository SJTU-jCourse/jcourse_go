package task

import (
	"github.com/hibiken/asynq"

	"jcourse_go/internal/app"
	"jcourse_go/internal/domain/task"
	"jcourse_go/internal/interface/task/handler"
)

func NewAsyncTaskServer(s *app.ServiceContainer) (*asynq.Server, *asynq.ServeMux) {
	srv := asynq.NewServerFromRedisClient(s.Redis, asynq.Config{})
	reviewAntiSpamHandler := handler.NewReviewAntiSpamHandler()
	mux := asynq.NewServeMux()
	mux.Handle(string(task.TypeReviewAntiSpam), reviewAntiSpamHandler)
	return srv, mux
}
