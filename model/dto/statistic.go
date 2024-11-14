package dto

import "jcourse_go/model/model"

type StatisticRequest struct {
	DateTime int64 `json:"date_time" form:"date_time"` // unix timestamp, 单位秒, 如果不指定则默认为当天
	InDetail bool  `json:"in_detail" form:"in_detail"`
}

type StatisticResponse struct {
	DailyInfos []model.DailyInfoDetail `json:"daily_infos"`
}
