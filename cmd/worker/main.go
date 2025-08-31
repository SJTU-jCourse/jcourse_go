package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"

	"jcourse_go/internal/dal"

	"jcourse_go/internal/infrastructure/repository"
	"jcourse_go/internal/interface/task"
	"jcourse_go/internal/interface/task/base"
	"jcourse_go/internal/service"
	"jcourse_go/pkg/util"
)

func Init() {
	_ = godotenv.Load()
	dal.InitRedisClient()
	dal.InitDBClient()
	repository.SetDefault(dal.GetDBClient())

	task.InitTaskManager(base.RedisConfig{
		DSN:      dal.GetRedisDSN(),
		Password: dal.GetRedisPassWord(),
	})

	if err := util.InitSegWord(); err != nil {
		panic(err)
	}

	err := service.InitLLM()
	if err != nil {
		panic(err)
	}
}

func main() {
	// 1. Initialize all components
	Init()

	// 2. Listen for signals to gracefully shut down
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	sig := <-c
	log.Printf("[main] Caught signal: %v. Starting graceful shutdown...", sig)
	err := task.GlobalTaskManager.Shutdown()
	if err != nil {
		log.Printf("[main] Error while shutting down TaskManager: %v\n", err)
	}
	log.Println("[main] Graceful shutdown complete. Exiting.")
	os.Exit(0)

}
