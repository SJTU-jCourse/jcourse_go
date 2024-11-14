package converter

import (
	"jcourse_go/model/model"
	"jcourse_go/model/po"
	"jcourse_go/util"

	"github.com/RoaringBitmap/roaring"
)

func ConvertDailyInfoMinimalFromPO(po po.StatisticPO) model.DailyInfoMinimal {
	return model.DailyInfoMinimal{
		ID:               int64(po.ID),
		Date:             util.FormatDate(po.Date),
		UVCount:          po.UVCount,
		PVCount:          po.PVCount,
		NewCourseCount:   po.NewCourses,
		NewUserCount:     po.NewUsers,
		NewReviewCount:   po.NewReviews,
		TotalCourseCount: po.TotalCourses,
		TotalReviewCount: po.TotalReviews,
		TotalUserCount:   po.TotalUsers,
	}
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
