package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"jcourse_go/internal/dal"
	"jcourse_go/internal/repository"
	"jcourse_go/internal/router"
	"jcourse_go/internal/service"
	"jcourse_go/pkg/util"
)

func Init() {
	_ = godotenv.Load()
	dal.InitRedisClient()
	dal.InitDBClient()
	repository.SetDefault(dal.GetDBClient())

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

	// 3. Start serving
	r := gin.Default()
	router.RegisterRouter(r)
	_ = r.Run()
}
