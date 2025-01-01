package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"jcourse_go/dal"
	"jcourse_go/task"
	"jcourse_go/task/base"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func Init() {
	_ = godotenv.Load()
	dal.InitRedisClient()
	dal.InitDBClient()

	task.InitTaskManager(base.RedisConfig{
		DSN:      dal.GetRedisDSN(),
		Password: dal.GetRedisPassWord(),
	})
}

func main() {
	// 1. Initialize all components
	Init()

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
	registerRouter(r)
	_ = r.Run()
}
