package po

import "gorm.io/gorm"

type ReviewPO struct {
	gorm.Model
	CourseID    int64 `gorm:"index;index:uniq_course_review,unique"`
	UserID      int64 `gorm:"index;index:uniq_course_review,unique"`
	Comment     string
	Rating      int64  `gorm:"index"`
	Semester    string `gorm:"index;index:uniq_course_review,unique"`
	IsAnonymous bool
	SearchIndex SearchIndex // `gorm:"index:idx_search, type:gin"`
}

func (po *ReviewPO) TableName() string {
	return "reviews"
}

type ReviewRevisionPO struct {
	gorm.Model
	ReviewID    int64 `gorm:"index"`
	UserID      int64 `gorm:"index"`
	Comment     string
	Rating      int64
	Semester    string `gorm:"index"`
	IsAnonymous bool
}

func (po *ReviewRevisionPO) TableName() string {
	return "review_revisions"
}

type ReviewReactionPO struct {
	gorm.Model
	ReviewID    int64  `gorm:"index"`
	UserID      int64  `gorm:"index"`
	Reaction    string `gorm:"index"`
	IsAnonymous bool
}

func (po *ReviewReactionPO) TableName() string {
	return "review_reactions"
}

type ReviewReplyPO struct {
	gorm.Model
	ReviewID       int64 `gorm:"index"`
	UserID         int64 `gorm:"index"`
	ReplyToReplyID int64 `gorm:"index"`
	Comment        string
	IsAnonymous    bool
}

func (po *ReviewReplyPO) TableName() string {
	return "review_replies"
}
