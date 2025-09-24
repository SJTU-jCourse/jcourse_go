package auth

import (
	"strings"
	"time"

	"jcourse_go/internal/domain/shared"
)

const (
	CodeTTL = 5 * time.Minute
)

type VerificationCode struct {
	Email     string
	Code      string
	ExpiresAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (v VerificationCode) IsExpired(now time.Time) bool {
	return v.ExpiresAt.Before(now)
}

func (v VerificationCode) IsValid(code string, now time.Time) bool {
	return v.Code == code && !v.IsExpired(now)
}

func (v VerificationCode) EmailTitle() string {
	return VerificationEmailTitle
}

func (v VerificationCode) EmailBody() string {
	return strings.ReplaceAll(VerificationEmailBody, "{code}", v.Code)
}

func RandomCode() string {
	return ""
}

func NewVerificationCode(email string, now time.Time) *VerificationCode {
	return &VerificationCode{
		Email:     email,
		Code:      RandomCode(),
		ExpiresAt: now.Add(CodeTTL),
		CreatedAt: now,
		UpdatedAt: now,
	}
}

type AuthUser struct {
	ID       shared.IDType
	Email    string
	Password string
	UserRole shared.UserRole

	LastSeenAt time.Time

	SuspendedAt   *time.Time
	SuspendedBy   *int64
	SuspendedTill *time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (a *AuthUser) ValidatePassword(password string) bool {
	return a.Password == password
}

func (a *AuthUser) ResetPassword(password string) error {
	a.Password = password
	return nil
}

func NewAuthUser(cmd RegisterUserCommand, now time.Time) AuthUser {
	return AuthUser{
		Email:    cmd.Email,
		Password: cmd.Password,

		CreatedAt:  now,
		UpdatedAt:  now,
		LastSeenAt: now,
	}
}

const (
	SessionTTL = time.Hour * 24 * 30
)

type Session struct {
	SessionID string
	UserID    shared.IDType
	CreatedAt time.Time
	ExpiresAt time.Time
}

func NewSession(userID shared.IDType, now time.Time) Session {
	return Session{
		UserID:    userID,
		CreatedAt: now,
		ExpiresAt: now.Add(SessionTTL),
	}
}
