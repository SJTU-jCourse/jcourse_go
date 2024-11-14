package task

import (
	"context"
	"jcourse_go/middleware"
	"jcourse_go/model/po"
	"jcourse_go/repository"
	"jcourse_go/util"
	"log"
	"math"
	"time"

	"gorm.io/gorm"

	"github.com/RoaringBitmap/roaring"
	"github.com/pkg/errors"
)

type DailyTriggerOff = time.Duration

// CalOffsetByDailyTrigger 计算从现在到下一次每日触发时间的偏移量
func CalOffsetByDailyTrigger(triggerOff DailyTriggerOff) time.Duration {
	now := time.Now()
	triggerTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).Add(triggerOff)
	if now.After(triggerTime) {
		triggerTime = triggerTime.AddDate(0, 0, 1)
	}
	return triggerTime.Sub(now)
}

const (
	DurationRefreshDupJudge = 30 * time.Minute
	// 2:33 AM 保存前一日的统计数据
	StatisticSaveTriggerOff DailyTriggerOff = 2*time.Hour + 33*time.Minute
	TimeCheckFailed                         = "time check failed"
)

// runScheduledTask 运行定时任务, 仅适用于执行的任务没有变化的情况(传递原型而非闭包)
func runScheduledTask(ctx context.Context, tasks []ScheduledTask) {
	const maxTaskQueueSize = 10000
	taskQueue := make(chan int64, maxTaskQueueSize)
	timers := make(map[int64]*time.Ticker, len(tasks))
	defer func() {
		for _, timer := range timers {
			if timer != nil {
				timer.Stop()
			}
		}
	}()
	startTimer := func(id int64) {
		for {
			select {
			case <-timers[id].C:
				curTime := time.Now()
				tasks[id].PrevTime = &curTime
				taskQueue <- id
			case <-ctx.Done():
				timers[id].Stop()
				return
			}
		}
	}
	for _, task := range tasks {
		if task.Offset < 0 {
			// no offset, start immediately
			curTime := time.Now()
			timers[task.ID] = time.NewTicker(task.Interval)
			tasks[task.ID].StartTime = &curTime
			go startTimer(task.ID)
		} else {
			// has offset, start after offset
			timers[task.ID] = nil
			time.AfterFunc(task.Offset, func() {
				curTime := time.Now()
				timers[task.ID] = time.NewTicker(task.Interval)
				tasks[task.ID].StartTime = &curTime
				go startTimer(task.ID)
			})
		}
	}
	// task event loop
	for {
		select {
		case id := <-taskQueue:
			task := tasks[id]
			go func() {
				err := task.TaskFunc(ctx, tasks[id].ScheduledTaskParams)
				if err != nil && task.ErrorHandler != nil {
					task.ErrorHandler(err, tasks[id].ScheduledTaskParams)
				}
			}()
		case <-ctx.Done():
			close(taskQueue)
			return
		}
	}
}

type ScheduledTaskParams struct {
	StartTime *time.Time // 在第一次触发时回写, 用于校验
	PrevTime  *time.Time // 上一次触发时间, 在每次触发时回写
	Interval  time.Duration
	Name      string
	Offset    time.Duration
}

func NewScheduledTaskParams() ScheduledTaskParams {
	return ScheduledTaskParams{
		StartTime: nil,
		PrevTime:  nil,
	}
}

type ScheduledTask struct {
	ID int64
	// 多次触发时可以通过param得到一些之前触发结果的信息, 也可以得到任务元信息，如名称
	TaskFunc     func(ctx context.Context, param ScheduledTaskParams) error
	ErrorHandler func(err error, param ScheduledTaskParams)
	ScheduledTaskParams
}

func InitStatistic(db *gorm.DB) context.CancelFunc {
	ctx, cancel := context.WithCancel(context.Background())
	saveStatisticOffset := CalOffsetByDailyTrigger(StatisticSaveTriggerOff)
	tasks := []ScheduledTask{
		{
			ID: 0, // ID: index in tasks
			TaskFunc: func(ctx context.Context, param ScheduledTaskParams) error {
				// 本次触发和理论触发时间差超过1小时，抛出错误
				if param.PrevTime != nil && math.Abs(float64(time.Since((*param.PrevTime).Add(param.Interval)))) > float64(time.Hour) {
					return errors.Errorf("%v: %v", param.Name, TimeCheckFailed)
				}
				return SaveStatistic(ctx, db, middleware.PV, middleware.UV, time.Now())
			},
			ErrorHandler: func(err error, param ScheduledTaskParams) {
				log.Fatalf("Error in %v: %v", param.Name, err)
			},
			ScheduledTaskParams: ScheduledTaskParams{
				Interval:  24 * time.Hour,
				Offset:    saveStatisticOffset,
				Name:      "SaveStatistic",
				StartTime: nil,
				PrevTime:  nil,
			},
		},
		{
			ID: 1,
			TaskFunc: func(ctx context.Context, param ScheduledTaskParams) error {
				middleware.UpdateDuplicateJudgeDuration(ctx)
				return nil
			},
			ErrorHandler: func(err error, param ScheduledTaskParams) {
				log.Printf("Error in %v: %v", param.Name, err)
			},
			ScheduledTaskParams: ScheduledTaskParams{
				Name:      "UpdateDuplicateJudgeDuration",
				Interval:  DurationRefreshDupJudge,
				Offset:    -1,
				StartTime: nil,
				PrevTime:  nil,
			},
		},
	}

	go runScheduledTask(ctx, tasks)
	return cancel
}

func SaveStatistic(ctx context.Context, db *gorm.DB, pvm middleware.IPVMiddleware, uvm middleware.IUVMiddleware, datetime time.Time) error {
	curPVCount := pvm.GetPVCount()
	var curUVData *roaring.Bitmap
	var curUVCount uint64
	curUVData = uvm.GetUVData()
	curUVCount = curUVData.GetCardinality()
	middleware.UV.ResetUV()

	today := util.FormatDate(time.Now())
	// save to db
	// 1. serialize bitmap to BLOB
	// 2. get different counts
	// 3. save statistic item
	// 4. save statisticData item
	log.Printf("today: %s", today)
	log.Printf("pv: %d, uv: %d", curPVCount, curUVCount)

	// TODO: save uv, pv, dau to db
	userQuery := repository.NewUserQuery(db)
	courseQuery := repository.NewCourseQuery(db)
	reviewQuery := repository.NewReviewQuery(db)
	statisticQuery := repository.NewStatisticQuery(db)
	statisticDataQuery := repository.NewStatisticDataQuery(db)
	// get counts
	dayStart, dayEnd := util.GetDayTimeRange(datetime)
	option := repository.WithCreatedAtBetween(dayStart, dayEnd)
	newUserCount, err := userQuery.GetUserCount(ctx, option)
	if err != nil {
		// TODO: save error handle
		log.Fatalf("failed to save statistic: %v", err)
		return err
	}
	newReviewCount, err := reviewQuery.GetReviewCount(ctx, option)
	if err != nil {
		// TODO: save error handle
		log.Fatalf("failed to save statistic: %v", err)
		return err
	}
	newCourseCount, err := courseQuery.GetCourseCount(ctx, option)
	if err != nil {
		// TODO: save error handle
		log.Fatalf("failed to save statistic: %v", err)
		return err
	}
	totalCourseCount, err := courseQuery.GetCourseCount(ctx)
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
		Date:         datetime,
		UVCount:      int64(curUVCount),
		PVCount:      curPVCount,
		NewCourses:   newCourseCount,
		NewReviews:   newReviewCount,
		NewUsers:     newUserCount,
		TotalCourses: totalCourseCount,
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
		Date:        datetime,
		UVData:      uvDataBytes,
		StatisticID: int64(newStatisticItem.ID),
	})
	if err != nil {
		return err
	}
	return nil
}
