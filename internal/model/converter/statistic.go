package converter

import (
	"fmt"

	"github.com/duke-git/lancet/v2/mathutil"

	"github.com/RoaringBitmap/roaring"

	"jcourse_go/internal/domain/statistic"
	entity2 "jcourse_go/internal/infrastructure/entity"
)

func ConvertDailyInfoFromPO(po *entity2.Statistic) statistic.DailyStatistic {
	return statistic.DailyStatistic{
		ID:               po.ID,
		Date:             po.Date,
		UVCount:          po.UVCount,
		PVCount:          po.PVCount,
		NewUserCount:     po.DailyNewUser,
		NewReviewCount:   po.DailyNewReview,
		TotalReviewCount: po.TotalReview,
		TotalUserCount:   po.TotalUser,
	}
}

// GetPeriodInfoFromPOs 从统计数据中获取指定的周期信息, 调用者保证pos中的数据是按时间增序排列的, 保证返回的数据是按时间增序排列的
func GetPeriodInfoFromPOs(pos []*entity2.Statistic, keys []statistic.PeriodInfoKey) (map[statistic.PeriodInfoKey][]statistic.PeriodInfo, error) {
	const week = 7
	const month = 30
	periodInfoMap := make(map[statistic.PeriodInfoKey][]statistic.PeriodInfo)
	total := len(pos)
	for _, key := range keys {
		periodInfoMap[key] = make([]statistic.PeriodInfo, 0)
		switch key {
		case statistic.PeriodInfoKeyMAU:
			months := total / month
			// 这里反向遍历, 保证返回的数据是按时间增序排列的
			for i := months - 1; i >= 0; i-- {
				end := total - 1 - i*month
				start := end - month + 1
				monthWindow := make([]int64, month)
				for j := start; j <= end; j++ {
					monthWindow[j-start] = pos[j].UVCount
				}
				newInfo := statistic.PeriodInfo{
					StartDate: pos[start].Date,
					EndDate:   pos[end].Date,
					Value:     mathutil.Average(monthWindow...),
					Key:       key,
				}
				periodInfoMap[key] = append(periodInfoMap[key], newInfo)
			}
		case statistic.PeriodInfoKeyWAU:
			weeks := total / week
			for i := weeks - 1; i >= 0; i-- {
				end := total - 1 - i*week
				start := end - week + 1
				weekWindow := make([]int64, week)
				for j := start; j <= end; j++ {
					weekWindow[j-start] = pos[j].UVCount
				}
				newInfo := statistic.PeriodInfo{
					StartDate: pos[start].Date,
					EndDate:   pos[end].Date,
					Value:     mathutil.Average(weekWindow...),
					Key:       key,
				}
				periodInfoMap[key] = append(periodInfoMap[key], newInfo)
			}
		default:
			return nil, fmt.Errorf(statistic.ErrInvalidPeriodInfoKey, key)
		}
	}
	return periodInfoMap, nil
}

func ConvertUVDataFromPO(data []byte) (statistic.UVData, error) {
	uv := roaring.New()
	err := uv.UnmarshalBinary(data)
	if err != nil {
		return nil, err
	}
	return uv, nil
}

func ConvertStatisticDataFromPO(po *entity.StatisticDataPO) (statistic.StatisticData, error) {
	uv, err := ConvertUVDataFromPO(po.UVData)
	if err != nil {
		return statistic.StatisticData{}, err
	}
	return statistic.StatisticData{
		ID:          int64(po.ID),
		StatisticID: po.StatisticID,
		Date:        po.Date,
		UVData:      uv,
	}, nil
}
