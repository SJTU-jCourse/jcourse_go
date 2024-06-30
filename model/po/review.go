package po

import "gorm.io/gorm"

type CourseReviewPO struct {
	gorm.Model
	CourseID        int64
	OfferedCourseID int64
	UserID          int64
	Comment         string
	Rate            int64
	Semester        string
	IsAnonymous     bool
	Version         int64
}

func (po *CourseReviewPO) TableName() string {
	return "course_reviews"
}

type ReviewReactionPO struct {
	gorm.Model
	CourseReviewID int64
	UserID         int64
	IsAnonymous    bool
	Reaction       string
}

func (po *ReviewReactionPO) TableName() string {
	return "review_reactions"
}

type ReviewReplyPO struct {
	gorm.Model
	CourseReviewID int64
	UserID         int64
	ReplyToReplyID int64
	Comment        string
	IsAnonymous    bool

	Version int64
}

func (po *ReviewReplyPO) TableName() string {
	return "review_replies"
}
