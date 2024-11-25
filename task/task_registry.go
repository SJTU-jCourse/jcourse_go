package task

const (
	ScheduleQueues = "periodic"
)

// Deprecated

// var deamonTaskRegistry = []struct {
// 	cronspec string
// 	taskType string
// 	opts     []TaskOption
// }{
// 	{
// 		cronspec: statistic.IntervalSaveDailyStatistic,
// 		taskType: statistic.TypeSaveDailyStatistic,
// 		opts:     []TaskOption{{WithQueue: ScheduleQueues}},
// 	},
// 	{
// 		cronspec: statistic.IntervalRefreshPVDupJudge,
// 		taskType: statistic.TypeRefreshPVDupJudge,
// 		opts:     []TaskOption{{WithQueue: ScheduleQueues}},
// 	},
// }
