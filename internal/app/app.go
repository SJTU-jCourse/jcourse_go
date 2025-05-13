package app

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"

	"jcourse_go/config"
	"jcourse_go/internal/infra"
	"jcourse_go/internal/infra/query"
	"jcourse_go/pkg/task"
	"jcourse_go/pkg/task/base"
	"jcourse_go/service"
	"jcourse_go/util"
)

func Init(conf *config.Config) {
	infra.InitRedisClient(&conf.Redis)
	infra.InitDBClient(&conf.DB)
	query.SetDefault(infra.GetDBClient())

	task.InitTaskManager(base.RedisConfig{
		DSN:      infra.GetRedisDSN(conf.Redis.Host, conf.Redis.Port),
		Password: conf.Redis.Password,
	})

	if err := util.InitSegWord(); err != nil {
		panic(err)
	}

	err := service.InitLLM(&conf.LLM)
	if err != nil {
		panic(err)
	}
}

func Run(conf *config.Config) {
	// 1. Initialize all components
	Init(conf)

	// 2. Listen for signals to gracefully shut down
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-c
		log.Printf("[main] Caught signal: %v. Starting graceful shutdown...", sig)
		err := task.GlobalTaskManager.Shutdown()
		if err != nil {
			log.Printf("[main] Error while shutting down TaskManager: %v\n", err)
		}
		log.Println("[main] Graceful shutdown complete. Exiting.")
		os.Exit(0)
	}()

	// 3. Start serving
	r := gin.Default()
	registerRouter(conf, r)
	_ = r.Run(fmt.Sprintf(":%d", conf.Server.Port))
}
