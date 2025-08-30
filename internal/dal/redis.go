package dal

import (
	"fmt"

	"github.com/go-redis/redismock/v9"
	"github.com/redis/go-redis/v9"

	"jcourse_go/pkg/util"
)

var rdb *redis.Client

func GetRedisDSN() string {
	host := util.GetRedisHost()
	port := util.GetRedisPort()
	return fmt.Sprintf("%s:%s", host, port)
}
func GetRedisPassWord() string {
	return util.GetRedisPassword()
}
func InitRedisClient() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     GetRedisDSN(),
		Password: GetRedisPassWord(),
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
