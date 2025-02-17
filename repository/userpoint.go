package repository

import (
	"context"

	"jcourse_go/model/po"
	"jcourse_go/model/types"

	"gorm.io/gorm"
)

type IUserPointDetailQuery interface {
	GetUserPoint(ctx context.Context, userID int64) (int64, error)
	GetUserPointDetail(ctx context.Context, opts ...DBOption) ([]po.UserPointDetailPO, error)
	GetUserPointDetailCount(ctx context.Context, opts ...DBOption) (int64, error)
	CreateUserPointDetail(ctx context.Context, userID int64, eventType types.PointEventType, value int64, description string) error
}

type UserPointDetailQuery struct {
	db *gorm.DB
}

func (q *UserPointDetailQuery) GetUserPoint(ctx context.Context, userID int64) (int64, error) {
	result := struct {
		Value int64 `json:"value"`
	}{}
	if err := q.db.WithContext(ctx).
		Model(&po.UserPointDetailPO{}).
		Select("sum(value) as value").
		Where("user_id = ?", userID).Find(&result).Error; err != nil {
		return 0, err
	}
	return result.Value, nil
}

func (q *UserPointDetailQuery) optionDB(ctx context.Context, opts ...DBOption) *gorm.DB {
	db := q.db.WithContext(ctx).Model(&po.UserPointDetailPO{})
	for _, opt := range opts {
		db = opt(db)
	}
	return db
}

func (q *UserPointDetailQuery) GetUserPointDetail(ctx context.Context, opts ...DBOption) ([]po.UserPointDetailPO, error) {
	db := q.optionDB(ctx, opts...)
	userPointDetailPOs := make([]po.UserPointDetailPO, 0)

	result := db.Find(&userPointDetailPOs)
	if result.Error != nil {
		return userPointDetailPOs, result.Error
	}
	return userPointDetailPOs, nil
}

func (q *UserPointDetailQuery) GetUserPointDetailCount(ctx context.Context, opts ...DBOption) (int64, error) {
	db := q.optionDB(ctx, opts...)
	var count int64
	result := db.Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}

func (q *UserPointDetailQuery) CreateUserPointDetail(ctx context.Context, userID int64, eventType types.PointEventType, value int64, description string) error {
	userPointDetail := po.UserPointDetailPO{
		UserID:      userID,
		EventType:   string(eventType),
		Value:       value,
		Description: description,
	}
	result := q.optionDB(ctx).Create(&userPointDetail)
	return result.Error
}

func NewUserPointDetailQuery(db *gorm.DB) IUserPointDetailQuery {
	return &UserPointDetailQuery{
		db: db,
	}
}
