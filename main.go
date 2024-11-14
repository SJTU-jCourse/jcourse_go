package main

import (
	"jcourse_go/dal"
	"jcourse_go/task"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func Init() {
	_ = godotenv.Load()
	dal.InitRedisClient()
	dal.InitDBClient()
	task.InitStatistic(dal.GetDBClient())
}

func main() {
	Init()
	r := gin.Default()
	registerRouter(r)
	_ = r.Run()
}
