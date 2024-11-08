package model

type UserRole = string

const (
	UserRoleNormal UserRole = "normal"
	UserRoleAdmin  UserRole = "admin"
)

type UserType = string

const (
	UserTypeStudent UserType = "student"
	UserTypeFaculty UserType = "faculty"
)

type UserFilterForQuery struct {
	PaginationFilterForQuery
}

type UserDetail struct {
	UserMinimal
	Bio        string `json:"bio"`
	Email      string `json:"email"`
	Role       string `json:"role"`
	Type       string `json:"type"`
	Department string `json:"department"`
	Major      string `json:"major"`
	Grade      string `json:"grade"`
	Points     int64  `json:"points"`
}

type UserActivity struct {
	ReviewCount          int64 `json:"review_count"`
	LikeReceive          int64 `json:"like_receive"`
	TipReceive           int64 `json:"tip_receive"`
	FollowingCourseCount int64 `json:"following_course_count"`
}

type UserMinimal struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
}
type PointEventType = string

const (
	PointEventReview      PointEventType = "review"
	PointEventLike        PointEventType = "like"
	PointEventBeLiked     PointEventType = "be_liked"
	PointEventAdminChange PointEventType = "admin_change"
	PointEventInit        PointEventType = "init"
	PointEventTransfer    PointEventType = "transfer"
	PointEventReward      PointEventType = "reward"
	PointEventPunish      PointEventType = "punish"
	PointEventWithdraw    PointEventType = "withdraw"
	PointEventConsume     PointEventType = "consume"
	PointEventRedeem      PointEventType = "redeem" // 兑换积分
)

// 用户积分明细
type UserPointDetailItem struct {
	Time        string `json:"time"`
	Value       int64  `json:"value"` // 积分变化值: +1, -3
	Description string `json:"description"`
}
type UserPointDetailFilter struct {
	Page              int64
	PageSize          int64
	UserPointDetailID int64
	UserID            int64
	StartTime         string
	EndTime           string
}
