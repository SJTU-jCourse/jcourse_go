package dto

import (
	"jcourse_go/internal/model/model"
)

type StatisticRequest struct {
	StartTime  int64                 `json:"start_time" form:"start_time"` // unix timestamp, 单位秒
	EndTime    int64                 `json:"end_time" form:"end_time"`     // unix timestamp, 单位秒
	PeriodKeys []model.PeriodInfoKey `json:"period_keys" form:"period_keys"`
}

type StatisticResponse struct {
	DailyInfos  []model.DailyInfo  `json:"daily_infos"`
	PeriodInfos []model.PeriodInfo `json:"period_infos"`
}
