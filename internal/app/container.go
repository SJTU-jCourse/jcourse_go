package app

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"jcourse_go/internal/config"
	dal2 "jcourse_go/internal/infrastructure/dal"
	"jcourse_go/internal/service/auth"
)

type ServiceContainer struct {
	DB    *gorm.DB
	Redis *redis.Client
	Auth  auth.AuthService
}

func NewServiceContainer(c *config.AppConfig) (*ServiceContainer, error) {
	db, err := dal2.NewPostgresSQL(&c.DB)
	if err != nil {
		return nil, err
	}
	rdb, err := dal2.NewRedisClient(&c.Redis)
	if err != nil {
		return nil, err
	}
	return &ServiceContainer{
		DB:    db,
		Redis: rdb,
	}, nil
}
