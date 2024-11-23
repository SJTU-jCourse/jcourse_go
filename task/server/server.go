package main

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/hibiken/asynq"

	"jcourse_go/task/biz/ping"
	"jcourse_go/task/biz/statistic"
	"jcourse_go/util"
)

func postEnqueueFunc(info *asynq.TaskInfo, err error) {
	if err != nil {
		log.Printf("Asynq Enqueue Error: %v", err)
	}
}

var server *asynq.Server
var scheduler *asynq.Scheduler

const (
	ScheduleQueues = "periodic"
)

func StartScheduler(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	host := util.GetRedisHost()
	port := util.GetRedisPort()
	redisOpt := asynq.RedisClientOpt{Addr: fmt.Sprintf("%s:%s", host, port), Password: util.GetRedisPassword()}
	sched := asynq.NewScheduler(redisOpt, nil)
	// ... Register tasks.
	saveStatisticTask := asynq.NewTask(statistic.TypeSaveDailyStatistic, nil)
	entryID, err := sched.Register(statistic.IntervalSaveDailyStatistic, saveStatisticTask, asynq.Queue(ScheduleQueues))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("registered an entry: %q\n", entryID)

	refreshPVDupJudgeTask := asynq.NewTask(statistic.TypeRefreshPVDupJudge, nil)
	entryID, err = sched.Register(statistic.IntervalRefreshPVDupJudge, refreshPVDupJudgeTask, asynq.Queue(ScheduleQueues))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("registered an entry: %q\n", entryID)

	go func() {
		if err := sched.Run(); err != nil {
			log.Fatal(err)
		}
		defer wg.Done()
	}()
	scheduler = sched
}

func StartServer(ctx context.Context, wg *sync.WaitGroup) {
	host := util.GetRedisHost()
	port := util.GetRedisPort()
	redisOpt := asynq.RedisClientOpt{Addr: fmt.Sprintf("%s:%s", host, port), Password: util.GetRedisPassword()}
	srv := asynq.NewServer(
		redisOpt,
		asynq.Config{
			// Specify how many concurrent workers to use
			Concurrency: 10,
			// Optionally specify multiple queues with different priority.
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
			// See the godoc for other configuration options
		},
	)
	// mux maps a type to a handler
	mux := asynq.NewServeMux()

	// ...register other handlers...
	mux.HandleFunc(statistic.TypeSaveDailyStatistic, statistic.HandleSaveStatisticTask)
	mux.HandleFunc(statistic.TypeRefreshPVDupJudge, statistic.HandleRefreshPVDupJudgeTask)
	mux.HandleFunc(ping.TypePing, ping.TaskPingHandler)
	// start server and scheduler
	go func() {
		if err := srv.Run(mux); err != nil {
			log.Fatalf("could not run server: %v", err)
		}
		defer wg.Done()
	}()
	server = srv
}
func main() {
	ctx := context.Background()
	var wg sync.WaitGroup
	wg.Add(2)
	StartServer(ctx, &wg)
	StartScheduler(ctx, &wg)
	wg.Wait()
}
