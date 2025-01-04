package lock

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type IDistributedLockManager interface {
	NewLock(key string, ttl time.Duration) (IDistributedLock, error)
}

type RedisDistributedLockManager struct {
	client *redis.Client
}

type IDistributedLock interface {
	TryLock() (bool, error)
	Unlock() (bool, error)
}

func NewRedisDistributedLockManager(client *redis.Client) *RedisDistributedLockManager {
	return &RedisDistributedLockManager{
		client: client,
	}
}

func (m *RedisDistributedLockManager) NewLock(key string, ttl time.Duration) (IDistributedLock, error) {
	if m == nil || m.client == nil || key == "" {
		return nil, errors.New("nil redis client or invalid parameters")
	}
	return &RedisDistributedLock{
		client: m.client,
		key:    key,
		ttl:    ttl,
		value:  uuid.NewString(),
	}, nil
}

type RedisDistributedLock struct {
	client *redis.Client
	key    string
	ttl    time.Duration
	value  string
}

func (l *RedisDistributedLock) TryLock() (bool, error) {
	ctx := context.Background()
	success, err := l.client.SetNX(ctx, l.key, l.value, l.ttl).Result()
	if err != nil {
		return false, err
	}
	return success, nil
}

func (l *RedisDistributedLock) Unlock() (bool, error) {
	ctx := context.Background()
	script := `
    if redis.call("get", KEYS[1]) == ARGV[1] then
        return redis.call("del", KEYS[1])
    else
        return 0
    end
    `
	result, err := l.client.Eval(ctx, script, []string{l.key}, l.value).Result()
	if err != nil {
		return false, err
	}
	return result.(int64) == 1, nil
}
