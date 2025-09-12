package dal

import (
	"github.com/redis/go-redis/v9"

	"jcourse_go/internal/config"
)

func NewRedisClient(c config.RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     c.Addr,
		Username: c.Username,
		Password: c.Password,
		DB:       c.DB,
	})
	return client, nil
}
