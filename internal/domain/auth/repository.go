package auth

import "context"

type VerificationCodeRepository interface {
	Get(ctx context.Context, email string) (*VerificationCode, error)
	Save(ctx context.Context, code *VerificationCode) error
	Delete(ctx context.Context, email string) error
}
