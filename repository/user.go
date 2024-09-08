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

type IUserProfileQuery interface {
	GetUserProfileByIDs(ctx context.Context, userIDs []int64) (map[int64]po.UserProfilePO, error)
	GetUserProfileByID(ctx context.Context, userID int64) (*po.UserProfilePO, error)
	GetUserProfileList(ctx context.Context, opts ...DBOption) ([]po.UserProfilePO, error)
	GetUserProfileCount(ctx context.Context, opts ...DBOption) (int64, error)
	UpdateUserProfileByID(ctx context.Context, userProfile *po.UserProfilePO) error
	// CreateUserProfile(ctx context.Context, userID int64) (*po.UserProfilePO, error)
}

type UserProfileQuery struct {
	db *gorm.DB
}

func (u *UserProfileQuery) GetUserProfileByIDs(ctx context.Context, userIDs []int64) (map[int64]po.UserProfilePO, error) {
	db := u.optionDB(ctx)
	userProfiles := make([]po.UserProfilePO, 0)
	userProfileMap := make(map[int64]po.UserProfilePO)
	result := db.Where("user_id in ?", userIDs).Find(&userProfiles)
	if result.Error != nil {
		return userProfileMap, result.Error
	}
	for _, userProfile := range userProfiles {
		userProfileMap[userProfile.UserID] = userProfile
	}
	return userProfileMap, nil
}

func (u *UserProfileQuery) GetUserProfileByID(ctx context.Context, userID int64) (*po.UserProfilePO, error) {
	db := u.optionDB(ctx, WithUserID(userID))
	userProfilePO := po.UserProfilePO{}
	result := db.Find(&userProfilePO)
	if result.Error != nil {
		return &userProfilePO, result.Error
	}
	return &userProfilePO, nil
}

func (u *UserProfileQuery) optionDB(ctx context.Context, opts ...DBOption) *gorm.DB {
	db := u.db.WithContext(ctx).Model(&po.UserProfilePO{})
	for _, opt := range opts {
		db = opt(db)
	}
	return db
}

func NewUserProfileQuery() IUserProfileQuery {
	return &UserProfileQuery{db: dal.GetDBClient()}
}

func NewUserQuery() IUserQuery {
	return &UserQuery{
		db: dal.GetDBClient(),
	}
}

func (q *UserProfileQuery) GetUserProfileList(ctx context.Context, opts ...DBOption) ([]po.UserProfilePO, error) {
	db := q.optionDB(ctx, opts...)
	userProfilePOs := make([]po.UserProfilePO, 0)
	result := db.Find(&userProfilePOs)
	if result.Error != nil {
		return userProfilePOs, result.Error
	}
	return userProfilePOs, nil
}

func (q *UserProfileQuery) GetUserProfileCount(ctx context.Context, opts ...DBOption) (int64, error) {
	db := q.optionDB(ctx, opts...)
	count := int64(0)
	result := db.Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}

func (q *UserProfileQuery) UpdateUserProfileByID(ctx context.Context, userProfile *po.UserProfilePO) error {
	result := q.optionDB(ctx, WithUserID(userProfile.UserID)).Save(&userProfile).Error
	return result
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
	result := q.optionDB(ctx, WithID(int64(user.ID))).Save(&user).Error
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
