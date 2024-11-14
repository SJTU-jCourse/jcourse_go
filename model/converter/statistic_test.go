package converter

import (
	"jcourse_go/model/po"
	"jcourse_go/util"
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
