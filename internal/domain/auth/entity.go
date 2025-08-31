package auth

import (
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

func NewVerificationCode(email, code string, now time.Time) *VerificationCode {
	return &VerificationCode{
		Email:     email,
		Code:      code,
		ExpiresAt: now.Add(CodeTTL),
		CreatedAt: now,
		UpdatedAt: now,
	}
}

type AuthUser struct {
	ID       shared.IDType
	Email    string
	Password string
	Role     shared.UserRole

	LastSeenAt time.Time

	SuspendedAt   *time.Time
	SuspendedBy   *int64
	SuspendedTill *time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}
