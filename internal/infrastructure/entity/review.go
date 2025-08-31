package entity

import (
	"time"
)

type ReviewPO struct {
	ID int64 `gorm:"primarykey"`

	CourseID    int64 `gorm:"index;index:uniq_course_review,unique"`
	UserID      int64 `gorm:"index;index:uniq_course_review,unique"`
	Comment     string
	Rating      int64  `gorm:"index"`
	Semester    string `gorm:"index"`
	IsAnonymous bool
	Grade       string // 成绩

	SearchIndex SearchIndex `gorm:"->:false;<-"`

	CreatedAt time.Time `gorm:"index"`
	UpdatedAt time.Time `gorm:"index"`
}

func (po *ReviewPO) TableName() string {
	return "reviews"
}

type ReviewRevisionPO struct {
	ID          int64 `gorm:"primarykey"`
	ReviewID    int64 `gorm:"index"`
	Comment     string
	Rating      int64
	Semester    string
	IsAnonymous bool
	Grade       string    // 成绩
	UpdatedBy   int64     `gorm:"index"`
	CreatedAt   time.Time `gorm:"index"`
}

func (po *ReviewRevisionPO) TableName() string {
	return "review_revisions"
}

type ReviewReactionPO struct {
	ID int64 `gorm:"primarykey"`

	ReviewID int64  `gorm:"index;index:uniq_review_reaction,unique"`
	UserID   int64  `gorm:"index;index:uniq_review_reaction,unique"`
	Reaction string `gorm:"index;index:uniq_review_reaction,unique"`

	CreatedAt time.Time `gorm:"index"`
}

func (po *ReviewReactionPO) TableName() string {
	return "review_reactions"
}
