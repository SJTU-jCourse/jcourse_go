package repository

import "jcourse_go/internal/domain/auth/model"

type UserRepository interface {
	FindUserByEmail(email string) (*model.UserDomain, error)
	Save(user *model.UserDomain) error
}
