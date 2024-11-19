package statistic

import (
	"context"
	"jcourse_go/dal"
	"jcourse_go/middleware"
	"jcourse_go/model/converter"
	"jcourse_go/repository"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSaveStatistic(t *testing.T) {
	t.Run("TestSaveOneWithData", func(t *testing.T) {
		ctx := context.Background()
		dal.InitTestMemDBClient()
		db := dal.GetDBClient()
		err := repository.Migrate(db)
		if err != nil {
			t.Errorf("Migrate() error = %v", err)
		}
		// uv, pv data
		userNum := 10000
		uvm := middleware.NewUVMiddleware()
		pvm := middleware.NewPVMiddleware()
		uris := []string{
			"/",
			"/course",
			"/review",
			"/test",
		}
		for i := 0; i < userNum; i++ {
			uvm.AddUser(int64(i))
			pvm.AddPageView(int64(i), uris[i%len(uris)])
		}
		// count data
		userQuery := repository.NewUserQuery(dal.GetDBClient())
		_, err = userQuery.CreateUser(ctx, "test@example.com", "password")
		if err != nil {
			return
		}
		err = SaveStatistic(ctx, db, pvm, uvm, time.Now())
		if err != nil {
			t.Errorf("SaveStatistic() error = %v", err)
		}
		statisticQuery := repository.NewStatisticQuery(dal.GetDBClient())
		statisticDataQuery := repository.NewStatisticDataQuery(dal.GetDBClient())
		statisticPOs, err := statisticQuery.GetStatistics(ctx, repository.WithDate(time.Now()))
		if err != nil {
			t.Errorf("SaveStatistic() error = %v", err)
		}
		if len(statisticPOs) != 1 {
			t.Errorf("SaveStatistic() error = %v", err)
		}
		item := statisticPOs[0]
		assert.Equal(t, int64(userNum), item.UVCount)
		assert.Equal(t, int64(userNum), item.PVCount)
		assert.Equal(t, int64(0), item.NewReviews)
		assert.Equal(t, item.NewUsers, int64(1))

		data, err := statisticDataQuery.GetUVDataList(ctx, repository.WithDate(time.Now()))
		if err != nil {
			t.Errorf("SaveStatistic() error = %v", err)
		}
		if len(data) != 1 {
			t.Errorf("SaveStatistic() error = %v", err)
		}
		bytes := data[0]
		uv, err := converter.ConvertUVDataFromPO(bytes)
		if err != nil {
			t.Errorf("SaveStatistic() error = %v", err)
		}
		assert.Equal(t, userNum, int(uv.GetCardinality()))
		for i := 0; i < userNum; i++ {
			assert.Equal(t, true, uv.Contains(uint32(i)))
		}

		if err != nil {
			t.Errorf("SaveStatistic() error = %v", err)
		}
	})
}

// TODO: Need Save mock: unique index constraint for date
//func BenchmarkSaveStatistic(b *testing.B) {
//	// 10000 user
//	userNum := 10000
//	uvm := middleware.NewUVMiddleware()
//	pvm := middleware.NewPVMiddleware()
//	uris := []string{
//		"/",
//		"/course",
//		"/review",
//		"/test",
//	}
//	for i := 0; i < userNum; i++ {
//		uvm.AddUser(int64(i))
//		pvm.AddPageView(int64(i), uris[i%len(uris)])
//	}
//	dal.InitTestMemDBClient()
//	db := dal.GetDBClient()
//	err := repository.Migrate(db)
//	if err != nil {
//		b.Errorf("Migrate() error = %v", err)
//	}
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//		err := SaveStatistic(context.Background(), db, time.Now())
//		if err != nil {
//			b.Errorf("SaveStatistic() error = %v", err)
//		}
//	}
//}
