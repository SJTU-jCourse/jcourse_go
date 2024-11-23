package task

import (
	"jcourse_go/task/asynq"
	"jcourse_go/task/base"
	"jcourse_go/task/biz/ping"
	"jcourse_go/task/biz/statistic"
	"log"
)

var (
	TaskRegistry = map[string]base.TaskHandler{
		statistic.TypeSaveDailyStatistic: statistic.HandleSaveStatisticTask,
		statistic.TypeRefreshPVDupJudge:  statistic.HandleRefreshPVDupJudgeTask,
		ping.TypePing:                    ping.TaskPingHandler,
	}
)

var TaskManager IAsyncTaskManager

// panic if failed
func MustInitTaskManager(redisConfig base.RedisConfig) {
	TaskManager = asynq.NewAsynqTaskManager(redisConfig)
	// 1. register task handler
	for taskType, handler := range TaskRegistry {
		TaskManager.RegisterTaskHandler(taskType, handler)
	}

	// 2. run server
	if err := TaskManager.RunServer(); err != nil {
		panic(err)
	}

	// 3. submit periodic task
	TaskManager.Submit(statistic.IntervalSaveDailyStatistic, TaskManager.CreateTask(statistic.TypeSaveDailyStatistic, nil))
	TaskManager.Submit(statistic.IntervalRefreshPVDupJudge, TaskManager.CreateTask(statistic.TypeRefreshPVDupJudge, nil))

	log.Println("[MustInitTaskManager] success")
}
