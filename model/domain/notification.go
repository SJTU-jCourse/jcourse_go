package domain

type NotificationRelatedType = string

const (
	NotificationRelatedTypeCourse NotificationRelatedType = "course"
	NotificationRelatedTypeReview NotificationRelatedType = "review"
	NotificationRelatedTypeReply  NotificationRelatedType = "reply"
)

type Notification struct {
	NotificationID int64
	UserID         int64
	RelatedType    NotificationRelatedType
	RelatedID      int64
	Title          string
	Content        string
	IsSolved       bool
}
