package statistic

import (
	"context"
	"log"
	"time"

	"github.com/RoaringBitmap/roaring"

	"jcourse_go/internal/interface/task/base"
	"jcourse_go/internal/interface/web/middleware"
	"jcourse_go/internal/service"
	"jcourse_go/pkg/util"
)

const (
	TypeSaveDailyStatistic     = "statistic:save_daily"
	TypeRefreshPVDupJudge      = "statistic:pv:refresh_dup_judge"
	IntervalSaveDailyStatistic = 24 * time.Hour // "@every 24h"
	IntervalRefreshPVDupJudge  = 1 * time.Hour  // "@every 1h"
)

func HandleSaveStatisticTask(ctx context.Context, t base.Task) error {
	err := SaveStatistic(ctx, middleware.PV, middleware.UV, time.Now())
	if err != nil {
		log.Fatalf("failed to save statistic: %v", err)
		return err
	}
	return nil
}

func SaveStatistic(ctx context.Context, pvm middleware.IPVMiddleware, uvm middleware.IUVMiddleware, datetime time.Time) error {
	curPVCount := pvm.GetPVCount()
	var curUVData *roaring.Bitmap
	var curUVCount int64
	curUVData = uvm.GetUVData()
	curUVCount = int64(curUVData.GetCardinality())
	today := util.FormatDate(datetime)
	// save to db
	// 1. serialize bitmap to BLOB
	// 2. get different counts
	// 3. save statistic item
	// 4. save statisticData item
	log.Printf("today: %s", today)
	log.Printf("pv: %d, uv: %d", curPVCount, curUVCount)

	// TODO: better error handle
	err := service.CreateStatistic(ctx, datetime, curPVCount, curUVCount, curUVData)
	if err != nil {
		return err
	}
	// reset only if succeed
	uvm.ResetUV()
	pvm.ResetPV()
	return nil
}

/*
func HandleRefreshPVDupJudgeTask(ctx context.Context, t base.Task) error {
	middleware.UpdateDuplicateJudgeDuration(ctx)
	return nil
}
*/
