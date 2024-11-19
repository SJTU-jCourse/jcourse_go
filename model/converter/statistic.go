package converter

import (
	"fmt"
	"jcourse_go/model/model"
	"jcourse_go/model/po"
	"jcourse_go/util"

	"github.com/duke-git/lancet/v2/mathutil"

	"github.com/RoaringBitmap/roaring"
)

func ConvertDailyInfoFromPO(po po.StatisticPO) model.DailyInfo {
	return model.DailyInfo{
		ID:               int64(po.ID),
		Date:             util.FormatDate(po.Date),
		UVCount:          po.UVCount,
		PVCount:          po.PVCount,
		NewUserCount:     po.NewUsers,
		NewReviewCount:   po.NewReviews,
		TotalReviewCount: po.TotalReviews,
		TotalUserCount:   po.TotalUsers,
	}
}

// GetPeriodInfoFromPOs 从统计数据中获取指定的周期信息, 调用者保证pos中的数据是按时间增序排列的
func GetPeriodInfoFromPOs(pos []po.StatisticPO, keys []model.PeriodInfoKey) (map[model.PeriodInfoKey][]model.PeriodInfo, error) {
	const week = 7
	const month = 30
	periodInfoMap := make(map[model.PeriodInfoKey][]model.PeriodInfo)
	for _, key := range keys {
		switch key {
		case model.PeriodInfoKeyMAU:
			for i := range pos {
				if i >= month-1 {
					monthWindow := make([]int64, month)
					for j := i - month + 1; j <= i; j++ {
						monthWindow[j] = pos[j].UVCount
					}
					newInfo := model.PeriodInfo{
						StartTime: pos[i-month+1].Date.Unix(),
						EndTime:   pos[i].Date.Unix(),
						Value:     mathutil.Average(monthWindow...),
						Key:       key,
					}
					periodInfoMap[key] = append(periodInfoMap[key], newInfo)
				}
			}
		case model.PeriodInfoKeyWAU:
			periodInfoMap[key] = make([]model.PeriodInfo, 0)
			for i := range pos {
				if i >= week-1 {
					weekWindow := make([]int64, week)
					for j := i - week + 1; j <= i; j++ {
						weekWindow[j] = pos[j].UVCount
					}
					newInfo := model.PeriodInfo{
						StartTime: pos[i-week+1].Date.Unix(),
						EndTime:   pos[i].Date.Unix(),
						Value:     mathutil.Average(weekWindow...),
						Key:       key,
					}
					periodInfoMap[key] = append(periodInfoMap[key], newInfo)
				}
			}
		default:
			return nil, fmt.Errorf(model.ErrInvalidPeriodInfoKey, key)
		}
	}
	return periodInfoMap, nil
}

func ConvertUVDataFromPO(data []byte) (model.UVData, error) {
	uv := roaring.New()
	err := uv.UnmarshalBinary(data)
	if err != nil {
		return nil, err
	}
	return uv, nil
}

func ConvertStatisticDataFromPO(po *po.StatisticDataPO) (model.StatisticData, error) {
	uv, err := ConvertUVDataFromPO(po.UVData)
	if err != nil {
		return model.StatisticData{}, err
	}
	return model.StatisticData{
		ID:          int64(po.ID),
		StatisticID: po.StatisticID,
		Date:        util.FormatDate(po.Date),
		UVData:      uv,
	}, nil
}
