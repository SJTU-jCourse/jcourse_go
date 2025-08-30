package statistic

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"jcourse_go/internal/dal"
	"jcourse_go/internal/middleware"
	"jcourse_go/internal/model/converter"
	"jcourse_go/internal/model/po"
	repository2 "jcourse_go/internal/repository"
	"jcourse_go/pkg/util"
)

func TestSaveStatistic(t *testing.T) {
	t.Run("TestSaveOneWithData", func(t *testing.T) {
		ctx := context.Background()
		dal.InitTestMemDBClient()
		db := dal.GetDBClient()
		err := dal.Migrate(db)
		if err != nil {
			t.Errorf("Migrate() error = %v", err)
		}
		repository2.SetDefault(db)
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
		u := repository2.Q.UserPO
		userPO := po.UserPO{
			Email:    "test@example.com",
			Password: "password",
		}
		err = u.WithContext(ctx).Create(&userPO)
		if err != nil {
			return
		}
		err = SaveStatistic(ctx, pvm, uvm, time.Now())
		if err != nil {
			t.Errorf("SaveStatistic() error = %v", err)
		}
		s := repository2.Q.StatisticPO
		d := repository2.Q.StatisticDataPO

		item, err := s.WithContext(ctx).Where(s.Date.Eq(util.FormatDate(time.Now()))).Take()
		if err != nil {
			t.Errorf("SaveStatistic() error = %v", err)
		}

		assert.Equal(t, int64(userNum), item.UVCount)
		assert.Equal(t, int64(userNum), item.PVCount)
		assert.Equal(t, int64(0), item.NewReviews)
		assert.Equal(t, item.NewUsers, int64(1))

		bytes, err := d.WithContext(ctx).Where(d.Date.Eq(util.FormatDate(time.Now()))).Take()
		if err != nil {
			t.Errorf("SaveStatistic() error = %v", err)
		}
		uv, err := converter.ConvertUVDataFromPO(bytes.UVData)
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
// func BenchmarkSaveStatistic(b *testing.B) {
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
// }
