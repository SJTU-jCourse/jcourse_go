package vo

type UserActivityVO struct {
	ReviewCount          int64 `json:"review_count"`
	LikeReceive          int64 `json:"like_receive"`
	TipReceive           int64 `json:"tip_receive"`
	FollowingCourseCount int64 `json:"following_course_count"`
}

type UserInfoVO struct {
	ID       int64  `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
	Role     string `json:"role,omitempty"`
}
