package deprecated

// import (
// 	"context"
// 	"fmt"
// 	"jcourse_go/dal"
// 	"jcourse_go/middleware"
// 	"jcourse_go/model/model"
// 	"jcourse_go/model/po"
// 	"jcourse_go/repository"
// 	"jcourse_go/task/biz/ping"
// 	"jcourse_go/task/biz/statistic"

// 	"jcourse_go/util"
// 	"log"
// 	"sync"
// 	"testing"
// 	"time"

// 	"github.com/stretchr/testify/assert"

// 	"github.com/hibiken/asynq"
// )

// func InitTaskClient(t *testing.T) *asynq.Client {
// 	host := util.GetRedisHost()
// 	port := util.GetRedisPort()
// 	redisOpt := asynq.RedisClientOpt{Addr: fmt.Sprintf("%s:%s", host, port), Password: util.GetRedisPassword()}
// 	client := asynq.NewClient(redisOpt)
// 	return client
// }
// func InitTestDB(t *testing.T) {
// 	dal.InitTestMemDBClient()
// 	db := dal.GetDBClient()
// 	err := repository.Migrate(db)
// 	if err != nil {
// 		t.Errorf("Migrate() error = %v", err)
// 		panic(err)
// 	}
// }
// func InitTaskServer(t *testing.T, ctx context.Context, wg *sync.WaitGroup) error {
// 	wg.Add(2)
// 	StartServer(ctx, wg)
// 	StartScheduler(ctx, wg)

// 	if server == nil {
// 		t.Errorf("server is nil")
// 		return fmt.Errorf("server is nil")
// 	}
// 	if scheduler == nil {
// 		t.Errorf("scheduler is nil")
// 		return fmt.Errorf("scheduler is nil")
// 	}
// 	t.Logf("Successfully started server and scheduler")
// 	return nil
// }
// func Test_Asynq(t *testing.T) {
// 	t.Run("TestRunQueue", func(t *testing.T) {
// 		ctx := context.Background()
// 		var wg sync.WaitGroup
// 		err := InitTaskServer(t, ctx, &wg)
// 		if err != nil {
// 			t.Errorf("InitTaskServer() error = %v", err)
// 			return
// 		}
// 		pingTask := asynq.NewTask(ping.TypePing, nil)
// 		client := InitTaskClient(t)
// 		defer func(client *asynq.Client) {
// 			err := client.Close()
// 			if err != nil {
// 				t.Errorf("Close() error = %v", err)
// 			}
// 		}(client)
// 		info, err := client.Enqueue(pingTask)
// 		if err != nil {
// 			t.Errorf("Enqueue() error = %v", err)
// 		}

// 		log.Printf("enqueued task: id=%s queue=%s", info.ID, info.Queue)
// 	})
// 	t.Run("TestRefreshPVDupJudge", func(t *testing.T) {
// 		ctx := context.Background()
// 		InitTestDB(t)
// 		db := dal.GetDBClient()
// 		duration := 30 * time.Minute
// 		db.Model(&po.SettingPO{}).Create(&po.SettingPO{Key: middleware.DuplicateJudgeDurationKey,
// 			Value: duration.String(), Type: model.SettingTypeString})
// 		time.Sleep(1 * time.Second) // wait for db to update
// 		var wg sync.WaitGroup
// 		err := InitTaskServer(t, ctx, &wg)
// 		if err != nil {
// 			t.Errorf("InitTaskServer() error = %v", err)
// 			return
// 		}
// 		client := InitTaskClient(t)
// 		defer func(client *asynq.Client) {
// 			err := client.Close()
// 			if err != nil {
// 				t.Errorf("Close() error = %v", err)
// 			}
// 		}(client)
// 		info, err := client.Enqueue(asynq.NewTask(statistic.TypeRefreshPVDupJudge, nil))
// 		if err != nil {
// 			t.Errorf("Enqueue() error = %v", err)
// 		}
// 		log.Printf("enqueued task: id=%s queue=%s", info.ID, info.Queue)
// 		// wait for task to finish
// 		time.Sleep(1 * time.Second)
// 		// check if duration is updated
// 		assert.NotEqual(t, middleware.DefaultDuplicateJudgeDuration, middleware.LastQuerySiteSettingDuration)
// 	})
// }
