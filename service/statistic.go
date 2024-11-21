package service

import (
	"context"
	"jcourse_go/dal"
	"jcourse_go/model/converter"
	"jcourse_go/model/model"
	"jcourse_go/repository"
	"jcourse_go/util"
	"log"
	"time"
)

// TODO: cache
const (
	GlobalInfoKey      = "global_info"
	DailyInfoKeyPrefix = "daily_info"
)

func getDailyInfoKeyFromTime(datetime time.Time) string {
	return DailyInfoKeyPrefix + ":" + util.FormatDate(datetime)
}
func getDailyInfoKeyFromStr(datetime string) string { return DailyInfoKeyPrefix + ":" + datetime }
func getDailyInfoInCache(ctx context.Context, datetime time.Time) (model.DailyInfo, error) {
	// TODO
	return model.DailyInfo{}, util.ErrNotFound
}
func updateDailyInfoCache(ctx context.Context, detail model.DailyInfo) error {
	key := getDailyInfoKeyFromStr(detail.Date)
	log.Printf("Update: %s", key)
	// TODO
	return nil
}

// CalDailyInfo 计算某一天0-24点的日活、新增课程数、新增点评数
func CalDailyInfo(ctx context.Context, datetime time.Time) (model.DailyInfo, error) {
	db := dal.GetDBClient()
	userQuery := repository.NewUserQuery(db)
	reviewQuery := repository.NewReviewQuery(db)
	dayStart, dayEnd := util.GetDayTimeRange(datetime)
	dailyOpt := repository.WithCreatedAtBetween(dayStart, dayEnd)

	dailyInfo := model.DailyInfo{}
	newUserCount, err := userQuery.GetUserCount(ctx, dailyOpt)
	dailyInfo.NewUserCount = newUserCount
	if err != nil {
		return model.DailyInfo{}, err
	}
	newReviewCount, err := reviewQuery.GetReviewCount(ctx, dailyOpt)
	dailyInfo.NewReviewCount = newReviewCount
	if err != nil {
		return model.DailyInfo{}, err
	}
	return dailyInfo, nil
}

// buildStatisticDBOptionsFromFilter filter保证传入str要么是空字符串, 要么是合法的日期字符串
func buildStatisticDBOptionsFromFilter(query repository.IStatisticQuery, filter model.StatisticFilter) []repository.DBOption {
	opts := make([]repository.DBOption, 0)
	if filter.StartDate != "" && filter.EndDate != "" {
		opts = append(opts, repository.WithDateBetween(filter.StartDate, filter.EndDate))
	} else if filter.StartDate != "" {
		opts = append(opts, repository.WithDateAfter(filter.StartDate))
	} else if filter.EndDate != "" {
		opts = append(opts, repository.WithDateBefore(filter.EndDate))
	}
	return opts
}
func GetStatistics(ctx context.Context, filter model.StatisticFilter) ([]model.DailyInfo, []model.PeriodInfo, error) {
	// TODO: cache
	statisticQuery := repository.NewStatisticQuery(dal.GetDBClient())
	options := buildStatisticDBOptionsFromFilter(statisticQuery, filter)
	statistics, err := statisticQuery.GetStatistics(ctx, options...)
	if err != nil {
		return []model.DailyInfo{}, []model.PeriodInfo{}, err
	}
	num := len(statistics)
	dailyInfos := make([]model.DailyInfo, num)
	if len(filter.PeriodInfoKeys) == 0 {
		for i, statisticPO := range statistics {
			dailyInfos[i] = converter.ConvertDailyInfoFromPO(statisticPO)
		}
		return dailyInfos, nil, nil
	}
	periodInfos := make([]model.PeriodInfo, 0)
	infoMap, err := converter.GetPeriodInfoFromPOs(statistics, filter.PeriodInfoKeys)
	if err != nil {
		return nil, nil, err
	}
	for _, infos := range infoMap {
		periodInfos = append(periodInfos, infos...)
	}
	for i, statisticPO := range statistics {
		dailyInfos[i] = converter.ConvertDailyInfoFromPO(statisticPO)
	}
	return dailyInfos, periodInfos, nil
}

func GetDailyUVData(ctx context.Context, datetime time.Time) (model.UVData, error) {
	db := dal.GetDBClient()
	statisticDataQuery := repository.NewStatisticDataQuery(db)
	date := util.FormatDate(datetime)
	data, err := statisticDataQuery.GetUVDataList(ctx, repository.WithDate(date))
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, util.ErrNotFound
	}
	uvData, err := converter.ConvertUVDataFromPO(data[0])
	if err != nil {
		return nil, err
	}
	return uvData, nil
}
