package olddto

import (
	"jcourse_go/internal/domain/statistic"
)

type StatisticRequest struct {
	StartTime  int64                     `json:"start_time" form:"start_time"` // unix timestamp, 单位秒
	EndTime    int64                     `json:"end_time" form:"end_time"`     // unix timestamp, 单位秒
	PeriodKeys []statistic.PeriodInfoKey `json:"period_keys" form:"period_keys"`
}

type StatisticResponse struct {
	DailyInfos  []statistic.DailyStatistic `json:"daily_infos"`
	PeriodInfos []statistic.PeriodInfo     `json:"period_infos"`
}
