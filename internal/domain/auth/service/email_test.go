package service

import (
	"context"
	"errors"
	"testing"
)

// Mock implementation of mail.EmailSender
type mockEmailSender struct {
	sendMailFunc func(ctx context.Context, recipient string, subject string, body string) error
}

func (m *mockEmailSender) SendMail(ctx context.Context, recipient string, subject string, body string) error {
	if m.sendMailFunc != nil {
		return m.sendMailFunc(ctx, recipient, subject, body)
	}
	return nil
}

func TestEmailService_SendVerificationCode(t *testing.T) {
	tests := []struct {
		name          string
		email         string
		code          string
		bodyTemplate  string
		title         string
		mockBehavior  func(recipient string, subject string, body string) error
		expectedError error
	}{
		{
			name:         "successful email send",
			email:        "test@example.com",
			code:         "123456",
			bodyTemplate: "Your code is: %s",
			title:        "Verification Code",
			mockBehavior: func(recipient string, subject string, body string) error {
				if recipient != "test@example.com" {
					return errors.New("recipient mismatch")
				}
				if subject != "Verification Code" {
					return errors.New("subject mismatch")
				}
				if body != "Your code is: 123456" {
					return errors.New("body content mismatch")
				}
				return nil
			},
			expectedError: nil,
		},
		{
			name:         "error from sender",
			email:        "error@example.com",
			code:         "654321",
			bodyTemplate: "Code: %s",
			title:        "Test Email",
			mockBehavior: func(recipient string, subject string, body string) error {
				return errors.New("failed to send email")
			},
			expectedError: errors.New("failed to send email"),
		},
		{
			name:         "empty email address",
			email:        "",
			code:         "000000",
			bodyTemplate: "Code: %s",
			title:        "Test",
			mockBehavior: func(recipient string, subject string, body string) error {
				return nil
			},
			expectedError: nil,
		},
		{
			name:         "empty code",
			email:        "nocode@example.com",
			code:         "",
			bodyTemplate: "Code: %s",
			title:        "Verification",
			mockBehavior: func(recipient string, subject string, body string) error {
				return nil
			},
			expectedError: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockSender := &mockEmailSender{
				sendMailFunc: func(ctx context.Context, recipient string, subject string, body string) error {
					return test.mockBehavior(recipient, subject, body)
				},
			}

			emailService := NewEmailService(test.bodyTemplate, test.title, mockSender)
			err := emailService.SendVerificationCode(context.Background(), test.email, test.code)

			if test.expectedError != nil {
				if err == nil || err.Error() != test.expectedError.Error() {
					t.Errorf("expected error: %v, got: %v", test.expectedError, err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}
