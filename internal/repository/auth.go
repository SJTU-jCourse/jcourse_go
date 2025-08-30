package repository

import (
	"context"
	"fmt"
	"time"

	"jcourse_go/internal/dal"
)

func makeSendVerifyCodeKey(email string) string {
	return fmt.Sprintf("send_verify_code:%s", email)
}

func makeVerifyCodeKey(email string) string {
	return fmt.Sprintf("auth_login_code:%s", email)
}

func GetSendVerifyCodeHistory(ctx context.Context, email string) bool {
	cli := dal.GetRedisClient()
	val, err := cli.Get(ctx, makeSendVerifyCodeKey(email)).Result()
	if err != nil {
		return false
	}
	if len(val) == 0 {
		return false
	}
	return true
}

func StoreSendVerifyCodeHistory(ctx context.Context, email string) error {
	cli := dal.GetRedisClient()
	_, err := cli.SetEx(ctx, makeSendVerifyCodeKey(email), 1, time.Minute).Result()
	return err
}

func GetVerifyCode(ctx context.Context, email string) (string, error) {
	cli := dal.GetRedisClient()
	code, err := cli.Get(ctx, makeVerifyCodeKey(email)).Result()
	return code, err
}

func StoreVerifyCode(ctx context.Context, email, code string) error {
	cli := dal.GetRedisClient()
	_, err := cli.SetEx(ctx, makeVerifyCodeKey(email), code, time.Minute*5).Result()
	return err
}

func ClearVerifyCodeHistory(ctx context.Context, email string) error {
	cli := dal.GetRedisClient()
	_, err := cli.Del(ctx, makeSendVerifyCodeKey(email), makeVerifyCodeKey(email)).Result()
	return err
}
