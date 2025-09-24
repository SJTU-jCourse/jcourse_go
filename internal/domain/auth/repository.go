package auth

import (
	"context"

	"jcourse_go/internal/domain/shared"
)

type VerificationCodeRepository interface {
	Get(ctx context.Context, email string) (*VerificationCode, error)
	Save(ctx context.Context, code *VerificationCode) error
	Delete(ctx context.Context, email string) error
}

type SessionRepository interface {
	Get(ctx context.Context, sessionID string) (*Session, error)
	GetByUser(ctx context.Context, userID shared.IDType) (*Session, error)
	Delete(ctx context.Context, sessionID string) error
	DeleteByUser(ctx context.Context, userID shared.IDType) error
	Save(ctx context.Context, session *Session) error
}

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*AuthUser, error)
	Save(ctx context.Context, user *AuthUser) error
}
