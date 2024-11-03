package repository

import (
	"context"
	"time"

	"gorm.io/gorm"

	"jcourse_go/constant"
	"jcourse_go/model/po"
)

type IUserQuery interface {
	GetUser(ctx context.Context, opts ...DBOption) ([]po.UserPO, error)
	GetUserCount(ctx context.Context, opts ...DBOption) (int64, error)
	GetUserByIDs(ctx context.Context, userIDs []int64) (map[int64]po.UserPO, error)
	UpdateUser(ctx context.Context, user po.UserPO) error
	CreateUser(ctx context.Context, email string, password string) (*po.UserPO, error)
	ResetUserPassword(ctx context.Context, userID int64, password string) error
}

func NewUserQuery(db *gorm.DB) IUserQuery {
	return &UserQuery{
		db: db,
	}
}

type UserQuery struct {
	db *gorm.DB
}

func (q *UserQuery) GetUserByIDs(ctx context.Context, userIDs []int64) (map[int64]po.UserPO, error) {
	db := q.optionDB(ctx)
	userPOs := make([]po.UserPO, 0)
	userMap := make(map[int64]po.UserPO)
	result := db.Where("id in ?", userIDs).Find(&userPOs)
	if result.Error != nil {
		return userMap, result.Error
	}
	for _, userPO := range userPOs {
		userMap[int64(userPO.ID)] = userPO
	}
	return userMap, nil
}

func (q *UserQuery) optionDB(ctx context.Context, opts ...DBOption) *gorm.DB {
	db := q.db.WithContext(ctx).Model(po.UserPO{})
	for _, opt := range opts {
		db = opt(db)
	}
	return db
}

func (q *UserQuery) GetUser(ctx context.Context, opts ...DBOption) ([]po.UserPO, error) {
	db := q.optionDB(ctx, opts...)
	userPOs := make([]po.UserPO, 0)
	result := db.Find(&userPOs)
	if result.Error != nil {
		return userPOs, result.Error
	}
	return userPOs, nil
}

func (q *UserQuery) GetUserCount(ctx context.Context, opts ...DBOption) (int64, error) {
	db := q.optionDB(ctx, opts...)
	count := int64(0)
	result := db.Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}

func (q *UserQuery) UpdateUser(ctx context.Context, user po.UserPO) error {
	result := q.optionDB(ctx, WithID(int64(user.ID))).Updates(&user).Error
	return result
}

func (q *UserQuery) CreateUser(ctx context.Context, email string, passwordStore string) (*po.UserPO, error) {
	user := po.UserPO{
		Username:   email,
		Email:      email,
		UserRole:   constant.UserRoleNormal,
		LastSeenAt: time.Now(),
		Password:   passwordStore,
	}
	result := q.optionDB(ctx).Create(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (q *UserQuery) ResetUserPassword(ctx context.Context, userID int64, passwordStore string) error {
	result := q.optionDB(ctx, WithID(userID)).Update("password", passwordStore)
	return result.Error
}

type IUserPointDetailQuery interface {
	GetUserPointDetail(ctx context.Context, opts ...DBOption) ([]po.UserPointDetailPO, error)
	GetUserPointDetailCount(ctx context.Context, opts ...DBOption) (int64, error)
	CreateUserPointDetail(ctx context.Context, userID int64, eventType po.PointEventType, value int64, description string) error
}

type UserPointDetailQuery struct {
	db *gorm.DB
}

func (q *UserPointDetailQuery) optionDB(opts ...DBOption) *gorm.DB {
	db := q.db.Model(po.UserPointDetailPO{})
	for _, opt := range opts {
		db = opt(db)
	}
	return db
}

func (q *UserPointDetailQuery) GetUserPointDetail(ctx context.Context, opts ...DBOption) ([]po.UserPointDetailPO, error) {
	db := q.optionDB(opts...)
	userPointDetailPOs := make([]po.UserPointDetailPO, 0)

	result := db.Find(&userPointDetailPOs)
	if result.Error != nil {
		return userPointDetailPOs, result.Error
	}
	return userPointDetailPOs, nil
}

func (q *UserPointDetailQuery) GetUserPointDetailCount(ctx context.Context, opts ...DBOption) (int64, error) {
	db := q.optionDB(opts...)
	var count int64
	result := db.Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}

func (q *UserPointDetailQuery) CreateUserPointDetail(ctx context.Context, userID int64, eventType po.PointEventType, value int64, description string) error {
	userPointDetail := po.UserPointDetailPO{
		UserID: userID,
		PointEvent: po.PointEvent{
			EventType:   eventType,
			Value:       value,
			Description: description,
		},
	}
	result := q.optionDB().Create(&userPointDetail)
	return result.Error
}

func NewUserPointDetailQuery(db *gorm.DB) IUserPointDetailQuery {
	return &UserPointDetailQuery{
		db: db,
	}
}
