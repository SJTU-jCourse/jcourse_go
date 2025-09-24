package entity

import "time"

type VerificationCode struct {
	ID        int64  `gorm:"primaryKey"`
	Email     string `gorm:"index;unique"`
	Code      string
	ExpiresAt time.Time `gorm:"index"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (po *VerificationCode) TableName() string {
	return "verification_code"
}

type Session struct {
	ID        int64     `gorm:"primaryKey"`
	SessionID string    `gorm:"index;unique"`
	UserID    int64     `gorm:"index"`
	CreatedAt time.Time `gorm:"index"`
	ExpiresAt time.Time `gorm:"index"`
}

func (po *Session) TableName() string {
	return "session"
}
