package repository

import (
	"context"

	"jcourse_go/internal/domain/auth/model"
	"jcourse_go/internal/infra/query"
	"jcourse_go/model/po"
	"jcourse_go/model/types"
)

type UserRepository interface {
	FindUserByEmail(ctx context.Context, email string) (*model.UserDomain, error)
	Save(ctx context.Context, user *model.UserDomain) error
}

type userRepository struct {
	q *query.Query
}

func (r *userRepository) FindUserByEmail(ctx context.Context, email string) (*model.UserDomain, error) {
	u := r.q.UserPO
	userPO, err := u.WithContext(ctx).Where(u.Email.Eq(email)).First()
	if err != nil {
		return nil, err
	}
	return r.fromPO(userPO), nil
}

func (r *userRepository) Save(ctx context.Context, user *model.UserDomain) error {
	u := r.q.UserPO
	userPO := r.toPO(user)
	err := u.WithContext(ctx).Save(userPO)
	if err != nil {
		return err
	}
	return nil
}

func (r *userRepository) toPO(user *model.UserDomain) *po.UserPO {
	if user == nil {
		return nil
	}
	return &po.UserPO{
		ID:            int64(user.ID),
		CreatedAt:     user.CreatedAt,
		Username:      user.Username,
		Email:         user.Email,
		Password:      user.Password.String(),
		UserRole:      user.Role.String(),
		LastSeenAt:    user.LastSeenAt,
		SuspendedAt:   user.SuspendedAt,
		SuspendedTill: user.SuspendedTill,
	}
}

func (r *userRepository) fromPO(user *po.UserPO) *model.UserDomain {
	if user == nil {
		return nil
	}
	return &model.UserDomain{
		ID:            int(user.ID),
		Username:      user.Username,
		Email:         user.Email,
		Password:      model.Password(user.Password),
		Role:          types.UserRole(user.UserRole),
		CreatedAt:     user.CreatedAt,
		LastSeenAt:    user.LastSeenAt,
		SuspendedAt:   user.SuspendedAt,
		SuspendedTill: user.SuspendedTill,
	}
}

func NewUserRepository(query *query.Query) UserRepository {
	return &userRepository{
		q: query,
	}
}
