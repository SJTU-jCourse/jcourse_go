package vo

import "jcourse_go/internal/infrastructure/entity"

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

func NewUserInfoVOFromEntity(e *entity.User) UserInfoVO {
	return UserInfoVO{
		ID:       e.ID,
		Username: e.Username,
		Role:     e.UserRole,
	}
}
