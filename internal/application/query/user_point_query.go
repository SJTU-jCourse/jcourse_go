package query

import (
	"context"

	"gorm.io/gorm"

	"jcourse_go/internal/application/vo"
	"jcourse_go/internal/domain/shared"
	"jcourse_go/internal/infrastructure/entity"
)

type UserPointQueryService interface {
	GetUserPoint(ctx context.Context, userID shared.IDType) (int, []vo.UserPointVO, error)
}

type userPointQueryService struct {
	db *gorm.DB
}

func (u *userPointQueryService) GetUserPoint(ctx context.Context, userID shared.IDType) (int, []vo.UserPointVO, error) {
	type row struct {
		entity.UserPoint
		TotalValue int `gorm:"column:total_value"`
	}

	up := make([]row, 0)
	if err := u.db.WithContext(ctx).
		Select("*, sum(value) as total_value").
		Model(&entity.UserPoint{}).
		Where("user_id = ?", userID).
		Find(&up).Error; err != nil {
		return 0, nil, err
	}

	res := make([]vo.UserPointVO, 0)
	if len(up) == 0 {
		return 0, res, nil
	}

	for _, e := range up {
		upVO := vo.NewUserPointFromEntity(&e.UserPoint)
		res = append(res, upVO)
	}
	return up[0].TotalValue, res, nil
}

func NewUserPointQueryService(db *gorm.DB) UserPointQueryService {
	return &userPointQueryService{db: db}
}
