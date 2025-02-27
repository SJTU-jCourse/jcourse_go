package po

import (
	"time"
)

type ReviewPO struct {
	ID        int64     `gorm:"primarykey"`
	CreatedAt time.Time `gorm:"index"`
	UpdatedAt time.Time `gorm:"index"`

	CourseID    int64     `gorm:"index;index:uniq_course_review,unique"`
	Course      *CoursePO `gorm:"constraint:OnDelete:CASCADE;"`
	UserID      int64     `gorm:"index;index:uniq_course_review,unique"`
	User        *UserPO   `gorm:"constraint:OnDelete:CASCADE;"`
	Comment     string
	Rating      int64  `gorm:"index"`
	Semester    string `gorm:"index"`
	IsAnonymous bool
	Grade       string // 成绩

	Revisions []ReviewRevisionPO `gorm:"foreignKey:ReviewID"`
	Reaction  []ReviewReactionPO `gorm:"foreignKey:ReviewID"`

	SearchIndex SearchIndex `gorm:"->:false;<-"`
}

func (po *ReviewPO) TableName() string {
	return "reviews"
}

type ReviewRevisionPO struct {
	ID        int64     `gorm:"primarykey"`
	CreatedAt time.Time `gorm:"index"`

	ReviewID    int64     `gorm:"index"`
	Review      *ReviewPO `gorm:"constraint:OnDelete:CASCADE;"`
	UserID      int64     `gorm:"index"`
	User        *UserPO   `gorm:"constraint:OnDelete:CASCADE;"`
	CourseID    int64     `gorm:"index"`
	Comment     string
	Rating      int64
	Semester    string
	IsAnonymous bool
	Grade       string // 成绩
}

func (po *ReviewRevisionPO) TableName() string {
	return "review_revisions"
}

type ReviewReactionPO struct {
	ID        int64     `gorm:"primarykey"`
	CreatedAt time.Time `gorm:"index"`

	ReviewID int64     `gorm:"index"`
	Review   *ReviewPO `gorm:"constraint:OnDelete:CASCADE;"`
	UserID   int64     `gorm:"index"`
	User     *UserPO   `gorm:"constraint:OnDelete:CASCADE;"`
	Reaction string    `gorm:"index"`
}

func (po *ReviewReactionPO) TableName() string {
	return "review_reactions"
}
