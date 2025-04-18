package task

import (
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/errgroup"

	"jcourse_go/task/asynq"
	"jcourse_go/task/base"
	"jcourse_go/task/biz/statistic"
)

var (
	taskRegistry = map[string]base.TaskHandler{
		statistic.TypeSaveDailyStatistic: statistic.HandleSaveStatisticTask,
		// statistic.TypeRefreshPVDupJudge:  statistic.HandleRefreshPVDupJudgeTask,
		// ping.TypePing:                    ping.TaskPingHandler,
	}

	taskSchedules = map[time.Duration]string{
		statistic.IntervalSaveDailyStatistic: statistic.TypeSaveDailyStatistic,
		// statistic.IntervalRefreshPVDupJudge:  statistic.TypeRefreshPVDupJudge,
		// ping.IntervalPing:                    ping.TypePing,
	}

	GlobalTaskManager IAsyncTaskManager
)

func InitTaskManager(redisConfig base.RedisConfig) {
	taskManager := asynq.NewAsynqTaskManager(redisConfig)

	// 1. Register task handlers
	for taskType, handler := range taskRegistry {
		if err := taskManager.RegisterTaskHandler(taskType, handler); err != nil {
			panic(err)
		}
	}

	// 2. Run server
	if err := taskManager.RunServer(); err != nil {
		panic(err)
	}

	// 3. Create scheduler and submit periodic tasks using errgroup
	GlobalTaskManager = asynq.NewDistributedAsynqTaskManager(taskManager, redis.NewClient(&redis.Options{
		Addr:     redisConfig.DSN,
		Password: redisConfig.Password,
	}))

	var g errgroup.Group

	for interval, taskType := range taskSchedules {
		g.Go(func() error {
			_, err := GlobalTaskManager.Submit(interval, GlobalTaskManager.CreateTask(taskType, nil))
			return err
		})
	}

	if err := g.Wait(); err != nil {
		panic(err)
	}

	log.Println("[MustInitTaskManager] success!")
}
