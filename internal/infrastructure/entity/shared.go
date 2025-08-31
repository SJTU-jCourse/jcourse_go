package entity

import "time"

type Semester struct {
	ID        int64
	Name      string
	Available bool
	CreatedAt time.Time
}

func (s *Semester) TableName() string {
	return "semesters"
}
