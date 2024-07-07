package domain

type UserActivity struct {
	User User

	CourseReviews []Review
	ReviewReplies []ReviewReply

	ReviewedCourse []OfferedCourse

	RewardRecords []RewardRecord

	Notifications []Notification

	FollowingUsers  []User
	FollowedUsers   []User
	FollowingCourse []OfferedCourse
}
