package main

import (
	"jcourse_go/dal"
	"jcourse_go/rpc"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func Init() {
	_ = godotenv.Load()
	dal.InitRedisClient()
	dal.InitDBClient()
	rpc.InitOpenAIClient()
}

func main() {
	Init()
	r := gin.Default()
	registerRouter(r)
	_ = r.Run()
}
