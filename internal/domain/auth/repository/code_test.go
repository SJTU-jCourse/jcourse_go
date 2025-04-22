package repository

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"

	"jcourse_go/internal/domain/auth/model"
)

func TestRedisVerificationCodeRepository_StoreVerifyCode(t *testing.T) {
	mockClient, mock := redismock.NewClientMock()
	repo := NewRedisRepository(mockClient, 5*time.Minute)

	tests := []struct {
		name     string
		email    string
		code     string
		mockFunc func()
		wantErr  bool
	}{
		{
			name:  "Success",
			email: "test@example.com",
			code:  "12345",
			mockFunc: func() {
				mock.ExpectSet("auth_login_code:test@example.com", "12345", 5*time.Minute).SetVal("OK")
			},
			wantErr: false,
		},
		{
			name:  "Redis error",
			email: "test@example.com",
			code:  "12345",
			mockFunc: func() {
				mock.ExpectSet("auth_login_code:test@example.com", "12345", 5*time.Minute).SetErr(errors.New("redis error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()
			err := repo.StoreVerifyCode(context.Background(), tt.email, tt.code)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreVerifyCode() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unmet redis expectations: %s", err)
			}
		})
	}
}

func TestRedisVerificationCodeRepository_GetVerifyCode(t *testing.T) {
	mockClient, mock := redismock.NewClientMock()
	repo := NewRedisRepository(mockClient, 5*time.Minute)

	tests := []struct {
		name     string
		email    string
		mockFunc func()
		want     model.VerificationCode
		wantErr  bool
	}{
		{
			name:  "Success",
			email: "test@example.com",
			mockFunc: func() {
				mock.ExpectGet("auth_login_code:test@example.com").SetVal("12345")
				mock.ExpectTTL("auth_login_code:test@example.com").SetVal(3 * time.Minute)
			},
			want: model.VerificationCode{
				Code:  "12345",
				Email: "test@example.com",
				TTL:   3 * time.Minute,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()
			got, err := repo.GetVerifyCode(context.Background(), tt.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetVerifyCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && (got.Code != tt.want.Code || got.Email != tt.want.Email || got.TTL != tt.want.TTL) {
				t.Errorf("GetVerifyCode() = %v, want %v", got, tt.want)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unmet redis expectations: %s", err)
			}
		})
	}
}

func TestRedisVerificationCodeRepository_ClearHistory(t *testing.T) {
	mockClient, mock := redismock.NewClientMock()
	repo := NewRedisRepository(mockClient, 5*time.Minute)

	tests := []struct {
		name     string
		email    string
		mockFunc func()
		wantErr  bool
	}{
		{
			name:  "Success",
			email: "test@example.com",
			mockFunc: func() {
				mock.ExpectDel("auth_login_code:test@example.com").SetVal(1)
			},
			wantErr: false,
		},
		{
			name:  "Redis error",
			email: "test@example.com",
			mockFunc: func() {
				mock.ExpectDel("auth_login_code:test@example.com").SetErr(errors.New("redis error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()
			err := repo.ClearHistory(context.Background(), tt.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("ClearHistory() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unmet redis expectations: %s", err)
			}
		})
	}
}
