package infra

import (
	"fmt"

	"github.com/go-redis/redismock/v9"
	"github.com/redis/go-redis/v9"

	"jcourse_go/config"
)

var rdb *redis.Client

func GetRedisDSN(host string, port int) string {
	return fmt.Sprintf("%s:%d", host, port)
}

func InitRedisClient(conf *config.Redis) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     GetRedisDSN(conf.Host, conf.Port),
		Password: conf.Password,
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
