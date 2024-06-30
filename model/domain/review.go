package domain

import "time"

type CourseReview struct {
	ID              int64
	CourseID        int64
	OfferedCourseID int64
	UserID          int64
	Comment         string
	Rate            int64
	IsAnonymous     bool
	CreatedAt       time.Time
	UpdatedAt       time.Time
	Version         int64

	Relies    []ReviewReply
	Reactions []ReviewReaction
}

type ReviewReaction struct {
	ID             uint
	CourseReviewID int64
	UserID         int64
	IsAnonymous    bool
	Reaction       string
}

type ReviewReply struct {
	ID             uint
	CourseReviewID int64
	UserID         int64
	ReplyToReplyID int64
	Comment        string
	IsAnonymous    bool
	Version        int64
}
