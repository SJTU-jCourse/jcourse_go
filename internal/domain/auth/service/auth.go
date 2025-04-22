package service

import (
	"context"

	"jcourse_go/internal/domain/auth/model"
	"jcourse_go/internal/domain/auth/repository"
)

type AuthService interface {
	Login(ctx context.Context, email string, password string) (*model.UserDomain, error)
	Register(ctx context.Context, email string, password string) (*model.UserDomain, error)
	ResetPassword(ctx context.Context, email string, password string) error
}

type authService struct {
	userRepo repository.UserRepository
}
