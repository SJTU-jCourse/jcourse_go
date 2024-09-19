package repository

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"jcourse_go/constant"
	"jcourse_go/model/po"
)

type IUserQuery interface {
	GetUserDetail(ctx context.Context, opts ...DBOption) (*po.UserPO, error)
	GetUserList(ctx context.Context, opts ...DBOption) ([]po.UserPO, error)
	GetUserCount(ctx context.Context, opts ...DBOption) (int64, error)
	GetUserByID(ctx context.Context, userID int64) (*po.UserPO, error)
	GetUserByIDs(ctx context.Context, userIDs []int64) (map[int64]po.UserPO, error)
	UpdateUserByID(ctx context.Context, user *po.UserPO) error
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

func (q *UserQuery) GetUserByID(ctx context.Context, userID int64) (*po.UserPO, error) {
	db := q.optionDB(ctx, WithID(userID))
	userPO := po.UserPO{}
	result := db.Find(&userPO)
	if result.Error != nil {
		return &userPO, result.Error
	}
	return &userPO, nil
}

func (q *UserQuery) optionDB(ctx context.Context, opts ...DBOption) *gorm.DB {
	db := q.db.WithContext(ctx).Model(po.UserPO{})
	for _, opt := range opts {
		db = opt(db)
	}
	return db
}

func (q *UserQuery) GetUserDetail(ctx context.Context, opts ...DBOption) (*po.UserPO, error) {
	db := q.optionDB(ctx, opts...)
	user := po.UserPO{}
	result := db.First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (q *UserQuery) GetUserList(ctx context.Context, opts ...DBOption) ([]po.UserPO, error) {
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

func (q *UserQuery) UpdateUserByID(ctx context.Context, user *po.UserPO) error {
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
	result := q.optionDB(ctx).Debug().Create(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (q *UserQuery) ResetUserPassword(ctx context.Context, userID int64, passwordStore string) error {
	result := q.optionDB(ctx, WithID(userID)).Debug().Update("password", passwordStore)
	return result.Error
}
