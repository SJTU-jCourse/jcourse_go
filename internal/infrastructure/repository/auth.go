package repository

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"jcourse_go/internal/domain/auth"
	"jcourse_go/internal/domain/shared"
	"jcourse_go/internal/infrastructure/entity"
)

type VerificationCodeRepository struct {
	db *gorm.DB
}

func (r *VerificationCodeRepository) Get(ctx context.Context, email string) (*auth.VerificationCode, error) {
	e := &entity.VerificationCode{}
	if err := r.db.Model(&entity.VerificationCode{}).Where("email = ?", email).First(&e).Error; err != nil {
		return nil, err
	}
	d := newVerificationCodeDomainFromEntity(e)
	return &d, nil
}

func (r *VerificationCodeRepository) Save(ctx context.Context, code *auth.VerificationCode) error {
	e := newVerificationCodeEntityFromDomain(code)
	if err := r.db.Model(&entity.VerificationCode{}).
		Clauses(clause.OnConflict{UpdateAll: true}).
		Create(&e).Error; err != nil {
		return err
	}
	return nil
}

func (r *VerificationCodeRepository) Delete(ctx context.Context, email string) error {
	if err := r.db.Delete(&entity.VerificationCode{}, "email = ?", email).Error; err != nil {
		return err
	}
	return nil
}

func NewVerificationCodeRepository(db *gorm.DB) auth.VerificationCodeRepository {
	return &VerificationCodeRepository{db: db}
}

type SessionRepository struct {
	db *gorm.DB
}

func (r *SessionRepository) Get(ctx context.Context, sessionID string) (*auth.Session, error) {
	e := &entity.Session{}
	if err := r.db.Model(&entity.Session{}).Where("session_id = ?", sessionID).First(&e).Error; err != nil {
		return nil, err
	}
	d := newSessionDomainFromEntity(e)
	return &d, nil
}

func (r *SessionRepository) GetByUser(ctx context.Context, userID shared.IDType) (*auth.Session, error) {
	e := &entity.Session{}
	if err := r.db.Model(&entity.Session{}).Where("user_id = ?", userID).First(&e).Error; err != nil {
		return nil, err
	}
	d := newSessionDomainFromEntity(e)
	return &d, nil
}

func (r *SessionRepository) Delete(ctx context.Context, sessionID string) error {
	if err := r.db.Delete(&entity.Session{}, "session_id = ?", sessionID).Error; err != nil {
		return err
	}
	return nil
}

func (r *SessionRepository) DeleteByUser(ctx context.Context, userID shared.IDType) error {
	if err := r.db.Delete(&entity.Session{}, "user_id = ?", userID).Error; err != nil {
		return err
	}
	return nil
}

func (r *SessionRepository) Save(ctx context.Context, session *auth.Session) error {
	e := newSessionEntityFromDomain(session)
	if err := r.db.Model(&entity.Session{}).Create(&e).Error; err != nil {
		return err
	}
	return nil
}

func NewSessionRepository(db *gorm.DB) auth.SessionRepository {
	return &SessionRepository{db: db}
}

type AuthUserRepository struct {
	db *gorm.DB
}

func (r *AuthUserRepository) FindByEmail(ctx context.Context, email string) (*auth.AuthUser, error) {
	e := &entity.User{}
	if err := r.db.Model(&entity.User{}).Where("email = ?", email).First(&e).Error; err != nil {
		return nil, err
	}
	d := newAuthUserDomainFromEntity(e)
	return &d, nil
}

func (r *AuthUserRepository) Save(ctx context.Context, user *auth.AuthUser) error {
	e := newAuthUserEntityFromDomain(user)
	if user.ID == 0 {
		if err := r.db.Model(&entity.User{}).Create(&e).Error; err != nil {
			return err
		}
		user.ID = shared.IDType(e.ID)
		return nil
	}

	if err := r.db.Model(&entity.User{}).Where("id = ?", user.ID).Updates(&e).Error; err != nil {
		return err
	}
	return nil
}

func NewAuthUserRepository(db *gorm.DB) auth.UserRepository {
	return &AuthUserRepository{db: db}
}

func newVerificationCodeDomainFromEntity(v *entity.VerificationCode) auth.VerificationCode {
	return auth.VerificationCode{
		Email:     v.Email,
		Code:      v.Code,
		ExpiresAt: v.ExpiresAt,
		CreatedAt: v.CreatedAt,
		UpdatedAt: v.UpdatedAt,
	}
}

func newVerificationCodeEntityFromDomain(v *auth.VerificationCode) entity.VerificationCode {
	return entity.VerificationCode{
		Email:     v.Email,
		Code:      v.Code,
		ExpiresAt: v.ExpiresAt,
		CreatedAt: v.CreatedAt,
		UpdatedAt: v.UpdatedAt,
	}
}

func newSessionDomainFromEntity(s *entity.Session) auth.Session {
	return auth.Session{
		SessionID: s.SessionID,
		UserID:    shared.IDType(s.UserID),
		CreatedAt: s.CreatedAt,
		ExpiresAt: s.ExpiresAt,
	}
}

func newSessionEntityFromDomain(s *auth.Session) entity.Session {
	return entity.Session{
		SessionID: s.SessionID,
		UserID:    int64(s.UserID),
		CreatedAt: s.CreatedAt,
		ExpiresAt: s.ExpiresAt,
	}
}

func newAuthUserDomainFromEntity(u *entity.User) auth.AuthUser {
	return auth.AuthUser{
		ID:            shared.IDType(u.ID),
		Email:         u.Email,
		Password:      u.Password,
		UserRole:      shared.UserRole(u.UserRole),
		LastSeenAt:    u.LastSeenAt,
		SuspendedTill: u.SuspendedTill,
		CreatedAt:     u.CreatedAt,
		UpdatedAt:     u.UpdatedAt,
	}
}

func newAuthUserEntityFromDomain(u *auth.AuthUser) entity.User {
	return entity.User{
		ID:            int64(u.ID),
		Email:         u.Email,
		Password:      u.Password,
		UserRole:      string(u.UserRole),
		LastSeenAt:    u.LastSeenAt,
		SuspendedTill: u.SuspendedTill,
		CreatedAt:     u.CreatedAt,
		UpdatedAt:     u.UpdatedAt,
	}
}
