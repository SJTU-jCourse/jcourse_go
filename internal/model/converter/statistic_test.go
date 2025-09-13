package converter

import (
	"math/rand"
	"testing"
	"time"

	"jcourse_go/internal/entity"

	"jcourse_go/internal/domain/statistic"
	entity2 "jcourse_go/internal/infrastructure/entity"
	"jcourse_go/pkg/util"

	"github.com/RoaringBitmap/roaring"
	"github.com/stretchr/testify/assert"
)

func TestConvertStatisticDataFromPO(t *testing.T) {
	t.Run("TestConvertStatisticDataFromPO", func(t *testing.T) {
		bitmap := roaring.New()
		bitmap.Add(1)
		bitmap.Add(2)
		bitmap.Add(10086)
		bytes, err := bitmap.ToBytes()
		if err != nil {
			t.Errorf("ConvertStatisticDataFromPO error: %v", err)
		}
		data := entity.StatisticDataPO{
			StatisticID: 1,
			Date:        util.FormatDate(time.Now()),
			UVData:      bytes,
		}
		statisticData, err := ConvertStatisticDataFromPO(&data)
		if err != nil {
			t.Errorf("ConvertStatisticDataFromPO error: %v", err)
		}
		assert.Equal(t, int64(1), statisticData.StatisticID)
		assert.Equal(t, uint64(3), statisticData.UVData.GetCardinality())
		assert.Equal(t, util.FormatDate(time.Now()), statisticData.Date)
		assert.True(t, statisticData.UVData.Contains(1))
		assert.True(t, statisticData.UVData.Contains(2))
		assert.True(t, statisticData.UVData.Contains(10086))
	})
}

func TestGetPeriodInfoFromPOs(t *testing.T) {
	// Helper function to create a Statistic with a given date and UVCount
	createStatisticPO := func(daysAgo int) entity2.Statistic {
		return entity2.Statistic{
			Date:           util.FormatDate(util.GetMidTime(time.Now()).AddDate(0, 0, -daysAgo)),
			UVCount:        int64(rand.Uint64() % 1000),
			PVCount:        int64(rand.Uint64() % 1000),
			TotalReview:    10000 - int64(daysAgo)*100,
			TotalUser:      1000 - int64(daysAgo)*10,
			DailyNewReview: 100,
			DailyNewUser:   10,
		}
	}

	// Test data
	pos := make([]*entity2.Statistic, 40)
	for i := 0; i < 40; i++ {
		p := createStatisticPO(i)
		pos[i] = &p
	}

	keys := []statistic.PeriodInfoKey{statistic.PeriodInfoKeyMAU, statistic.PeriodInfoKeyWAU}

	// Call the function
	periodInfoMap, err := GetPeriodInfoFromPOs(pos, keys)
	assert.Nil(t, err)

	// Verify the results for MAU
	mauInfo := periodInfoMap[statistic.PeriodInfoKeyMAU]
	assert.Equal(t, 1, len(mauInfo))

	assert.Equal(t, pos[10].Date, mauInfo[0].StartDate, 0)
	assert.Equal(t, pos[39].Date, mauInfo[0].EndDate)
	mau := int64(0)
	for i := 10; i < 40; i++ {
		mau += pos[i].UVCount
	}
	mau /= 30
	assert.Equal(t, mau, mauInfo[0].Value) // Average of UVCounts from 0 to 29

	// Verify the results for WAU
	wauInfo := periodInfoMap[statistic.PeriodInfoKeyWAU]
	assert.Equal(t, 5, len(wauInfo))
	for i, info := range wauInfo {
		start := 7*i + 5 // 6, 13, 20, 27, 34
		end := start + 6
		expectedValue := int64(0)
		for j := start; j <= end; j++ {
			expectedValue += pos[j].UVCount
		}
		expectedValue /= 7
		assert.Equal(t, pos[start].Date, info.StartDate, 0)
		assert.Equal(t, pos[end].Date, info.EndDate)
		assert.Equal(t, expectedValue, info.Value)
	}
}
