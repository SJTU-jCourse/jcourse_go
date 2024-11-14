package service

import (
	"context"
	"errors"
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
func getDailyInfoInCache(ctx context.Context, datetime time.Time) (model.DailyInfoDetail, error) {
	// TODO
	return model.DailyInfoDetail{}, util.ErrNotFound
}
func getAllStatisticInCache(ctx context.Context) ([]model.DailyInfoDetail, error) {
	return []model.DailyInfoDetail{}, util.ErrNotFound
}
func updateDailyInfoCache(ctx context.Context, detail model.DailyInfoDetail) error {
	key := getDailyInfoKeyFromStr(detail.Date)
	log.Printf("Update: %s", key)
	// TODO update
	return nil
}

// CalDailyInfo 计算某一天0-24点的日活、新增课程数、新增点评数
func CalDailyInfo(ctx context.Context, datetime time.Time) (model.DailyInfoMinimal, error) {
	db := dal.GetDBClient()
	userQuery := repository.NewUserQuery(db)
	courseQuery := repository.NewCourseQuery(db)
	reviewQuery := repository.NewReviewQuery(db)
	dayStart, dayEnd := util.GetDayTimeRange(datetime)
	dailyOpt := repository.WithCreatedAtBetween(dayStart, dayEnd)

	dailyInfo := model.DailyInfoMinimal{}
	newUserCount, err := userQuery.GetUserCount(ctx, dailyOpt)
	dailyInfo.NewUserCount = newUserCount
	if err != nil {
		return model.DailyInfoMinimal{}, err
	}
	newCourseCount, err := courseQuery.GetCourseCount(ctx, dailyOpt)
	dailyInfo.NewCourseCount = newCourseCount
	if err != nil {
		return model.DailyInfoMinimal{}, err
	}
	newReviewCount, err := reviewQuery.GetReviewCount(ctx, dailyOpt)
	dailyInfo.NewReviewCount = newReviewCount
	if err != nil {
		return model.DailyInfoMinimal{}, err
	}
	return dailyInfo, nil
}
func getDailyInfoMinimalInDB(ctx context.Context, datetime time.Time) (model.DailyInfoMinimal, error) {
	db := dal.GetDBClient()
	statisticQuery := repository.NewStatisticQuery(db)
	option := repository.WithDate(datetime)
	statistics, err := statisticQuery.GetStatistics(ctx, option)
	if err != nil {
		return model.DailyInfoMinimal{}, err
	}
	if len(statistics) == 0 {
		return model.DailyInfoMinimal{}, util.ErrNotFound
	}
	dailyInfo := converter.ConvertDailyInfoMinimalFromPO(statistics[0])
	return dailyInfo, nil
}

func getDailyInfoInDB(ctx context.Context, datetime time.Time) (model.DailyInfoDetail, error) {
	db := dal.GetDBClient()
	statistcQuery := repository.NewStatisticQuery(db)

	_, dateEnd := util.GetDayTimeRange(datetime)
	monthStart := util.GetTimeSubMonth(dateEnd, 1)
	options := []repository.DBOption{
		repository.WithDateBetween(monthStart, dateEnd),
		repository.WithDateOrder(false),
	}

	statistics, err := statistcQuery.GetStatistics(ctx, options...)
	if err != nil {
		return model.DailyInfoDetail{}, err
	}
	if len(statistics) <= 7 {
		return model.DailyInfoDetail{}, util.ErrNotFound
	}
	var wau int64 = 0
	var mau int64 = 0
	total := len(statistics)

	for i, statistic := range statistics {
		if total-i <= 7 {
			wau += statistic.UVCount
		}
		mau += statistic.UVCount
	}
	dayStatistic := statistics[total-1]
	dailyInfo := model.DailyInfoDetail{
		DailyInfoMinimal: converter.ConvertDailyInfoMinimalFromPO(dayStatistic),
		WAU:              wau,
		MAU:              mau,
	}
	return dailyInfo, nil

}

func GetDailyStatisticDetail(ctx context.Context, datetime time.Time) (model.DailyInfoDetail, error) {
	dailyInfo, err := getDailyInfoInCache(ctx, datetime)
	if err == nil {
		return dailyInfo, nil
	}
	if !errors.Is(err, util.ErrNotFound) {
		return model.DailyInfoDetail{}, err
	}
	detail, err := getDailyInfoInDB(ctx, datetime)
	if err != nil {
		return model.DailyInfoDetail{}, err
	}
	err = updateDailyInfoCache(ctx, detail)
	if err != nil {
		return model.DailyInfoDetail{}, err
	}
	return detail, err
}
func GetDailyStatisticMinimal(ctx context.Context, datetime time.Time) (model.DailyInfoMinimal, error) {
	dailyInfo, err := getDailyInfoInCache(ctx, datetime)
	if err == nil {
		return dailyInfo.DailyInfoMinimal, nil
	}
	if !errors.Is(err, util.ErrNotFound) {
		return model.DailyInfoMinimal{}, err
	}
	return getDailyInfoMinimalInDB(ctx, datetime)
}
func avgInt64(items []int64) int64 {
	var avgCount int64 = 0
	var total int64 = 0
	for _, item := range items {
		avgCount += item
		total += 1
	}
	if total == 0 {
		return 0
	}
	return avgCount / total
}

func GetAllStatistics(ctx context.Context) ([]model.DailyInfoDetail, error) {
	statisticQuery := repository.NewStatisticQuery(dal.GetDBClient())
	statistics, err := statisticQuery.GetStatistics(ctx)
	if err != nil {
		return []model.DailyInfoDetail{}, err
	}
	num := len(statistics)
	details := make([]model.DailyInfoDetail, num)

	for i, statisticPO := range statistics {
		weekWindow := make([]int64, 7)
		monthWindow := make([]int64, 30)
		minimal := converter.ConvertDailyInfoMinimalFromPO(statisticPO)
		details[i] = model.DailyInfoDetail{
			DailyInfoMinimal: minimal,
			WAU:              0,
			MAU:              0,
		}
		if i >= 6 {
			for j := i - 6; j <= i; j++ {
				weekWindow[i] = details[j].DailyInfoMinimal.UVCount
			}
			details[i].WAU = avgInt64(weekWindow)
		}
		if i >= 29 {
			for j := i - 29; j <= i; j++ {
				monthWindow[i] = details[j].DailyInfoMinimal.UVCount
			}
			details[i].MAU = avgInt64(monthWindow)
		}
	}
	return details, nil
}

func GetDailyUVData(ctx context.Context, datetime time.Time) (model.UVData, error) {
	db := dal.GetDBClient()
	statisticDataQuery := repository.NewStatisticDataQuery(db)
	data, err := statisticDataQuery.GetUVDataList(ctx, repository.WithDate(datetime))
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
