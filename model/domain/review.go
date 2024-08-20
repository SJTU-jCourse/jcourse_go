package domain

import "time"

type ReviewFilter struct {
	Page     int64
	PageSize int64
	CourseID int64
	Semester string
	UserID   int64
	ReviewID int64
}

type Review struct {
	ID          int64
	Course      Course
	User        User
	Comment     string
	Rating      int64
	Semester    string
	IsAnonymous bool
	CreatedAt   time.Time
	UpdatedAt   time.Time

	Relies     []ReviewReply
	ReplyCount int64

	Reactions []ReviewReaction
}

type ReviewReaction struct {
	ID        uint
	ReviewID  int64
	UserID    int64
	Reaction  string
	CreatedAt time.Time
}

type ReviewReply struct {
	ID             uint
	ReviewID       int64
	UserID         int64
	ReplyToReplyID int64
	Comment        string
}
