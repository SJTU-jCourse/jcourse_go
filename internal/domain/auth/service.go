package auth

import "context"

type AuthService interface {
	Login(ctx context.Context, cmd LoginCommand) error
	Register(ctx context.Context, cmd RegisterUserCommand) error
	ResetPassword(ctx context.Context, cmd ResetPasswordCommand) error
	SendVerificationCode(ctx context.Context, cmd SendVerificationCodeCommand) error
}
