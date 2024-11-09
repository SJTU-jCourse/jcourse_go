package middleware

import (
	"context"
	"jcourse_go/dal"
	"jcourse_go/repository"
	"jcourse_go/util"
	"log"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/RoaringBitmap/roaring"
	"github.com/gin-gonic/gin"
)

const (
	DefaultDuplicateJudgeDuration = time.Second
	DuplicateJudgeDurationKey     = "duplicate_judge_duration"
	DurationRefreshDupJudge       = 30 * time.Minute
)

func InitStatistic(r *gin.Engine) {
	timer := time.NewTicker(DurationRefreshDupJudge)
	go func() {
		for {
			<-timer.C
			UpdateDuplicateJudgeDuration(context.Background())
		}
	}()
}

var lastQuerySiteSettingTime time.Time
var lastQuerySiteSettingDuration time.Duration = DefaultDuplicateJudgeDuration

func GetDuplicateJudgeDuration() time.Duration {
	return lastQuerySiteSettingDuration
}
func UpdateDuplicateJudgeDuration(ctx context.Context) time.Duration {
	lastQuerySiteSettingTime = time.Now()
	siteQuery := repository.NewSettingQuery(dal.GetDBClient())
	siteSetting, err := siteQuery.GetSetting(ctx, DuplicateJudgeDurationKey)
	if err != nil {
		lastQuerySiteSettingDuration = DefaultDuplicateJudgeDuration
		return DefaultDuplicateJudgeDuration
	}
	return siteSetting.GetValue().(time.Duration)
}

func (u *UserRequestCache) IsDuplicate(req string) bool {
	u.mutex.RLock()
	defer u.mutex.RUnlock()
	lastReq, ok := u.lastReqMap[req]
	if !ok {
		return false
	}
	return time.Since(lastReq) < GetDuplicateJudgeDuration()
}

func IsRobot(c *gin.Context) bool {
	return false
}

func IsInternalRequest(c *gin.Context) bool {
	return false
}

type AbstractPaginateRequest struct {
	Page     int64 `json:"page" form:"page"`
	PageSize int64 `json:"page_size" form:"page_size"`
}

func IsPaginateRequest(c *gin.Context) int64 {
	var req AbstractPaginateRequest
	if c.ShouldBind(&req) != nil {
		return -1
	}
	return req.Page
}

// 如果只用单个bitmap, 可以不考虑并发问题，并发set一样是ok的，省去了mutex的开销
// 可以考虑后续优化
var rbm = roaring.New()
var rbmMutex sync.RWMutex

func UVStatisticWithLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 在登录的时候，将用户id加入到bitmap中
		// 这样避免了在每次请求的时候都去拿锁
		// 但需要处理用户一直保持登录状态的情况, 目前还没想好怎么处理
		user := GetCurrentUser(c)

		if user == nil {
			// 对应登录的情况
			c.Next()
			user = GetCurrentUser(c)
			if user == nil {
				return
			}
			rbmMutex.Lock()
			defer rbmMutex.Unlock()
			rbm.Add(uint32(GetCurrentUser(c).ID)) // 如果id需要用64位int,则修改rbm;add 自带去重
		} else {
			c.Next()
		}
	}
}
func UVStatistic() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 简单加放锁的方式，性能约为10000次请求7ms
		user := GetCurrentUser(c)
		if user == nil {
			log.Printf("user is nil")
			return
		}
		rbmMutex.Lock()
		defer rbmMutex.Unlock()
		rbm.Add(uint32(GetCurrentUser(c).ID)) // 如果id需要用64位int,则修改rbm;add 自带去重
		c.Next()
	}
}
func GetUVCount() uint64 {
	rbmMutex.RLock()
	defer rbmMutex.RUnlock()
	return rbm.GetCardinality()
}

var pvCount atomic.Int64

type UserRequestCache struct {
	mutex      sync.RWMutex
	lastReqMap map[string]time.Time
}

func (u *UserRequestCache) Request(req string) {
	u.mutex.Lock()
	defer u.mutex.Unlock()
	u.lastReqMap[req] = time.Now()
}
func ShouldBeCountInPV(c *gin.Context) bool {
	if c.Request.Method != "GET" {
		return false
	}
	if IsRobot(c) || IsInternalRequest(c) {
		return false
	}
	return SetIfNoDuplicate(c)
}

var requestUserCache = make(map[int64]*UserRequestCache)
var requestUserCacheMutex sync.RWMutex

func PVStatistic() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Next()
		if !ShouldBeCountInPV(c) {
			return
		}
		pvCount.Add(1)
	}
}
func SetIfNoDuplicate(c *gin.Context) bool {
	// 目前简单将不同URI的请求和同一个paginate URI下不同page的请求视为不同请求
	user := GetCurrentUser(c)
	if user == nil {
		return false
	}
	requestUserCacheMutex.RLock()
	userCache, ok := requestUserCache[user.ID]
	requestUserCacheMutex.RUnlock()
	// HINT: 2-stage check
	if !ok {
		requestUserCacheMutex.Lock()
		defer requestUserCacheMutex.Unlock()
		userCache, ok = requestUserCache[user.ID]
		if !ok {
			userCache = &UserRequestCache{lastReqMap: make(map[string]time.Time)}
			requestUserCache[user.ID] = userCache
		}
	}
	req := c.Request.RequestURI
	pageNo := IsPaginateRequest(c)
	if pageNo > 0 {
		req = c.Request.RequestURI + strconv.FormatInt(pageNo, 10)
	}
	if userCache.IsDuplicate(req) {
		return false
	}
	userCache.Request(req)
	return true
}
func GetPVMap() map[int64]*UserRequestCache {
	return requestUserCache
}
func GetPVCount() int64 {
	return pvCount.Load()
}

func SaveStatistic() {
	log.Printf("save statistic")

	curPVCount := GetPVCount()
	var curUVData *roaring.Bitmap
	var curUVCount uint64
	{
		rbmMutex.Lock()
		defer rbmMutex.Unlock()
		curUVData = rbm.Clone()
		curUVCount = curUVData.GetCardinality()
		rbm.Clear()
	}

	today := util.FormatDate(time.Now())
	// save to db
	// 1. serialize bitmap to BLOB
	// 2. save count item
	// 3. read the prev day's count, and calculate dau
	// 4. save dau item
	log.Printf("today: %s", today)
	log.Printf("pv: %d, uv: %d", curPVCount, curUVCount)
	log.Printf("uv data(first 10): %v", curUVData.ToArray()[:10])

	// TODO: save uv, pv, dau to db

}
func ScheduleSaveStatistic(gap time.Duration, stopChan <-chan struct{}) {
	timer := time.NewTicker(gap)
	go func() {
		for {
			select {
			case <-timer.C:
				SaveStatistic()
			// elegant quit: 关闭时正常退出
			case <-stopChan:
				log.Printf("stopped statistic saver")
				timer.Stop()
				return
			}
		}
	}()
}
