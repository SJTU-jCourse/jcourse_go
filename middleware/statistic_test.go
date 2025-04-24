package middleware

import (
	"context"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"sync"
	"testing"
	"time"

	"jcourse_go/internal/constant"
	"jcourse_go/model/model"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/RoaringBitmap/roaring"
)

type userKeyType = struct{ key string } // make ci happy
func TestUVStatistic(t *testing.T) {
	t.Run("TestRBMAddDup", func(t *testing.T) {
		testRbm := roaring.New()
		testRbm.Add(1)
		testRbm.Add(1)
		if testRbm.GetCardinality() != 1 {
			t.Errorf("RBM Add duplicate failed")
		}
	})

	t.Run("TestUV", func(t *testing.T) {
		testNum := 10000
		var testParams []*gin.Context
		for i := 0; i < testNum; i++ {
			c := &gin.Context{}
			user := &model.UserDetail{
				UserMinimal: model.UserMinimal{
					ID: int64(i),
				},
			}
			c.Set(constant.CtxKeyUser, user)
			testParams = append(testParams, c)
		}

		// Test UV
		m := NewUVMiddleware()
		handler := m.UVStatisticMock()
		var wg sync.WaitGroup
		for i := 0; i < testNum; i++ {
			wg.Add(1)
			go func() {
				handler(testParams[i])
				wg.Done()
			}()
		}
		wg.Wait()
		assert.Equal(t, uint64(testNum), m.GetUVCount())
		for i := 0; i < testNum; i++ {
			assert.True(t, m.ContainsUser(int64(i)))
		}
		for i := testNum; i < testNum*2; i++ {
			assert.False(t, m.ContainsUser(int64(i)))
		}
	})
}

func BenchmarkRbm(b *testing.B) {
	// 10000 users: 1.9 ns/op
	b.Run("CountSeq", func(b *testing.B) {
		rbm := roaring.New()
		userNum := 10000
		for j := 0; j < userNum; j++ {
			rbm.Add(uint32(j))
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			rbm.GetCardinality()
		}
	})
	// 10000 users: 40 us/op
	b.Run("CountRandom", func(b *testing.B) {
		rbm := roaring.New()
		userNum := 10000
		for j := 0; j < userNum; j++ {
			rbm.Add(rand.Uint32())
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			rbm.GetCardinality()
		}
	})
	// 10000 users: 4.2 us/op
	b.Run("CloneSeq", func(b *testing.B) {
		rbm := roaring.New()
		userNum := 10000
		for j := 0; j < userNum; j++ {
			rbm.Add(uint32(j))
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			rbm.Clone()
		}
	})
	// 10000 users: 2.3 ms/op
	b.Run("CloneRandom", func(b *testing.B) {
		rbm := roaring.New()
		userNum := 10000
		for j := 0; j < userNum; j++ {
			rbm.Add(rand.Uint32())
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			rbm.Clone()
		}
	})
	b.Run("SerializeSeq", func(b *testing.B) {
		// 10000 users: 11 us/op
		rbm := roaring.New()
		userNum := 10000
		for j := 0; j < userNum; j++ {
			rbm.Add(uint32(j))
		}
		tmpFile, _ := os.Create("tmp")
		defer func() {
			err := tmpFile.Close()
			if err != nil {
				b.Errorf("Close tmp file failed")
			}
			err = os.Remove("tmp")
			if err != nil {
				b.Errorf("Remove tmp file failed")
			}
		}()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			bytes, err := rbm.ToBytes()
			if err != nil {
				b.Errorf("Serialize failed")
			}
			_, err = tmpFile.Write(bytes)
			if err != nil {
				b.Errorf("Write failed")
			}
		}
	})
	b.Run("SerializeRandom", func(b *testing.B) {
		// 10000 users: 1.3 ms/op
		rbm := roaring.New()
		userNum := 10000
		for j := 0; j < userNum; j++ {
			rbm.Add(rand.Uint32())
		}
		tmpFile, _ := os.Create("tmp")
		defer func() {
			err := tmpFile.Close()
			if err != nil {
				b.Errorf("Close tmp file failed")
			}
			err = os.Remove("tmp")
			if err != nil {
				b.Errorf("Remove tmp file failed")
			}
		}()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			bytes, err := rbm.ToBytes()
			if err != nil {
				b.Errorf("Serialize failed")
			}
			_, err = tmpFile.Write(bytes)
			if err != nil {
				b.Errorf("Write failed")
			}
		}
	})
}

func BenchmarkUVStatistic(b *testing.B) {
	// prepare test data
	// 10000 users: 6.7 ms/op
	userNum := 10000
	testCtxs := make([]*gin.Context, userNum)
	for i := 0; i < userNum; i++ {
		c := &gin.Context{}
		user := &model.UserDetail{
			UserMinimal: model.UserMinimal{
				ID: int64(i),
			},
		}
		c.Set(constant.CtxKeyUser, user)
		testCtxs[i] = c
	}
	m := NewUVMiddleware()
	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup
		for _, c := range testCtxs {
			fn := m.UVStatisticMock()
			wg.Add(1)
			go func() {
				defer wg.Done()
				fn(c)
			}()
		}
		wg.Wait()
		// t.Logf("UVStatistic: %d", rbm.GetCardinality())
	}
}

func GenerateRandomRequest(N int, paths []string) []*http.Request {
	// 定义一些测试用的参数
	baseURL := "http://example.com"
	methods := []string{"GET", "POST", "PUT", "DELETE"}
	type queryParamKV struct {
		key   string
		value string
	}
	queryParamKVs := []queryParamKV{
		{"key1", "value1"},
		{"page", "1"},
	}

	reqs := make([]*http.Request, N)
	// 随机生成请求
	for i := 0; i < N; i++ {
		path := paths[rand.Intn(len(paths))]
		method := methods[rand.Intn(len(methods))]
		queryParam := queryParamKVs[rand.Intn(len(queryParamKVs))]

		// 构建 URL
		u, _ := url.Parse(baseURL)
		u.Path = path
		query := u.Query()
		query.Set(queryParam.key, queryParam.value)
		u.RawQuery = query.Encode()

		// 创建 http.Request 对象
		req, _ := http.NewRequest(method, u.String(), nil)
		reqs[i] = req
	}
	return reqs
}

func TestPVStatistic(t *testing.T) {
	testPaths := []string{
		"/api/user/login",
		"/api/user/logout",
		"/api/user/register",
		"/api/user/common",
		"/api/teacher",
		"/api/teacher/filter",
		"/api/teacher/1",
		"/api/base_course/1",
		"/api/course",
		"/api/course/filter",
		"/api/course/1",
		"/api/training_plan",
		"/api/training_plan/filter",
	}
	testNum := 10000
	var testParams []*gin.Context
	testReqs := GenerateRandomRequest(testNum, testPaths)
	for i := 0; i < testNum; i++ {
		c := &gin.Context{}
		user := &model.UserDetail{
			UserMinimal: model.UserMinimal{
				ID: int64(i),
			},
		}
		c.Set(constant.CtxKeyUser, user)
		c.Request = testReqs[i]
		testParams = append(testParams, c)
	}

	m := NewPVMiddleware()
	// Test PV
	handler := m.PVStatisticMock()
	var wg sync.WaitGroup
	for i := 0; i < testNum; i++ {
		wg.Add(1)
		go func() {
			handler(testParams[i])
			wg.Done()
		}()
	}
	wg.Wait()
	t.Logf("PVStatistic: %d", m.GetPVCount())
	t.Logf("It should be approximately to %d", testNum/4) // 4 is the number of different methods
}
func BenchmarkPVStatistic(b *testing.B) {
	testPaths := []string{
		"/api/user/login",
		"/api/user/logout",
		"/api/user/register",
		"/api/user/common",
		"/api/teacher",
		"/api/teacher/filter",
		"/api/teacher/1",
		"/api/base_course/1",
		"/api/course",
		"/api/course/filter",
		"/api/course/1",
		"/api/training_plan",
		"/api/training_plan/filter",
	}
	testNum := 10000 * 4
	var testParams []*gin.Context
	testReqs := GenerateRandomRequest(testNum, testPaths)
	for i := 0; i < testNum; i++ {
		c := &gin.Context{}
		user := &model.UserDetail{
			UserMinimal: model.UserMinimal{
				ID: int64(i),
			},
		}
		c.Set(constant.CtxKeyUser, user)
		c.Request = testReqs[i]
		testParams = append(testParams, c)
	}

	// Test PV
	m := NewPVMiddleware()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.ClearPVCache()
		handler := m.PVStatisticMock()
		var wg sync.WaitGroup
		for j := 0; j < testNum; j++ {
			wg.Add(1)
			go func() {
				handler(testParams[j])
				wg.Done()
			}()
		}
		wg.Wait()
	}

}
func TestCtxPass(t *testing.T) {
	t.Run("TestCtxPass", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		r := gin.Default()

		// 模拟登录中间件
		mockLogin := func(c *gin.Context) {
			// 从gin.Context中获取请求的上下文
			reqCtx := c.Request.Context()
			userKey := userKeyType{key: constant.CtxKeyUser}
			user, ok := reqCtx.Value(userKey).(*model.UserDetail)
			assert.True(t, ok)

			// 将用户信息设置回gin.Context的上下文中
			c.Set(constant.CtxKeyUser, user)
			c.Next()
		}
		r.Use(mockLogin)
		r.GET("/some-path", func(c *gin.Context) {
			_, ok := c.Get(constant.CtxKeyUser)
			assert.True(t, ok)
			c.Status(http.StatusOK)
		})
		req := httptest.NewRequest(http.MethodGet, "/some-path", nil)
		var userKey = userKeyType{key: constant.CtxKeyUser}
		req = req.WithContext(context.WithValue(context.Background(), userKey, &model.UserDetail{}))

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
	})
}

func GetTestGinEngine() *gin.Engine {
	gin.SetMode(gin.TestMode)
	gin.DefaultWriter = io.Discard // disable gin's log output
	r := gin.New()
	r.Use(gin.Recovery())
	return r
}

// 只能模拟在c.Next()之前的完成放锁的高并发情况, 否则中间件时延依赖于具体业务时延
func SimulateHighConcurrency(pvMock gin.HandlerFunc, uvMock gin.HandlerFunc, qps int, seconds int, initUserNum int, totalUserNum int, testFn func(b *testing.B, iter int)) func(b *testing.B) time.Duration {
	return func(b *testing.B) time.Duration {
		r := GetTestGinEngine()
		middlewares := make([]gin.HandlerFunc, 0)
		// 因为ServeHTTP只传入Request, 所以将登录的context放在Request中, 需要取出放回context
		mockLogin := func(c *gin.Context) {
			reqCtx := c.Request.Context()
			userKey := userKeyType{key: constant.CtxKeyUser}
			user, ok := reqCtx.Value(userKey).(*model.UserDetail)
			assert.True(b, ok)
			c.Set(constant.CtxKeyUser, user)
			c.Next()
		}
		endMiddleWare := func(c *gin.Context) { c.Status(http.StatusOK) }
		middlewares = append(middlewares, mockLogin)
		middlewares = append(middlewares, pvMock)
		middlewares = append(middlewares, pvMock)
		middlewares = append(middlewares, endMiddleWare) // 避免c.Next()带来一些问题
		for _, m := range middlewares {
			r.Use(m)
		}

		// prepare test data
		testCtxs := make([]*gin.Context, totalUserNum)
		testPaths := []string{
			"/api/user/login",
			"/api/user/logout",
			"/api/user/register",
			"/api/user/common",
			"/api/teacher",
			"/api/teacher/filter",
			"/api/teacher/1",
			"/api/base_course/1",
			"/api/course",
			"/api/course/filter",
			"/api/course/1",
			"/api/training_plan",
			"/api/training_plan/filter",
		}
		testReqs := GenerateRandomRequest(totalUserNum, testPaths)
		mockHandler := func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{})
		}
		w := httptest.NewRecorder()
		for _, path := range testPaths {
			for _, method := range []string{"GET", "POST", "PUT", "DELETE"} {
				r.Handle(method, path, mockHandler)
			}
		}

		for i := 0; i < totalUserNum; i++ {
			c := &gin.Context{}
			user := &model.UserDetail{
				UserMinimal: model.UserMinimal{
					ID: int64(i),
				},
			}
			userKey := userKeyType{key: constant.CtxKeyUser}
			reqCtx := context.WithValue(context.Background(), userKey, user)
			c.Request = testReqs[i].WithContext(reqCtx)
			testCtxs[i] = c
		}
		// login initUserNum users
		var wg sync.WaitGroup
		for i := 0; i < initUserNum; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				c := testCtxs[i]
				r.ServeHTTP(w, c.Request)
			}()
		}
		wg.Wait()
		ticker := time.NewTicker(time.Second / 10)
		// HINT: 默认运行一秒，此时无效果，所以/10
		// simulate count seconds
		totalTime := time.Second * time.Duration(seconds)
		finTicker := time.NewTicker(totalTime)
		defer func() {
			finTicker.Stop()
			ticker.Stop()
		}()
		iter := 0
		for {
			select {
			// mock qps(to simulate hard to get lock)
			// also mock login if initUserNum < totalUserNum
			case <-ticker.C:
				for j := 0; j < qps/10; j++ {
					c := testCtxs[rand.Intn(totalUserNum)]
					go r.ServeHTTP(httptest.NewRecorder(), c.Request)
				}
			case <-finTicker.C:
				b.StopTimer()
				return totalTime / time.Duration(iter)
			default:
				testFn(b, iter)
				iter += 1
			}
		}
	}
}

func BenchmarkGetPVCount(b *testing.B) {
	// <23 us / op
	b.Run("TestHighConcurrency", func(b *testing.B) {
		pvm := NewPVMiddleware()
		countLog := map[int64]bool{}
		mu := sync.Mutex{}
		GetPVCountWrapper := func(b *testing.B, iter int) {
			count := pvm.GetPVCount()
			b.StopTimer()
			mu.Lock()
			countLog[count] = true
			mu.Unlock()
			b.StartTimer()
		}
		uvm := NewUVMiddleware()
		avgTime := SimulateHighConcurrency(uvm.UVStatisticMock(), pvm.PVStatisticMock(), 1000, 5,
			10000, 10000, GetPVCountWrapper)(b)
		mu.Lock()
		defer mu.Unlock()
		for k := range countLog {
			b.Logf("PV Count: %d", k)
		}
		b.Logf("Average time: %v", avgTime)
	})

}
func BenchmarkGetUVCount(b *testing.B) {
	// <23 us / op
	b.Run("TestHighConcurrency", func(b *testing.B) {
		pvm := NewPVMiddleware()
		uvm := NewUVMiddleware()
		countLog := map[int64]bool{}
		mu := sync.Mutex{}
		GetUVCountWrapper := func(b *testing.B, iter int) {
			count := uvm.GetUVCount()
			b.StopTimer()
			mu.Lock()
			countLog[int64(count)] = true
			mu.Unlock()
			b.StartTimer()
		}
		avgTime := SimulateHighConcurrency(pvm.PVStatisticMock(), uvm.UVStatisticMock(), 1000, 5,
			5000, 10000, GetUVCountWrapper)(b)
		mu.Lock()
		defer mu.Unlock()
		for k := range countLog {
			b.Logf("UV Count: %d", k)
		}
		b.Logf("Average time: %v", avgTime)
	})
}
