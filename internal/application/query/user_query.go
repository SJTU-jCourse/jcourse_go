package query

import (
	"context"

	"gorm.io/gorm"

	"jcourse_go/internal/application/vo"
	"jcourse_go/internal/domain/shared"
)

type UserQueryService interface {
	GetUserInfo(ctx context.Context, userID shared.IDType) (*vo.UserInfoVO, error)
}

type userQueryService struct {
	db *gorm.DB
}

func (u *userQueryService) GetUserInfo(ctx context.Context, userID shared.IDType) (*vo.UserInfoVO, error) {
	panic("implement me")
}

func NewUserQueryService(db *gorm.DB) UserQueryService {
	return &userQueryService{db: db}
}
