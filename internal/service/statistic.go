package service

import (
	"context"
	"log"
	"time"

	"github.com/RoaringBitmap/roaring"
	"github.com/pkg/errors"

	"jcourse_go/internal/domain/statistic"
	entity2 "jcourse_go/internal/infrastructure/entity"
	"jcourse_go/internal/infrastructure/repository"
	"jcourse_go/internal/model/converter"
	util2 "jcourse_go/pkg/util"
)

// TODO: cache
const (
	GlobalInfoKey      = "global_info"
	DailyInfoKeyPrefix = "daily_info"
)

func getDailyInfoKeyFromTime(datetime time.Time) string {
	return DailyInfoKeyPrefix + ":" + util2.FormatDate(datetime)
}
func getDailyInfoKeyFromStr(datetime string) string { return DailyInfoKeyPrefix + ":" + datetime }
func getDailyInfoInCache(ctx context.Context, datetime time.Time) (statistic.DailyInfo, error) {
	// TODO
	return statistic.DailyInfo{}, util2.ErrNotFound
}
func updateDailyInfoCache(ctx context.Context, detail statistic.DailyInfo) error {
	key := getDailyInfoKeyFromStr(detail.Date)
	log.Printf("Update: %s", key)
	// TODO
	return nil
}

// CalDailyInfo 计算某一天0-24点的日活、新增课程数、新增点评数
func CalDailyInfo(ctx context.Context, datetime time.Time) (statistic.DailyInfo, error) {
	u := repository.Q.UserPO
	r := repository.Q.ReviewPO

	dayStart, dayEnd := util2.GetDayTimeRange(datetime)

	dailyInfo := statistic.DailyInfo{}
	newUserCount, err := u.WithContext(ctx).Where(u.CreatedAt.Between(dayStart, dayEnd)).Count()
	dailyInfo.NewUserCount = newUserCount
	if err != nil {
		return statistic.DailyInfo{}, err
	}
	newReviewCount, err := r.WithContext(ctx).Where(r.CreatedAt.Between(dayStart, dayEnd)).Count()
	dailyInfo.NewReviewCount = newReviewCount
	if err != nil {
		return statistic.DailyInfo{}, err
	}
	return dailyInfo, nil
}

// buildStatisticDBOptionsFromFilter filter保证传入str要么是空字符串, 要么是合法的日期字符串
func buildStatisticDBOptionsFromFilter(ctx context.Context, q *repository.Query, filter statistic.StatisticFilter) repository.IStatisticPODo {
	builder := q.StatisticPO.WithContext(ctx)
	s := q.StatisticPO
	if filter.StartDate != "" && filter.EndDate != "" {
		builder = builder.Where(s.Date.Between(filter.StartDate, filter.EndDate))
	} else if filter.StartDate != "" {
		builder = builder.Where(s.Date.Gte(filter.StartDate))
	} else if filter.EndDate != "" {
		builder = builder.Where(s.Date.Lte(filter.EndDate))
	}
	return builder
}
func GetStatistics(ctx context.Context, filter statistic.StatisticFilter) ([]statistic.DailyInfo, []statistic.PeriodInfo, error) {
	// TODO: cache

	q := buildStatisticDBOptionsFromFilter(ctx, repository.Q, filter)
	statistics, err := q.Find()
	if err != nil {
		return []statistic.DailyInfo{}, []statistic.PeriodInfo{}, err
	}
	num := len(statistics)
	dailyInfos := make([]statistic.DailyInfo, num)
	if len(filter.PeriodInfoKeys) == 0 {
		for i, statisticPO := range statistics {
			dailyInfos[i] = converter.ConvertDailyInfoFromPO(statisticPO)
		}
		return dailyInfos, nil, nil
	}
	periodInfos := make([]statistic.PeriodInfo, 0)
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

func GetDailyUVData(ctx context.Context, datetime time.Time) (statistic.UVData, error) {
	s := repository.Q.StatisticDataPO
	date := util2.FormatDate(datetime)
	data, err := s.WithContext(ctx).Where(s.Date.Eq(date)).Find()
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, util2.ErrNotFound
	}
	uvData, err := converter.ConvertUVDataFromPO(data[0].UVData)
	if err != nil {
		return nil, err
	}
	return uvData, nil
}

func CreateStatistic(ctx context.Context, datetime time.Time, uvCount, pvCount int64, uvData *roaring.Bitmap) error {
	u := repository.Q.UserPO
	r := repository.Q.ReviewPO

	// get counts
	dayStart, dayEnd := util2.GetDayTimeRange(datetime)
	newUserCount, err := u.WithContext(ctx).Where(u.CreatedAt.Between(dayStart, dayEnd)).Count()
	if err != nil {
		log.Fatalf("failed to save statistic: %v", err)
		return err
	}
	newReviewCount, err := r.WithContext(ctx).Where(r.CreatedAt.Between(dayStart, dayEnd)).Count()
	if err != nil {
		log.Fatalf("failed to save statistic: %v", err)
		return err
	}
	if err != nil {
		log.Fatalf("failed to save statistic: %v", err)
		return err
	}
	totalReviewCount, err := r.WithContext(ctx).Count()
	if err != nil {
		log.Fatalf("failed to save statistic: %v", err)
		return err
	}
	totalUserCount, err := u.WithContext(ctx).Count()
	if err != nil {
		log.Fatalf("failed to save statistic: %v", err)
		return err
	}
	// create statistic item
	newStatisticItem := entity2.StatisticPO{
		Date:         util2.FormatDate(datetime),
		UVCount:      uvCount,
		PVCount:      pvCount,
		NewReviews:   newReviewCount,
		NewUsers:     newUserCount,
		TotalReviews: totalReviewCount,
		TotalUsers:   totalUserCount,
	}

	err = repository.Q.StatisticPO.WithContext(ctx).Create(&newStatisticItem)
	if err != nil {
		log.Fatalf("failed to save statistic: %v", err)
		return err
	}
	// save uv data
	uvDataBytes, err := uvData.ToBytes()
	if err != nil {
		log.Fatalf("failed to save statistic: %v", err)
		return err
	}
	if newStatisticItem.ID <= 0 {
		return errors.Errorf("failed to save statistic: Not write ID back")
	}
	err = repository.Q.StatisticDataPO.WithContext(ctx).Create(&entity.StatisticDataPO{
		Date:        util2.FormatDate(datetime),
		UVData:      uvDataBytes,
		StatisticID: newStatisticItem.ID,
	})
	if err != nil {
		return err
	}
	return nil
}
