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

// GetPeriodInfoFromPOs 从统计数据中获取指定的周期信息, 调用者保证pos中的数据是按时间增序排列的, 保证返回的数据是按时间增序排列的
func GetPeriodInfoFromPOs(pos []po.StatisticPO, keys []model.PeriodInfoKey) (map[model.PeriodInfoKey][]model.PeriodInfo, error) {
	const week = 7
	const month = 30
	periodInfoMap := make(map[model.PeriodInfoKey][]model.PeriodInfo)
	total := len(pos)
	for _, key := range keys {
		periodInfoMap[key] = make([]model.PeriodInfo, 0)
		switch key {
		case model.PeriodInfoKeyMAU:
			months := total / month
			// 这里反向遍历, 保证返回的数据是按时间增序排列的
			for i := months - 1; i >= 0; i-- {
				end := total - 1 - i*month
				start := end - month + 1
				monthWindow := make([]int64, month)
				for j := start; j <= end; j++ {
					monthWindow[j-start] = pos[j].UVCount
				}
				newInfo := model.PeriodInfo{
					StartTime: pos[start].Date.Unix(),
					EndTime:   pos[end].Date.Unix(),
					Value:     mathutil.Average(monthWindow...),
					Key:       key,
				}
				periodInfoMap[key] = append(periodInfoMap[key], newInfo)
			}
		case model.PeriodInfoKeyWAU:
			weeks := total / week
			for i := weeks - 1; i >= 0; i-- {
				end := total - 1 - i*week
				start := end - week + 1
				weekWindow := make([]int64, week)
				for j := start; j <= end; j++ {
					weekWindow[j-start] = pos[j].UVCount
				}
				newInfo := model.PeriodInfo{
					StartTime: pos[start].Date.Unix(),
					EndTime:   pos[end].Date.Unix(),
					Value:     mathutil.Average(weekWindow...),
					Key:       key,
				}
				periodInfoMap[key] = append(periodInfoMap[key], newInfo)
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
