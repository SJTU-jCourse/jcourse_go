package middleware

import (
	"context"
	"jcourse_go/dal"
	"jcourse_go/model/model"
	"jcourse_go/repository"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/RoaringBitmap/roaring"
	"github.com/gin-gonic/gin"
)

var UV = NewUVMiddleware()
var PV = NewPVMiddleware()

const (
	DefaultDuplicateJudgeDuration = 10 * time.Minute
	DuplicateJudgeDurationKey     = "duplicate_judge_duration"
)

var LastQuerySiteSettingTime time.Time
var LastQuerySiteSettingDuration time.Duration = DefaultDuplicateJudgeDuration

func UpdateDuplicateJudgeDuration(ctx context.Context) time.Duration {
	LastQuerySiteSettingTime = time.Now()
	siteQuery := repository.NewSettingQuery(dal.GetDBClient())
	siteSetting, err := siteQuery.GetSetting(ctx, DuplicateJudgeDurationKey)
	if err != nil {
		log.Printf("failed to get site setting: %v", err)
		LastQuerySiteSettingDuration = DefaultDuplicateJudgeDuration
		return DefaultDuplicateJudgeDuration
	}
	// assert setting type is string
	duration, err := time.ParseDuration(siteSetting.GetValue().(string))
	if err != nil {
		log.Printf("failed to parse setting: %v", siteSetting.GetValue())
		LastQuerySiteSettingDuration = DefaultDuplicateJudgeDuration
		return DefaultDuplicateJudgeDuration
	}
	if duration > 0 {
		log.Printf("update duration to: %v", duration)
		LastQuerySiteSettingDuration = duration
		return duration
	}
	LastQuerySiteSettingDuration = DefaultDuplicateJudgeDuration
	return DefaultDuplicateJudgeDuration
}
func (u *UserRequestCache) IsDuplicate(req string) bool {
	u.mutex.RLock()
	defer u.mutex.RUnlock()
	lastReq, ok := u.lastReqMap[req]
	if !ok {
		return false
	}
	return time.Since(lastReq) < LastQuerySiteSettingDuration
}

func IsRobot(c *gin.Context) bool {
	return false
}

func IsInternalRequest(c *gin.Context) bool {
	return false
}

func IsPaginateRequest(c *gin.Context) int64 {
	var req model.PaginationFilterForQuery
	if c.ShouldBind(&req) != nil {
		return -1
	}
	return req.Page
}

type IUVMiddleware interface {
	UVStatistic() gin.HandlerFunc
	UVStatisticMock() gin.HandlerFunc // // 用于测试, 仅删去底部的c.Next()
	GetUVCount() uint64
	GetUVData() *roaring.Bitmap
	ResetUV()
	ContainsUser(id int64) bool
	AddUser(id int64) // 用于测试, 实际上的逻辑在UVStatistic中自动完成
}

func NewUVMiddleware() IUVMiddleware {
	return &UVMiddleware{
		Rbm:      roaring.New(),
		RbmMutex: sync.RWMutex{},
	}
}

type UVMiddleware struct {
	Rbm      *roaring.Bitmap
	RbmMutex sync.RWMutex
}

func (m *UVMiddleware) AddUser(id int64) {
	m.RbmMutex.Lock()
	defer m.RbmMutex.Unlock()
	m.Rbm.Add(uint32(id))
}

func (m *UVMiddleware) ContainsUser(id int64) bool {
	m.RbmMutex.RLock()
	defer m.RbmMutex.RUnlock()
	return m.Rbm.Contains(uint32(id))
}

// 如果只用单个bitmap, 可以不考虑并发问题，并发set一样是ok的，省去了mutex的开销
// 可以考虑后续优化

func (m *UVMiddleware) UVStatisticWithLogin() gin.HandlerFunc {
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
			m.AddUser(user.ID)
		} else {
			c.Next()
		}
	}
}
func (m *UVMiddleware) UVStatisticMock() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 简单加放锁的方式，性能约为10000次请求7ms
		user := GetCurrentUser(c)
		if user == nil {
			return
		}
		m.AddUser(user.ID) // 如果id需要用64位int,则修改rbm;add 自带去重
	}
}
func (m *UVMiddleware) UVStatistic() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 简单加放锁的方式，性能约为10000次请求7ms
		user := GetCurrentUser(c)
		if user == nil {
			return
		}
		m.AddUser(user.ID) // 如果id需要用64位int,则修改rbm;add 自带去重
		c.Next()
	}
}
func (m *UVMiddleware) GetUVCount() uint64 {
	m.RbmMutex.RLock()
	defer m.RbmMutex.RUnlock()
	return m.Rbm.GetCardinality()
}
func (m *UVMiddleware) GetUVData() *roaring.Bitmap {
	m.RbmMutex.RLock()
	defer m.RbmMutex.RUnlock()
	return m.Rbm.Clone()
}
func (m *UVMiddleware) ResetUV() {
	m.RbmMutex.Lock()
	defer m.RbmMutex.Unlock()
	m.Rbm.Clear()
}

type IPVMiddleware interface {
	PVStatistic() gin.HandlerFunc
	PVStatisticMock() gin.HandlerFunc
	GetPVCount() int64
	GetPVCache() map[int64]*UserRequestCache
	ClearPVCache()
	ResetPV()
	AddPageView(userID int64, uri string) // 用于测试, 实际上的逻辑在PVStatistic中自动完成
}
type UserRequestCache struct {
	mutex      sync.RWMutex
	lastReqMap map[string]time.Time
}

func (u *UserRequestCache) Request(req string) {
	u.mutex.Lock()
	defer u.mutex.Unlock()
	u.lastReqMap[req] = time.Now()
}

type PVMiddleware struct {
	PVCount               atomic.Int64
	RequestUserCache      map[int64]*UserRequestCache
	RequestUserCacheMutex sync.RWMutex
}

func (m *PVMiddleware) AddPageView(userID int64, req string) {
	m.RequestUserCacheMutex.Lock()
	defer m.RequestUserCacheMutex.Unlock()
	userCache, ok := m.RequestUserCache[userID]
	if !ok {
		userCache = &UserRequestCache{lastReqMap: make(map[string]time.Time)}
		m.RequestUserCache[userID] = userCache
	}
	userCache.Request(req)
	m.PVCount.Add(1)
}

func (m *PVMiddleware) ResetPV() {
	m.PVCount.Store(0)
	m.ClearPVCache()
}

func NewPVMiddleware() IPVMiddleware {
	return &PVMiddleware{
		PVCount:               atomic.Int64{},
		RequestUserCache:      make(map[int64]*UserRequestCache),
		RequestUserCacheMutex: sync.RWMutex{},
	}
}

func (m *PVMiddleware) ShouldBeCountInPV(c *gin.Context) bool {
	if c.Request.Method != "GET" {
		return false
	}
	if IsRobot(c) || IsInternalRequest(c) {
		return false
	}
	return m.SetIfNoDuplicate(c)
}

func (m *PVMiddleware) ClearPVCache() {
	m.RequestUserCacheMutex.Lock()
	defer m.RequestUserCacheMutex.Unlock()
	m.RequestUserCache = make(map[int64]*UserRequestCache)
}
func (m *PVMiddleware) PVStatistic() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Next()
		if !m.ShouldBeCountInPV(c) {
			return
		}
		m.PVCount.Add(1)
	}
}
func (m *PVMiddleware) PVStatisticMock() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !m.ShouldBeCountInPV(c) {
			return
		}
		m.PVCount.Add(1)
	}
}

func (m *PVMiddleware) SetIfNoDuplicate(c *gin.Context) bool {
	// 目前简单将不同URI的请求和同一个paginate URI下不同page的请求视为不同请求
	user := GetCurrentUser(c)
	if user == nil {
		return false
	}
	m.RequestUserCacheMutex.RLock()
	userCache, ok := m.RequestUserCache[user.ID]
	m.RequestUserCacheMutex.RUnlock()
	// HINT: 2-stage check
	if !ok {
		m.RequestUserCacheMutex.Lock()
		defer m.RequestUserCacheMutex.Unlock()
		userCache, ok = m.RequestUserCache[user.ID]
		if !ok {
			userCache = &UserRequestCache{lastReqMap: make(map[string]time.Time)}
			m.RequestUserCache[user.ID] = userCache
		}
	}
	req := c.Request.RequestURI
	if userCache.IsDuplicate(req) {
		return false
	}
	userCache.Request(req)
	return true
}
func (m *PVMiddleware) GetPVCache() map[int64]*UserRequestCache {
	return m.RequestUserCache
}
func (m *PVMiddleware) GetPVCount() int64 {
	return m.PVCount.Load()
}
