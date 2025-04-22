package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"jcourse_go/internal/domain/auth/model"
)

type mockRepo struct {
	mock.Mock
}

func (m *mockRepo) StoreVerifyCode(ctx context.Context, email, code string) error {
	args := m.Called(ctx, email, code)
	return args.Error(0)
}

func (m *mockRepo) GetVerifyCode(ctx context.Context, email string) (model.VerificationCode, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(model.VerificationCode), args.Error(1)
}

func (m *mockRepo) ClearHistory(ctx context.Context, email string) error {
	args := m.Called(ctx, email)
	return args.Error(0)
}

type mockEmailService struct {
	mock.Mock
}

func (m *mockEmailService) SendVerificationCode(ctx context.Context, email string, code string) error {
	args := m.Called(ctx, email, code)
	return args.Error(0)
}

func TestSendCode(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func(repo *mockRepo, emailService *mockEmailService)
		email          string
		expectedErrMsg string
	}{
		{
			name: "Error fetching verification code",
			mockSetup: func(repo *mockRepo, emailService *mockEmailService) {
				repo.On("GetVerifyCode", mock.Anything, "test@example.com").Return(model.VerificationCode{}, errors.New("fetch error"))
			},
			email:          "test@example.com",
			expectedErrMsg: "fetch error",
		},
		{
			name: "Rate limit error",
			mockSetup: func(repo *mockRepo, emailService *mockEmailService) {
				repo.On("GetVerifyCode", mock.Anything, "test@example.com").Return(model.VerificationCode{CreatedAt: time.Now()}, nil)
			},
			email:          "test@example.com",
			expectedErrMsg: ErrorRateLimit.Error(),
		},
		{
			name: "Error storing verification code",
			mockSetup: func(repo *mockRepo, emailService *mockEmailService) {
				repo.On("GetVerifyCode", mock.Anything, "test@example.com").Return(model.VerificationCode{}, nil)
				repo.On("StoreVerifyCode", mock.Anything, "test@example.com", mock.Anything).Return(errors.New("storage error"))
			},
			email:          "test@example.com",
			expectedErrMsg: "storage error",
		},
		{
			name: "Error sending verification email",
			mockSetup: func(repo *mockRepo, emailService *mockEmailService) {
				repo.On("GetVerifyCode", mock.Anything, "test@example.com").Return(model.VerificationCode{}, nil)
				repo.On("StoreVerifyCode", mock.Anything, "test@example.com", mock.Anything).Return(nil)
				emailService.On("SendVerificationCode", mock.Anything, "test@example.com", mock.Anything).Return(errors.New("email error"))
			},
			email:          "test@example.com",
			expectedErrMsg: "email error",
		},
		{
			name: "Success",
			mockSetup: func(repo *mockRepo, emailService *mockEmailService) {
				repo.On("GetVerifyCode", mock.Anything, "test@example.com").Return(model.VerificationCode{}, nil)
				repo.On("StoreVerifyCode", mock.Anything, "test@example.com", mock.Anything).Return(nil)
				emailService.On("SendVerificationCode", mock.Anything, "test@example.com", mock.Anything).Return(nil)
			},
			email:          "test@example.com",
			expectedErrMsg: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockRepo{}
			emailService := &mockEmailService{}
			tt.mockSetup(repo, emailService)

			service := &verificationCodeService{
				repo:         repo,
				emailService: emailService,
				rateLimit:    1 * time.Minute,
				codeLength:   6,
				codeCharSet:  "1234567890",
			}

			err := service.SendCode(context.Background(), tt.email)
			if tt.expectedErrMsg == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tt.expectedErrMsg)
			}

			repo.AssertExpectations(t)
			emailService.AssertExpectations(t)
		})
	}
}

func TestVerifyCode(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func(repo *mockRepo)
		email          string
		code           string
		expectedErrMsg string
	}{
		{
			name: "Error fetching verification code",
			mockSetup: func(repo *mockRepo) {
				repo.On("GetVerifyCode", mock.Anything, "test@example.com").Return(model.VerificationCode{}, errors.New("fetch error"))
			},
			email:          "test@example.com",
			code:           "123456",
			expectedErrMsg: "fetch error",
		},
		{
			name: "Error clearing history",
			mockSetup: func(repo *mockRepo) {
				repo.On("GetVerifyCode", mock.Anything, "test@example.com").Return(model.VerificationCode{Code: "123456"}, nil)
				repo.On("ClearHistory", mock.Anything, "test@example.com").Return(errors.New("clear error"))
			},
			email:          "test@example.com",
			code:           "123456",
			expectedErrMsg: "clear error",
		},
		{
			name: "Code mismatch",
			mockSetup: func(repo *mockRepo) {
				repo.On("GetVerifyCode", mock.Anything, "test@example.com").Return(model.VerificationCode{Code: "654321"}, nil)
				repo.On("ClearHistory", mock.Anything, "test@example.com").Return(nil)
			},
			email:          "test@example.com",
			code:           "123456",
			expectedErrMsg: ErrorCodeNotMatch.Error(),
		},
		{
			name: "Success",
			mockSetup: func(repo *mockRepo) {
				repo.On("GetVerifyCode", mock.Anything, "test@example.com").Return(model.VerificationCode{Code: "123456", CreatedAt: time.Now(), TTL: time.Minute}, nil)
				repo.On("ClearHistory", mock.Anything, "test@example.com").Return(nil)
			},
			email:          "test@example.com",
			code:           "123456",
			expectedErrMsg: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockRepo{}
			tt.mockSetup(repo)

			service := &verificationCodeService{
				repo: repo,
			}

			err := service.VerifyCode(context.Background(), tt.email, tt.code)
			if tt.expectedErrMsg == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tt.expectedErrMsg)
			}

			repo.AssertExpectations(t)
		})
	}
}

func TestGenerateCode(t *testing.T) {
	service := &verificationCodeService{
		codeLength:  6,
		codeCharSet: "1234567890",
	}

	code := service.GenerateCode()
	require.Len(t, code, service.codeLength)

	for _, char := range code {
		require.Contains(t, service.codeCharSet, string(char))
	}
}
