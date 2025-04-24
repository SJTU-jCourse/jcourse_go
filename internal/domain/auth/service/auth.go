package service

import (
	"context"
	"errors"

	"jcourse_go/internal/domain/auth/model"
	"jcourse_go/internal/domain/auth/repository"
	"jcourse_go/pkg/password_hash"
	"jcourse_go/pkg/validator"
)

type AuthService interface {
	Login(ctx context.Context, email string, password string) (*model.UserDomain, error)
	Register(ctx context.Context, email string, password string, code string) (*model.UserDomain, error)
	ResetPassword(ctx context.Context, email string, password string, code string) error
	SendVerificationCode(ctx context.Context, email string) error
}

var (
	ErrorInvalidEmail    = errors.New("invalid email")
	ErrorInvalidPassword = errors.New("invalid email or password")
	ErrorSuspendedUser   = errors.New("suspended user")
	ErrorExistingEmail   = errors.New("existing email")
)

type authService struct {
	verificationCodeService VerificationCodeService
	emailValidator          validator.EmailValidator
	userRepo                repository.UserRepository
	passwordHasher          password_hash.Hasher
}

func (s *authService) SendVerificationCode(ctx context.Context, email string) error {
	if !s.emailValidator.Validate(email) {
		return ErrorInvalidEmail
	}

	err := s.verificationCodeService.SendCode(ctx, email)
	if err != nil {
		return err
	}
	return nil
}

func (s *authService) Login(ctx context.Context, email string, password string) (*model.UserDomain, error) {
	user, err := s.userRepo.FindUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if !user.ValidatePassword(password) {
		return nil, ErrorInvalidPassword
	}

	if !user.CanLogin() {
		return nil, ErrorSuspendedUser
	}

	return user, nil
}

func (s *authService) Register(ctx context.Context, email string, password string, code string) (*model.UserDomain, error) {
	if !s.emailValidator.Validate(email) {
		return nil, ErrorInvalidEmail
	}

	err := s.verificationCodeService.VerifyCode(ctx, email, code)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.FindUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if user != nil {
		return nil, ErrorExistingEmail
	}

	newPassword := model.NewPassword(password, s.passwordHasher)
	user = model.NewUser(email, newPassword)

	err = s.userRepo.Save(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) ResetPassword(ctx context.Context, email string, password string, code string) error {
	if !s.emailValidator.Validate(email) {
		return ErrorInvalidEmail
	}

	err := s.verificationCodeService.VerifyCode(ctx, email, code)
	if err != nil {
		return err
	}

	user, err := s.userRepo.FindUserByEmail(ctx, email)
	if err != nil {
		return err
	}

	newPassword := model.NewPassword(password, s.passwordHasher)
	user.SetPassword(newPassword)

	err = s.userRepo.Save(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func NewAuthService() AuthService {
	return &authService{
		verificationCodeService: nil,
		emailValidator:          nil,
		userRepo:                nil,
	}
}
