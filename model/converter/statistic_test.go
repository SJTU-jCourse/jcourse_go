package converter

import (
	"jcourse_go/model/model"
	"jcourse_go/model/po"
	"jcourse_go/util"
	"math/rand"
	"testing"
	"time"

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
		data := po.StatisticDataPO{
			StatisticID: 1,
			Date:        time.Now(),
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
	// Helper function to create a StatisticPO with a given date and UVCount
	createStatisticPO := func(daysAgo int) po.StatisticPO {
		return po.StatisticPO{
			Date:         util.GetMidTime(time.Now()).AddDate(0, 0, -daysAgo),
			UVCount:      int64(rand.Uint64() % 1000),
			PVCount:      int64(rand.Uint64() % 1000),
			TotalReviews: 10000 - int64(daysAgo)*100,
			TotalUsers:   1000 - int64(daysAgo)*10,
			NewReviews:   100,
			NewUsers:     10,
		}
	}

	// Test data
	pos := make([]po.StatisticPO, 40)
	for i := 0; i < 40; i++ {
		pos[i] = createStatisticPO(i)
	}

	keys := []model.PeriodInfoKey{model.PeriodInfoKeyMAU, model.PeriodInfoKeyWAU}

	// Call the function
	periodInfoMap, err := GetPeriodInfoFromPOs(pos, keys)
	assert.Nil(t, err)

	// Verify the results for MAU
	mauInfo := periodInfoMap[model.PeriodInfoKeyMAU]
	assert.Equal(t, 1, len(mauInfo))
	assert.Equal(t, pos[10].Date.Unix(), mauInfo[0].StartTime)
	assert.Equal(t, pos[39].Date.Unix(), mauInfo[0].EndTime)
	mau := int64(0)
	for i := 10; i < 40; i++ {
		mau += pos[i].UVCount
	}
	mau /= 30
	assert.Equal(t, mau, mauInfo[0].Value) // Average of UVCounts from 0 to 29

	// Verify the results for WAU
	wauInfo := periodInfoMap[model.PeriodInfoKeyWAU]
	assert.Equal(t, 5, len(wauInfo))
	for i, info := range wauInfo {
		start := 7*i + 5 // 6, 13, 20, 27, 34
		end := start + 6
		expectedValue := int64(0)
		for j := start; j <= end; j++ {
			expectedValue += pos[j].UVCount
		}
		expectedValue /= 7
		assert.Equal(t, pos[start].Date.Unix(), info.StartTime)
		assert.Equal(t, pos[end].Date.Unix(), info.EndTime)
		assert.Equal(t, expectedValue, info.Value)
	}
}
