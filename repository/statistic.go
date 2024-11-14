package repository

import (
	"context"
	"jcourse_go/model/po"
	"jcourse_go/util"

	"gorm.io/gorm"
)

type IStatisticQuery interface {
	GetStatistics(ctx context.Context, opts ...DBOption) ([]po.StatisticPO, error)
	GetStatisticsByIDs(ctx context.Context, statisticIDs []int64) (map[int64]po.StatisticPO, error)
	GetStatisticsCount(ctx context.Context, opts ...DBOption) (int64, error)
	CreateStatistic(ctx context.Context, statistic *po.StatisticPO) error
	UpdateStatistic(ctx context.Context, statistic *po.StatisticPO) error
}

type StatisticQuery struct {
	db *gorm.DB
}

func NewStatisticQuery(db *gorm.DB) IStatisticQuery {
	return &StatisticQuery{db: db}
}
func (q *StatisticQuery) optionDB(ctx context.Context, opts ...DBOption) *gorm.DB {
	db := q.db.WithContext(ctx).Model(&po.StatisticPO{})
	for _, opt := range opts {
		db = opt(db)
	}
	return db
}

func (q *StatisticQuery) GetStatistics(ctx context.Context, opts ...DBOption) ([]po.StatisticPO, error) {
	db := q.optionDB(ctx, opts...)
	statistics := make([]po.StatisticPO, 0)
	result := db.Find(&statistics)
	return statistics, result.Error
}

func (q *StatisticQuery) GetStatisticsCount(ctx context.Context, opts ...DBOption) (int64, error) {
	db := q.optionDB(ctx, opts...)
	var count int64
	result := db.Count(&count)
	return count, result.Error
}

func (q *StatisticQuery) GetStatisticsByIDs(ctx context.Context, statisticIDs []int64) (map[int64]po.StatisticPO, error) {
	db := q.optionDB(ctx, WithIDs(statisticIDs))
	statistics := make([]po.StatisticPO, 0)
	statisticMap := make(map[int64]po.StatisticPO)
	result := db.Find(&statistics)
	if result.Error != nil {
		return statisticMap, result.Error
	}
	for _, statistic := range statistics {
		statisticMap[int64(statistic.ID)] = statistic
	}
	return statisticMap, nil
}

func (q *StatisticQuery) CreateStatistic(ctx context.Context, statistic *po.StatisticPO) error {
	statistic.Date = util.GetMidTime(statistic.Date) // keep midday
	return q.db.WithContext(ctx).Create(statistic).Error
}

func (q *StatisticQuery) UpdateStatistic(ctx context.Context, statistic *po.StatisticPO) error {
	return q.db.WithContext(ctx).Save(statistic).Error
}

type IStatisticDataQuery interface {
	GetDataList(ctx context.Context, opts ...DBOption) ([]po.StatisticDataPO, error)
	GetUVDataList(ctx context.Context, opts ...DBOption) ([][]byte, error)
	GetDataByIDs(ctx context.Context, ids []int64) (map[int64]po.StatisticDataPO, error)
	CreateData(ctx context.Context, data *po.StatisticDataPO) error
}

type StatisticDataQuery struct {
	db *gorm.DB
}

func (q *StatisticDataQuery) optionDB(ctx context.Context, opts ...DBOption) *gorm.DB {
	db := q.db.WithContext(ctx).Model(&po.StatisticDataPO{})
	for _, opt := range opts {
		db = opt(db)
	}
	return db
}
func (q *StatisticDataQuery) GetDataList(ctx context.Context, opts ...DBOption) ([]po.StatisticDataPO, error) {
	db := q.optionDB(ctx, opts...)
	data := make([]po.StatisticDataPO, 0)
	result := db.Find(&data)
	return data, result.Error
}

func (q *StatisticDataQuery) GetUVDataList(ctx context.Context, opts ...DBOption) ([][]byte, error) {
	db := q.optionDB(ctx, opts...)
	data := make([][]byte, 0)
	result := db.Select("uv_data").Find(&data)
	return data, result.Error
}

func (q *StatisticDataQuery) GetDataByIDs(ctx context.Context, ids []int64) (map[int64]po.StatisticDataPO, error) {
	db := q.optionDB(ctx, WithIDs(ids))
	data := make([]po.StatisticDataPO, 0)
	dataMap := make(map[int64]po.StatisticDataPO)
	result := db.Find(&data)
	if result.Error != nil {
		return dataMap, result.Error
	}
	for _, d := range data {
		dataMap[int64(d.ID)] = d
	}
	return dataMap, nil
}

func (q *StatisticDataQuery) CreateData(ctx context.Context, data *po.StatisticDataPO) error {
	data.Date = util.GetMidTime(data.Date) // keep midday
	return q.db.WithContext(ctx).Create(data).Error
}

func NewStatisticDataQuery(db *gorm.DB) IStatisticDataQuery {
	return &StatisticDataQuery{
		db: db,
	}
}
