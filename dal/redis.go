package dal

import (
	"fmt"
	"os"

	"github.com/go-redis/redismock/v9"
	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client

func GetRedisDSN() string {
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	return fmt.Sprintf("%s:%s", host, port)
}

func InitRedisClient() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     GetRedisDSN(),
		Password: "",
		DB:       0,
	})
}

func GetRedisClient() *redis.Client {
	return rdb
}

func InitMockRedisClient() redismock.ClientMock {
	var mock redismock.ClientMock
	rdb, mock = redismock.NewClientMock()
	return mock
}
