package command

import (
	"context"

	"jcourse_go/internal/domain/auth"
)

type AuthService interface {
	Login(ctx context.Context, cmd auth.LoginCommand) error
	Register(ctx context.Context, cmd auth.RegisterUserCommand) error
	ResetPassword(ctx context.Context, cmd auth.ResetPasswordCommand) error
	SendVerificationCode(ctx context.Context, cmd auth.SendVerificationCodeCommand) error
}
