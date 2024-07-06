package repository

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
	"jcourse_go/constant"
	"jcourse_go/dal"
	"jcourse_go/model/po"
)

type DBOption func(*gorm.DB) *gorm.DB

type IUserQuery interface {
	GetUserDetail(ctx context.Context, opts ...DBOption) (*po.UserPO, error)
	GetUserList(ctx context.Context, opts ...DBOption) ([]po.UserPO, error)
	WithID(id int64) DBOption
	WithEmail(email string) DBOption
	WithPassword(password string) DBOption
	CreateUser(ctx context.Context, email string, password string) (*po.UserPO, error)
	ResetUserPassword(ctx context.Context, userID int64, password string) error
}

func NewUserQuery() IUserQuery {
	return &UserQuery{
		db: dal.GetDBClient(),
	}
}

type UserQuery struct {
	db *gorm.DB
}

func (q *UserQuery) WithEmail(email string) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("email = ?", email)
	}
}

func (q *UserQuery) WithID(id int64) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	}
}

func (q *UserQuery) WithPassword(password string) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("password = ?", password)
	}
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
	result := q.optionDB(ctx, q.WithID(userID)).Debug().Update("password", passwordStore)
	return result.Error
}
