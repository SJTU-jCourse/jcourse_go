package statistic

import (
	"context"
	"jcourse_go/dal"
	"jcourse_go/middleware"
	"jcourse_go/model/po"
	"jcourse_go/repository"
	"jcourse_go/task/base"
	"jcourse_go/util"
	"log"
	"time"

	"github.com/RoaringBitmap/roaring"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const (
	TypeSaveDailyStatistic     = "statistic:save_daily"
	TypeRefreshPVDupJudge      = "statistic:pv:refresh_dup_judge"
	IntervalSaveDailyStatistic = 24 * time.Hour //"@every 24h"
	IntervalRefreshPVDupJudge  = 1 * time.Hour  //"@every 1h"
)

func HandleSaveStatisticTask(ctx context.Context, t base.Task) error {
	db := dal.GetDBClient()
	err := SaveStatistic(ctx, db, middleware.PV, middleware.UV, time.Now())
	if err != nil {
		log.Fatalf("failed to save statistic: %v", err)
		return err
	}
	return nil
}

func SaveStatistic(ctx context.Context, db *gorm.DB, pvm middleware.IPVMiddleware, uvm middleware.IUVMiddleware, datetime time.Time) error {
	curPVCount := pvm.GetPVCount()
	var curUVData *roaring.Bitmap
	var curUVCount uint64
	curUVData = uvm.GetUVData()
	curUVCount = curUVData.GetCardinality()
	today := util.FormatDate(time.Now())
	// save to db
	// 1. serialize bitmap to BLOB
	// 2. get different counts
	// 3. save statistic item
	// 4. save statisticData item
	log.Printf("today: %s", today)
	log.Printf("pv: %d, uv: %d", curPVCount, curUVCount)

	// TODO: better error handle
	userQuery := repository.NewUserQuery(db)
	reviewQuery := repository.NewReviewQuery(db)
	statisticQuery := repository.NewStatisticQuery(db)
	statisticDataQuery := repository.NewStatisticDataQuery(db)
	// get counts
	dayStart, dayEnd := util.GetDayTimeRange(datetime)
	option := repository.WithCreatedAtBetween(dayStart, dayEnd)
	newUserCount, err := userQuery.GetUserCount(ctx, option)
	if err != nil {
		log.Fatalf("failed to save statistic: %v", err)
		return err
	}
	newReviewCount, err := reviewQuery.GetReviewCount(ctx, option)
	if err != nil {
		log.Fatalf("failed to save statistic: %v", err)
		return err
	}
	if err != nil {
		log.Fatalf("failed to save statistic: %v", err)
		return err
	}
	totalReviewCount, err := reviewQuery.GetReviewCount(ctx)
	if err != nil {
		log.Fatalf("failed to save statistic: %v", err)
		return err
	}
	totalUserCount, err := userQuery.GetUserCount(ctx)
	if err != nil {
		log.Fatalf("failed to save statistic: %v", err)
		return err
	}
	// create statistic item
	newStatisticItem := po.StatisticPO{
		Date:         util.FormatDate(datetime),
		UVCount:      int64(curUVCount),
		PVCount:      curPVCount,
		NewReviews:   newReviewCount,
		NewUsers:     newUserCount,
		TotalReviews: totalReviewCount,
		TotalUsers:   totalUserCount,
	}
	err = statisticQuery.CreateStatistic(ctx, &newStatisticItem)
	if err != nil {
		log.Fatalf("failed to save statistic: %v", err)
		return err
	}
	// save uv data
	uvDataBytes, err := curUVData.ToBytes()
	if err != nil {
		log.Fatalf("failed to save statistic: %v", err)
		return err
	}
	if newStatisticItem.ID <= 0 {
		return errors.Errorf("failed to save statistic: Not write ID back")
	}
	err = statisticDataQuery.CreateData(ctx, &po.StatisticDataPO{
		Date:        util.FormatDate(datetime),
		UVData:      uvDataBytes,
		StatisticID: int64(newStatisticItem.ID),
	})
	if err != nil {
		return err
	}
	// reset only if succeed
	uvm.ResetUV()
	pvm.ResetPV()
	return nil
}

func HandleRefreshPVDupJudgeTask(ctx context.Context, t base.Task) error {
	middleware.UpdateDuplicateJudgeDuration(ctx)
	return nil
}
