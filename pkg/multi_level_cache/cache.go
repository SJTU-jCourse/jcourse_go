package multi_level_cache

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

// Cache 接口定义了缓存的基本操作
type Cache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, expiration time.Duration) error
}

// RedisCache 实现了 Cache 接口
type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(rdb *redis.Client) *RedisCache {
	return &RedisCache{client: rdb}
}

func (r *RedisCache) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *RedisCache) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}

// FetchFunc 回源函数类型
type FetchFunc = func(ctx context.Context, key string) (string, error)

// MultiLevelCache 结构体，包含 Redis 和回源函数
type MultiLevelCache struct {
	redisCache Cache
	fetchFunc  FetchFunc
	expiration time.Duration
}

func NewMultiLevelCache(redisCache Cache, fetchFunc FetchFunc, expiration time.Duration) *MultiLevelCache {
	return &MultiLevelCache{
		redisCache: redisCache,
		fetchFunc:  fetchFunc,
		expiration: expiration,
	}
}

func (m *MultiLevelCache) Get(ctx context.Context, key string) (string, error) {
	// 从 Redis 获取缓存
	value, err := m.redisCache.Get(ctx, key)
	if err == nil {
		return value, nil
	}

	// 如果 Redis 缓存未命中，调用回源函数获取数据
	value, err = m.fetchFunc(ctx, key)
	if err != nil {
		return "", err
	}

	// 将数据缓存到 Redis
	err = m.redisCache.Set(ctx, key, value, m.expiration)
	if err != nil {
		log.Printf("Failed to set redis cache: %v", err)
	}

	return value, nil
}
