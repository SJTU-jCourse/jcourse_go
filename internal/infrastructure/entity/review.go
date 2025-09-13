package entity

import (
	"time"
)

type Review struct {
	ID int64 `gorm:"primaryKey"`

	CourseID    int64 `gorm:"index;index:uniq_course_review,unique"`
	Course      *Course
	UserID      int64 `gorm:"index;index:uniq_course_review,unique"`
	User        *User
	Comment     string
	Rating      int64  `gorm:"index"`
	Semester    string `gorm:"index"`
	IsAnonymous bool
	Score       string // 成绩

	Revisions []*ReviewRevision
	Reactions []*ReviewReaction

	CreatedAt time.Time `gorm:"index"`
	UpdatedAt time.Time `gorm:"index"`
}

func (po *Review) TableName() string {
	return "review"
}

type ReviewRevision struct {
	ID          int64 `gorm:"primaryKey"`
	UserID      int64 `gorm:"index"`
	User        *User
	ReviewID    int64 `gorm:"index"`
	Review      *Review
	Comment     string
	Rating      int64
	Semester    string
	IsAnonymous bool
	Grade       string    // 成绩
	UpdatedBy   int64     `gorm:"index"`
	CreatedAt   time.Time `gorm:"index"`
}

func (po *ReviewRevision) TableName() string {
	return "review_revision"
}

type ReviewReaction struct {
	ID int64 `gorm:"primaryKey"`

	ReviewID int64 `gorm:"index;index:uniq_review_reaction,unique"`
	Review   *Review
	UserID   int64 `gorm:"index;index:uniq_review_reaction,unique"`
	User     *User
	Reaction string `gorm:"index;index:uniq_review_reaction,unique"`

	CreatedAt time.Time `gorm:"index"`
}

func (po *ReviewReaction) TableName() string {
	return "review_reaction"
}
