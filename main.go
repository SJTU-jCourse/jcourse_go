package main

import (
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
	Init()
	r := gin.Default()
	registerRouter(r)
	_ = r.Run()
}
