package repository

import (
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"

	"jcourse_go/internal/dal"
)

func TestClearVerifyCodeHistory(t *testing.T) {
	t.Run("has sent", func(t *testing.T) {
		TestStoreSendVerifyCodeHistory(t)
		TestStoreVerifyCode(t)
		ctx := context.TODO()
		mock := dal.InitMockRedisClient()
		email := "test@example.com"
		mock.ExpectDel(makeSendVerifyCodeKey(email), makeVerifyCodeKey(email)).SetVal(2)
		err := ClearVerifyCodeHistory(ctx, email)
		assert.Nil(t, err)
	})
	t.Run("no sent", func(t *testing.T) {
		ctx := context.TODO()
		mock := dal.InitMockRedisClient()
		email := "test@example.com"
		mock.ExpectDel(makeSendVerifyCodeKey(email), makeVerifyCodeKey(email)).SetVal(0)
		err := ClearVerifyCodeHistory(ctx, email)
		assert.Nil(t, err)
	})
}

func TestGetSendVerifyCodeHistory(t *testing.T) {
	t.Run("has sent", func(t *testing.T) {
		TestStoreSendVerifyCodeHistory(t)
		ctx := context.TODO()
		mock := dal.InitMockRedisClient()
		email := "test@example.com"
		mock.ExpectGet(makeSendVerifyCodeKey(email)).SetVal("1")
		sent := GetSendVerifyCodeHistory(ctx, email)
		assert.True(t, sent)
	})
	t.Run("no sent", func(t *testing.T) {
		ctx := context.TODO()
		mock := dal.InitMockRedisClient()
		email := "test@example.com"
		mock.ExpectGet(makeSendVerifyCodeKey(email)).RedisNil()
		sent := GetSendVerifyCodeHistory(ctx, email)
		assert.False(t, sent)
	})
}

func TestGetVerifyCode(t *testing.T) {
	t.Run("has sent", func(t *testing.T) {
		TestStoreVerifyCode(t)
		ctx := context.TODO()
		mock := dal.InitMockRedisClient()
		email := "test@example.com"
		mock.ExpectGet(makeVerifyCodeKey(email)).SetVal("123456")
		code, err := GetVerifyCode(ctx, email)
		assert.Nil(t, err)
		assert.Equal(t, "123456", code)
	})
	t.Run("no sent", func(t *testing.T) {
		ctx := context.TODO()
		mock := dal.InitMockRedisClient()
		email := "test@example.com"
		mock.ExpectGet(makeVerifyCodeKey(email)).RedisNil()
		_, err := GetVerifyCode(ctx, email)
		assert.Equal(t, redis.Nil, err)
	})
}

func TestStoreSendVerifyCodeHistory(t *testing.T) {
	ctx := context.TODO()
	mock := dal.InitMockRedisClient()
	email := "test@example.com"
	mock.ExpectSetEx(makeSendVerifyCodeKey(email), 1, time.Minute).SetVal("Ok")
	err := StoreSendVerifyCodeHistory(ctx, email)
	assert.Nil(t, err)
}

func TestStoreVerifyCode(t *testing.T) {
	ctx := context.TODO()
	mock := dal.InitMockRedisClient()
	email := "test@example.com"
	code := "123456"

	mock.ExpectSetEx(makeVerifyCodeKey(email), code, time.Minute*5).SetVal("Ok")
	err := StoreVerifyCode(ctx, email, code)
	assert.Nil(t, err)
}
