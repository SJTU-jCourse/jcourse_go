package user

import (
	"context"

	"jcourse_go/internal/domain/shared"
)

type UserRepository interface {
	Get(ctx context.Context, id shared.IDType) (*User, error)
	Save(ctx context.Context, user *User) error
}

type UserPointRepository interface {
	FindByUserID(ctx context.Context, userID shared.IDType) ([]UserPoint, error)
	Save(ctx context.Context, userPoint *UserPoint) error
	Transaction(ctx context.Context, transaction *UserPointTransaction) error
}
