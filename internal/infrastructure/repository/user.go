package repository

import (
	"context"

	"gorm.io/gorm"

	"jcourse_go/internal/domain/shared"
	"jcourse_go/internal/domain/user"
	"jcourse_go/internal/infrastructure/entity"
)

type UserProfileRepository struct {
	db *gorm.DB
}

func (r *UserProfileRepository) Get(ctx context.Context, id shared.IDType) (*user.User, error) {
	e := &entity.User{}
	if err := r.db.Model(&entity.User{}).Where("id = ?", id).First(&e).Error; err != nil {
		return nil, err
	}
	d := newUserDomainFromEntity(e)
	return &d, nil
}

func (r *UserProfileRepository) Save(ctx context.Context, user *user.User) error {
	e := newUserEntityFromDomain(user)
	if err := r.db.Model(&entity.User{}).
		Where("id = ?", user.ID).
		Updates(&e).Error; err != nil {
		return err
	}
	return nil
}

func (r *UserProfileRepository) UpdateUserInfo(ctx context.Context, userID shared.IDType, cmd user.UpdateUserInfoCommand) error {
	updates := map[string]interface{}{
		"nickname": cmd.Nickname,
		"bio":      cmd.Bio,
	}
	if cmd.Avatar != "" {
		updates["avatar"] = cmd.Avatar
	}
	if err := r.db.Model(&entity.User{}).
		Where("id = ?", userID).
		Updates(updates).Error; err != nil {
		return err
	}
	return nil
}

func NewUserProfileRepository(db *gorm.DB) user.UserProfileRepository {
	return &UserProfileRepository{db: db}
}

type UserPointRepository struct {
	db *gorm.DB
}

func (r *UserPointRepository) FindByUserID(ctx context.Context, userID shared.IDType) ([]user.UserPoint, error) {
	var entities []entity.UserPoint
	if err := r.db.Model(&entity.UserPoint{}).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&entities).Error; err != nil {
		return nil, err
	}
	var domains []user.UserPoint
	for _, e := range entities {
		domains = append(domains, newUserPointDomainFromEntity(&e))
	}
	return domains, nil
}

func (r *UserPointRepository) Save(ctx context.Context, userPoint *user.UserPoint) error {
	e := newUserPointEntityFromDomain(userPoint)
	if err := r.db.Model(&entity.UserPoint{}).Create(&e).Error; err != nil {
		return err
	}
	return nil
}

func (r *UserPointRepository) Transaction(ctx context.Context, transaction *user.UserPointTransaction) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		fromPoint := entity.UserPoint{
			UserID:      int64(transaction.FromID),
			Type:        string(user.PointEventTransferOut),
			Description: transaction.Reason,
			Value:       -transaction.Value,
			CreatedAt:   transaction.CreatedAt,
		}

		toPoint := entity.UserPoint{
			UserID:      int64(transaction.ToID),
			Type:        string(user.PointEventTransferIn),
			Description: transaction.Reason,
			Value:       transaction.Value,
			CreatedAt:   transaction.CreatedAt,
		}

		if err := tx.Model(&entity.UserPoint{}).Create(&fromPoint).Error; err != nil {
			return err
		}

		if err := tx.Model(&entity.UserPoint{}).Create(&toPoint).Error; err != nil {
			return err
		}

		return nil
	})
}

func NewUserPointRepository(db *gorm.DB) user.UserPointRepository {
	return &UserPointRepository{db: db}
}

func newUserDomainFromEntity(u *entity.User) user.User {
	return user.User{
		ID:        shared.IDType(u.ID),
		Username:  u.Username,
		Bio:       "",
		Email:     "",
		Role:      shared.UserRole(u.UserRole),
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func newUserEntityFromDomain(u *user.User) entity.User {
	return entity.User{
		ID:        int64(u.ID),
		Username:  u.Username,
		UserRole:  string(u.Role),
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func newUserPointDomainFromEntity(u *entity.UserPoint) user.UserPoint {
	t, _ := user.PointEventString(u.Type)
	return user.UserPoint{
		ID:          shared.IDType(u.ID),
		Value:       u.Value,
		Event:       t,
		Description: u.Description,
		CreatedAt:   u.CreatedAt,
	}
}

func newUserPointEntityFromDomain(u *user.UserPoint) entity.UserPoint {
	return entity.UserPoint{
		ID:          int64(u.ID),
		UserID:      int64(u.UserID),
		Type:        u.Event.String(),
		Description: u.Description,
		Value:       u.Value,
		CreatedAt:   u.CreatedAt,
	}
}
