package query

import (
	"context"

	"gorm.io/gorm"

	"jcourse_go/internal/application/vo"
	"jcourse_go/internal/domain/shared"
)

type UserPointQueryService interface {
	GetUserPoint(ctx context.Context, userID shared.IDType) (int, []vo.UserPointVO, error)
}

type userPointQueryService struct {
	db *gorm.DB
}

func (u *userPointQueryService) GetUserPoint(ctx context.Context, userID shared.IDType) (int, []vo.UserPointVO, error) {
	// TODO implement me
	panic("implement me")
}

func NewUserPointQueryService(db *gorm.DB) UserPointQueryService {
	return &userPointQueryService{db: db}
}
