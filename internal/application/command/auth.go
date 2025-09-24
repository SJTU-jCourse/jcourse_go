package command

import (
	"context"
	"time"

	"jcourse_go/internal/domain/auth"
	"jcourse_go/internal/domain/email"
	"jcourse_go/internal/domain/shared"
	"jcourse_go/pkg/apperror"
)

type AuthService interface {
	Login(ctx context.Context, cmd auth.LoginCommand) error
	Logout(ctx context.Context, reqCtx shared.RequestCtx) error
	Register(ctx context.Context, cmd auth.RegisterUserCommand) error
	ResetPassword(ctx context.Context, reqCtx shared.RequestCtx, cmd auth.ResetPasswordCommand) error
	SendVerificationCode(ctx context.Context, cmd auth.SendVerificationCodeCommand) error
}

type authService struct {
	verificationCodeRepo auth.VerificationCodeRepository
	sessionRepo          auth.SessionRepository
	userRepo             auth.UserRepository
	emailSender          email.EmailSender
}

func (a *authService) Login(ctx context.Context, cmd auth.LoginCommand) error {
	u, err := a.userRepo.FindByEmail(ctx, cmd.Email)
	if err != nil {
		return err
	}
	if u == nil || !u.ValidatePassword(cmd.Password) {
		return apperror.ErrNoPermission
	}
	session := auth.NewSession(u.ID, time.Now())
	err = a.sessionRepo.Save(ctx, &session)
	if err != nil {
		return err
	}
	return nil
}

func (a *authService) Logout(ctx context.Context, reqCtx shared.RequestCtx) error {
	s, err := a.sessionRepo.GetByUser(ctx, reqCtx.User.UserID)
	if err != nil || s == nil {
		return err
	}
	return a.sessionRepo.Delete(ctx, s.SessionID)
}

func (a *authService) Register(ctx context.Context, cmd auth.RegisterUserCommand) error {
	code, err := a.verificationCodeRepo.Get(ctx, cmd.Email)
	if err != nil {
		return err
	}
	if code == nil {
		return apperror.ErrNoPermission
	}

	if err = a.verificationCodeRepo.Delete(ctx, cmd.Email); err != nil {
		return err
	}

	now := time.Now()
	user := auth.NewAuthUser(cmd, now)

	if err = a.userRepo.Save(ctx, &user); err != nil {
		return err
	}

	session := auth.NewSession(user.ID, now)
	if err = a.sessionRepo.Save(ctx, &session); err != nil {
		return err
	}
	return nil
}

func (a *authService) ResetPassword(ctx context.Context, reqCtx shared.RequestCtx, cmd auth.ResetPasswordCommand) error {
	code, err := a.verificationCodeRepo.Get(ctx, cmd.Email)
	if err != nil {
		return err
	}
	if code == nil {
		return apperror.ErrNoPermission
	}

	if err = a.verificationCodeRepo.Delete(ctx, cmd.Email); err != nil {
		return err
	}

	user, err := a.userRepo.FindByEmail(ctx, cmd.Email)
	if err != nil {
		return err
	}
	if err = user.ResetPassword(cmd.Password); err != nil {
		return err
	}
	now := time.Now()
	if err = a.userRepo.Save(ctx, user); err != nil {
		return err
	}

	if err = a.sessionRepo.DeleteByUser(ctx, user.ID); err != nil {
		return err
	}

	session := auth.NewSession(user.ID, now)
	if err = a.sessionRepo.Save(ctx, &session); err != nil {
		return err
	}
	return nil
}

func (a *authService) SendVerificationCode(ctx context.Context, cmd auth.SendVerificationCodeCommand) error {
	now := time.Now()
	verificationCode := auth.NewVerificationCode(cmd.Email, now)
	if err := a.verificationCodeRepo.Save(ctx, verificationCode); err != nil {
		return err
	}

	emailReq := email.Request{
		Title:     verificationCode.EmailTitle(),
		Body:      verificationCode.EmailBody(),
		Recipient: cmd.Email,
	}

	if err := a.emailSender.SendEmail(ctx, emailReq); err != nil {
		return err
	}
	return nil
}

func NewAuthService(
	verificationCodeRepo auth.VerificationCodeRepository,
	sessionRepo auth.SessionRepository,
	userRepo auth.UserRepository,
	emailSender email.EmailSender,
) AuthService {
	return &authService{
		verificationCodeRepo: verificationCodeRepo,
		sessionRepo:          sessionRepo,
		userRepo:             userRepo,
		emailSender:          emailSender,
	}
}
