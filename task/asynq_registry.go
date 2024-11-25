package impl

import (
	"context"

	"github.com/hibiken/asynq"
)

var muxRegistry = map[string]func(context.Context, *asynq.Task) error{
	statistic.TypeSaveDailyStatistic: statistic.HandleSaveStatisticTask,
	statistic.TypeRefreshPVDupJudge:  statistic.HandleRefreshPVDupJudgeTask,
	ping.TypePing:                    ping.TaskPingHandler,
}

const (
	ScheduleQueues = "periodic"
)

var deamonTaskRegistry = []struct {
	cronspec string
	task     *asynq.Task
	opts     []asynq.Option
}{
	{
		cronspec: statistic.IntervalSaveDailyStatistic,
		task:     asynq.NewTask(statistic.TypeSaveDailyStatistic, nil),
		opts:     []asynq.Option{asynq.Queue(ScheduleQueues)},
	},
	{
		cronspec: statistic.IntervalRefreshPVDupJudge,
		task:     asynq.NewTask(statistic.TypeRefreshPVDupJudge, nil),
		opts:     []asynq.Option{asynq.Queue(ScheduleQueues)},
	},
}
