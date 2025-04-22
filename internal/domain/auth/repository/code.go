package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"jcourse_go/internal/domain/auth/model"
)

type VerificationCodeRepository interface {
	StoreVerifyCode(ctx context.Context, email, code string) error
	GetVerifyCode(ctx context.Context, email string) (model.VerificationCode, error)
	ClearHistory(ctx context.Context, email string) error
}

type redisVerificationCodeRepository struct {
	client   *redis.Client
	duration time.Duration
}

func (r *redisVerificationCodeRepository) makeVerifyCodeKey(email string) string {
	return fmt.Sprintf("auth_login_code:%s", email)
}

func (r *redisVerificationCodeRepository) StoreVerifyCode(ctx context.Context, email, code string) error {
	key := r.makeVerifyCodeKey(email)
	_, err := r.client.Set(ctx, key, code, r.duration).Result()
	return err
}

func (r *redisVerificationCodeRepository) GetVerifyCode(ctx context.Context, email string) (model.VerificationCode, error) {
	key := r.makeVerifyCodeKey(email)

	pipe := r.client.Pipeline()

	codeCmd := pipe.Get(ctx, key)
	ttlCmd := pipe.TTL(ctx, key)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return model.VerificationCode{}, err
	}

	code := codeCmd.Val()
	ttl := ttlCmd.Val()
	createdAt := time.Now().Add(-ttl)
	return model.VerificationCode{
		Code:      code,
		Email:     email,
		CreatedAt: createdAt,
		TTL:       ttl,
	}, err
}

func (r *redisVerificationCodeRepository) ClearHistory(ctx context.Context, email string) error {
	key := r.makeVerifyCodeKey(email)
	_, err := r.client.Del(ctx, key).Result()
	return err
}

func NewRedisRepository(client *redis.Client, duration time.Duration) VerificationCodeRepository {
	return &redisVerificationCodeRepository{
		client:   client,
		duration: duration,
	}
}
