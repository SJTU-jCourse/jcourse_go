package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type VerifyCodeRepository struct {
	rdb *redis.Client
}

func NewAuthRepository(redis *redis.Client) *VerifyCodeRepository {
	return &VerifyCodeRepository{rdb: redis}
}

func (r *VerifyCodeRepository) makeSendVerifyCodeKey(email string) string {
	return fmt.Sprintf("send_verify_code:%s", email)
}

func (r *VerifyCodeRepository) makeVerifyCodeKey(email string) string {
	return fmt.Sprintf("auth_login_code:%s", email)
}

func (r *VerifyCodeRepository) GetSendVerifyCodeHistory(ctx context.Context, email string) bool {
	val, err := r.rdb.Get(ctx, r.makeSendVerifyCodeKey(email)).Result()
	if err != nil {
		return false
	}
	if len(val) == 0 {
		return false
	}
	return true
}

func (r *VerifyCodeRepository) StoreSendVerifyCodeHistory(ctx context.Context, email string) error {
	_, err := r.rdb.SetEx(ctx, r.makeSendVerifyCodeKey(email), 1, time.Minute).Result()
	return err
}

func (r *VerifyCodeRepository) GetVerifyCode(ctx context.Context, email string) (string, error) {
	code, err := r.rdb.Get(ctx, r.makeVerifyCodeKey(email)).Result()
	return code, err
}

func (r *VerifyCodeRepository) StoreVerifyCode(ctx context.Context, email, code string) error {
	_, err := r.rdb.SetEx(ctx, r.makeVerifyCodeKey(email), code, time.Minute*5).Result()
	return err
}

func (r *VerifyCodeRepository) ClearVerifyCodeHistory(ctx context.Context, email string) error {
	_, err := r.rdb.Del(ctx, r.makeSendVerifyCodeKey(email), r.makeVerifyCodeKey(email)).Result()
	return err
}
