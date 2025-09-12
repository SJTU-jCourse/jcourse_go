package entity

import (
	"time"
)

type Semester struct {
	ID        int64
	Name      string `gorm:"uniqueIndex"`
	Available bool
	CreatedAt time.Time
}

func (s *Semester) TableName() string {
	return "semester"
}

type Department struct {
	ID        int64
	Name      string `gorm:"uniqueIndex"`
	CreatedAt time.Time
}

func (d *Department) TableName() string {
	return "department"
}
