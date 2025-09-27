package entity

import "time"

type Report struct {
	ID        int64 `gorm:"primaryKey"`
	UserID    int64
	Content   string
	Reply     string
	Solved    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
